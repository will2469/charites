# 06-ROADMAP: 04 - Phase 4 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-04-ENGINE`
> **Tahapan:** Fase 4 - Konfigurasi, Concurrency Scanner & Traversal Engine
> **Status:** Ready for Review

Dokumen ini menetapkan deliverable wajib, kriteria kelulusan (*exit criteria*), dan gerbang transisi (*transition gate*) untuk **Fase 4 (Konfigurasi, Concurrency Scanner & Traversal Engine)** sebelum melangkah ke **Fase 5 (Reporter Output & CLI Entrypoint)**.

---

## 1. Deliverables Wajib Fase 4

1. **`internal/config/config.go`**: Parser `charites.yaml` dengan penegakan prinsip **Default: YES** (Model Argus) serta resolusi override status dan severity rule.
2. **`internal/config/ignore.go`**: Matcher pola `.charitesignore` yang mendukung wildcard glob, negasi `!`, dan fungsi *early directory pruning* instan.
3. **`internal/scanner/walker.go`**: Traversal direktori cepat yang mendukung target berkas langsung (Ergonomi A) dan penyaringan ekstensi berkas (Ergonomi B).
4. **`internal/scanner/pool.go`**: Concurrency worker pool berbasis goroutine dengan kapasitas default `runtime.NumCPU()`.
5. **`internal/analyzer/context.go`**: Context terisolasi per-berkas dengan parser direktif penekanan inline ignore (`// charites:ignore` dan `<!-- charites:ignore -->`).
6. **`internal/analyzer/engine.go`**: Traversal engine pohon AST berbasis iterator Go 1.26 (`root.Walk()`) dan modul pengurutan deterministik `(File, Line, Column, Rule)`.
7. **Suite Pengujian Konkurensi**: Unit test dan race detection test untuk seluruh paket Fase 4 (`config_test.go`, `ignore_test.go`, `walker_test.go`, `pool_test.go`, `engine_test.go`).

---

## 2. Checklist Definition of Done (DoD) Fase 4

- [ ] Invarian Default: YES terbukti bekerja (seluruh rule aktif 100% saat `charites.yaml` absen).
- [ ] Override penonaktifan rule (`off`/`false`) dan perubahan severity (`warn`/`error`) terbukti efektif.
- [ ] Early directory pruning terbukti mencegah pembacaan folder `node_modules/`, `dist/`, dan `.git/`.
- [ ] Worker pool terbukti mampu memproses ribuan berkas secara konkuren tanpa goroutine leak atau deadlock.
- [ ] Pengujian `go test -race ./internal/config/... ./internal/scanner/... ./internal/analyzer/...` lulus 100% tanpa satupun data race.
- [ ] Direktif inline ignore terbukti mampu menyaring temuan diagnostic pada baris yang sama maupun baris tepat di bawahnya.
- [ ] Hasil pelaporan diagnostic terbukti selalu terurut deterministik.
- [ ] Total test coverage paket Fase 4 mencapai minimal $85\%$.
- [ ] `golangci-lint run ./internal/...` bersih 100%.

---

## 3. Gerbang Transisi ke Fase 5 (CLI & Reporters)

Begitu seluruh kriteria DoD di atas terpenuhi:
1. Buat git commit: `feat(engine): implement config resolver, parallel scanner, worker pool, and AST traversal engine`.
2. Melangkah ke Fase 5: Buka dokumen `docs/01-spec/05-cli.md` untuk merancang CLI entrypoint (`charites scan`, alias `check`/`run`), flag parsing (A-E), serta reporter ANSI terminal berwarna dan JSON payload formatter.
