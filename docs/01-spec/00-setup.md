# 01-SPEC: 00 - Project Setup & Toolchain Specification

> **Kode Dokumen:** `SPEC-00-SETUP`
> **Tahapan:** Fase 0 - Inisialisasi Repositori & Toolchain
> **Peran Pilar:** SPEC = WHAT (Spesifikasi Fungsional & Kontrak Input/Output)
> **Status:** Graduated (All Phase Gates Passed)

Dokumen ini mendefinisikan spesifikasi formal untuk penyiapan awal (_setup_) repositori, modul Go, struktur direktori, dan aturan ignorasi berkas proyek **Charites**.

---

## 1. Spesifikasi Modul Go, Toolchain & Standar Kompilasi

- **Module Path:** `github.com/will2469/charites`
- **Go Version & Compatibility:**
  - **Language / Module Compatibility:** `Go 1.26` (dideklarasikan di `go.mod` sebagai `go 1.26`).
  - **Supported Local Toolchain:** `go1.26.x` (mendukung minor patch toolchain Go 1.26 lokal pengembang).
  - **CI / Reproducibility Baseline:** `go1.26.0` (dipin secara deterministik pada pipeline CI GitHub Actions dan release workflow).
- **Vendor Policy:** Dilarang mengimpor library pihak ketiga (_third-party dependencies_) pada Fase 0. Seluruh scaffolding entrypoint murni menggunakan **Go Standard Library**.
- **SPEC-00-BUILD-001 (Zero CGO Production & Verification Boundary Invariant):**
  - **Production & Release Artifacts Boundary:** Biner Charites (`bin/charites`) dan seluruh target kompilasi silang **MUST** dikompilasi secara penuh dengan CGO dinonaktifkan (`CGO_ENABLED=0`) menghasilkan biner statis murni tanpa ketergantungan pustaka C sistem (glibc/musl). Kode sumber repositori dilarang keras memuat berkas CGO (`.c`, `.h`) atau direktif `import "C"`.
  - **Verification & Dynamic Analysis Boundary:** Pengujian konkurensi dengan Go Race Detector (`-race`) dan cakupan kode berbasis atomic (`-covermode=atomic`) diizinkan dan beroperasi pada *verification environment* menggunakan `CGO_ENABLED=1` semata-mata untuk mengaktifkan instrumentasi Go ThreadSanitizer runtime, tanpa mengubah status zero-CGO kode sumber maupun artefak biner produksi.
- **SPEC-00-BUILD-002 (Cross-Platform Compilation Targets):**
  - Binary Charites **MUST** mendukung kompilasi silang (*cross-compilation*) native dengan `CGO_ENABLED=0` untuk 4 target platform rilis resmi:
    1. **Linux x86_64:** `GOOS=linux GOARCH=amd64`
    2. **Linux ARM64:** `GOOS=linux GOARCH=arm64`
    3. **macOS Apple Silicon:** `GOOS=darwin GOARCH=arm64`
    4. **Windows x86_64:** `GOOS=windows GOARCH=amd64`
- **Dependency File Policy:** Berkas `go.sum` **MUST NOT** be required pada Fase 0 saat dependensi eksternal bernilai nol (*zero external dependencies*).

---

## 2. Struktur Direktori Wajib & Reservasi Skeleton (Directory Skeleton)

Struktur repositori memisahkan secara tegas antara **berkas implementasi wajib Fase 0** dengan **reservasi direktori arsitektur** (*repository skeleton reservation*) untuk fase-fase berikutnya:

### A. Berkas Implementasi Wajib Fase 0 (Mandatory Phase 0 Files)
Fase 0 **MUST** menginisialisasi dan menyediakan berkas-berkas berikut:
1. `cmd/charites/main.go` - Entrypoint trampoline `os.Exit(cli.Execute(os.Args[1:]))`.
2. `internal/cli/root.go` - Root command dispatcher, routing flags `-v, --version`, `-h, --help`, subcommand `version`, dan unknown command handler.
3. `internal/cli/version.go` - Metadata rilis & runtime version info.
4. `internal/cli/root_test.go` - Unit test routing flag, subcommand, dan exit codes.
5. `tests/e2e/smoke_test.go` - Subprocess E2E smoke test binary terkompilasi.
6. Berkas konfigurasi root:
   - `go.mod` (Go 1.26, zero external dependencies)
   - `Makefile` (target `all`, `build`, `test`, `lint`, `clean`)
   - `.charitesignore` (pola ignore default Gitignore-compatible)
   - `.golangci.yml` (konfigurasi baseline linter)
   - `charites.example.yaml` (template konfigurasi rule & severity)

