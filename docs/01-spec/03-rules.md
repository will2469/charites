# 01-SPEC: 03 - Rule Interface, Registry & Proving Ground Specification

> **Kode Dokumen:** `SPEC-03-RULES`
> **Tahapan:** Fase 3 - Rule Contract & Proving Ground Rule (`theme.hardcode-opacity-color`)
> **Status:** Ready for Review
> **Standar Rujukan:** IETF RFC 2119 / Semgrep Rule Design Specification

Dokumen ini mendefinisikan spesifikasi kebutuhan fungsional untuk **Interface Evaluasi Rule**, **In-Memory Registry**, dan implementasi rule pembuktian pertama (**`theme.hardcode-opacity-color`**) beserta standar pengujian semantik model Argus.

---

## 1. Kontrak Interface Rule (`internal/rules/rule.go`)

Seluruh modul audit pada Charites **MUST** mengimplementasikan interface murni berikut:

```go
package rules

import (
    "github.com/will2469/charites/internal/ir"
)

type Rule interface {
    ID() string                             // Semgrep ID tunggal (contoh: "theme.hardcode-opacity-color")
    Description() string                    // Penjelasan maksud dan tujuan rule
    Category() string                       // Kategori bidang: theme, a11y, perf, layout, seo
    DefaultSeverity() ir.Severity           // ir.SeverityError | ir.SeverityWarn | ir.SeverityInfo
    Evaluate(node *ir.Node) []ir.Diagnostic // Fungsi evaluasi murni tanpa efek samping I/O disk
}
```

### Invarian Evaluasi Murni:
- **No Side Effects:** Fungsi `Evaluate()` dilarang memodifikasi pointer pohon `*ir.Node` dan dilarang melakukan operasi disk atau jaringan.
- **Idempotent:** Pemanggilan `Evaluate(node)` berulang kali pada node yang sama **MUST** menghasilkan slice `[]ir.Diagnostic` yang identik.

---

## 2. Spesifikasi In-Memory Registry (`internal/rules/registry.go`)

Registry bertindak sebagai manajer katalog seluruh rule yang terdaftar di dalam binary:

```go
type Registry struct {
    rules map[string]Rule
    mu    sync.RWMutex
}

func NewRegistry() *Registry
func (r *Registry) Register(rule Rule) error
func (r *Registry) Get(id string) (Rule, bool)
func (r *Registry) All() []Rule
func (r *Registry) ByCategory(category string) []Rule
```

### Aturan Registrasi:
1. Pendaftaran ID rule yang duplikat **MUST** menghasilkan error registrasi.
2. Registry **MUST** mendukung lookup cepat $O(1)$ berdasarkan Semgrep ID string.
3. Fungsi `ByCategory()` mengembalikan seluruh rule aktif dalam kategori tertentu (contoh: `"theme"`).

---

## 3. Spesifikasi Rule #1: `theme.hardcode-opacity-color`

Sebagai bukti kerja compiler pipeline (*proving ground reference implementation*), rule pertama menegakkan kepatuhan token semantik warna dengan opacity:

### 3.1. Identitas Rule
- **ID:** `theme.hardcode-opacity-color`
- **Kategori:** `theme`
- **Default Severity:** `error` (`ir.SeverityError`)
- **Deskripsi:** Mendeteksi utility class warna dengan slash opacity langsung yang wajib diganti dengan token semantik dari `global.css`.

### 3.2. Pola Pelanggaran
Rule memeriksa setiap token di dalam `node.Classes`. Pelanggaran terjadi jika token cocok dengan pola:
```text
(bg|text|border|ring)-(primary|secondary|accent|destructive|warning|muted|amber|emerald)/\d{1,2}
```
dan pasangan token tersebut terdaftar dalam kamus `OPACITY_TOKEN_MAP`.

### 3.3. Kamus Token Pengganti (`OPACITY_TOKEN_MAP`)

| Utility dengan Slash | Token Semantik Resmi `global.css` |
| :--- | :--- |
| `primary/10`, `primary/20` | `primary-light` |
| `primary/5` | `primary-subtle` |
| `secondary/10` | `muted-light` |
| `secondary/5` | `muted-subtle` |
| `destructive/10`, `destructive/20` | `destructive-light` |
| `destructive/5` | `destructive-subtle` |
| `accent/10`, `accent/20` | `accent-light` |
| `accent/5` | `accent-subtle` |
| `warning/10` | `warning-light` |
| `warning/5` | `warning-subtle` |
| `muted/10` | `muted-light` |
| `muted/5` | `muted-subtle` |
| `amber/10` | `amber-light` |
| `amber/5` | `amber-subtle` |
| `emerald/10` | `emerald-light` |
| `emerald/5` | `emerald-subtle` |

### 3.4. Payload Temuan Diagnostic
Jika ditemukan pelanggaran pada `node`:
- `File`: Path berkas sumber.
- `Line`: `node.Span.Line`.
- `Column`: `node.Span.Column`.
- `Rule`: `"theme.hardcode-opacity-color"`.
- `Severity`: `ir.SeverityError`.
- `Message`: `"Hardcode opacity color (<class>) - wajib pakai semantic token dari global.css"`.
- `Hint`: `"Ganti dengan token semantik: primary/10 → primary-light, destructive/10 → destructive-light"`.

---

## 4. Spesifikasi Tri-Corpus Argus (`tests/correctness/theme.hardcode-opacity-color/`)

Pengujian semantik rule ini **MUST** dibagi ke dalam 3 sub-korpus:

1. **`positive/` (True Positives):**
   - Berkas memuat `bg-primary/10`, `border-destructive/20`.
   - **Ekspektasi:** `PositiveCount > 0`.
2. **`negative/` (True Negatives / Zero-Noise Invariant):**
   - Berkas memuat `bg-primary-light`, `border-destructive-light`, `text-muted`.
   - **Ekspektasi:** `NegativeCount == 0` (Wajib 0 pelanggaran).
3. **`adversarial/` (False Positive Bait & Evasion Tests):**
   - Berkas memuat utilitas slash non-warna: `w-1/2`, `aspect-16/9`, `grid-cols-2/3`.
   - Berkas dengan template literal dinamis kompleks: `` `p-4 ${isActive ? "bg-primary-light" : ""}` ``.
   - Berkas dengan inline ignore: `// charites:ignore theme.hardcode-opacity-color`.
   - **Ekspektasi:** Engine kebal dari false positive dan menghormati direktif inline ignore.
