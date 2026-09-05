# 03-TESTING: 04 - Configuration, Concurrency Scanner & Traversal Engine Verification Plan

> **Kode Dokumen:** `TEST-04-ENGINE`
> **Tahapan:** Fase 4 - Konfigurasi, Concurrency Scanner & Traversal Engine
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Pembuktian Batas & Asersi Keamanan)
> **Status:** Ready for Review
> **Standar Rujukan:** Go Concurrency Testing, Stress Benchmarks & Race Detection Standards

Dokumen ini mendefinisikan strategi pengujian menyeluruh untuk memvalidasi paket konfigurasi (`internal/config`), pemindai direktori paralel (`internal/scanner`), serta mesin traversal AST (`internal/analyzer`).

---

## 1. Skenario Pengujian Paket Konfigurasi (`internal/config/`)

### 1.1. Uji Resolusi & Presedensi Rule (`internal/config/config_test.go`)

- **Test Case 1 (Default YES Invariant):**
  - Kondisi: Berkas `charites.yaml` absen (`nil`).
  - Ekspektasi: `ResolveActiveRules()` mengembalikan seluruh rule di registry dengan severity default masing-masing.
- **Test Case 2 (Severity Override):**
  - Input: YAML memuat `theme.hardcode-opacity-color: warn`.
  - Ekspektasi: Rule tetap aktif dalam `ActiveRule`, namun `EffectiveSeverity` berubah menjadi `ir.SeverityWarn`.
- **Test Case 3 (Presedensi CLI Scope vs Config Policy):**
  - Kondisi: Flag CLI `--rule=theme.hardcode-opacity-color`, namun config menetapkan `theme.hardcode-opacity-color: off`.
  - Ekspektasi: Rule **TIDAK AKTIF** (kebijakan config mengalahkan seleksi CLI).

### 1.2. Uji Matcher Ignore & Semantik Negasi (`internal/config/ignore_test.go`)

- **Test Case 1 (Builtin Hard Exclusion Immunity):**
  - Input `.charitesignore`: `!node_modules/my-pkg/**`.
  - Ekspektasi: `node_modules` tetap diabaikan (`ShouldIgnoreDir == true`). Negasi tidak dapat membuka direktori builtin.
- **Test Case 2 (Sequential Evaluation & Last-Rule-Wins):**
  - Input:
    ```gitignore
    vendor/**
    !vendor/special/**
    vendor/special/private/**
    ```
  - Ekspektasi:
    - `vendor/other.tsx` $\rightarrow$ Ignored.
    - `vendor/special/public.tsx` $\rightarrow$ Allowed.
    - `vendor/special/private/secret.tsx` $\rightarrow$ Ignored.

---

## 2. Skenario Pengujian Concurrency Scanner (`internal/scanner/`)

### 2.1. Uji Fast Directory Walker (`internal/scanner/walker_test.go`)

- **Test Case 1 (Direct Target Safety):**
  - Input: `charites scan node_modules/react/index.d.ts`.
  - Ekspektasi: Walker mendeteksi path target memiliki leluhur terlarang via `HasBuiltinAncestor(target)` dan mengembalikan error validasi sebelum traversal (0 jobs diantrekan).
- **Test Case 2 (Symlink Safety Guard):**
  - Kondisi: Buat symlink siklis antar-folder temporer (`t.TempDir()`).
  - Ekspektasi: Walker melewati symlink direktori (`DO NOT FOLLOW`) tanpa memicu infinite loop.
- **Test Case 3 (Max File Size Guard):**
  - Kondisi: Buat dummy file sebesar 12 MB berekstensi `.tsx`.
  - Ekspektasi: Berkas diabaikan dari antrean `jobs` dan tidak dibaca ke memori.

### 2.2. Uji Concurrency Worker Pool & Pembatalan (`internal/scanner/pool_test.go`)

- **Test Case 1 (Beban Konkuren & Race Detection):**
  - Mengirim 1.000 berkas mock ke worker pool ($N = \text{runtime.GOMAXPROCS(0)}$).
  - Verifikasi: `go test -race ./internal/scanner/...` menghasilkan **0 data race**.
- **Test Case 2 (Context Cancellation Clean Exit):**
  - Batalkan `context.Context` di tengah pemrosesan 500 berkas.
  - Ekspektasi:
    - Seluruh goroutine worker berhenti dengan bersih (*zero goroutine leak*).
    - Channel ditutup rapi dan tidak ada diagnostic parsial yang dikirim ke reporter.

---

## 3. Skenario Pengujian Traversal Analyzer Engine (`internal/analyzer/`)

### 3.1. Uji Direktif Inline Ignore & Span Scope (`internal/analyzer/context_test.go`)

- **Test Case 1 (Grammar Parsing: Whitespace & Deduplication):**
  - Direktif: `// charites:ignore  theme.hardcode-opacity-color  ,  theme.hardcode-opacity-color  `.
  - Ekspektasi: Spasi dipangkas bersih dan ID duplikat disaring.
- **Test Case 2 (Wildcard Suppression):**
  - Direktif: `// charites:ignore *`.
  - Ekspektasi: Seluruh temuan rule pada baris/node tersebut ditekan.
- **Test Case 3 (Node Span Scope: Multi-Line JSX Element):**
  - Kode sumber:
    ```tsx
    // charites:ignore theme.hardcode-opacity-color
    <div
      id="hero"
      className="
        bg-primary/10
      "
    />
    ```
  - Ekspektasi: Karena komentar berada di baris $N$ dan opening tag berada di baris $N+1$, rentang node mencakup baris attribute `bg-primary/10`, sehingga diagnostic berhasil ditekan (0 diagnostic).

### 3.2. Uji Pengurutan Determinis Total Ordering (`internal/analyzer/sort_test.go`)

```go
func TestSortDiagnostics_TotalOrdering(t *testing.T) {
    diags := []ir.Diagnostic{
        {File: "a.tsx", Line: 10, Column: 5, Rule: "theme.b", Severity: ir.SeverityWarn, Message: "msg B", Hint: "hint B"},
        {File: "a.tsx", Line: 10, Column: 5, Rule: "theme.a", Severity: ir.SeverityError, Message: "msg A", Hint: "hint A"},
        {File: "a.tsx", Line: 10, Column: 5, Rule: "theme.a", Severity: ir.SeverityError, Message: "msg A2", Hint: "hint A"},
    }

    ir.SortDiagnostics(diags)

    // Verifikasi total ordering deterministik
    if diags[0].Rule != "theme.a" || diags[0].Message != "msg A" {
        t.Errorf("urutan tidak memenuhi total ordering: got %+v", diags[0])
    }
}
```

---

## 4. Metodologi Benchmark Throughput (`TEST-04-BENCH-001`)

Benchmark performa dipecah menjadi modul terisolasi untuk mengukur hotspot secara akurat:

1. **`BenchmarkWalker_DirectoryTraversal`**: Mengukur throughput pembacaan struktur direktori tanpa parsing.
2. **`BenchmarkIgnore_Matcher`**: Mengukur throughput evaluasi pola glob sekuensial pada 1.000 path.
3. **`BenchmarkAnalyzer_EngineTraversal`**: Mengukur throughput penelusuran AST dan evaluasi active rules pada AST in-memory.
4. **`BenchmarkEndToEnd_ScanPipeline`**: Mengukur throughput end-to-end pada korpus fixture standar.

### Protokol Benchmark:
- Toolchain: Go 1.26 (`go1.26.x`), `CGO_ENABLED=0`.
- Perintah: `go test -bench=. -benchmem ./internal/...`.
