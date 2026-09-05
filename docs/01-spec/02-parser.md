# 01-SPEC: 02 - Frontend Parsers & IR Builder Specification

> **Kode Dokumen:** `SPEC-02-PARSER`
> **Tahapan:** Fase 2 - Parser Frontend & IR Builder
> **Status:** Ready for Review
> **Standar Rujukan:** IETF RFC 2119 / W3C HTML5 / JSX Specification / Tailwind CSS v4 Spec

Dokumen ini mendefinisikan spesifikasi kebutuhan fungsional untuk tiga parser frontend murni Go (**Tailwind CSS v4 `@theme`**, **Astro Component**, dan **TypeScript/JSX**) serta **IR Builder** yang menormalkan seluruh sintaks menjadi pohon `ir.Node`.

---

## 1. Ruang Lingkup Parsing Frontend

Mesin Charites **MUST** mendukung pemrosesan berkas frontend tanpa bergantung pada runtime Node.js atau compiler C:

1. **CSS Parser (`internal/parser/tailwind`):**
   - Membaca berkas `global.css` proyek.
   - Mengekstrak blok aturan `@theme` (spesifikasi Tailwind CSS v4).
   - Menghasilkan peta token semantik warna dan opacity yang diizinkan (*whitelisted design tokens*).
2. **Astro Parser (`internal/parser/astro`):**
   - Memisahkan blok *frontmatter script* (antara pembatas `---` pembuka dan penutup) dari blok *markup template*.
   - Menjaga offset baris (*line number*) template tetap sinkron 100% dengan berkas sumber `.astro`.
3. **TSX Parser (`internal/parser/tsx`):**
   - Membaca berkas `.tsx` dan `.jsx`.
   - Mengekstrak struktur elemen JSX, atribut `class` / `className`, serta ekspresi string literal dan template literals.
4. **IR Builder (`internal/ir/builder.go`):**
   - Menyatukan output dari parser Astro dan TSX ke dalam pohon `*ir.Node`.
   - Melakukan tokenisasi string class menjadi `Classes []string`.
   - Menghubungkan pointer `Parent` secara bi-direksional.

---

## 2. Spesifikasi Parser Tailwind CSS v4 (`@theme`)

Sistem **MUST** memindai deklarasi CSS custom properties di dalam blok `@theme`:

```css
@theme {
  --color-primary: #2563eb;
  --color-primary-light: rgba(37, 99, 235, 0.1);
  --color-destructive: #dc2626;
  --color-destructive-light: rgba(220, 38, 38, 0.1);
}
```

### Aturan Ekstraksi Token:
- Properti `--color-<name>` diekstrak menjadi token warna bernama `<name>`.
- Token yang berakhiran `-light` atau `-subtle` dicatat sebagai **Semantic Opacity Token** resmi.
- Parser **MUST** mengabaikan komentar CSS `/* ... */` dan whitespace.

---

## 3. Spesifikasi Parser Astro (`.astro`)

Berkas Astro memiliki struktur dua bagian:

```astro
---
// 1. Frontmatter JS/TS (Line 1 - 4)
import Button from '../components/Button.astro';
const title = "Dashboard";
---
<!-- 2. Markup Template (Line 5+) -->
<div class="p-4 bg-primary/10">
  <h1>{title}</h1>
</div>
```

### Invarian Parsing Astro:
1. **Frontmatter Isolation:** Parser **MUST** mendeteksi pembatas `---` di baris pertama dan penutup `---`.
2. **Presisi Penomoran Baris (Line Offset Preservation):** Baris pertama dari template HTML **MUST** mempertahankan nomor baris aslinya di file `.astro` (tidak di-reset ke line 1). Jika template dimulai di baris 5, `Span.Line` node pertama **MUST** bernilai `5`.
3. **Tag Support:** Mendukung elemen HTML standar, custom components (PascalCase), dan Astro built-ins (`<Fragment>`, `<slot />`).

---

## 4. Spesifikasi Parser TSX / JSX (`.tsx`, `.jsx`)

### Aturan Ekstraksi Elemen JSX:
1. **Atribut Class:**
   - Mendukung nama atribut `className` (React) dan `class` (Preact/HTML).
   - Ekstraksi nilai string statis: `className="p-4 bg-primary"`
   - Ekstraksi template literals statis: ``className={`p-4 bg-primary`}``
   - Deteksi interpolasi dinamis untuk rule audit kebersihan sintaks.
2. **Span Presisi:**
   - `Span.Line` dan `Span.Column` **MUST** menunjuk ke posisi pembuka tag elemen JSX.

---

## 5. Spesifikasi IR Builder (`internal/ir/builder.go`)

IR Builder bertindak sebagai normalisator akhir:

```go
package ir

// BuildTree mengubah representasi mentah parser menjadi pohon ir.Node yang valid
func BuildTree(rawElements []RawElement) *Node
```

### Invarian Pembentukan Pohon:
- **`Classes` Slicing:** String class asli di-split berdasarkan whitespace `\s+` dan nilai kosong dibuang.
- **Bi-Directional Pointer:** Seluruh elemen anak (`Children`) memiliki pointer `Parent` yang menunjuk ke elemen induknya.
- **Panic-Free Guarantee:** Input malformed (tag tidak tertutup, tag unclosed, atribut tanpa tanda kutip) **DILARANG KERAS** memicu panic pada Go runtime. Parser harus melakukan *graceful recovery* atau melewatkan node yang korup.
