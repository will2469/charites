# 01-SPEC: 06 - Full Pipeline Integration, Golden Snapshots & Fuzzing Specification

> **Kode Dokumen:** `SPEC-06-GOLDEN`
> **Tahapan:** Fase 6 - Validasi Penuh & Golden Snapshots (Milestone Selesai Pipa)
> **Peran Pilar:** SPEC = WHAT (Spesifikasi Validasi Pipa Lengkap, Skema Golden Master & Fuzzing)
> **Status:** Ready for Review
> **Standar Rujukan:** Golden Master Testing Pattern / Go 1.26 Native Fuzzing Specification

Dokumen ini mendefinisikan spesifikasi pengujian integrasi pipa compiler dari ujung ke ujung (*end-to-end pipeline*), standarisasi skema kanonikal **Golden Snapshots** (`tests/golden/*`), struktur korpus berkas percontohan (*fixtures*), protokol ketahanan *native fuzzing*, serta kriteria kelulusan **Gerbang Stabilitas Pipa** (*Pipeline Stability Gate*).

---

## 1. Spesifikasi Golden Master Testing (`tests/golden/`)

Untuk mencegah terjadinya regresi diagnosis atau pergeseran pelaporan (*diagnostic drift*):

### 1.1. Skema Kanonikal Dokumen JSON Golden (ScanResult JSON v1)
Setiap berkas snapshot JSON (`<scenario>.golden.json`) mematuhi skema kanonikal tunggal:

```json
{
  "schema_version": 1,
  "tool_version": "1.0.0",
  "summary": {
    "scanned_files": 28,
    "error_count": 1,
    "warning_count": 1,
    "info_count": 0,
    "passed": false
  },
  "diagnostics": [
    {
      "file": "src/pages/index.astro",
      "line": 14,
      "column": 8,
      "rule": "theme.hardcode-opacity-color",
      "category": "theme",
      "severity": "error",
      "message": "Hardcode opacity color: \"bg-primary/10\"",
      "hint": "Use semantic token \"primary-light\"."
    }
  ]
}
```

### 1.2. Isolasi Metadata Non-Deterministik (`duration_ms`)
- Nilai `duration_ms` merupakan stopwatch runtime yang bervariasi di setiap eksekusi mesin.
- **Invarian Komparasi Golden:** Selama evaluasi uji snapshot golden, field `duration_ms` dinormalisasi (dihilangkan atau diatur ke nilai konstan 0) sehingga pengujian regresi murni memverifikasi semantik diagnosis, bukan fluktuasi CPU.

### 1.3. Format Terminal Text Snapshot (`<scenario>.golden.txt`)
- Menyimpan representasi visual pelaporan teks terminal.
- **Normalisasi Lintas Platform:** Wajib menggunakan encoding UTF-8, pemisah baris LF murni (`\n`, bukan `\r\n`), dan path POSIX forward slash (`/`).
- Disimpan dalam mode `ColorNever` untuk menjamin kesetaraan biner lintas runner CI (Linux, macOS, Windows).

### 1.4. Mekanisme Pembaruan Terkontrol (`-update` Flag)
- Secara bawaan, eksekusi `go test ./tests/...` membandingkan output scanner terkini dengan berkas `.golden.*`. Jika terdapat perbedaan sekecil apa pun, test **MUST FAIL**.
- Pembaruan berkas golden snapshot hanya diizinkan via flag lokal eksplisit:
  ```bash
  go test ./tests/... -run TestPipeline_GoldenSnapshots -update
  ```
- **Larangan CI:** Lingkungan CI **DILARANG KERAS** menjalankan test dengan flag `-update`.
- Ketika `-update` dijalankan, harness wajib mencetak daftar berkas yang diperbarui dan menginstruksikan pengembang untuk meninjau `git diff` sebelum melakukan commit.

---

## 2. Struktur Korpus Fixtures Terstandarisasi (`tests/fixtures/`)

Direktori fixtures memisahkan sampel level berkas dan sampel level proyek:

```
tests/fixtures/
├── astro/                  # Sampel komponen Astro individual
│   ├── clean.astro
│   ├── opacity_violations.astro
│   ├── complex_frontmatter.astro
│   └── inline_ignore.astro
├── tsx/                    # Sampel komponen React TSX individual
│   ├── clean.tsx
│   ├── opacity_violations.tsx
│   ├── template_literals.tsx
│   └── inline_ignore.tsx
└── projects/               # Sampel repositori mini terintegrasi
    ├── clean/              # Repositori bersih (Zero Noise Invariant)
    ├── opacity_violations/ # Repositori dengan beragam pelanggaran
    ├── config_override/    # Repositori dengan charites.yaml overrides
    └── ignore_patterns/    # Repositori dengan .charitesignore & inline ignore
```

---

## 3. Spesifikasi Native Fuzzing Suite (`tests/fuzz/`)

Untuk memastikan parser dan IR builder kebal dari input acak (*chaos resilience*):

1. **Dua Tingkat Pengujian Fuzzing:**
   - **Level 1 (Parser Invariants):** `FuzzAstroParser` dan `FuzzTSXParser` (dari Fase 2) untuk membuktikan ketahanan lexer/parser murni.
   - **Level 2 (Pipeline Integration):** `FuzzAstroPipeline` dan `FuzzTSXPipeline` untuk menguji aliran penuh (Parser $\rightarrow$ IR $\rightarrow$ Engine Traversal $\rightarrow$ Sorting).
2. **Kriteria Kelulusan Fuzzing:**
   Target fuzzing wajib berjalan terus-menerus selama minimal **60 detik per modul** tanpa memicu panic runtime, fatal error stack overflow, atau penghentian abnormal proses Go.

---

## 4. Kriteria Kelulusan Gerbang Stabilitas Pipa (Stability Gate)

Fase 6 adalah titik pembekuan pipa compiler inti:
1. Seluruh pengujian Golden Master lulus 100%.
2. Pengujian Fuzzing lulus $\ge 60\text{ detik}$ tanpa crash.
3. Seluruh alur CLI terbukti terhubung mulus dengan Engine, Rules, Parser, dan IR.
4. **Pembekuan Arsitektur Inti (Architecture Freeze $\neq$ Bug Fix Freeze):**
   - Batas antarmuka modul inti (`internal/ir`, `internal/parser`, `internal/scanner`, `internal/analyzer`, `internal/reporter`) dibekukan dari perombakan struktural.
   - Perbaikan bug (*bug fixes*) dan optimasi non-breaking yang mematuhi kontrak tetap diizinkan.
   - Penambahan rule baru di Fase 8 murni bersifat modular (*pluggable rules*) tanpa memodifikasi engine inti.
