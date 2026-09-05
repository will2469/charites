# 04-QUALITY: 05 - CLI Reliability, Terminal Ergonomics & Presentation Invariants

> **Kode Dokumen:** `QUAL-05-CLI`
> **Tahapan:** Fase 5 - Reporter Output & CLI Entrypoint
> **Status:** Ready for Review
> **Standar Rujukan:** POSIX.1-2017 CLI Standards & The 12-Factor App CLI Ergonomics

Dokumen ini mendefinisikan batasan kualitas untuk antarmuka CLI, penanganan aliran data (*stream separation*), keamanan karakter escape terminal, serta invarian integritas *exit code*.

---

## 1. Invarian Aliran Data & Keamanan Terminal (Terminal Ergonomics)

1. **Pemisahan Aliran Standar (Stream Separation):**
   - **`stdout`**: Hanya digunakan untuk mencetak hasil pelaporan analisis (*inline report* atau *JSON payload*).
   - **`stderr`**: Hanya digunakan untuk mencetak pesan kesalahan sistem fatal, peringatan kegagalan parsing file konfigurasi, atau pesan bantuan (*usage/help*).
2. **Invarian Sanitasi Non-TTY (Piping & Redirection Safety):**
   Ketika `stdout` dialihkan ke pipa atau berkas (`charites scan . | grep ...` atau `charites scan . > report.txt`), CLI secara otomatis menonaktifkan kode warna ANSI escape sequence. Hal ini mencegah polusi karakter kontrol yang dapat merusak parsing script hilir (*downstream script*).
3. **Kepatuhan Terhadap Standar `NO_COLOR`:**
   Sistem wajib menghormati inisiatif industri `NO_COLOR` (jika `os.Getenv("NO_COLOR") != ""`, pewarnaan ANSI dinonaktifkan tanpa perlu flag manual).
4. **Restorasi State Terminal pada Interupsi:**
   Jika proses menerima sinyal interupsi sistem (`SIGINT`/`Ctrl+C`), aplikasi wajib merestorasi status terminal (termasuk visibilitas kursor dan atribut warna) sebelum keluar.

---

## 2. Invarian Integritas Exit Code

Exit code adalah kontrak deterministik mesin dengan ekosistem CI/CD:

1. **No Silent Failures:** Dilarang keras mengembalikan exit code `0` apabila terdapat minimal satu temuan bertingkat `error`.
2. **Determinisme Kode Keluar:**
   - **`0`**: Sukses murni (bersih).
   - **`1`**: Kegagalan kualitas kode (ditemukan pelanggaran).
   - **`2`**: Kegagalan operasional sistem (path tidak ditemukan, argumen tidak dikenal, dll.).
3. Nilai exit code di luar rentang `[0, 2]` dilarang digunakan untuk menjaga portabilitas antar sistem operasi (Linux/macOS/Windows).

---

## 3. Invarian Efisiensi Alokasi Memori Reporter

1. **Streaming JSON Encoding:**
   Presenter JSON dilarang melakukan konversi seluruh objek diagnostik ke buffer string besar di memori (`json.Marshal`). Reporter wajib menggunakan `json.NewEncoder(io.Writer)` untuk mengalirkan byte secara langsung ke output stream.
2. **Nol Pembengkakan Waktu (Sub-5ms Rendering):**
   Waktu pemformatan dan pencetakan hasil laporan (baik inline ANSI maupun JSON) untuk 10.000 temuan harus selesai dalam waktu $\le 5\text{ milidetik}$.

---

## 4. Ambang Batas Kualitas & Metrik Kelulusan

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Line Coverage Paket CLI** | $\ge 85\%$ | `go test -cover ./internal/cli/...` |
| **Line Coverage Paket Reporter** | $\ge 90\%$ | `go test -cover ./internal/reporter/...` |
| **Subprocess E2E Tests** | 100% lulus | `go test -v ./tests/e2e/...` |
| **Linter Compliance** | $0$ issues | `golangci-lint run ./internal/cli/... ./internal/reporter/...` |
