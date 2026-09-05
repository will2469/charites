# 06-ROADMAP: 03 - Phase 3 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-03-RULES`
> **Tahapan:** Fase 3 - Rule Contract & Proving Ground Rule (`theme.hardcode-opacity-color`)
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Ready for Execution

Dokumen ini menetapkan kriteria kelulusan (*exit criteria*) dan gerbang transisi (*phase gate*) untuk **Fase 3 (Rule Contract & Proving Ground Rule: `theme.hardcode-opacity-color`)** sebelum tim diizinkan melangkah ke **Fase 4 (Konfigurasi, Concurrency Scanner & Traversal Engine)**. Sesuai prinsip pemisahan otoritas arsitektur:
- **SPEC** = WHAT (Kontrak Antarmuka Rule, Detection Contract & Nomenklatur ID)
- **ARCH** = HOW (Enkapsulasi Registry Deterministik, Normalisasi Variant & Invarian Kemurnian)
- **TEST** = PROOF (Matriks Uji Batas Boundary, Tri-Corpus & Asersi Determinis Registry)
- **QUALITY** = QUALITY THRESHOLD (Invarian Pure Function, Allocation Gate & Performance Budget)
- **ROADMAP** = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)

---

## 1. Deliverables Berkas Fase 3

1. **`internal/rules/rule.go`**: Interface `Rule` (`ID`, `Description`, `Category`, `DefaultSeverity`, `Evaluate`).
2. **`internal/rules/registry.go`**: In-memory registry thread-safe dengan `sync.RWMutex`, pencarian $O(1)$, dan metode `All()` / `ByCategory()` dengan pengurutan deterministik berdasarkan `Rule.ID()`.
3. **`internal/rules/registry_test.go`**: Unit test pendaftaran rule, validasi penolakan ID duplikat, pengujian urutan deterministik, dan concurrent read/write test (`-race`).
4. **`internal/rules/theme/hardcode_opacity_color.go`**: Implementasi Rule #1 (`theme.hardcode-opacity-color`) dengan normalisasi Tailwind variant (`stripVariants`), lookup `OPACITY_TOKEN_MAP`, dan dynamic hint generator.
5. **`internal/rules/theme/hardcode_opacity_color_test.go`**: Table-driven unit test mencakup 20 skenario batas (in-scope, variants, clean negative, out-of-scope baits, arbitrary colors) serta benchmark alokasi memori.
6. **`tests/correctness/theme.hardcode-opacity-color/`**: Tiga sub-korpus uji nyata:
   - `positive/`: Berkas contoh dengan pelanggaran `bg-primary/10`, `border-destructive/20`, variant `hover:bg-primary/10`.
   - `negative/`: Berkas contoh bersih dengan token semantik resmi `bg-primary-light`, `text-muted`.
   - `adversarial/`: Berkas jebakan `w-1/2`, `aspect-16/9`, `text-sm/6`, `bg-primary/30`, `bg-[#123456]/10`.
7. **`tests/correctness_gate_test.go`**: Runner evaluasi otomatis yang menghitung dan memvalidasi `RuleCorrectnessMetric`.

---

## 2. Gerbang Evaluasi Kelulusan (Phase Gate DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [ ] **`ROAD-03-GATE-001` (SPEC-03 Compliance = PASS):**
  - Interface `Rule` terdefinisi bersih di `internal/rules/rule.go` tanpa dependensi ke scanner/engine.
  - Menggunakan canonical Charites Rule ID (`theme.hardcode-opacity-color`).
  - Detection Contract Rule #1 terkunci: hanya mendeteksi utility color dengan pemetaan token semantik resmi (`OPACITY_TOKEN_MAP`).
  - Out-of-Scope ditegakkan: layout fraction (`w-1/2`), aspect ratio (`aspect-16/9`), grid fraction (`grid-cols-2/3`), line-height (`text-sm/6`), unmapped opacity (`bg-primary/30`), dan arbitrary color (`bg-[#123456]/10`) diabaikan (0 diagnostic).
  - Diagnostic message dan hint di-generate secara dinamis sesuai mapping token pengganti.
  - Rule `Evaluate()` terpisah dari inline comment suppression (didelegasikan ke engine layer Fase 4).

- [ ] **`ROAD-03-GATE-002` (ARCH-03 Compliance = PASS):**
  - Rule bersifat murni (*pure function*) dan *stateless*.
  - Registri rule aman dari race condition menggunakan `sync.RWMutex`.
  - `Registry.All()` dan `Registry.ByCategory()` menjamin urutan deterministik (sorted by `Rule.ID()`).
  - Normalisasi varian (`stripVariants`) mengekstrak base utility untuk deteksi multi-prefix (`hover:`, `dark:`, `md:hover:`).

- [ ] **`ROAD-03-GATE-003` (TEST-03 Compliance = PASS):**
  - Matriks pengujian batas 20 skenario lolos 100% pada unit test.
  - Charites 1-SSOT Tri-Corpus Verification lolos:
    - `PositiveViolations > 0` (Terbukti mendeteksi pelanggaran).
    - `NegativeViolations == 0` (Zero Noise Invariant terpenuhi).
    - `AdversarialViolations == 0` (Bait Immunity Invariant terpenuhi).
  - `TestRegistry_DeterministicOrder` dan concurrent test lolos tanpa kegagalan atau race condition.

- [ ] **`ROAD-03-GATE-004` (QUAL-03 Compliance = PASS):**
  - Invarian fungsi evaluasi murni terpenuhi (0 disk/network I/O, AST read-only, idempotensi).
  - Allocation Invariant terpenuhi: `0 B/op` dan `0 allocs/op` saat evaluasi node bersih.
  - Benchmark dijalankan sesuai metodologi `QUAL-03-PERF-001` / `TEST-03-BENCH-001`.
  - Cakupan pengujian (*statement coverage*) paket `internal/rules/...` mencapai minimal $90\%$.
  - `golangci-lint run ./internal/rules/...` lulus 100% tanpa isu.

---

## 3. Gerbang Transisi ke Fase 4 (Engine & Scanner)

Begitu keempat gerbang di atas berstatus **PASS**:
1. Buat git commit: `feat(rules): implement rule interface, deterministic registry, and theme.hardcode-opacity-color rule with tri-corpus suite`.
2. Melangkah ke Fase 4: Buka dokumen `docs/01-spec/04-engine.md` untuk merancang parser konfigurasi `charites.yaml`, `.charitesignore` matcher, worker pool concurrency scanner, diagnostic suppression filter, dan traversal engine.
