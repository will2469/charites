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
  - [ ] Kontrak CLI Entrypoint (Version, Usage, Unknown Commands, Exit Codes 0 & 2) terpenuhi.
  - [ ] Kontrak Stream Routing (Stdout vs Stderr terisolasi penuh) terpenuhi.
  - [ ] Kontrak Ignorasi Gitignore-compatible (`.charitesignore`) terpenuhi.
  - [ ] Kontrak Modul Go & Zero External Dependencies (`go.sum` not required) terpenuhi.
  - [ ] Kontrak Kompilasi `SPEC-00-BUILD-001` (Zero-CGO) & `SPEC-00-BUILD-002` (4 Target Resmi) terpenuhi.

- [ ] **`ROAD-00-GATE-002` (ARCH-00 Compliance = PASS):**
  - [ ] Enkapsulasi `internal/` terpenuhi tanpa kebocoran kontrak publik.
  - [ ] Pola entrypoint trampoline `cmd/charites/main.go` terpenuhi.
  - [ ] Isolasi fase skeleton directory reservations terpenuhi tanpa kebocoran logika fase masa depan.
  - [ ] Automasi Makefile mengimplementasikan urutan dependensi deterministik (`ARCH-00-BUILD-ORDER`).

- [ ] **`ROAD-00-GATE-003` (TEST-00 Verification = PASS):**
  - [ ] `TEST-00-CLI-001` (Unit Test CLI Dispatcher) = PASS.
  - [ ] `TEST-00-SMOKE-001` (Subprocess E2E Smoke & Stream Routing) = PASS.
  - [ ] `TEST-00-BUILD-001` (Zero-CGO Binary Compilation) = PASS.
  - [ ] `TEST-00-BUILD-002` (Cross-Platform 4 Target Compilation) = PASS.
  - [ ] `TEST-00-BUILD-ORDER` (Fresh Checkout Build Before Test Verification) = PASS.

- [ ] **`ROAD-00-GATE-004` (QUAL-00 Threshold = PASS):**
  - [ ] `QUAL-00-LINT-001` (Baseline 9 Linters + Govet Enable-All) = PASS.
  - [ ] `QUAL-00-DEPS-001` (Zero Third-Party Dependencies Policy) = PASS.
  - [ ] `QUAL-00-INVAR-001` (Zero CGO Repository Surface Invariant) = PASS.
  - [ ] `QUAL-00-REPRO-001` (GitHub Actions CI Pinned Toolchain `go1.26.0`) = PASS.

---

## 3. Gerbang Transisi ke Fase 1 (IR Contract)

Begitu seluruh kriteria evaluasi gerbang `ROAD-00-GATE-001` s/d `ROAD-00-GATE-004` di atas tercentang hijau:
1. Rekam checkpoint git commit: `feat(init): complete Phase 0 repository setup, CLI entrypoint and toolchain baseline`.
2. Buka dokumen [docs/01-spec/01-contract.md](file:///home/will/Monorepo/charites/docs/01-spec/01-contract.md) untuk memulai implementasi kontrak data Intermediate Representation (**Fase 1**).

