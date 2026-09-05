# 02-ARCHITECTURE: 08 - Repetitive Architecture Flow: Shared Components vs Rule-Specific Logic

> **Kode Dokumen:** `ARCH-08-EXPANSION`
> **Tahapan:** Fase 8 - Repetitive Pattern Flow Guide & Rule Authoring Template (Core Assessment)
> **Peran Pilar:** ARCH = HOW (Rancangan Arsitektur Ekstensibilitas, Rule SSOT & 3 Touchpoints)
> **Status:** Ready for Review
> **Standar Rujukan:** Micro-Kernel Plugin Architecture & Open-Closed Principle (OCP)

Dokumen ini mendefinisikan arsitektur penambahan rule secara berulang (*repetitive architecture flow*), menguraikan isolasi mutlak **komponen bersama yang dibekukan (*Frozen Shared Components*)** serta menyatukan definisi metadata rule ke dalam **Single Source of Truth (SSOT)**.

---

## 1. Pemetaan Arsitektur: Komponen Bersama vs Titik Sentuh Rule

Sesuai prinsip *Open-Closed Principle* (terbuka untuk ekstensi, tertutup untuk modifikasi), pipa compiler inti yang telah dibekukan pada Fase 6 **DILARANG KERAS** dimodifikasi saat menambahkan rule baru:

### 1.1. Matriks Komponen Bersama (100% Shared - Immutable / Nol Perubahan Kode)

| Komponen | Status | Peran & Alasan Arsitektural |
| :--- | :---: | :--- |
| **`internal/ir`** | **FROZEN** | Kontrak data `*ir.Node`, `Diagnostic`, `Span`, dan `Severity` telah final. |
| **`internal/parser`** | **FROZEN** | Parser Astro, TSX, dan Tailwind `@theme` mengekstrak representasi markup secara netral. |
| **`internal/scanner`** | **FROZEN** | Directory walker dan worker pool mendistribusikan berkas secara generik tanpa mengenal jenis rule. |
| **`internal/config`** | **FROZEN** | Parser `charites.yaml` bekerja otomatis dengan model Default: YES. Rule baru otomatis aktif. |
| **`internal/analyzer`** | **FROZEN** | Traversal engine menjalankan iterator `root.Walk()` dan memanggil `Evaluate()` secara generik. |
| **`internal/reporter`** | **FROZEN** | Presenter ANSI terminal dan JSON merender slice `ir.Diagnostic` apa saja tanpa perubahan kode. |
| **`internal/mcp`** | **FROZEN** | Tool `charites_scan`, `charites_explain_rule`, dan `charites_list_rules` membaca katalog rule dari registry secara dinamis. |
| **`internal/wiki`** | **FROZEN** | Generator ensiklopedia otomatis mengekstrak metadata rule baru dari registry ke `wiki/*.md`. |

### 1.2. Tiga Titik Sentuh Spesifik Rule (The 3 Touchpoints)

Untuk menambahkan rule baru, pengembang **HANYA** menyentuh 3 lokasi terisolasi:
1. **`internal/rules/<domain>/<rule_slug>.go`** (*Berkas Baru*): Mengimplementasikan interface `rules.Rule` beserta metadatanya.
2. **`internal/rules/registry.go`** (*Modifikasi 1 Baris*): Mendaftarkan instance singleton rule ke registry in-memory.
3. **`tests/correctness/<category>.<slug>/`** (*Direktori Baru*): Menyediakan berkas uji Tri-Corpus (`positive/`, `negative/`, `adversarial/`) dan matriks kasus ekspektasi.

---

## 2. Arsitektur Single Source of Truth (SSOT) Metadata Rule

Agar `charites_explain_rule` pada MCP dan generator `charites wiki` tidak memerlukan pembaruan kode saat ada rule baru, interface `Rule` menyediakan metadata terpadu:

```go
package rules

import "github.com/will2469/charites/internal/ir"

type RuleMetadata struct {
    Explanation string // Alasan arsitektural dan dampak buruk pelanggaran
    BadExample  string // Contoh potongan kode yang melanggar
    GoodExample string // Contoh potongan kode rekomendasi yang benar
    Remediation string // Panduan remedi spesifik atau padanan token semantik
}

type Rule interface {
    ID() string
    Category() string
    DefaultSeverity() ir.Severity
    Description() string
    Metadata() RuleMetadata
    Evaluate(node *ir.Node) []ir.Diagnostic
}
```

Dengan desain ini, baik server MCP maupun generator Wiki murni bertindak sebagai konsumen (*consumers*) dari metadata rule yang didefinisikan secara mandiri oleh masing-masing rule.

---

## 3. Template Struktur Kode Baku Rule (`<rule_slug>.go`)

```go
package theme

import (
    "strings"
    "github.com/will2469/charites/internal/ir"
    "github.com/will2469/charites/internal/rules"
)

type HardcodeColorRule struct{}

func NewHardcodeColorRule() rules.Rule {
    return &HardcodeColorRule{}
}

func (r *HardcodeColorRule) ID() string {
    return "theme.hardcode-color"
}

func (r *HardcodeColorRule) Category() string {
    return "theme"
}

func (r *HardcodeColorRule) DefaultSeverity() ir.Severity {
    return ir.SeverityWarn
}

func (r *HardcodeColorRule) Description() string {
    return "Mendeteksi penggunaan kode warna mentah yang wajib diganti dengan token semantik dari global.css"
}

func (r *HardcodeColorRule) Metadata() rules.RuleMetadata {
    return rules.RuleMetadata{
        Explanation: "Penggunaan warna mentah (hex/rgb) merusak konsistensi dark mode dan mempersulit rebranding desain sistem.",
        BadExample:  `<div className="bg-[#2563eb] text-[#000]" />`,
        GoodExample: `<div className="bg-primary text-foreground" />`,
        Remediation: "Ganti warna mentah dengan utility token semantik dari global.css (misal: bg-primary, text-muted).",
    }
}

func (r *HardcodeColorRule) Evaluate(node *ir.Node) []ir.Diagnostic {
    // Fast path: jika node tidak memuat class atau inline style, return nil (0 B/op)
    if len(node.Classes) == 0 && len(node.Attributes["style"]) == 0 {
        return nil
    }

    var diags []ir.Diagnostic
    for _, class := range node.Classes {
        // Pengecekan cepat karakter '#' atau 'rgb'
        if strings.Contains(class, "#") || strings.Contains(class, "rgb") {
            if isRawColorClass(class) {
                diags = append(diags, ir.Diagnostic{
                    File:     node.Span.File,
                    Line:     node.Span.Line,
                    Column:   node.Span.Column,
                    Rule:     r.ID(),
                    Category: r.Category(),
                    Severity: r.DefaultSeverity(),
                    Message:  "Hardcode color: \"" + class + "\"",
                    Hint:     r.Metadata().Remediation,
                })
            }
        }
    }
    return diags
}
```
