---
name: charites-wikis-builder
description: "Standard operational procedure (SOP) and automated transformation engine to generate clean, authoritative, public-facing open-source Rule Wikis (wiki/Home.md, wiki/<category>.md, wiki/<category>/<slug>.md) conforming to Charites 8-Pillars documentation standards and 100% automated SSOT code generation."
compatibility: "Markdown, Go 1.26+, text/template, embed.FS, W3C/WCAG documentation standards"
metadata:
  version: "2.1.0"
  author: "Charites Team / Will"
  license: "MIT"
  citations:
    - "https://github.com/will2469/charites/wiki"
    - "W3C CSS Color Module Level 4"
    - "W3C Web Content Accessibility Guidelines (WCAG 2.2)"
    - "Core Web Vitals Performance Standards"
---

# Charites Automated Wiki Builder & SSOT Engine (`charites-wikis-builder`)

> **Mandat Utama:** Menyediakan dokumentasi static analysis bertaraf industri yang **100% otomatis dihasilkan dari kode Go (Single Source of Truth)**. Developer **DILARANG KERAS** membuat atau mengedit berkas `.md` wiki secara manual baris-demi-baris; seluruh dokumentasi di [`wiki/`](wiki/) di-*compile* secara dinamis oleh package `internal/wiki` menggunakan embedded Go templates (`//go:embed templates/*.tmpl`).

---

## 1. Topologi Direktori Wiki Resmi

Sesuai arsitektur GitHub Wiki Gollum (Flat Root Namespace & Collapsible Sidebar Navigation):

```text
wiki/
├── Home.md                                 # Master Catalog: Tabel kategori, rule counts & indeks seluruh rule
├── _Sidebar.md                             # Navigasi Collapsible: <details open><summary><b>Category</b> (Count)</summary>
├── theme.md                                # Domain Overview: Ringkasan kategori Theme & Design Tokens
├── theme.hardcode-opacity-color.md         # Spesifikasi Lengkap 8-Pillars Rule (Flat Namespace)
├── theme.hardcode-color.md                 # Spesifikasi Lengkap 8-Pillars Rule (Flat Namespace)
├── a11y.md                                 # Domain Overview: Ringkasan kategori Accessibility
├── a11y.alt-text.md                        # Spesifikasi Lengkap 8-Pillars Rule (Flat Namespace)
├── responsive.md                           # Domain Overview: Ringkasan kategori Responsive Design
├── responsive.touch-target.md              # Spesifikasi Lengkap 8-Pillars Rule (Flat Namespace)
├── perf.md                                 # Domain Overview: Ringkasan kategori Performance & Web Vitals
└── perf.inline-css.md                      # Spesifikasi Lengkap 8-Pillars Rule (Flat Namespace)
```

> [!IMPORTANT]
> **Invarian Format Tautan GitHub Wiki:**
> 1. **Flat Namespace (`<category>.<slug>.md`):** Gollum (engine GitHub Wiki) memberlakukan flat root namespace. Subdirektori (`theme/foo.md`) tidak didukung sebagai halaman wiki dan akan di-redirect ke `raw.githubusercontent.com`. Seluruh rule wajib berada langsung di root `wiki/`.
> 2. **Tanpa Ekstensi `.md` pada Tautan Internal:** Seluruh tautan internal wajib menghilangkan ekstensi `.md` (misal: `[theme.hardcode-color](theme.hardcode-color)` dan `[Home](Home)`). Menyertakan `.md` pada URL markdown membuat GitHub menganggapnya sebagai raw blob download.
> 3. **Collapsible Sidebar Accordion:** Di `_Sidebar.md`, setiap kategori dibungkus dengan `<details open><summary><b>{{.Title}}</b> ({{.Count}} rules)</summary>...</details>` untuk mencegah bloat navigasi ketika ada ratusan rules.

---

## 2. Prinsip Single Source of Truth (Zero Manual Editing)

Untuk mencegah desinkronisasi dokumentasi (*documentation drift*):

1. **Definisi di Kode Go:** Setiap rule mendefinisikan seluruh metadata dan penjelasannya langsung pada method `Doc() ir.RuleDocumentation` di struct rule-nya (`internal/rules/<category>/<snake_slug>.go`).
2. **Generasi Otomatis via `make wiki`:**
   ```bash
   make wiki
   ```
   Command ini mengeksekusi `internal/wiki/generator.go` yang me-render:
   - `wiki/Home.md` (menghitung jumlah rule per kategori secara otomatis dan mendaftar seluruh rule).
   - `wiki/_Sidebar.md` (navigasi collapsible accordion dengan link extension-less).
   - `wiki/<category>.md` (membuat tabel index rule per kategori dengan tautan relatif ke spesifikasi rule).
   - `wiki/<category>.<slug>.md` (me-render spesifikasi lengkap 8-Pillars langsung dari method `Doc()`).
3. **Penyelarasan CI/CD & Testing:** Generator wiki divalidasi pada pengujian integrasi (`internal/wiki/generator_test.go`). Pengujian `TestGenerator_RegenerateWiki` menjamin output deterministik biner dan zero diff terhadap rule yang terdaftar.

---

## 3. Struktur Kontrak Dokumentasi Rule di Kode Go (`ir.RuleDocumentation`)

Setiap struct rule mengimplementasikan interface `rules.Rule` dan `rules.DocumentedRule` (`internal/rules/doc.go` & `internal/ir/doc.go`):

