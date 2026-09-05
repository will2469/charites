---
name: charites-rule-scaffold
description: "MANDATORY SCAFFOLDING & EXTENSION HARNESS FOR CHARITES: Rapid, zero-defect scaffolding engine and exact code template generator for adding new Charites static analyzer rules (e.g. theme.*, a11y.*, responsive.*, perf.*, tailwind.*). Enforces Semgrep canonical rule identifiers (<category>.<slug>, e.g. theme.hardcode-opacity-color, strictly rejecting arbitrary txx/axx numbering) and standard charites:ignore directives. Generates internal/rules/<category>_<slug>.go, 1-SSOT Tri-Corpus golden fixtures under tests/correctness/<category>.<slug>/, rule registration, and authoritative 8-Pillars wiki documentation ('scaffold rule', 'new rule', 'add rule', 'create rule', 'bikin rule', 'scaffold-rule')."
compatibility: "Go 1.26+, Leaf IR Parser, bash"
metadata:
  version: "1.0.0"
  author: "Charites Team / Will"
  license: "MIT"
---

# Charites Rule Scaffolding Harness (`charites-rule-scaffold`)

Ekstensi dan pembuatan aturan linter baru untuk **Charites** (Go 1.26+, Astro, TSX, CSS) wajib mematuhi arsitektur **Zero-Divergence Rule Engine**, **Semgrep Canonical ID**, dan **1-SSOT Tri-Corpus Golden Test Standard**.

---

## 1. Aturan Penamaan Canonical Semgrep ID (Bukan `txx` atau `axx`)

> [!CRITICAL]
> **Charites TIDAK menggunakan kode penomoran arbitrer seperti `T01`/`txx` ataupun warisan Argus `A01`/`axx`!**
>
> Seluruh rule Charites wajib menggunakan **Semgrep Canonical Identifier** berformat **`<category>.<slug>`**:
> - `theme.hardcode-color`
> - `theme.hardcode-opacity-color`
> - `a11y.alt-text`
> - `responsive.touch-target`
> - `perf.inline-css`
> - `tailwind.arbitrary-color`
>
> Alasan: Charites adalah *domain-driven web static analyzer*. Nama rule harus *self-describing* dan mencerminkan kategori fungsionalnya agar mudah dibaca pada diagnostik CLI, IDE MCP, maupun konfigurasi `charites.yaml`.

---

## 2. Standar Direktif Supresi (`charites:ignore`)

Ketika pengguna atau fixture pengujian mengecualikan baris/blok dari pemeriksaan rule, sintaks supresi **WAJIB** merujuk pada Semgrep ID lengkap:

```astro
<!-- Astro / HTML: Single-Line Ignore -->
<!-- charites:ignore theme.hardcode-opacity-color intentional brand exception -->
<div class="bg-primary/42">Brand Opacity</div>
```

```tsx
// TSX / JSX / JS: Single-Line Ignore
// charites:ignore theme.hardcode-opacity-color intentional brand exception
<div className="bg-primary/42">Brand Opacity</div>
```

```css
/* CSS / PostCSS: Single-Line Ignore */
/* charites:ignore theme.hardcode-opacity-color intentional brand exception */
.badge {
  background-color: rgb(var(--primary) / 0.42);
}
```

```astro
<!-- Block Range Ignore -->
<!-- charites:ignore-start theme.hardcode-opacity-color -->
<div class="bg-primary/10">Legacy Widget 1</div>
<div class="bg-primary/20">Legacy Widget 2</div>
<!-- charites:ignore-end -->
```

> **DILARANG KERAS:**
> - `// charites:ignore T01`  *(Format txx ditolak)*
> - `// charites:ignore A01`  *(Format axx ditolak)*
> - `// charites:ignore` tanpa nama rule  *(Wajib menyebutkan Semgrep ID lengkap)*

---

## 3. Quick Start Generator

Untuk men-generate kerangka aturan baru secara otomatis:

```bash
# Men-generate rule baru (misal: theme.hardcode-opacity-color dengan severity HIGH)
./.agents/skills/charites-rule-scaffold/scripts/scaffold_rule.sh theme hardcode-opacity-color HIGH "Detects hardcoded opacity color slash modifiers"

# Men-generate rule accessibility (misal: a11y.alt-text dengan severity CRITICAL)
./.agents/skills/charites-rule-scaffold/scripts/scaffold_rule.sh a11y alt-text CRITICAL "Requires descriptive alt attributes on images"
```

Skrip ini akan secara otomatis membuat:
1. File implementasi rule di `internal/rules/<category>_<snake_slug>.go` (misal: `internal/rules/theme_hardcode_opacity_color.go`).
2. 1-SSOT Tri-Corpus di `tests/correctness/<category>.<slug>/` (`positive/`, `negative/`, `adversarial/`, dan `rule_test.go`).
3. Berkas dokumentasi 8-Pillars resmi di `wiki/<category>.<slug>.md`.

---

## 4. Directory & Asset Topology

```text
.agents/skills/charites-rule-scaffold/
├── SKILL.md                          # Panduan inti & checklist verifikasi (berkas ini)
├── assets/                           # Templat kode standar
│   ├── rule.go.tmpl                  # Kerangka rule analyzer Charites
│   ├── rule_test.go.tmpl             # Test runner suite 1-SSOT Tri-Corpus
│   └── wiki_rule.md.tmpl             # Templat dokumentasi 8-Pillars Matrix
├── references/                       # Dokumentasi pendukung
│   └── wiring_guide.md               # Snippet registrasi rule ke registry sentral
└── scripts/
    └── scaffold_rule.sh              # Generator bash otomatis (executable)
```

---

## 5. The 6-Step Atomic Rule Authoring Checklist

Setiap pembuatan rule baru di Charites wajib menuntaskan 6 langkah berurutan:

- [ ] **Langkah 1: Scaffolding Otomatis:** Jalankan `./.agents/skills/charites-rule-scaffold/scripts/scaffold_rule.sh <category> <slug> [severity] [description]`.
- [ ] **Langkah 2: Registrasi ke Registry Sentral:** Daftarkan rule di `internal/rules/registry.go` dalam fungsi `RegisterBuiltinRules()`.
- [ ] **Langkah 3: Implementasi Logika AST Visitor:** Lengkapi fungsi inspeksi AST pada `internal/rules/<category>_<snake_slug>.go` dengan traversal `ir.Walk(root)`.
- [ ] **Langkah 4: Kelengkapan Tri-Corpus:** Isi skenario P1-P5 (`positive/`), N1-N5 (`negative/` dengan `charites:ignore <category>.<slug>`), dan A1-A7 (`adversarial/`).
- [ ] **Langkah 5: Verifikasi Pengujian:** Jalankan `go test -v ./tests/correctness/<category>.<slug>/...` hingga seluruh tes PASS.
- [ ] **Langkah 6: Validasi Dokumentasi 8-Pillars:** Pastikan `wiki/<category>.<slug>.md` memuat penjelasan grounding, bad code, good code, dan panduan mitigasi.
