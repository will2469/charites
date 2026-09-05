# 06-ROADMAP: 07 - Phase 7 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-07-MCP`
> **Tahapan:** Fase 7 - Ekosistem Lanjutan (MCP Server & Wiki Generator)
> **Status:** Ready for Review

Dokumen ini menetapkan deliverable wajib, kriteria kelulusan (*exit criteria*), dan gerbang transisi (*transition gate*) untuk **Fase 7 (Ekosistem Lanjutan: MCP Server & Wiki Generator)** sebelum melangkah ke fase penutup **Fase 8 (Ekspansi 30+ Rules & Dokumentasi Ensiklopedia)**.

---

## 1. Deliverables Wajib Fase 7

1. **`internal/mcp/stdio.go`**: Implementasi loop komunikasi Stdio JSON-RPC 2.0 dengan pemisahan aliran ketat (log ke `stderr`, JSON-RPC ke `stdout`).
2. **`internal/mcp/dispatcher.go`**: Router JSON-RPC yang menangani `initialize`, `notifications/initialized`, `tools/list`, dan `tools/call`.
3. **`internal/mcp/handlers.go`**: Handler untuk 3 tools resmi (`charites_scan`, `charites_explain_rule`, `charites_list_rules`).
4. **`internal/mcp/stdio_test.go`**: Unit test pengujian handshake protokol, eksekusi tool call, dan validasi framing JSON.
5. **`internal/wiki/generator.go`**: Modul pengekspor katalog rule ke direktori `wiki/` (`Home.md`, `theme.md`, `a11y.md`, `perf.md`, `layout.md`, `seo.md`).
6. **`internal/wiki/generator_test.go`**: Unit test pembuatan berkas wiki markdown dan verifikasi isi remedi.
7. **`scripts/install.sh`**: Skrip instalasi satu baris POSIX-compliant yang mendukung arsitektur Linux/macOS dengan verifikasi checksum SHA256.

---

## 2. Checklist Definition of Done (DoD) Fase 7

- [ ] Subcommand `charites mcp` berjalan stabil dan mematuhi spesifikasi JSON-RPC 2.0.
- [ ] Server MCP terbukti dapat dihubungkan ke klien AI Agent (Cursor / Claude / Antigravity).
- [ ] Tool `charites_scan` berhasil mengeksekusi pemindaian dan mengembalikan JSON diagnostik ke AI client.
- [ ] Tool `charites_explain_rule` berhasil mengembalikan penjelasan dan contoh remedi untuk rule yang diminta.
- [ ] Subcommand `charites wiki` berhasil mengekspor seluruh dokumentasi rule ke `wiki/*.md` secara atomik dan deterministik (tanpa timestamp chrun).
- [ ] Skrip `scripts/install.sh` lolos pengujian `shellcheck` tanpa peringatan dan terbukti mengunduh serta memvalidasi binary rilis.
- [ ] Cakupan pengujian (*test coverage*) paket `internal/mcp/...` dan `internal/wiki/...` mencapai minimal $85\%$.
- [ ] `golangci-lint run ./...` bersih 100%.

---

## 3. Gerbang Transisi ke Fase 8 (Ekspansi Rules Massal)

Begitu seluruh kriteria DoD di atas terpenuhi:
1. Buat git commit: `feat(ecosystem): implement MCP stdio server, automated wiki generator, and one-line installer`.
2. Melangkah ke fase terakhir: Buka dokumen `docs/01-spec/08-expansion.md` untuk merinci katalog lengkap 30+ rule warisan (*legacy rules*) yang akan di-porting ke `internal/rules/<domain>/` beserta ensiklopedia markdown resminya di `wiki/`.
