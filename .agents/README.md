# .agents/ - Charites AI Agent Governance & Skills Suite

> **Single Source of Truth (SSOT)** untuk seluruh *Skills*, *Rules*, dan *Guardrails* kecerdasan buatan dalam pengembangan **Charites** (`github.com/will2469/charites`).
>
> Seluruh aset dalam direktori ini berlisensi **MIT** ([LICENSE.md](LICENSE.md)).

---

## 1. Filosofi & Tujuan Arsitektur

Charites adalah *compiler-grade static analyzer & AST linter* berperforma tinggi yang dibangun dengan Go 1.26+ untuk ekosistem web modern (Astro, React/TSX, dan CSS / Tailwind CSS design tokens).

AI Agent yang berkontribusi pada repositori ini tunduk pada prinsip **Compiler-Grade Quality Governance**:

1. **Zero Hallucination & Evidence-Guided Reasoning:** Setiap klaim performa, keabsahan sintaks, atau optimasi alokasi memori wajib dibuktikan secara empiris melalui unit test, benchmark (`go test -bench`), atau bukti traversal Leaf IR nyata.
2. **Single Source of Truth (1-SSOT) Tri-Corpus:** Integritas rule static analyzer diukur dari ketahanannya terhadap matriks 17-pola: *Positive (P1-P5)*, *Negative (N1-N5)*, dan *Adversarial (A1-A7)*.
3. **Semgrep Canonical Identifiers:** Seluruh aturan menggunakan format `<category>.<slug>` (seperti `theme.hardcode-color`, `a11y.alt-text`). Dilarang keras menggunakan penomoran arbitrary seperti `txx` atau `axx`.
4. **Pure Stateless Architecture:** Protokol eksternal seperti MCP (Model Context Protocol revisi 2026-07-28) diimplementasikan secara **Pure Stateless** tanpa dual-track atau dependensi status sesi masa lalu.

---

## 2. Katalog Core Skills Charites (`.agents/skills/`)

```text
.agents/skills/
├── charites-anti-sycophancy/
│   └── SKILL.md                  # Mandatory: Evidence-guided engineering, RFC 2119, WCAG 2.2, Core Web Vitals, OKLCH
├── charites-golden-corpus/
│   └── SKILL.md                  # 1-SSOT Tri-Corpus standard (P1-P5, N1-N5, A1-A7) & 4-Layer Quality Pyramid
├── charites-rule-scaffold/
│   ├── SKILL.md                  # Scaffolding engine & 6-step atomic checklist untuk pembuatan rule baru
│   ├── assets/                   # Templat rule.go, rule_test.go, wiki_rule.md
│   ├── references/               # wiring_guide.md untuk registrasi ke internal/rules/registry.go
│   └── scripts/
│       └── scaffold_rule.sh      # Generator bash otomatis (executable)
├── charites-mcp/
│   ├── SKILL.md                  # Pure Stateless MCP 2026-07-28 (charites_scan, charites_explain_rule, charites_list_rules)
│   └── references/
│       └── spec_2026_07_28.md    # Spesifikasi teknis stateless JSON-RPC 2.0
├── charites-versioning/
│   ├── SKILL.md                  # Deterministic SemVer 2.0.0 evaluator & 5 public API boundaries
│   ├── assets/ & references/     # Criteria matrix, SemVer spec clauses, changelog curation
│   └── scripts/
│       ├── charites-diff-inspector.sh   # Analisis perbandingan git diff LATEST_TAG...HEAD
│       └── generate-release-notes.sh    # Generator rilis terkurasi & updater CHANGELOG.md
├── charites-wikis-builder/
│   └── SKILL.md                  # SOP pembuatan dokumentasi 8-Pillars (wiki/<category>.<slug>.md & wiki/Home.md)
└── charites-responsive-guidelines/
    └── SKILL.md                  # Panduan responsive UI/UX, Fitts's law (44x44px), container queries, WCAG 2.2
```

### Rincian Peran Masing-Masing Skill:

| Nama Skill | Sasaran & Fungsi Utama |
| :--- | :--- |
| **[`charites-anti-sycophancy`](skills/charites-anti-sycophancy/SKILL.md)** | Menolak *submissive alignment*, mencegah *hardcoded secret whitelists*, dan menegakkan validasi berbasis spesifikasi resmi W3C & Go 1.26. |
| **[`charites-golden-corpus`](skills/charites-golden-corpus/SKILL.md)** | Menegakkan struktur uji coba `tests/correctness/<category>.<slug>/` (P1-P5 untuk want triggers, N1-N5 untuk zero false-positives, A1-A7 untuk stress cases). |
| **[`charites-rule-scaffold`](skills/charites-rule-scaffold/SKILL.md)** | Mempercepat authoring rule baru dengan scaffolding otomatis (`./scaffold_rule.sh theme hardcode-opacity-color HIGH`), registrasi registry, dan pengujian. |
| **[`charites-mcp`](skills/charites-mcp/SKILL.md)** | Menyediakan panduan arsitektur server Model Context Protocol (MCP) revisi 2026-07-28 murni *stateless* tanpa sesi lama (*no dual-track*). |
| **[`charites-versioning`](skills/charites-versioning/SKILL.md)** | Menginspeksi git diff secara matematis untuk penentuan kenaikan versi `MAJOR.MINOR.PATCH` dan kurasi `release_notes.md`. |
| **[`charites-wikis-builder`](skills/charites-wikis-builder/SKILL.md)** | Memastikan dokumentasi terbuka bertaraf industri sesuai matriks 8 pilar (Overview, Grounding, Bad Code, Good Code, Remediation, dll.). |
| **[`charites-responsive-guidelines`](skills/charites-responsive-guidelines/SKILL.md)** | Menjaga kepatuhan ergonomi sentuh (44x44px), `@container` queries, keyboard handling, dan layout zero-JS Astro SSR. |

---

## 3. Invarian Penamaan & Direktif Supresi

1. **Semgrep Canonical ID (Kewajiban Mutlak):**
   - Format penamaan rule: `<category>.<slug>` (contoh: `theme.hardcode-color`, `theme.hardcode-opacity-color`, `a11y.alt-text`).
   - Dilarang menggunakan penomoran arbitrary seperti `T01`/`txx` ataupun `A01`/`axx`.
2. **Sintaks Direktif `charites:ignore`:**
   - **Astro / HTML:** `<!-- charites:ignore <category>.<slug> [alasan opsional] -->`
   - **TSX / JSX / JS:** `// charites:ignore <category>.<slug> [alasan opsional]`
   - **CSS / PostCSS:** `/* charites:ignore <category>.<slug> [alasan opsional] */`
   - **Blok:** `<!-- charites:ignore-start <category>.<slug> --> ... <!-- charites:ignore-end -->`

---

## 4. Alur Kerja Kontributor AI

Saat menerima tugas di repositori Charites:

1. **Konsultasikan Skill Terkait:** Baca `SKILL.md` yang relevan sebelum menulis atau memodifikasi kode produksi.
2. **Uji Mandiri:** Jalankan suite uji coba lokal via `make test` atau `make test-full`.
3. **Verifikasi Kepatuhan SemVer:** Sebelum rilis tag baru, jalankan `./.agents/skills/charites-versioning/scripts/charites-diff-inspector.sh`.
