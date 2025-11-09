# Pink ElevenLabs

Text-to-speech and voice transformation using ElevenLabs API.

## Requirements

- macOS 12+
- Python 3.12
- uv
- ElevenLabs API key

## Install

```bash
# Install uv if needed
curl -LsSf https://astral.sh/uv/install.sh | sh

# Clone repository
git clone https://github.com/pinkhairedboy/pink-elevenlabs.git
cd pink-elevenlabs

# Install
uv tool install -e . --python python3.12

# Configure
cp .env.example .env
# Edit .env and add your API key and voice IDs
```

## Configuration

Create `.env` file in project directory:

```bash
ELEVENLABS_API_KEY=your_api_key_here
ELEVENLABS_TTS_VOICE_ID=your_tts_voice_id_here
ELEVENLABS_VOICE_CHANGE_ID=your_voice_change_id_here
```

**Voice IDs:**
- `ELEVENLABS_TTS_VOICE_ID` - Default voice for text-to-speech
- `ELEVENLABS_VOICE_CHANGE_ID` - Default voice for voice transformation
```

## Usage

### Text-to-Speech

```bash
# Basic usage (uses ELEVENLABS_TTS_VOICE_ID from .env)
pink-elevenlabs tts "Your text here"

# Custom output path
pink-elevenlabs tts "Your text here" -o /path/to/output.ogg

# Custom voice
pink-elevenlabs tts "Your text here" --voice VOICE_ID

# Full example with all options
pink-elevenlabs tts "Your text here" \
  --voice VOICE_ID \
  --output /path/to/output.ogg \
  --stability 0.0 \
  --similarity-boost 0.75 \
  --style 0.5 \
  --speed 1.0

# Using emotion tags (Eleven v3 model)
pink-elevenlabs tts "[excited] Hello! [casual] How are you?"

# Disable speaker boost
pink-elevenlabs tts "Your text" --no-speaker-boost
```

**Output format:** Opus 48kHz @ 96kbps (.ogg)

**Default values:**
- Output: `/tmp/speech.ogg`
- Voice: from `ELEVENLABS_TTS_VOICE_ID` env (required)
- Stability: `0.0` (Creative)
- Similarity boost: `0.75`
- Style: `0.5`
- Speed: `1.0`
- Speaker boost: enabled

**Parameters:**
- `-o, --output` - Output file path
- `-v, --voice` - Voice ID (overrides env)
- `--stability` - Voice stability (0.0=Creative, 0.5=Natural, 1.0=Robust)
- `--similarity-boost` - Similarity boost (0.0-1.0)
- `--style` - Style exaggeration (0.0-1.0)
- `--speed` - Speech speed (0.7-1.2)
- `--no-speaker-boost` - Disable speaker boost

### Voice Changer

```bash
# Basic usage (uses ELEVENLABS_VOICE_CHANGE_ID from .env)
pink-elevenlabs voice input.ogg

# Custom output path
pink-elevenlabs voice input.ogg -o /path/to/output.ogg

# Custom target voice
pink-elevenlabs voice input.ogg --voice VOICE_ID

# Full example
pink-elevenlabs voice input.ogg \
  --voice VOICE_ID \
  --output /path/to/output.ogg
```

**Output format:** Opus 48kHz @ 96kbps (.ogg)

**Default values:**
- Output: `/tmp/voice_changed.ogg`
- Voice: from `ELEVENLABS_VOICE_CHANGE_ID` env (required)

**Parameters:**
- `-o, --output` - Output file path
- `-v, --voice` - Target voice ID (overrides env)

## Architecture

```
src/pink_elevenlabs/
├── __init__.py              # Version
├── config.py                # Environment variables and constants
├── cli/
│   ├── __init__.py
│   └── main.py              # CLI entry point with subcommands
└── core/
    ├── __init__.py
    └── elevenlabs.py        # ElevenLabs API client and logic
```

## Models

- **Text-to-Speech**: `eleven_v3` (supports emotion tags, 70+ languages)
- **Voice Changer**: `eleven_multilingual_sts_v2` (speech-to-speech)
