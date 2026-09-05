# 06-ROADMAP: 02 - Phase 2 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-02-PARSER`
> **Tahapan:** Fase 2 - Parser Frontend & IR Builder
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Ready for Execution

Dokumen ini menetapkan kriteria kelulusan (*exit criteria*) dan gerbang transisi (*phase gate*) untuk **Fase 2 (Parser Frontend & IR Builder)** sebelum tim diizinkan melangkah ke **Fase 3 (Rule Contract & Rule #1 Proving Ground: `theme.hardcode-opacity-color`)**. Sesuai prinsip pemisahan otoritas arsitektur:
- **SPEC** = WHAT (Spesifikasi Ekstraktor Frontend & Kontrak Normalisasi)
- **ARCH** = HOW (Rancangan Pipeline Parsing, Boundary Netral & Recovery)
- **TEST** = PROOF (Harness Pengujian, Skenario Smoke & Asersi Pembuktian)
- **QUALITY** = QUALITY THRESHOLD (Ambang Batas Ketahanan, Security & Resource Budgets)
- **ROADMAP** = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)

---

## 1. Deliverables Berkas Fase 2

1. **`internal/parser/tailwind/`**: Ekstraktor blok `@theme` di `global.css` untuk warna dan token opacity semantik.
2. **`internal/parser/astro/`**: Splitter frontmatter `---` dan template markup dengan offset baris presisi.
3. **`internal/parser/tsx/`**: JSX Structural Extractor untuk elemen JSX, atribut target, dan template literal classes.
4. **`internal/ir/builder.go`**: Perakit pohon `*ir.Node` terpadu dengan relasi bi-direksional `Parent`/`Children` dan stack guard 256.
5. **`tests/fixtures/`**: Kumpulan sampel berkas sumber frontend (`global.css`, `sample.astro`, `sample.tsx`, dan korpus regresi).
6. **`tests/fuzz/`**: Native fuzzing suite untuk parser Astro dan TSX (`astro_fuzz_test.go`, `tsx_fuzz_test.go`).

---

## 2. Gerbang Evaluasi Kelulusan (Phase Gate DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [ ] **`ROAD-02-GATE-001` (SPEC-02 Compliance = PASS):**
  - Ekstraksi token Tailwind `@theme` mematuhi konvensi default opacity semantik.
  - Parser Astro terbukti mempertahankan offset baris template (tidak ter-reset ke line 1).
  - JSX Structural Extractor mendukung grammar subset target (elemen, fragment, static/dynamic classes di `RawClasses`).
  - Semantik pemulihan kegagalan sintaks (*recovery semantics*) berjalan deterministik tanpa partial corrupted node di IR.
  - Substrat murni netral (*rule-agnostic*), tanpa evaluasi rule apa pun di layer parser.

- [ ] **`ROAD-02-GATE-002` (ARCH-02 Compliance = PASS):**
  - Seluruh sub-parser terimplementasi menggunakan Go Standard Library tanpa runtime Node.js atau CGO.
  - Mekanisme assembly IR Builder mengunci batas kedalaman hierarki hingga maksimal 256 tingkat (*stack guard*).

- [ ] **`ROAD-02-GATE-003` (TEST-02 Compliance = PASS):**
  - Pengujian unit lolos 100%: verifikasi line offset, kedalaman nesting 255/256/257, semantik recovery, dan template literal dinamis.
  - Fuzz test `FuzzAstroParser` dan `FuzzTSXParser` lolos minimal 60 detik tanpa panic.
  - Korpus regresi permanen terbentuk di `tests/fixtures/regression/`.

- [ ] **`ROAD-02-GATE-004` (QUAL-02 Compliance = PASS):**
  - Invarian Zero-Panic terpenuhi pada input arbitrer.
  - Scanning linear tanpa ketergantungan regex rawan ReDoS.
  - Anggaran memori terpenuhi (peak heap live bytes $\le 4\times$ ukuran berkas sumber).
  - Unit test mencapai coverage $\ge 85\%$ dan `golangci-lint` menghasilkan exit code `0`.

---

## 3. Gerbang Transisi ke Fase 3 (Rule #1 Proving Ground)

Begitu keempat gerbang di atas berstatus **PASS**:
1. Rekam checkpoint git commit: `feat(parser): complete Phase 2 frontend parsers, recovery semantics and AST-to-IR tree builder`.
2. Buka dokumen [docs/01-spec/03-rules.md](file:///home/will/Monorepo/charites/docs/01-spec/03-rules.md) untuk memulai implementasi kernel rule dan pembuktian Rule #1: **`theme.hardcode-opacity-color`** (**Fase 3**).

