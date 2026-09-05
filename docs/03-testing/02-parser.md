# 03-TESTING: 02 - Parser & IR Builder Verification Plan

> **Kode Dokumen:** `TEST-02-PARSER`
> **Tahapan:** Fase 2 - Parser Frontend & IR Builder
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Skenario Smoke & Asersi Pembuktian)
> **Status:** Ready for Execution
> **Standar Rujukan:** Modern Testing Principles & Go 1.26 Native Fuzzing Specification

Dokumen ini mendefinisikan strategi pengujian menyeluruh untuk parser Tailwind CSS, Astro, TSX, dan IR Builder, membuktikan invarian batas kedalaman stack, semantik pemulihan kegagalan sintaks, serta pengelolaan korpus regresi mutasi (*permanent regression corpus*).

---

## 1. Skenario Uji Coba Unit (`internal/parser/...`)

### 1.1. Uji Coba Tailwind Parser (`internal/parser/tailwind/theme_test.go`)
- **Input:** Berkas `tests/fixtures/global.css` yang memuat blok `@theme`.
- **Verifikasi:**
  - Token warna `--color-primary: #2563eb` berhasil diekstrak dengan key `primary`.
  - Token opacity `--color-primary-light` berhasil diekstrak dan dipetakan sesuai konvensi default.
  - Deklarasi di luar blok `@theme` diabaikan dengan benar.

### 1.2. Uji Coba Astro Component Lexer (`internal/parser/astro/lexer_test.go`)
- **Test Case 1 (Line Offset Preservation):**
  - Input: File `.astro` dengan 10 baris frontmatter.
  - Ekspektasi: Tag pembuka pertama template di baris 11 memiliki `Span.Line == 11`.
- **Test Case 2 (Nested Components):**
  - Input: Komponen `<Card><Button class="p-2">Klik</Button></Card>`.
  - Ekspektasi: Terbentuk hierarki 2 tingkat dengan pointer `Parent` yang valid.
- **Test Case 3 (Batas Kedalaman Nesting 255, 256, 257):**
  - Input dokumen dengan tag bersarang:
    - 255 tingkat nesting: berhasil diparse utuh (*accepted*).
    - 256 tingkat nesting: berhasil diparse utuh (*accepted*, batas puncak tercapai).
    - 257+ tingkat nesting: diproses secara deterministik tanpa stack overflow panic (*depth bounded*).
- **Test Case 4 (Semantik Pemulihan / Recovery Semantics):**
  - Input:
    ```html
    <div>
      <span>
        <broken
        <button class="valid">Klik</button>
      </span>
    </div>
    ```
  - Ekspektasi: Token `<broken` dibuang, parser melakukan resinkronisasi pada `<button`, dan elemen `button` dengan kelas `valid` berhasil diekstrak normal ke dalam IR tree.

### 1.3. Uji Coba JSX Structural Extractor (`internal/parser/tsx/extractor_test.go`)
- **Test Case 1 (JSX Attributes & Self-Closing):**
  - Input: `<div className="p-4 bg-primary" id="main" />`.
  - Ekspektasi: Terbaca tag `div`, atribut `className` dan `id`.
- **Test Case 2 (Static Template Literal Classes):**
  - Input: ``<div className={`p-4 bg-primary`} />``.
  - Ekspektasi: String di dalam backtick berhasil diekstrak ke dalam `Classes`.
- **Test Case 3 (Dynamic Template Literal & Interpolation):**
  - Input: ``<div className={`p-4 ${cond ? "bg-red" : "bg-blue"}`} />``.
  - Ekspektasi:
    - `Classes` memuat `["p-4"]`.
    - `RawClasses` memuat string mentah lengkap.
    - Flag penanda kelas dinamis aktif.
- **Test Case 4 (Ternary JSX Expressions):**
  - Input: `const x = a < b ? <Foo /> : <Bar />;`.
  - Ekspektasi: Simbol `<` dalam ekspresi ternary tidak merusak deteksi elemen `<Foo />` dan `<Bar />`.

---

## 2. Pengujian Ketahanan dengan Native Fuzzing & Korpus Regresi Permanen (`tests/fuzz/`)

Fuzz testing membuktikan bahwa parser memiliki ketahanan mutlak (*zero panic invariant*):

```go
package fuzz_test

import (
    "testing"
    "github.com/will2469/charites/internal/parser/astro"
    "github.com/will2469/charites/internal/parser/tsx"
)

func FuzzAstroParser(f *testing.F) {
    f.Add([]byte("---\nconst x = 1;\n---\n<div class=\"p-4\">Hello</div>"))
    f.Add([]byte("<div class=\"unclosed\"><span>text"))
    f.Add([]byte("--- malformed frontmatter"))
    f.Add([]byte("<broken <tag>"))
    f.Add([]byte(""))

    f.Fuzz(func(t *testing.T, data []byte) {
        // MUST NOT PANIC under any circumstances
        _, _ = astro.Parse(data)
    })
}

func FuzzTSXParser(f *testing.F) {
    f.Add([]byte("export default function App() { return <div className=\"p-4\">Hi</div>; }"))
    f.Add([]byte("const x = <Button className={`dynamic-${val}`} />"))
    f.Add([]byte("<>fragment unclosed"))
    f.Add([]byte("a < b ? <Foo /> : <Bar />"))

    f.Fuzz(func(t *testing.T, data []byte) {
        // MUST NOT PANIC
        _, _ = tsx.Extract(data)
    })
}
```

### Kebijakan Pengujian & Korpus Regresi:
1. **Safety Invariant:** Parser **MUST NOT** panic pada sembarang input byte (termasuk input acak, biner, atau terpotong).
2. **CI Verification:** Fuzzing dijalankan minimal **60 detik per target** pada pipa pengujian rilis.
3. **Permanent Regression Corpus:** Setiap input anomali atau kasus kegagalan yang ditemukan melalui fuzzing **WAJIB** disimpan secara permanen di direktori `tests/fixtures/regression/` sebagai pengujian unit regresi deterministik.

---

## 3. Fixtures Pengujian (`tests/fixtures/`)

Fase 2 menyiapkan berkas fixture nyata:
- `tests/fixtures/global.css`: Definisi `@theme` resmi.
- `tests/fixtures/sample.astro`: Komponen Astro valid dan kompleks.
- `tests/fixtures/sample.tsx`: Komponen React TSX dengan variasi sintaks JSX.
- `tests/fixtures/regression/`: Sampel input anomali hasil fuzzing terdokumentasi.

