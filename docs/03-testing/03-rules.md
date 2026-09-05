# 03-TESTING: 03 - Rule Correctness & Tri-Corpus Semantic Verification Plan

> **Kode Dokumen:** `TEST-03-RULES`
> **Tahapan:** Fase 3 - Rule Contract & Proving Ground Rule (`theme.hardcode-opacity-color`)
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Skenario Smoke & Asersi Pembuktian)
> **Status:** Graduated (All Phase Gates Passed)
> **Standar Rujukan:** Argus Tri-Corpus Semantic Verification & Table-Driven Unit Testing

Dokumen ini mendefinisikan strategi pengujian ketat untuk modul rule Charites, mencakup **Argus Tri-Corpus Semantic Verification**, matriks pembuktian batas deteksi (*boundary matrix*), dan pengujian determinisme registri in-memory.

---

## 1. Metodologi Argus Tri-Corpus & Pemisahan Tanggung Jawab

Setiap rule diuji terhadap 3 sub-korpus di bawah direktori `tests/correctness/<category>/<slug>/`:

```text
tests/correctness/theme/hardcode-opacity-color/
├── positive/      # True Positives: Kasus pelanggaran nyata (termasuk varian Tailwind)
│   ├── basic.astro
│   └── variants.tsx
├── negative/      # True Negatives: Kasus legal (Zero Noise Invariant)
│   ├── semantic.astro
│   └── clean.tsx
└── adversarial/   # False Positive Bait: Jebakan slash non-warna & out-of-scope utilities
    ├── slash_layout.astro
    ├── line_height.tsx
    ├── arbitrary_color.tsx
    └── unmapped_opacity.astro
```

### 1.1. Pemisahan Pengujian (Rule Correctness vs Suppression):
1. **Rule Correctness:** Menguji fungsi murni `rule.Evaluate(node)` terhadap token kelas untuk memastikan hanya pola yang memiliki pengganti semantik yang dilaporkan.
2. **Suppression Correctness:** Pengujian pemfilteran komentar `charites:ignore` diisolasi pada layer engine analyzer (Fase 4), memastikan rule murni bebas dari dependensi parser komentar.

### 1.2. Metrik Kelulusan Evaluasi (`RuleCorrectnessMetric`):
```text
Pass = (PositiveViolations > 0) && (NegativeViolations == 0) && (AdversarialViolations == 0)
```

---

## 2. Matriks Uji Lengkap (Boundary & Correctness Matrix)

Matriks berikut mengunci batas fungsional Rule #1:

| Pola Input Class | Klasifikasi | Ekspektasi Deteksi | Ekspektasi Diagnostic / Hint |
| :--- | :---: | :---: | :--- |
| `bg-primary/10` | In-Scope |  Violation | Hint: `Use semantic token "primary-light".` |
| `text-primary/20` | In-Scope |  Violation | Hint: `Use semantic token "primary-light".` |
| `border-destructive/20` | In-Scope |  Violation | Hint: `Use semantic token "destructive-light".` |
| `ring-warning/10` | In-Scope |  Violation | Hint: `Use semantic token "warning-light".` |
| `bg-primary/5` | In-Scope |  Violation | Hint: `Use semantic token "primary-subtle".` |
| `hover:bg-primary/10` | In-Scope (Variant) |  Violation | Hint: `Use semantic token "primary-light".` |
| `dark:bg-primary/10` | In-Scope (Variant) |  Violation | Hint: `Use semantic token "primary-light".` |
| `md:hover:bg-primary/10` | In-Scope (Chained) |  Violation | Hint: `Use semantic token "primary-light".` |
| `bg-primary-light` | Negative (Clean) |  Pass (0 diag) | Token semantik resmi |
| `bg-primary-subtle` | Negative (Clean) |  Pass (0 diag) | Token semantik resmi |
| `text-muted` | Negative (Clean) |  Pass (0 diag) | Token semantik resmi |
| `w-1/2`, `h-1/3`, `max-w-1/2` | Out-of-Scope |  Pass (0 diag) | Pecahan dimensi tata letak |
| `aspect-16/9`, `aspect-4/3` | Out-of-Scope |  Pass (0 diag) | Rasio aspek CSS |
| `grid-cols-2/3` | Out-of-Scope |  Pass (0 diag) | Fraksi template grid |
| `text-sm/6`, `text-xs/relaxed` | Out-of-Scope |  Pass (0 diag) | Modifier line-height typography |
| `bg-primary/30`, `bg-primary/50` | Out-of-Scope |  Pass (0 diag) | Opacity tanpa pemetaan token resmi |
| `bg-primary/100`, `bg-primary/[0.1]` | Out-of-Scope |  Pass (0 diag) | Format nilai arbitrary / non-mapped |
| `bg-[#123456]/10` | Out-of-Scope |  Pass (0 diag) | Arbitrary hex color (tanpa token pengganti) |
| `bg-red-500/10` | Out-of-Scope |  Pass (0 diag) | Raw palette color (milik rule terpisah) |
| `bg-black/10` | Out-of-Scope |  Pass (0 diag) | Raw black color (milik rule terpisah) |

---

## 3. Pengujian In-Memory Registry (`internal/rules/registry_test.go`)

Registry diuji terhadap thread-safety dan **stabilitas urutan deterministik**:

```go
func TestRegistry_DeterministicOrder(t *testing.T) {
    reg := rules.NewRegistry()
    // Daftarkan rule secara acak
    _ = reg.Register(&mockRule{id: "theme.hardcode-palette-color"})
    _ = reg.Register(&mockRule{id: "a11y.missing-alt"})
    _ = reg.Register(&mockRule{id: "theme.hardcode-opacity-color"})

    all := reg.All()
    if len(all) != 3 {
        t.Fatalf("expected 3 rules, got %d", len(all))
    }

    // Wajib terurut secara leksikografis
    expected := []string{
        "a11y.missing-alt",
        "theme.hardcode-opacity-color",
        "theme.hardcode-palette-color",
    }
    for i, r := range all {
        if r.ID() != expected[i] {
            t.Errorf("urutan indeks %d tidak sesuai: dapat %s, ingin %s", i, r.ID(), expected[i])
        }
    }
}
```

---

## 4. Benchmark Kinerja Evaluasi Rule

Kecepatan evaluasi diukur untuk memastikan fungsi tidak memicu alokasi heap saat memeriksa node bersih:

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

### Evaluasi Target (Performance Budget):
- **Node Bersih:** Target `0 B/op` dan `0 allocs/op`.
- **Waktu Eksekusi:** Target evaluasi berada dalam rentang wajar nanodetik pada runner target.

