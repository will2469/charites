# 03-TESTING: 02 - Parser & IR Builder Verification Plan

> **Kode Dokumen:** `TEST-02-PARSER`
> **Tahapan:** Fase 2 - Parser Frontend & IR Builder
> **Status:** Ready for Review
> **Standar Rujukan:** Modern Testing Principles & Go 1.26 Native Fuzzing Specification

Dokumen ini mendefinisikan strategi pengujian menyeluruh untuk parser Tailwind CSS, Astro, TSX, dan IR Builder, termasuk skenario **Native Fuzzing** untuk menjamin ketahanan dari *unhandled panic*.

---

## 1. Skenario Uji Coba Unit (`internal/parser/...`)

### 1.1. Uji Coba Tailwind Parser (`internal/parser/tailwind/theme_test.go`)
- **Input:** Berkas `tests/fixtures/global.css` yang memuat blok `@theme`.
- **Verifikasi:**
  - Token warna `--color-primary: #2563eb` berhasil diekstrak dengan key `primary`.
  - Token opacity `--color-primary-light` berhasil diekstrak dan dipetakan.
  - Deklarasi di luar blok `@theme` diabaikan dengan benar.

### 1.2. Uji Coba Astro Parser (`internal/parser/astro/lexer_test.go`)
- **Test Case 1 (Line Offset Preservation):**
  - Input: File `.astro` dengan 10 baris frontmatter.
  - Ekspektasi: Tag pembuka pertama template di baris 11 memiliki `Span.Line == 11`.
- **Test Case 2 (Nested Components):**
  - Input: Komponen `<Card><Button class="p-2">Klik</Button></Card>`.
  - Ekspektasi: Terbentuk hirarki 2 tingkat dengan pointer `Parent` yang valid.

### 1.3. Uji Coba TSX Parser (`internal/parser/tsx/visitor_test.go`)
- **Test Case 1 (JSX Attributes):**
  - Input: `<div className="p-4 bg-primary" id="main" />`.
  - Ekspektasi: Terbaca tag `div`, atribut `className` dan `id`.
- **Test Case 2 (Template Literal Classes):**
  - Input: ``<div className={`p-4 bg-primary`} />``.
  - Ekspektasi: String di dalam backtick berhasil diekstrak ke dalam `Classes`.

---

## 2. Pengujian Ketahanan dengan Native Fuzzing (`tests/fuzz/`)

Kode frontend di dunia nyata seringkali mengandung sintaks yang tidak valid, tag tidak tertutup, atau karakter khusus yang tidak lazim. Fuzz testing wajib membuktikan bahwa parser **tidak pernah mengalami panic**:

```go
package fuzz_test

import (
    "testing"
    "github.com/will2469/charites/internal/parser/astro"
    "github.com/will2469/charites/internal/parser/tsx"
)

func FuzzAstroParser(f *testing.F) {
    // Seed corpus
    f.Add([]byte("---\nconst x = 1;\n---\n<div class=\"p-4\">Hello</div>"))
    f.Add([]byte("<div class=\"unclosed\"><span>text"))
    f.Add([]byte("--- malformed frontmatter"))
    f.Add([]byte(""))

    f.Fuzz(func(t *testing.T, data []byte) {
        // MUST NOT PANIC under any circumstances
        _, _ = astro.Parse(data)
    })
}

func FuzzTSXParser(f *testing.F) {
    // Seed corpus
    f.Add([]byte("export default function App() { return <div className=\"p-4\">Hi</div>; }"))
    f.Add([]byte("const x = <Button className={`dynamic-${val}`} />"))
    f.Add([]byte("<>fragment unclosed"))

    f.Fuzz(func(t *testing.T, data []byte) {
        // MUST NOT PANIC
        _, _ = tsx.Parse(data)
    })
}
```

### Kriteria Lolos Fuzzing:
- Dijalankan selama minimal **60 detik per target fuzzing** di pipeline CI lokal tanpa crash (`exit status 0`).

---

## 3. Fixtures Pengujian (`tests/fixtures/`)

Fase 2 menyiapkan berkas fixture nyata:
- `tests/fixtures/global.css`: Definisi `@theme` resmi.
- `tests/fixtures/sample.astro`: Komponen Astro valid dan kompleks.
- `tests/fixtures/sample.tsx`: Komponen React TSX dengan berbagai ragam penulisan class.
