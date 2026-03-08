# pink-elevenlabs

Text-to-speech and voice transformation via ElevenLabs API.

    pink-elevenlabs tts "text"
    pink-elevenlabs tts "text" -o out.ogg --voice VOICE_ID
    pink-elevenlabs voice input.ogg
    pink-elevenlabs voice input.ogg -o out.ogg --voice VOICE_ID
    pink-elevenlabs --health     Check API key
    pink-elevenlabs --version    Show version

Output: Opus 48kHz. Default paths: /tmp/speech.ogg, /tmp/voice_changed.ogg
