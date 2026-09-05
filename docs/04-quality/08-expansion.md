# 04-QUALITY: 08 - Rule Quality Standards, Anti-Sycophancy & Non-Regression Invariants

> **Kode Dokumen:** `QUAL-08-EXPANSION`
> **Tahapan:** Fase 8 - Repetitive Pattern Flow Guide & Rule Authoring Template (Core Assessment)
> **Status:** Ready for Review
> **Standar Rujukan:** OpenSSF Best Practices & Anti-Sycophancy Static Analyzer Guidelines

Dokumen ini mendefinisikan batasan kualitas ketat, penegakan integritas fungsi murni, pelarangan jalan pintas rahasia (*anti-sycophancy / zero-bypass invariant*), serta jaminan nol-regresi (*non-regression invariant*) saat menambahkan rule baru.

---

## 1. Invarian Integritas Evaluasi (Anti-Sycophancy Invariant)

Dalam pengembangan linter statis, pengembang sering kali tergoda membuat "jalan pintas rahasia" (*hardcoded secret whitelists*) agar kode tertentu tidak memicu error tanpa memperbaiki akar masalah. Pada Charites:

1. **Zero Secret Bypass Invariant:**
   **DILARANG KERAS** menyisipkan pengecualian nama berkas tersembunyi di dalam kode rule (contoh dosa linter legacy: `if strings.Contains(file, "Vendor")` atau `if file == "OpenApiDocs.tsx"`).
2. **Pengecualian Transparan:**
   Seluruh bentuk pengecualian berkas atau direktori **WAJIB** dikonfigurasikan secara transparan melalui:
   - Pola pengabaian `.charitesignore` resmi, ATAU
   - Blok `ignore:` pada berkas konfigurasi `charites.yaml`, ATAU
   - Komentar direktif inline ignore: `// charites:ignore <rule-id>`.
3. **Pesan Diagnostik yang Jujur & Konstruktif:**
   Setiap rule wajib menyertakan alasan teknis yang jujur pada field `Message` dan solusi konkret pada field `Hint`. Rule dilarang menghasilkan pesan ambigu tanpa petunjuk perbaikan.

---

## 2. Invarian Performa & Alokasi Memori Rule

Penambahan puluhan rule baru tidak boleh mendegradasi kecepatan pemindaian Charites:

1. **Fast-Path Character Check:**
   Sebelum menjalankan perulangan atau operasi pemecahan string, rule wajib melakukan pengecekan awal cepat (*quick filter*). Contoh pada `theme.hardcode-color`: jika string tidak memuat karakter `#` atau `rgb`, lewati token tersebut secara instan.
2. **Zero Allocation pada Kasus Bersih:**
   Jika sebuah node tidak melanggar aturan, pemanggilan `Evaluate(node)` wajib mengembalikan `nil` tanpa alokasi heap (`0 B/op`).
3. **Batas Kompleksitas Siklomatik:**
   Fungsi `Evaluate()` wajib memiliki *Cyclomatic Complexity* $\le 10$. Jika logika pemeriksaan terlalu kompleks, pecah menjadi fungsi-fungsi pembantu (*helper functions*) yang teruji secara terpisah.

---

## 3. Invarian Nol Regresi (Non-Regression Invariant)

1. **Integritas Registry:**
   Pendaftaran rule baru tidak boleh menimpa atau merusak urutan pendaftaran rule yang sudah ada di dalam `Registry`.
2. **Throughput Preservation:**
   Penambahan 1 rule baru dilarang meningkatkan waktu pemindaian total repositori lebih dari **3%**.
3. **Core Isolation:**
   File rule baru dilarang mengimpor atau memanipulasi komponen di luar `internal/ir` dan `internal/rules`.

---

## 4. Ambang Batas Kualitas Kode untuk Rule Baru

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Branch Coverage Logika Rule** | $100\%$ | Menguji seluruh cabang `if/else` pada file `<rule>.go` |
| **Tri-Corpus Gate Test** | $100\%$ Lulus | `RuleCorrectnessMetric == Pass` |
| **Kompleksitas Siklomatik** | $\le 10$ per fungsi | `gocyclo -over 10 internal/rules/<domain>/<rule>.go` |
| **Alokasi Heap Node Bersih** | `0 B/op` | `go test -bench=BenchmarkEvaluate_<Rule> -benchmem` |
| **Linter Compliance** | $0$ issues | `golangci-lint run internal/rules/<domain>/...` |
