# 06-ROADMAP: 06 - Phase 6 Milestone & Pipeline Freeze Gate

> **Kode Dokumen:** `ROAD-06-GOLDEN`
> **Tahapan:** Fase 6 - Validasi Penuh & Golden Snapshots (Milestone Selesai Pipa)
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Graduated / Frozen (Pipeline Locked)

Dokumen ini menetapkan kriteria kelulusan (*exit criteria*) serta deklarasi pembekuan pipa compiler (*pipeline freeze gate*) untuk **Fase 6 (Validasi Penuh & Golden Snapshots)** sebelum melangkah ke **Fase 7 (Ekosistem MCP & Wiki Generator)** dan **Fase 8 (Ekspansi 30+ Rules)**. Sesuai prinsip tata kelola [docs/00-CONTRACT.md](https://github.com/will2469/charites/blob/main/docs/00-CONTRACT.md):
- **SPEC** = WHAT (Spesifikasi Validasi Pipa Lengkap, Korpus Fixture & Golden Snapshots)
- **ARCH** = HOW (Arsitektur Runner Snapshot, Harness Fuzzing & Pembekuan Pipa)
- **TEST** = PROOF (Suite Uji Regresi Golden Master & Native Fuzzing CI)
- **QUALITY** = QUALITY THRESHOLD (Invarian Nol-Regresi, Ketahanan Crash & Batas Kualitas)
- **ROADMAP** = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)

---

## 1. Deliverables Berkas Fase 6

1. **`tests/fixtures/`**: Direktori korpus berkas percontohan dunia nyata yang mencakup:
   - `astro/`: Kasus pelanggaran opacity, line offset frontmatter panjang, inline ignore, dan komponen bersih.
   - `tsx/`: Kasus pelanggaran atribut JSX, template literals, inline ignore, dan komponen bersih.
   - `projects/`: Sampel repositori mini (`clean/`, `opacity_violations/`, `config_override/`, `ignore_patterns/`).
2. **`tests/golden/projects/`**: Kumpulan berkas snapshot kebenaran mutlak:
   - `clean.golden.json` & `clean.golden.txt`
   - `opacity_violations.golden.json` & `opacity_violations.golden.txt`
   - `config_override.golden.json` & `config_override.golden.txt`
   - `ignore_patterns.golden.json` & `ignore_patterns.golden.txt`
3. **`tests/golden_test.go`**: Suite pengujian regresi snapshot otomatis dengan dukungan flag pembaruan lokal `-update`.
4. **`tests/fuzz/`**: Native fuzzing suite Go 1.26 bertingkat (`astro_fuzz_test.go`, `tsx_fuzz_test.go`, `pipeline_fuzz_test.go`).
5. **`tests/benchmark_test.go`**: Pengujian benchmark latensi pemindaian menyeluruh (`BENCH-06-E2E-001`).

---

## 2. Gerbang Evaluasi Kelulusan (Phase Gate DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [x] **`ROAD-06-GATE-001` (SPEC-06 Compliance = PASS):**
  - Seluruh skenario pengujian Golden Master (`.golden.json` dan `.golden.txt`) cocok 100% tanpa perbedaan byte.
  - Mekanisme pembaruan terkontrol `-update` terbukti bekerja dan menghasilkan diff bersih.
  - Korpus fixtures lengkap merepresentasikan berkas `.astro`, `.tsx`, `global.css`, dan konfigurasi.

- [x] **`ROAD-06-GATE-002` (ARCH-06 Compliance = PASS):**
  - Arsitektur pipa inti (`internal/ir`, `internal/parser`, `internal/scanner`, `internal/analyzer`, `internal/reporter`) resmi dibekukan (*Frozen*).
  - Penambahan rule baru di Fase 8 tidak memerlukan modifikasi pada layer inti.

- [x] **`ROAD-06-GATE-003` (TEST-06 Compliance = PASS):**
  - Seluruh suite regresi `TestPipeline_GoldenSnapshots` lulus 100%.
  - Native fuzzing berjalan minimal 60 detik per modul tanpa memicu crash atau unhandled panic.
  - Uji proyek bersih (*clean_project*) membuktikan invarian *Zero Noise* (0 diagnostic).

- [x] **`ROAD-06-GATE-004` (QUAL-06 Compliance = PASS):**
  - Memenuhi seluruh ambang batas `QUAL-06`:
    - Total line coverage repositori $\ge 85\%$.
    - Verifikasi data race: $0$ data race detected (`go test -race ./...`).
    - Linter compliance: `golangci-lint run ./...` lulus 100% tanpa isu.

---

## 3. Gerbang Transisi ke Fase 7 (Ekosistem MCP & Wiki) & Fase 8 (Ekspansi Rules)

Begitu keempat gerbang di atas berstatus **PASS**:
1. Buat git commit: `chore(pipeline): lock core compiler pipeline with golden regression and native fuzzing suites`.
2. Binary Charites resmi dinyatakan **STABIL & SIAP EKSPANSI**.
3. Melangkah ke Fase 7: Buka dokumen `docs/01-spec/07-mcp.md` untuk merancang server MCP Stdio JSON-RPC 2.0 (`charites mcp`) dan generator ensiklopedia otomatis (`charites wiki`).
