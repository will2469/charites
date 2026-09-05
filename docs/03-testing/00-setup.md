# 03-TESTING: 00 - Toolchain & Scaffolding Smoke Verification

> **Kode Dokumen:** `TEST-00-SETUP`
> **Tahapan:** Fase 0 - Inisialisasi Repositori & Toolchain
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Skenario Smoke & Asersi Pembuktian)
> **Status:** Ready for Execution

Dokumen ini mendefinisikan skenario uji coba awal (*smoke tests*) dan pembuktian empiris untuk memverifikasi bahwa seluruh kontrak fungsional pada **Fase 0** (SPEC-00) dan rancangan arsitektur (ARCH-00) terbukti berfungsi 100% tanpa cacat.

---

## 1. Skenario Uji Coba Unit CLI (`internal/cli/root_test.go`)

Paket `internal/cli` wajib memiliki unit test deterministik yang menguji seluruh kontrak CLI Fase 0:

### Test Case 1: Flag Version (`--version`, `-v`) & Subcommand `version`
- **Input Argumen:** `[]string{"--version"}`, `[]string{"-v"}`, dan `[]string{"version"}`
- **Ekspektasi:**
  - Fungsi `cli.Execute()` mengembalikan exit code `0`.
  - Output `stdout` mencakup string `charites version 0.1.0-dev`.

### Test Case 2: Flag Help (`--help`, `-h`)
- **Input Argumen:** `[]string{"--help"}` dan `[]string{"-h"}`
- **Ekspektasi:**
  - Fungsi `cli.Execute()` mengembalikan exit code `0`.
  - Output `stdout` mencakup panduan penggunaan `Usage: charites`.

### Test Case 3: Pemanggilan Tanpa Argumen (`[]string{}`)
- **Input Argumen:** `[]string{}` (slice kosong)
- **Ekspektasi:**
  - Fungsi `cli.Execute()` mengembalikan exit code `0`.
  - Output `stdout` mencakup panduan penggunaan dasar `Usage: charites`.

### Test Case 4: Subcommand / Flag Tidak Dikenal (*Unknown Command / Invalid Flag*)
- **Input Argumen:** `[]string{"unknown-command"}` dan `[]string{"--bogus-flag"}`
- **Ekspektasi:**
  - Fungsi `cli.Execute()` mengembalikan exit code `2` (CLI error).
  - Pesan error deskriptif dicetak ke `stderr`.

---

## 2. Skenario Subprocess End-to-End Smoke Test (`tests/e2e/smoke_test.go`)

Pengujian end-to-end membuktikan bahwa binary terkompilasi (`bin/charites`) dapat dieksekusi dari subprocess OS dengan perilaku kontrak yang presisi, termasuk **pemisahan saluran stream routing** (`stdout` vs `stderr`):

