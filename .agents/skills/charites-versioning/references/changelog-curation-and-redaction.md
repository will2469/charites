# Changelog Detection, Curation, and Redaction Guide
**Standard:** Keep a Changelog 1.1.0 & GitHub Releases Best Practices
**Target Project:** Charites Compiler & Static Linter (`github.com/will2469/charites`)

---

## 1. Taksonomi: Mana yang Harus Ditampilkan vs Dieliminasi

Release Notes ditulis untuk **manusia (frontend engineer, tech lead, CI/CD maintainer)** yang menggunakan Charites, **bukan dump mentah git log**.

```
                           [ SEMUA COMMIT GIT DIFF ]
                                      │
                 ┌────────────────────┴────────────────────┐
                 ▼                                         ▼
     [ USER-FACING VALUE ]                       [ INTERNAL NOISE / SLOP ]
    (WAJIB DITAMPILKAN)                         (REDAKSI / GABUNG / DROP)
                 │                                         │
  • Breaking Changes (`feat!:`, `BREAKING`)      • In-flight bug fixes (fix fitur unreleased)
  • Fitur Baru & Rules Baru (`feat:`)            • Micro-typo commits (`docs: fix typo`)
  • Fix bug nyata di produksi (`fix:`)           • Lint / Formatting / Whitespace (`gofmt`)
  • Optimasi performa signifikan (`perf:`)       • Dependabot bump (`build(deps): bump...`)
  • Deprekasi fungsi / rule (`Deprecated:`)      • Minor refactor tanpa efek samping (`clean code`)
```

### Matriks Klasifikasi Tampilan

| Kategori | Status | Kriteria & Tindakan | Contoh Commit Charites |
| :--- | :--- | :--- | :--- |
| ** Breaking Changes** | **WAJIB (Paling Atas)** | Perubahan API publik, flag CLI, format konfigurasi `charites.yaml` yang memutus workflow konsumen. Wajib sertakan panduan migrasi. | `feat!: change --format flag to accept json|markdown|inline only` |
| ** New Features** | **WAJIB** | Rule baru (`theme.*`, `a11y.*`), CLI subcommand (`wiki`, `mcp`), atau tool MCP baru. | `feat: implement native Model Context Protocol (MCP) server` |
| ** Bug Fixes** | **TERKURASI** | Memperbaiki false-positive (FP), false-negative (FN), panic traversal, atau frontmatter offset bug. | `fix: correct line span calculation for Astro frontmatter expressions` |
| ** Performance** | **TERKURASI** | Pengurangan latensi atau alokasi memori AST traversal yang terukur. | `perf: utilize sync.Pool for node allocator in Leaf IR builder` |
| ** Deprecations** | **WAJIB** | Rule atau flag yang ditandai deprecated untuk persiapan rilis Major berikutnya. | `feat: deprecate --ignore-dir in favor of .charitesignore` |
| ** Internal Hygiene** | **DIRINGKAS (1 Baris)** | Adopsi golden test, penambahan Makefile target, linting. | *Gabungkan 15+ commit test menjadi 1 bullet point.* |
| ** Micro-noise / Typo** | **DROP TOTAL** | Typos, whitespace, commit perbaikan coba-coba saat develop branch. | `fix: typo in comment`, `fix: test failing on CI` |

---

## 2. Aturan Redaksi Bug Kecil (*Minor Bugs Redaction*)

### Aturan 1: *The In-Flight Fix Elimination Rule*
> **Prinsip:** Jika sebuah bug muncul dan diperbaiki dalam siklus pengembangan fitur yang **belum pernah dirilis ke publik**, bug tersebut **DILARANG** dimasukkan ke dalam release notes sebagai "Bug Fix".

- *Kasus Nyata:*
  - Commit A: `feat: implement MCP server tool registry`
  - Commit B: `fix: fix nil pointer in MCP server tool registry` (dibuat 1 jam setelah commit A)
