# 01-SPEC: 02 - Frontend Parsers & IR Builder Specification

> **Kode Dokumen:** `SPEC-02-PARSER`
> **Tahapan:** Fase 2 - Parser Frontend & IR Builder
> **Peran Pilar:** SPEC = WHAT (Spesifikasi Ekstraktor Frontend & Kontrak Normalisasi)
> **Status:** Ready for Execution
> **Standar Rujukan:** IETF RFC 2119 / W3C HTML5 / JSX Specification / Tailwind CSS v4 Spec

Dokumen ini mendefinisikan spesifikasi kebutuhan fungsional untuk tiga ekstraktor frontend murni Go (**Tailwind CSS v4 `@theme` Token Extractor**, **Astro Component Lexer**, dan **TSX / JSX Structural Extractor**) serta **IR Builder** yang menormalkan seluruh sintaks menjadi pohon `ir.Node` yang netral (*rule-agnostic substrate*).

---

## 1. Ruang Lingkup Parsing & Substrat Netral (Rule-Agnostic Substrate)

Mesin Charites **MUST** mendukung pemrosesan berkas frontend tanpa bergantung pada runtime Node.js atau compiler C:

1. **Invarian Netralitas Rule (Rule-Agnostic Boundary):**
   Paket parser (`internal/parser/*`) dan IR Builder murni bertugas mengekstrak struktur markup, atribut, dan token tema ke dalam representasi `ir.Node`. Parser **DILARANG KERAS** melakukan evaluasi pelanggaran rule, mendeteksi token terlarang (misal `bg-primary/10`), atau menghasilkan temuan diagnosis. Evaluasi pelanggaran sepenuhnya menjadi domain `internal/rules` pada Fase 3.
2. **Komponen Layer Parsing:**
   - **CSS Parser (`internal/parser/tailwind`):** Mengekstrak definisi variabel di dalam blok `@theme` pada `global.css`.
   - **Astro Lexer (`internal/parser/astro`):** Memisahkan frontmatter script (`---`) dari template markup sambil menjaga offset baris akurat 100%.
   - **JSX Structural Extractor (`internal/parser/tsx`):** Mengekstrak struktur elemen JSX, atribut, kelas, dan segmen string/template literals (Option B: structural scanning tanpa overhead kompilator TypeScript penuh).
   - **IR Builder (`internal/ir/builder.go`):** Merakit pohon `*ir.Node` terpadu dengan relasi bi-direksional `Parent`/`Children` dan tokenisasi kelas.

---

## 2. Spesifikasi Ekstraktor Token Tailwind CSS v4 (`@theme`)

Sistem **MUST** memindai deklarasi custom properties CSS di dalam blok `@theme`:

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
- **Konvensi Token Opacity (Charites Opinionated Default):**
  Secara default bawaan, Charites menerapkan konvensi opini desain bahwa token yang memiliki akhiran `-light` atau `-subtle` (misal `--color-primary-light`) dicatat sebagai **Semantic Opacity Token** resmi. Proyek dapat menyesuaikan atau menimpa pemetaan ini melalui konfigurasi `charites.yaml` pada Fase 4.
- Parser **MUST** mengabaikan komentar CSS `/* ... */` dan spasi kosong.

---

## 3. Spesifikasi Astro Component Lexer (`.astro`)

Berkas Astro memisahkan frontmatter dan markup template:

```astro
---
// 1. Frontmatter JS/TS (Baris 1 - 4)
import Button from '../components/Button.astro';
const title = "Dashboard";
---
<!-- 2. Markup Template (Baris 5+) -->
<div class="p-4 bg-primary/10">
  <h1>{title}</h1>
</div>
```

