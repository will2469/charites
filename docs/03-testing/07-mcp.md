# 03-TESTING: 07 - MCP Server Protocol, Wiki Generator & Installer Verification Plan

> **Kode Dokumen:** `TEST-07-MCP`
> **Tahapan:** Fase 7 - Ekosistem Lanjutan (MCP Server & Wiki Generator)
> **Status:** Ready for Review
> **Standar Rujukan:** JSON-RPC 2.0 Conformance Testing & Shellcheck Automation

Dokumen ini mendefinisikan strategi pengujian untuk server **Model Context Protocol (MCP)** berbasis Stdio, pengujian integritas ekspor **Wiki Generator**, serta validasi skrip instalasi shell.

---

## 1. Skenario Pengujian Server MCP (`internal/mcp/`)

Pengujian MCP memverifikasi kepatuhan terhadap standar JSON-RPC 2.0 dengan menyuntikkan aliran byte ke `stdin` mock dan membaca `stdout` mock:

### 1.1. Pengujian Handshake & Discovery (`internal/mcp/stdio_test.go`)
- **Test Case 1 (Inisialisasi Protokol):**
  - Input: `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`.
  - Ekspektasi: Respon memuat `protocolVersion: "2026-07-28"` dan `serverInfo.name: "charites-mcp"`.
- **Test Case 2 (Tools Listing):**
  - Input: `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`.
  - Ekspektasi: Respon memuat 3 tool resmi: `charites_scan`, `charites_explain_rule`, dan `charites_list_rules`.

### 1.2. Pengujian Eksekusi Tool Call
- **Test Case 3 (`charites_scan` Execution):**
  - Input:
    ```json
    {
      "jsonrpc": "2.0",
      "id": 3,
      "method": "tools/call",
      "params": {
        "name": "charites_scan",
        "arguments": { "path": "tests/fixtures/astro_opacity" }
      }
    }
    ```
  - Ekspektasi: Respon mengembalikan payload JSON diagnostik yang mendeteksi pelanggaran opacity.
- **Test Case 4 (`charites_explain_rule` Execution):**
  - Input: Tool call dengan argumen `rule_id: "theme.hardcode-opacity-color"`.
  - Ekspektasi: Respon memuat string markdown deskripsi rule dan rekomendasi token semantik dari `global.css`.
- **Test Case 5 (Unknown Method / Bad Request):**
  - Input: `{"jsonrpc":"2.0","id":99,"method":"non_existent_method"}`.
  - Ekspektasi: Mengembalikan error code JSON-RPC `-32601` (*Method not found*).

### 1.3. Invarian Kemurnian Output (Zero Stream Contamination)
- Seluruh byte yang keluar dari `os.Stdout` wajib di-decode menggunakan `json.NewDecoder`. Jika terdapat baris yang gagal di-decode, pengujian dinyatakan **FAIL**.

---

## 2. Skenario Pengujian Wiki Generator (`internal/wiki/`)

Pengujian memastikan generator markdown menghasilkan struktur dokumentasi yang rapi dan dapat ditautkan:

```go
func TestWikiGenerator_Export(t *testing.T) {
    tmpDir := t.TempDir()
    reg := rules.NewRegistry()
    _ = reg.Register(theme.NewHardcodeOpacityColorRule())

    err := wiki.Generate(reg, tmpDir)
    if err != nil {
        t.Fatalf("failed to generate wiki: %v", err)
    }

    // 1. Verifikasi Home.md
    homeContent, err := os.ReadFile(filepath.Join(tmpDir, "Home.md"))
    if err != nil || !strings.Contains(string(homeContent), "theme.hardcode-opacity-color") {
        t.Errorf("Home.md missing rule entry")
    }

    // 2. Verifikasi theme.md
    themeContent, err := os.ReadFile(filepath.Join(tmpDir, "theme.md"))
    if err != nil || !strings.Contains(string(themeContent), "primary-light") {
        t.Errorf("theme.md missing remediation hint")
    }
}
```

---

## 3. Pengujian Kualitas Skrip Installer (`scripts/install.sh`)

1. **Static Analysis Shellcheck:**
   - Menjalankan linter shell: `shellcheck -s sh scripts/install.sh`.
   - Menjamin skrip bebas dari deklarasi bashism yang tidak kompatibel dengan `/bin/sh` murni.
2. **Uji Coba Deteksi Platform:**
   - Memvalidasi pemetaan `uname -s` dan `uname -m` ke target binary GitHub Releases (`charites-linux-amd64.tar.gz`, `charites-darwin-arm64.tar.gz`).
