# 01-SPEC: 03 - Rule Interface, Registry & Proving Ground Specification

> **Kode Dokumen:** `SPEC-03-RULES`
> **Tahapan:** Fase 3 - Rule Contract & Proving Ground Rule (`theme.hardcode-opacity-color`)
> **Peran Pilar:** SPEC = WHAT (Spesifikasi Kontrak Rule, Registry & Detection Domain)
> **Status:** Ready for Execution
> **Standar Rujukan:** IETF RFC 2119 / Deterministic Rule Engine Specification

Dokumen ini mendefinisikan spesifikasi kebutuhan fungsional untuk **Interface Evaluasi Rule**, **In-Memory Registry**, dan implementasi rule pembuktian pertama (**`theme.hardcode-opacity-color`**) beserta **Detection Contract** dan standar pengujian semantik model Argus.

---

## 1. Kontrak Interface Rule (`internal/rules/rule.go`)

Seluruh modul audit pada Charites **MUST** mengimplementasikan interface murni berikut:

```go
package rules

import (
    "github.com/will2469/charites/internal/ir"
)

type Rule interface {
    ID() string                             // Charites Rule ID tunggal (contoh: "theme.hardcode-opacity-color")
    Description() string                    // Penjelasan maksud dan tujuan rule
    Category() string                       // Kategori bidang: theme, a11y, perf, layout, seo
    DefaultSeverity() ir.Severity           // ir.SeverityError | ir.SeverityWarn | ir.SeverityInfo
    Evaluate(node *ir.Node) []ir.Diagnostic // Fungsi evaluasi murni tanpa efek samping I/O disk
}

// DocumentedRule adalah interface opsional untuk rule yang menyediakan dokumentasi kaya 8-Pillars.
type DocumentedRule interface {
    Rule
    Doc() ir.RuleDocumentation
}
```

### Invarian Evaluasi Murni:
- **No Side Effects:** Fungsi `Evaluate()` dilarang memodifikasi pointer pohon `*ir.Node` dan dilarang melakukan operasi disk atau jaringan.
- **Idempotent:** Pemanggilan `Evaluate(node)` berulang kali pada node yang sama **MUST** menghasilkan slice `[]ir.Diagnostic` yang identik.
- **Rule Independence from Comments:** Fungsi `Evaluate()` murni mendeteksi pelanggaran pada `node`. Penyaringan komentar direktif (`charites:ignore`) dilakukan oleh lapisan traversal analyzer di Fase 4 (*separation of concerns*).

### Invarian Single Source of Truth (SSOT) Dokumentasi:
- Setiap rule yang diimplementasikan wajib menyertakan method `Doc() ir.RuleDocumentation` (`internal/ir/doc.go`).
- **Dilarang Menulis Wiki Manual:** Seluruh dokumentasi ensiklopedia di `wiki/` (`Home.md`, `<category>.md`, `<category>/<slug>.md`) dihasilkan secara otomatis via `make wiki` menggunakan paket `internal/wiki`.


---

## 2. Spesifikasi In-Memory Registry (`internal/rules/registry.go`)

Registry bertindak sebagai manajer katalog seluruh rule yang terdaftar di dalam binary:

```go
type Registry struct {
    rules      map[string]Rule
    categories map[string][]Rule
    mu         sync.RWMutex
}

func NewRegistry() *Registry
func (r *Registry) Register(rule Rule) error
func (r *Registry) Get(id string) (Rule, bool)
func (r *Registry) All() []Rule
func (r *Registry) ByCategory(category string) []Rule
```

### Aturan Registrasi & Determinisme Urutan:
1. Pendaftaran ID rule yang duplikat **MUST** menghasilkan error registrasi.
2. Registry **MUST** mendukung lookup cepat $O(1)$ berdasarkan Charites Rule ID string.
3. **Deterministic Ordering:** Metode `All()` dan `ByCategory()` **MUST** mengembalikan slice rule dalam urutan yang deterministik (diurutkan secara leksikografis berdasarkan `Rule.ID()`), untuk mencegah inkonsistensi urutan eksekusi akibat iterasi map bawaan Go.

---

## 3. Spesifikasi Rule #1: `theme.hardcode-opacity-color`

Sebagai bukti kerja compiler pipeline (*proving ground reference implementation*), rule pertama menegakkan kepatuhan token semantik warna dengan opacity yang memiliki pemetaan pengganti resmi:

### 3.1. Identitas Rule
- **ID:** `theme.hardcode-opacity-color`
- **Kategori:** `theme`
- **Default Severity:** `error` (`ir.SeverityError`)
- **Deskripsi:** Mendeteksi utility class warna dengan slash opacity langsung yang memiliki pemetaan token semantik resmi dari `global.css`.

