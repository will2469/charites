# 01-SPEC: 04 - Configuration, Concurrency Scanner & Traversal Engine Specification

> **Kode Dokumen:** `SPEC-04-ENGINE`
> **Tahapan:** Fase 4 - Konfigurasi, Concurrency Scanner & Traversal Engine
> **Peran Pilar:** SPEC = WHAT (Spesifikasi Kebutuhan Fungsional & Kontrak Engine)
> **Status:** Ready for Review (Implementation Locked: DO NOT START YET)
> **Standar Rujukan:** IETF RFC 2119 / Semgrep Engine Architecture / Gitignore Pattern Specification

Dokumen ini mendefinisikan spesifikasi kebutuhan fungsional untuk sistem konfigurasi **`charites.yaml`** (Model Argus: Invarian Default YES & Precedence), sistem pengabaian **`.charitesignore`**, proteksi traversal filesystem, batas sumber daya I/O, *goroutine worker pool*, *AST traversal engine*, serta kontrak determinisme pelaporan.

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
- **Nonaktif:** `off`, `false`, `disable`, `disabled` $\rightarrow$ Rule dimatikan (dikeluarkan dari active rules).
- **Error:** `error` $\rightarrow$ Pelanggaran memicu severity Error (exit code 1 pada CLI).
- **Warning:** `warn`, `warning` $\rightarrow$ Pelanggaran dilaporkan sebagai Warning.
- **Info:** `info` $\rightarrow$ Pelanggaran dilaporkan sebagai Info/Notice.

### 1.4. Kontrak Presedensi Resolusi Rule (Candidate Selection vs Policy)
Untuk mencegah ambiguitas ketika flag CLI bertabrakan dengan berkas konfigurasi, urutan presedensi ditetapkan sebagai berikut:

$$\text{Registry} \longrightarrow \text{CLI Candidate Scope} \longrightarrow \text{Config Policy} \longrightarrow \text{Active Rules}$$

1. **Tahap 1 (Base Candidates):** Seluruh rule terdaftar di `Registry` menjadi kandidat awal.
2. **Tahap 2 (CLI Selection / Scope):** Jika user menyertakan flag `--category` atau `--rule`, kandidat disaring hanya ke sub-koleksi yang dipilih.
3. **Tahap 3 (Config Policy / Safety):** Konfigurasi `charites.yaml` dievaluasi terhadap kandidat yang lolos:
   - Jika konfigurasi menetapkan status `off` / `false` / `disable`, rule **TIDAK AKTIF** (*Policy* mengalahkan *CLI Selection* demi kepatuhan kebijakan organisasi).
   - Jika konfigurasi menetapkan override severity (`warn`, `error`, `info`), nilai ini menjadi `EffectiveSeverity` rule tersebut.
   - Jika rule tidak disebutkan di konfigurasi, rule tetap aktif dengan `DefaultSeverity()`.
4. **Immutability Invariant:** Rule singleton di registry **DILARANG DIMUTASI**. Override severity dikemas dalam struktur pembungkus `ActiveRule` di level engine.

---

## 2. Spesifikasi Pengabaian Berkas (`.charitesignore` & Boundary Invariants)

### 2.1. Direktori Default Builtin (Immutable Hard Exclusion)
Meskipun berkas `.charitesignore` tidak ada atau kosong, scanner **WAJIB** mengabaikan direktori berikut:
- `.git/`
- `node_modules/`
- `dist/`
- `.astro/`
- `.next/`
- `.turbo/`
- `build/`
- `coverage/`

**Invarian Hard Exclusion:** Direktori builtin bersifat mutlak. Pola negasi (`!`) pada `.charitesignore` **DILARANG KERAS** membuka kembali (*re-include*) direktori builtin atau sub-berkas di dalamnya.

### 2.2. Sintaks & Semantik Evaluasi `.charitesignore`
Format `.charitesignore` mendukung subset spesifikasi `.gitignore` yang deterministik:
- Komentar diawali `#` dan baris kosong diabaikan.
- Trailing slash `dir/` hanya mencocokkan direktori.
- Leading slash `/path` atau `path` mengacu relatif terhadap root workspace.
- Single wildcard `*` mencocokkan sembarang karakter dalam satu segmen direktori.
- Deep wildcard `**` mencocokkan nol atau lebih tingkat hierarki direktori.
- Negasi `!` membatalkan pengabaian berkas yang cocok dengan aturan sebelumnya.

