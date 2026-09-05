# 01-SPEC: 04 - Configuration, Concurrency Scanner & Traversal Engine Specification

> **Kode Dokumen:** `SPEC-04-ENGINE`
> **Tahapan:** Fase 4 - Konfigurasi, Concurrency Scanner & Traversal Engine
> **Status:** Ready for Review
> **Standar Rujukan:** IETF RFC 2119 / Semgrep Engine Architecture / Gitignore Pattern Specification

Dokumen ini mendefinisikan spesifikasi kebutuhan fungsional untuk sistem konfigurasi **`charites.yaml`** (Model Argus: Invarian Default YES), sistem pengabaian **`.charitesignore`**, *fast directory walker*, *goroutine worker pool*, dan *AST traversal engine*.

---

## 1. Spesifikasi Konfigurasi `charites.yaml` (Model Argus: Default YES)

Sistem konfigurasi Charites mengadopsi model **Opt-Out / Override Only**:

### 1.1. Invarian Default YES (Default-ON)
1. **Zero-Config Operability:** Jika berkas `charites.yaml` tidak ditemukan di root workspace, **seluruh rule aktif 100%** menggunakan `DefaultSeverity` bawaan masing-masing rule.
2. **Implicit Rule Inclusion:** Setiap rule yang terdaftar di binary namun tidak disebutkan di dalam `charites.yaml` dianggap aktif secara otomatis tanpa perlu deklarasi manual.

### 1.2. Format Berkas Konfigurasi
Berkas `charites.yaml` (atau `charites.yml`) murni digunakan untuk melakukan penyesuaian (*overrides*):

```yaml
# charites.yaml - Optional override configuration
rules:
  # Menonaktifkan rule tertentu
  theme.hardcode-opacity-color: off

  # Mengubah severity default
  theme.hardcode-color: warn
  a11y.button-accessible-name: error

# Menambahkan pola path yang diabaikan di luar .charitesignore
ignore:
  - "legacy-vendor/**"
  - "tests/fixtures/**"
  - "generated/**/*.ts"
```

### 1.3. Nilai Severity Override yang Sah
- **Nonaktif:** `off`, `false`, `disable`, `disabled` $\rightarrow$ Rule tidak dievaluasi.
- **Error:** `error` $\rightarrow$ Pelanggaran memicu exit code 1.
- **Warning:** `warn`, `warning` $\rightarrow$ Pelanggaran dicetak sebagai peringatan.
- **Info:** `info` $\rightarrow$ Pelanggaran dicetak sebagai informasi catatan.

---

## 2. Spesifikasi Pengabaian Berkas (`.charitesignore` & Early Pruning)

Sistem pengabaian berkas menggunakan format standar yang kompatibel dengan `.gitignore`:

### 2.1. Direktori Default Builtin yang Wajib Diabaikan (Zero Overhead)
Meskipun berkas `.charitesignore` tidak ada, scanner **MUST** mengabaikan direktori berikut secara otomatis:
- `.git/`
- `node_modules/`
- `dist/`
- `.astro/`
- `.next/`
- `.turbo/`
- `build/`
- `coverage/`

### 2.2. Sintaks Pola `.charitesignore`
- Komentar diawali dengan tanda `#`.
- Baris kosong diabaikan.
- Wildcard `*` (mencocokkan sembarang karakter dalam satu segmen direktori).
- Deep wildcard `**` (mencocokkan nol atau lebih segmen direktori bertingkat).
- Trailing slash `dir/` (hanya mencocokkan direktori, bukan berkas bernama sama).
- Negasi `!` (membatalkan pengabaian berkas yang cocok dengan pola sebelumnya).

### 2.3. Invarian Early Directory Pruning
Jika sebuah path direktori cocok dengan pola ignore (contoh: `node_modules/`), walker **DILARANG KERAS** memanggil fungsi baca isi direktori (`os.ReadDir`) pada direktori tersebut. Direktori tersebut langsung dipangkas (*pruned*) dari pohon pemindaian.

---

## 3. Spesifikasi Concurrency Scanner & Worker Pool

### 3.1. Pemetaan Target Pemindaian (Target Mapping)
1. **Mode Default (`charites scan .`):**
   - Menelusuri seluruh direktori kerja.
   - Menyaring berkas berekstensi: `.astro`, `.tsx`, `.jsx`.
   - Membaca `global.css` (jika ada) untuk mengekstrak kamus token `@theme`.
2. **Mode Target Berkas Langsung (Ergonomi A):**
   - `charites scan src/components/Button.tsx`: Langsung memindai berkas tunggal tersebut tanpa melakukan pemindaian pohon direktori.
3. **Mode Filter Ekstensi (Ergonomi B):**
   - Flag `--ext=astro`: Hanya memasukkan berkas berekstensi `.astro` ke antrean pemrosesan.
4. **Mode Filter Kategori & Rule (Ergonomi C & D):**
   - Flag `--category=theme`: Hanya mengeksekusi rule dalam kategori `theme`.
   - Flag `--rule=theme.hardcode-opacity-color`: Hanya mengeksekusi satu rule spesifik.

### 3.2. Worker Pool Concurrency Model
- Scanner menggunakan arsitektur antrean kerja konkuren berbasis Go Channel:
  ```go
  type Job struct {
      FilePath string
  }
  ```
- **Kapasitas Konkurensi:** Secara default mengalokasikan $N = \text{runtime.NumCPU()}$ goroutine pekerja.
- **Fan-Out / Fan-In:** File yang lolos filter ignore didistribusikan ke worker pool. Setiap worker membaca berkas, mem-parse ke `*ir.Node`, dan mengeksekusi engine traversal.
- **Buffer Non-Blocking:** Pengumpulan hasil diagnostic dikumpulkan ke channel hasil thread-safe atau dikelompokkan per-worker untuk menghindari contention mutex global.

---

## 4. Spesifikasi AST Traversal Engine & Direktif Inline Ignore

### 4.1. AST Traversal Loop
Engine menelusuri pohon `*ir.Node` menggunakan iterator standar Go 1.26:
```go
for node := range root.Walk() {
    for _, rule := range activeRules {
        diags := rule.Evaluate(node)
        for _, d := range diags {
            if !ctx.IsIgnored(d.Line, d.Rule) {
                ctx.AddDiagnostic(d)
            }
        }
    }
}
```
Diagnostic yang lolos dari filter inline ignore diserahkan ke pipeline pelaporan (Fase 5).

### 4.2. Sintaks Direktif Inline Ignore
Charites mendukung pengabaian pelanggaran pada baris kode tertentu melalui komentar kode sumber:
1. **Format JavaScript / TypeScript / TSX:**
   ```tsx
   // charites:ignore theme.hardcode-opacity-color
   <div className="bg-primary/10">Legacy Widget</div>
   ```
2. **Format Multi-Rule (Dipisahkan koma):**
   ```tsx
   // charites:ignore theme.hardcode-opacity-color, theme.hardcode-color
   <div className="bg-primary/10 text-[#2563eb]">Custom Exception</div>
   ```
3. **Format Template Astro HTML:**
   ```astro
   <!-- charites:ignore theme.hardcode-opacity-color -->
   <div class="bg-primary/10">Hero Section</div>
   ```

### 4.3. Lingkup Baris (Line Scoping)
- **Same-Line:** Jika direktif diletakkan di baris yang sama dengan elemen pelanggar, pelanggaran pada baris tersebut diabaikan.
- **Next-Line (Previous Line Suppression):** Jika direktif diletakkan tepat satu baris di atas elemen pelanggar (`line_directive + 1 == line_violation`), pelanggaran pada baris di bawahnya diabaikan.
