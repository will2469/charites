# 04-QUALITY: 00 - Linter Baseline & Zero-Dependency Hygiene

> **Kode Dokumen:** `QUAL-00-SETUP`
> **Tahapan:** Fase 0 - Inisialisasi Repositori & Toolchain
> **Peran Pilar:** QUALITY = QUALITY THRESHOLD (Ambang Batas Kualitas Linter, Zero Dep & Hygiene)
> **Status:** Graduated (All Phase Gates Passed)

Dokumen ini mendefinisikan ambang batas kualitas (*quality threshold*), konfigurasi linter awal, gerbang analisis statis, dan aturan kebersihan ketergantungan (*dependency hygiene*) pada **Fase 0**.

---

## 1. Konfigurasi Baku Linter (`.golangci.yml`)

Eksekusi Fase 0 **MUST** menyertakan berkas `.golangci.yml` di akar repositori dengan konfigurasi minimal berikut:

```yaml
# yaml-language-server: $schema=https://golangci-lint.run/jsonschema/golangci.v1.jsonschema.json
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
> **Lint Baseline Scope (`QUAL-00-LINT-001`):** *Lint baseline applies to all code present at the current phase.*
> Konfigurasi 9 linter dan `govet enable-all` dievaluasi secara ketat terhadap seluruh kode Go yang eksis pada Fase 0 (`cmd/charites` dan `internal/cli`). Aturan ini tidak menuntut keberadaan kode untuk fase mendatang yang belum dibangun.

---

## 2. Kebijakan Kualitas, Kebersihan Dependensi & Baseline Reproduktifitas

- **QUAL-00-DEPS-001 (Zero External Dependencies):**
  - Berkas `go.mod` hanya memuat deklarasi modul dan versi bahasa Go (`go 1.26`).
  - Dilarang keras menambahkan blok `require` pihak ketiga (*zero third-party dependencies*).
  - Berkas `go.sum` **MUST NOT** be required pada Fase 0 saat dependensi eksternal bernilai nol.

- **QUAL-00-INVAR-001 (Zero CGO Surface Invariant):**
  - *Distinction of Authority:* Bila `SPEC-00-BUILD-001` mewajibkan biner produksi dan artefak rilis dikompilasi murni dengan `CGO_ENABLED=0`, maka `QUAL-00-INVAR-001` bertindak sebagai *defense-in-depth hygiene barrier* yang memverifikasi bahwa struktur kode sumber repositori **murni 100% Go** tanpa berkas `.c`, `.h`, atau blok `import "C"`. Keberadaan CGO hanya diizinkan secara terbatas pada runtime *dynamic verification harness* (`go test -race` via ThreadSanitizer) tanpa mencemari kode sumber produksi.
  - Prosedur verifikasi:
    ```bash
    test -z "$(go list -f '{{range .CgoFiles}}{{.}} {{end}}' ./...)"
    ```

- **QUAL-00-REPRO-001 (Deterministic CI Reproducibility Baseline):**
  - Versi Go pada pipeline CI GitHub Actions **MUST** dipin secara presisi pada `1.26.0` (`go-version: '1.26.0'`) untuk menjamin reproduktifitas deterministik build CI yang tidak terpengaruh oleh drift rilis patch toolchain di runner.

---

## 3. Quality Gate Commands

Sebelum transisi fase disetujui, seluruh perintah verifikasi kualitas berikut **MUST** menghasilkan output bersih tanpa warning/error dan keluar dengan exit code `0`:

```bash
# 1. Analisis statis linter (QUAL-00-LINT-001)
golangci-lint run ./...

# 2. Verifikasi kebersihan surface CGO (QUAL-00-INVAR-001)
test -z "$(go list -f '{{range .CgoFiles}}{{.}} {{end}}' ./...)"

# 3. Verifikasi kompilasi silang tanpa CGO (SPEC-00-BUILD-002 / TEST-00-BUILD-002)
make cross-compile
```

