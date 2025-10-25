# VOICE CHAT WEBSOCKET FLOW 

Frontend (React)
┌──────────────────────────────┐
│ 1. Record audio chunks        │
│ 2. Send binary frames         │
│ 3. On pause: send JSON event  │
└─────────────┬────────────────┘
              │ WebSocket (WSS)
              ▼
Backend (Go)
┌──────────────────────────────┐
│ 1. Receives audio chunks     │
│ 2. Forwards to AI/ML WS      │
│ 3. Receives streamed AI text │
│ 4. Receives streamed TTS     │
│ 5. Relays text/audio to FE   │
└─────────────┬────────────────┘
              │ WebSocket (internal)
              ▼
AI/ML Service (Python)
┌─────────────────────────────          ─┐
│ 1. Receive audio chunks                │
│ 2. STT engine → partial/final text     │
│ 3. LLM engine → generate text response │
│ 4. TTS engine → audio chunks           │
│ 5. Stream text/audio back              │
└──────────────────────────────          ┘


---

## 🎙️ VOICE CHAT API DOCUMENTATION

### Base URLs

| Environment | URL                                    |
| ----------- | -------------------------------------- |
| Development | `http://localhost:8080/api/v1`         |
| Production  | `https://api.arogyasahayak.com/api/v1` |

---

## 1. 🎧 Voice Chat Session Lifecycle

| Step | Description                                                                         |
| ---- | ----------------------------------------------------------------------------------- |
| 1    | Frontend requests to **start a voice session**                                      |
| 2    | Backend issues a **session ID** and establishes a **WebSocket** for streaming       |
| 3    | Frontend records audio chunks (PCM or Opus) and sends via WebSocket to backend      |
| 4    | Backend forwards chunks over another WebSocket to **AI/ML service**                 |
| 5    | AI/ML service streams **partial transcriptions (STT)** → forwards to AI model       |
| 6    | AI/ML service streams **AI text responses** → converts them to **TTS audio chunks** |
| 7    | Backend relays the streamed **audio output chunks** to frontend                     |
| 8    | Frontend plays them in real-time (like ChatGPT Voice)                               |

---

## 2. 🚀 API ENDPOINTS — BACKEND (Go)

### `POST /voice/session/start`

Start a new voice chat session.

**Request**

```json
{
  "language": "hi",
  "model": "mistral-7b",
  "session_type": "voice"
}
```

**Response (200)**

```json
{
  "session_id": "vsn_8df91e",
  "ws_url": "wss://api.arogyasahayak.com/api/v1/voice/session/vsn_8df91e/ws"
}
```

---

### `GET /voice/session/:id/ws` (WebSocket)

Bi-directional streaming endpoint for audio and AI responses.

#### 🔄 WebSocket Message Types

**Client → Server**

| Type           | Description                                | Example                                                       |
| -------------- | ------------------------------------------ | ------------------------------------------------------------- |
| `audio_chunk`  | Binary PCM/Opus data                       | (binary data)                                                 |
| `end_of_input` | Marks user finished speaking               | `{"type": "end_of_input"}`                                    |
| `text_message` | Optional: Send text query instead of voice | `{"type": "text_message", "content": "Show my blood report"}` |

**Server → Client**

| Type                 | Description                     | Example                                                                     |
| -------------------- | ------------------------------- | --------------------------------------------------------------------------- |
| `partial_transcript` | STT partial result              | `{"type": "partial_transcript", "text": "show me my"}`                      |
| `final_transcript`   | Full user query after pause     | `{"type": "final_transcript", "text": "show me my blood report"}`           |
| `ai_text`            | AI model streamed text response | `{"type": "ai_text", "text": "Here’s what your blood report indicates..."}` |
| `ai_audio`           | AI response audio chunks (TTS)  | (binary audio data)                                                         |
| `end_of_response`    | Marks end of AI response        | `{"type": "end_of_response"}`                                               |

---

### `POST /voice/session/end`

End the current voice session manually.

**Request**

```json
{ "session_id": "vsn_8df91e" }
```

**Response**

```json
{ "message": "Session ended successfully" }
```

---

## 3. 🧠 API ENDPOINTS — AI/ML SERVICE (Python)

> All streaming handled via WebSocket between Backend ↔ AI/ML Service.

### `POST /internal/ai/session/start`

Called by Go backend to initialize AI pipeline (STT → AI → TTS).

**Request**

```json
{
  "session_id": "vsn_8df91e",
  "language": "hi",
  "model": "mistral-7b"
}
```

**Response**

```json
{
  "session_id": "vsn_8df91e",
  "ws_url": "ws://ai-service:8000/session/vsn_8df91e/ws"
}
```

---

### `GET /session/:id/ws` (WebSocket — Internal)

Handles real-time AI conversation.

#### Backend → AI/ML

| Type           | Description                                     |
| -------------- | ----------------------------------------------- |
| `audio_chunk`  | Raw audio stream from user                      |
| `end_of_input` | Pause detected → begin STT → AI inference → TTS |
| `text_message` | Text query instead of voice                     |

#### AI/ML → Backend

| Type                 | Description             |
| -------------------- | ----------------------- |
| `partial_transcript` | Realtime STT            |
| `final_transcript`   | Final STT text          |
| `ai_text`            | AI streamed text output |
| `ai_audio`           | TTS audio chunks        |
| `end_of_response`    | Marks end               |

---

## 4. 🎤 AI/ML INTERNAL PIPELINE (Python service)

**Pipeline** inside AI service for every `end_of_input` event:

```
AUDIO CHUNKS
   ↓
STT Model (e.g., Whisper small)
   ↓
Text → AI Model (Mistral, Llama, etc.)
   ↓
TTS Model (e.g., VITS / XTTS / Bark / Coqui)
   ↓
Stream audio chunks → Backend
```

---

## 5. 🧩 PROTOCOL SUMMARY

| Layer                   | Protocol  | Direction | Purpose           |
| ----------------------- | --------- | --------- | ----------------- |
| Frontend ↔ Backend      | WebSocket | Duplex    | Audio in, TTS out |
| Backend ↔ AI/ML Service | WebSocket | Duplex    | Stream STT + TTS  |
| Backend ↔ Frontend      | HTTP      | Control   | Start/end session |
| Backend ↔ Auth DB       | HTTP/SQL  | -         | JWT, user info    |

---

## 6. 🎙️ AUDIO SPECIFICATIONS

| Type   | Format                 | Sample Rate | Encoding |
| ------ | ---------------------- | ----------- | -------- |
| Input  | 16-bit PCM / Opus      | 16kHz       | mono     |
| Output | 16-bit PCM / MP3 / OGG | 16kHz       | mono     |

---

