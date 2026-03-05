# pink-elevenlabs

Text-to-speech and voice transformation via ElevenLabs API.

```bash
{{PINK_TOOLS}}/pink-elevenlabs/pink-elevenlabs tts "text"
{{PINK_TOOLS}}/pink-elevenlabs/pink-elevenlabs tts "text" -o out.ogg --voice VOICE_ID
{{PINK_TOOLS}}/pink-elevenlabs/pink-elevenlabs voice input.ogg
{{PINK_TOOLS}}/pink-elevenlabs/pink-elevenlabs voice input.ogg -o out.ogg --voice VOICE_ID
{{PINK_TOOLS}}/pink-elevenlabs/pink-elevenlabs --health    # Check API key
{{PINK_TOOLS}}/pink-elevenlabs/pink-elevenlabs --version   # Show version
```

Output: Opus 48kHz (default paths: `/tmp/speech.ogg`, `/tmp/voice_changed.ogg`)
