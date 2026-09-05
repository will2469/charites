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

Pengujian end-to-end membuktikan bahwa binary terkompilasi (`bin/charites`) dapat dieksekusi dari subprocess OS dengan perilaku kontrak yang presisi:

```go
package e2e_test

import (
    "os/exec"
    "strings"
    "testing"
)

func TestBinarySmoke(t *testing.T) {
    binPath := "../../bin/charites"

    tests := []struct {
        name         string
        args         []string
        expectedCode int
        containsOut  string
    }{
        {
            name:         "flag --version",
            args:         []string{"--version"},
            expectedCode: 0,
            containsOut:  "charites version",
        },
        {
            name:         "flag -v",
            args:         []string{"-v"},
            expectedCode: 0,
            containsOut:  "charites version",
        },
        {
            name:         "flag --help",
            args:         []string{"--help"},
            expectedCode: 0,
            containsOut:  "Usage: charites",
        },
        {
            name:         "flag -h",
            args:         []string{"-h"},
            expectedCode: 0,
            containsOut:  "Usage: charites",
        },
        {
            name:         "unknown command exits 2",
            args:         []string{"unknown-command"},
            expectedCode: 2,
            containsOut:  "unknown command",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            cmd := exec.Command(binPath, tc.args...)
            out, err := cmd.CombinedOutput()
            outStr := string(out)

            if tc.expectedCode == 0 && err != nil {
                t.Fatalf("ekspektasi exit 0, didapat error: %v, output: %s", err, outStr)
            }
            if tc.expectedCode != 0 {
                if exitErr, ok := err.(*exec.ExitError); ok {
                    if exitErr.ExitCode() != tc.expectedCode {
                        t.Fatalf("ekspektasi exit %d, didapat %d", tc.expectedCode, exitErr.ExitCode())
                    }
                } else {
                    t.Fatalf("ekspektasi exit code non-zero, didapat err: %v", err)
                }
            }

            if tc.containsOut != "" && !strings.Contains(outStr, tc.containsOut) {
                t.Errorf("output tidak memuat substring yang diharapkan '%s': %s", tc.containsOut, outStr)
            }
        })
    }
}
```

---

## 3. Scaffolding Reservasi Struktur Tri-Corpus Correctness (`tests/correctness/`)

Mengadopsi pola pengujian semantik **Argus**, struktur direktori pengujian rule dipersiapkan sebagai *directory reservation*:
- Direktori `tests/correctness/` dan `tests/golden/` dibentuk sebagai fondasi scaffolding.
- Evaluasi pengujian rule (Positive, Negative, Adversarial) mulai aktif secara formal pada **Fase 3** saat interface rule dan Rule #1 (`theme.hardcode-opacity-color`) diimplementasikan.
- Fase 0 hanya memverifikasi keberadaan reservasi folder tanpa menuntut suite test rule aktif.

---

## 4. Checklist Verifikasi Manual Fase 0

Jalankan rangkaian perintah pembuktian (*proof execution*) di terminal:
```bash
# 1. Pastikan build binary berhasil dengan zero CGO
make build

# 2. Pastikan binary merespons flag versi dan help
./bin/charites --version
./bin/charites -v
./bin/charites --help
./bin/charites -h

# 3. Pastikan argumen tidak dikenal menghasilkan exit code 2
./bin/charites unknown-command || [ $? -eq 2 ]

# 4. Pastikan unit test dan subprocess smoke test lolos
make test
```
Seluruh perintah di atas **wajib** membuktikan pemenuhan kontrak SPEC-00.

