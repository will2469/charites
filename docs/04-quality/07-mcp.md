# 04-QUALITY: 07 - MCP Protocol Integrity, Wiki Formatting & Security Invariants

> **Kode Dokumen:** `QUAL-07-MCP`
> **Tahapan:** Fase 7 - Ekosistem Lanjutan (MCP Server, Wiki Generator & Secure Installer)
> **Peran Pilar:** QUALITY = QUALITY THRESHOLD (Ambang Batas Kualitas, Integritas Protokol & Keamanan Pasokan)
> **Status:** Approved / Completed
> **Standar Rujukan:** Model Context Protocol Security Standards & POSIX Shell Guidelines

Dokumen ini mendefinisikan batasan kualitas, integritas kanal protokol MCP, determinisme generasi dokumentasi wiki, serta standar keamanan skrip instalasi shell.

---

## 1. Invarian Protokol MCP & Integritas Stdio

1. **Zero Stream Contamination Invariant:**
   `os.Stdout` dialokasikan secara eksklusif untuk aliran biner pesan JSON-RPC 2.0. Satu karakter pun di luar framing JSON (seperti logging, banner teks, atau panic trace) **DILARANG KERAS** bocor ke `os.Stdout`. Seluruh log diagnostik wajib dialihkan ke `os.Stderr`.
2. **Batas Ukuran Frame Protokol (Frame Size Bound):**
   Ukuran satu frame pesan JSON-RPC dibatasi maksimal **4 Megabytes** ($4 \times 1024 \times 1024\text{ bytes}$) untuk mencegah serangan *Denial of Service* (DoS) atau kehabisan memori (*Out-of-Memory*).
3. **Batas Waktu Eksekusi Tool (Execution Timeout):**
   Setiap pemanggilan `charites_scan` dibatasi batas waktu maksimum **30 detik** menggunakan `context.WithTimeout` untuk melindungi AI Host dari proses yang membeku (*hung execution*).
4. **Batas Keamanan Direktori (Trust Boundary Invariant):**
   Pemindaian MCP tidak boleh mengeksekusi path di luar direktori workspace yang sah (*path traversal defense*).

---

## 2. Invarian Determinis Wiki Generator

1. **Nol Pergeseran Git Diff (Zero Git Churn):**
   Generator markdown dilarang mencetak timestamp generasi (seperti `"Generated on 2026-09-05"`) di dalam berkas dokumentasi. Dua kali pengeksekusian `charites wiki` pada versi binary yang sama **PASTI** menghasilkan berkas markdown yang identik secara biner.
2. **Determinis Pengurutan (Total Ordering Invariant):**
   Daftar kategori diurutkan leksikografis menaik, dan rule di dalam setiap kategori diurutkan menaik berdasarkan `Rule.ID()`.
3. **Integritas Tautan Relatif:**
   Tautan navigasi di `Home.md` wajib menggunakan relative link yang valid tanpa path absolut disk lokal pengembang.

---

## 3. Invarian Keamanan Skrip Instalasi Shell (`scripts/install.sh`)

1. **Enkripsi HTTPS Mutlak:** Seluruh pengunduhan binary dan manifest rilis wajib melalui protokol HTTPS resmi (`https://github.com/will2469/charites/releases/...`).
2. **Verifikasi Checksum SHA256:**
   Skrip **WAJIB** mengunduh berkas `checksums.txt` dan memverifikasi hash SHA256 tarball sebelum ekstraksi. Jika hash tidak cocok, instalasi langsung dibatalkan dengan kode keluar 1.
3. **Pencegahan Path Traversal pada Arsip:**
   Arsip tarball diverifikasi bebas dari entri relatif berbahaya (`../`) sebelum dipindahkan ke direktori bin tujuan.
4. **POSIX-Compliant Shell:**
   Skrip wajib dapat dieksekusi oleh `/bin/sh` standar tanpa bergantung pada fitur non-standar Bash (*zero bashism*).

---

## 4. Ambang Batas Kualitas & Metrik Kelulusan

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Line Coverage Paket MCP** | $\ge 85\%$ | `go test -cover ./internal/mcp/...` |
| **Line Coverage Paket Wiki** | $\ge 85\%$ | `go test -cover ./internal/wiki/...` |
| **Linter Skrip Shell** | $0$ warning / error | `shellcheck -s sh scripts/install.sh` |
| **Linter Go Compliance** | $0$ issues | `golangci-lint run ./internal/mcp/... ./internal/wiki/...` |
