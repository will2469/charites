# 04-QUALITY: 00 - Linter Baseline & Zero-Dependency Hygiene

> **Kode Dokumen:** `QUAL-00-SETUP`
> **Tahapan:** Fase 0 - Inisialisasi Repositori & Toolchain
> **Status:** Ready for Execution

Dokumen ini mendefinisikan konfigurasi linter awal, gerbang analisis statis, dan aturan kebersihan ketergantungan (*dependency hygiene*) pada **Fase 0**.

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

---

## 2. Kebijakan Zero-Dependency di `go.mod`

Pada Fase 0:
- File `go.mod` hanya memuat deklarasi modul dan versi Go:
  ```text
  module github.com/will2469/charites

  go 1.26.0
  ```
- **Larangan Keras:** Tidak boleh ada blok `require` pihak ketiga yang belum disetujui. File `go.sum` harus kosong atau belum ada jika tidak ada dependensi eksternal.

---

## 3. Quality Gate Command

Sebelum commit pertama di-push:
```bash
golangci-lint run ./...
```
Perintah ini **MUST** menghasilkan output kosong dan exit code `0`.
