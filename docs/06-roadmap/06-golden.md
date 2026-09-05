# 06-ROADMAP: 06 - Phase 6 Milestone & Pipeline Freeze Gate

> **Kode Dokumen:** `ROAD-06-GOLDEN`
> **Tahapan:** Fase 6 - Validasi Penuh & Golden Snapshots (Milestone Selesai Pipa)
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Ready for Execution

Dokumen ini menetapkan kriteria kelulusan (*exit criteria*) serta deklarasi pembekuan pipa compiler (*pipeline freeze gate*) untuk **Fase 6 (Validasi Penuh & Golden Snapshots)** sebelum melangkah ke **Fase 7 (Ekosistem MCP & Wiki Generator)** dan **Fase 8 (Ekspansi 30+ Rules)**. Sesuai prinsip tata kelola [docs/00-CONTRACT.md](file:///home/will/Monorepo/charites/docs/00-CONTRACT.md):
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
   - `config/`: Berkas `global.css` (@theme Tailwind v4), `charites.yaml`, dan `.charitesignore`.
2. **`tests/golden/`**: Kumpulan berkas snapshot kebenaran mutlak:
   - `astro_opacity.golden.json` & `astro_opacity.golden.txt`
   - `tsx_opacity.golden.json` & `tsx_opacity.golden.txt`
   - `clean_project.golden.json` & `clean_project.golden.txt`
   - `config_override.golden.json` & `config_override.golden.txt`
3. **`tests/golden_test.go`**: Suite pengujian regresi snapshot otomatis dengan dukungan flag pembaruan `-update`.
4. **`tests/fuzz/`**: Native fuzzing suite Go 1.26 (`astro_fuzz_test.go`, `tsx_fuzz_test.go`).
5. **`tests/benchmark_test.go`**: Pengujian benchmark latensi pemindaian menyeluruh (*end-to-end latency benchmark*).

---

## 2. Gerbang Evaluasi Kelulusan (Phase Gate DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [ ] **`ROAD-06-GATE-001` (SPEC-06 Compliance = PASS):**
  - Seluruh skenario pengujian Golden Master (`.golden.json` dan `.golden.txt`) cocok 100% tanpa perbedaan byte.
  - Mekanisme pembaruan terkontrol `-update` terbukti bekerja dan menghasilkan diff bersih.
  - Korpus fixtures lengkap merepresentasikan berkas `.astro`, `.tsx`, `global.css`, dan konfigurasi.

- [ ] **`ROAD-06-GATE-002` (ARCH-06 Compliance = PASS):**
  - Arsitektur pipa inti (`internal/ir`, `internal/parser`, `internal/scanner`, `internal/analyzer`, `internal/reporter`) resmi dibekukan (*Frozen*).
  - Penambahan rule baru di Fase 8 tidak memerlukan modifikasi pada layer inti.

- [ ] **`ROAD-06-GATE-003` (TEST-06 Compliance = PASS):**
  - Seluruh suite regresi `TestPipeline_GoldenSnapshots` lulus 100%.
  - Native fuzzing berjalan minimal 60 detik per modul tanpa memicu crash atau unhandled panic.
  - Uji proyek bersih (*clean_project*) membuktikan invarian *Zero Noise* (0 diagnostic).

- [ ] **`ROAD-06-GATE-004` (QUAL-06 Compliance = PASS):**
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
