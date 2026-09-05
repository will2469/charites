# 06-ROADMAP: 04 - Phase 4 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-04-ENGINE`
> **Tahapan:** Fase 4 - Konfigurasi, Concurrency Scanner & Traversal Engine
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Graduated (All Phase Gates Passed)

Dokumen ini menetapkan kriteria kelulusan (*exit criteria*) dan gerbang transisi (*phase gate*) untuk **Fase 4 (Konfigurasi, Concurrency Scanner & Traversal Engine)** sebelum tim diizinkan melangkah ke **Fase 5 (Reporter Output & CLI Entrypoint)**. Sesuai prinsip pemisahan otoritas arsitektur:
- **SPEC** = WHAT (Spesifikasi Kebutuhan Fungsional, Presedensi & Kontrak Engine)
- **ARCH** = HOW (Rancangan Pipeline, Enkapsulasi ActiveRule & Total Ordering)
- **TEST** = PROOF (Skenario Pengujian Unit, Race Detection & Pembuktian Batas)
- **QUALITY** = QUALITY THRESHOLD (Invarian Keamanan, Proteksi I/O & Ambang Batas Paket)
- **ROADMAP** = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)

---

## 1. Deliverables Berkas Fase 4

1. **`internal/config/config.go`**: Parser `charites.yaml` dengan penegakan prinsip **Default: YES**, struktur `ActiveRule` pembungkus `EffectiveSeverity`, dan resolusi presedensi (Registry $\rightarrow$ CLI Scope $\rightarrow$ Config Policy).
2. **`internal/config/ignore.go`**: Matcher sekuensial `.charitesignore` dengan kekebalan *builtin hard exclusion* dan fungsi *early directory pruning*.
3. **`internal/scanner/walker.go`**: Traversal direktori cepat dengan proteksi symlink (`DO NOT FOLLOW`), batas ukuran berkas 10 MB, dan keamanan target berkas langsung.
4. **`internal/scanner/pool.go`**: Concurrency worker pool berbasis Go channel dengan batas kapasitas $[1, 256]$ goroutine dan siklus pembatalan atomik via `context.Context`.
5. **`internal/analyzer/context.go`**: Konteks analisis terisolasi per-berkas dengan deep-copy caller map, slice cloning `DiagnosticsList`, dan parser direktif leksikal penekanan inline ignore multi-rule/node span.
6. **`internal/analyzer/engine.go`**: Traversal engine AST berbasis iterator Go 1.26 (`root.Walk()`).
7. **`internal/analyzer/sort.go`**: Modul pengurutan deterministik berbasis relasi pengurutan total (menggunakan langsung fungsi 1-SSOT `ir.SortDiagnostics(diags)`).
8. **Suite Pengujian Konkurensi**: Kumpulan test unit dan benchmark performa Fase 4 (`config_test.go`, `ignore_test.go`, `walker_test.go`, `pool_test.go`, `context_test.go`, `sort_test.go`).

Status Evaluasi Tata Kelola:
- **Implementation Status:** PASS (Config Parser, Ignore Matcher, Parallel Scanner, Worker Pool, AST Engine, Total Ordering Sorter, Modernized Go 1.26 Benchmarks).
- **Phase Gate Status:** PASS (All Phase Gates Passed & Verification Evidenced).

---

## 2. Gerbang Evaluasi Kelulusan (Phase Gate DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [x] **`ROAD-04-GATE-001` (SPEC-04 Compliance = PASS):**
  - Invarian Default: YES terbukti bekerja (seluruh rule aktif 100% saat `charites.yaml` absen).
  - Presedensi rule terbukti: Kebijakan config (`off`) mengalahkan seleksi CLI `--rule`.
  - Builtin hard exclusions terbukti kebal dari negasi `!`.
  - Proteksi symlink dan batas ukuran berkas 10 MB terbukti menolak traversal yang berbahaya atau boros memori.
  - Direktif inline ignore terbukti mendukung tata bahasa multi-rule, wildcard `*`, serta cakupan baris dan AST node span.
  - Parser leksikal mengenali komentar sah (`//`, `/* ... */`, `<!-- ... -->`, `${ ... }`) dan menolak substring dalam literal string, template literal, atribut JSX, dan plain HTML text.
  - Kontrak pembatalan interupsi (SIGINT/SIGTERM) menghentikan worker secara bersih dan membuang temuan parsial dengan exit code 130.

- [x] **`ROAD-04-GATE-002` (ARCH-04 Compliance = PASS):**
  - Rule singleton di registry tetap murni dan immutable melalui enkapsulasi `ActiveRule`.
  - Matcher `.charitesignore` mengevaluasi pola secara sekuensial (*last matching rule wins*).
  - Worker pool beroperasi dengan prinsip isolasi mandiri (*share-nothing worker*).
  - Isolasi memori caller terjamin dengan deep-copy `InlineIgnores` pada `NewContext`.
  - Ownership diagnosis terproteksi via `slices.Clone` pada `DiagnosticsList()`.
  - Pengurutan temuan diagnostic menggunakan komparator pengurutan total 7 tingkat.

