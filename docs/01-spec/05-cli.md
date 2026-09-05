# 01-SPEC: 05 - CLI Interface, Commands, Ergonomics & Reporter Specification

> **Kode Dokumen:** `SPEC-05-CLI`
> **Tahapan:** Fase 5 - Reporter Output & CLI Entrypoint
> **Status:** Ready for Review
> **Standar Rujukan:** POSIX CLI Standards / 12-Factor App CLI / GNU Coding Standards

Dokumen ini mendefinisikan spesifikasi antarmuka baris perintah (*Command-Line Interface* / CLI), sintaks subcommand, opsi ergonomi pemindaian (A-E), mekanisme deteksi TTY, format pelaporan (*inline ANSI* dan *JSON stream*), serta kontrak *exit codes*.

---

## 1. Spesifikasi Perintah & Subcommand CLI

Binary utama bernama **`charites`**.

```text
Usage:
  charites [command] [flags] [path]

Available Commands:
  scan        Pindai file frontend untuk audit kualitas dan token semantik (Default)
  check       Alias identik untuk 'scan'
  run         Alias identik untuk 'scan'
  version     Cetak versi kompilasi binary, commit git, dan Go runtime
  help        Bantuan penggunaan perintah
```

### 1.1. Aturan Resolusi Perintah Bawaan (Default Command)
Jika `charites` dipanggil langsung dengan argumen path atau flag tanpa subcommand (misal: `charites .` atau `charites src/Button.tsx`), sistem **MUST** mengarahkan eksekusi ke subcommand `scan`.

---

## 2. Spesifikasi Opsi Pemindaian & Ergonomi CLI (Opsi A-E)

Subcommand `scan` (beserta alias `check` dan `run`) **MUST** mendukung opsi berikut:

### 2.1. Opsi A: Direct File Targeting
- Pengguna dapat menentukan target berkas spesifik:
  ```bash
  charites scan src/components/Header.astro
  ```
- Scanner langsung mengevaluasi berkas tersebut tanpa melakukan penelusuran pohon direktori.

### 2.2. Opsi B: Filter Ekstensi (`--ext`)
- Menyaring berkas yang diproses berdasarkan ekstensi:
  ```bash
  charites scan . --ext=astro
  charites scan . --ext=tsx,jsx
  ```

### 2.3. Opsi C: Filter Kategori (`--category`)
- Membatasi evaluasi rule hanya pada domain/kategori tertentu:
  ```bash
  charites scan . --category=theme
  charites scan . --category=a11y
  ```

### 2.4. Opsi D: Filter Single Rule (`--rule`)
- Mengeksekusi satu rule spesifik berdasarkan Semgrep ID:
  ```bash
  charites scan . --rule=theme.hardcode-opacity-color
  ```

### 2.5. Opsi E: Subcommand Aliases
- `charites check <path>` dan `charites run <path>` memiliki perilaku, flag, dan output yang identik 100% dengan `charites scan <path>`.

---

## 3. Matriks Flag Subcommand `scan`

| Flag | Shorthand | Tipe | Nilai Bawaan | Deskripsi |
| :--- | :---: | :---: | :---: | :--- |
| `--format` | `-f` | `string` | `inline` | Format output: `inline` (ANSI) atau `json` |
| `--ext` | `-e` | `string` | `astro,tsx,jsx` | Filter ekstensi yang dipindai (koma terpisah) |
| `--category` | `-c` | `string` | `""` (semua) | Filter kategori rule (`theme`, `a11y`, `perf`, dll.) |
| `--rule` | `-r` | `string` | `""` (semua) | Filter satu Semgrep ID rule spesifik |
| `--config` | | `string` | `charites.yaml` | Path kustom berkas konfigurasi |
| `--ignore` | | `string` | `""` | Pola glob ignore tambahan via CLI |
| `--no-color` | | `bool` | `false` | Matikan pewarnaan ANSI di terminal |
| `--fail-on-warn` | | `bool` | `false` | Return exit code 1 jika hanya ada warning |

---

## 4. Spesifikasi Reporter Output

Charites menyediakan dua presenter output:

### 4.1. Inline ANSI Reporter (`--format=inline`, Default)
Ditujukan untuk kenyamanan mata pengembang di terminal interaktif:

```text
[ERROR] src/pages/index.astro:14:8 [theme.hardcode-opacity-color]
  Hardcode opacity color (bg-primary/10) - wajib pakai semantic token dari global.css
  Hint: Ganti dengan token semantik: primary/10 → primary-light

[WARN] src/components/Card.tsx:42:12 [theme.hardcode-color]
  Hardcode hex color (#2563eb) - gunakan token warna dari global.css
  Hint: Ganti dengan token semantik: #2563eb → bg-primary

 Scanned 28 files in 18ms. Found 1 error, 1 warning.
```

- **Palet Warna ANSI:**
  - `[ERROR]`: Merah tebal (*Bold Red*, ANSI 31;1).
  - `[WARN]`: Kuning tebal (*Bold Yellow*, ANSI 33;1).
  - `[INFO]`: Biru/Sian tebal (*Bold Cyan*, ANSI 36;1).
  - File path & baris: Teks putih/terang dengan underline atau dim.
  - Hint: Abu-abu redup (*Dim*, ANSI 2).
- **Deteksi TTY & No-Color:**
  Jika `stdout` bukan terminal interaktif (misal dialihkan ke file/pipe: `charites scan . > report.txt`), atau terdapat variabel lingkungan `NO_COLOR=1`, sistem **MUST** menonaktifkan kode escape ANSI secara otomatis.

### 4.2. JSON Stream Reporter (`--format=json`)
Ditujukan untuk integrasi CI/CD, tool pelaporan otomatis, dan bot PR:

```json
{
  "version": "1.0.0",
  "summary": {
    "scanned_files": 28,
    "duration_ms": 18,
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
      "message": "Hardcode opacity color (bg-primary/10) - wajib pakai semantic token dari global.css",
      "hint": "Ganti dengan token semantik: primary/10 → primary-light"
    }
  ]
}
```

---

## 5. Kontrak Exit Codes

Untuk kepatuhan standar POSIX dan otomasi pipeline CI:

| Exit Code | Kondisi Keluar | Keterangan |
| :---: | :--- | :--- |
| **`0`** | **CLEAN / SUCCESS** | Tidak ada pelanggaran bertingkat `error`. Jika ada `warning`, tetap `0` kecuali `--fail-on-warn` aktif. |
| **`1`** | **VIOLATIONS DETECTED** | Ditemukan minimal satu pelanggaran bertingkat `error`, atau ada `warning` saat `--fail-on-warn` diaktifkan. |
| **`2`** | **FATAL SYSTEM ERROR** | Kesalahan konfigurasi, argumen tidak valid, atau berkas target tidak dapat diakses. |
