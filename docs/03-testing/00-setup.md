# 03-TESTING: 00 - Toolchain & Scaffolding Smoke Verification

> **Kode Dokumen:** `TEST-00-SETUP`
> **Tahapan:** Fase 0 - Inisialisasi Repositori & Toolchain
> **Status:** Ready for Execution

Dokumen ini mendefinisikan skenario uji coba awal (*smoke tests*) untuk memverifikasi bahwa scaffolding repositori, kompilasi binary, dan CLI trampoline pada **Fase 0** berfungsi 100% tanpa cacat.

---

## 1. Skenario Uji Coba Unit (`internal/cli/root_test.go`)

Paket `internal/cli` wajib memiliki unit test deterministik:

### Test Case 1: Flag Version (`--version` & `-v`)
- **Input Argumen:** `[]string{"--version"}` dan `[]string{"-v"}`
- **Ekspektasi:**
  - Fungsi `cli.Execute()` mengembalikan exit code `0`.
  - Output mencakup string `charites version`.

### Test Case 2: Flag Help (`--help` & `-h`)
- **Input Argumen:** `[]string{"--help"}` dan `[]string{"-h"}`
- **Ekspektasi:**
  - Fungsi `cli.Execute()` mengembalikan exit code `0`.
  - Output mencakup panduan penggunaan `Usage: charites`.

### Test Case 3: Command Tidak Dikenal (*Unknown Command*)
- **Input Argumen:** `[]string{"unknown-command"}`
- **Ekspektasi:**
  - Fungsi `cli.Execute()` mengembalikan exit code `2` (CLI error).
  - Pesan error dicetak ke `stderr`.

---

## 2. Skenario Subprocess End-to-End Smoke Test (`tests/e2e/smoke_test.go`)

Pengujian memastikan binary terkompilasi dapat dipanggil dari luar proses:

```go
package e2e_test

import (
    "os/exec"
    "strings"
    "testing"
)

func TestBinarySmoke(t *testing.T) {
    // 1. Jalankan binary hasil build
    cmd := exec.Command("../../bin/charites", "--version")
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("gagal mengeksekusi binary: %v, output: %s", err, string(out))
    }

    // 2. Verifikasi konten
    if !strings.Contains(string(out), "charites version") {
        t.Errorf("output tidak memuat string versi yang valid: %s", string(out))
    }
}
```

---

## 3. Scaffolding Struktur Tri-Corpus Correctness (`tests/correctness/`)

Mengadopsi pola pengujian semantik **Argus**, struktur direktori pengujian rule dipersiapkan sejak Fase 0 dengan standar:
- Setiap rule memiliki folder sendiri di `tests/correctness/<rule_id>/`.
- Setiap folder rule wajib memiliki 3 sub-direktori:
  1. `positive/` - Sampel kode yang dengan sengaja melanggar rule.
  2. `negative/` - Sampel kode valid yang wajib menghasilkan 0 deteksi (*Zero-Noise Invariant*).
  3. `adversarial/` - Kasus ekstrem, edge syntax, dynamic interpolation, dan jebakan *false positive*.

Struktur snapshot golden regression disiapkan di `tests/golden/` untuk menampung expected JSON & ANSI output.

---

## 4. Checklist Verifikasi Manual Fase 0

Jalankan perintah berikut di terminal:
```bash
# 1. Pastikan build berhasil
make build

# 2. Pastikan binary merespons
./bin/charites --version
./bin/charites --help

# 3. Pastikan test suite lolos
make test
```
Seluruh perintah di atas **wajib** keluar dengan exit code `0`.

