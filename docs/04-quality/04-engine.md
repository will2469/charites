# 04-QUALITY: 04 - Concurrency Safety, I/O Limits & Engine Performance Invariants

> **Kode Dokumen:** `QUAL-04-ENGINE`
> **Tahapan:** Fase 4 - Konfigurasi, Concurrency Scanner & Traversal Engine
> **Peran Pilar:** QUALITY = QUALITY THRESHOLD (Ambang Batas Kualitas, Keamanan & Anggaran Sumber Daya)
> **Status:** Ready for Review
> **Standar Rujukan:** Go Concurrency Guidelines & Defensive Systems Programming

Dokumen ini mendefinisikan batasan kualitas, ketahanan konkurensi antar-goroutine, pencegahan kebocoran goroutine (*goroutine leak*), batas aman I/O disk, serta kriteria kelulusan cakupan pengujian per-paket pada Fase 4.

---

## 1. Invarian Keamanan Konkurensi & Pencegahan Deadlock

1. **Zero Data Race Invariant:**
   Seluruh interaksi antrean kerja (`jobs`), pengiriman diagnostic (`results`), dan pemetaan konfigurasi wajib bebas dari *data race*. Pengujian `go test -race ./internal/...` wajib lulus 100%.
2. **Share-Nothing Worker Isolation:**
   Setiap goroutine worker memproses satu berkas dalam konteks terisolasi (`*analyzer.Context`). Dilarang membagi pointer state mutable antar worker.
3. **Pencegahan Kebocoran Goroutine (Goroutine Leak Prevention):**
   - Channel `jobs` dan `results` dikelola secara ketat dengan kepemilikan produsen (*producer-owned*).
   - Seluruh worker wajib merespons sinyal `ctx.Done()` untuk memastikan terminasi bersih saat pembatalan interupsi (`SIGINT`/`SIGTERM`).
   - Penutupan channel hanya dilakukan setelah seluruh `sync.WaitGroup` selesai dieksekusi.

---

## 2. Invarian Batasan I/O Disk & Proteksi DoS

1. **Strict Early Directory Pruning:**
   Direktori yang masuk dalam daftar builtin exclusion (seperti `.git/`, `node_modules/`, `dist/`) atau pola ignore direktori **DILARANG KERAS** dibuka dengan `os.ReadDir`. Walker wajib langsung mengembalikan `filepath.SkipDir`.
2. **Batas Maksimum Ukuran Berkas (Max File Size Limit):**
   Charites membatasi ukuran berkas sumber frontend maksimal **10 Megabytes** ($10 \times 1024 \times 1024\text{ bytes}$). Berkas yang melebihi batas ini dilewati tanpa parsing untuk mencegah Out-of-Memory (OOM).
3. **Proteksi Traversal Symlink:**
   Walker tidak mengikuti direktori symlink (`DO NOT FOLLOW`) dan melewati file symlink secara default, mencegah risiko siklus tak terbatas (*circular symlink loops*) dan kebocoran akses di luar workspace.
4. **Proteksi Target Langsung (Direct Target Safety):**
   Pemanggilan target file langsung dilarang menembus *builtin hard exclusion*. Target yang berada di dalam `node_modules` atau `.git` wajib ditolak.

---

## 3. Invarian Output Deterministik (Total Ordering Invariant)

1. **Byte-Level Determinism:**
   Hasil temuan diagnostic wajib diurutkan menggunakan relasi pengurutan total (*total ordering*):
   $$\text{File} \longrightarrow \text{Line} \longrightarrow \text{Column} \longrightarrow \text{RuleID} \longrightarrow \text{Severity} \longrightarrow \text{Message} \longrightarrow \text{Hint}$$
2. **Idempotensi Pelaporan:**
   Dua kali pemindaian berturut-turut pada basis kode yang sama **PASTI** menghasilkan byte output yang identik secara biner, independen terhadap urutan penjadwalan goroutine.

---

## 4. Anggaran Performa & Throughput (`QUAL-04-PERF-001`)

Anggaran performa Fase 4 dipecah per-komponen sebagai target desain terukur:

| Komponen Engine | Metrik Indikator | Target Desain | Klasifikasi |
| :--- | :--- | :---: | :--- |
| **Directory Walker** | Throughput traversal direktori | $\ge 20.000\text{ entries/s}$ | Performance Baseline |
| **Ignore Matcher** | Evaluasi 1.000 path sekuensial | $\le 1\text{ ms}$ total | Performance Baseline |
| **Analyzer Traversal** | Throughput evaluasi node AST in-memory | $\ge 50.000\text{ nodes/s/core}$ | Performance Target |
| **End-to-End Scan** | Throughput pemindaian korpus standar | Baseline terukur | Performance Budget |

---

## 5. Ambang Batas Kualitas Kode & Metrik Kelulusan

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Line Coverage `internal/config/...`** | $\ge 90\%$ | `go test -cover ./internal/config/...` |
| **Line Coverage `internal/scanner/...`** | $\ge 85\%$ | `go test -cover ./internal/scanner/...` |
| **Line Coverage `internal/analyzer/...`** | $\ge 90\%$ | `go test -cover ./internal/analyzer/...` |
| **Verifikasi Data Race** | $0$ data race detected | `go test -race ./internal/...` |
| **Batas Kompleksitas Siklomatik** | $\le 12$ per fungsi | `gocyclo -over 12 ./internal/scanner ./internal/analyzer` |
| **Linter Compliance** | $0$ issues | `golangci-lint run ./internal/...` |
