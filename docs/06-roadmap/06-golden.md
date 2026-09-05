# 06-ROADMAP: 06 - Phase 6 Milestone & Pipeline Freeze Gate

> **Kode Dokumen:** `ROAD-06-GOLDEN`
> **Tahapan:** Fase 6 - Validasi Penuh & Golden Snapshots (Milestone Selesai Pipa)
> **Status:** Ready for Review

Dokumen ini menetapkan deliverable wajib, kriteria kelulusan (*exit criteria*), serta deklarasi pembekuan pipa compiler (*pipeline freeze gate*) untuk **Fase 6 (Validasi Penuh & Golden Snapshots)** sebelum melangkah ke **Fase 7 (Ekosistem MCP & Wiki Generator)** dan **Fase 8 (Ekspansi 30+ Rules)**.

---

## 1. Deliverables Wajib Fase 6

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

## 2. Checklist Definition of Done (DoD) Fase 6

- [ ] Seluruh skenario Golden Snapshot lulus 100% tanpa perbedaan satu byte pun (`go test -v ./tests/... -run TestPipeline_GoldenSnapshots`).
- [ ] Native Fuzzing berjalan minimal 60 detik per target tanpa memicu satupun unhandled panic atau memory leak.
- [ ] Pengujian proyek bersih (*clean_project*) terbukti menghasilkan 0 diagnostic dan exit code 0 (*Zero Noise Invariant*).
- [ ] Pengujian override konfigurasi (`charites.yaml`) terbukti berhasil mematikan rule pada snapshot.
- [ ] Pengujian benchmark membuktikan monorepo mock (1.000 file) selesai dipindai dalam waktu $< 100\text{ milidetik}$.
- [ ] Total test coverage keseluruhan repositori mencapai minimal $85\%$.
- [ ] `golangci-lint run ./...` bersih 100%.
- [ ] Arsitektur pipa inti (`internal/ir`, `internal/parser`, `internal/scanner`, `internal/analyzer`, `internal/reporter`) resmi dibekukan (*Frozen*).

---

## 3. Gerbang Transisi ke Fase 7 (Ekosistem MCP & Wiki) & Fase 8 (Ekspansi Rules)

Begitu seluruh kriteria DoD di atas terpenuhi:
1. Buat git commit: `chore(pipeline): lock core compiler pipeline with golden regression and native fuzzing suites`.
2. Binary Charites resmi dinyatakan **STABIL & SIAP EKSPANSI**.
3. Melangkah ke Fase 7: Buka dokumen `docs/01-spec/07-mcp.md` untuk merancang server MCP Stdio JSON-RPC 2.0 (`charites mcp`) dan generator ensiklopedia otomatis (`charites wiki`).
