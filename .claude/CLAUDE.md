# pink-elevenlabs

Text-to-speech and voice transformation via ElevenLabs API.

```bash
/Users/pink-tools/pink-elevenlabs/pink-elevenlabs tts "text"
/Users/pink-tools/pink-elevenlabs/pink-elevenlabs tts "text" -o /tmp/out.ogg --voice VOICE_ID
/Users/pink-tools/pink-elevenlabs/pink-elevenlabs voice input.ogg
/Users/pink-tools/pink-elevenlabs/pink-elevenlabs voice input.ogg -o output.ogg --voice VOICE_ID
/Users/pink-tools/pink-elevenlabs/pink-elevenlabs --health  # Check API key
```

Output: `/tmp/speech.ogg`, `/tmp/voice_changed.ogg` (Opus 48kHz)
