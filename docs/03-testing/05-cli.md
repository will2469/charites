# 03-TESTING: 05 - CLI Dispatcher, End-to-End & Reporter Verification Plan

> **Kode Dokumen:** `TEST-05-CLI`
> **Tahapan:** Fase 5 - Reporter Output & CLI Entrypoint
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Matriks E2E & Snapshot Reporter)
> **Status:** Ready / Approved for Implementation
> **Standar Rujukan:** CLI Subprocess Testing & Golden Output Regression Testing

Dokumen ini mendefinisikan strategi pengujian untuk antarmuka baris perintah (`internal/cli/*`), format dokumen presenter laporan (`internal/reporter/*`), pengujian terintegrasi *subprocess* E2E (`tests/e2e/*`), dan verifikasi kontrak *exit code*.

---

## 1. Matriks Pengujian Eksekutabel 13 Skenario Kritis

Fase 5 mewajibkan pembuktian empiris melalui matriks 13 skenario uji eksekutabel:

| No | Kasus Uji | Perintah CLI | Kondisi Masukan | Ekspektasi Hasil & Exit Code |
| :---: | :--- | :--- | :--- | :--- |
| **1** | `charites` $\rightarrow$ `scan .` | `charites` | 0 argumen pada root repositori | Menjalankan `scan .`, memindai berkas default, exit `0` (clean) atau `1` (violations). |
| **2** | Path/Flag tanpa subcommand | `charites src/` atau `charites --format=json` | Argumen diawali path atau flag `-` | Ditransformasikan otomatis ke `scan src/` atau `scan --format=json .`. |
| **3** | Ekuivalensi `scan`, `check`, `run` | `charites scan .`<br>`charites check .`<br>`charites run .` | Target direktori identik | Output stdout, stderr, dan exit code identik 100% secara biner. |
| **4** | Valid vs Invalid `--ext` | `charites scan --ext=astro,tsx .`<br>`charites scan --ext=vue .` | Ekstensi terdaftar vs ekstensi tidak didukung | Valid: memindai hanya `.astro` & `.tsx`.<br>Invalid: stderr error `unsupported extension "vue"`, exit `2`. |
| **5** | Repeated + Comma-separated `--ext` | `charites scan --ext=astro,tsx --ext jsx .` | Kombinasi koma dan flag berulang | Dinormalisasi menjadi himpunan unik: `[".astro", ".tsx", ".jsx"]`. |
| **6** | `--category` $\times$ `--rule` Intersection | `charites scan --category=theme --rule=theme.hardcode-opacity-color .`<br>`charites scan --category=theme --rule=a11y.alt .` | Irisan valid vs konflik antar-kategori | Cocok: dieksekusi normal.<br>Konflik: stderr error `rule does not belong to category`, exit `2`. |
| **7** | Config `off` vs Explicit `--rule` | `charites scan --rule=theme.color .` dengan `charites.yaml` (`theme.color: "off"`) | Rule dimatikan via policy config | Invarian 3-Tier Precedence: Policy `off` mengalahkan CLI flag, 0 rule dievaluasi, exit `0`. |
| **8** | Builtin Hard Exclusion + `--ignore` | `charites scan --ignore="!node_modules" .` | Mencoba meniadakan hard exclusion | Builtin exclusion kebal negasi `!`, `node_modules` tetap di-skip, 0 pelanggaran dari `node_modules`. |
| **9** | Inline ANSI Reporter | `charites scan --format=inline .` | Repo dengan temuan pelanggaran | Format teks ANSI, path forward slash POSIX, problem summary, durasi waktu `ms`. |
| **10** | JSON Reporter & Schema | `charites scan --format=json .` | Format dokumen JSON tunggal | Valid terhadap skema JSON: `version`, `summary` (`scanned_files`, `duration_ms`, `passed`), `diagnostics`. |
| **11** | TTY / Pipe / `NO_COLOR` | `NO_COLOR=1 charites scan .`<br>`charites scan . \| cat`<br>`charites scan --no-color .` | Lingkungan non-TTY atau flag/env no-color | ANSI escape codes dinonaktifkan (Mode `ColorNever`), output teks murni tanpa ANSI bytes. |
| **12** | Deterministic Ordering | 2x eksekusi `charites scan --format=json .` pada korpus sama | Urutan scheduling paralel acak | Output biner identik 100% (`bytes.Equal`), urutan slice mengikuti 7-tier total ordering. |
| **13** | Exit Codes Taxonomy | Berbagai kondisi pemindaian | Skenario clean, violations, operational error, cancellation | `0`: Clean repo (atau warning tanpa fail-on-warn).<br>`1`: Violation error (atau warning + fail-on-warn).<br>`2`: Operational / CLI error.<br>`130`: SIGINT cancel. |

