# 04-QUALITY: 06 - Stability Gate, Zero-Regression & Integrity Invariants

> **Kode Dokumen:** `QUAL-06-GOLDEN`
> **Tahapan:** Fase 6 - Validasi Penuh & Golden Snapshots (Milestone Selesai Pipa)
> **Peran Pilar:** QUALITY = QUALITY THRESHOLD (Ambang Batas Kualitas, Invarian Nol-Regresi & Keamanan)
> **Status:** Ready for Review
> **Standar Rujukan:** Continuous Verification Standards & Zero-Regression Principles

Dokumen ini mendefinisikan batasan kualitas akhir (*quality gates*), penegakan invarian nol-regresi (*zero-regression invariant*), integritas keamanan pemindaian direktori, serta tata kelola pembekuan arsitektur inti (*core architecture freeze*).

---

## 1. Invarian Nol Regresi (Zero-Regression Invariant)

1. **Zero Snapshot Drift:**
   Perubahan sekecil apa pun pada struktur keluaran diagnosis-baik nomor baris, kolom, severity, teks pesan, maupun hint-**DILARANG KERAS** lolos tanpa pembaruan golden snapshot yang disetujui secara eksplisit oleh tim inti melalui peninjauan *git diff*.
2. **Zero-Noise Invariant pada Korpus Bersih:**
   Proyek frontend yang mematuhi panduan token dan didefinisikan pada korpus `tests/fixtures/projects/clean/` **WAJIB** menghasilkan tepat **0 diagnostic** (exit code 0).
3. **Determinis Lintas Platform (Cross-Platform Parity):**
   Hasil pemindaian dan pelaporan golden snapshot wajib menghasilkan byte yang identik saat dijalankan di Linux, macOS, maupun Windows melalui penegakan pemisah baris LF murni (`\n`) dan path POSIX forward slash (`/`).

---

## 2. Invarian Keamanan File Traversal & Ketahanan Fuzzing

1. **Pencegahan Path Traversal Escape:**
   Walker dilarang mengikuti tautan simbolik direktori (`DO NOT FOLLOW`) dan wajib membatasi evaluasi berkas di dalam root workspace pemindaian.
2. **Crash-Resilience Invariant:**
   Fuzzing suite wajib berjalan terus-menerus selama minimal **60 detik per modul** tanpa memicu panic runtime (*unhandled panic*), segmentation fault, atau terminasi abnormal pada proses Go.

---

## 3. Tata Kelola Pembekuan Arsitektur Inti (Core Architecture Freeze)

Fase 6 merupakan gerbang penentu kestabilan arsitektur:
- **Architecture Freeze $\neq$ Bug Fix Freeze:**
  - Batas antarmuka `internal/ir`, `internal/parser`, `internal/scanner`, `internal/analyzer`, dan `internal/reporter` dibekukan dari perombakan arsitektural mayor.
  - Perbaikan cacat kode (*bug fixes*), penajaman deteksi, dan penguatan keamanan yang mematuhi kontrak antarmuka tetap diperbolehkan.
- **Pluggable Modular Expansion:**
  - Seluruh penambahan rule audit di Fase 8 murni mengimplementasikan interface `Rule` di dalam folder mandiri tanpa menyentuh satu baris pun logika traversal engine inti.

---

## 4. Ambang Batas Kualitas & Metrik Kelulusan Gerbang Fase 6

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Total Test Coverage Proyek** | $\ge 85\%$ | `go test -cover ./...` |
| **Kelulusan Golden Snapshots** | $100\%$ lulus | `go test -v ./tests -run TestPipeline_GoldenSnapshots` |
| **Ketahanan Fuzzing** | $0$ crash dalam $\ge 60\text{s}$ per modul | `go test -fuzz=. -fuzztime=60s ./tests/fuzz/...` |
| **Verifikasi Data Race** | $0$ race condition | `go test -race ./...` |
| **Linter Compliance** | $0$ issues | `golangci-lint run ./...` |
