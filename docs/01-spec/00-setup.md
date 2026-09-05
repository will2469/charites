# 01-SPEC: 00 - Project Setup & Toolchain Specification

> **Kode Dokumen:** `SPEC-00-SETUP`
> **Tahapan:** Fase 0 - Inisialisasi Repositori & Toolchain
> **Status:** Ready for Execution

Dokumen ini mendefinisikan spesifikasi formal untuk penyiapan awal (_setup_) repositori, modul Go, struktur direktori, dan aturan ignorasi berkas proyek **Charites**.

---

## 1. Spesifikasi Modul Go

- **Module Path:** `github.com/will2469/charites`
- **Go Toolchain Version:** Go `1.26.0`.
- **Vendor Policy:** Dilarang mengimpor library pihak ketiga (_third-party dependencies_) pada Fase 0. Seluruh scaffolding entrypoint murni menggunakan **Go Standard Library**.
- **CGO Policy:** Wajib mendukung kompilasi `CGO_ENABLED=0` secara native tanpa memerlukan library C sistem.

---

## 2. Struktur Direktori Wajib (Directory Skeleton)

Eksekusi Fase 0 **MUST** menginisialisasi arsitektur repositori secara utuh sesuai **Single Source of Truth (SSOT)** pada [docs/stratch.md](file:///home/will/Monorepo/charites/docs/stratch.md):

```text
.
├── cmd/
│   └── charites/                  # Entrypoint binary utama (`charites`)
│       └── main.go                # [FASE 0] Entrypoint trampoline os.Exit(cli.Execute(...))
├── internal/                      # Private core packages (Enkapsulasi Go internal)
│   ├── cli/                       # Routing subcommands & flag parsing
│   │   ├── root.go                # [FASE 0] Root dispatcher, --version & --help flags
│   │   ├── scan.go                # [FASE 5] Subcommand: charites scan (aliases: check, run)
│   │   ├── mcp.go                 # [FASE 7] Subcommand daemon: charites mcp (JSON-RPC 2.0)
│   │   ├── wiki.go                # [FASE 8] Subcommand: charites wiki (Markdown catalog)
│   │   └── version.go             # [FASE 0] Metadata rilis & runtime version info
│   ├── config/                    # Konfigurasi & ignore engine
│   │   ├── config.go              # [FASE 4] Parser charites.yaml (severity, overrides)
│   │   └── ignore.go              # [FASE 4] Engine pattern matching (.charitesignore)
│   ├── ir/                        # Intermediate Representation (SSOT Data Contract)
│   │   ├── node.go                # [FASE 1] Unified node: Tag, Attribute, ClassList, Span
│   │   ├── diagnostic.go          # [FASE 1] Struct Diagnostic & Severity (SSOT Rule, CLI, & MCP)
│   │   └── builder.go             # [FASE 2] Normalisasi AST Astro & TSX ke IR tree
│   ├── parser/                    # Layer ekstraksi AST mentah
│   │   ├── tailwind/              # [FASE 2] Parser CSS @theme di global.css
│   │   ├── astro/                 # [FASE 2] Lexer frontmatter (---) & template Astro
│   │   └── tsx/                   # [FASE 2] TSX & JSX AST visitor
│   ├── rules/                     # Rule Registry & Interface
│   │   ├── registry.go            # [FASE 3] Rule lookup table & filter
│   │   ├── rule.go                # [FASE 3] Interface baku Rule (Evaluate pure function)
│   │   └── theme/
│   │       └── hardcode_opacity_color.go # [FASE 3] Rule #1: theme.hardcode-opacity-color
│   ├── analyzer/                  # Traversal engine & diagnostic accumulator
│   │   ├── context.go             # [FASE 3] State per-file, token scope, diagnostic buffer
│   │   └── engine.go              # [FASE 4] Traversal IR -> Rule Dispatcher (Go 1.26 iter)
│   ├── scanner/                   # Fast walker & konkurensi worker
│   │   ├── walker.go              # [FASE 5] Dirwalker patuh terhadap .charitesignore
│   │   └── pool.go                # [FASE 5] Worker pool distribusi file ke analyzer
│   ├── reporter/                  # Format output presentasi
│   │   ├── inline.go              # [FASE 4] Output terminal (ANSI color, file:line:col)
│   │   ├── json.go                # [FASE 4] Structured JSON datar untuk tooling/jq
│   │   └── markdown.go            # [FASE 5] Ringkasan tabel markdown
│   ├── mcp/                       # Server MCP Stdio (Offline mode AI Agent)
│   │   ├── server.go              # [FASE 7] Loop JSON-RPC 2.0 via os.Stdin / os.Stdout
│   │   ├── protocol.go            # [FASE 7] Protocol mapping spec 2026-07-28
│   │   └── tools.go               # [FASE 7] Tools: charites_scan, charites_explain_rule
│   ├── wiki/                      # Generator katalog rule otomatis
│   │   ├── generator.go           # [FASE 8] Ekspor metadata rule ke markdown di wiki/
│   │   └── templates/             # [FASE 8] Template panduan remediasi
│   └── lifecycle/                 # Distribusi & update binary
│       ├── installer.go           # [FASE 8] Manajemen path & symlink
│       ├── updater.go             # [FASE 8] Self-binary replace via GitHub Release
│       └── uninstaller.go         # [FASE 8] Pembersihan binary dari PATH
├── wiki/                          # Ensiklopedia detail rules ringkas per bidang (theme, a11y, perf, dll.)
│   └── .gitkeep                   # Diisi pada Fase 8
├── scripts/                       # Shell automation helper
│   ├── install.sh                 # [FASE 8] One-liner curl installer
│   └── uninstall.sh               # [FASE 8] Uninstaller script
├── tests/                         # Test fixtures, correctness corpus, & suites
│   ├── correctness/               # [FASE 3] Model Evaluasi Semantik Argus (Tri-Corpus)
│   │   └── theme.hardcode-opacity-color/
│   │       ├── positive/          # True violations (Wajib terdeteksi > 0)
│   │       ├── negative/          # Clean valid code (Zero-Noise Invariant == 0)
│   │       └── adversarial/       # False positive bait & syntax stress tests
│   ├── golden/                    # [FASE 4] Snapshot regression (JSON & ANSI output)
│   ├── fixtures/                  # [FASE 2] Berkas mentah (.astro, .tsx, global.css)
│   ├── integration/               # [FASE 5] Test pipeline end-to-end
│   ├── fuzz/                      # [FASE 2] Go 1.26 native fuzzing untuk Parser & IR
│   └── e2e/                       # [FASE 0] Subprocess CLI runner smoke tests
│       └── smoke_test.go          # [FASE 0] Verifikasi binary terkompilasi
├── .charitesignore                # [FASE 0] Default ignore patterns (Semgrep-compatible)
├── .golangci.yml                  # [FASE 0] Konfigurasi 9 linter wajib
├── charites.example.yaml          # [FASE 0] Template konfigurasi rule & severity
├── Makefile                       # [FASE 0] Perintah build, test, lint, format
├── go.mod                         # [FASE 0] Modul github.com/will2469/charites (Go 1.26)
└── go.sum
```

---

## 3. Spesifikasi Default Ignore Pattern (`.charitesignore`)

Sistem pemindai **MUST** mendukung berkas `.charitesignore` dengan standar sintaks **Semgrep / Gitignore**:

1. **Aturan Sintaksis Semgrep-Compatible:**
   - Komentar diawali dengan tanda pagar `#`.
   - Pola direktori diakhiri dengan garis miring `/` (hanya mencocokkan folder).
   - Pola yang diawali garis miring `/` terikat ke akar repositori (_root-anchored_), sedangkan tanpa garis miring mencocokkan secara rekursif.
   - Karakter wildcard `*` mencocokkan karakter dalam satu tingkat direktori, dan globstar `**` mencocokkan direktori bersarang.
   - Tanda seru `!` mendukung _negation / un-ignore_ (contoh: `!important.min.css`).
2. **Kategori Default Bawaan (`.charitesignore`):**

```text
# ==============================================================================
# CHARITES DEFAULT IGNORE PATTERNS (.charitesignore)
# Semgrep-compatible ignore specification for frontend static analysis
# ==============================================================================

# 1. Package Manager & Third-Party Dependencies
node_modules/
vendor/

# 2. Build Outputs, Transpiled Bundles & Framework Caches
dist/
build/
out/
.astro/
.next/
.nuxt/
.output/
.cache/
coverage/

# 3. Minified & Generated Production Bundles
*.min.js
*.min.css
*.bundle.js
*.map

# 4. Dependency Lockfiles
package-lock.json
pnpm-lock.yaml
yarn.lock
bun.lockb

# 5. Version Control Systems & OS Metadata
.git/
.svn/
.hg/
.DS_Store
Thumbs.db

# 6. Temporary Files & Crash Logs
*.tmp
*.temp
*.log
*.swp
*.swo
```

---

## 4. Spesifikasi Entrypoint Binary (`cmd/charites/main.go`)

- **Nama Binary Output:** `charites`
- **Behavior Fase 0:**
  - Jika dipanggil tanpa argumen atau dengan flag `-v, --version`: Mencetak informasi versi:
    ```text
    charites version 0.1.0-dev (go1.26)
    ```
    dan keluar dengan exit code `0`.
  - Jika dipanggil dengan flag `-h, --help`: Mencetak panduan penggunaan dasar:

    ```text
    Usage: charites <command> [options] [path]

    Commands:
      scan       Scan frontend files for design system, a11y, and performance issues
                 Aliases: check, run
      version    Print binary version

    Options for 'scan':
      -f, --format string      Output format: inline (default ANSI) or json
      --ext string             Filter by extension: astro, tsx, jsx
      --category string        Filter by category: theme, a11y, perf, layout, seo
      --rule string            Filter by single rule ID: theme.hardcode-opacity-color
      --ignore string          Additional custom ignore pattern
    ```

    dan keluar dengan exit code `0`.

---

## 5. Acceptance Criteria (Kriteria Lolos Fase 0)

1. Perintah `go build -o bin/charites ./cmd/charites` berhasil tanpa error dan tanpa warning.
2. Binary `bin/charites --version` mencetak string versi dan keluar dengan exit code `0`.
3. Berkas `.charitesignore` ada dan memuat seluruh direktori build/dependency umum.
4. Seluruh skeleton folder di `internal/` dan `tests/` telah dibuat dan terdaftar.