### Invarian Parsing Astro:
1. **Frontmatter Isolation:** Parser **MUST** mendeteksi pembatas `---` pembuka di awal berkas dan pembatas `---` penutup.
2. **Presisi Penomoran Baris (Line Offset Preservation):** Baris pertama dari template markup **MUST** mempertahankan nomor baris fisiknya di dalam berkas sumber `.astro`. Jika template dimulai pada baris 5, `Span.Line` elemen pertama **MUST** bernilai `5`.
3. **Dukungan Tag:** Mendukung tag HTML standar, komponen kustom PascalCase, dan elemen bawaan Astro (`<Fragment>`, `<slot />`).
4. **Batas Kedalaman Nesting:** Parser **MUST** menerapkan batas kedalaman nesting berhingga sebesar **256 tingkat** untuk mencegah *call-stack overflow* pada input abnormal.

---

## 4. Spesifikasi JSX Structural Extractor (`.tsx`, `.jsx`)

Mengadopsi **Option B (JSX Structural Extractor)**: Sistem tidak bertindak sebagai TypeScript compiler penuh, melainkan melakukan ekstraksi struktural deterministik terhadap hierarki JSX untuk kebutuhan audit statis.

### Grammar Subset JSX yang Didukung:
1. **Elemen JSX:** Tag pembuka (`<Tag ...>`), tag self-closing (`<Tag ... />`), tag penutup (`</Tag>`), dan fragment (`<>...</>`).
2. **Atribut Target:** Ekstraksi atribut `className` (React) dan `class` (Preact/HTML), serta atribut struktural (`id`, `role`, `aria-*`, `data-*`).
3. **Representasi String & Template Literals di IR:**
   - **Static Literal:** `className="p-4 bg-primary"` $\rightarrow$ `Classes = ["p-4", "bg-primary"]`, `RawClasses = "p-4 bg-primary"`.
   - **Static Template Literal:** ``className={`p-4 bg-primary`}`` $\rightarrow$ `Classes = ["p-4", "bg-primary"]`.
   - **Dynamic Interpolated Literal:**
     Contoh: ``className={`p-4 ${condition ? "bg-red" : "bg-blue"}`}`` atau ``className={`p-4 ${foo}`}``
     - Token statis yang dapat dipastikan diekstrak ke `Classes` (misal: `["p-4"]`).
     - Seluruh teks mentah termasuk ekspresi interpolasi `${...}` disimpan utuh pada `RawClasses`.
     - Node ditandai memiliki kelas dinamis (`HasDynamicClasses = true` atau atribut flag) agar rule evaluator di Fase 3 dapat menangani ekspresi dinamis tanpa *false positive*.
4. **Presisi Span:** `Span.Line` dan `Span.Column` **MUST** menunjuk ke karakter awal pembuka tag elemen JSX (`<`).

---

## 5. Spesifikasi IR Builder & Semantik Pemulihan (*Recovery Semantics*)

IR Builder (`internal/ir/builder.go`) menormalkan keluaran parser menjadi pohon `*ir.Node`:

### 5.1. Semantik Penanganan Input Cacat (*Recovery Semantics*):
Ketika menghadapi sintaks HTML/JSX yang cacat (*malformed input*), sistem **MUST** mengikuti semantik deterministik berikut:
1. **Token Discard & Resynchronization:** Jika tag pembuka tidak terbentuk sempurna (contoh: `<broken` tanpa `>` yang langsung diikuti `<button`), token yang rusak dibuang dan parser melakukan resinkronisasi ke karakter pembuka tag berikutnya (`<`).
2. **Partial Node Policy:** Tag malformed yang tidak memiliki nama tag atau sintaks dasar yang valid **DILARANG** dimasukkan ke dalam pohon IR. Elemen saudara berikutnya yang valid **MUST** tetap diekstrak normal.
3. **Stack Unwinding & Void Elements:**
   - Tag penutup void elements (`<img>`, `<input>`, `<br>`, `<hr>`) ditangani secara implisit tanpa memerlukan closing tag pasangan.
   - Jika ditemukan penutup tag yang melompati tingkat (misal `</span>` menutup konteks di atas unclosed element), stack di-pop hingga elemen yang cocok.
4. **Silent Non-Panicking Execution:** Parser **DILARANG KERAS** memicu panic runtime pada input apa pun. Seluruh kegagalan sintaks dipulihkan secara hening (*silent recovery*) agar pemindaian berkas tidak terhenti.

