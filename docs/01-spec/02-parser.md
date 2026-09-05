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

### Aturan Ekstraksi Token & Netralitas Substrat:
1. **Raw Token Metadata Extraction:**
   - Properti `--color-<name>` diekstrak menjadi token tema mentah dalam `ThemeTokenRegistry` (`Variables map[string]string`).
   - Parser mencatat nama variabel asli dan nilai CSS mentah (misal `"--color-primary": "#2563eb"`, `"--color-primary-light": "rgba(37, 99, 235, 0.1)"`).
2. **Strict Rule-Agnostic Boundary:**
   - Parser **DILARANG KERAS** melakukan pemetaan semantik ekuivalensi opacity (contoh: memutuskan `primary/10 -> primary-light`).
   - Interpretasi bahwa akhiran `-light` atau `-subtle` merupakan ekuivalen dari modifier opacity adalah domain evaluasi **Rule #1 (`theme.hardcode-opacity-color`)** pada Fase 3, bukan domain parser.
3. **Hygiene:** Parser **MUST** mengabaikan komentar CSS `/* ... */` dan spasi kosong.

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
3. **Dukungan Tag Lengkap:**
   - Tag HTML standar (termasuk void elements).
   - Komponen kustom PascalCase (contoh: `<Card>`, `<HeaderNavigation>`).
   - Elemen bawaan Astro (`<Fragment>`, `<slot />`, `<slot name="..." />`).
4. **Batas Kedalaman Nesting (256 Level Stack Guard):**
   - Parser **MUST** membatasi kedalaman tumpukan parsing hingga maksimal **256 tingkat**.
   - **Definisi Flattening Normatif:** Jika tumpukan telah mencapai kedalaman 256 (`len(stack) == 256`), elemen anak berikutnya tetap diekstrak sebagai `*ir.Node` yang sah dan disematkan sebagai anak (*child*) di bawah node tingkat-256, tetapi **TIDAK DI-PUSH** ke dalam stack (*attached as flat siblings under the depth-256 parent*). Kedalaman stack parsing dijamin tidak pernah melebihi 256.

---

## 4. Spesifikasi JSX Structural Extractor (`.tsx`, `.jsx`)

Mengadopsi **Option B (JSX Structural Extractor)**: Sistem tidak bertindak sebagai TypeScript compiler penuh, melainkan melakukan ekstraksi struktural deterministik terhadap hierarki JSX untuk kebutuhan audit statis.

### Grammar Subset JSX yang Didukung:
1. **Elemen JSX:** Tag pembuka (`<Tag ...>`), tag self-closing (`<Tag ... />`), tag penutup (`</Tag>`), dan fragment (`<>...</>`).
2. **Atribut Target:** Ekstraksi atribut `className` (React) dan `class` (Preact/HTML), serta atribut struktural (`id`, `role`, `aria-*`, `data-*`).
3. **Kontrak Ekstraksi Template Literal Dinamis (Deterministic Dynamic Contract):**
   - **Teks Literal di Luar `${...}`:** Diekstrak sebagai token kelas statis yang pasti ke dalam `Classes` (dipisahkan oleh spasi).
   - **Teks di Dalam `${...}` (Opaque Dynamic Region):** Diperlakukan sebagai *opaque dynamic region*. Parser **DILARANG KERAS** mengevaluasi AST JavaScript, ekspresi ternary (`cond ? 'a' : 'b'`), atau logika boolean di dalam `${...}`.
   - **Dynamic Flag:** Keberadaan interpolasi `${...}` mengaktifkan penanda kelas dinamis pada node (misal `HasDynamicClasses = true` atau atribut penanda) agar rule engine di Fase 3 dapat menangani ekspresi dinamis tanpa false positive.
   - **Preservasi String Utuh:** Seluruh isi atribut template literal mentah (termasuk backtick dan blok `${...}`) disimpan secara verbatim pada `RawClasses`.
   - *Contoh Konkret:*
     ``className={`p-4 ${foo ? "bg-red" : "bg-blue"} text-sm`}``
     - `Classes` = `["p-4", "text-sm"]`
     - Wilayah `${...}` tidak dievaluasi.
     - `HasDynamicClasses` = `true`
     - `RawClasses` = ``"`p-4 ${foo ? \"bg-red\" : \"bg-blue\"} text-sm`"``
4. **Disambiguasi Karakter `<` Non-Tag:**
   - Di dalam komentar HTML/JSX (`<!-- ... -->`, `{/* ... */}`): diabaikan, tidak menghasilkan node elemen.
   - Di dalam nilai atribut string (`placeholder="a < b"`): diperlakukan sebagai teks literal atribut.
   - Di dalam ekspresi JSX kurung kurawal (`{count < 10 && ...}`): operator perbandingan `<` tidak memicu pembukaan tag elemen.
5. **Presisi Span:** `Span.Line` dan `Span.Column` **MUST** menunjuk ke karakter awal pembuka tag elemen JSX (`<`).

---

## 5. Spesifikasi IR Builder & Semantik Pemulihan (*Recovery Semantics*)

IR Builder (`internal/ir/builder.go`) menormalkan keluaran parser menjadi pohon `*ir.Node`:

### 5.1. Semantik Penanganan Input Cacat (*Recovery Semantics*):
Ketika menghadapi sintaks HTML/JSX yang cacat (*malformed input*), sistem **MUST** mematuhi aturan deterministik berikut:
1. **Token Discard & Resynchronization:**
   Jika tag pembuka tidak terbentuk sempurna (contoh: `<broken` tanpa `>` yang langsung diikuti `<button`), token yang rusak dibuang dan parser melakukan resinkronisasi ke karakter pembuka tag berikutnya (`<`). Elemen saudara berikutnya yang valid **MUST** tetap diekstrak normal.
2. **Partial Node Policy:** Tag malformed yang tidak memiliki nama tag atau sintaks dasar yang valid **DILARANG** dimasukkan ke dalam pohon IR.
3. **Aturan Resolusi Closing Tag & Stack Unwinding (`</X>`):**
   Ketika menemukan tag penutup `</X>`:
   - **Jika `X` ditemukan di dalam stack:** Pop seluruh elemen dari puncak stack hingga elemen `X`. Elemen perantara yang tidak ditutup secara eksplisit tetap menjadi anak sah dari induknya masing-masing.
   - **Jika `X` TIDAK ditemukan di dalam stack:** Buang token penutup `</X>` secara hening (*discard silently*), dan tumpukan stack **TIDAK BERUBAH** (*stack untouched*).
4. **Void Elements Autoclosing:**
   Tag void elements (`<img>`, `<input>`, `<br>`, `<hr>`, `<meta>`, `<link>`) ditangani secara mandiri (*self-contained*) tanpa memerlukan closing tag pasangan dan tidak memasukkan elemen anak ke dalam stack.
5. **Silent Non-Panicking Invariant:** Parser **DILARANG KERAS** memicu panic runtime pada input apa pun. Seluruh kegagalan sintaks dipulihkan secara hening (*silent recovery*) agar pemindaian berkas tidak terhenti.
