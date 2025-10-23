

---

## 🏥 Project Name: **Swasth-AI**

### *(Voice-First Multilingual Health Assistant for Rural India)*

---

## 🌍 Problem

Healthcare access in rural India faces three main barriers:

* **Connectivity Gaps:** Many regions lack reliable internet access.
* **Language Barriers:** Most healthcare resources exist only in English or Hindi.
* **Limited Medical Guidance:** People often ignore early symptoms or don’t know what to do in emergencies like snake bites or burns.

---

## 💡 Solution

**Arogya Sahayak** is a **voice-first multilingual AI healthcare assistant** that works **offline and online** — bringing medical guidance to every rural user’s phone.

* **Offline:** Voice-based interaction, local language understanding, emergency video guides.
* **Online:** Vision-based AI for analyzing reports and images, advanced medical reasoning, and cloud-assisted symptom checking.

The app acts as a **digital frontline health helper**, not replacing doctors but helping users understand when professional consultation is necessary.

---

## ⚙️ Core User Flow

1. **Login / Setup:**
   Simple OTP-based login or offline local user profile creation.

2. **Voice Interaction (Offline or Online):**

   * User taps the mic button and speaks in their **local language** (Hindi, Tamil, Marathi, Bengali, etc.).
   * Offline speech model converts the speech to text.
   * AI asks a few follow-up questions to narrow down possible causes.
   * Provides a **spoken and visual response** — simple, actionable advice.

3. **Vision-Based Medical Analysis (Online):**

   * 🩸 Upload **blood report** → AI reads values (OCR) and highlights potential issues.
   * 🩻 Upload **X-ray image** → Cloud model detects abnormalities.
   * 🔥 Upload **skin/burn images** → Vision model gives infection or burn severity rating and first-aid advice.

4. **Offline Health Education Library:**

   * Pre-downloaded **video guides** available in local languages.
   * Topics:

     * What to do in a **snake bite** 🐍
     * How to perform **CPR** ❤️‍🔥
     * **Stop bleeding**, **handle burns**, or **treat fever** at home
   * Helps users handle emergencies safely even without internet access.

5. **Online Doctor Connect (Optional):**

   * If AI detects serious health conditions, it suggests or connects the user to nearby clinics or telemedicine options.

---

## 🧠 AI Model Architecture

### **Offline Models**

| Task                   | Model                          | Description                                                   |
| ---------------------- | ------------------------------ | ------------------------------------------------------------- |
| Speech-to-Text         | **Whisper.cpp / Vosk**         | Lightweight offline transcription of voice in local languages |
| Text-to-Speech         | **Coqui TTS / gTTS**           | Local voice output                                            |
| Language Understanding | **IndicBERT / Bhashini**       | Offline multilingual question understanding                   |
| Basic Health Q&A       | **Tiny local LLM** (quantized) | Runs simple health triage on-device                           |
| Video Library          | **Preloaded Media**            | Plays local health & emergency videos                         |

### **Online Models**

| Task                  | Model                            | Description                                      |
| --------------------- | -------------------------------- | ------------------------------------------------ |
| Symptom Understanding | **MedPaLM / ClinicalBERT**       | Contextual analysis of user’s responses          |
| Blood Report Reading  | **Tesseract OCR + ClinicalBERT** | Extracts and interprets key health indicators    |
| X-Ray Analysis        | **CheXNet / DenseNet121**        | Detects abnormalities like pneumonia or TB       |
| Skin/Burn Detection   | **MobileNet / EfficientNet**     | Classifies burn severity or infection patterns   |
| Doctor Connect        | **FastAPI Backend + Database**   | Sends referrals and connects to verified doctors |

---

## 🧩 Tech Stack Overview

* **Frontend:** Flutter (cross-platform, local storage, multilingual UI)
* **Backend:** FastAPI / Flask for API endpoints
* **Database:**

  * Local: SQLite (for offline logs & video library)
  * Cloud: Firebase / Supabase (for user data & reports)
* **Model Deployment:**

  * Offline: `.tflite` or `.gguf` (quantized, efficient models)
  * Online: Hugging Face / Custom cloud inference

---

## 🌐 Offline vs Online Mode Comparison

| Mode             | Key Features                                                          | Models Used                      |
| ---------------- | --------------------------------------------------------------------- | -------------------------------- |
| **Offline Mode** | Voice-based health Q&A, emergency video library, multilingual support | Whisper.cpp, IndicBERT, TinyLLM  |
| **Online Mode**  | Vision-based diagnostics, report analysis, doctor connect             | MedPaLM, CheXNet, MobileNet, OCR |

---

## 🧩 Additional Features

* **User History & Logs:** Store offline chat and report summaries locally.
* **Voice Feedback Loop:** Speak back answers for low-literacy users.
* **Multilingual Interface:** Switchable between major Indian languages.
* **Health Tips Section:** Daily offline tips for hygiene, nutrition, and first aid.

---

## 📊 Real-World Impact

* **Empowers non-literate users** through voice-based interaction.
* **Bridges connectivity gaps** by working both online and offline.
* **Promotes early medical action**, reducing preventable emergencies.
* **Spreads verified healthcare awareness** through localized videos.
