# 06-ROADMAP: 02 - Phase 2 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-02-PARSER`
> **Tahapan:** Fase 2 - Parser Frontend & IR Builder
> **Status:** Ready for Review

Dokumen ini menetapkan kriteria kelulusan (*exit criteria*) dan gerbang transisi (*transition gate*) untuk **Fase 2 (Parser Frontend & IR Builder)** sebelum melangkah ke **Fase 3 (Rule Contract & Rule #1: `theme.hardcode-opacity-color`)**.

---

## 1. Deliverables Wajib Fase 2

1. **`internal/parser/tailwind/`**: Parser `@theme` di `global.css` untuk ekstraksi warna dan token opacity.
2. **`internal/parser/astro/`**: Splitter frontmatter `---` dan template markup dengan offset baris presisi.
3. **`internal/parser/tsx/`**: Visitor JSX tag, atribut `class`/`className`, dan template literal strings.
4. **`internal/ir/builder.go`**: Perakit pohon `*ir.Node` terpadu dengan relasi `Parent`/`Children` dan tokenisasi class.
5. **`tests/fixtures/`**: Kumpulan sampel berkas sumber frontend (`global.css`, `sample.astro`, `sample.tsx`).
6. **`tests/fuzz/`**: Native fuzzing suite untuk parser Astro dan TSX (`astro_fuzz_test.go`, `tsx_fuzz_test.go`).

---

## 2. Checklist Definition of Done (DoD) Fase 2

- [ ] Seluruh sub-parser terimplementasi menggunakan Go Standard Library tanpa ketergantungan runtime Node.js atau CGO.
- [ ] Parser Astro terbukti mempertahankan nomor baris asli template (tidak di-reset ke line 1).
- [ ] Parser TSX terbukti mengekstrak nama tag, atribut `className`, dan template literal tanpa panic.
- [ ] `internal/ir/builder.go` berhasil membentuk pohon hirarki lengkap dengan pointer `Parent` yang valid.
- [ ] Fuzz test `FuzzAstroParser` dan `FuzzTSXParser` berjalan minimal 60 detik tanpa satupun *unhandled panic*.
- [ ] Unit test paket parser mencapai coverage $\ge 85\%$ (`go test -cover ./internal/parser/...`).
- [ ] `golangci-lint run ./internal/parser/... ./internal/ir/...` bersih 100%.

---

## 3. Gerbang Transisi ke Fase 3 (Rule #1 Proving Ground)

Begitu seluruh checklist di atas terpenuhi:
1. Buat git commit: `feat(parser): implement tailwind v4, astro, tsx parsers and AST-to-IR tree builder`.
2. Buka dokumen `docs/01-spec/03-rules.md` untuk memulai spesifikasi evaluasi Rule #1: **`theme.hardcode-opacity-color`** beserta infrastruktur Tri-Corpus Argus.
