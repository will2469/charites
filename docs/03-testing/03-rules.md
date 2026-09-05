# 03-TESTING: 03 - Rule Correctness & Tri-Corpus Semantic Verification Plan

> **Kode Dokumen:** `TEST-03-RULES`
> **Tahapan:** Fase 3 - Rule Contract & Proving Ground Rule (`theme.hardcode-opacity-color`)
> **Status:** Ready for Review
> **Standar Rujukan:** Argus Tri-Corpus Semantic Verification & Table-Driven Unit Testing

Dokumen ini mendefinisikan strategi pengujian ketat untuk modul rule Charites, mencakup **Argus Tri-Corpus Semantic Verification**, verifikasi regresi unit testing, dan pengujian konkurensi registri in-memory.

---

## 1. Metodologi Argus Tri-Corpus Semantic Verification

Setiap rule yang diimplementasikan pada Charites wajib lolos uji 3 sub-korpus terpisah di bawah direktori `tests/correctness/<rule_id>/`:

```text
tests/correctness/theme.hardcode-opacity-color/
├── positive/      # True Positives: Kasus pelanggaran nyata
│   ├── basic.astro
│   └── complex.tsx
├── negative/      # True Negatives: Kasus legal (Zero Noise Invariant)
│   ├── semantic.astro
│   └── clean.tsx
└── adversarial/   # False Positive Bait: Pola jebakan mirip tapi sah
    ├── slash_layout.astro
    ├── line_height.tsx
    └── ignored.astro
```

### 1.1. Metrik Kelulusan Evaluasi (`RuleCorrectnessMetric`)
Sebuah rule dinyatakan **CORRECT & VERIFIED** jika dan hanya jika:
```text
Pass = (PositiveViolations > 0) && (NegativeViolations == 0) && (AdversarialViolations == 0)
```
- **Zero Noise Invariant:** `NegativeViolations == 0`. Tidak boleh ada peringatan/error palsu pada kode yang mematuhi pedoman.
- **Bait Immunity Invariant:** `AdversarialViolations == 0`. Kode yang menggunakan utilitas slash non-warna atau diberi inline ignore tidak boleh memicu diagnostic.

---

## 2. Test Harness Otomatis Tri-Corpus (`tests/correctness_gate_test.go`)

Runner pengujian memuat AST dan mengevaluasi rule terhadap ketiga sub-korpus secara otomatis:

```go
package tests

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/will2469/charites/internal/parser/astro"
    "github.com/will2469/charites/internal/parser/tsx"
    "github.com/will2469/charites/internal/rules"
    "github.com/will2469/charites/internal/rules/theme"
)

func TestTriCorpus_ThemeHardcodeOpacityColor(t *testing.T) {
    rule := theme.NewHardcodeOpacityColorRule()
    baseDir := filepath.Join("correctness", rule.ID())

    // 1. Verifikasi Positive Corpus (Wajib mendeteksi pelanggaran)
    t.Run("Positive_Corpus", func(t *testing.T) {
        files := loadTestFiles(t, filepath.Join(baseDir, "positive"))
        totalDiags := runRuleOnFiles(t, rule, files)
        if totalDiags == 0 {
            t.Fatalf("Rule %s GAGAL mendeteksi pelanggaran pada positive corpus", rule.ID())
        }
    })

    // 2. Verifikasi Negative Corpus (Wajib 0 pelanggaran - Zero Noise)
    t.Run("Negative_Corpus", func(t *testing.T) {
        files := loadTestFiles(t, filepath.Join(baseDir, "negative"))
        totalDiags := runRuleOnFiles(t, rule, files)
        if totalDiags != 0 {
            t.Fatalf("Rule %s menghasilkan false positive (%d temuan) pada negative corpus", rule.ID(), totalDiags)
        }
    })

    // 3. Verifikasi Adversarial Corpus (Wajib 0 pelanggaran - Bait Immunity)
    t.Run("Adversarial_Corpus", func(t *testing.T) {
        files := loadTestFiles(t, filepath.Join(baseDir, "adversarial"))
        totalDiags := runRuleOnFiles(t, rule, files)
        if totalDiags != 0 {
            t.Fatalf("Rule %s termakan jebakan false positive (%d temuan) pada adversarial corpus", rule.ID(), totalDiags)
        }
    })
}
```

