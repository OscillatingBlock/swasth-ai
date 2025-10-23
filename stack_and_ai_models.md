> **Stack Recap**
>
> * **Backend:** Go (for APIs, WebSocket streaming, auth, data)
> * **AI/ML Microservice:** Flask (Python, runs AI models)
> * **Frontend:** React Native or Flutter (mobile, voice UI)
> * **Realtime:** WebSockets (for live voice + response streaming)

---

## 🏗️ Overall Architecture

```
[Frontend App]
   ↓ WebSocket / HTTPS
[Go Backend]  ←→  [Flask AI Microservice]
   ↓
[Database] (SQLite + Firebase)
```

* **Go Backend** → Handles auth, user management, video sync, and WebSocket bridge.
* **Flask Microservice** → Runs AI tasks (speech, NLP, vision, OCR).
* **Frontend (React Native/Flutter)** → Mic input, speech output, local video playback.

---

## 🧠 Feature-by-Feature Breakdown

---

### 1️⃣ **Voice-based Interaction (Offline + Online)**

#### Implementation Flow

1. User taps mic → app records audio.
2. If offline: audio processed locally using embedded STT & TTS models.
3. If online: audio streamed to Flask microservice via WebSocket → transcribed + analyzed.
4. Flask sends back response text/audio stream → streamed via Go backend to frontend.

#### Open Source Models

