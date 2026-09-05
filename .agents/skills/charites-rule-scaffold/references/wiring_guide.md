# Charites Rule Wiring Guide & Registry Reference
**Target System:** Charites Rule Registry (`internal/rules/registry.go`)

---

## 1. Aturan Penamaan Canonical Semgrep ID

Setiap rule di Charites menggunakan **Semgrep Canonical ID** `<category>.<slug>`:
- Kategori resmi: `theme`, `a11y`, `responsive`, `perf`, `tailwind`
- Contoh: `theme.hardcode-color`, `theme.hardcode-opacity-color`, `a11y.alt-text`
- **Dilarang:** Menggunakan penomoran singkatan seperti `T01`/`txx` atau `A01`/`axx`.

---

## 2. Menambahkan Rule ke Registry

Setiap rule baru yang di-scaffold perlu didaftarkan ke `internal/rules/registry.go`:

```go
package rules

func RegisterBuiltinRules(r *Registry) {
    // Built-in theme rules
    r.Register(NewThemeHardcodeColorRule())
    r.Register(NewThemeHardcodeOpacityColorRule())

    // Built-in accessibility rules
    r.Register(NewA11yAltTextRule())

    // Rule baru yang di-scaffold:
    r.Register(New<Category><Slug>Rule())
}
```

---

## 3. Struktur Data Rule Interface

Setiap analyzer rule mengimplementasikan interface `Rule`:

```go
type Rule interface {
    ID() string                     // Mengembalikan canonical Semgrep ID, misal "theme.hardcode-opacity-color"
    Category() string               // Mengembalikan "theme", "a11y", dll.
    DefaultSeverity() ir.Severity   // ir.SeverityCritical, ir.SeverityHigh, dll.
    Description() string            // Ringkasan 1 kalimat
    Analyze(file *ir.File) []ir.Diagnostic
}
```

---

## 4. Penanganan Direktif Supresi (`charites:ignore`)

Engine direktif Charites (`internal/directives/`) mem-parsing komentar kode:
- `// charites:ignore <category>.<slug> [reason]`
- `/* charites:ignore <category>.<slug> [reason] */`
- `<!-- charites:ignore <category>.<slug> [reason] -->`

Directive parser memvalidasi kesesuaian ID secara persis terhadap `rule.ID()`. Jika komentar menggunakan nama rule yang tidak terdaftar (seperti `T01` atau `A01`), directive dianggap tidak cocok dan pelanggaran tetap dilaporkan (*fail-closed*).

---

## 5. Eksekusi Pengujian Otomatis

Setelah scaffolding dan registrasi selesai:

```bash
# Menjalankan test spesifik untuk rule
go test -v ./tests/correctness/<category>.<slug>/...

# Menjalankan verifikasi matriks adopsi Tri-Corpus
go test -v ./tests -run TestGoldenCorpus_AdoptionMatrix
```
