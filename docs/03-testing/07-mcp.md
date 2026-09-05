# 03-TESTING: 07 - MCP Server Protocol, Wiki Generator & Installer Verification Plan

> **Kode Dokumen:** `TEST-07-MCP`
> **Tahapan:** Fase 7 - Ekosistem Lanjutan (MCP Server, Wiki Generator & Secure Installer)
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Matriks Protokol MCP & Validasi Pemasang)
> **Status:** Ready for Review
> **Standar Rujukan:** JSON-RPC 2.0 Conformance Testing & Shellcheck Automation

Dokumen ini mendefinisikan strategi pengujian menyeluruh untuk server **Model Context Protocol (MCP)** berbasis Stdio, pengujian integritas ekspor **Wiki Generator**, serta validasi keamanan skrip instalasi shell.

---

## 1. Matriks Pengujian Protokol Server MCP (`internal/mcp/`)

Pengujian menyuntikkan aliran byte ke `stdin` mock dan memvalidasi `stdout` mock baris per baris:

| ID Kasus Uji | Kondisi Masukan JSON-RPC | Ekspektasi Respon / Kode Status |
| :--- | :--- | :--- |
| **`MCP-TEST-001`** | Handshake normal: `initialize` lalu `notifications/initialized` | Respon sukses dengan `protocolVersion: "2026-07-28"` |
| **`MCP-TEST-002`** | `tools/list` atau `tools/call` dipanggil sebelum `initialize` | Error `-32002` (*Server not initialized*) |
| **`MCP-TEST-003`** | `initialize` dipanggil kedua kali setelah server `READY` | Error `-32600` (*Invalid Request: already initialized*) |
| **`MCP-TEST-004`** | JSON sintaks cacat (malformed JSON string) | Error `-32700` (*Parse Error*) |
| **`MCP-TEST-005`** | Frame pesan melebihi batas 4 Megabytes | Error `-32700` (*Parse Error / Frame Exceeded*) |
| **`MCP-TEST-006`** | Request tanpa field `"jsonrpc": "2.0"` | Error `-32600` (*Invalid Request*) |
| **`MCP-TEST-007`** | Method RPC tidak dikenal (misal `foo/bar`) | Error `-32601` (*Method not found*) |
| **`MCP-TEST-008`** | Tool call dengan nama tool tidak terdaftar | Error `-32601` (*Tool not found*) |
| **`MCP-TEST-009`** | Tool call `charites_scan` tanpa parameter `path` | Error `-32602` (*Invalid Params: missing required path*) |
| **`MCP-TEST-010`** | `charites_scan` dengan path traversal `../../etc/passwd` | Error `-32602` (*Invalid Params: path traversal denied*) |
| **`MCP-TEST-011`** | Preservasi Request ID presisi (string `"req-abc"` dan integer `42`) | Respon memuat field `id` yang identik 100% tanpa konversi float |
| **`MCP-TEST-012`** | Klien mengirim notifikasi pembatalan `notifications/cancelled` | Proses pemindaian dihentikan dan context dibatalkan bersih |
| **`MCP-TEST-013`** | Invarian Kemurnian Output (Zero Stream Contamination) | Setiap baris di `stdout` lolos parsing JSON-RPC tunggal (0 polusi teks) |

---

## 2. Skenario Pengujian Wiki Generator (`internal/wiki/`)

```go
func TestWikiGenerator_DynamicCategoriesAndAtomic(t *testing.T) {
    tmpTarget := t.TempDir()
    reg := rules.NewRegistry()
    _ = reg.Register(theme.NewHardcodeOpacityColorRule())

    // 1. Uji Penemuan Kategori Dinamis & Generasi 3-Tier
    gen := wiki.NewGenerator(reg)
    err := gen.Generate(tmpTarget)
    if err != nil {
        t.Fatalf("failed to generate wiki: %v", err)
    }

    homeContent, _ := os.ReadFile(filepath.Join(tmpTarget, "Home.md"))
    if !strings.Contains(string(homeContent), "theme.hardcode-opacity-color") {
        t.Errorf("Home.md missing rule entry")
    }

    ruleDoc, err := os.ReadFile(filepath.Join(tmpTarget, "theme", "hardcode-opacity-color.md"))
    if err != nil || !strings.Contains(string(ruleDoc), "## 1. Overview & Core Invariant") {
        t.Errorf("rule 8-pillars document missing or corrupted: %v", err)
    }

    // 2. Uji Determinisme Biner (Dua kali generasi menghasilkan byte identik)
    secondTarget := t.TempDir()
    _ = gen.Generate(secondTarget)

    firstBytes, _ := os.ReadFile(filepath.Join(tmpTarget, "theme", "hardcode-opacity-color.md"))
    secondBytes, _ := os.ReadFile(filepath.Join(secondTarget, "theme", "hardcode-opacity-color.md"))
    if !bytes.Equal(firstBytes, secondBytes) {
        t.Errorf("Wiki output is not byte-for-byte identical across runs")
    }
}
```

---

## 3. Pengujian Kualitas & Keamanan Skrip Pemasang (`scripts/install.sh`)

1. **Analisis Statis Shellcheck:**
   - Menjalankan linter shell: `shellcheck -s sh scripts/install.sh`.
   - Wajib bebas 100% dari peringatan bashism (*pure POSIX sh compliance*).
2. **Uji Coba Penolakan Checksum Palsu:**
   - Menyediakan tarball dummy dengan hash yang tidak sesuai dengan `checksums.txt`.
   - Skrip pemasang **WAJIB** menolak ekstraksi dan keluar dengan exit code `1`.
3. **Uji Coba Pencegahan Path Traversal:**
   - Menyediakan tarball tiruan yang memuat entri path `../usr/bin`.
   - Skrip **WAJIB** menggagalkan instalasi dan menghapus folder sementara.
