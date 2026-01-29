# pink-elevenlabs

Text-to-speech and voice transformation using ElevenLabs API.

## Install

Download the latest binary from [Releases](https://github.com/pink-tools/pink-elevenlabs/releases).

## Usage

```bash
# Text-to-speech
pink-elevenlabs tts "Hello world"
pink-elevenlabs tts "Text" -o output.ogg --stability 0.5

# Voice transformation
pink-elevenlabs voice input.ogg
pink-elevenlabs voice input.ogg -o output.ogg

# Check API connectivity
pink-elevenlabs --health
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ELEVENLABS_API_KEY` | Yes | ElevenLabs API key |
| `ELEVENLABS_TTS_VOICE_ID` | No | Default voice for TTS |
| `ELEVENLABS_VOICE_CHANGE_ID` | No | Default voice for transformation |

## TTS Options

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output` | /tmp/speech.ogg | Output file |
| `-v, --voice` | env | Voice ID |
| `-f, --format` | opus | Format (opus, mp3, pcm) |
| `--stability` | 0.0 | Voice stability (0.0-1.0) |
| `--similarity-boost` | 0.75 | Similarity boost (0.0-1.0) |
| `--style` | 0.5 | Style exaggeration (0.0-1.0) |
| `--speed` | 1.0 | Speech speed (0.7-1.2) |
| `--no-speaker-boost` | false | Disable speaker boost |

## Voice Options

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output` | /tmp/voice_changed.ogg | Output file |
| `-v, --voice` | env | Target voice ID |
| `-f, --format` | opus | Format (opus, mp3, pcm) |

## Build from Source

```bash
git clone https://github.com/pink-tools/pink-elevenlabs.git
cd pink-elevenlabs
go build .
```