| Task                     | Model                                                                                     | Description                                         |
| ------------------------ | ----------------------------------------------------------------------------------------- | --------------------------------------------------- |
| Speech-to-Text (Offline) | [**Whisper.cpp**](https://github.com/ggerganov/whisper.cpp)                               | Local inference (Hindi + English supported)         |
| Speech-to-Text (Online)  | [**Whisper-large-v3**](https://huggingface.co/openai/whisper-large-v3)                    | High-accuracy transcription                         |
| Text-to-Speech           | [**Coqui TTS**](https://github.com/coqui-ai/TTS)                                          | Fast multilingual TTS (supports Hindi, Tamil, etc.) |
| Language Understanding   | [**IndicBERT**](https://huggingface.co/ai4bharat/indic-bert)                              | Handles intent + context in Indian languages        |
| Basic Health Q&A         | [**tiny-llama-1.1B-chat.gguf**](https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-GGUF) | On-device quantized model for quick health triage   |

#### Tech Integration

* **Frontend:** Mic input → WebSocket audio streaming
* **Go:** Manages streaming events, fallback to offline mode
* **Flask:** Transcribes, processes NLP, sends back incremental responses

---

### 2️⃣ **Vision-based Medical Analysis (Online Only)**

#### Implementation Flow

1. User uploads image (X-ray, burn, report).
2. Go backend uploads to Flask via REST endpoint `/analyze`.
3. Flask routes to the relevant model:

   * OCR → for text reports
   * CNN → for medical images
4. Flask returns a structured JSON response with insights.

#### Open Source Models

| Feature              | Model                                                                         | Description                                    |
| -------------------- | ----------------------------------------------------------------------------- | ---------------------------------------------- |
| Blood Report OCR     | [**Tesseract OCR**](https://github.com/tesseract-ocr/tesseract)               | Extracts lab readings from PDFs/images         |
| Report Analysis      | [**BioClinicalBERT**](https://huggingface.co/emilyalsentzer/Bio_ClinicalBERT) | Understands medical terms and reference ranges |
| X-Ray Classification | [**CheXNet (DenseNet121)**](https://github.com/arnoweng/CheXNet)              | Detects pneumonia, TB, etc.                    |
| Skin/Burn Detection  | [**EfficientNet-B0**](https://github.com/rwightman/pytorch-image-models)      | Trained on skin disease datasets like HAM10000 |

#### Tech Integration

* **Frontend:** Image picker → upload to Go backend
* **Go:** Forwards file to Flask API `/analyze/xray` or `/analyze/skin`
* **Flask:** Runs inference → returns diagnosis JSON → Go sends to frontend

---

### 3️⃣ **Offline Emergency Video Library**

#### Implementation Flow

1. When app installs (or during sync), Go backend sends list of available health videos.
2. User downloads videos for offline access (snake bite, CPR, bleeding, etc.).
3. Stored locally via React Native/Flutter storage (e.g., `video_player` or `react-native-video`).

#### Data Source

* Curated verified medical emergency videos (in Hindi, Tamil, Bengali, etc.)
* Hosted on IPFS, Firebase Storage, or your own backend CDN.

#### Tech Integration

* **Frontend:** Offline player + category list
* **Go Backend:** Serves video metadata and files
* **Database:** SQLite for local video index

---

### 4️⃣ **Multilingual Support**

#### Implementation

* Frontend UI uses i18n with local `.json` strings (React Native: `react-i18next`, Flutter: `flutter_localizations`).
* Flask models (IndicBERT) handle query understanding.
* TTS models (Coqui) output speech in selected language.

#### Languages Targeted

* Hindi, Marathi, Tamil, Telugu, Bengali, English (extendable with Bhashini API).

---

### 5️⃣ **Doctor Connect / Teleconsult Suggestion (Online)**

#### Implementation Flow

1. AI flags potential serious condition.
2. Go backend fetches nearest available doctors (stored in DB or from external API).
3. Displays “Consult Doctor” button → triggers in-app call or WhatsApp link.

#### Tech Integration

* **Backend:** Go REST endpoint `/doctors/nearby?lat=...&lng=...`
* **Database:** PostGIS / Firebase geolocation
* **Frontend:** Shows doctor list and call option

---

### 6️⃣ **User Authentication & Sync**

#### Implementation

* **Go Backend:**

  * `/auth/otp` (OTP verification using Firebase Auth or custom SMS API)
  * `/user/profile` (sync profile + local logs)
* **Frontend:** Stores user profile locally for offline access.
* **Database:** SQLite (local), Firebase (cloud backup)

---

## ⚙️ Communication Between Components

| Channel                         | Purpose                                               |
| ------------------------------- | ----------------------------------------------------- |
| **WebSocket (Go ↔ Frontend)**   | Real-time voice transcription and streaming responses |
| **REST API (Go ↔ Flask)**       | Vision tasks, OCR, report analysis                    |
| **gRPC (optional, Go ↔ Flask)** | Low-latency structured data exchange                  |
| **Firebase / Supabase**         | Cloud sync for profiles, reports, logs                |

---

## 🧩 Suggested Open Source Libraries

| Layer      | Library / Tool                                                              |
| ---------- | --------------------------------------------------------------------------- |
| Go Backend | `gorilla/websocket`, `gin-gonic/gin`, `go-fiber/fiber`, `gorm`              |
| Flask AI   | `transformers`, `torch`, `tesseract`, `opencv`, `fastapi` (if migrating)    |
| Frontend   | `react-native-voice`, `react-native-video`, `flutter_tts`, `speech_to_text` |
| Database   | `SQLite` (offline), `Firebase` (sync)                                       |

---

## 🧠 Example Data Flow (Voice Q&A)

1. User → speaks → audio stream via WebSocket to Go
2. Go → forwards to Flask → Whisper.cpp transcribes
3. Flask → IndicBERT interprets query → TinyLlama generates response
4. Flask → streams back text chunks → Go → frontend → Coqui TTS speaks out

---

## 🚀 Summary of Model Usage

| Task                         | Offline Model | Online Model               |
| ---------------------------- | ------------- | -------------------------- |
| Speech Recognition           | Whisper.cpp   | Whisper-large-v3           |
| Text-to-Speech               | Coqui TTS     | Coqui TTS (cloud voice)    |
| NLP / Language Understanding | IndicBERT     | MedPaLM or ClinicalBERT    |
| Health Reasoning             | TinyLlama     | MedPaLM                    |
| Vision Analysis              | –             | CheXNet, EfficientNet, OCR |
| Video Library                | Local storage | –                          |