```go
package theme

import (
    "github.com/will2469/charites/internal/ir"
)

type HardcodeOpacityColorRule struct{}

func (r *HardcodeOpacityColorRule) ID() string {
    return "theme.hardcode-opacity-color"
}

func (r *HardcodeOpacityColorRule) Category() string {
    return "theme"
}

func (r *HardcodeOpacityColorRule) DefaultSeverity() ir.Severity {
    return ir.SeverityError
}

func (r *HardcodeOpacityColorRule) Description() string {
    return "Detects utility classes with hardcoded slash opacity modifiers that have official semantic token replacements"
}

func (r *HardcodeOpacityColorRule) Doc() ir.RuleDocumentation {
    return ir.RuleDocumentation{
        TargetStandards: []string{
            "W3C Design Tokens Community Group (DTCG)",
            "Tailwind CSS Design Token Architecture",
            "WCAG 2.2 Relative Contrast",
        },
        CoreInvariant: "Every color opacity variation that represents a semantic state or visual elevation must use a centralized semantic design token rather than an arbitrary slash modifier.",
        Grounding: "In modern design token architecture, semantic colors like primary and destructive are calibrated for foreground/background contrast against explicit color stops...",
        BadExamples: []ir.CodeExample{
            {
                Language: "astro",
                Comment:  "Direct slash opacity modifiers on semantic colors",
                Code:     `<div class="bg-primary/10">Bad</div>`,
            },
            {
                Language: "tsx",
                Comment:  "Chained and single variants with hardcoded opacity",
                Code:     `<div className="hover:bg-primary/10">Bad</div>`,
            },
        },
        GoodExamples: []ir.CodeExample{
            {
                Language: "astro",
                Comment:  "Using official semantic tokens from global.css",
                Code:     `<div class="bg-primary-light">Good</div>`,
            },
            {
                Language: "tsx",
                Comment:  "Using semantic tokens with variants",
                Code:     `<div className="hover:bg-primary-light">Good</div>`,
            },
        },
        Risks: []ir.RiskItem{
            {
                Vector:   "Accessibility Degradation",
                Severity: "HIGH",
                Impact:   "Contrast ratio drops below 4.5:1 under dark mode themes due to uncalibrated alpha blending.",
            },
        },
    }
}
```

---

## 4. Taksonomi 8-Pillars Dokumentasi Charites

Template [`internal/wiki/templates/rule.md.tmpl`](internal/wiki/templates/rule.md.tmpl) secara deterministik memproduksi 8 pilar dokumentasi lengkap untuk setiap rule:

| Pilar | Bagian Dokumen | Sumber Data (`ir.RuleDocumentation` / `rules.Rule`) |
| :---: | :--- | :--- |
| **Pilar 1** | **Metadata Header & Target Standards** | `r.ID()`, `r.Category()`, `r.DefaultSeverity()`, `doc.TargetStandards` |
| **Pilar 2** | **1. Overview & Core Invariant** | `r.Description()`, `doc.CoreInvariant` |
| **Pilar 3** | **2. Technical Grounding & Engine Realities** | `doc.Grounding` |
| **Pilar 4** | **3. Vulnerability & Risk Taxonomy** | `doc.Risks` (`Vector`, `Severity`, `Impact`) |
| **Pilar 5** | **4. Non-Compliant Code Patterns (Bad Examples)** | `doc.BadExamples` (`Language`, `Comment`, `Code`) |
| **Pilar 6** | **5. Compliant Implementation Patterns (Good Examples)** | `doc.GoodExamples` (`Language`, `Comment`, `Code`) |
| **Pilar 7** | **6. How to Suppress (Ignore Directives)** | Dihasilkan dinamis untuk sintaks Astro (`<!-- charites:ignore -->`) & TSX (`// charites:ignore`) |
| **Pilar 8** | **7. Configuration Reference (`charites.yaml`)** | Dihasilkan dinamis dengan cuplikan konfigurasi YAML sesuai severity rule |

---

## 5. Arsitektur Engine Generator (`internal/wiki/`)

```text
internal/wiki/
├── generator.go           # Logika kompilasi, staging atomik & copyTree
├── generator_test.go      # Pengujian determinisme biner & regresi wiki (>= 90% coverage)
└── templates/             # Embedded templates via Go 1.26 embed.FS
    ├── home.md.tmpl       # Template untuk wiki/Home.md
    ├── category.md.tmpl   # Template untuk wiki/<category>.md
    └── rule.md.tmpl       # Template 8-Pillars untuk wiki/<category>/<slug>.md
```

### Invarian Engine Wiki:
1. **Determinis Penuh (Zero Churn):** Dilarang mencetak timestamp atau path absolut mesin build lokal. Pengeksekusian berulang menghasilkan output biner identik.
2. **Staging Atomik:** Seluruh berkas dirender ke direktori temporer (`os.MkdirTemp("", "charites-wiki-staging-*")`). Direktori tujuan `wiki/` hanya disinkronkan setelah semua template berhasil dirender tanpa error.
3. **Pengurutan Leksikografis (Total Ordering):** Kategori diurutkan alfabetis (`cat ASC`), dan rule diurutkan leksikografis berdasarkan `rule.ID() ASC`.

---

## 6. Standar Nada Tulisan (Voice & Tone)

- **Otoritatif & Berbasis Standar:** Gunakan terminologi resmi web platform (Fitts's Law, WCAG Contrast Minimum, Composite Layers, Cumulative Layout Shift, DTCG Tokens).
- **Edukatif & Siap Pakai:** Setiap contoh bad pattern wajib dipasangkan dengan good pattern yang siap pakai (*copy-pasteable*).
- **Ringkas & To-the-Point:** Hindari narasi berbunga-bunga; utamakan kejelasan teknis bagi software engineer.
