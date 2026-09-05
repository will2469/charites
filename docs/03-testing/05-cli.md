# 03-TESTING: 05 - CLI Dispatcher, End-to-End & Reporter Verification Plan

> **Kode Dokumen:** `TEST-05-CLI`
> **Tahapan:** Fase 5 - Reporter Output & CLI Entrypoint
> **Status:** Ready for Review
> **Standar Rujukan:** CLI Subprocess Testing & Golden Output Regression Testing

Dokumen ini mendefinisikan strategi pengujian untuk antarmuka baris perintah (`internal/cli/*`), format keluaran presenter (`internal/reporter/*`), pengujian terintegrasi *subprocess* E2E (`tests/e2e/*`), dan verifikasi kontrak *exit code*.

---

## 1. Skenario Pengujian Unit Dispatcher & CLI Flags (`internal/cli/`)

### 1.1. Uji Routing Subcommand & Alias (`internal/cli/root_test.go`)
- **Test Case 1 (Default Command Mapping):**
  - Input argumen: `[]string{"."}` atau `[]string{"src/Button.tsx"}`.
  - Ekspektasi: Ter-routing otomatis ke handler pemindaian `runScan()`.
- **Test Case 2 (Subcommand Aliases - Ergonomi E):**
  - Input argumen: `[]string{"check", "."}` dan `[]string{"run", "."}`.
  - Ekspektasi: Menghasilkan konfigurasi `ScanOptions` dan alur eksekusi yang identik dengan `[]string{"scan", "."}`.
- **Test Case 3 (Subcommand Version):**
  - Input argumen: `[]string{"version"}` atau `[]string{"-v"}`.
  - Ekspektasi: Menampilkan info versi dan mengembalikan exit code 0.

### 1.2. Uji Parsing Flag Ergonomi (A-D) (`internal/cli/scan_test.go`)
- **Test Case 1 (Direct File Target - Ergonomi A):**
  - Input: `charites scan src/pages/index.astro`.
  - Ekspektasi: `opts.TargetPath == "src/pages/index.astro"` dan `opts.IsDirectFile == true`.
- **Test Case 2 (Extension Filter - Ergonomi B):**
  - Input: `charites scan . --ext=astro,tsx`.
  - Ekspektasi: `opts.Extensions` memuat `[".astro", ".tsx"]`.
- **Test Case 3 (Category & Rule Filter - Ergonomi C & D):**
  - Input: `charites scan . --category=theme --rule=theme.hardcode-opacity-color`.
  - Ekspektasi: `opts.CategoryFilter == "theme"` dan `opts.RuleFilter == "theme.hardcode-opacity-color"`.

---

## 2. Skenario Pengujian Unit Reporter (`internal/reporter/`)

### 2.1. Uji Coba Inline ANSI Reporter (`internal/reporter/inline_test.go`)
- **Test Case 1 (ANSI Color Codes Rendering):**
  - Input: `ScanResult` memuat 1 error dan 1 warning.
  - Ekspektasi: Buffer keluaran memuat ANSI escape code merah (`\x1b[31;1m`) untuk error dan kuning (`\x1b[33;1m`) untuk warning.
- **Test Case 2 (No-Color Invariant):**
  - Kondisi: Mengaktifkan flag `--no-color` atau mengeset variabel lingkungan `NO_COLOR=1`.
  - Ekspektasi: Buffer keluaran bersih dari karakter escape ANSI (`\x1b[`). Teks berupa string polos yang mudah dibaca.
- **Test Case 3 (Footer Summary):**
  - Ekspektasi: Mencetak baris ringkasan yang memuat jumlah berkas yang dipindai, waktu eksekusi dalam milidetik, dan total temuan.

### 2.2. Uji Coba JSON Stream Reporter (`internal/reporter/json_test.go`)
- **Test Case 1 (Schema & Deserialization Invariant):**
  - Input: `ScanResult` dengan data diagnostic lengkap.
  - Ekspektasi: Output JSON berhasil di-unmarshal kembali ke struct Go tanpa error sintaks.
- **Test Case 2 (Field Preservations):**
  - Ekspektasi: Seluruh field (`file`, `line`, `column`, `rule`, `category`, `severity`, `message`, `hint`) tidak terpotong dan memiliki tipe data yang presisi.

---

## 3. Pengujian Terintegrasi E2E Subprocess (`tests/e2e/cli_test.go`)

Pengujian E2E mengompilasi binary `charites` ke direktori sementara dan mengeksekusinya via `os/exec`:

```go
package e2e_test

import (
    "os/exec"
    "testing"
)

func TestCLI_SubprocessExecution(t *testing.T) {
    binPath := buildTestBinary(t)

    // Test 1: Clean Directory -> Exit Code 0
    t.Run("Clean_Directory_Exit_0", func(t *testing.T) {
        cmd := exec.Command(binPath, "scan", "tests/fixtures/clean")
        err := cmd.Run()
        if err != nil {
            t.Fatalf("expected exit code 0, got %v", err)
        }
    })

    // Test 2: Violations Found -> Exit Code 1
    t.Run("Violations_Found_Exit_1", func(t *testing.T) {
        cmd := exec.Command(binPath, "scan", "tests/fixtures/violations")
        err := cmd.Run()
        if exitErr, ok := err.(*exec.ExitError); ok {
            if exitErr.ExitCode() != 1 {
                t.Fatalf("expected exit code 1, got %d", exitErr.ExitCode())
            }
        } else {
            t.Fatalf("expected exit error, got nil")
        }
    })

    // Test 3: Non-Existent Path -> Exit Code 2
    t.Run("Fatal_Error_Exit_2", func(t *testing.T) {
        cmd := exec.Command(binPath, "scan", "non/existent/path")
        err := cmd.Run()
        if exitErr, ok := err.(*exec.ExitError); ok {
            if exitErr.ExitCode() != 2 {
                t.Fatalf("expected exit code 2, got %d", exitErr.ExitCode())
            }
        } else {
            t.Fatalf("expected exit error, got nil")
        }
    })
}
```
