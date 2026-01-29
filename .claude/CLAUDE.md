# pink-elevenlabs

Text-to-speech and voice transformation via ElevenLabs API.

```bash
pink-elevenlabs tts "text"
pink-elevenlabs tts "text" -o /tmp/out.ogg --voice VOICE_ID
pink-elevenlabs voice input.ogg
pink-elevenlabs voice input.ogg -o output.ogg --voice VOICE_ID
```

Output: `/tmp/speech.ogg`, `/tmp/voice_changed.ogg` (Opus 48kHz)
