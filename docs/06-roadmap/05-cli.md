# 06-ROADMAP: 05 - Phase 5 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-05-CLI`
> **Tahapan:** Fase 5 - Reporter Output & CLI Entrypoint
> **Status:** Ready for Review

Dokumen ini menetapkan deliverable wajib, kriteria kelulusan (*exit criteria*), dan gerbang transisi (*transition gate*) untuk **Fase 5 (Reporter Output & CLI Entrypoint)** sebelum melangkah ke **Fase 6 (Validasi Penuh & Golden Snapshots)**.

---

## 1. Deliverables Wajib Fase 5

1. **`internal/reporter/reporter.go`**: Definisi interface `Reporter` dan struktur data transfer `ScanResult` serta `ScanSummary`.
2. **`internal/reporter/inline.go`**: Presenter terminal berbasis ANSI dengan deteksi TTY otomatis dan kepatuhan terhadap `NO_COLOR`.
3. **`internal/reporter/inline_test.go`**: Unit test pemformatan teks ANSI, peringatan, error, hint, dan mode no-color.
4. **`internal/reporter/json.go`**: Presenter streaming JSON yang efisien dan mematuhi skema mesin.
5. **`internal/reporter/json_test.go`**: Unit test serialisasi dan deserialisasi validitas payload JSON.
6. **`internal/cli/root.go`**: Router utama CLI yang menangani delegasi subcommand dan pemetaan otomatis argumen default ke `scan`.
7. **`internal/cli/scan.go`**: Controller subcommand `scan` yang mengurai flag ergonomi (A-E: direct file target, `--ext`, `--category`, `--rule`, `--format`).
8. **`cmd/charites/main.go`**: Berkas titik masuk (*entrypoint*) binary yang memanggil `cli.Execute(os.Args[1:])` dan mengalirkan exit code ke `os.Exit()`.
9. **`tests/e2e/cli_test.go`**: Suite pengujian subprocess E2E yang memvalidasi interaksi terminal dan kontrak exit code 0, 1, dan 2.

---

## 2. Checklist Definition of Done (DoD) Fase 5

- [ ] Binary `charites` dapat dikompilasi secara bersih (`go build ./cmd/charites`).
- [ ] Eksekusi `charites .` atau `charites src/` otomatis menjalankan pemindaian (default command `scan`).
- [ ] Subcommand alias `charites check` dan `charites run` bekerja identik dengan `charites scan`.
- [ ] Opsi ergonomi pemindaian (A-E) teruji dan bekerja sesuai spesifikasi:
  - Target berkas langsung tanpa traversal folder.
  - Filter `--ext` menyaring ekstensi dengan benar.
  - Filter `--category` dan `--rule` membatasi rule yang dieksekusi.
- [ ] Format `--format=inline` mencetak teks ANSI berwarna yang rapi, dan otomatis mematikan warna saat dialihkan ke pipe (`|`) atau saat `NO_COLOR=1`.
- [ ] Format `--format=json` menghasilkan JSON valid yang dapat diuraikan oleh `jq`.
- [ ] Kontrak exit code POSIX terbukti di pengujian E2E:
  - Return `0` jika repositori bersih.
  - Return `1` jika terdeteksi pelanggaran `error`.
  - Return `2` jika terdapat kesalahan argumen fatal atau path tidak ada.
- [ ] Cakupan pengujian (*test coverage*) paket `internal/cli/...` dan `internal/reporter/...` mencapai minimal $85\%$.
- [ ] `golangci-lint run ./...` bersih 100%.

---

## 3. Gerbang Transisi ke Fase 6 (Validasi Penuh & Golden Tests)

Begitu seluruh kriteria DoD di atas terpenuhi:
1. Buat git commit: `feat(cli): implement CLI entrypoint, ergonomic flags A-E, ANSI and JSON reporters, and E2E test suite`.
2. Melangkah ke Fase 6: Buka dokumen `docs/01-spec/06-golden.md` untuk merancang penutupan pipa (*milestone selesai pipa*), mencakup uji regresi *golden snapshots* (`tests/golden/*`), fixtures komprehensif, dan verifikasi akhir sebelum binary dinyatakan stabil untuk ekspansi rule massal.
