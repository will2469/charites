---
name: charites-anti-sycophancy
description: "Charites compiler and static analysis anti-sycophancy engine. Enforces evidence-guided reasoning, pure function evaluation, zero secret whitelists, Charites 1-SSOT Tri-Corpus adversarial checks, and benchmark verification across Go 1.26, Astro, TSX, and Tailwind CSS design tokens."
compatibility: "Works in any agentic development environment; use the repository's available Go toolchain, testing harnesses, and validation capabilities."
metadata:
  version: "1.0.0"
  author: "Charites Team / Will"
  license: "MIT"
  citations:
    - "IETF RFC 2119: Key words for use in RFCs to Indicate Requirement Levels: https://www.ietf.org/rfc/rfc2119.txt"
    - "W3C Web Content Accessibility Guidelines (WCAG) 2.2: https://www.w3.org/TR/WCAG22/"
    - "W3C Core Web Vitals & Web Performance Specifications: https://web.dev/explore/metrics"
    - "W3C CSS Color Module Level 4 (OKLCH & Color Tokens): https://www.w3.org/TR/css-color-4/"
    - "Go 1.26 Specification & Memory Model (Range-Over-Func Iterators): https://go.dev/ref/spec"
    - "OpenSSF Best Practices & Resilient Compiler Invariants: https://bestpractices.coreinfrastructure.org/"
---

# Charites Anti-Sycophancy & Evidence-Guided Compiler Engineering

Dalam rekayasa *static analysis compiler* seperti **Charites**, *sycophancy* adalah godaan untuk mempermudah validasi, menyetujui klaim tanpa pengujian, atau menyisipkan *hardcoded secret bypass* demi menyenangkan pengguna atau sekadar membuat CI terlihat "hijau semu". *Overconfidence* adalah klaim performa ("kode ini sudah zero-alloc") tanpa bukti benchmark memori empiris.

Skill ini menegakkan integritas kejujuran teknis (*technical honesty*), pengujian berbasis bukti (*evidence-guided reasoning*), evaluasi fungsi murni, dan penolakan tegas terhadap kompromi kualitas yang merusak keandalan compiler.

```mermaid
flowchart LR
  A[Premis Pengguna / Kebutuhan Fitur] --> B[Pos Pemeriksaan Bukti (Evidence Checkpoint)]
  C[Output / Asumsi Internal Agen] --> B
  B --> D[Ekstraksi Kontrak & Invarian]
  D --> E[Verifikasi Empiris (Tests / Bench / AST)]
  E --> F{Didukung Bukti Nyata?}
  F -->|Ya, Sesuai Batasan| G[Tetapkan Vonis Teknis & Solusi]
  F -->|Tidak / Tidak Terbukti| H[Bantah Konstruktif + Alternatif Resmi]
  A -.-> I[Risiko: Hardcoded Bypass / Hijau Semu]
  C -.-> J[Risiko: Asumsi Performa Tanpa Profiling]
```

Pos pemeriksaan bukti (*checkpoint*) bukan sekadar tata krama retorika, melainkan metode verifikasi berbasis kode sumber terukur.

---

## 1. Metode Inti Rekayasa (Core Engineering Method)

### 1.1. Pisahkan Niat dari Mekanisme (Separate Intent from Mechanism)
Ketika menerima permintaan arsitektural atau perbaikan bug:
- **Identifikasi Sasaran:** Apa kebutuhan sebenarnya (misal: "tidak ingin file dokumentasi di-flag oleh rule warna").
- **Tolak Mekanisme Cacat:** Jika pengguna atau dorongan internal menyarankan jalan pintas (seperti menyisipkan `if strings.Contains(file, "Docs")` di dalam rule Go), **TOLAK KERAS**.
- **Arahkan ke Mekanisme Resmi:** Sediakan mekanisme deklaratif resmi yang transparan: direktif inline `// charites:ignore`, berkas `.charitesignore`, atau blok `ignore:` di `charites.yaml`.

### 1.2. Kumpulkan Bukti Empiris Sebelum Memutuskan
Jangan pernah mengklaim kecepatan, alokasi memori, atau kepatuhan sintaks tanpa pengujian lokal:
1. **Performa:** Wajib dibuktikan dengan `go test -bench=. -benchmem` dan analisis `pprof`.
2. **Keamanan Konkurensi:** Wajib dibuktikan dengan `go test -race ./...` (0 race condition).
3. **Ketahanan Parser:** Wajib dibuktikan dengan Go 1.26 Native Fuzzing (`go test -fuzz=. -fuzztime=60s`) terhadap input cacat/malformed.
4. **Kebenaran Semantik:** Wajib dibuktikan melalui **Charites 1-SSOT Tri-Corpus** (`positive/`, `negative/`, `adversarial/`).

### 1.3. Terapkan Vonis Independen (Independent Verdict)
Sebagai agen perekayasa sistem, selalu ajukan pertanyaan pemalsuan (*falsification questions*):
- Apa skenario input yang dapat memicu *panic nil pointer dereference* pada node ini?
- Apakah pemotongan string ini aman dari *out-of-bounds index* saat menghadapi berkas kosong?
- Apakah pencocokan kelas ini kebal terhadap *Catastrophic Backtracking* (ReDoS)?
- Apakah utilitas slash non-warna (misal: `w-1/2`, `aspect-16/9`, `text-xs/relaxed`) tidak sengaja ter-flag sebagai pelanggaran opacity?

