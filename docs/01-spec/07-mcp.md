# 01-SPEC: 07 - MCP Server, Wiki Generator & Installer Specification

> **Kode Dokumen:** `SPEC-07-MCP`
> **Tahapan:** Fase 7 - Ekosistem Lanjutan (MCP Server, Wiki Generator & Secure Installer)
> **Peran Pilar:** SPEC = WHAT (Spesifikasi Protokol MCP, Generator Wiki & Keamanan Pasokan)
> **Status:** Ready for Review
> **Standar Rujukan:** Model Context Protocol (MCP) JSON-RPC 2.0 Specification / POSIX Shell Standards

Dokumen ini mendefinisikan spesifikasi kebutuhan fungsional untuk tiga sub-ekosistem Charites:
1. **Server Model Context Protocol (MCP)** berbasis Stdio (`charites mcp`)
2. **Generator Dokumentasi Ensiklopedia** (`charites wiki`)
3. **Skrip Pemasang Mandiri Aman** (*Secure One-Line Installer*)

---

## BAGIAN A: Server Model Context Protocol (MCP) (`internal/mcp/`)

Charites menyediakan antarmuka terstandarisasi untuk berinteraksi langsung dengan AI Agent (seperti Cursor, Claude Desktop, Antigravity) melalui protokol MCP:

### 1. Transport & Framing Protokol
- **Subcommand:** `charites mcp [--workspace=<path>]`
- **Transport Streams:**
  - `stdin`: Aliran masuk eksklusif untuk request dan notification JSON-RPC dari client.
  - `stdout`: Aliran keluar eksklusif untuk response dan notification JSON-RPC ke client. **Invarian Mutlak:** Dilarang mencetak log diagnostik, banner, atau teks non-protokol apa pun ke `stdout`.
  - `stderr`: Aliran eksklusif untuk logging diagnostik internal, debugging, dan error traces.
- **Framing & Encoding:**
  - Format pesan: UTF-8 tunggal per baris (*single-line JSON-RPC*), diakhiri karakter LF (`\n`, karakter CR `\r` otomatis dipangkas saat ingestion).
  - **Batas Maksimum Ukuran Frame (Protocol Safety Limit):** Maksimal **4 Megabytes** ($4 \times 1024 \times 1024\text{ bytes}$). Pesan yang melebihi batas ini langsung ditolak dengan JSON-RPC Parse Error (`-32700`).

### 2. State Machine Siklus Hidup Protokol (Lifecycle State Machine)
Server MCP mengimplementasikan *finite state machine* ketat:

$$\text{NEW} \xrightarrow{\text{initialize}} \text{INITIALIZING} \xrightarrow{\text{notifications/initialized}} \text{READY} \xrightarrow{\text{EOF / exit}} \text{TERMINATED}$$

1. **State `NEW`:**
   - Server hanya menerima request `initialize` atau `ping`.
   - Pemanggilan method lain (seperti `tools/list` atau `tools/call`) sebelum inisialisasi **MUST** menghasilkan error code `-32002` (*Server not initialized*).
2. **State `INITIALIZING`:**
   - Server telah merespons `initialize` dengan versi protokol `2026-07-28` dan kapabilitas server. Server menunggu notifikasi `notifications/initialized`.
3. **State `READY`:**
   - Server siap memproses `tools/list`, `tools/call`, dan notifikasi pembatalan `notifications/cancelled`.
   - Pemanggilan ulang `initialize` pada state `READY` **MUST** menghasilkan error code `-32600` (*Invalid Request: already initialized*).
   - Notifikasi yang tidak dikenal diabaikan secara senyap tanpa respon.
4. **State `TERMINATED`:**
   - Deteksi EOF pada `stdin` memicu shutdown anggun (*graceful shutdown*) dengan exit code 0.

### 3. Kontrak JSON-RPC 2.0 & Matriks Error
- Format request wajib menyertakan `"jsonrpc": "2.0"`.
- **Preservasi Request ID Presisi:** Field `id` dalam respon wajib mempertahankan format request secara verbatim (string atau integer). Dilarang melakukan mutasi tipe (misal: integer `1` berubah menjadi float `1.0`).
- **Matriks Kode Error Normatif:**
  | Kode Error | Kategori Standar | Kondisi Pemicu |
  | :---: | :--- | :--- |
  | **`-32700`** | Parse Error | JSON malformed atau frame melebihi 4 MB. |
  | **`-32600`** | Invalid Request | Versi jsonrpc bukan "2.0", missing method, atau pelanggaran state machine. |
  | **`-32601`** | Method Not Found | Method RPC yang dipanggil tidak terdaftar. |
  | **`-32602`** | Invalid Params | Parameter tool call tidak sesuai skema JSON Schema. |
  | **`-32000`** | Internal Tool Error | Eksekusi internal tool gagal atau crash. |
  | **`-32002`** | Server Not Initialized | Pemanggilan tool sebelum handshake `initialize` selesai. |

### 4. Pemisahan MCP Tool Registry vs Rule Registry
- **MCP Tool Registry:** Katalog tool level protokol MCP yang diekspos ke AI agent. Charites mengekspos tepat 3 tool:
  1. `charites_scan`
  2. `charites_explain_rule`
  3. `charites_list_rules`
- **Rule Registry:** Katalog rule analisis statis internal (`rules.Registry`). Tool MCP memanggil Rule Registry sebagai pustaka internal. Penambahan puluhan rule di Fase 8 **DILARANG** mengekspos rule individual sebagai tool MCP baru.