```go
package e2e_test

import (
    "bytes"
    "os/exec"
    "strings"
    "testing"
)

func TestBinarySmoke(t *testing.T) {
    binPath := "../../bin/charites"

    tests := []struct {
        name           string
        args           []string
        expectedCode   int
        containsStdout string
        containsStderr string
    }{
        {
            name:           "flag --version",
            args:           []string{"--version"},
            expectedCode:   0,
            containsStdout: "charites version",
        },
        {
            name:           "flag -v",
            args:           []string{"-v"},
            expectedCode:   0,
            containsStdout: "charites version",
        },
        {
            name:           "subcommand version",
            args:           []string{"version"},
            expectedCode:   0,
            containsStdout: "charites version",
        },
        {
            name:           "flag --help",
            args:           []string{"--help"},
            expectedCode:   0,
            containsStdout: "Usage: charites",
        },
        {
            name:           "flag -h",
            args:           []string{"-h"},
            expectedCode:   0,
            containsStdout: "Usage: charites",
        },
        {
            name:           "empty args usage",
            args:           []string{},
            expectedCode:   0,
            containsStdout: "Usage: charites",
        },
        {
            name:           "unknown command exits 2 to stderr",
            args:           []string{"unknown-command"},
            expectedCode:   2,
            containsStderr: "unknown command",
        },
        {
            name:           "unknown flag exits 2 to stderr",
            args:           []string{"--bogus-flag"},
            expectedCode:   2,
            containsStderr: "unknown",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            cmd := exec.Command(binPath, tc.args...)
            var stdout, stderr bytes.Buffer
            cmd.Stdout = &stdout
            cmd.Stderr = &stderr

            err := cmd.Run()
            stdoutStr := stdout.String()
            stderrStr := stderr.String()

            if tc.expectedCode == 0 {
                if err != nil {
                    t.Fatalf("ekspektasi exit 0, didapat error: %v, stderr: %s", err, stderrStr)
                }
                if stderrStr != "" {
                    t.Errorf("ekspektasi stderr bersih pada exit 0, didapat: %s", stderrStr)
                }
            } else {
                if exitErr, ok := err.(*exec.ExitError); ok {
                    if exitErr.ExitCode() != tc.expectedCode {
                        t.Fatalf("ekspektasi exit %d, didapat %d (stderr: %s)", tc.expectedCode, exitErr.ExitCode(), stderrStr)
                    }
                } else {
                    t.Fatalf("ekspektasi exit code %d, didapat err: %v", tc.expectedCode, err)
                }
            }

            if tc.containsStdout != "" && !strings.Contains(stdoutStr, tc.containsStdout) {
                t.Errorf("stdout tidak memuat substring yang diharapkan '%s': %s", tc.containsStdout, stdoutStr)
            }
            if tc.containsStderr != "" && !strings.Contains(stderrStr, tc.containsStderr) {
                t.Errorf("stderr tidak memuat substring yang diharapkan '%s': %s", tc.containsStderr, stderrStr)
            }
        })
    }
}
```

---

## 3. Prosedur Verifikasi Kompilasi Silang (`TEST-00-BUILD-002`)

Untuk membuktikan pemenuhan kontrak `SPEC-00-BUILD-002` (Cross-Platform Targets), pengujian otomatis memverifikasi bahwa kompilasi silang berhasil untuk 4 platform resmi tanpa ketergantungan CGO:

```bash
# Target 1: Linux x86_64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /dev/null ./cmd/charites

# Target 2: Linux ARM64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o /dev/null ./cmd/charites

# Target 3: macOS Apple Silicon
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o /dev/null ./cmd/charites

# Target 4: Windows x86_64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o /dev/null ./cmd/charites
```

Setiap perintah di atas **MUST** menghasilkan exit code `0` tanpa peringatan linker.

---

## 4. Scaffolding Reservasi Struktur Tri-Corpus Correctness (`tests/correctness/`)

Mengadopsi pola pengujian semantik **Argus**, struktur direktori pengujian rule dipersiapkan sebagai *directory reservation*:
- Direktori `tests/correctness/` dan `tests/golden/` dibentuk sebagai fondasi scaffolding.
- Evaluasi pengujian rule (Positive, Negative, Adversarial) mulai aktif secara formal pada **Fase 3** saat interface rule dan Rule #1 (`theme.hardcode-opacity-color`) diimplementasikan.
- Fase 0 hanya memverifikasi keberadaan reservasi folder tanpa menuntut suite test rule aktif.

---

## 5. Checklist Verifikasi Manual Fase 0

Jalankan rangkaian perintah pembuktian (*proof execution*) di terminal:
```bash
# 1. Pastikan rantai Makefile all berjalan mulus dari awal (build -> test -> lint)
make all

# 2. Pastikan binary merespons flag versi dan help dengan stream routing presisi
./bin/charites --version
./bin/charites -v
./bin/charites --help
./bin/charites -h

# 3. Pastikan argumen tidak dikenal menghasilkan error murni ke stderr dengan exit code 2
./bin/charites unknown-command 2> err.log || [ $? -eq 2 ]
grep -q "unknown command" err.log && rm err.log

# 4. Pastikan verifikasi 4 target kompilasi silang lolos
make cross-compile
```
Seluruh perintah di atas **wajib** membuktikan pemenuhan kontrak SPEC-00.

