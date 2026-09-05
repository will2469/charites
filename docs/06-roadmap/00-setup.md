# 06-ROADMAP: 00 - Phase 0 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-00-SETUP`
> **Tahapan:** Fase 0 - Inisialisasi Repositori & Toolchain
> **Status:** Ready for Execution

Dokumen ini mendefinisikan kriteria kelulusan (*exit criteria*) dan gerbang transisi (*transition gate*) untuk **Fase 0 (Setup)** sebelum tim diizinkan melangkah ke **Fase 1 (IR Contract)**.

---

## 1. Deliverables Wajib Fase 0

1. Berkas `go.mod` terinisialisasi dengan Go 1.26.
2. Seluruh direktori skeleton (`cmd/charites`, `internal/*`, `tests/*`) terbentuk.
3. Berkas konfigurasi `.charitesignore`, `.golangci.yml`, dan `Makefile` tersedia di root.
4. Entrypoint `cmd/charites/main.go` dan `internal/cli/root.go` dapat dieksekusi.

---

## 2. Checklist Definition of Done (DoD) Fase 0

- [ ] `make build` menghasilkan binary yang berfungsi di `bin/charites`.
- [ ] `./bin/charites --version` mencetak string versi `charites version 0.1.0-dev` dan exit code `0`.
- [ ] `./bin/charites --help` mencetak usage manual dan exit code `0`.
- [ ] `make test` lulus 100% pada unit test CLI dan E2E smoke test.
- [ ] `make lint` lulus tanpa satu pun warning dari `golangci-lint`.
- [ ] Verifikasi cross-compile (Linux, macOS, Windows) berhasil tanpa error CGO.

---

## 3. Gerbang Transisi ke Fase 1 (IR Contract)

Begitu seluruh checklist di atas tercentang hijau:
- Buat git commit inisial `feat(init): initialize repository skeleton, Go 1.26 toolchain and CLI entrypoint`.
- Buka dokumen [`docs/01-spec/01-contract.md`](../01-spec/01-contract.md) untuk memulai pendefinisian struktur data Intermediate Representation.