---

## 2. Spesifikasi Bukti Kontrak Golden & Snapshot (`tests/golden/reporters/`)

Fase 5 mengunci kontrak visual dan format mesin melalui berkas snapshot referensi biner:
1. **`inline_clean.golden`**: Snapshot teks untuk pemindaian repositori bersih (0 error, 0 warning).
2. **`inline_violations.golden`**: Snapshot teks berwarna ANSI saat ditemukan error & warning.
3. **`inline_no_color.golden`**: Snapshot teks polos saat mode `ColorNever` aktif (`--no-color` / `NO_COLOR`).
4. **`json_clean.golden`**: Snapshot dokumen JSON terformat lengkap saat repositori bersih (`passed: true`).
5. **`json_violations.golden`**: Snapshot dokumen JSON lengkap dengan array temuan diagnostic (`passed: false`).

### 2.1. Uji Asersi Determinisme Biner Reporter (`reporter_test.go`)

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

---

## 3. Pengujian Terintegrasi E2E Subprocess (`tests/e2e/cli_test.go`)

Pengujian E2E mengompilasi binary `charites` ke direktori sementara dan memvalidasi seluruh matriks interaksi terminal:

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
        {"Direct_Path_Without_Subcommand", []string{"tests/fixtures/clean"}, 0},
        {"Direct_Flag_Without_Subcommand", []string{"--format=json", "tests/fixtures/clean"}, 0},
        {"Subcommand_Check_Alias", []string{"check", "tests/fixtures/clean"}, 0},
        {"Subcommand_Run_Alias", []string{"run", "tests/fixtures/clean"}, 0},
        {"Version_Flag", []string{"-v"}, 0},
        {"Clean_Repo_Exit_0", []string{"scan", "tests/fixtures/clean"}, 0},
        {"Warning_Only_Defaults_To_0", []string{"scan", "tests/fixtures/warning_only"}, 0},
        {"Warning_Only_With_FailOnWarn_Exits_1", []string{"scan", "tests/fixtures/warning_only", "--fail-on-warn"}, 1},
        {"Violations_Found_Exits_1", []string{"scan", "tests/fixtures/violations"}, 1},
        {"Unsupported_Ext_Exits_2", []string{"scan", ".", "--ext=vue"}, 2},
        {"Empty_Ext_Exits_2", []string{"scan", ".", "--ext="}, 2},
        {"Category_Rule_Conflict_Exits_2", []string{"scan", ".", "--category=theme", "--rule=a11y.alt"}, 2},
        {"Unsupported_Format_Exits_2", []string{"scan", ".", "--format=xml"}, 2},
        {"Unknown_Command_Exits_2", []string{"foobar"}, 2},
        {"Unknown_Flag_Exits_2", []string{"scan", ".", "--unknown-flag"}, 2},
        {"Non_Existent_Target_Exits_2", []string{"scan", "non/existent/path"}, 2},
        {"Builtin_Ancestor_Target_Exits_2", []string{"scan", "node_modules/foo/Card.tsx"}, 2},
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
    for b.Loop() {
        _ = rep.Render(io.Discard, result)
    }
}
```
Pengujian ini menggunakan idiom benchmark modern Go 1.26 (`b.Loop()`) dan menguji throughput formatting terisolasi murni.
