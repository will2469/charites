# 06-ROADMAP: 07 - Phase 7 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-07-MCP`
> **Tahapan:** Fase 7 - Ekosistem Lanjutan (MCP Server, Wiki Generator & Secure Installer)
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Ready for Execution

Dokumen ini menetapkan deliverable wajib, kriteria kelulusan (*exit criteria*), dan gerbang transisi (*phase gate*) untuk **Fase 7 (Ekosistem Lanjutan: MCP Server, Wiki Generator & Secure Installer)** sebelum melangkah ke fase penutup **Fase 8 (Ekspansi 30+ Rules & Dokumentasi Ensiklopedia)**. Sesuai prinsip tata kelola [docs/00-CONTRACT.md](file:///home/will/Monorepo/charites/docs/00-CONTRACT.md):
- **SPEC** = WHAT (Spesifikasi Protokol MCP, Generator Wiki & Keamanan Pasokan)
- **ARCH** = HOW (Rancangan Protokol Stdio, Router MCP, Generator Wiki & Keamanan Pemasang)
- **TEST** = PROOF (Harness Pengujian, Matriks Protokol MCP & Validasi Pemasang)
- **QUALITY** = QUALITY THRESHOLD (Ambang Batas Kualitas, Integritas Protokol & Keamanan Pasokan)
- **ROADMAP** = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)

---

## 1. Deliverables Berkas Fase 7

1. **`internal/mcp/stdio.go`**: Implementasi loop komunikasi Stdio JSON-RPC 2.0 dengan pemisahan aliran mutlak (log ke `stderr`, JSON-RPC ke `stdout`) dan batas buffer 4 MB.
2. **`internal/mcp/dispatcher.go`**: Router JSON-RPC dengan state machine (`NEW` $\rightarrow$ `INITIALIZING` $\rightarrow$ `READY`) dan preservasi request ID presisi.
3. **`internal/mcp/tools.go`**: Katalog Tool Registry MCP yang mengekspos tepat 3 tool (`charites_scan`, `charites_explain_rule`, `charites_list_rules`) dengan enkapsulasi trust boundary dan batas waktu 30 detik.
4. **`internal/mcp/stdio_test.go`**: Suite pengujian 13 skenario protokol JSON-RPC 2.0 dan pembuktian kemurnian output stream.
5. **`internal/wiki/generator.go`**: Generator ensiklopedia markdown dinamis berbasis kategori rule dengan mekanisme pembuatan direktori atomik.
6. **`internal/wiki/generator_test.go`**: Unit test pembuatan berkas wiki markdown, penemuan kategori dinamis, dan verifikasi determinisme biner.
7. **`scripts/install.sh`**: Skrip instalasi aman satu baris berstandar POSIX dengan verifikasi hash SHA256 dan ekstraksi terisolasi.

---

## 2. Gerbang Evaluasi Kelulusan (Phase Gate DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [ ] **`ROAD-07-GATE-001` (SPEC-07 Compliance = PASS):**
  - Protokol MCP mematuhi JSON-RPC 2.0 dan siklus state machine (`NEW` $\rightarrow$ `INITIALIZING` $\rightarrow$ `READY`).
  - Request ID dipreservasi secara verbatim tanpa mutasi tipe data.
  - Batas keamanan trust boundary `charites_scan` menolak path traversal di luar workspace.
  - Batas waktu eksekusi 30 detik dan pembatalan interaktif via `notifications/cancelled` terbukti efektif.
  - Generator wiki mengekstrak kategori secara dinamis dan menghasilkan dokumen terurut deterministik.
  - Skrip installer memverifikasi checksum SHA256 sebelum ekstraksi tarball.

- [ ] **`ROAD-07-GATE-002` (ARCH-07 Compliance = PASS):**
  - Zero stdout contamination: `os.Stdout` murni memuat frame JSON-RPC valid.
  - Batas frame protokol 4 Megabytes ditegakkan.
  - Pemisahan arsitektur bersih: Tool Registry MCP terisolasi dari Rule Registry internal.
  - Generasi wiki bersifat atomik (staging direktori sementara sebelum pemindahan final).

- [ ] **`ROAD-07-GATE-003` (TEST-07 Compliance = PASS):**
  - Seluruh 13 skenario matriks protokol MCP (`MCP-TEST-001` s/d `MCP-TEST-013`) lulus 100%.
  - Pengujian determinisme wiki membuktikan output biner identik antar-run tanpa timestamp drift.
  - Pengujian keamanan pemasang membuktikan penolakan tarball korup dan manipulasi path traversal.

- [ ] **`ROAD-07-GATE-004` (QUAL-07 Compliance = PASS):**
  - Cakupan pengujian memenuhi ambang batas:
    - `internal/mcp`: $\ge 85\%$ line coverage.
    - `internal/wiki`: $\ge 85\%$ line coverage.
  - Skrip shell lulus uji linter `shellcheck -s sh scripts/install.sh` dengan 0 peringatan/error.
  - `golangci-lint run ./...` lulus 100% tanpa isu.

---

## 3. Gerbang Transisi ke Fase 8 (Ekspansi Rules Massal)

Begitu keempat gerbang di atas berstatus **PASS**:
1. Buat git commit: `feat(ecosystem): implement MCP stdio server, automated wiki generator, and secure installer`.
2. Melangkah ke fase penutup: Buka dokumen `docs/01-spec/08-expansion.md` untuk mengeksekusi ekspansi bertahap (Batch 1 s/d 5) porting 30+ rule warisan ke dalam direktori modular `internal/rules/<domain>/` tanpa menyentuh pipa inti engine.
