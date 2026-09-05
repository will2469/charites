# 01-SPEC: 05 - CLI Interface, Commands, Ergonomics & Reporter Specification

> **Kode Dokumen:** `SPEC-05-CLI`
> **Tahapan:** Fase 5 - Reporter Output & CLI Entrypoint
> **Peran Pilar:** SPEC = WHAT (Spesifikasi Antarmuka CLI, Format Reporter & Kontrak Exit Code)
> **Status:** Ready / Approved for Implementation
> **Standar Rujukan:** POSIX CLI Standards / 12-Factor App CLI / GNU Coding Standards

Dokumen ini mendefinisikan spesifikasi antarmuka baris perintah (*Command-Line Interface* / CLI), sintaks subcommand, opsi ergonomi pemindaian (A-E), tata bahasa flag dan normalisasi, aturan resolusi warna/TTY, format dokumen pelaporan (*inline ANSI* dan *JSON Document*), serta kontrak *exit codes*.

---

## 1. Spesifikasi Perintah & Subcommand CLI

Binary utama bernama **`charites`**.

```text
Usage:
  charites [command] [flags] [path]

Available Commands:
  scan        Pindai berkas frontend untuk audit kualitas dan token semantik (Default)
  check       Alias identik untuk 'scan'
  run         Alias identik untuk 'scan'
  update      Periksa dan perbarui biner Charites ke versi terbaru (Alias: 'upgrade')
  upgrade     Alias identik untuk 'update'
  uninstall   Copot pemasangan Charites dari sistem secara tuntas tanpa residu
  version     Cetak versi kompilasi binary, commit git, dan Go runtime
  help        Bantuan penggunaan perintah
```

### 1.1. Aturan Resolusi Perintah Bawaan (Default Command & Aliases)
1. **Pemanggilan Kosong (0 Argumen):** Pemanggilan `charites` tanpa argumen sama sekali **MUST** secara otomatis mengeksekusi `charites scan .`.
2. **Argumen Path / Flag Langsung Tanpa Subcommand:**
   - Pemanggilan dengan path langsung (misal: `charites .` atau `charites src/Button.tsx`) **MUST** diarahkan ke subcommand `scan` dengan path tersebut sebagai target.
   - Pemanggilan dengan flag langsung (misal: `charites --format=json` atau `charites --ext=astro src/`) **MUST** diarahkan ke subcommand `scan` dengan meneruskan seluruh flag dan target path default `.` jika path tidak ditentukan.
3. **Subcommand Aliases:**
   - Subcommand `check` dan `run` memiliki perilaku, flag, output, dan kode keluar yang identik 100% dengan `scan`.
   - Subcommand `upgrade` memiliki perilaku, flag, output, dan kode keluar yang identik 100% dengan `update`.
4. **Subcommand Bantuan & Versi:**
   - `charites version`, `charites -v`, `charites --version` mencetak string versi ke `stdout` dan mengembalikan exit code `0`.
   - `charites help`, `charites -h`, `charites --help` mencetak petunjuk penggunaan ke `stdout` dan mengembalikan exit code `0`.
5. **Perintah / Flag Tidak Dikenal:**
   - Argumen subcommand yang tidak dikenal (misal: `charites foobar`) **MUST** mencetak pesan kesalahan ke `stderr` dan keluar dengan exit code `2`:
     `charites: error: unknown command "foobar". Run 'charites --help' for usage.`
   - Flag yang tidak dikenal (misal: `charites scan --unknown`) **MUST** mencetak pesan kesalahan ke `stderr` dan keluar dengan exit code `2`.
6. **Batasan Target Path:**
   - Subcommand `scan` menerima maksimal 1 target path posisi (`charites scan [flags] [path]`). Jika target path tidak disertakan, nilai bawaan adalah `.` (root workspace).
   - Jika pengguna menyertakan lebih dari satu target path (misal `charites scan path1 path2`), CLI **MUST** keluar dengan exit code `2` dan pesan error di `stderr`:
     `charites: error: multiple scan targets not supported. Specify a single path.`
   - Jika path target tidak ada di filesystem, CLI **MUST** keluar dengan exit code `2` dan pesan error di `stderr`:
     `charites: error: scan target "<path>" does not exist.`
   - Jika path target berada di dalam direktori *builtin hard exclusion* (misal `charites scan node_modules/foo/Card.tsx`), sistem keamanan *direct-target safety* **MUST** menolak pemindaian dan keluar dengan exit code `2` dan pesan error di `stderr`:
     `charites: error: scan target "<path>" is within excluded directory (builtin hard exclusion).`

