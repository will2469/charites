# 06-ROADMAP: 05 - Phase 5 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-05-CLI`
> **Tahapan:** Fase 5 - Reporter Output & CLI Entrypoint
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Ready / Approved for Implementation

Dokumen ini menetapkan kriteria kelulusan (*exit criteria*) dan gerbang transisi (*phase gate*) untuk **Fase 5 (Reporter Output & CLI Entrypoint)** sebelum tim diizinkan melangkah ke **Fase 6 (Validasi Penuh & Golden Tests)**. Sesuai prinsip tata kelola [docs/00-CONTRACT.md](https://github.com/will2469/charites/blob/main/docs/00-CONTRACT.md):
- **SPEC** = WHAT (Spesifikasi Antarmuka CLI, Format Reporter & Kontrak Exit Code)
- **ARCH** = HOW (Rancangan Arsitektur CLI, Dispatcher & Abstraksi Reporter)
- **TEST** = PROOF (Harness Pengujian, Matriks E2E & Snapshot Reporter)
- **QUALITY** = QUALITY THRESHOLD (Ambang Batas Kualitas, Keamanan Terminal & Anggaran Sumber Daya)
- **ROADMAP** = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)

---

## 1. Deliverables Berkas Fase 5

1. **`internal/reporter/reporter.go`**: Definisi interface `Reporter` dan struktur data transfer `ScanResult` serta `ScanSummary` (dengan field `version` dan `duration_ms`).
2. **`internal/reporter/inline.go`**: Presenter terminal berbasis ANSI dengan path POSIX, penanganan clean scan banner, dan kepatuhan `NO_COLOR`.
3. **`internal/reporter/json.go`**: Presenter dokumen tunggal JSON lengkap yang mematuhi skema deterministik biner.
4. **`internal/reporter/color.go`**: Resolver `ColorMode` portabel (`ColorAuto`, `ColorNever`) berbasis flag, `NO_COLOR`, dan TTY.
5. **`internal/cli/root.go`**: Router utama CLI yang menangani delegasi subcommand, pemanggilan 0 argumen ke `scan .`, dan alias (`check`, `run`).
6. **`internal/cli/scan.go`**: Controller subcommand `scan` yang mengurai dan memvalidasi flag ergonomi (`--ext`, `--category`, `--rule`, `--ignore`, `--fail-on-warn`).
7. **`internal/cli/exit.go`**: Modul resolusi exit code POSIX terpadu (0, 1, 2, 130).
8. **`cmd/charites/main.go`**: Berkas titik masuk (*entrypoint*) binary yang memanggil `cli.Execute(os.Args[1:])`.
9. **`tests/golden/reporters/`**: Berkas referensi snapshot reporter (`inline_clean.golden`, `inline_violations.golden`, `inline_no_color.golden`, `json_clean.golden`, `json_violations.golden`).
10. **`tests/e2e/cli_test.go`**: Suite pengujian terintegrasi subprocess E2E yang memvalidasi seluruh matriks interaksi CLI.

---

## 2. Gerbang Kesiapan Implementasi (Readiness Gate DoD)

Sebelum implementasi kode Fase 5 dimulai, seluruh pilar arsitektur wajib menyelesaikan rantai gerbang kesiapan (*readiness gate chain*):

$$\mathbf{SPEC\ PASS} \longrightarrow \mathbf{ARCH\ PASS} \longrightarrow \mathbf{TEST\ PASS} \longrightarrow \mathbf{QUALITY\ PASS} \longrightarrow \mathbf{PHASE\ 5\ IMPLEMENTATION\ READY}$$

- [x] **`ROAD-05-READINESS-001` (SPEC-05 Compliance = PASS):**
  - Kontrak CLI dibekukan tanpa kebutuhan implisit (default command `scan .`, alias `check`/`run`).
  - Normalisasi `--ext` (case-insensitive, optional leading dot, comma, repeated) dan penolakan nilai tidak sah (exit 2).
  - Validasi keberadaan & irisan `--category` $\times$ `--rule` (exit 2 jika konflik).
  - Presedensi `--ignore` & kepatuhan invarian builtin hard exclusions (kekebalan negasi `!`).
  - Kontrak format `inline` ANSI dan `json` tunggal lengkap.
  - Resolusi TTY dan `NO_COLOR` (Mode `ColorNever`).
  - Total ordering 7-tingkat untuk determinisme output.
  - Taksonomi exit code 0, 1, 2, 130 terkunci.
  - Kontrak Zero Residual Footprint & Clean Lifecycle (SPEC-05-LIFECYCLE-001) terkunci.

