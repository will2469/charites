# 01-SPEC: Core System Requirements Specification

> **Dokumen Status:** Active / Draft
> **Standar Rujukan:** IETF RFC 2119 / RFC 8174
> **Domain:** Rekayasa Kebutuhan Sistem & Spesifikasi Inti Mesin Charites

Dokumen ini mendefinisikan spesifikasi kebutuhan fungsional dan non-fungsional untuk **Charites**, sebuah static analysis engine mandiri performa tinggi berbasis Go 1.26 untuk mendeteksi pelanggaran desain sistem, web vitals, a11y, dan pola frontend pada kode sumber Astro, TypeScript/JSX, dan Tailwind CSS.

---

## 1. Konvensi Notasi (IETF RFC 2119)

Kata kunci **"MUST"** (WAJIB), **"MUST NOT"** (DILARANG), **"REQUIRED"** (HARUS), **"SHALL"** (AKAN), **"SHOULD"** (SEBAIKNYA), **"SHOULD NOT"** (SEBAIKNYA TIDAK), dan **"MAY"** (BOLEH) dalam dokumen ini ditafsirkan sesuai dengan pedoman resmi [IETF RFC 2119](https://www.ietf.org/rfc/rfc2119.txt).

---

## 2. Ruang Lingkup Sistem & Target Ingestion

Charites memproses berkas sumber frontend pada direktori proyek lokal:

1. **Target Ekstensi Berkas:**
   - Sistem **MUST** mendukung pemindaian berkas berekstensi `.astro`, `.tsx`, `.jsx`.
   - Sistem **MAY** mendukung pemindaian berkas `.ts` dan `.js` jika berisi deklarasi markup/template literal.
   - Sistem **MUST** membaca berkas `global.css` untuk mengekstrak definisi semantic token Tailwind CSS v4 (`@theme`).
   - Sistem **MUST** membaca konfigurasi kustom `charites.yaml` (opsional) untuk menyesuaikan severity atau menonaktifkan rule.
2. **Aturan Pengabaian (Ignore System):**
   - Sistem **MUST** secara *default* mengabaikan direktori `node_modules/`, `.git/`, `dist/`, `.astro/`, dan build artifacts lainnya.
   - Sistem **MUST** membaca berkas `.charitesignore` (dengan fallback ke `.gitignore` atau `.ignore`) pada akar repositori untuk menyaring berkas yang tidak perlu dipindai.

---

## 3. Kebutuhan Fungsional (Functional Requirements)

### FR-01: Ekstraksi Semantic Design Tokens (Tailwind v4 `@theme`)
- Sistem **MUST** memiliki parser CSS ringan untuk mengekstrak custom property tokens di dalam blok `@theme` pada `global.css`.
- Token yang diekstrak mencakup warna semantik (misal: `primary`, `secondary`, `destructive`, `muted`, `accent`, `border`, dsb.).
- Token ini **MUST** dijadikan whitelist validasi untuk rule validasi tema.

### FR-02: Normalisasi AST ke Intermediate Representation (IR)
- Parser Astro dan TSX **MUST** mengekstrak struktur tag elemen, atribut, dan posisi baris/kolom.
- Sistem **MUST** menormalkan AST mentah ke dalam kontrak data tunggal: `ir.Node`.
- Elemen markup dari `.astro` (template HTML/JSX) dan `.tsx` (React JSX) **MUST** menghasilkan representasi `ir.Node` yang identik sehingga rule evaluasi tidak terikat pada format bahasa asal.

### FR-03: Kontrak Evaluasi Rule (Single Source of Truth)
- Seluruh rule audit **MUST** mengimplementasikan interface Go tunggal:
  ```go
  type Rule interface {
      ID() string
      Category() string
      DefaultSeverity() ir.Severity
      Evaluate(node *ir.Node) []ir.Diagnostic
  }
  ```
- Evaluasi rule **MUST** bersifat fungsi murni (*pure function*): menerima `*ir.Node`, mengembalikan irisan (*slice*) `[]ir.Diagnostic` tanpa efek samping I/O disk.
- Sistem penamaan rule **MUST** menggunakan format Charites Rule ID `<category>.<rule-slug>` (contoh: `theme.hardcode-opacity-color`).

### FR-04: Reference Proving Rule (`theme.hardcode-opacity-color`)
Sebagai pembuktian stabilitas pipa (*proving ground*), sistem **MUST** mengimplementasikan Rule #1:
- Mendeteksi kelas warna dengan slash opacity langsung (contoh: `bg-primary/10`, `text-destructive/20`, `border-accent/10`).
- Memvalidasi terhadap kamus token semantik:
  - `primary/10` $\rightarrow$ `primary-light`
  - `primary/5` $\rightarrow$ `primary-subtle`
  - `destructive/10` $\rightarrow$ `destructive-light`
  - `accent/10` $\rightarrow$ `accent-light`, dll.
- Menghasilkan diagnosis `ERROR` dengan rekomendasi remedi eksplisit.
- Mendukung pengabaian inline (*inline ignore directives*) multi-rule:
  `// charites:ignore theme.hardcode-opacity-color, theme.hardcode-color`
  atau di Astro template:
  `<!-- charites:ignore theme.hardcode-opacity-color -->`.

### FR-05: Command-Line Interface (CLI) & Exit Codes
Sistem **MUST** menyediakan antarmuka baris perintah dengan subcommands:
- `charites scan [options] [path]`: Memindai berkas atau direktori proyek.
  - **Auto-Detection (Default):** Secara default memindai seluruh berkas frontend (`.astro`, `.tsx`, `.jsx`) dan mengekstrak `global.css` secara otomatis dalam satu kali jalan (*single pass*).
  - **A. Direct File Target:** Mendukung pemindaian berkas tunggal (contoh: `charites scan src/pages/index.astro`).
  - **B. Extension Filter (`--ext`):** Memfilter pemindaian berdasarkan tipe ekstensi (contoh: `--ext=astro`, `--ext=tsx`).
  - **C. Category Filter (`--category`):** Memfilter evaluasi berdasarkan kategori rule (contoh: `--category=a11y`, `--category=theme`, `--category=perf`).
  - **D. Single Rule Filter (`--rule`):** Memfilter evaluasi hanya pada satu rule spesifik (contoh: `--rule=theme.hardcode-opacity-color`).
  - **E. Subcommand Aliases:** Menyediakan alias `charites check` dan `charites run` yang ekuivalen 100% dengan `charites scan`.
  - **Output Formatting (`-f, --format`):** Mendukung format `inline` (default terminal teks berwarna ANSI) dan `json` (flat structured array untuk CI/CD dan jq).
  - **Ignore Override (`--ignore`):** Menambahkan pola ignore kustom saat eksekusi CLI.
- `charites version`: Menampilkan versi rilis, git commit hash, dan Go runtime version.
- **Standar Exit Codes:**
  - `0`: Sukses, tidak ada violation ditemukan (*Clean*).
  - `1`: Pemindaian selesai, ditemukan violation berkategori `ERROR` atau `WARN`.
  - `2`: Kegagalan fatal sistem (panic tertangkap, konfigurasi invalid, atau I/O error).

### FR-06: Antarmuka Server MCP (Model Context Protocol)
- Sistem **MUST** menyediakan subcommand `charites mcp` yang menjalankan server Stdio interaktif.
- Server **MUST** mematuhi protokol **JSON-RPC 2.0** (spesifikasi MCP 2026-07-28).
- Server **MUST** mengekspos tool:
  - `charites_scan`: Menjalankan pemindaian pada path tertentu dan mengembalikan diagnostic dalam payload JSON-RPC.
  - `charites_explain_rule`: Mengembalikan deskripsi rule dan rekomendasi perbaikan (*remediation*).

### FR-07: Konfigurasi Proyek (`charites.yaml`) & Invarian Default-ON
Sistem **MUST** mendukung pemuatan konfigurasi kustom `charites.yaml` (atau `charites.yml`) dengan aturan:
1. **Invarian Default-ON (Model Argus):**
   - Jika berkas `charites.yaml` tidak ada, atau sebuah rule tidak disebutkan di dalam file konfigurasi, maka rule tersebut **OTOMATIS AKTIF (Default: YES)** menggunakan `DefaultSeverity` bawaannya.
2. **Mekanisme Override Rule:**
   - Konfigurasi murni digunakan untuk penyesuaian (*overrides*):
     - **Nonaktifkan Rule:** `theme.hardcode-opacity-color: off` (atau `false`)
     - **Ubah Severity:** `theme.hardcode-color: warn`
3. **Mekanisme Path Ignore Tambahan:**
   - Block `ignore:` dapat mendaftarkan pola path tambahan yang diabaikan di luar `.charitesignore`:
     ```yaml
     ignore:
       - "legacy-vendor/**"
       - "tests/fixtures/**"
     ```

---

## 4. Kebutuhan Non-Fungsional (Non-Functional Requirements)

| ID | Metrik | Target Spesifikasi |
| :--- | :--- | :--- |
| **NFR-01** | **Integritas Adversarial (No Happy-Path Gaming)** | **Zero tolerance for false passes.** Dilarang keras "mengakali" kode rule atau test runner agar sekadar terlihat aman/hijau. Rule **MUST** mengevaluasi invarian secara murni dan jujur tanpa *hardcoded bypass/whitelist rahasia* (contoh dosa legacy: `if file.includes("OpenApiDocs")`). Seluruh eksepsi atau pengecualian berkas **MUST** dideklarasikan secara transparan melalui **Ignore Pattern resmi** (`.ignore`, `.charites.yaml`, atau inline ignore directive). |
| **NFR-02** | **Kecepatan Pindai (Latency)** | < **100 milidetik** untuk memindai 1.000 berkas frontend di disk lokal SSD. |
| **NFR-03** | **Batas Memori (RAM RSS)** | < **50 MB** resident set size saat memindai repositori skala menengah. |
| **NFR-04** | **Zero Runtime Dependency** | Dikeluarkan sebagai **single static binary** tanpa butuh Node.js, npm, Java, atau runtime eksternal. |
| **NFR-05** | **Presisi Lokasi (Accuracy)** | Nomor `Line` dan `Column` pada temuan `Diagnostic` **MUST** 100% cocok dengan posisi kode sumber asli. |
| **NFR-06** | **Deterministik & Idempoten** | Pemindaian berulang terhadap berkas yang sama **MUST** menghasilkan diagnostic identik dengan urutan stabil. |

---

## 5. Batasan Sistem (Out of Scope)

Untuk menjaga performa dan fokus mesin, Charites secara tegas **TIDAK** mencakup:
1. **Bukan Code Formatter**: Charites tidak memformat whitespace/indentasi kode (tugas Prettier / Biome).
2. **Bukan Type Checker**: Charites tidak melakukan evaluasi type checking lengkap TypeScript (tugas `tsc`).
3. **Bukan Bundler**: Charites tidak melakukan bundling aset, tree-shaking, atau minifikasi.
4. **Bukan Dynamic Runtime Evaluator**: Seluruh analisis bersifat statis pada kode sumber mentah (compile-time static analysis).