### 1.2. Kontrak Siklus Hidup Biner & Zero Residual Footprint (SPEC-05-LIFECYCLE-001)
1. **Zero Host Pollution Invariant:** Biner CLI `charites` beroperasi secara murni *in-memory* dan *in-workspace*. CLI dilarang keras membuat berkas state, cache, atau konfigurasi global di luar target direktori yang sedang dipindai (dilarang membuat `~/.config/charites`, `~/.cache/charites`, `%APPDATA%\charites`, dsb.).
2. **Ephemeral Run-to-Completion:** Setiap pemanggilan CLI (subcommand `scan`, `version`, `help`) bersifat *run-to-completion*. Saat proses mengembalikan kode keluar, seluruh alokasi memori dilepaskan oleh OS dan tidak ada proses latar belakang, daemon, socket, atau berkas sementara di `/tmp` yang ditinggalkan.
3. **Pembaruan & Pencopotan Bersih (Clean Update & Uninstall Contract):**
   - **Subcommand `update` (Alias: `upgrade`):**
     - Memeriksa ketersediaan versi rilis terbaru di repository GitHub resmi.
     - Jika biner sudah merupakan versi terbaru atau tidak ditemukan versi baru:
       Mencetak: `No update found. Charites is up to date.` ke `stdout` dan keluar dengan exit code `0`.
     - Jika ditemukan versi baru:
       Mengunduh biner artefak yang sesuai dengan OS/Arsitektur saat ini, menimpa biner lama secara atomik (`os.Rename`), dan mencetak: `Charites updated to <version> successfully.` dengan exit code `0`.
     - Flag opsional `--check` / `-c`: Hanya memeriksa versi tanpa mengunduh.
   - **Subcommand `uninstall`:**
     - Menghapus satu-satunya berkas biner `charites` dari sistem operasi host (`os.Remove(selfPath)`).
     - Menjamin sistem pengguna kembali 100% bersih tanpa ada *leftover* berkas atau cache (*Zero Leftover Guarantee*).
     - Mencetak:
       ```text
       Charites uninstalled successfully.
       0 residual files or caches remaining.
       ```
     - Keluar dengan exit code `0` (atau exit code `2` jika terdapat kegagalan permission OS).
     - Flag opsional `--yes` / `-y`: Melewati konfirmasi prompt interaktif.

---

## 2. Matriks Flag Subcommand `scan` & Aturan Normalisasi

| Flag | Shorthand | Tipe | Nilai Bawaan | Deskripsi |
| :--- | :---: | :---: | :---: | :--- |
| `--format` | `-f` | `string` | `inline` | Format output: `inline` (ANSI) atau `json` (Dokumen JSON) |
| `--ext` | `-e` | `string` | `astro,tsx,jsx` | Filter ekstensi yang dipindai (koma terpisah atau berulang) |
| `--category` | `-c` | `string` | `""` (semua) | Filter kategori rule (`theme`, `a11y`, `perf`, dll.) |
| `--rule` | `-r` | `string` | `""` (semua) | Filter satu Charites Rule ID spesifik (`<category>.<slug>`) |
| `--config` | | `string` | `charites.yaml` | Path kustom berkas konfigurasi |
| `--ignore` | | `string` | `""` | Pola glob ignore tambahan (koma terpisah atau berulang) |
| `--no-color` | | `bool` | `false` | Matikan pewarnaan ANSI di terminal |
| `--fail-on-warn` | | `bool` | `false` | Return exit code 1 jika hanya ada warning |

