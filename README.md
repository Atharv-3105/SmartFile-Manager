# 🧠 Smart File Manager (Local Semantic Search Engine)

A **local AI-powered file management system** that enables **semantic search over your files** — find documents based on meaning, not just filenames.

---

## 🚀 Overview

Traditional file systems rely on **exact filename matching** or keyword search, which often fails when:

* File names don’t reflect content
* You forget exact keywords
* Documents contain rich semantic information

This project solves that by building a **real-time ingestion + semantic retrieval pipeline**.

> 🔍 Search files like:
> *“documents about reinforcement learning”*
> instead of remembering filenames.

---

## ⚙️ What It Does

### ✅ Real-time File Monitoring

* Watches a directory for new or modified files
* Uses **fsnotify + polling fallback** for reliability across environments (Docker, Windows)

### ✅ Automatic Content Extraction

* Extracts text from:

  * Images (OCR via EasyOCR)
  * Documents (PDFs, text files)
* Uses a Python microservice

### ✅ Embedding Generation

* Converts extracted text into vector embeddings using:

  * `sentence-transformers (all-MiniLM-L6-v2)`
* CPU-optimized for local machines

### ✅ Persistent Storage

* Stores:

  * File metadata
  * Extracted text
  * Vector embeddings
* Uses **SQLite (3NF schema design)**

### ✅ Semantic Search API

* Accepts natural language queries
* Converts query → embedding
* Computes similarity (Cosine Similarity)
* Returns top-k most relevant files

---

## 🧱 Architecture

```
          ┌──────────────┐
          │  File System │
          └──────┬───────┘
                 │
        (fsnotify + poller)
                 │
         ┌───────▼────────┐
         │   Watcher (Go) │
         └───────┬────────┘
                 │
         ┌───────▼────────┐
         │ Worker Pool     │
         │ (Concurrency)   │
         └───────┬────────┘
                 │ HTTP
         ┌───────▼──────────────┐
         │ Extraction Service   │
         │ (Python + AI Models) │
         └───────┬──────────────┘
                 │
         ┌───────▼────────┐
         │ SQLite Storage │
         │ (Metadata + Vectors)
         └───────┬────────┘
                 │
         ┌───────▼────────┐
         │ Search API (Go)│
         └────────────────┘
```

---

## 🧠 Core Concepts & Logic

### 🔁 Event-driven + Polling Hybrid Watcher

* Uses `fsnotify` for real-time events
* Adds **polling fallback** to handle Docker/Windows filesystem inconsistencies
* Ensures **no file changes are missed**

---

### ⚡ Worker Pool (Concurrency)

* Implements a **bounded worker pool in Go**
* Handles multiple file events concurrently
* Prevents blocking of the watcher thread

---

### 🔄 Transactional DB Writes

Each file ingestion follows:

```
Begin Transaction
  → Upsert File
  → Insert Extraction
  → Insert Embedding
Commit
```

Ensures:

* Consistency
* No partial writes
* Safe concurrent processing

---

### 🧮 Semantic Search (Vector Similarity)

* Query → embedding vector
* Stored embeddings → loaded from DB
* Similarity computed using:

```
Cosine Similarity
```

* Top-K results returned based on ranking

---

### 🧩 Microservice Separation

| Service            | Responsibility              |
| ------------------ | --------------------------- |
| Watcher (Go)       | File monitoring + ingestion |
| Extractor (Python) | OCR + embeddings            |
| API (Go)           | Semantic search             |

---

## 🛠 Tech Stack

### Backend

* **Go (Golang)**

  * Concurrency (goroutines, channels)
  * Worker pool design
  * HTTP server

### AI / ML

* **Python**
* **FastAPI**
* **Sentence-Transformers**
* **EasyOCR**

### Database

* **SQLite**

  * Relational schema (3NF)
  * Vector storage (BLOB)

### Infrastructure

* **Docker + Docker Compose**
* Volume mounts for:

  * `watched/` (input files)
  * `data/` (database)
  * `models/` (AI models)

---

## 📦 Features

* 🔍 Semantic file search (meaning-based)
* ⚡ Real-time ingestion pipeline
* 🧵 Concurrent processing with worker pool
* 🧠 AI-powered embeddings
* 💾 Persistent local storage
* 🐳 Fully containerized setup
* 🔁 Robust file watching (polling fallback)

---

## ▶️ Getting Started

### 1. Clone the repository

```bash
git clone <repo-url>
cd AIFileManager
```

---

### 2. Start services

```bash
docker compose up --build
```

---

### 3. Add files

Place files inside:

```
watched/
```

They will be automatically indexed.

---

### 4. Search

```bash
curl -X POST http://localhost:8080/search \
-H "Content-Type: application/json" \
-d '{"query":"machine learning","top_k":3}'
```

---

## 📁 Project Structure

```
cmd/
  ├── watcher/
  └── api/

internal/
  ├── watcher/
  ├── worker/
  ├── storage/
  ├── search/
  ├── client/
  ├── model/
  └── debounce/

extraction-service/
  ├── app/
  └── Dockerfile

data/
watched/
models/
docker-compose.yml
```

---

## ⚠️ Known Challenges & Solutions

### ❌ File events not detected in Docker

✔ Solved using **polling fallback**

---

### ❌ Extractor not ready at startup

✔ Add retry logic / health checks

---

### ❌ Large Docker image size

✔ Optimized using:

* CPU-only PyTorch
* Slim base images
* External model volume

---

## 🚀 Future Improvements

* ANN indexing (FAISS / HNSW)
* UI for browsing/search
* File preview support
* Incremental re-indexing
* Distributed storage
* Streaming ingestion pipeline

---

## 🎯 Key Learnings

* Building **event-driven systems**
* Handling **cross-platform filesystem issues**
* Designing **concurrent pipelines in Go**
* Integrating **AI into backend systems**
* Managing **multi-service architectures**

---

## 📌 Summary

This project demonstrates how to build a **production-style AI system** combining:

* Backend engineering (Go)
* AI/ML (Python)
* Data systems (SQLite)
* Infrastructure (Docker)

👉 Result: A **local semantic search engine for your files**

---

## 🤝 Acknowledgement

Built as a hands-on project to explore:

* Backend system design
* AI integration
* Real-world data pipelines

---
