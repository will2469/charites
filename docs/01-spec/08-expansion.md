# 01-SPEC: 08 - Repetitive Pattern Flow Guide & Rule Authoring Specification

> **Kode Dokumen:** `SPEC-08-EXPANSION`
> **Tahapan:** Fase 8 - Repetitive Pattern Flow Guide & Rule Authoring Template (Core Assessment)
> **Peran Pilar:** SPEC = WHAT (Spesifikasi Kontrak Penulisan Rule & Template Studi Kasus)
> **Status:** Ready for Review
> **Standar Rujukan:** Standardized Static Analysis Rule Specification / Micro-Kernel Extensibility Pattern

Dokumen ini mendefinisikan spesifikasi kontrak penulisan rule baru (*Rule Authoring Contract*), struktur metadata terpadu, penegakan matriks verifikasi semantik, serta studi kasus spesifikasi percontohan (**`theme.hardcode-color`**).

---

## 1. Kontrak Baku Penulisan Rule Baru (Rule Authoring Contract)

Setiap rule baru yang akan dikembangkan untuk Charites **WAJIB** memenuhi 7 elemen kontrak baku berikut sebelum diimplementasikan:

### 1.1. Identitas & Penamaan Kanonikal (Canonical Rule ID)
- Format ID tunggal: `<category>.<slug>` (huruf kecil, dipisahkan titik dan tanda hubung minus).
- Domain kategori yang sah:
  - `theme`: Token desain, warna, spacing, tipografi, `@theme`.
  - `a11y`: Aksesibilitas, ARIA, semantic HTML, kontras.
  - `perf`: Core Web Vitals (LCP, CLS, INP), asset priority, rendering.
  - `layout`: Responsive design, container queries, grid.
  - `seo`: Metadata dokumen, canonical link, social tags.
- Dilarang keras menggunakan penomoran numerik acak (seperti `T01`, `A02`).

### 1.2. Kontrak Cakupan & Batas Deteksi (Detection Boundary Contract)
Setiap rule wajib mendefinisikan tabel batas eksplisit:
- **In-Scope:** Daftar tag target, atribut target, keluarga utilitas, atau pola deklarasi CSS yang dievaluasi.
- **Out-of-Scope:** Daftar pola sintaks serupa yang secara sengaja diabaikan (misal: anchor link `#`, variabel CSS `var(--*)`, pecahan dimensi layout).

### 1.3. Normalisasi Varian Tailwind (Variant Normalization)
Rule yang memeriksa class Tailwind wajib mendukung normalisasi prefix/variant (`hover:`, `dark:`, `focus:`, `md:`, dll.) agar deteksi tidak menghasilkan *false negative* pada class ber-modifier.

### 1.4. Kontrak Payload Diagnostik & Hint Dinamis
Setiap temuan diagnostik wajib memuat:
- Lokasi presisi: `File`, `Line`, `Column` dari `node.Span`.
- Identitas: `RuleID`, `Category`, `Severity`.
- Pesan Edukatif: `Message` menjelaskan pola pelanggaran secara objektif.
- Rekomendasi Remedi Konkret: `Hint` menghasilkan token semantik pengganti yang spesifik secara dinamis.

### 1.5. Kontrak Penekanan Direktif Inline Ignore
Mendukung sintaks penekanan standar:
- JS / TSX: `// charites:ignore <rule-id>`
- Astro HTML: `<!-- charites:ignore <rule-id> -->`
- Cakupan penekanan: berlaku untuk baris yang sama (same-line) dan rentang elemen AST berikutnya (node span).

### 1.6. Matriks Verifikasi Semantik Ekspektasi (Case-by-Case Matrix)
Pengujian rule tidak boleh hanya memeriksa `PositiveViolations > 0`. Rule wajib mendefinisikan matriks kasus ber-ID unik (`POS-001`, `NEG-001`, `ADV-001`) dengan ekspektasi hasil yang diverifikasi 100% cocok (*exact match*).

### 1.7. Korpus Uji Tri-Corpus Terstandarisasi
Menyediakan berkas uji nyata di `tests/correctness/<category>.<slug>/`:
- `positive/`: Berkas yang memuat variasi pelanggaran nyata.
- `negative/`: Berkas sah yang mematuhi desain sistem resmi.
- `adversarial/`: Berkas jebakan false-positive dan inline ignore.

---

## 2. Studi Kasus Spesifikasi Rule: `theme.hardcode-color` (Reference Specification)

Sebagai contoh penerapan kontrak penulisan rule baru:

### 2.1. Identitas & Metadata
- **ID:** `theme.hardcode-color`
- **Kategori:** `theme`
- **Default Severity:** `warn` (`ir.SeverityWarning`)
- **Deskripsi:** Mendeteksi penggunaan kode warna mentah yang wajib diganti dengan token semantik dari `global.css`.

### 2.2. Batas Deteksi (Detection Boundary Matrix)

| Pola Input | Klasifikasi | Ekspektasi Deteksi | Rekomendasi Hint / Alasan |
| :--- | :---: | :---: | :--- |
| `bg-[#2563eb]` | In-Scope |  Violation | Ganti dengan token semantik `bg-primary` |
| `text-[#000000]` | In-Scope |  Violation | Ganti dengan token semantik `text-foreground` |
| `hover:border-[rgb(255,0,0)]` | In-Scope (Variant) |  Violation | Ganti dengan token semantik `border-destructive` |
| `style="color: #2563eb;"` | In-Scope (Inline Style)|  Violation | Ganti dengan token semantik dari `global.css` |
| `bg-primary` | Negative (Clean) |  Pass (0 diag) | Token semantik resmi |
| `text-muted` | Negative (Clean) |  Pass (0 diag) | Token semantik resmi |
| `<a href="#section-hero">` | Out-of-Scope (Bait) |  Pass (0 diag) | Tanda `#` adalah anchor link, bukan warna |
| `<div id="card-1">` | Out-of-Scope (Bait) |  Pass (0 diag) | Nilai ID, bukan deklarasi warna |
| `var(--custom-color)` | Out-of-Scope |  Pass (0 diag) | Variabel CSS (domain rule terpisah) |

---

## 3. Metadata Edukasi untuk MCP & Wiki Generator

Setiap rule wajib menyertakan metadata dokumentasi terpadu:
- **Explanation:** Menjelaskan dampak negatif warna mentah terhadap *dark mode consistency*, *rebranding flexibility*, dan *design token governance*.
- **Bad Example:** Contoh kode nyata yang melanggar.
- **Good Example:** Contoh kode yang telah diperbaiki menggunakan token semantik resmi.
- **Remediation:** Panduan teknis pemetaan nilai warna mentah ke token semantik `global.css`.
