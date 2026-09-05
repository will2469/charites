# 04-QUALITY: 08 - Rule Quality Standards, Anti-Sycophancy & Non-Regression Invariants

> **Kode Dokumen:** `QUAL-08-EXPANSION`
> **Tahapan:** Fase 8 - Repetitive Pattern Flow Guide & Rule Authoring Template (Core Assessment)
> **Peran Pilar:** QUALITY = QUALITY THRESHOLD (Ambang Batas Kualitas, Anti-Sycophancy & Anggaran Kinerja)
> **Status:** Ready for Review
> **Standar Rujukan:** OpenSSF Best Practices & Anti-Sycophancy Static Analyzer Guidelines

Dokumen ini mendefinisikan batasan kualitas, penegakan integritas fungsi murni, pelarangan jalan pintas rahasia (*anti-sycophancy invariant*), serta pemisahan tegas antara invarian mutlak, target kualitas, dan anggaran performa.

---

## 1. Invarian Mutlak (Hard Non-Negotiable Invariants)

1. **Anti-Sycophancy / Zero Secret Bypass Invariant:**
   **DILARANG KERAS** menyisipkan pengecualian nama berkas tersembunyi di dalam kode rule (contoh larangan: `if strings.Contains(file, "Vendor")` atau `if file == "SpecialPage.tsx"`). Seluruh pengecualian wajib melalui mekanisme transparan: `.charitesignore`, config `charites.yaml`, atau direktif inline `// charites:ignore`.
2. **Pure Function Invariant:**
   Fungsi `rule.Evaluate(*ir.Node)` wajib murni tanpa disk/network I/O, memperlakukan pointer `*ir.Node` sebagai *immutable*, dan bersifat idempoten deterministik.
3. **Concurrency Safety Invariant:**
   Rule singleton wajib aman dievaluasi secara paralel oleh puluhan goroutine tanpa mutable state sharing (`0 data races` di bawah flag `-race`).

---

## 2. Target Kualitas Kode (Quality Targets)

1. **Cakupan Pengujian Pernyataan (Statement Coverage):**
   Setiap berkas rule baru di `internal/rules/<domain>/` wajib memiliki coverage $\ge 90\%$.
2. **Kompleksitas Siklomatik Terkendali:**
   Fungsi `Evaluate()` dan helper-nya wajib memiliki *Cyclomatic Complexity* $\le 10$ (`gocyclo -over 10`).
3. **Pesan Diagnostik Konstruktif:**
   Setiap pelanggaran wajib menyertakan deskripsi masalah objektif di `Message` dan saran perbaikan konkret di `Hint`.

---

## 3. Anggaran Performa & Alokasi Memori (`QUAL-08-PERF-001`)

1. **Fast-Path Check pada Node Bersih:**
   Rule wajib menyediakan pemeriksaan cepat (*quick filter*) di awal fungsi untuk menghindari alokasi heap saat mengevaluasi node legal (Target Desain: `0 B/op` dan `0 allocs/op` pada node bersih).
2. **Throughput Preservation Budget:**
   Penambahan satu rule baru tidak boleh mendegradasi throughput pemindaian repositori lebih dari batas wajar yang terukur pada benchmark pipeline standar.

---

## 4. Ambang Batas Kualitas & Metrik Kelulusan

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran | Klasifikasi |
| :--- | :---: | :--- | :--- |
| **Data Race Safety** | $0$ data race detected | `go test -race ./internal/rules/...` | Hard Invariant |
| **Zero Secret Bypass** | $0$ hardcoded whitelists | Code audit & adversarial corpus | Hard Invariant |
| **Tri-Corpus Matrix Match** | $100\%$ exact match | `go test -run TestTriCorpus` | Hard Invariant |
| **Statement Coverage** | $\ge 90\%$ | `go test -cover ./internal/rules/...` | Quality Target |
| **Kompleksitas Siklomatik** | $\le 10$ per fungsi | `gocyclo -over 10 ./internal/rules` | Quality Target |
| **Linter Compliance** | $0$ issues | `golangci-lint run ./internal/rules/...` | Quality Target |
