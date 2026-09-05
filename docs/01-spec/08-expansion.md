# 01-SPEC: 08 - Repetitive Pattern Flow Guide & Rule Authoring Specification

> **Kode Dokumen:** `SPEC-08-EXPANSION`
> **Tahapan:** Fase 8 - Repetitive Pattern Flow Guide & Rule Authoring Template (Core Assessment)
> **Status:** Ready for Review
> **Standar Rujukan:** Standardized Static Analysis Rule Specification / Standardized Extensibility Pattern

Dokumen ini mendefinisikan panduan spesifikasi baku (*standardized specification template*) untuk menambahkan rule audit baru ke dalam Charites setelah fondasi compiler selesai, dengan studi kasus konkret implementasi **Rule #2 (`theme.hardcode-color`)**.

---

## 1. Standar Penamaan & Identitas Rule Baru

Setiap rule baru yang akan dikembangkan wajib mematuhi aturan penamaan Charites Rule ID tunggal:

1. **Format Identifier Tunggal:**
   ```text
   <category>.<rule-slug>
   ```
   - `<category>`: Domain bidang yang sah (`theme`, `a11y`, `perf`, `layout`, `seo`).
   - `<rule-slug>`: Kata deskriptif ringkas menggunakan huruf kecil dan tanda hubung minus (`kebab-case`).
   - Contoh: `theme.hardcode-color`, `a11y.html-missing-lang`, `perf.lcp-priority`.
2. **Larangan Kode Ganda:**
   Dilarang menggunakan kode numerik acak seperti `T01` atau `A02`. Semua referensi konfigurasi, CLI filter, dan inline ignore wajib menggunakan Charites Rule ID tunggal.

---

## 2. Template Spesifikasi Rule Baru (Contoh Studi Kasus: `theme.hardcode-color`)

Setiap penambahan rule baru wajib diawali dengan mendefinisikan spesifikasi kebutuhan dalam format standar berikut:

### 2.1. Identitas & Metadata
- **ID:** `theme.hardcode-color`
- **Kategori:** `theme`
- **Default Severity:** `warn` (`ir.SeverityWarn`)
- **Deskripsi:** Mendeteksi penggunaan kode warna mentah (Hexadecimal `#RGB`/`#RRGGBB`, RGB/RGBA, atau HSL) di dalam inline style atau utility class Tailwind yang wajib diganti dengan token semantik dari `global.css`.

### 2.2. Pola Pelanggaran (Violation Pattern)
Pemeriksaan dilakukan terhadap `node.Classes` dan atribut `style`:
- Pola kelas warna sembarang: `bg-[#2563eb]`, `text-[#000]`, `border-[rgb(255,0,0)]`.
- Pola inline CSS style: `style="color: #2563eb;"`.

### 2.3. Payload Temuan Diagnostic
- **File:** `node.Span.File`
- **Line & Column:** Titik awal tag elemen (`node.Span.Line`, `node.Span.Column`).
- **Rule:** `"theme.hardcode-color"`
- **Severity:** `ir.SeverityWarn` (dapat di-override via `charites.yaml`).
- **Message:** `"Hardcode color (<value>) - gunakan semantic token dari global.css"`.
- **Hint:** `"Ganti dengan token semantik: #2563eb → bg-primary / text-primary"`.

### 2.4. Format Direktif Inline Ignore
Pengembang dapat mengecualikan baris tertentu menggunakan:
- JS/TSX: `// charites:ignore theme.hardcode-color`
- Astro HTML: `<!-- charites:ignore theme.hardcode-color -->`

---

## 3. Spesifikasi Kebutuhan Charites 1-SSOT Tri-Corpus

Setiap rule baru wajib merinci 3 sub-korpus uji coba yang akan ditempatkan di `tests/correctness/<rule_id>/`:

1. **`positive/` (True Positives):**
   - Berkas `.astro` dan `.tsx` yang sengaja memuat warna mentah (`bg-[#2563eb]`, `style="background: #ffffff"`).
   - **Kriteria:** Wajib memicu temuan diagnostic (`PositiveViolations > 0`).
2. **`negative/` (True Negatives / Zero-Noise Invariant):**
   - Berkas sah yang memuat token semantik resmi (`bg-primary`, `text-muted`, `border-destructive`).
   - **Kriteria:** Wajib menghasilkan 0 diagnostic (`NegativeViolations == 0`).
3. **`adversarial/` (False Positive Bait & Evasion):**
   - Berkas dengan tanda `#` yang bukan warna (misal anchor link `href="#main"`, ID selector `id="#header"`).
   - Berkas dengan komentar `// charites:ignore theme.hardcode-color`.
   - **Kriteria:** Wajib menghasilkan 0 diagnostic (`AdversarialViolations == 0`).
