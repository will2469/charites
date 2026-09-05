# 04-QUALITY: 00 - Linter Baseline & Zero-Dependency Hygiene

> **Kode Dokumen:** `QUAL-00-SETUP`
> **Tahapan:** Fase 0 - Inisialisasi Repositori & Toolchain
> **Peran Pilar:** QUALITY = QUALITY THRESHOLD (Ambang Batas Kualitas Linter, Zero Dep & Hygiene)
> **Status:** Ready for Execution

Dokumen ini mendefinisikan ambang batas kualitas (*quality threshold*), konfigurasi linter awal, gerbang analisis statis, dan aturan kebersihan ketergantungan (*dependency hygiene*) pada **Fase 0**.

---

## 1. Konfigurasi Baku Linter (`.golangci.yml`)

Eksekusi Fase 0 **MUST** menyertakan berkas `.golangci.yml` di akar repositori dengan konfigurasi minimal berikut:

```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - revive
    - gocritic
    - gosec
    - prealloc
    - unconvert
    - unused

linters-settings:
  govet:
    enable-all: true
    disable:
      - fieldalignment
  revive:
    rules:
      - name: exported
        disabled: false
      - name: package-comments
        disabled: true
  gosec:
    excludes:
      - G104 # errcheck sudah menangani unhandled error

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

> [!IMPORTANT]
> **Lint Baseline Scope:** *Lint baseline applies to all code present at the current phase.*
> Konfigurasi 9 linter dan `govet enable-all` dievaluasi secara ketat terhadap seluruh kode Go yang eksis pada Fase 0 (`cmd/charites` dan `internal/cli`). Aturan ini tidak menuntut keberadaan kode untuk fase mendatang yang belum dibangun.

---

## 2. Kebijakan Zero-Dependency di `go.mod`, Status `go.sum` & Versi Pinned

Pada Fase 0:
- **Go Version & Compatibility:**
  - **Language / Module Compatibility:** `go 1.26` di `go.mod`.
  - **Supported Local Toolchain:** `go1.26.x`.
  - **CI / Reproducibility Baseline:** `go1.26.0` (dipin secara deterministik pada GitHub Actions runner).
- **Kebijakan Dependensi Eksternal:** Dilarang keras menambahkan blok `require` pihak ketiga (*zero third-party dependencies*).
- **Status `go.sum`:** Berkas `go.sum` **MUST NOT** be required pada Fase 0 saat dependensi eksternal bernilai nol.
- **Invarian Zero CGO:** Seluruh paket dilarang menyertakan file CGO (`.c`, `.h`, atau blok `import "C"`). Verifikasi:
  ```bash
  test -z "$(go list -f '{{range .CgoFiles}}{{.}} {{end}}' ./...)"
  ```

---

## 3. Quality Gate Commands

Sebelum transisi fase disetujui, seluruh perintah verifikasi kualitas berikut **MUST** menghasilkan output bersih tanpa warning/error dan keluar dengan exit code `0`:

```bash
# 1. Analisis statis linter (9 linter aktif)
golangci-lint run ./...

# 2. Verifikasi kompilasi silang tanpa CGO (SPEC-00-BUILD-002)
make cross-compile
```

