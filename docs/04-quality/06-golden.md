# 04-QUALITY: 06 - Stability Gate, Zero-Regression & Integrity Invariants

> **Kode Dokumen:** `QUAL-06-GOLDEN`
> **Tahapan:** Fase 6 - Validasi Penuh & Golden Snapshots (Milestone Selesai Pipa)
> **Status:** Ready for Review
> **Standar Rujukan:** Continuous Verification Standards & Zero-Regression Principles

Dokumen ini mendefinisikan batasan kualitas akhir (*quality gates*), penegakan invarian nol-regresi (*zero-regression invariant*), integritas keamanan pemindaian direktori, serta pembekuan arsitektur inti (*core architecture freeze*).

---

## 1. Invarian Nol Regresi (Zero-Regression Invariant)

1. **Zero Snapshot Drift:**
   Perubahan sekecil apapun pada keluaran pelaporan-baik nomor baris, nomor kolom, teks pesan peringatan, maupun tata letak tabel terminal-**DILARANG KERAS** lolos tanpa pembaruan golden snapshot yang disetujui secara eksplisit oleh tim inti melalui git review.
2. **Zero-Noise Invariant pada Kode Sah:**
   Kode frontend yang valid dan mematuhi panduan desain wajib menghasilkan **0 diagnostic** (exit code 0). Tingkat toleransi terhadap false-positive pada proyek bersih adalah **0%**.
3. **Deterministik Lintas Lingkungan (Cross-Platform Parity):**
   Hasil pemindaian dan pelaporan golden snapshot wajib menghasilkan byte yang identik saat dijalankan di Linux, macOS, maupun Windows (termasuk penanganan karakter line ending `\r\n` vs `\n`).

---

## 2. Invarian Keamanan File Traversal (Security Boundaries)

Untuk memastikan binary aman dieksekusi di server CI/CD bersama atau lingkungan multi-tenant:
1. **Pencegahan Path Traversal Escape:**
   Walker dilarang mengikuti tautan simbolik (*symbolic links*) yang mengarah ke luar root direktori repositori (contoh: symlink yang mengarah ke `/etc/passwd` atau root direktori sistem).
2. **Crash-Resilience Invariant:**
   Jutaan variasi byte acak dari suite fuzzing tidak boleh memicu panic runtime atau terminasi abnormal pada proses Go.

---

## 3. Gerbang Pembekuan Arsitektur Inti (Core Architecture Freeze)

Fase 6 merupakan gerbang penentu kestabilan arsitektur:
- Kontrak data `internal/ir` dibekukan (*frozen*).
- Pipa traversal `internal/analyzer` dibekukan.
- Antarmuka registri `internal/rules/registry.go` dibekukan.
- Seluruh penambahan fungsionalitas di Fase 8 murni berupa *plugin rule modular* yang mengimplementasikan interface `Rule` tanpa menyentuh satu barispun logika inti engine.

---

## 4. Ambang Batas Kualitas & Metrik Kelulusan Gerbang Fase 6

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Total Test Coverage Proyek** | $\ge 85\%$ | `go test -cover ./...` |
| **Kelulusan Golden Snapshots** | $100\%$ lulus | `go test -v ./tests -run TestPipeline_GoldenSnapshots` |
| **Ketahanan Fuzzing** | $0$ crash dalam 60s per target | `go test -fuzz=FuzzAstroPipeline -fuzztime=60s ./tests/fuzz/...` |
| **Verifikasi Data Race** | $0$ race condition | `go test -race ./...` |
| **Linter Compliance** | $0$ issues | `golangci-lint run ./...` |
