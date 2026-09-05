# 03-TESTING: 05 - CLI Dispatcher, End-to-End & Reporter Verification Plan

> **Kode Dokumen:** `TEST-05-CLI`
> **Tahapan:** Fase 5 - Reporter Output & CLI Entrypoint
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Matriks E2E & Snapshot Reporter)
> **Status:** Ready for Review
> **Standar Rujukan:** CLI Subprocess Testing & Golden Output Regression Testing

Dokumen ini mendefinisikan strategi pengujian untuk antarmuka baris perintah (`internal/cli/*`), format dokumen presenter laporan (`internal/reporter/*`), pengujian terintegrasi *subprocess* E2E (`tests/e2e/*`), dan verifikasi kontrak *exit code*.

---

## 1. Skenario Pengujian Unit Dispatcher & Flag (`internal/cli/`)

### 1.1. Uji Routing Subcommand & Default Invocation (`root_test.go`)
- **Test Case 1 (Empty Invocation Mapping):**
  - Input: `[]string{}`.
  - Ekspektasi: Ter-routing otomatis ke `RunScan([]string{"."})`.
- **Test Case 2 (Subcommand Aliases):**
  - Input: `[]string{"check", "."}` dan `[]string{"run", "."}`.
  - Ekspektasi: Menghasilkan alur eksekusi identik dengan `[]string{"scan", "."}`.
- **Test Case 3 (Unknown Command):**
  - Input: `[]string{"foobar"}`.
  - Ekspektasi: Menghasilkan exit code `2` dan pesan kesalahan di `stderr`.

### 1.2. Uji Normalisasi & Validasi Flag (`scan_test.go`)
- **Test Case 1 (Extension Normalization & Invalid Ext):**
  - Input valid: `--ext=ASTRO,.tsx`. Ekspektasi: `[".astro", ".tsx"]`.
  - Input tidak valid: `--ext=vue`. Ekspektasi: Exit code `2` dengan error `unsupported extension`.
  - Input kosong: `--ext=`. Ekspektasi: Exit code `2` dengan error `empty extension flag`.
- **Test Case 2 (Category & Rule Conflict):**
  - Input konflik: `--category=theme --rule=a11y.alt-text`.
  - Ekspektasi: Exit code `2` dengan error `rule does not belong to category`.

---

## 2. Skenario Pengujian Unit Reporter (`internal/reporter/`)

### 2.1. Uji Determinis Biner Reporter (`reporter_test.go`)
```go
func TestReporter_Determinism(t *testing.T) {
    result := createSampleScanResult()

    // 1. Uji Determinisme JSON
    var bufJSON1, bufJSON2 bytes.Buffer
    jsonRep := reporter.NewJSONReporter()
    _ = jsonRep.Render(&bufJSON1, result)
    _ = jsonRep.Render(&bufJSON2, result)

    if !bytes.Equal(bufJSON1.Bytes(), bufJSON2.Bytes()) {
        t.Fatalf("JSON reporter output is not byte-for-byte identical")
    }

    // 2. Uji Determinisme Inline ANSI
    var bufInline1, bufInline2 bytes.Buffer
    inlineRep := reporter.NewInlineReporter(reporter.ColorNever)
    _ = inlineRep.Render(&bufInline1, result)
    _ = inlineRep.Render(&bufInline2, result)

    if !bytes.Equal(bufInline1.Bytes(), bufInline2.Bytes()) {
        t.Fatalf("Inline reporter output is not byte-for-byte identical")
    }
}
```

### 2.2. Golden Contract Snapshots (`tests/golden/reporters/`)
Fase 5 mengunci kontrak visual dan format mesin melalui berkas snapshot referensi:
1. `inline_clean.golden`: Format teks bersih saat 0 pelanggaran.
2. `inline_violations.golden`: Format teks berwarna ANSI saat ditemukan error & warning.
3. `inline_no_color.golden`: Format teks polos saat mode `ColorNever` aktif.
4. `json_clean.golden`: Format dokumen JSON saat 0 pelanggaran.
5. `json_violations.golden`: Format dokumen JSON lengkap dengan temuan diagnostic.

---

## 3. Pengujian Terintegrasi E2E Subprocess (`tests/e2e/cli_test.go`)

Pengujian E2E mengompilasi binary `charites` dan memvalidasi seluruh matriks interaksi terminal:

```go
package e2e_test

import (
    "os/exec"
    "testing"
)

func TestCLI_Matrix(t *testing.T) {
    bin := buildBinary(t)

    matrix := []struct {
        name     string
        args     []string
        wantExit int
    }{
        {"Empty_Invocation_Defaults_To_Scan", []string{}, 0},
        {"Subcommand_Check_Alias", []string{"check", "tests/fixtures/clean"}, 0},
        {"Subcommand_Run_Alias", []string{"run", "tests/fixtures/clean"}, 0},
        {"Version_Flag", []string{"-v"}, 0},
        {"Clean_Repo_Exit_0", []string{"scan", "tests/fixtures/clean"}, 0},
        {"Warning_Only_Defaults_To_0", []string{"scan", "tests/fixtures/warning_only"}, 0},
        {"Warning_Only_With_FailOnWarn_Exits_1", []string{"scan", "tests/fixtures/warning_only", "--fail-on-warn"}, 1},
        {"Violations_Found_Exits_1", []string{"scan", "tests/fixtures/violations"}, 1},
        {"Unsupported_Ext_Exits_2", []string{"scan", ".", "--ext=vue"}, 2},
        {"Category_Rule_Conflict_Exits_2", []string{"scan", ".", "--category=theme", "--rule=a11y.alt"}, 2},
        {"Unknown_Flag_Exits_2", []string{"scan", ".", "--unknown-flag"}, 2},
        {"Non_Existent_Target_Exits_2", []string{"scan", "non/existent/path"}, 2},
    }

    for _, tt := range matrix {
        t.Run(tt.name, func(t *testing.T) {
            cmd := exec.Command(bin, tt.args...)
            err := cmd.Run()
            gotExit := 0
            if exitErr, ok := err.(*exec.ExitError); ok {
                gotExit = exitErr.ExitCode()
            }
            if gotExit != tt.wantExit {
                t.Fatalf("args %v: want exit code %d, got %d", tt.args, tt.wantExit, gotExit)
            }
        })
    }
}
```

---

## 4. Metodologi Benchmark Formatting Reporter (`TEST-05-BENCH-001`)

Pengujian kinerja pelaporan diisolasi dari I/O terminal nyata:

```go
func BenchmarkReporter_JSON_Format(b *testing.B) {
    result := generateMockResult(1000) // 1.000 findings
    rep := reporter.NewJSONReporter()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = rep.Render(io.Discard, result)
    }
}
```
Metodologi ini menjamin pengukuran murni terhadap efisiensi serialisasi data Go tanpa gangguan fluktuasi emulator terminal atau sistem I/O disk.
