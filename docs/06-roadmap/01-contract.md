# 06-ROADMAP: 01 - Phase 1 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-01-CONTRACT`
> **Tahapan:** Fase 1 - Kunci Kontrak Data (IR & Diagnostic)
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Graduated (All Phase Gates Passed)

Dokumen ini menetapkan kriteria kelulusan (*exit criteria*) dan gerbang transisi (*phase gate*) untuk **Fase 1 (Kunci Kontrak Data)** sebelum tim diizinkan melangkah ke **Fase 2 (Parser Frontend & IR Builder)**. Sesuai prinsip pemisahan otoritas arsitektur:
- **SPEC** = WHAT (Spesifikasi Kontrak Data & Invarian Fungsional)
- **ARCH** = HOW (Rancangan Arsitektur, Desain Iterator & Batasan Teknis)
- **TEST** = PROOF (Harness Pengujian, Skenario Smoke & Asersi Pembuktian)
- **QUALITY** = QUALITY THRESHOLD (Ambang Batas Kualitas Linter, Zero Dep & Hygiene)
- **ROADMAP** = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)

---

## 1. Deliverables Berkas Fase 1

1. **`internal/ir/node.go`**: Definisi struct `Node`, `Span`, enum `NodeType`, dan iterator `Walk()`.
2. **`internal/ir/diagnostic.go`**: Definisi struct `Diagnostic`, enum `Severity`, dan tag JSON.
3. **`internal/ir/node_test.go`**: Suite pengujian unit `ir.Node` dan benchmark traversal.
4. **`internal/ir/diagnostic_test.go`**: Suite pengujian serialisasi JSON `Diagnostic`.
5. **`internal/rules/rule.go`**: Interface baku `rules.Rule` yang mengonsumsi `*ir.Node` dan menghasilkan `[]ir.Diagnostic`.

---

## 2. Gerbang Evaluasi Kelulusan (Phase Gate DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [x] **`ROAD-01-GATE-001` (SPEC-01 Compliance = PASS):**
  - Seluruh invarian fungsional `ir.Node` (deterministic pre-order traversal, immediate early exit, post-construction immutability, parent-child bidirectional pointer, whitespace class tokenization, 1-indexed span) terpenuhi.
  - Penamaan Charites Rule ID, serialisasi flat JSON, dan 7-level Canonical Diagnostic Total Ordering (`DiagnosticOrderKey`) terpenuhi.

- [x] **`ROAD-01-GATE-002` (ARCH-01 Compliance = PASS):**
  - Paket `internal/ir` terbukti asiklik murni sebagai *pure leaf package* tanpa impor paket internal lain.
  - Traversal `Walk()` memanfaatkan `iter.Seq` Go 1.26 native.
  - Batasan siklus kepemilikan dipatuhi (`internal/parser` memiliki hak konstruksi, konsumen read-only).
  - Penyortir total `CompareDiagnostics` dan `SortDiagnostics` tersedia di `internal/ir`.

- [x] **`ROAD-01-GATE-003` (TEST-01 Compliance = PASS):**
  - Seluruh unit test lolos 100%: `TestNode_Walk`, `TestNode_EarlyExit`, `TestNode_ParentChildInvariant`, `TestNode_ClassTokenization`, `TestNode_SpanIndexing`, `TestNode_StructSize`, `TestNode_Helpers`, `TestDiagnostic_JSONDeterminism`, dan `TestDiagnostic_CollectionOrdering`.
  - `BenchmarkNode_Walk` berhasil dieksekusi dan dicatat dalam baseline performa.

- [x] **`ROAD-01-GATE-004` (QUAL-01 Compliance = PASS):**
  - Code coverage `internal/ir` mencapai $\ge 90\%$.
  - `golangci-lint run ./internal/ir/...` menghasilkan exit code `0` tanpa peringatan.
  - Ukuran struct `ir.Node` terbukti $\le 136$ bytes pada target 64-bit.
  - Invarian determinisme total ordering `QUAL-01-INVAR-001` terbukti 100% permutation-invariant.
  - Zero third-party dependencies di `go.mod`.

---

## 3. Gerbang Transisi ke Fase 2 (Parser & IR Builder)

Begitu keempat gerbang di atas berstatus **PASS**:
1. Rekam checkpoint git commit: `feat(ir): complete Phase 1 unified node, diagnostic contracts and iterator verification`.
2. Buka dokumen [docs/01-spec/02-parser.md](https://github.com/will2469/charites/blob/main/docs/01-spec/02-parser.md) untuk memulai spesifikasi parser AST Tailwind CSS v4 `@theme`, Astro, dan React/TSX (**Fase 2**).

