#!/usr/bin/env python3
"""
Pink ElevenLabs CLI
Main entry point with subcommands.
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

from pink_elevenlabs import __version__
from pink_elevenlabs.config import (
    get_tts_voice_id,
    get_voice_change_voice_id,
    DEFAULT_TTS_OUTPUT,
    DEFAULT_STABILITY,
    DEFAULT_SIMILARITY_BOOST,
    DEFAULT_STYLE,
    DEFAULT_SPEED,
    DEFAULT_SPEAKER_BOOST,
    DEFAULT_VOICE_OUTPUT,
)
from pink_elevenlabs.core.elevenlabs import text_to_speech, voice_change


def cmd_tts(args: argparse.Namespace) -> None:
    """Handle TTS subcommand."""
    voice_id = args.voice if args.voice else get_tts_voice_id()

    try:
        text_to_speech(
            text=args.text,
            output_path=Path(args.output),
            voice_id=voice_id,
            stability=args.stability,
            similarity_boost=args.similarity_boost,
            style=args.style,
            speed=args.speed,
            speaker_boost=not args.no_speaker_boost,
        )
        print(args.output)
    except Exception as e:
        print(f"ERROR: {e}", file=sys.stderr)
        sys.exit(1)


def cmd_voice(args: argparse.Namespace) -> None:
    """Handle voice change subcommand."""
    input_file = Path(args.input)
    if not input_file.exists():
        print(f"ERROR: Input file not found: {input_file}", file=sys.stderr)
        sys.exit(1)

    voice_id = args.voice if args.voice else get_voice_change_voice_id()

    try:
        voice_change(
            input_path=input_file,
            output_path=Path(args.output),
            voice_id=voice_id,
        )
        print(args.output)
    except Exception as e:
        print(f"ERROR: {e}", file=sys.stderr)
        sys.exit(1)


def main() -> None:
    """CLI entry point."""
    parser = argparse.ArgumentParser(
        prog='pink-elevenlabs',
        description='Text-to-speech and voice transformation using ElevenLabs API',
    )

    parser.add_argument(
        '--version',
        action='version',
        version=f'%(prog)s {__version__}'
    )

    subparsers = parser.add_subparsers(dest='command', required=True)

    # TTS subcommand
    parser_tts = subparsers.add_parser('tts', help='Text-to-speech synthesis')
    parser_tts.add_argument('text', help='Text to synthesize')
    parser_tts.add_argument(
        '--output', '-o',
        default=DEFAULT_TTS_OUTPUT,
        help=f'Output file path (default: {DEFAULT_TTS_OUTPUT})'
    )
    parser_tts.add_argument(
        '--stability',
        type=float,
        default=DEFAULT_STABILITY,
        help=f'Voice stability (0.0=Creative, 0.5=Natural, 1.0=Robust, default: {DEFAULT_STABILITY})'
    )
    parser_tts.add_argument(
        '--similarity-boost',
        type=float,
        default=DEFAULT_SIMILARITY_BOOST,
        help=f'Similarity boost (0.0-1.0, default: {DEFAULT_SIMILARITY_BOOST})'
    )
    parser_tts.add_argument(
        '--style',
        type=float,
        default=DEFAULT_STYLE,
        help=f'Style exaggeration (0.0-1.0, default: {DEFAULT_STYLE})'
    )
    parser_tts.add_argument(
        '--speed',
        type=float,
        default=DEFAULT_SPEED,
        help=f'Speech speed (0.7-1.2, default: {DEFAULT_SPEED})'
    )
    parser_tts.add_argument(
        '--no-speaker-boost',
        action='store_true',
        help='Disable speaker boost'
    )
    parser_tts.add_argument(
        '--voice', '-v',
        default=None,
        help='Voice ID (default: from ELEVENLABS_TTS_VOICE_ID env)'
    )
    parser_tts.set_defaults(func=cmd_tts)

    # Voice change subcommand
    parser_voice = subparsers.add_parser('voice', help='Voice transformation')
    parser_voice.add_argument('input', help='Input audio file')
    parser_voice.add_argument(
        '--output', '-o',
        default=DEFAULT_VOICE_OUTPUT,
        help=f'Output file path (default: {DEFAULT_VOICE_OUTPUT})'
    )
    parser_voice.add_argument(
        '--voice', '-v',
        default=None,
        help='Target voice ID (default: from ELEVENLABS_VOICE_CHANGE_ID env)'
    )
    parser_voice.set_defaults(func=cmd_voice)

    args = parser.parse_args()
    args.func(args)


if __name__ == "__main__":
    main()