---

### 3.2. Detection Contract (Tabel Kontrak Deteksi)

Prinsip Deteksi: **"Slash opacity yang memiliki semantic token replacement adalah pelanggaran (violation)."**

#### A. Cakupan Masuk (IN-SCOPE):
1. **Utility Families:**
   - `bg-`
   - `text-`
   - `border-`
   - `ring-`
2. **Named Semantic Colors:**
   `primary`, `secondary`, `accent`, `destructive`, `warning`, `muted`, `amber`, `emerald`
3. **Dukungan Tailwind Variants (Lexical Normalization):**
   Rule **MUST** melakukan normalisasi leksikal membuang prefix variant tunggal maupun bersarang sebelum mencocokkan utility dasar:
   - Single variant: `hover:bg-primary/10`, `dark:text-primary/20`, `focus:ring-warning/10`
   - Chained variants: `md:hover:bg-primary/10`, `dark:hover:border-destructive/20`, `sm:dark:bg-primary/5`
4. **Supported Opacity Mappings (`OPACITY_TOKEN_MAP`):**
   Hanya pasangan warna dan angka opacity berikut yang memicu pelanggaran:

| Base Utility Pattern | Replacement Semantic Token |
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

#### B. Di Luar Cakupan (OUT-OF-SCOPE):
Pola-pola berikut **BUKAN** pelanggaran Rule #1 dan **DILARANG** menghasilkan diagnostic:
1. **Layout & Dimension Fractions:** `w-1/2`, `w-2/3`, `h-1/3`, `h-2/4`, `max-w-1/2` (pecahan dimensi lebar/tinggi).
2. **Aspect Ratio:** `aspect-16/9`, `aspect-4/3`, `aspect-1/1` (rasio aspek CSS).
3. **Grid Column Fractions:** `grid-cols-2/3` (fraksi template grid).
4. **Typography Line-Height:** `text-sm/6`, `text-xs/relaxed` (modifier line-height Tailwind).
5. **Arbitrary Color Values:** `bg-[#123456]/10`, `text-[#ff0000]/20` (karena tidak memiliki token pengganti deterministik; ditangani rule terpisah).
6. **Non-Mapped Opacities:** `bg-primary/30`, `bg-primary/50`, `bg-primary/100`, `bg-primary/[0.1]` (tidak memiliki pemetaan semantik resmi pada default theme).
7. **Raw Palette Colors:** `bg-red-500/10` (dikelola oleh rule terpisah: `theme.hardcode-palette-color`).
8. **Dynamic Unresolved Classes:** Ekspresi runtime yang belum terevaluasi.

---

### 3.3. Payload Temuan Diagnostic (Dynamic Generation)

Jika ditemukan pelanggaran pada `node`:
- `File`: Path berkas sumber.
- `Line`: `node.Span.Line`.
- `Column`: `node.Span.Column`.
- `Rule`: `"theme.hardcode-opacity-color"`.
- `Severity`: `ir.SeverityError`.
- `Message`: Dihasilkan dinamis: `"Hardcode opacity color: \"<matched_class>\""`
  (contoh: `Hardcode opacity color: "bg-primary/10"` atau `Hardcode opacity color: "hover:bg-primary/10"`).
- `Hint`: Dihasilkan dinamis dari mapping: `"Use semantic token \"<replacement_token>\"."`
  (contoh: `Use semantic token "primary-light".` atau `Use semantic token "destructive-subtle".`).

---

## 4. Spesifikasi Argus Tri-Corpus (`tests/correctness/theme/hardcode-opacity-color/`)

Pengujian semantik rule ini **MUST** dibagi ke dalam 3 sub-korpus terpisah:

1. **`positive/` (True Positives):**
   - Berkas memuat `bg-primary/10`, `hover:bg-primary/10`, `dark:border-destructive/20`.
   - **Ekspektasi:** `PositiveCount > 0`.
2. **`negative/` (True Negatives / Zero-Noise Invariant):**
   - Berkas memuat `bg-primary-light`, `border-destructive-light`, `text-muted`, `bg-primary/30`, `w-1/2`.
   - **Ekspektasi:** `NegativeCount == 0` (Wajib 0 pelanggaran).
3. **`adversarial/` (False Positive Bait & Evasion Tests):**
   - Berkas memuat `aspect-16/9`, `text-sm/6`, `grid-cols-2/3`, `bg-[#123456]/10`.
   - **Ekspektasi:** `AdversarialViolations == 0`.