### 5. Batas Keamanan & Trust Boundary (`charites_scan`)
- **Workspace-Root Scoped:** Server MCP terikat pada direktori kerja resmi (`WorkspaceRoot`, default: direktori aktif saat perintah dijalankan, atau via flag `--workspace`).
- **Pencegahan Path Traversal:**
  - Parameter `path` yang bernilai relatif diresolusikan terhadap `WorkspaceRoot`.
  - Parameter `path` bernilai absolut **WAJIB** berada di dalam `WorkspaceRoot` (diverifikasi melalui `filepath.Clean`).
  - Path yang mencoba menembus keluar (misal: `../../etc/passwd` atau symlink ke luar workspace) **MUST** ditolak dengan error keamanan.
  - Path yang menunjuk langsung ke direktori builtin terlarang (`node_modules/`, `.git/`) ditolak sesuai kebijakan keamanan Fase 4.
- **Batas Waktu Eksekusi (Execution Timeout):**
  - Setiap pemanggilan `charites_scan` dibatasi batas waktu maksimum **30 detik** (`context.WithTimeout`). Jika melebihi 30 detik, eksekusi dibatalkan dan mengembalikan error timeout deterministik.
- **Pembatalan Interaktif (Cancellation Contract):**
  - Jika klien mengirimkan notifikasi `notifications/cancelled` dengan parameter `{"requestId": <id>}`, server segera membatalkan context eksekusi scan terkait. Hasil parsial yang belum lengkap dibuang.

### 6. Detail Skema Parameter Tool MCP
- **`charites_scan`**:
  - `path` (`string`, required): Path sasaran relatif di dalam workspace.
  - `category` (`string`, optional): Filter kategori rule.
  - `rule` (`string`, optional): Filter Charites Rule ID spesifik.
  - `ext` (`string`, optional): Filter ekstensi berkas target.
- **`charites_explain_rule`**:
  - `rule_id` (`string`, required): Charites Rule ID (contoh: `theme.hardcode-opacity-color`).
  - Output: Mengembalikan dokumentasi komprehensif, alasan larangan, contoh salah (*bad code*), contoh benar (*good code*), dan token pengganti semantik resmi dari SSOT `RuleMetadata`.
- **`charites_list_rules`**:
  - `category` (`string`, optional): Filter daftar berdasarkan kategori.

---

## BAGIAN B: Generator Ensiklopedia Dokumentasi (`charites wiki`)

### 1. Single Source of Truth (SSOT) Rule Metadata
Dokumentasi wiki dan respon tool `charites_explain_rule` bersumber dari metadata rule tunggal:
- Identitas: `ID`, `Category`, `DefaultSeverity`.
- Konten Edukasi: `Description`, `Explanation` (alasan arsitektural), `BadExample`, `GoodExample`, dan `Remediation` (token pengganti).

### 2. Penemuan Kategori Dinamis & Pengurutan Determinis
- Subcommand `charites wiki [output_dir]` (default: `./wiki/`) secara dinamis mengekstrak seluruh kategori unik dari `rules.Registry`.
- **Determinis Pengurutan (Total Ordering):**
  - Berkas kategori diurutkan alfabetis: `category ASC` $\rightarrow$ `wiki/<category>.md`.
  - Di dalam setiap berkas, entri rule diurutkan alfabetis berdasarkan `Rule.ID() ASC`.
  - Berkas indeks utama `wiki/Home.md` memuat tabel master seluruh rule terdaftar yang terurut rapi.
- **Invarian Nol Pergeseran Git (Zero Churn):** Generator dilarang menyertakan timestamp pembuatan yang berubah-ubah di dalam output markdown.

### 3. Generasi Berkas Bersifat Atomik (Atomic Directory Generation)
- Generator merender seluruh berkas ke dalam direktori staging sementara (`.wiki.tmp.<pid>`).
- Jika terjadi kegagalan rendering pada salah satu berkas, proses dibatalkan dan direktori sementara dihapus tanpa mengubah direktori `wiki/` yang sudah ada sebelumnya.
- Jika seluruh berkas berhasil dirender 100%, direktori staging dipindahkan secara atomik ke direktori sasaran.

---

## BAGIAN C: Skrip Pemasang Mandiri Aman (`scripts/install.sh`)

Skrip instalasi satu baris (*one-line installer*) wajib mematuhi standar keamanan rantai pasokan (*supply-chain security*):

```bash
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/scripts/install.sh | sh
```

### Kebutuhan Fungsional & Invarian Keamanan:
1. **HTTPS-Only Invariant:** Seluruh pengunduhan tarball dan manifest wajib menggunakan protokol HTTPS resmi ke GitHub Releases.
2. **Verifikasi Checksum SHA256:**
   - Skrip mengunduh tarball biner `charites-<version>-<os>-<arch>.tar.gz` dan berkas `checksums.txt`.
   - Skrip memverifikasi hash SHA256 biner terhadap manifest sebelum melakukan ekstraksi. Jika hash tidak cocok, instalasi langsung dibatalkan dengan exit code 1.
3. **Ekstraksi Terisolasi & Proteksi Path Traversal:**
   - Ekstraksi arsip tar dilakukan di direktori sementara terisolasi (`mktemp -d`).
   - Ekstraktor memvalidasi bahwa isi tarball tidak memuat relative path escape (`../`) atau link berbahaya.
4. **Pemasangan Atomik:**
   - Binary dipindahkan ke `/usr/local/bin/charites` (jika root/sudo) atau `$HOME/.local/bin/charites` dengan izin akses `0755`.
   - Trap `EXIT` dan `ERR` menjamin pembersihan file sementara secara otomatis saat skrip selesai atau gagal.
