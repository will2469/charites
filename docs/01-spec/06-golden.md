# 01-SPEC: 06 - Full Pipeline Integration, Golden Snapshots & Fuzzing Specification

> **Kode Dokumen:** `SPEC-06-GOLDEN`
> **Tahapan:** Fase 6 - Validasi Penuh & Golden Snapshots (Milestone Selesai Pipa)
> **Status:** Ready for Review
> **Standar Rujukan:** Golden Master Testing Pattern / Go 1.26 Native Fuzzing Specification

Dokumen ini mendefinisikan spesifikasi pengujian integrasi pipa compiler dari ujung ke ujung (*end-to-end pipeline*), standarisasi berkas **Golden Snapshots** (`tests/golden/*`), korpus berkas percontohan (*fixtures*), ketahanan *native fuzzing*, serta kriteria kelulusan **Gerbang Stabilitas Pipa** (*Pipeline Stability Gate*).

---

## 1. Spesifikasi Golden Master Testing (`tests/golden/`)

Untuk mencegah terjadinya regresi diagnosis atau pergeseran pelaporan (*diagnostic drift*):

### 1.1. Format Berkas Golden Snapshot
Setiap skenario pemindaian uji coba wajib memiliki snapshot kebenaran mutlak (*ground truth*):
1. **JSON Snapshot (`<scenario>.golden.json`):**
   Menyimpan seluruh payload struktur `ScanResult` (array `diagnostics`, `summary`, `passed`) secara presisi hingga baris dan kolom.
2. **Terminal ANSI Snapshot (`<scenario>.golden.txt`):**
   Menyimpan representasi teks mentah yang dicetak ke terminal (termasuk kode warna ANSI atau mode teks tanpa warna) untuk memastikan tidak ada perubahan visual atau degradasi keterbacaan yang tidak disengaja.

### 1.2. Mekanisme Pembaruan Terkontrol (`-update` Flag)
- Secara bawaan, eksekusi `go test ./tests/...` membandingkan output scanner terkini dengan isi berkas `.golden.*`.
- Jika terdapat perbedaan sekecil apapun (pergeseran nomor baris, perubahan kata pesan, dsb.), test **MUST FAIL**.
- Pembaruan berkas golden snapshot hanya diizinkan jika pengembang menambahkan flag eksplisit:
  ```bash
  go test ./tests/... -run TestGolden -update
  ```
- Setiap perubahan pada berkas golden snapshot **MUST** ditinjau dalam git diff sebelum di-commit.

---

## 2. Spesifikasi Korpus Fixtures Komprehensif (`tests/fixtures/`)

Direktori fixture merepresentasikan variasi kode dunia nyata yang akan dipindai oleh Charites:

### 2.1. Korpus Astro (`tests/fixtures/astro/`)
- `clean.astro`: Komponen Astro bersih yang mematuhi seluruh aturan desain (token semantik, accessible markup).
- `opacity_violations.astro`: Komponen yang memuat variasi kelas terlarang: `bg-primary/10`, `border-destructive/20`, `text-accent/5`.
- `complex_frontmatter.astro`: Komponen dengan blok `---` lebih dari 50 baris untuk memvalidasi presisi line offset.
- `inline_ignore.astro`: Komponen dengan komentar `<!-- charites:ignore theme.hardcode-opacity-color -->`.
- `nested_slots.astro`: Template bersarang kompleks dengan `<slot />` dan komponen Astro internal.

### 2.2. Korpus TSX / JSX (`tests/fixtures/tsx/`)
- `clean.tsx`: Komponen React bersih tanpa pelanggaran.
- `opacity_violations.tsx`: Komponen dengan pelanggaran pada atribut `className`.
- `template_literals.tsx`: Komponen yang memanfaatkan string backtick: `` `p-4 ${active ? "bg-primary/10" : ""}` ``.
- `inline_ignore.tsx`: Komponen dengan komentar `// charites:ignore theme.hardcode-opacity-color`.

### 2.3. Korpus Konfigurasi & CSS (`tests/fixtures/config/`)
- `global.css`: Definisi blok `@theme` resmi Tailwind CSS v4.
- `charites.yaml`: Contoh konfigurasi override status dan penambahan ignore path.
- `.charitesignore`: Contoh pola pengabaian file dan direktori kustom.

---

## 3. Spesifikasi Native Fuzzing Suite (`tests/fuzz/`)

Untuk memastikan parser dan IR builder kebal dari input acak (*chaos resilience*):

1. **Target Fuzzing:**
   - `FuzzAstroPipeline`: Mengalirkan byte acak ke Astro frontmatter splitter, HTML lexer, hingga IR builder.
   - `FuzzTSXPipeline`: Mengalirkan byte acak ke TSX JSX visitor hingga IR builder.
2. **Invarian Tanpa Panic:**
   Input apapun-termasuk kode HTML tidak ditutup, kutipan gantung, sintaks JavaScript cacat, atau byte biner acak-**DILARANG KERAS** memicu panic di Go runtime. Parser wajib pulih secara anggun (*graceful recovery*).
3. **Durasi Eksekusi CI:**
   Fuzz testing wajib dijalankan minimal **60 detik per modul** sebelum rilis binary.

---

## 4. Kriteria Kelulusan Gerbang Stabilitas Pipa (Stability Gate)

Fase 6 adalah titik balik arsitektur Charites. Seluruh pipa pemindaian dinyatakan **SELESAI & STABIL** jika:
1. Seluruh pengujian Golden Snapshot lulus 100%.
2. Pengujian Fuzzing lulus tanpa crash.
3. Seluruh alur CLI (Fase 5) terbukti terhubung mulus dengan Scanner (Fase 4), Rules (Fase 3), Parser (Fase 2), dan IR (Fase 1).
4. Binary dikunci (*locked*) dan siap memasuki Fase 7 (Ekosistem MCP & Wiki) dan Fase 8 (Ekspansi 30+ Rules).
