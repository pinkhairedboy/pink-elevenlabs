"""
ElevenLabs API client and business logic.
"""

from __future__ import annotations

from pathlib import Path
from elevenlabs.client import ElevenLabs

from pink_elevenlabs.config import (
    get_api_key,
    DEFAULT_TTS_MODEL,
    DEFAULT_TTS_FORMAT,
    DEFAULT_VOICE_MODEL,
    DEFAULT_VOICE_FORMAT,
)


def get_client() -> ElevenLabs:
    """
    Initialize and return ElevenLabs client.

    Returns:
        ElevenLabs client instance
    """
    return ElevenLabs(api_key=get_api_key())


def text_to_speech(
    text: str,
    output_path: Path,
    voice_id: str,
    stability: float,
    similarity_boost: float,
    style: float,
    speed: float,
    speaker_boost: bool,
) -> None:
    """
    Convert text to speech using ElevenLabs API.

    Args:
        text: Text to synthesize
        output_path: Output file path
        voice_id: Voice ID to use
        stability: Voice stability (0.0-1.0)
        similarity_boost: Similarity boost (0.0-1.0)
        style: Style exaggeration (0.0-1.0)
        speed: Speech speed (0.7-1.2)
        speaker_boost: Enable speaker boost
    """
    client = get_client()

    voice_settings = {
        "stability": stability,
        "similarity_boost": similarity_boost,
        "style": style,
        "speed": speed,
        "use_speaker_boost": speaker_boost,
    }

    audio_bytes = client.text_to_speech.convert(
        text=text,
        voice_id=voice_id,
        model_id=DEFAULT_TTS_MODEL,
        output_format=DEFAULT_TTS_FORMAT,
        voice_settings=voice_settings
    )

    output_path.parent.mkdir(parents=True, exist_ok=True)

    with open(output_path, "wb") as f:
        for chunk in audio_bytes:
            f.write(chunk)


def voice_change(
    input_path: Path,
    output_path: Path,
    voice_id: str,
) -> None:
    """
    Transform voice using ElevenLabs Speech-to-Speech API.

    Args:
        input_path: Input audio file path
        output_path: Output file path
        voice_id: Target voice ID
    """
    client = get_client()

    with open(input_path, "rb") as f:
        audio_bytes = client.speech_to_speech.convert(
            voice_id=voice_id,
            audio=f,
            model_id=DEFAULT_VOICE_MODEL,
            output_format=DEFAULT_VOICE_FORMAT
        )

        output_path.parent.mkdir(parents=True, exist_ok=True)

        with open(output_path, "wb") as out:
            for chunk in audio_bytes:
                out.write(chunk)
