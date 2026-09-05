# 06-ROADMAP: 00 - Phase 0 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-00-SETUP`
> **Tahapan:** Fase 0 - Inisialisasi Repositori & Toolchain
> **Peran Pilar:** ROADMAP = PHASE GATE (Gerbang Evaluasi Transisi Antar-Fase)
> **Status:** Ready for Execution

Dokumen ini mendefinisikan gerbang evaluasi kelulusan (*phase gate*) untuk **Fase 0 (Setup)** sebelum proyek diizinkan melangkah ke **Fase 1 (IR Contract)**. Sesuai prinsip pemisahan otoritas arsitektur:
- **SPEC** = WHAT (Spesifikasi Kebutuhan & Kontrak Fungsional)
- **ARCH** = HOW (Rancangan Arsitektur, Boundaries & Metode)
- **TEST** = PROOF (Harness Pengujian & Asersi Pembuktian)
- **QUALITY** = QUALITY THRESHOLD (Ambang Batas Linter, Zero Dep & Hygiene)
- **ROADMAP** = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)

---

## 1. Deliverables Wajib Fase 0

1. **Berkas Modul Go:** `go.mod` terinisialisasi dengan Go `1.26` (toolchain `go1.26.x`) tanpa dependensi eksternal. Berkas `go.sum` **MUST NOT** be required.
2. **Implementasi Entrypoint & CLI:** `cmd/charites/main.go`, `internal/cli/root.go`, dan `internal/cli/version.go` aktif dan fungsional.
3. **Konfigurasi Lingkungan:** `.charitesignore`, `.golangci.yml`, `charites.example.yaml`, dan `Makefile` di root repositori.
4. **Harness Uji:** Unit test `internal/cli/root_test.go` dan subprocess E2E `tests/e2e/smoke_test.go`.
5. **Reservasi Direktori Arsitektur:** Direktori stub `internal/` dan `tests/` terbentuk rapi untuk isolasi fase berikutnya.

---

## 2. Checklist Definition of Done (Phase Gate DoD) Fase 0

- [ ] **SPEC (WHAT) Verification:**
  - [ ] `./bin/charites --version`, `./bin/charites -v`, dan `./bin/charites version` mencetak `charites version 0.1.0-dev` dan exit `0`.
  - [ ] `./bin/charites --help` dan `./bin/charites -h` mencetak usage manual dan exit `0`.
  - [ ] `./bin/charites` (tanpa argumen) mencetak usage manual dan exit `0`.
  - [ ] `./bin/charites unknown-command` mencetak error ke `stderr` dan exit `2`.
  - [ ] `.charitesignore` ada di root dan mematuhi sintaks Gitignore-compatible.
- [ ] **ARCH (HOW) Verification:**
  - [ ] Binary dikompilasi dengan `CGO_ENABLED=0` secara native.
  - [ ] `cmd/charites/main.go` murni bertindak sebagai trampoline ke `cli.Execute()`.
  - [ ] Struktur paket mematuhi enkapsulasi `internal/` tanpa kebocoran logika fase masa depan.
- [ ] **TEST (PROOF) Verification:**
  - [ ] `make test` lolos 100% pada unit test CLI dan subprocess E2E smoke test (`go test -v -race ./...`).
- [ ] **QUALITY (THRESHOLD) Verification:**
  - [ ] `make lint` (`golangci-lint run ./...`) lolos dengan exit code `0` tanpa peringatan pada seluruh kode Fase 0.
  - [ ] `go.mod` bebas dependensi pihak ketiga (*zero third-party dependencies*).
- [ ] **Cross-Platform Compilation Gate:**
  - [ ] Verifikasi kompilasi silang (*cross-compile*) berhasil tanpa CGO untuk Linux (`amd64`, `arm64`), macOS Apple Silicon (`darwin/arm64`), dan Windows (`amd64`).

---

## 3. Gerbang Transisi ke Fase 1 (IR Contract)

Begitu seluruh kriteria evaluasi gerbang di atas tercentang hijau:
1. Rekam checkpoint git commit: `feat(init): complete Phase 0 repository setup, CLI entrypoint and toolchain baseline`.
2. Buka dokumen [docs/01-spec/01-contract.md](file:///home/will/Monorepo/charites/docs/01-spec/01-contract.md) untuk memulai implementasi kontrak data Intermediate Representation (**Fase 1**).

