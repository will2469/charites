# 04-QUALITY: 07 - MCP Protocol Integrity, Wiki Formatting & Security Invariants

> **Kode Dokumen:** `QUAL-07-MCP`
> **Tahapan:** Fase 7 - Ekosistem Lanjutan (MCP Server & Wiki Generator)
> **Status:** Ready for Review
> **Standar Rujukan:** Model Context Protocol Security Standards & POSIX Verification

Dokumen ini mendefinisikan batasan kualitas, penegakan integritas framing protokol MCP, determinisme generasi dokumentasi wiki, serta standar keamanan skrip instalasi shell.

---

## 1. Invarian Protokol MCP & Integritas Stdio

1. **Zero Stream Contamination Invariant:**
   `os.Stdout` dialokasikan secara eksklusif untuk aliran byte JSON-RPC 2.0. Satu karakter pun di luar framing JSON (seperti logging `slog`, `fmt.Print`, atau panic trace) **DILARANG KERAS** bocor ke `os.Stdout`. Seluruh log aplikasi wajib dialihkan ke `os.Stderr`.
2. **Kepatuhan JSON Schema:**
   Setiap parameter tool yang diekspos (`charites_scan`, `charites_explain_rule`, `charites_list_rules`) wajib lolos validasi skema JSON Schema resmi agar tidak memicu kegagalan komunikasi pada AI Agent host.
3. **Batas Waktu Eksekusi Tool (Execution Timeout):**
   Setiap pemanggilan `charites_scan` dibatasi batas waktu maksimum **30 detik** menggunakan `context.WithTimeout` untuk melindungi AI Agent dari proses yang membeku (*frozen process*).

---

## 2. Invarian Deterministik Wiki Generator

1. **Nol Pergeseran Git Diff (Zero Git Churn):**
   Generator markdown dilarang mencetak timestamp generasi (seperti `"Generated on 2026-09-05"`) di dalam berkas dokumentasi. Dua kali pengeksekusian `charites wiki` pada versi binary yang sama **PASTI** menghasilkan berkas markdown yang identik secara byte.
2. **Integritas Tautan Lokal (Relative Link Integrity):**
   Tautan markdown di `Home.md` wajib menggunakan relative link yang valid (misal `[Theme Rules](theme.md)`) tanpa tautan absolut disk lokal pengembang.

---

## 3. Invarian Keamanan Skrip Instalasi Shell (`scripts/install.sh`)

Untuk melindungi mesin pengguna dari manipulasi pihak ketiga (*supply-chain attacks*):
1. **Enkripsi HTTPS Mutlak:** Seluruh pengunduhan binary dan manifest rilis wajib melalui protokol HTTPS resmi (`https://github.com/will2469/charites/releases/...`).
2. **Verifikasi Checksum SHA256:**
   Skrip **WAJIB** mengunduh berkas `checksums.txt` dan memverifikasi integritas hash SHA256 dari tarball sebelum mengekstrak binary ke dalam `/usr/local/bin` atau `$HOME/.local/bin`. Jika checksum tidak cocok, skrip wajib membatalkan proses instalasi dengan kode keluar 1.
3. **POSIX-Compliant Shell:**
   Skrip wajib dapat dieksekusi oleh `/bin/sh` standar tanpa bergantung pada fitur non-standar Bash (*zero bashism*).

---

## 4. Ambang Batas Kualitas & Metrik Kelulusan

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Line Coverage Paket MCP** | $\ge 85\%$ | `go test -cover ./internal/mcp/...` |
| **Line Coverage Paket Wiki** | $\ge 85\%$ | `go test -cover ./internal/wiki/...` |
| **Linter Skrip Shell** | $0$ warning / error | `shellcheck -s sh scripts/install.sh` |
| **Linter Go Compliance** | $0$ issues | `golangci-lint run ./internal/mcp/... ./internal/wiki/...` |
