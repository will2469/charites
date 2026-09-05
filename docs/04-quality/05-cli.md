# 04-QUALITY: 05 - CLI Reliability, Terminal Ergonomics & Presentation Invariants

> **Kode Dokumen:** `QUAL-05-CLI`
> **Tahapan:** Fase 5 - Reporter Output & CLI Entrypoint
> **Peran Pilar:** QUALITY = QUALITY THRESHOLD (Ambang Batas Kualitas, Keamanan Terminal & Anggaran Sumber Daya)
> **Status:** Ready / Approved for Implementation
> **Standar Rujukan:** POSIX.1-2017 CLI Standards & The 12-Factor App CLI Ergonomics

Dokumen ini mendefinisikan batasan kualitas untuk antarmuka CLI, penanganan pemisahan aliran data (*stream separation*), sanitasi karakter escape terminal, determinisme kode keluar, serta anggaran kinerja pelaporan.

---

## 1. Invarian Aliran Data & Keamanan Terminal (Terminal Ergonomics)

1. **Pemisahan Aliran Standar (Stream Separation):**
   - **`stdout`**: Hanya digunakan untuk mencetak dokumen hasil analisis (*inline report* atau *dokumen JSON*).
   - **`stderr`**: Hanya digunakan untuk mencetak pesan kesalahan sistem fatal, pesan validasi argumen CLI, atau pesan panduan bantuan (*usage/help*).
2. **Invarian Sanitasi Non-TTY (Piping & Redirection Safety):**
   Ketika `stdout` dialihkan ke pipa atau berkas (`charites scan . | grep ...` atau `charites scan . > report.txt`), CLI secara otomatis beralih ke mode `ColorNever`. Hal ini mencegah polusi karakter escape ANSI yang dapat merusak parsing script hilir.
3. **Kepatuhan Terhadap Standar `NO_COLOR`:**
   Sistem wajib menghormati variabel lingkungan `NO_COLOR` (jika `os.Getenv("NO_COLOR") != ""`, pewarnaan ANSI dinonaktifkan secara otomatis).
4. **Restorasi State Terminal pada Interupsi:**
   Jika proses menerima sinyal interupsi sistem (`SIGINT`/`Ctrl+C`), aplikasi wajib merestorasi status terminal (visibilitas kursor dan atribut warna) sebelum mengakhiri proses.

---

## 2. Invarian Integritas Exit Code

1. **No Silent Failures:** Dilarang keras mengembalikan exit code `0` apabila terdapat minimal satu temuan bertingkat `error`.
2. **Determinisme Kode Keluar:**
   - **`0`**: Sukses murni / bersih (atau hanya warning tanpa `--fail-on-warn`).
   - **`1`**: Kegagalan kualitas kode (ditemukan `error`, atau `warning` dengan `--fail-on-warn`).
   - **`2`**: Kegagalan operasional sistem / CLI (argumen salah, konflik category/rule, berkas target tidak dapat diakses, dll.).
   - **`130`**: Proses dihentikan oleh pengguna via sinyal terminal.
3. **Pemisahan Domain Pelanggaran:** Temuan diagnostik kode **DILARANG KERAS** menghasilkan exit code `2`.

---

## 3. Anggaran Performa Formatting Reporter (`QUAL-05-PERF-001`)

Untuk mencegah pembengkakan waktu pemrosesan laporan:

| Komponen Reporter | Metrik Indikator | Target Desain | Klasifikasi |
| :--- | :--- | :---: | :--- |
| **JSON Reporter Render** | Serialisasi 1.000 finding ke `io.Discard` | $\le 5\text{ ms}$ | Performance Baseline |
| **Inline Reporter Render** | Format teks 1.000 finding ke `io.Discard` | $\le 10\text{ ms}$ | Performance Baseline |
| **Alokasi Memori Reporter** | Alokasi heap per-finding | $\le 512\text{ B/finding}$ | Memory Budget |

Pengukuran wajib dilakukan secara terisolasi terhadap `io.Discard` untuk menghilangkan bias fluktuasi emulator terminal atau I/O disk.

---

## 4. Ambang Batas Kualitas & Metrik Kelulusan

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Line Coverage Paket CLI** | $\ge 85\%$ | `go test -cover ./internal/cli/...` |
| **Line Coverage Paket Reporter** | $\ge 90\%$ | `go test -cover ./internal/reporter/...` |
| **Verifikasi Data Race** | $0$ data race detected | `go test -race ./...` |
| **Batas Kompleksitas Siklomatik** | $\le 12$ per fungsi | `gocyclo -over 12 ./internal/cli ./internal/reporter` |
| **Subprocess E2E Tests** | 100% lulus | `go test -v ./tests/e2e/...` |
| **Deterministic Byte Equality** | 100% identik biner | `TestReporter_Determinism` |
| **Linter Compliance** | $0$ issues | `golangci-lint run ./internal/cli/... ./internal/reporter/...` |
| **Cross-Platform Compilation** | 100% build pass (4 target) | `make cross-compile` |

