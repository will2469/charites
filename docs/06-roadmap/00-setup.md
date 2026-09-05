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

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [ ] **`ROAD-00-GATE-001` (SPEC-00 Compliance = PASS):**
  - [ ] Flag `--version`, `-v`, dan subcommand `version` mencetak info ke `stdout` (bersih dari `stderr`) dan keluar dengan exit code `0`.
  - [ ] Flag `--help`, `-h`, dan pemanggilan tanpa argumen mencetak panduan usage ke `stdout` (bersih dari `stderr`) dan keluar dengan exit code `0`.
  - [ ] Subcommand/flag tidak dikenal (`unknown-command`, `--bogus`) mencetak pesan error ke `stderr` (bersih dari `stdout`) dan keluar dengan exit code `2`.
  - [ ] Berkas `.charitesignore` tersedia di root dan memuat seluruh default pattern ignorasi Gitignore-compatible.
  - [ ] Berkas `go.sum` **MUST NOT** be required pada Fase 0 saat dependensi eksternal bernilai nol.

- [ ] **`ROAD-00-GATE-002` (ARCH-00 Compliance = PASS):**
  - [ ] Rantai dependensi Makefile `all: build test lint` terbukti berhasil dieksekusi secara deterministik pada *fresh checkout*.
  - [ ] Entrypoint `cmd/charites/main.go` murni bertindak sebagai trampoline ke `cli.Execute()`.
  - [ ] Struktur direktori mematuhi enkapsulasi `internal/` dan pemisahan *skeleton directory reservations* tanpa kebocoran logika fase masa depan.
  - [ ] Implementasi build secara native menggunakan `CGO_ENABLED=0` (`SPEC-00-BUILD-001`).

- [ ] **`ROAD-00-GATE-003` (TEST-00 Compliance = PASS):**
  - [ ] Unit test `internal/cli/root_test.go` lolos 100% (`go test -v -race ./...`).
  - [ ] Subprocess E2E smoke test `tests/e2e/smoke_test.go` lolos dengan pembuktian pemisahan buffer `stdout` vs `stderr`.
  - [ ] Prosedur verifikasi kompilasi silang `TEST-00-BUILD-002` lolos tanpa CGO untuk 4 platform resmi: `linux/amd64`, `linux/arm64`, `darwin/arm64`, dan `windows/amd64`.

- [ ] **`ROAD-00-GATE-004` (QUAL-00 Compliance = PASS):**
  - [ ] `golangci-lint run ./...` lolos dengan exit code `0` tanpa peringatan pada seluruh kode Fase 0.
  - [ ] `go.mod` bebas dependensi pihak ketiga (*zero third-party dependencies*).
  - [ ] Invarian Zero CGO terverifikasi (bebas dari berkas `.c`, `.h`, atau blok `import "C"`).
  - [ ] Konfigurasi pipeline CI memin versi Go pada baseline reproduktifitas `go1.26.0`.

---

## 3. Gerbang Transisi ke Fase 1 (IR Contract)

Begitu seluruh kriteria evaluasi gerbang `ROAD-00-GATE-001` s/d `ROAD-00-GATE-004` di atas tercentang hijau:
1. Rekam checkpoint git commit: `feat(init): complete Phase 0 repository setup, CLI entrypoint and toolchain baseline`.
2. Buka dokumen [docs/01-spec/01-contract.md](file:///home/will/Monorepo/charites/docs/01-spec/01-contract.md) untuk memulai implementasi kontrak data Intermediate Representation (**Fase 1**).