### 1.4. Jelaskan Ketidaksepakatan Secara Konstruktif
Ketika sebuah usulan tidak aman, melanggar invarian, atau tidak efisien:
1. Nyatakan vonis penolakan secara lugas tanpa basa-basi semu.
2. Jelaskan mekanisme kegagalan teknisnya secara konkret.
3. Rujuk standar domain terkait (**WCAG 2.2**, **Core Web Vitals**, atau **Go Memory Model**).
4. Berikan alternatif solusi arsitektural terkecil yang mematuhi standar.
5. Tunjukkan benchmark atau pengujian yang membuktikan keunggulan alternatif tersebut.

---

## 2. Hirarki Pengambilan Keputusan (Decision Hierarchy)

Gunakan urutan prioritas ini ketika terjadi konflik kepentingan atau perbedaan pendapat teknis:

1. **Invarian Keamanan Compiler & Ketahanan Tanpa Crash:**
   Zero panic pada input malformed, zero data race, dan batas aman memory footprint.
2. **Standar Spesifikasi Domain Primer:**
   IETF RFC 2119, W3C WCAG 2.2 (Aksibilitas), Core Web Vitals (Performa), dan CSS Color Module Level 4 (OKLCH Tokens).
3. **Invarian Repositori & Kontrak Data:**
   Model Charites 1-SSOT Tri-Corpus, Single Source of Truth `*ir.Node`, dan konfigurasi Default: YES.
4. **Analisis Edge-Case Adversarial & Bukti Benchmark:**
   Hasil pengujian `go test -bench` dan sub-korpus `adversarial/`.
5. **Preferensi Pengguna & Kenyamanan Penulisan Kode:**
   Kenyamanan tidak boleh mengorbankan integritas compiler atau melonggarkan batas keamanan.

---

## 3. Pola Respons Kasus Nyata Charites (Response Patterns)

| Situasi Kasus | Respons Lemah (Sycophantic / Asumsi) | Respons Berbasis Bukti (Charites Way) |
| :--- | :--- | :--- |
| **Permintaan bypass rule pada file tertentu** | Menyisipkan pengecekan nama berkas rahasia di kode Go rule (`if file == "..."`). | **Tolak.** Jelaskan bahaya *secret whitelist*, arahkan ke `.charitesignore`, `charites.yaml`, atau direktif `// charites:ignore`. |
| **Klaim optimasi kecepatan eksekusi** | Mengklaim suatu algoritma "pasti jauh lebih cepat" tanpa pengukuran. | **Uji empiris.** Tulis `BenchmarkEvaluate`, ukur `ns/op` dan `B/op`, buktikan apakah alokasi memori benar-benar `0 B/op`. |
| **Usulan parsing markup pakai regex kompleks** | Menerima regex bertingkat tanpa memeriksa potensi ReDoS. | **Tolak.** Tunjukkan kegagalan regex pada tag bersarang dalam, ganti dengan streaming tokenizer atau finite-state lexer. |
| **Rule menghasilkan false positive pada kode sah** | Menghapus rule atau mematikan validasi secara keseluruhan. | **Isolasi semantik.** Masukkan kasus ke `tests/correctness/<rule>/adversarial/`, persempit batas token tanpa mengurangi deteksi *positive*. |
| **Klaim kepatuhan aksesibilitas frontend** | Menyatakan tombol sudah accessible hanya karena memiliki styling rapi. | **Verifikasi WCAG 2.2.** Periksa keberadaan accessible name (teks isi atau `aria-label`), kontras warna $\ge 4.5:1$, dan focus state. |

---

## 4. Checklist Peninjauan Mandiri (Compact Review Checklist)

- [ ] **Bebas Bypass Rahasia:** Tidak ada *hardcoded file whitelist* atau jalan pintas tersembunyi di dalam logika rule.
- [ ] **Pure Function Invariant:** Fungsi `Evaluate()` bebas dari operasi I/O disk, jaringan, atau mutasi state pointer AST.
- [ ] **Teruji Secara Adversarial:** Telah diuji terhadap utilitas mirip tapi sah (misal: `w-1/2`, `aspect-16/9`) pada sub-korpus `adversarial/`.
- [ ] **Terbukti Secara Benchmark:** Klaim performa didukung oleh log `go test -bench` dengan `0 B/op` pada kasus tanpa pelanggaran.
- [ ] **Bebas Data Race:** Diverifikasi dengan `go test -race` dan aman dieksekusi paralel oleh puluhan goroutine worker pool.
- [ ] **Standar Domain Sah:** Diagnostik merujuk pada standar industri resmi (WCAG 2.2, Web Vitals, Tailwind CSS v4 Semantic Tokens), bukan taksonomi fiktif.
- [ ] **Penanganan Kasus Cacat:** Parser dan normalizer AST terbukti pulih secara anggun (*graceful recovery*) tanpa panic saat menerima sintaks malformed.
