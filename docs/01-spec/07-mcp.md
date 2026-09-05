# 01-SPEC: 07 - MCP Server, Wiki Generator & Installer Specification

> **Kode Dokumen:** `SPEC-07-MCP`
> **Tahapan:** Fase 7 - Ekosistem Lanjutan (MCP Server & Wiki Generator)
> **Status:** Ready for Review
> **Standar Rujukan:** Model Context Protocol (MCP) JSON-RPC 2.0 Specification / POSIX Shell Standards

Dokumen ini mendefinisikan spesifikasi kebutuhan fungsional untuk server **Model Context Protocol (MCP)** berbasis Stdio (`charites mcp`), generator dokumentasi ensiklopedia otomatis (**`charites wiki`**), serta skrip instalasi otomatis satu baris (*one-line installer*).

---

## 1. Spesifikasi Server Model Context Protocol (MCP) (`internal/mcp/`)

Charites menyediakan antarmuka terstandarisasi untuk berinteraksi langsung dengan AI Agent (seperti Claude Desktop, Cursor, Antigravity) melalui MCP:

### 1.1. Transport & Protokol
- **Subcommand:** `charites mcp`
- **Transport:** Standard Input/Output (`stdio`)
- **Protokol:** **JSON-RPC 2.0** (kompatibel dengan MCP Specification 2026-07-28).
- **Format Payload:** Pesan dibatasi dengan newline (`\n`).

### 1.2. Siklus Hidup Protokol (Lifecycle Handshake)
1. **`initialize` Request:**
   Mengembalikan info server:
   ```json
   {
     "protocolVersion": "2026-07-28",
     "serverInfo": {
       "name": "charites-mcp",
       "version": "1.0.0"
     },
     "capabilities": {
       "tools": {}
     }
   }
   ```
2. **`notifications/initialized`:** Konfirmasi inisialisasi selesai.
3. **`tools/list`:** Mengembalikan daftar tool yang tersedia beserta skema JSON Schema.
4. **`tools/call`:** Mengeksekusi tool dan mengembalikan konten hasil analisis.

### 1.3. Katalog Tool yang Diekspos

#### Tool 1: `charites_scan`
- **Deskripsi:** Memindai berkas frontend pada path tertentu dan mengembalikan daftar pelanggaran kode.
- **Parameter Skema:**
  ```json
  {
    "type": "object",
    "properties": {
      "path": { "type": "string", "description": "Path relatif atau absolut direktori/berkas yang dipindai" },
      "category": { "type": "string", "description": "Filter kategori rule (misal: theme, a11y)" },
      "rule": { "type": "string", "description": "Filter Charites Rule ID spesifik" },
      "ext": { "type": "string", "description": "Filter ekstensi berkas (misal: astro, tsx)" }
    },
    "required": ["path"]
  }
  ```
- **Output:** Konten teks berformat JSON terstruktur yang memuat ringkasan dan array `diagnostics`.

#### Tool 2: `charites_explain_rule`
- **Deskripsi:** Mengembalikan penjelasan mendalam, dampak buruk, dan rekomendasi perbaikan untuk suatu rule.
- **Parameter Skema:**
  ```json
  {
    "type": "object",
    "properties": {
      "rule_id": { "type": "string", "description": "Charites Rule ID (contoh: theme.hardcode-opacity-color)" }
    },
    "required": ["rule_id"]
  }
  ```
- **Output:** Dokumentasi markdown lengkap mengenai rule terkait beserta contoh *bad code* vs *good code*.

#### Tool 3: `charites_list_rules`
- **Deskripsi:** Mendaftar seluruh rule yang terkompilasi di dalam binary beserta kategori dan severity bawaan.
- **Parameter Skema:**
  ```json
  {
    "type": "object",
    "properties": {
      "category": { "type": "string", "description": "Opsional: filter berdasarkan bidang" }
    }
  }
  ```

---

## 2. Spesifikasi Wiki Generator (`internal/wiki/`)

Charites menyertakan subperintah otomatis untuk mengekspor seluruh dokumentasi rule ke direktori `wiki/`:

### 2.1. Sintaks Perintah
```bash
charites wiki [output_directory] # Default: ./wiki/
```

### 2.2. Struktur Berkas yang Dihasilkan
Dokumentasi disusun ringkas per-bidang (bukan 30+ berkas terpisah):
```text
wiki/
├── Home.md       # Indeks utama & tabel ringkasan seluruh rule
├── theme.md      # Rules bidang Theme & Token Semantik
├── a11y.md       # Rules bidang Aksesibilitas
├── perf.md       # Rules bidang Web Vitals & Performa
├── layout.md     # Rules bidang Layout & Responsive Design
└── seo.md        # Rules bidang SEO & Metadata
```

### 2.3. Format Konten Per-Rule di dalam Berkas Bidang
Setiap rule di dalam berkas bidang memuat format standar:
1. **Header:** `### <rule-slug> (Rule ID: <category>.<rule-slug>)`
2. **Badge Metadata:** Severity bawaan (`ERROR` / `WARN` / `INFO`).
3. **Deskripsi Teknis:** Mengapa pola ini dilarang dan dampaknya terhadap UX/performa.
4. **Contoh Pelanggaran (*Bad Practice*):** Snippet kode yang melanggar.
5. **Solusi Perbaikan (*Good Practice*):** Snippet kode rekomendasi yang benar.

---

## 3. Spesifikasi Skrip Instalasi Otomatis (`scripts/install.sh`)

Skrip bash instalasi satu baris (*one-line installer*):
```bash
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/scripts/install.sh | bash
```

### Kebutuhan Fungsional Skrip:
1. **Deteksi Platform Otomatis:**
   - OS: `Linux` (`linux`), `Darwin` (`darwin`).
   - Arsitektur: `x86_64` (`amd64`), `aarch64`/`arm64` (`arm64`).
2. **Pemberian Hak Akses & Lokasi Bin:**
   - Mengunduh tarball rilis dari GitHub Releases.
   - Memverifikasi checksum SHA256 berkas.
   - Meletakkan binary ke `/usr/local/bin/charites` (jika memiliki akses root) atau `$HOME/.local/bin/charites`.
   - Menyetel permission eksekusi `chmod +x`.
