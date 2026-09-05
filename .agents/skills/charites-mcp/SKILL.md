---
name: charites-mcp
description: "Authoritative architectural guidance and operational engineering harness for Model Context Protocol (MCP) version 2026-07-28 pure stateless implementation in Charites. Auto-triggers when authoring, reviewing, debugging, or upgrading Charites MCP server (internal/mcp), implementing stateless protocol handlers, server/discover RPC, _meta response caching, charites_scan, charites_explain_rule, or charites_list_rules ('mcp 2026-07-28', 'stateless mcp', 'charites mcp', 'charites_scan', 'charites_explain_rule', 'charites_list_rules')."
compatibility: "Requires modern agentic IDE environment, Go 1.26+, bash, and git"
metadata:
  version: "1.0.0"
  author: "Charites Team / Will"
  license: "MIT"
  citations:
    - "Model Context Protocol Specification (Revision 2026-07-28): https://modelcontextprotocol.io/specification/2026-07-28"
    - "JSON-RPC 2.0 Specification (2010): https://www.jsonrpc.org/specification"
    - "RFC 8259: The JavaScript Object Notation (JSON) Data Interchange Format (IETF)"
---

# Charites MCP Engineering Skill (`2026-07-28` Pure Stateless Standard)

> **Core Thesis:** Sebagai compiler-linter generasi terbaru (Go 1.26+), **Charites mengadopsi secara murni arsitektur Pure Stateless MCP (Spesifikasi 2026-07-28)** tanpa mewarisi beban historis (*no legacy baggage*) dan **tanpa dual-track**.
>
> Setiap request JSON-RPC 2.0 bersifat **self-describing** membawa blok metadata `_meta`, sepenuhnya menghilangkan *stateful session handshake* (`initialize`, `notifications/initialized`, `Mcp-Session-Id`). Server tidak mempertahankan sesi memori antar-request, memungkinkan eksekusi deterministik, bebas race-condition, dan ramah concurrency.

---

## 1. Arsitektur Pure Stateless MCP (`2026-07-28`)

Arsitektur Charites MCP memproses setiap pesan masuk secara independen:

```
                            [Incoming JSON-RPC Request]
                                         │
                                         ▼
                           ┌───────────────────────────┐
                           │   ValidateJSONRPC Gate    │
                           │  - JSON-RPC 2.0 schema    │
                           │  - Validates _meta block  │
                           └─────────────┬─────────────┘
                                         │
                                         ▼
                           ┌───────────────────────────┐
                           │   Pure Stateless Engine   │
                           │  - Zero session locks     │
                           │  - Idempotent execution   │
                           │  - Zero prior memory      │
                           └─────────────┬─────────────┘
                                         │
                                         ▼
                           ┌───────────────────────────┐
                           │   Tool Dispatcher Pool    │
                           │  - charites_scan          │
                           │  - charites_explain_rule  │
                           │  - charites_list_rules    │
                           └─────────────┬─────────────┘
                                         │
                                         ▼
                           [Outgoing JSON-RPC Response]
```

### Karakteristik Utama Pure Stateless:
1. **Zero Session Machine:** Tidak ada state machine `PreInit -> Init -> Ready`. Request valid langsung diproses.
2. **Self-Describing Metadata:** Protokol mengandalkan field `_meta.protocolVersion: "2026-07-28"` di setiap frame request.
3. **Idempotent RPC Handlers:** Subcommand `server/discover`, `tools/list`, dan `tools/call` dapat dipanggil secara bebas tanpa urutan saklek.
4. **Lockless Concurrency:** Karena tidak ada state sesi yang dibagikan antar-request, penanganan tool sepenuhnya lockless kecuali pembacaan registry aturan (`sync.RWMutex`).

---

## 2. Tools Resmi Charites MCP

Charites mengekspos 3 tool terdaftar:

### 1. `charites_scan`
- **Tujuan:** Menjalankan pemindaian AST statis pada path target (file atau direktori).
- **Input Schema:**
  ```json
  {
    "type": "object",
    "properties": {
      "path": { "type": "string", "description": "Absolute or relative path to file/directory to scan" },
      "category": { "type": "string", "description": "Optional category filter (theme, a11y, responsive, perf, tailwind)" },
      "rule": { "type": "string", "description": "Optional rule identifier (e.g. theme.hardcode-color)" }
    },
    "required": ["path"]
  }
  ```
- **Output:** Array terstruktur dari objek `ir.Diagnostic`.

### 2. `charites_explain_rule`
- **Tujuan:** Mengembalikan dokumentasi komprehensif 8-Pillars untuk rule yang diminta.
- **Input Schema:**
  ```json
  {
    "type": "object",
    "properties": {
      "rule_id": { "type": "string", "description": "Canonical Semgrep rule identifier (e.g. theme.hardcode-color)" }
    },
    "required": ["rule_id"]
  }
  ```
- **Output:** Markdown lengkap berformat 8-Pillars (Overview, Grounding, Bad Code, Good Code, Remediation).

### 3. `charites_list_rules`
- **Tujuan:** Menampilkan katalog seluruh rule yang tersedia beserta metadata keparahan dan kategorinya.
- **Input Schema:**
  ```json
  {
    "type": "object",
    "properties": {}
  }
  ```
- **Output:** JSON array memuat seluruh rule metadata (ID, Category, Severity, Description).

---

## 3. Penghapusan Pola Usang (*No Legacy Anti-Patterns*)

Karena Charites tidak memiliki klien legacy 2024, pola-pola berikut **DILARANG** diimplementasikan di `internal/mcp/`:

| Pola Usang (Deprecated) | Alasan Penolakan di Charites | Standar Pure Stateless Charites |
| :--- | :--- | :--- |
| **Session Tracking** (`Mcp-Session-Id`) | Menciptakan mutable shared state yang rawan memory leak. | **Stateless Per-Request** (`_meta`). |
| **Handshake State Machine** (`initialize`) | Memblokir pemanggilan tool sebelum handshake selesai. | **Immediate Readiness** pada saat server bootstrap. |
| **Sampling Callback** (`sampling/createMessage`) | Mengharuskan reverse connection yang kompleks. | **Standard Request/Response** deterministik. |
| **roots/list Protocol Gate** | Tidak aman mempercayai path boundaries dari klien. | **Server-side containment** via path traversal sanitizer. |

---

## 4. Guardrail Implementasi Charites MCP (`internal/mcp/`)

1. **Pure Stdio Stream:** Membaca baris tunggal JSON-RPC dari `stdin` dan menulis respons ke `stdout` dengan newline delimiter (`\n`). Semua log diagnostik internal diarahkan ke `stderr`.
2. **Buffer Flush Disiplin:** Setiap respons `stdout` wajib segera di-`Flush()` untuk mencegah deadlock buffering pada IDE client.
3. **Fail-Safe Isolation:** Kegagalan saat mem-parsing berkas web (misal template sintaks rusak) dilaporkan sebagai diagnostik error dalam payload response, tanpa mematikan proses MCP server.