- *Tindakan:* Hapus Commit B dari release notes! Bagi publik, fitur MCP dirilis dalam keadaan langsung berfungsi, bukan "kami buat rusak lalu kami perbaiki sebelum rilis".

### Aturan 2: *Cohesive Synthesis (Penggabungan Commit Terdistribusi)*
> **Prinsip:** Jangan mengekspos 10 commit mikro berturut-turut. Sintesiskan menjadi satu kalimat teknis yang berbobot.

- *Mentah (AI Slop / Raw Git Log):*
  ```text
  - feat(theme): adopt Tri-Corpus golden tests
  - feat(a11y): adopt Tri-Corpus golden tests
  - ... (10 commit berulang)
  - feat(perf): adopt Tri-Corpus golden tests
  ```
- *Redaksi Bersih (Curated):*
  ```markdown
  * **1-SSOT Tri-Corpus Golden Test Adoption:** Migrated all Charites analyzer rules to the standardized Positive (P1-P5), Negative (N1-N5), and Adversarial (A1-A7) resilience validation harness.
  ```

---

## 3. Siklus Hidup Dokumen: `release_notes.md` vs `CHANGELOG.md`

### Arsitektur Dua Dokumen (Two-Document Architecture)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           ARSITEKTUR DOKUMEN RILIS                          │
│                                                                             │
│  [release_notes.md]  ──► HANYA BERISI DRAF RILIS SAAT INI (Active Target)   │
│                          Digunakan untuk GitHub Release body, copy-paste    │
│                          pengumuman rilis, atau release automation bot.     │
│                                      │                                      │
│                                      ▼ (Saat rilis disahkan / git tag)      │
│  [CHANGELOG.md]      ──► BUKU BESAR AKUMULATIF HISTORIS (Living Ledger)     │
│                          Rilis baru di-PREPEND di paling atas.              │
│                          Semua rilis terdahulu (v1.0.0, v0.9.0) TETAP ADA   │
│                          dan tersimpan utuh di bawahnya.                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Prosedur "Prepend" Rilis Baru:
1. Ketika rilis baru (misal `v1.1.0`) siap diterbitkan:
2. Baca draf dari `release_notes.md`.
3. Buka `CHANGELOG.md`. Sisipkan blok `## [v1.1.0] - YYYY-MM-DD` tepat di bawah header `# Changelog` dan di atas entri rilis lama `## [v1.0.0]`.
4. Rilis `v1.0.0` tidak pernah ditimpa atau dihapus!
5. Perbarui `release_notes.md` agar mencerminkan ringkasan rilis `v1.1.0`.

---

## 4. Format Standar Release Notes (`release_notes.md`)

```markdown
# Release Notes - Charites v1.1.0 (2026-09-05)

Charites v1.1.0 introduces native Model Context Protocol (MCP) server support, sub-millisecond AST caching, and complete 1-SSOT Tri-Corpus test adoption across all design token and layout rules.

---

###  New Features & Capabilities
* **Model Context Protocol (MCP 2026-07-28):** Native stateless JSON-RPC / stdio server exposing `charites_scan`, `charites_explain_rule`, and `charites_list_rules`.
* **Adoption of 1-SSOT Tri-Corpus:** Standardized Positive (P1-P5), Negative (N1-N5), and Adversarial (A1-A7) testing harness across all rules.

---

###  Bug Fixes & Diagnostic Precision
* **Astro Frontmatter Parser:** Fixed line number offset when expressions span multiline code blocks.
* **OKLCH Color Normalizer:** Corrected hue angle boundary wrapping for 360-degree color spaces.

---

###  Performance Improvements
* Memory pool optimization in Leaf IR builder reducing allocations by 40%.

---

###  Installation & Upgrade
```bash
# In-place self update
charites update # or: charites -u

# Via Go Toolchain
go install github.com/will2469/charites/cmd/charites@v1.1.0
charites --version
```

```
