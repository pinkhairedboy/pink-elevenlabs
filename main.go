package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pink-tools/pink-core"
	"github.com/pink-tools/pink-otel"
)

var version = "dev"

const (
	serviceName = "pink-elevenlabs"
	apiBaseURL  = "https://api.elevenlabs.io/v1"

	defaultTTSModel   = "eleven_v3"
	defaultVoiceModel = "eleven_multilingual_sts_v2"

	defaultStability       = 0.0
	defaultSimilarityBoost = 0.75
	defaultStyle           = 0.5
	defaultSpeed           = 1.0
)

var outputFormats = map[string]string{
	"opus": "opus_48000_96",
	"mp3":  "mp3_44100_128",
	"pcm":  "pcm_44100",
}

func getDefaultTTSOutput() string {
	return filepath.Join(os.TempDir(), "speech.ogg")
}

func getDefaultVoiceOutput() string {
	return filepath.Join(os.TempDir(), "voice_changed.ogg")
}

func getAPIKey() string {
	key := os.Getenv("ELEVENLABS_API_KEY")
	if key == "" {
		otel.Error(context.Background(), "ELEVENLABS_API_KEY not found")
		fmt.Fprintln(os.Stderr, "ERROR: ELEVENLABS_API_KEY not found in environment")
		os.Exit(1)
	}
	return key
}

func getTTSVoiceID() string {
	id := os.Getenv("ELEVENLABS_TTS_VOICE_ID")
	if id == "" {
		otel.Error(context.Background(), "ELEVENLABS_TTS_VOICE_ID not found")
		fmt.Fprintln(os.Stderr, "ERROR: ELEVENLABS_TTS_VOICE_ID not found in environment")
		os.Exit(1)
	}
	return id
}

func getVoiceChangeID() string {
	id := os.Getenv("ELEVENLABS_VOICE_CHANGE_ID")
	if id == "" {
		otel.Error(context.Background(), "ELEVENLABS_VOICE_CHANGE_ID not found")
		fmt.Fprintln(os.Stderr, "ERROR: ELEVENLABS_VOICE_CHANGE_ID not found in environment")
		os.Exit(1)
	}
	return id
}

func checkHealth() bool {
	key := os.Getenv("ELEVENLABS_API_KEY")
	if key == "" {
		return false
	}

	req, err := http.NewRequest("GET", apiBaseURL+"/user", nil)
	if err != nil {
		return false
	}
	req.Header.Set("xi-api-key", key)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

type ttsRequest struct {
	Text          string        `json:"text"`
	ModelID       string        `json:"model_id"`
	VoiceSettings voiceSettings `json:"voice_settings"`
}

type voiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style"`
	Speed           float64 `json:"speed"`
	UseSpeakerBoost bool    `json:"use_speaker_boost"`
}