**Semantik Evaluasi Berurutan (Sequential Evaluation):**
Aturan dievaluasi secara sekuensial dari baris pertama hingga terakhir. Aturan yang cocok paling akhir adalah yang menentukan (*last matching rule wins*). Namun, jika direktori leluhur (*ancestor directory*) telah dipangkas pada level traversal, aturan negasi pada berkas di dalamnya tidak berlaku.

### 2.3. Invarian Early Directory Pruning
Jika sebuah path direktori cocok dengan pola ignore (misal `node_modules/` atau `vendor/`), walker **DILARANG KERAS** memanggil fungsi baca isi direktori (`os.ReadDir`). Direktori langsung dipangkas (*pruned*) dari traversal.

### 2.4. Kebijakan Traversal Symlink (Security-First Policy)
Untuk mencegah serangan *directory traversal*, *symlink race condition*, dan *infinite loop*:
1. **Directory Symlink:** `DO NOT FOLLOW`. Symlink yang menunjuk ke direktori langsung diabaikan dan tidak ditelusuri.
2. **File Symlink:** Diabaikan secara default (`skip by default`).

### 2.5. Invarian Target Berkas Langsung (Explicit Target Safety)
Jika user memanggil Charites dengan path berkas langsung (contoh: `charites scan node_modules/foo/Button.tsx` atau `charites scan .git/hooks/pre-commit`):
- Panggilan **DILARANG KERAS** menjadi celah (*escape hatch*) untuk menerobos *builtin hard exclusion*.
- **Algoritma Ancestry Inspection:** Sebelum memulai scanning atau traversal, scanner memecah `filepath.Clean(target)` menjadi segmen path. Jika ADA segmen yang cocok dengan salah satu `builtinExclusions` (`.git`, `node_modules`, `dist`, `.astro`, `.next`, `.turbo`, `build`, `coverage`):
  - Scanner menolak pemindaian secara langsung sebelum pembacaan I/O atau parsing dilakukan.
  - Walker mengembalikan error keamanan dan tepat **0 jobs** dimasukkan ke dalam antrean kerja `jobs`.
  - CLI keluar dengan status error non-zero (exit code 2).
- **Target Berkas Tunggal yang Sah:** Jika target adalah berkas tunggal yang sah (di luar direktori terlarang), walker memvalidasi ekstensi (`.astro`, `.tsx`, `.jsx`) dan ukuran berkas ($\le 10\text{ MB}$), langsung memasukkannya ke `jobs`, tanpa menelusuri pohon direktori secara rekursif.

### 2.6. Batas Maksimal Ukuran Berkas (Resource Invariant)
- Batas ukuran berkas sumber frontend: **Maksimal 10 Megabytes** ($10 \times 1024 \times 1024\text{ bytes}$).
- Berkas yang melebihi batas ini dilewati secara otomatis dengan catatan informasi agar tidak memicu Out-of-Memory (OOM) akibat berkas bundel terkompilasi yang salah ditempatkan.

---

## 3. Spesifikasi Concurrency Scanner & Worker Pool

### 3.1. Pemetaan Target Pemindaian (Target Mapping)
1. **Mode Default (`charites scan .`):** Memindai seluruh direktori kerja untuk berkas berekstensi: `.astro`, `.tsx`, `.jsx`.
2. **Mode Target Berkas Langsung:** `charites scan src/components/Button.tsx`.
3. **Mode Filter Ekstensi:** Flag `--ext=astro` (hanya memproses ekstensi target).
4. **Mode Filter Kategori & Rule:** Flag `--category=theme` atau `--rule=theme.hardcode-opacity-color`.

### 3.2. Worker Pool Concurrency Model
- Scanner menggunakan model antrean kerja konkuren berbasis Go Channel:
  ```go
  type Job struct {
      FilePath string
  }
  ```
- **Kapasitas Konkurensi:** Default mengalokasikan $N = \text{runtime.GOMAXPROCS(0)}$ goroutine pekerja, dibatasi dalam rentang $[1, 256]$. Dapat dikonfigurasi melalui flag `--workers=N`.
- **Fan-Out / Fan-In:** File yang lolos filter ignore didistribusikan ke worker pool. Setiap worker membaca berkas, mem-parse ke `*ir.Node`, dan mengeksekusi traversal engine.

### 3.3. Kontrak Pembatalan Interupsi & Kepemilikan Channel (Channel Ownership & Cancellation)
Sistem pemindai mendukung terminasi bersih saat menerima sinyal interupsi terminal (`SIGINT`/`SIGTERM`) melalui `context.Context` dengan **Invarian Single Producer = Single Closer**:
- **Kepemilikan Channel `jobs` (`chan string`):**
  - **Producer Tunggal:** `Walker`.
  - **Closer Tunggal:** `Walker` via `defer close(jobs)` saat traversal selesai atau `ctx.Done()`.