### 2.1. Normalisasi & Validasi Flag `--ext`
- **Ekstensi yang Didukung:** `.astro`, `.tsx`, `.jsx`.
- **Normalisasi:** Case-insensitive (`ASTRO` $\rightarrow$ `astro`), tanda titik di awal bersifat opsional (`astro` $\equiv$ `.astro`). Nilai dapat dipisahkan koma (`--ext=astro,tsx`) atau berulang (`--ext astro --ext tsx`).
- **Penolakan Nilai Tidak Sah:**
  - Jika pengguna memasukkan ekstensi yang tidak didukung (misal `--ext=vue` atau `--ext=foo`), CLI **MUST** menghentikan eksekusi dengan exit code `2` dan pesan error di `stderr`:
    `charites: error: unsupported extension "foo". Supported extensions: .astro, .tsx, .jsx.`
  - Jika nilai flag kosong (`--ext=""` atau `--ext=`), CLI **MUST** keluar dengan exit code `2`:
    `charites: error: empty extension flag.`

### 2.2. Validasi Konflik & Keberadaan Flag `--category` dan `--rule`
- **Validasi Keberadaan Kategori:** Jika `--category` ditentukan dan kategori tersebut tidak terdaftar pada registry, CLI **MUST** keluar dengan exit code `2`:
  `charites: error: unknown category "<category>".`
- **Validasi Keberadaan Rule:** Jika `--rule` ditentukan dan Rule ID tersebut tidak terdaftar pada registry, CLI **MUST** keluar dengan exit code `2`:
  `charites: error: unknown rule "<rule>".`
- **Validasi Irisan (Intersection Check):** Jika pengguna menentukan kedua flag sekaligus (misal: `--category=theme --rule=a11y.alt-text`):
  - Sistem melakukan validasi irisan: apakah rule tersebut terdaftar di dalam kategori yang ditentukan?
  - Jika rule tidak termasuk dalam kategori tersebut, CLI **MUST** keluar dengan exit code `2` dan pesan error di `stderr`:
    `charites: error: rule "a11y.alt-text" does not belong to category "theme".`

### 2.3. Validasi Flag `--format`
- **Format yang Didukung:** `inline` (default) dan `json`.
- Jika format tidak didukung (misal `--format=xml` atau `--format=yaml`), CLI **MUST** keluar dengan exit code `2` dan pesan error di `stderr`:
  `charites: error: unsupported format "<format>". Supported formats: inline, json.`

### 2.4. Validasi Flag `--config`
- Jika path konfigurasi kustom ditentukan via `--config` namun berkas tidak ditemukan atau tidak dapat dibaca, CLI **MUST** keluar dengan exit code `2` dan pesan error di `stderr`:
  `charites: error: config file not found: "<path>".`
- Jika isi berkas konfigurasi tidak dapat di-parse karena sintaks tidak sah, CLI **MUST** keluar dengan exit code `2` dan pesan error parsing di `stderr`.