func textToSpeech(text, outputPath, voiceID, format string, stability, similarityBoost, style, speed float64, speakerBoost bool) error {
	apiKey := getAPIKey()

	apiFormat, ok := outputFormats[format]
	if !ok {
		return fmt.Errorf("unsupported format: %s", format)
	}

	reqBody := ttsRequest{
		Text:    text,
		ModelID: defaultTTSModel,
		VoiceSettings: voiceSettings{
			Stability:       stability,
			SimilarityBoost: similarityBoost,
			Style:           style,
			Speed:           speed,
			UseSpeakerBoost: speakerBoost,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/text-to-speech/%s?output_format=%s", apiBaseURL, voiceID, apiFormat)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", apiKey)

	otel.Info(context.Background(), "tts_request", map[string]any{
		"voice_id": voiceID,
		"format":   format,
		"text_len": len(text),
	})

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	otel.Info(context.Background(), "tts_complete", map[string]any{"output": outputPath})
	return nil
}

func voiceChange(inputPath, outputPath, voiceID, format string) error {
	apiKey := getAPIKey()

	apiFormat, ok := outputFormats[format]
	if !ok {
		return fmt.Errorf("unsupported format: %s", format)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("audio", filepath.Base(inputPath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, inputFile)
	if err != nil {
		return fmt.Errorf("failed to copy audio data: %w", err)
	}

	writer.WriteField("model_id", defaultVoiceModel)
	writer.Close()

	url := fmt.Sprintf("%s/speech-to-speech/%s?output_format=%s", apiBaseURL, voiceID, apiFormat)
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("xi-api-key", apiKey)

	otel.Info(context.Background(), "voice_change_request", map[string]any{
		"voice_id": voiceID,
		"format":   format,
		"input":    inputPath,
	})

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	otel.Info(context.Background(), "voice_change_complete", map[string]any{"output": outputPath})
	return nil
}

func main() {
	core.LoadEnv(serviceName)

	core.Run(core.Config{
		Name:    serviceName,
		Version: version,
		Usage: fmt.Sprintf(`pink-elevenlabs v%s - Text-to-speech and voice transformation using ElevenLabs API

Usage:
  pink-elevenlabs tts "text" [options]     Text-to-speech synthesis
  pink-elevenlabs voice <input> [options]  Voice transformation
  pink-elevenlabs --health                 Check API key validity
  pink-elevenlabs --version                Show version

TTS options:
  -o, --output <path>         Output file (default: %s)
  -v, --voice <id>            Voice ID (default: ELEVENLABS_TTS_VOICE_ID env)
  -f, --format <fmt>          Output format: opus, mp3, pcm (default: opus)
  --stability <0.0-1.0>       Voice stability (default: %.1f)
  --similarity-boost <0.0-1.0> Similarity boost (default: %.2f)
  --style <0.0-1.0>           Style exaggeration (default: %.1f)
  --speed <0.7-1.2>           Speech speed (default: %.1f)
  --no-speaker-boost          Disable speaker boost

Voice options:
  -o, --output <path>         Output file (default: %s)
  -v, --voice <id>            Target voice ID (default: ELEVENLABS_VOICE_CHANGE_ID env)
  -f, --format <fmt>          Output format: opus, mp3, pcm (default: opus)
`, version, getDefaultTTSOutput(), defaultStability, defaultSimilarityBoost, defaultStyle, defaultSpeed, getDefaultVoiceOutput()),
		Commands: map[string]core.Command{
			"tts": {
				Desc: "Text-to-speech synthesis",
				Run:  cmdTTS,
			},
			"voice": {
				Desc: "Voice transformation",
				Run:  cmdVoice,
			},
		},
	}, nil) // nil main = CLI-only tool
}

func cmdTTS(args []string) error {
	fs := flag.NewFlagSet("tts", flag.ExitOnError)

	output := fs.String("output", getDefaultTTSOutput(), "Output file path")
	fs.StringVar(output, "o", getDefaultTTSOutput(), "Output file path")

	voice := fs.String("voice", "", "Voice ID")
	fs.StringVar(voice, "v", "", "Voice ID")

	format := fs.String("format", "opus", "Output format (opus, mp3, pcm)")
	fs.StringVar(format, "f", "opus", "Output format")

	stability := fs.Float64("stability", defaultStability, "Voice stability (0.0-1.0)")
	similarityBoost := fs.Float64("similarity-boost", defaultSimilarityBoost, "Similarity boost (0.0-1.0)")
	style := fs.Float64("style", defaultStyle, "Style exaggeration (0.0-1.0)")
	speed := fs.Float64("speed", defaultSpeed, "Speech speed (0.7-1.2)")
	noSpeakerBoost := fs.Bool("no-speaker-boost", false, "Disable speaker boost")

	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("text argument required")
	}

	text := fs.Arg(0)
	voiceID := *voice
	if voiceID == "" {
		voiceID = getTTSVoiceID()
	}

	err := textToSpeech(text, *output, voiceID, *format, *stability, *similarityBoost, *style, *speed, !*noSpeakerBoost)
	if err != nil {
		otel.Error(context.Background(), "tts_failed", map[string]any{"error": err.Error()})
		return err
	}

	fmt.Println(*output)
	return nil
}

func cmdVoice(args []string) error {
	fs := flag.NewFlagSet("voice", flag.ExitOnError)

	output := fs.String("output", getDefaultVoiceOutput(), "Output file path")
	fs.StringVar(output, "o", getDefaultVoiceOutput(), "Output file path")

	voice := fs.String("voice", "", "Target voice ID")
	fs.StringVar(voice, "v", "", "Target voice ID")

	format := fs.String("format", "opus", "Output format (opus, mp3, pcm)")
	fs.StringVar(format, "f", "opus", "Output format")

	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("input file argument required")
	}

	inputPath := fs.Arg(0)
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", inputPath)
	}

	voiceID := *voice
	if voiceID == "" {
		voiceID = getVoiceChangeID()
	}

	err := voiceChange(inputPath, *output, voiceID, *format)
	if err != nil {
		otel.Error(context.Background(), "voice_change_failed", map[string]any{"error": err.Error()})
		return err
	}

	fmt.Println(*output)
	return nil
}