- **Kepemilikan Channel `results` (`chan []ir.Diagnostic`):**
  - **Producers:** $N$ goroutine pekerja dalam `WorkerPool` (fan-in).
  - **Closer Tunggal:** Goroutine koordinator tersinkronisasi `sync.WaitGroup` (`go func() { wg.Wait(); close(results) }()`). Pekerja **DILARANG** menutup channel `results`.
- **Siklus Hidup Pembatalan (`State: RUNNING -> CANCELLING -> STOPPED`):**
  1. Ketika `ctx.Done()` diterima:
     - `Walker` mendeteksi pembatalan, membatalkan traversal, dan menutup channel `jobs`.
     - Worker yang membaca `jobs` mendeteksi `ctx.Done()` atau channel closed, membatalkan pemrosesan aktif, tidak mengirim hasil parsial, dan memanggil `defer wg.Done()`.
     - Koordinator menutup `results` setelah seluruh worker berhenti bersih (*zero goroutine leak*).
     - Hasil diagnostik parsial **DIBUANG** (temuan tidak lengkap dilarang dianggap sah).
     - Program keluar dengan kode status interupsi non-zero (exit code 130).

---

## 4. Spesifikasi AST Traversal Engine & Direktif Inline Ignore

### 4.1. AST Traversal Loop & Effective Severity
Engine menelusuri pohon `*ir.Node` menggunakan iterator Go 1.26:
```go
for node := range root.Walk() {
    for _, active := range activeRules {
        diags := active.Rule.Evaluate(node)
        for _, d := range diags {
            // Terapkan EffectiveSeverity dari konfigurasi
            d.Severity = active.EffectiveSeverity

            // Evaluasi penekanan inline ignore (signature baku: IsIgnored(d, node))
            if !ctx.IsIgnored(d, node) {
                ctx.AddDiagnostic(d)
            }
        }
    }
}
```

### 4.2. Tata Bahasa Direktif Inline Ignore (Grammar & Syntax)
1. **Format JavaScript / TypeScript / TSX:**
   `// charites:ignore <rule-list>`
2. **Format Template Astro HTML:**
   `<!-- charites:ignore <rule-list> -->`
3. **Aturan Sintaks Direktif:**
   - **Wildcard:** `// charites:ignore *` menekan seluruh temuan rule pada cakupan target.
   - **Multi-Rule:** `// charites:ignore theme.hardcode-opacity-color, a11y.alt-text`. Spasi di sekitar koma dipangkas (*trimmed*).
   - **Deduplikasi:** Rule ID duplikat disatukan tanpa error.
   - **Unknown Rule:** Rule ID yang tidak dikenal diabaikan secara senyap tanpa menggagalkan pemindaian.
   - **Empty Directive:** Direktif kosong (`// charites:ignore`) dianggap tidak valid dan tidak menekan pelanggaran apa pun.

### 4.3. Cakupan Penekanan (Scope: Line & AST Node Span)
Untuk mengakomodasi elemen JSX multi-baris:
1. **Same-Line:** Direktif di baris $N$ menekan temuan pada baris $N$ (inline trailing comment).
2. **Next-Line / Node Span Scope:** Jika direktif berada di baris $N$, dan elemen/node AST dimulai pada baris $N+1$ (`Node.Span.StartLine == N+1`), cakupan penekanan meluas ke seluruh rentang baris node tersebut (`Node.Span.StartLine` hingga `Node.Span.EndLine`). Ini menjamin atribut atau class multi-baris ditekan secara akurat.

---

## 5. Spesifikasi Pengurutan Determinis (Total Ordering)

Hasil diagnostic yang dikumpulkan dari berbagai goroutine wajib diurutkan sebelum diserahkan ke reporter menggunakan relasi pengurutan total (*total ordering relation*):

1. **`File`** (string leksikografis menaik)
2. **`Span.StartLine`** (integer menaik)
3. **`Span.StartColumn`** (integer menaik)
4. **`RuleID`** (string leksikografis menaik)
5. **`Severity`** (nilai numerik bobot: Error > Warning > Info)
6. **`Message`** (string leksikografis menaik)
7. **`Hint`** (string leksikografis menaik)

Relasi pengurutan total ini menjamin byte output $100\%$ identik dan idempoten pada input yang sama, tanpa terpengaruh penjadwalan goroutine.