- [x] **`ROAD-05-READINESS-002` (ARCH-05 Compliance = PASS):**
  - Matriks kepemilikan komponen terpetakan secara ketat (`internal/cli/root.go`, `internal/cli/scan.go`, `internal/config`, `internal/scanner`, `internal/analyzer`, `internal/reporter`, `internal/cli/exit.go`).
  - ARCH hanya menjelaskan HOW tanpa menambah behavior baru di luar SPEC.
  - Zero External Dependencies invariant dipertahankan.
  - Realisasi arsitektur bebas residu (stateless memory, ephemeral lifecycle, hermetic host isolation) terpetakan.

- [x] **`ROAD-05-READINESS-003` (TEST-05 Compliance = PASS):**
  - Matriks 14 skenario eksekutabel didefinisikan secara lengkap (termasuk audit bebas residu host).
  - Golden snapshot contracts (`inline_clean`, `inline_violations`, `inline_no_color`, `json_clean`, `json_violations`) terspesifikasi.
  - Pengujian determinisme biner dan E2E subprocess exit code matrix terspesifikasi.

- [x] **`ROAD-05-READINESS-004` (QUALITY-05 Compliance = PASS):**
  - Ambang batas eksplisit dikunci: `internal/cli` $\ge 85\%$, `internal/reporter` $\ge 90\%$.
  - 0 Data Race Invariant, $\le 12$ cyclomatic complexity per fungsi, linter compliance 0 issues.
  - Invarian kebersihan host dan 0 proses tertinggal terkunci (`QUAL-05-HOST-CLEANLINESS`).
  - CI Baseline terverifikasi pada commit final Fase 4.

- [x] **`ROAD-05-READINESS-005` (CI Evidence Verification = PASS):**
  - Hasil CI lokal: `make all` PASS (Build + Race Test + Lint + Gofmt).
  - Cross-compilation 4 target: `make cross-compile` PASS (`linux/amd64`, `linux/arm64`, `darwin/arm64`, `windows/amd64`).
  - Bukti eksekusi terlampir pada commit `12d4ed2`.

---

## 3. Gerbang Evaluasi Kelulusan Fase (Exit Criteria DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [ ] **`ROAD-05-GATE-001` (SPEC-05 Compliance = PASS):**
  - Pemanggilan 0 argumen (`charites`) otomatis menjalankan `scan .`.
  - Subcommand alias `check` dan `run` bekerja identik dengan `scan`.
  - Normalisasi `--ext` mendukung format case-insensitive dan menolak ekstensi yang tidak didukung dengan exit code `2`.
  - Konflik `--category` dan `--rule` terdeteksi dan menghasilkan exit code `2`.
  - Format `--format=inline` dan `--format=json` memenuhi skema dokumen lengkap.
  - Taksonomi exit code dipatuhi (0 = Clean, 1 = Violations, 2 = Operational error, 130 = Interrupted).

- [ ] **`ROAD-05-GATE-002` (ARCH-05 Compliance = PASS):**
  - Binary mempertahankan Zero Dependency Invariant (hanya menggunakan pustaka standar Go).
  - Skema DTO `ScanResult` memuat field `version` dan `duration_ms` secara presisi.
  - Abstraksi `ColorMode` memungkinkan isolasi pengujian terminal tanpa ketergantungan device OS.

- [ ] **`ROAD-05-GATE-003` (TEST-05 Compliance = PASS):**
  - Pengujian determinisme membuktikan output reporter 100% identik secara biner (`bytes.Equal`).
  - Golden contract snapshots (`tests/golden/reporters/`) terverifikasi 100%.
  - Seluruh skenario matriks E2E subprocess (termasuk kasus `--fail-on-warn`) lulus tanpa kegagalan.

- [ ] **`ROAD-05-GATE-004` (QUAL-05 Compliance = PASS):**
  - Seluruh ambang batas cakupan pengujian paket terpenuhi:
    - `internal/cli`: $\ge 85\%$ line coverage.
    - `internal/reporter`: $\ge 90\%$ line coverage.
  - Anggaran performa rendering formatting terisolasi terpenuhi (`QUAL-05-PERF-001`).
  - `golangci-lint run ./...` lulus 100% tanpa isu.

---

## 4. Gerbang Transisi ke Fase 6 (Validasi Penuh & Golden Tests)

Begitu keempat gerbang di atas berstatus **PASS**:
1. Buat git commit: `feat(cli): implement CLI entrypoint, ergonomic flags A-E, ANSI and JSON reporters, and E2E test suite`.
2. Melangkah ke Fase 6: Buka dokumen `docs/01-spec/06-golden.md` untuk merancang penutupan pipa (*milestone selesai pipa*), mencakup uji regresi *golden snapshots* komprehensif (`tests/golden/*`), fixtures lintas framework, dan verifikasi akhir sebelum binary dinyatakan stabil untuk ekspansi rule massal.
