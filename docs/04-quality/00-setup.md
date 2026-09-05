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

## 2. Kebijakan Zero-Dependency di `go.mod` & Status `go.sum`

Pada Fase 0:
- Berkas `go.mod` hanya memuat deklarasi modul dan versi Go:
  ```text
  module github.com/will2469/charites

  go 1.26
  ```
- **Go Toolchain:** Dikelola via standar `go1.26.x`.
- **Kebijakan Dependensi Eksternal:** Dilarang keras menambahkan blok `require` pihak ketiga (*zero third-party dependencies*).
- **Status `go.sum`:** Berkas `go.sum` **MUST NOT** be required pada Fase 0 saat dependensi eksternal bernilai nol.

---

## 3. Quality Gate Command

Sebelum transisi fase disetujui:
```bash
golangci-lint run ./...
```
Perintah ini **MUST** menghasilkan output bersih tanpa warning/error dan keluar dengan exit code `0`.

