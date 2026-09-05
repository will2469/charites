# 06-ROADMAP: 08 - Phase 8 Milestone, Rule Authoring Checklist & Expansion Phasing

> **Kode Dokumen:** `ROAD-08-EXPANSION`
> **Tahapan:** Fase 8 - Repetitive Pattern Flow Guide & Rule Authoring Template (Core Assessment)
> **Peran Pilar:** ROADMAP = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)
> **Status:** Ready for Review

Dokumen ini menetapkan *checklist* baku pengerjaan rule baru (*Standard Rule Authoring Checklist*), gerbang transisi kelulusan template otorisasi rule, serta peta jalan ekspansi bertahap (*batch expansion roadmap*). Sesuai prinsip tata kelola [docs/00-CONTRACT.md](https://github.com/will2469/charites/blob/main/docs/00-CONTRACT.md):
- **SPEC** = WHAT (Spesifikasi Kontrak Penulisan Rule & Template Studi Kasus)
- **ARCH** = HOW (Rancangan Arsitektur Ekstensibilitas, Rule SSOT & 3 Touchpoints)
- **TEST** = PROOF (Harness Pengujian, Matriks Semantik Ekspektasi & Tri-Corpus)
- **QUALITY** = QUALITY THRESHOLD (Ambang Batas Kualitas, Anti-Sycophancy & Anggaran Kinerja)
- **ROADMAP** = PHASE GATE (Otoritas Gerbang Evaluasi Kelulusan Transisi)

---

## 1. Rantai Prasyarat Gerbang Bertahap (Sequential Gate Prerequisite Chain)

Fase 8 adalah template ekspansi setelah fondasi compiler dan Rule #1 terbukti stabil. Pelaksanaan batch expansion rule baru **DILARANG KERAS** melompati tahapan sebelumnya:

$$\text{Fase 0--6 (Core Pipeline Frozen)} \longrightarrow \text{Rule \#1 ACTIVE \& Lolos Tri-Corpus} \longrightarrow \text{Fase 7 (MCP \& Wiki Stable)} \longrightarrow \text{Fase 8 (Batch Expansion Diizinkan)}$$

Fase 8 bukan tombol jalan pintas (*bypass button*) untuk langsung menulis kode, melainkan kerangka tata kelola penambahan aturan secara modular setelah seluruh gerbang prasyarat berstatus **PASS**.

---

## 2. Deliverables Berkas Fase 8

1. **`docs/01-spec/08-expansion.md`**: Dokumen spesifikasi kontrak penulisan rule baru dan matriks studi kasus `theme.hardcode-color`.
2. **`docs/02-architecture/08-expansion.md`**: Arsitektur 3 touchpoints dan interface `RuleMetadata` SSOT.
3. **`docs/03-testing/08-expansion.md`**: Protokol matriks ekspektasi kasus-per-kasus dan harness Tri-Corpus otomatis.
4. **`docs/04-quality/08-quality.md`**: Pemisahan invarian mutlak, target kualitas, dan anggaran performa rule baru.
5. **Peta Jalan Ekspansi Batch 1 s/d 5**: Rencana kerja terstruktur untuk porting 30+ rule warisan.

---

## 3. Gerbang Evaluasi Kelulusan (Phase Gate DoD)

Sebuah fase dinyatakan lulus (*graduated*) jika dan hanya jika seluruh evaluasi gerbang berikut berstatus **PASS**:

- [ ] **`ROAD-08-GATE-001` (SPEC-08 Compliance = PASS):**
  - Kontrak penulisan rule baru mendefinisikan 7 elemen wajib (ID kanonikal, batas cakupan, normalisasi varian, payload diagnostik dinamis, penekanan inline ignore, matriks semantik kasus, dan Tri-Corpus).
  - Studi kasus `theme.hardcode-color` diposisikan sebagai spesifikasi referensi percontohan.

- [ ] **`ROAD-08-GATE-002` (ARCH-08 Compliance = PASS):**
  - Seluruh komponen inti compiler (`internal/ir`, `parser`, `scanner`, `analyzer`, `reporter`, `mcp`, `wiki`) terkonfirmasi dibekukan (100% Shared Frozen Components).
  - Penambahan rule murni mengisolasi perubahan pada 3 Touchpoints (`<rule>.go`, `registry.go`, `tests/correctness/`).
  - Definisi `RuleMetadata` terintegrasi pada interface `Rule` sebagai SSOT untuk MCP dan Wiki.

- [ ] **`ROAD-08-GATE-003` (TEST-08 Compliance = PASS):**
  - Protokol pengujian mewajibkan matriks kasus-per-kasus (`matrix.json`) dengan kecocokan mutlak ($\text{Actual} \equiv \text{Expected}$).
  - Menghapus verifikasi longgar `PositiveViolations > 0`.

- [ ] **`ROAD-08-GATE-004` (QUAL-08 Compliance = PASS):**
  - Invarian anti-sycophancy (zero secret whitelists) dan fungsi evaluasi murni ditegakkan.
  - Alokasi memori dan performa terpisah antara Invarian, Target Kualitas ($\ge 90\%$ coverage, $\le 10$ cyclomatic), dan Anggaran Performa (`QUAL-08-PERF-001`).

---

## 4. Peta Jalan Ekspansi Bertahap (Batch Expansion Roadmap)

Begitu seluruh fondasi Fase 0 s/d 7 berstatus hijau dan diotorisasi:
1. **Batch 1 (Theme & Design Tokens):** `theme.hardcode-color`, `theme.apply-bloat`.
2. **Batch 2 (Accessibility):** `a11y.html-missing-lang`, `a11y.heading-skip`, `a11y.multiple-h1`, `a11y.img-missing-alt`, `a11y.button-accessible-name`.
3. **Batch 3 (Web Vitals & Performance):** `perf.lcp-priority`, `perf.cls-aspect-ratio`, `perf.inp-listeners`.
4. **Batch 4 (Responsive & Layout):** `layout.missing-breakpoint`, `layout.redundant-classes`.
5. **Batch 5 (SEO & Metadata):** `seo.meta-description`, `seo.canonical-url`.