- [x] **`ROAD-04-GATE-003` (TEST-04 Compliance = PASS):**
  - Seluruh skenario pengujian unit lolos 100% pada paket `internal/config`, `internal/scanner`, dan `internal/analyzer`.
  - Pengujian adversarial non-comment directives lolos 100% (0 false directives).
  - Verifikasi pembatalan `context.Context` membuktikan tidak ada kebocoran goroutine (*goroutine leak*).
  - Pengujian urutan deterministik membuktikan pengurutan total bebas kolisi acak.

- [x] **`ROAD-04-GATE-004` (QUAL-04 Compliance = PASS):**
  - Verifikasi bebas data race: `go test -race ./internal/...` lolos $100\%$ ($0$ race detected).
  - Seluruh ambang batas cakupan pengujian per-paket terpenuhi:
    - `internal/config`: $\ge 90\%$ line coverage (Aktual: $93.9\%$).
    - `internal/scanner`: $\ge 85\%$ line coverage (Aktual: $91.2\%$).
    - `internal/analyzer`: $\ge 90\%$ line coverage (Aktual: $98.3\%$).
  - Kompleksitas siklomatik $\le 12$ per fungsi (`gocyclo` pass).
  - `golangci-lint run ./internal/...` lolos 100% tanpa isu ($0$ error).

---

## 3. Final Verification Record

**Commit:** `1f15b3b`

### SPEC:
**PASS**
- Default: YES invariant verified (100% active rules when `charites.yaml` absent).
- 3-Tier Precedence verified (Config policy `off` overrides CLI flag).
- Builtin hard exclusions immune to negation (`!node_modules`).
- Direct-target safety rejects targets inside builtin exclusions.
- Symlink protection & MaxScanFileSize (10 MB) verified.
- Inline ignore directives verified for multi-rule, wildcard `*`, and AST node span.
- Comment-context correctness verified: string literals, template literals, JSX attribute strings, and plain text rejected as directives.

### ARCH:
**PASS**
- Rule singleton immutability guaranteed via `ActiveRule`.
- Concurrency worker pool adheres to share-nothing worker isolation.
- `NewContext()` deep-copies `InlineIgnores` to guarantee caller map isolation.
- `DiagnosticsList()` clones internal diagnostics (`slices.Clone`) to protect slice ownership.
- 7-tier total ordering deterministic sorting verified via `ir.SortDiagnostics`.

### TEST:
**PASS**
- `go test -race ./internal/...` -> PASS (0 race, 0 failure, 0 panic across all packages)
- `go test ./internal/...` -> PASS (100% pass)
- Context cancellation terminates workers cleanly with 0 goroutine leaks.

### QUALITY:
**PASS**
- `internal/config` coverage: **93.9%** (Threshold $\ge 90\%$)
- `internal/scanner` coverage: **91.2%** (Threshold $\ge 85\%$)
- `internal/analyzer` coverage: **98.3%** (Threshold $\ge 90\%$)
- Data race: **0**
- Linter: **0** (`golangci-lint run ./internal/...` exit code 0)
- Cyclomatic complexity: **$\le 12$** per function (`gocyclo -over 12` clean across all production and test functions)

### BENCHMARKS:
**PASS**
- `BenchmarkIgnore_Matcher-8`: 672 ops, 1,669,397 ns/op, 831,594 B/op, 17,000 allocs/op
- `BenchmarkWalker_DirectoryTraversal-8`: 1,419 ops, 848,216 ns/op, 323,431 B/op, 3,667 allocs/op
- `BenchmarkEndToEnd_ScanPipeline-8`: 6,828 ops, 147,570 ns/op, 49,797 B/op, 437 allocs/op
- `BenchmarkAnalyzer_EngineTraversal-8`: 170,152 ops, 6,786 ns/op, 13,392 B/op, 68 allocs/op

### CI / BUILD:
**PASS**
- `make all`: PASS (Deterministic Build + Race Test + Linter + Gofmt)
- `make cross-compile`: PASS (linux/amd64, linux/arm64, darwin/arm64, windows/amd64)

### Reviewer:
**PASS**

---

## 4. Gerbang Transisi ke Fase 5 (CLI Entrypoint & Reporters)

Dengan seluruh evidence empiris diverifikasi secara ketat:

$$\text{Phase 4} \longrightarrow \mathbf{GRADUATED} \longrightarrow \mathbf{Phase\ 5\ OPEN}$$

1. Buka dokumen `docs/01-spec/05-cli.md` untuk merancang CLI entrypoint (`charites scan`, alias `check`/`run`).
2. Implementasikan flag parsing (A-E), runner scan engine terintegrasi, formatters (ANSI terminal berwarna dan JSON payload formatter).