### B. Reservasi Direktori Arsitektur (Repository Skeleton Directory Reservation)
Direktori-direktori berikut dipersiapkan sebagai reservasi struktur arsitektur repositori (*directory reservation*). Direktori ini **TIDAK** diwajibkan telah memiliki implementasi logika pada Fase 0, melainkan diimplementasikan secara terisolasi per fase sesuai roadmap:
- `internal/ir/` - Diimplementasikan pada **Fase 1** (Leaf IR Data Contract: node, diagnostic).
- `internal/parser/` - Diimplementasikan pada **Fase 2** (Parser Tailwind CSS v4 `@theme`, Astro, TSX).
- `internal/rules/` - Diimplementasikan pada **Fase 3** (Rule interface, registry & Rule #1).
- `internal/analyzer/` - Diimplementasikan pada **Fase 3 & 4** (Context buffer & traversal engine).
- `internal/config/` - Diimplementasikan pada **Fase 4** (Parser charites.yaml & engine ignore).
- `internal/reporter/` - Diimplementasikan pada **Fase 4 & 5** (Terminal ANSI, JSON, markdown reporter).
- `internal/scanner/` - Diimplementasikan pada **Fase 5** (Dirwalker & worker pool).
- `internal/cli/scan.go` - Subcommand `scan` (aliases: `check`, `run`) diimplementasikan pada **Fase 5**.
- `internal/mcp/` & `internal/cli/mcp.go` - Subcommand `mcp` (JSON-RPC 2.0 daemon) diimplementasikan pada **Fase 7**.
- `internal/wiki/`, `internal/lifecycle/`, & `internal/cli/wiki.go` - Subcommand `wiki` & binary updater diimplementasikan pada **Fase 8**.
- `tests/fixtures/`, `tests/fuzz/` - Diimplementasikan pada **Fase 2**.
- `tests/correctness/` - Diimplementasikan pada **Fase 3** (Tri-Corpus test harness).
- `tests/golden/` - Diimplementasikan pada **Fase 4** (Regression golden snapshots).
- `tests/integration/` - Diimplementasikan pada **Fase 5**.

Pemisahan ini menjamin isolasi fase (*phase isolation*) yang bersih tanpa kebocoran kode masa depan (*future code leakage*) ke dalam build Fase 0.

---

## 3. Spesifikasi Default Ignore Pattern (`.charitesignore`)

Sistem pemindai **MUST** mendukung berkas `.charitesignore` dengan standar sintaks **Gitignore**:

1. **Aturan Sintaksis Gitignore-Compatible:**
   - Komentar diawali dengan tanda pagar `#`.
   - Pola direktori diakhiri dengan garis miring `/` (hanya mencocokkan folder).
   - Pola yang diawali garis miring `/` terikat ke akar repositori (_root-anchored_), sedangkan tanpa garis miring mencocokkan secara rekursif.
   - Karakter wildcard `*` mencocokkan karakter dalam satu tingkat direktori, dan globstar `**` mencocokkan direktori bersarang.
   - Tanda seru `!` mendukung _negation / un-ignore_ (contoh: `!important.min.css`).
2. **Kategori Default Bawaan (`.charitesignore`):**

```text
# ==============================================================================
# CHARITES DEFAULT IGNORE PATTERNS (.charitesignore)
# Gitignore-compatible ignore specification for frontend static analysis
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

## 4. Spesifikasi Kontrak Entrypoint Binary (`cmd/charites/main.go` & `internal/cli`)

- **Nama Binary Output:** `charites`
- **Kontrak CLI Fase 0 (Input, Output, Stream Routing, Exit Codes):**
  1. **Flag Versi (`--version` dan `-v`) serta Subcommand `version`:**
     - Mencetak informasi versi secara eksklusif ke `stdout`:
       ```text
       charites version 0.1.0-dev (go1.26.x)
       ```
     - Saluran `stderr` **MUST** tetap bersih (kosong).
     - Keluar dengan exit code `0`.
  2. **Flag Bantuan (`--help` dan `-h`):**
     - Mencetak panduan penggunaan dasar secara eksklusif ke `stdout`:
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
     - Saluran `stderr` **MUST** tetap bersih (kosong).
     - Keluar dengan exit code `0`.
  3. **Pemanggilan Tanpa Argumen (`[]string{}`):**
     - Mencetak panduan penggunaan dasar (Usage) ke `stdout`.
     - Saluran `stderr` **MUST** tetap bersih (kosong).
     - Keluar dengan exit code `0`.
  4. **Subcommand / Flag Tidak Dikenal (*Unknown Command / Invalid Flag*):**
     - Jika pengguna memberikan argumen atau subcommand yang tidak dikenali (contoh: `charites unknown-command` atau `charites --bogus`), CLI **MUST** mencetak pesan kesalahan deskriptif ke `stderr`.
     - Saluran `stdout` dilarang mencemari pesan kesalahan.
     - Keluar dengan exit code `2` (CLI argument syntax error).

---

## 5. Acceptance Criteria (Kriteria Fungsional Lolos Fase 0)

1. Perintah `go build -o bin/charites ./cmd/charites` berhasil tanpa error dan tanpa warning dengan `CGO_ENABLED=0` (`SPEC-00-BUILD-001`).
2. Kompilasi silang (*cross-compilation*) native berhasil untuk 4 target resmi: `linux/amd64`, `linux/arm64`, `darwin/arm64`, dan `windows/amd64` (`SPEC-00-BUILD-002`).
3. Binary `./bin/charites --version`, `./bin/charites -v`, serta `./bin/charites version` mencetak string versi ke `stdout` (bersih dari `stderr`) dan keluar dengan exit code `0`.
4. Binary `./bin/charites --help` dan `./bin/charites -h` mencetak panduan penggunaan ke `stdout` dan keluar dengan exit code `0`.
5. Binary `./bin/charites` (tanpa argumen) mencetak panduan penggunaan ke `stdout` dan keluar dengan exit code `0`.
6. Binary `./bin/charites unknown-command` dan `./bin/charites --bogus` mencetak pesan kesalahan secara presisi ke `stderr` dan keluar dengan exit code `2`.
7. Berkas `.charitesignore` tersedia di root dan memuat seluruh default pattern ignorasi Gitignore-compatible.
8. Berkas `go.sum` **MUST NOT** be required pada Fase 0 karena ketiadaan dependensi pihak ketiga (*zero external dependencies*).
9. Rantai dependensi Makefile `all: build test lint` terbukti berhasil dieksekusi secara berurutan pada checkout repositori baru (*fresh checkout*).
10. Seluruh skeleton folder di `internal/` dan `tests/` telah dibuat sebagai *directory reservations*.

