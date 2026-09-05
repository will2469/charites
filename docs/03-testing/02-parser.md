# 03-TESTING: 02 - Parser & IR Builder Verification Plan

> **Kode Dokumen:** `TEST-02-PARSER`
> **Tahapan:** Fase 2 - Parser Frontend & IR Builder
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Skenario Smoke & Asersi Pembuktian)
> **Status:** Graduated (All Phase Gates Passed)
> **Standar Rujukan:** Modern Testing Principles & Go 1.26 Native Fuzzing Specification

Dokumen ini mendefinisikan strategi pengujian menyeluruh untuk parser Tailwind CSS, Astro, TSX, dan IR Builder, membuktikan invarian batas kedalaman stack, semantik pemulihan kegagalan sintaks, serta pengelolaan korpus regresi mutasi (*permanent regression corpus*).

---

## 1. Skenario Uji Coba Unit (`internal/parser/...`)

### 1.1. Uji Coba Tailwind Parser (`internal/parser/tailwind/theme_test.go`)
- **Input:** Berkas `tests/fixtures/global.css` yang memuat blok `@theme`.
- **Verifikasi:**
  - Token `--color-primary: #2563eb` berhasil diekstrak ke dalam `ThemeTokenRegistry.Variables["--color-primary"]`.
  - Token `--color-primary-light: rgba(37, 99, 235, 0.1)` berhasil diekstrak mentah tanpa pemetaan ekuivalensi semantik prematur.
  - Deklarasi di luar blok `@theme` diabaikan dengan benar.

### 1.2. Uji Coba Astro Component Lexer (`internal/parser/astro/lexer_test.go`)
- **Test Case 1 (Line Offset Preservation):**
  - Input: File `.astro` dengan 10 baris frontmatter.
  - Ekspektasi: Tag pembuka pertama template di baris 11 memiliki `Span.Line == 11`.
- **Test Case 2 (Nested & Custom Components):**
  - Input: Komponen `<Card><HeaderNavigation /><Button class="p-2">Klik</Button></Card>`.
  - Ekspektasi: Komponen PascalCase teridentifikasi dengan relasi `Parent`/`Children` valid.
- **Test Case 3 (Astro Slots & Fragments):**
  - Input: `<Fragment><slot name="header" /><slot /></Fragment>`.
  - Ekspektasi: Elemen slot dan fragment terurai menjadi node IR yang tepat.
- **Test Case 4 (Void Elements Autoclosing):**
  - Input: `<div><img src="pic.jpg"><input type="text"><p>Teks</p></div>`.
  - Ekspektasi: `<img>` dan `<input>` tidak menjebak `<p>` sebagai elemen anak; ketiga elemen bersaudara sejajar di bawah `<div>`.
- **Test Case 5 (Batas Kedalaman Nesting 255, 256, 257+):**
  - 255 tingkat: diparse normal (*accepted*).
  - 256 tingkat: diparse normal (*accepted* pada kedalaman puncak).
  - 257+ tingkat: elemen anak disematkan sebagai anak flat di bawah node tingkat-256 (*depth-256 flat siblings*), stack tidak pernah melebihi 256.

### 1.3. Uji Coba JSX Structural Extractor (`internal/parser/tsx/extractor_test.go`)
- **Test Case 1 (JSX Attributes & Self-Closing):**
  - Input: `<div className="p-4 bg-primary" id="main" />`.
  - Ekspektasi: Terbaca tag `div`, atribut `className` dan `id`.
- **Test Case 2 (Fragments `<>...</>`):**
  - Input: `<><span>A</span><span>B</span></>`.
  - Ekspektasi: Menghasilkan `NodeFragment` dengan dua anak `span`.
- **Test Case 3 (Static Template Literal Classes):**
  - Input: ``<div className={`p-4 bg-primary`} />``.
  - Ekspektasi: String di dalam backtick diekstrak ke dalam `Classes`.
- **Test Case 4 (Dynamic Template Literal & Opaque Interpolation):**
  - Input: ``<div className={`p-4 ${cond ? "bg-red" : "bg-blue"} text-sm`} />``.
  - Ekspektasi:
    - `Classes` tepat memuat `["p-4", "text-sm"]`.
    - Wilayah `${...}` diisolasi buram (*opaque*), tidak memicu parsing JS AST.
    - `RawClasses` menyimpan string mentah lengkap.
    - Flag kelas dinamis aktif pada node.
- **Test Case 5 (Disambiguasi Karakter `<` Non-Tag):**
  - Sub-test 5a (Komentar): `{/* <div class="ignore"> */}` dan `<!-- <div /> -->` tidak menghasilkan node elemen.
  - Sub-test 5b (String Atribut): `<input placeholder="a < b" />` tidak memicu pembukaan tag baru.
  - Sub-test 5c (Ekspresi JSX): `{count < 10 && <Badge />}` membedakan operator perbandingan `<` dengan tag `<Badge />`.

### 1.4. Uji Coba Semantik Pemulihan (*Recovery Semantics* / `internal/ir/builder_test.go`)
- **Test Case 1 (Single Malformed Opening Tag):**
  - Input: `<div><span><broken <button class="valid">Klik</button></span></div>`.
  - Ekspektasi: `<broken` dibuang, resinkronisasi ke `<button`, tombol diekstrak normal.
- **Test Case 2 (Multiple Consecutive Sibling Malformations):**
  - Input: `<broken1 <broken2 <div class="ok">Text</div>`.
  - Ekspektasi: Seluruh token cacat dilewati, `<div class="ok">` berhasil dirakit.
- **Test Case 3 (Unmatched Closing Tag Discard):**
  - Input: `<div><p>Hello</p></unknown><span>World</span></div>`.
  - Ekspektasi: Token `</unknown>` dibuang secara hening (*discarded*), stack `div` tidak terganggu, `<span>` tetap menjadi anak dari `div`.
- **Test Case 4 (Stack Unwinding pada Unclosed Intermediate Elements):**
  - Input: `<div><span><button>Teks</div>`.
  - Ekspektasi: Tag `</div>` mem-pop `<button>` dan `<span>` secara implisit hingga menutup `<div>`. Stack kembali ke kondisi stabil tanpa crash.

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