### 2.5. Tata Bahasa Flag `--ignore` & Presedensi
- Pola ignore tambahan dapat diberikan via argumen `--ignore=pattern` atau berulang `--ignore p1 --ignore p2`.
- Pola dievaluasi secara relatif terhadap root direktori pemindaian menggunakan sintaks glob `.charitesignore`.
- Pola CLI ditambahkan ke aturan pengguna dan **DILARANG** membuka kembali (*override/negate*) direktori builtin hard exclusion (`.git`, `node_modules`, `dist`, dll.). Negasi `!node_modules` diabaikan dan direktori tetap dikecualikan.
- **Presedensi Kebijakan Konfigurasi vs CLI Scope:** Sesuai kontrak [docs/00-CONTRACT.md](https://github.com/will2469/charites/blob/main/docs/00-CONTRACT.md), flag CLI `--rule` dan `--category` berfungsi sebagai *Candidate Scope*. Konfigurasi `charites.yaml` berfungsi sebagai *Policy*. Jika konfigurasi menetapkan status `off` untuk rule tertentu, rule tersebut **TIDAK AKTIF** meskipun pengguna secara eksplisit menentukan `--rule` pada CLI.

---

## 3. Spesifikasi Reporter Output

### 3.1. Inline ANSI Reporter (`--format=inline`, Default)
Format teks interaktif untuk konsol terminal:

```text
[ERROR] src/pages/index.astro:14:8 [theme.hardcode-opacity-color]
  Hardcode opacity color: "bg-primary/10"
  Hint: Use semantic token "primary-light".

[WARN] src/components/Card.tsx:42:12 [theme.hardcode-color]
  Hardcode hex color: "#2563eb"
  Hint: Use semantic token "bg-primary".

 2 problems found (1 error, 1 warning)
  Scanned 28 files in 18ms.
```

- **Format Repositori Bersih (Clean Scan):**
  ```text
   0 problems found (0 errors, 0 warnings)
    Scanned 28 files in 12ms.
  ```
- **Kontrak Formatting:**
  - Path berkas ditampilkan sebagai format relatif terhadap workspace dengan pemisah POSIX forward slash (`/`), bahkan di sistem operasi Windows.
  - Teks menggunakan encoding UTF-8 dengan karakter newline trailing akhir (`\n`).
  - Pewarnaan ANSI: Bold Red (`[ERROR]`), Bold Yellow (`[WARN]`), Bold Cyan (`[INFO]`), Dim (`Hint:`).

### 3.2. Resolusi Warna ANSI & Deteksi TTY
Pewarnaan ANSI escape codes dikendalikan oleh resolusi deterministik:
- Pewarnaan ANSI **DIMATIKAN** (`ColorNever`) jika salah satu kondisi berikut terpenuhi:
  1. Pengguna menyertakan flag `--no-color`.
  2. Variabel lingkungan `NO_COLOR` ada dan tidak kosong (`os.Getenv("NO_COLOR") != ""`).
  3. `stdout` bukan terminal interaktif TTY (misal: dialihkan ke pipe `|` atau file `>`).
- Selain kondisi di atas, pewarnaan ANSI diaktifkan secara otomatis (`ColorAuto`).

### 3.3. JSON Document Reporter (`--format=json`)
Format dokumen JSON tunggal lengkap (*complete JSON document*) yang dicetak di akhir pemindaian untuk konsumsi mesin, PR bot, atau pipeline CI/CD:

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
      "message": "Hardcode opacity color: \"bg-primary/10\"",
      "hint": "Use semantic token \"primary-light\"."
    }
  ]
}
```

- **Skema Waktu:** Durasi pemindaian dicatat dalam field `"duration_ms"` sebagai integer milidetik.
- **Determinis Biner:** Urutan isi slice `diagnostics` mengikuti *total ordering* Fase 4 (`File` $\rightarrow$ `Line` $\rightarrow$ `Col` $\rightarrow$ `RuleID` $\rightarrow$ `Severity` $\rightarrow$ `Message` $\rightarrow$ `Hint`).

---

## 4. Taksonomi & Kontrak Exit Codes

Exit code adalah kontrak deterministik mesin dengan lingkungan sistem operasi dan runner CI/CD:

| Exit Code | Nama Status | Kondisi Pemicu |
| :---: | :--- | :--- |
| **`0`** | **CLEAN / SUCCESS** | Pemindaian selesai tanpa temuan bertingkat `error`. Temuan bertingkat `warning` diizinkan lewat, kecuali jika flag `--fail-on-warn` aktif. |
| **`1`** | **VIOLATIONS FOUND** | Ditemukan minimal 1 temuan bertingkat `error`, ATAU ditemukan minimal 1 temuan bertingkat `warning` saat flag `--fail-on-warn` diaktifkan. |
| **`2`** | **CLI / OPERATIONAL ERROR** | Kesalahan pemanggilan CLI (flag tidak dikenal, nilai flag tidak valid, ekstensi tidak didukung, konflik category/rule, path target tidak dapat diakses, atau sintaks `charites.yaml` korup). |
| **`130`** | **TERMINATED BY SIGNAL** | Proses diinterupsi oleh pengguna via sinyal terminal (`SIGINT`/`SIGTERM`). |

**Invarian Mutlak:** Temuan pelanggaran kode (diagnostic violation) **DILARANG KERAS** menghasilkan exit code `2`.
