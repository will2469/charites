# 03-TESTING: 04 - Configuration, Concurrency Scanner & Traversal Engine Verification Plan

> **Kode Dokumen:** `TEST-04-ENGINE`
> **Tahapan:** Fase 4 - Konfigurasi, Concurrency Scanner & Traversal Engine
> **Status:** Ready for Review
> **Standar Rujukan:** Go Concurrency Testing, Stress Benchmarks & Race Detection Standards

Dokumen ini mendefinisikan strategi pengujian komprehensif untuk paket konfigurasi (`internal/config`), pemindai direktori paralel (`internal/scanner`), dan mesin traversal AST (`internal/analyzer`), termasuk pengujian *data race* dan *inline ignore suppression*.

---

## 1. Skenario Pengujian Unit Konfigurasi & Ignore (`internal/config/`)

### 1.1. Uji Coba Konfigurasi Default YES (`internal/config/config_test.go`)
- **Test Case 1 (Nil Config / Absent File):**
  - Kondisi: Berkas `charites.yaml` tidak ada pada root direktori.
  - Ekspektasi: `ResolveActiveRules()` mengembalikan seluruh rule yang terdaftar di registry tanpa ada yang terpotong.
- **Test Case 2 (Rule Override - Off):**
  - Input: YAML memuat `theme.hardcode-opacity-color: off`.
  - Ekspektasi: Rule tersebut tidak masuk ke dalam slice `[]rules.Rule` hasil resolusi.
- **Test Case 3 (Severity Override):**
  - Input: YAML memuat `theme.hardcode-color: warn`.
  - Ekspektasi: Rule tetap aktif, namun diagnosis yang dihasilkan menggunakan severity `warn`.

### 1.2. Uji Coba Engine Ignore (`internal/config/ignore_test.go`)
- **Test Case 1 (Builtin Pruning):**
  - Path: `node_modules/lodash/index.js`, `.git/HEAD`, `dist/bundle.js`.
  - Ekspektasi: `ShouldIgnoreDir()` mengembalikan `true` secara instan.
- **Test Case 2 (Custom Glob Patterns):**
  - Pola: `legacy-vendor/**`, `*.test.tsx`, `temp_*`.
  - Ekspektasi: Seluruh berkas yang cocok diabaikan secara akurat.
- **Test Case 3 (Negation Rules):**
  - Pola: `vendor/**` diikuti `!vendor/special.tsx`.
  - Ekspektasi: `vendor/foo.tsx` diabaikan, namun `vendor/special.tsx` diizinkan lewat.

---

## 2. Skenario Pengujian Concurrency Scanner (`internal/scanner/`)

### 2.1. Uji Fast Directory Walker (`internal/scanner/walker_test.go`)
- **Test Case 1 (Direct File Targeting - Ergonomi A):**
  - Input: Path langsung ke file `src/Button.tsx`.
  - Ekspektasi: Channel pekerjaan menerima tepat 1 path file tanpa melakukan `ReadDir` pada folder lain.
- **Test Case 2 (Extension Filtering - Ergonomi B):**
  - Input: Direktori memuat `.astro`, `.tsx`, `.json`, `.md`. Flag `--ext=astro`.
  - Ekspektasi: Hanya berkas `.astro` yang dikirim ke channel pekerjaan.
- **Test Case 3 (Early Directory Pruning Verification):**
  - Kondisi: Buat folder tiruan `node_modules` berisi 5.000 file dummy di direktori temporer (`t.TempDir()`).
  - Ekspektasi: Pemindaian selesai dalam waktu $< 2\text{ ms}$ karena folder `node_modules` sama sekali tidak dibaca isinya.

### 2.2. Uji Beban Concurrency Worker Pool (`internal/scanner/pool_test.go`)
- **Uji Stres 1.000 Berkas:**
  - Mengirim 1.000 file mock ke worker pool yang dikonfigurasi dengan $N = \text{runtime.NumCPU()}$ goroutine.
  - Memastikan seluruh berkas selesai diproses dan channel tertutup bersih tanpa ada goroutine yang menggantung (*leak-free*).
- **Deteksi Race Condition:**
  - Menjalankan seluruh pengujian scanner dengan flag balap:
    ```bash
    go test -race -v ./internal/scanner/...
    ```
  - Ekspektasi: **0 data race detected**.

---

## 3. Skenario Pengujian Traversal Analyzer Engine (`internal/analyzer/`)

### 3.1. Uji Inline Ignore Suppression (`internal/analyzer/context_test.go`)
- **Test Case 1 (Same-Line Directive):**
  - Kode sumber:
    ```tsx
    <div className="bg-primary/10">Test</div> // charites:ignore theme.hardcode-opacity-color
    ```
  - Ekspektasi: Diagnostic pada baris tersebut disaring (hasil akhir 0 diagnostic).
- **Test Case 2 (Next-Line Directive):**
  - Kode sumber:
    ```astro
    <!-- charites:ignore theme.hardcode-opacity-color -->
    <div class="bg-primary/10">Hero</div>
    ```
  - Ekspektasi: Diagnostic pada baris tepat di bawah comment disaring.
- **Test Case 3 (Multi-Rule Directive):**
  - Komentar: `// charites:ignore rule.a, rule.b`.
  - Ekspektasi: Baik pelanggaran `rule.a` maupun `rule.b` pada baris tersebut ditekan.

### 3.2. Uji Determinisme Pengurutan Diagnostic (`internal/analyzer/engine_test.go`)
- Jalankan analisis pada AST dengan beberapa pelanggaran yang terdistribusi acak.
- Ekspektasi: Hasil akhir selalu terurut stabil: `File` $\rightarrow$ `Line` $\rightarrow$ `Column` $\rightarrow$ `Rule`.

---

## 4. Benchmark Throughput Engine & Scanner

```go
func BenchmarkScannerAndEngineThroughput(b *testing.B) {
    // Siapkan 500 file AST fixtures di memori
    fixtures := prepareMockASTFiles(500)
    engine := analyzer.NewEngine()
    rules := loadActiveRules()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = runBenchmarkScan(engine, fixtures, rules)
    }
}
```

### Ambang Batas Throughput:
- Mampu memproses $\ge 5.000$ file per detik per core CPU SSD modern.
