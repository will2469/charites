# 04-QUALITY: 04 - Concurrency Safety, I/O Limits & Engine Performance Invariants

> **Kode Dokumen:** `QUAL-04-ENGINE`
> **Tahapan:** Fase 4 - Konfigurasi, Concurrency Scanner & Traversal Engine
> **Status:** Ready for Review
> **Standar Rujukan:** Go Concurrency Guidelines & Defensive Systems Programming

Dokumen ini mendefinisikan batasan kualitas, ketahanan konkurensi antar-goroutine, pencegahan kebocoran goroutine (*goroutine leak*), serta batas aman I/O disk pada paket konfigurasi, pemindai direktori, dan mesin traversal AST.

---

## 1. Invarian Keamanan Konkurensi & Pencegahan Deadlock

Operasi pemindaian paralel pada `internal/scanner` dan `internal/analyzer` wajib mematuhi panduan konkurensi berikut:

1. **Zero Data Race Invariant:**
   Seluruh interaksi antrean kerja, pertukaran data hasil pemindaian, dan resolusi status rule wajib bebas dari *data race*. Pengujian `go test -race` wajib lulus 100%.
2. **Share-Nothing Worker Architecture:**
   Setiap goroutine worker mengeksekusi satu berkas dalam konteks terisolasi (`*analyzer.Context`). Worker dilarang memodifikasi state bersama tanpa sinkronisasi eksplisit.
3. **Pencegahan Kebocoran Goroutine (Goroutine Leak Prevention):**
   - Semua channel pekerjaan (`jobs`) dan channel hasil (`results`) memiliki kepemilikan produsen (*producer-owned*) yang jelas.
   - Penutupan channel hanya dilakukan oleh produsen setelah seluruh `sync.WaitGroup` selesai (`wg.Wait()`).
   - Seluruh siklus kerja worker pool dapat dihentikan lebih awal (*early termination*) menggunakan `context.Context` untuk menangani sinyal interupsi terminal (`SIGINT`/`SIGTERM`).

---

## 2. Invarian Batasan I/O Disk & Proteksi DoS

Untuk melindungi sistem host dari kehabisan memori atau saturasi *file descriptor*:

1. **Strict Early Directory Pruning:**
   Direktori yang terdeteksi diabaikan oleh `.charitesignore` atau builtin ignore (seperti `node_modules/`, `.git/`, `.next/`) **DILARANG KERAS** dibuka dengan `os.ReadDir`. Ini memotong hingga 95% panggilan sistem (*system calls*) I/O disk pada proyek berskala besar.
2. **Batas Maksimum Ukuran Berkas (Max File Size Limit):**
   Charites membatasi ukuran berkas frontend yang dipindai maksimal **10 Megabytes**. Berkas yang melebihi ambang batas ini dianggap sebagai aset biner atau bundle minifikasi otomatis, dan dilewati dengan catatan peringatan agar tidak memicu Out-of-Memory (OOM).
3. **Pencegahan Symlink Loops:**
   Walker melacak inode direktori nyata untuk mendeteksi siklus symlink (*circular symbolic links*), mencegah loop tak terbatas (*infinite loop*).

---

## 3. Invarian Output Deterministik

Karena urutan penyelesaian berkas di worker pool bersifat non-deterministik (bergantung pada scheduler thread Go):
- Hasil temuan dari seluruh worker wajib dikumpulkan dan diurutkan secara deterministik menggunakan kriteria stabil: `(File, Line, Column, RuleID)`.
- Dua kali pemindaian berturut-turut pada repositori yang sama **PASTI** menghasilkan urutan diagnosis yang identik hingga ke tingkat byte.

---

## 4. Ambang Batas Kualitas & Metrik Kelulusan

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Line Coverage Paket Config** | $\ge 90\%$ | `go test -cover ./internal/config/...` |
| **Line Coverage Paket Scanner** | $\ge 85\%$ | `go test -cover ./internal/scanner/...` |
| **Line Coverage Paket Analyzer** | $\ge 90\%$ | `go test -cover ./internal/analyzer/...` |
| **Verifikasi Data Race** | $0$ race condition | `go test -race ./internal/...` |
| **Batas Kompleksitas Siklomatik** | $\le 12$ per fungsi | `gocyclo -over 12 ./internal/scanner ./internal/analyzer` |
| **Linter Compliance** | $0$ issues | `golangci-lint run ./internal/...` |