---

## 3. Matriks Kasus Uji `theme.hardcode-opacity-color`

### 3.1. Positive Test Cases
| Input Class | Deteksi Pelanggaran | Rekomendasi Hint |
| :--- | :---: | :--- |
| `bg-primary/10` | Ya | `primary-light` |
| `border-destructive/20` | Ya | `destructive-light` |
| `text-accent/5` | Ya | `accent-subtle` |
| `ring-warning/10` | Ya | `warning-light` |
| `bg-secondary/10` | Ya | `muted-light` |

### 3.2. Negative Test Cases (Valid & Sah)
| Input Class | Alasan Valid |
| :--- | :--- |
| `bg-primary-light` | Token semantik resmi dari `global.css` |
| `border-destructive-light` | Token semantik resmi dari `global.css` |
| `text-muted-subtle` | Token semantik resmi dari `global.css` |
| `p-4 flex flex-col gap-2` | Utilitas tata letak standar tanpa relasi warna |

### 3.3. Adversarial Test Cases (Jebakan Slash Non-Warna)
| Input Class | Sifat Jebakan | Ekspektasi |
| :--- | :--- | :---: |
| `w-1/2`, `h-1/3`, `w-2/3` | Slash digunakan untuk rasio lebar/tinggi pecahan | **Abaikan (0 diag)** |
| `aspect-16/9`, `aspect-4/3` | Slash digunakan untuk rasio aspek | **Abaikan (0 diag)** |
| `grid-cols-2/3` | Slash digunakan untuk grid template column | **Abaikan (0 diag)** |
| `text-xs/relaxed`, `text-sm/6` | Slash digunakan untuk modifier line-height Tailwind | **Abaikan (0 diag)** |
| `bg-primary/10` + inline ignore | Terdapat comment `// charites:ignore theme.hardcode-opacity-color` | **Abaikan (0 diag)** |

---

## 4. Pengujian Registri In-Memory (`internal/rules/registry_test.go`)

Paket registri diuji terhadap keandalan konkurensi dan validitas integritas data:

```go
func TestRegistry_Concurrency(t *testing.T) {
    reg := rules.NewRegistry()
    rule := theme.NewHardcodeOpacityColorRule()
    _ = reg.Register(rule)

    const workers = 50
    var wg sync.WaitGroup
    wg.Add(workers)

    for i := 0; i < workers; i++ {
        go func() {
            defer wg.Done()
            // Akses serentak get & list dilarang race condition atau panic
            r, ok := reg.Get("theme.hardcode-opacity-color")
            if !ok || r == nil {
                t.Errorf("expected rule found")
            }
            all := reg.All()
            if len(all) == 0 {
                t.Errorf("expected non-empty rules")
            }
        }()
    }
    wg.Wait()
}
```

---

## 5. Benchmark Kinerja Evaluasi Rule

Pengujian kecepatan evaluasi node wajib diukur menggunakan tool benchmark Go:

```go
func BenchmarkEvaluateHardcodeOpacityColor_Clean(b *testing.B) {
    rule := theme.NewHardcodeOpacityColorRule()
    node := &ir.Node{
        Tag:     "div",
        Classes: []string{"p-4", "flex", "items-center", "justify-between", "rounded-lg"},
    }

    b.ResetTimer()
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = rule.Evaluate(node)
    }
}
```

### Ambang Batas Kinerja:
- **Zero Allocations:** `0 B/op` dan `0 allocs/op` saat mengevaluasi node bersih tanpa pelanggaran.
- **Throughput:** Waktu eksekusi $\le 50\text{ ns/op}$ per node.
