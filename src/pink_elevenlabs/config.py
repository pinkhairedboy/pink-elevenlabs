"""
Pink ElevenLabs Configuration
All constants, environment variables, and defaults.
"""

from __future__ import annotations

import os
import sys
from pathlib import Path
from dotenv import load_dotenv


def _load_env() -> None:
    """Load environment variables from .env file."""
    project_root = Path(__file__).resolve().parent.parent.parent
    load_dotenv(project_root / ".env")


def get_api_key() -> str:
    """
    Get ElevenLabs API key from environment.

    Returns:
        API key string

    Exits:
        If ELEVENLABS_API_KEY not found
    """
    _load_env()
    api_key = os.getenv("ELEVENLABS_API_KEY")
    if not api_key:
        project_root = Path(__file__).resolve().parent.parent.parent
        print("ERROR: ELEVENLABS_API_KEY not found in environment", file=sys.stderr)
        print(f"Create .env file at: {project_root}/.env", file=sys.stderr)
        sys.exit(1)
    return api_key


def get_tts_voice_id() -> str:
    """
    Get ElevenLabs TTS voice ID from environment.

    Returns:
        Voice ID string for text-to-speech

    Exits:
        If ELEVENLABS_TTS_VOICE_ID not found
    """
    _load_env()
    voice_id = os.getenv("ELEVENLABS_TTS_VOICE_ID")
    if not voice_id:
        project_root = Path(__file__).resolve().parent.parent.parent
        print("ERROR: ELEVENLABS_TTS_VOICE_ID not found in environment", file=sys.stderr)
        print(f"Add ELEVENLABS_TTS_VOICE_ID to .env file at: {project_root}/.env", file=sys.stderr)
        sys.exit(1)
    return voice_id


def get_voice_change_voice_id() -> str:
    """
    Get ElevenLabs voice change voice ID from environment.

    Returns:
        Voice ID string for voice transformation

    Exits:
        If ELEVENLABS_VOICE_CHANGE_ID not found
    """
    _load_env()
    voice_id = os.getenv("ELEVENLABS_VOICE_CHANGE_ID")
    if not voice_id:
        project_root = Path(__file__).resolve().parent.parent.parent
        print("ERROR: ELEVENLABS_VOICE_CHANGE_ID not found in environment", file=sys.stderr)
        print(f"Add ELEVENLABS_VOICE_CHANGE_ID to .env file at: {project_root}/.env", file=sys.stderr)
        sys.exit(1)
    return voice_id


# TTS defaults
DEFAULT_TTS_OUTPUT = "/tmp/speech.ogg"
DEFAULT_TTS_MODEL = "eleven_v3"
DEFAULT_TTS_FORMAT = "opus_48000_96"
DEFAULT_STABILITY = 0.0
DEFAULT_SIMILARITY_BOOST = 0.75
DEFAULT_STYLE = 0.5
DEFAULT_SPEED = 1.0
DEFAULT_SPEAKER_BOOST = True

# Voice change defaults
DEFAULT_VOICE_OUTPUT = "/tmp/voice_changed.ogg"
DEFAULT_VOICE_MODEL = "eleven_multilingual_sts_v2"
DEFAULT_VOICE_FORMAT = "opus_48000_96"
