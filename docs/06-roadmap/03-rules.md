# 06-ROADMAP: 03 - Phase 3 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-03-RULES`
> **Tahapan:** Fase 3 - Rule Contract & Proving Ground Rule (`theme.hardcode-opacity-color`)
> **Status:** Ready for Review

Dokumen ini menetapkan deliverable wajib, kriteria kelulusan (*exit criteria*), dan gerbang transisi (*transition gate*) untuk **Fase 3 (Rule Contract & Proving Ground Rule)** sebelum melanjutkan ke **Fase 4 (Konfigurasi, Concurrency Scanner & Traversal Engine)**.

---

## 1. Deliverables Wajib Fase 3

1. **`internal/rules/rule.go`**: Interface `Rule` (`ID`, `Description`, `Category`, `DefaultSeverity`, `Evaluate`).
2. **`internal/rules/registry.go`**: In-memory registry thread-safe dengan `sync.RWMutex`, lookup $O(1)$, filter kategori, dan defensive copy.
3. **`internal/rules/registry_test.go`**: Unit test pendaftaran rule, validasi error ID duplikat, dan concurrent read/write test (`-race`).
4. **`internal/rules/theme/hardcode_opacity_color.go`**: Implementasi penuh Rule #1 (`theme.hardcode-opacity-color`) berbasis `OPACITY_TOKEN_MAP`.
5. **`internal/rules/theme/hardcode_opacity_color_test.go`**: Unit test table-driven mencakup kasus sah, pelanggaran, jebakan slash, dan benchmark performa.
6. **`tests/correctness/theme.hardcode-opacity-color/`**: Tiga sub-korpus uji nyata:
   - `positive/`: Berkas contoh dengan pelanggaran `bg-primary/10`, `border-destructive/20`.
   - `negative/`: Berkas contoh bersih dengan token semantik `bg-primary-light`, `text-muted`.
   - `adversarial/`: Berkas jebakan `w-1/2`, `aspect-16/9`, `line-height`, dan inline ignore.
7. **`tests/correctness_gate_test.go`**: Runner evaluasi otomatis yang menghitung dan memvalidasi `RuleCorrectnessMetric`.

---

## 2. Checklist Definition of Done (DoD) Fase 3

- [ ] Interface `Rule` terdefinisi bersih di `internal/rules/rule.go` tanpa ketergantungan sirkular ke scanner/analyzer.
- [ ] In-Memory Registry mampu menangani pencarian $O(1)$ berdasarkan Semgrep ID dan aman dari race condition saat diakses oleh banyak goroutine (`go test -race ./internal/rules/...`).
- [ ] Rule `theme.hardcode-opacity-color` terimplementasi penuh dengan fast-path scan character `/` dan lookup map token semantik.
- [ ] Argus Tri-Corpus Verification lolos sempurna:
  - `PositiveViolations > 0` (Terbukti mendeteksi pelanggaran).
  - `NegativeViolations == 0` (Zero Noise Invariant terpenuhi).
  - `AdversarialViolations == 0` (Bait Immunity Invariant terpenuhi).
- [ ] Benchmark membuktikan `0 B/op` dan `0 allocs/op` saat mengevaluasi AST node yang tidak memiliki pelanggaran.
- [ ] Cakupan pengujian (*test coverage*) paket `internal/rules/...` mencapai minimal $95\%$.
- [ ] `golangci-lint run ./internal/rules/...` lulus 100% tanpa peringatan.

---

## 3. Gerbang Transisi ke Fase 4 (Engine & Scanner)

Begitu seluruh kriteria DoD di atas terpenuhi:
1. Buat git commit: `feat(rules): implement rule interface, registry, and theme.hardcode-opacity-color rule with tri-corpus suite`.
2. Melangkah ke Fase 4: Buka dokumen `docs/01-spec/04-engine.md` untuk merancang parser konfigurasi `charites.yaml` (prinsip default: YES), `.charitesignore` matcher, worker pool concurrency scanner, dan traversal analyzer engine.
