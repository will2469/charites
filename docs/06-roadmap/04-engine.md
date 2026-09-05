# 06-ROADMAP: 04 - Phase 4 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-04-ENGINE`
> **Tahapan:** Fase 4 - Konfigurasi, Concurrency Scanner & Traversal Engine
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Ready for Execution

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
5. **`internal/analyzer/context.go`**: Konteks analisis terisolasi per-berkas dengan parser direktif penekanan inline ignore multi-rule dan cakupan node span.
6. **`internal/analyzer/engine.go`**: Traversal engine AST berbasis iterator Go 1.26 (`root.Walk()`).
7. **`internal/analyzer/sort.go`**: Modul pengurutan deterministik berbasis relasi pengurutan total (*total ordering comparator*).
8. **Suite Pengujian Konkurensi**: Kumpulan test unit dan benchmark performa Fase 4 (`config_test.go`, `ignore_test.go`, `walker_test.go`, `pool_test.go`, `context_test.go`, `sort_test.go`).

---

## 2. Gerbang Evaluasi Kelulusan (Phase Gate DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [ ] **`ROAD-04-GATE-001` (SPEC-04 Compliance = PASS):**
  - Invarian Default: YES terbukti bekerja (seluruh rule aktif 100% saat `charites.yaml` absen).
  - Presedensi rule terbukti: Kebijakan config (`off`) mengalahkan seleksi CLI `--rule`.
  - Builtin hard exclusions terbukti kebal dari negasi `!`.
  - Proteksi symlink dan batas ukuran berkas 10 MB terbukti menolak traversal yang berbahaya atau boros memori.
  - Direktif inline ignore terbukti mendukung tata bahasa multi-rule, wildcard `*`, serta cakupan baris dan AST node span.
  - Kontrak pembatalan interupsi (SIGINT/SIGTERM) menghentikan worker secara bersih dan membuang temuan parsial dengan exit code 130.

- [ ] **`ROAD-04-GATE-002` (ARCH-04 Compliance = PASS):**
  - Rule singleton di registry tetap murni dan immutable melalui enkapsulasi `ActiveRule`.
  - Matcher `.charitesignore` mengevaluasi pola secara sekuensial (*last matching rule wins*).
  - Worker pool beroperasi dengan prinsip isolasi mandiri (*share-nothing worker*).
  - Pengurutan temuan diagnostic menggunakan komparator pengurutan total 7 tingkat.

- [ ] **`ROAD-04-GATE-003` (TEST-04 Compliance = PASS):**
  - Seluruh skenario pengujian unit lolos 100% pada paket `internal/config`, `internal/scanner`, dan `internal/analyzer`.
  - Verifikasi pembatalan `context.Context` membuktikan tidak ada kebocoran goroutine (*goroutine leak*).
  - Pengujian urutan deterministik membuktikan pengurutan total bebas kolisi acak.

- [ ] **`ROAD-04-GATE-004` (QUAL-04 Compliance = PASS):**
  - Verifikasi bebas data race: `go test -race ./internal/...` lolos $100\%$ ($0$ race detected).
  - Seluruh ambang batas cakupan pengujian per-paket terpenuhi:
    - `internal/config`: $\ge 90\%$ line coverage.
    - `internal/scanner`: $\ge 85\%$ line coverage.
    - `internal/analyzer`: $\ge 90\%$ line coverage.
  - Kompleksitas siklomatik $\le 12$ per fungsi.
  - `golangci-lint run ./internal/...` lolos 100% tanpa isu.

---

## 3. Gerbang Transisi ke Fase 5 (CLI Entrypoint & Reporters)

Begitu keempat gerbang di atas berstatus **PASS**:
1. Buat git commit: `feat(engine): implement config resolver, parallel scanner, worker pool, and AST traversal engine`.
2. Melangkah ke Fase 5: Buka dokumen `docs/01-spec/05-cli.md` untuk merancang CLI entrypoint (`charites scan`, alias `check`/`run`), flag parsing (A-E), serta reporter ANSI terminal berwarna dan JSON payload formatter.
