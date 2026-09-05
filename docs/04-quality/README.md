# 04-QUALITY: Quality Assurance & Standards Plan

> **Dokumen Status:** Active / Draft
> **Standar Rujukan:** OpenSSF Best Practices Badge & Go Code Review Guidelines
> **Domain:** Standar Mutu Kode, Konfigurasi Linter, & Keamanan Rantai Pasok Charites

Dokumen ini mendefinisikan aturan baku penulisan kode sumber, gerbang inspeksi statis (*quality gating*), dan kebijakan keamanan perangkat lunak pada proyek **Charites**.

---

## 1. Pedoman Gaya Kode (Go 1.26 Idiomatic Guidelines)

Seluruh kontributor dan maintainer wajib mematuhi standar penulisan kode Go modern:

1. **Error Handling Terpadu:**
   - Wajib menggunakan wrapping error dengan `%w` (`fmt.Errorf("parse astro failed: %w", err)`).
   - Pemeriksaan error wajib menggunakan `errors.Is()` atau `errors.As()`, dilarang membandingkan string error secara mentah (`err.Error() == "..."`).
   - Setiap error harus ditangani tepat satu kali: catat ke log atau teruskan (*return*) ke pemanggil, jangan lakukan keduanya sekaligus (*the single handling rule*).
2. **Zero-Alloc Iterators (Go 1.26):**
   - Traversal tree atau struktur data koleksi wajib menggunakan pola `iter.Seq` atau `iter.Seq2`.
   - Hindari membuat slice penampung sementara hanya untuk melakukan loop (*avoid intermediate allocations*).
3. **Defensive Slices & Maps:**
   - Nilai slice atau map yang diterima dari luar atau diekspos keluar harus disalin (*defensive copy*) untuk mencegah *race condition* atau modifikasi tak sengaja.
4. **Panic Isolation:**
   - Package internal (parser, IR builder, analyzer) **DILARANG KERAS** membiarkan `panic` lolos ke pemanggil. Seluruh panic internal wajib diisolasi dengan `recover()` dan diubah menjadi `error` atau `Diagnostic` terstruktur.

---

## 2. Konfigurasi Static Analysis (`golangci-lint`)

Seluruh perubahan kode wajib lolos inspeksi `golangci-lint` tanpa peringatan (*zero warnings policy*).

### Linters Wajib Aktif (`.golangci.yml`):
- **`govet`**: Memeriksa keselarasan struct, shadow variables, dan format printf.
- **`staticcheck`**: Analisis statis canggih untuk mendeteksi bug fungsional dan kode mati.
- **`errcheck`**: Memastikan tidak ada nilai kembalian `error` yang diabaikan.
- **`revive`**: Pengganti modern untuk `golint` dengan aturan penamaan idiomatik.
- **`gocritic`**: Mendeteksi kode tidak efisien dan potensi bottleneck performa.
- **`gosec`**: Memeriksa kerentanan keamanan (injeksi path, integer overflow, dsb.).
- **`prealloc`**: Mengingatkan alokasi slice kapasitas awal saat ukuran sudah diketahui.

---

## 3. Tata Kelola Dependensi (Supply Chain Security)

Mengadopsi prinsip **OpenSSF Minimum Attack Surface**:
1. **Prioritaskan Standard Library:** Sebisa mungkin manfaatkan paket bawaan Go (`strings`, `bytes`, `iter`, `sync`, `path/filepath`, `encoding/json`).
2. **Larangan CGO Tidak Perlu:** Menjaga `CGO_ENABLED=0` untuk memastikan binary murni statis dan terbebas dari kerentanan memory corruption C library runtime.
3. **Audit Dependensi Eksternal:**
   - Setiap penambahan dependensi di `go.mod` harus melalui justifikasi teknis pada PR.
   - Rutin menjalankan `govulncheck ./...` di CI untuk mendeteksi Known Vulnerabilities (CVE) pada dependency tree.

---

## 4. Batasan Keamanan (Security Boundaries)

Sebagai linter yang memindai berkas di sistem berkas pengguna, Charites menerapkan perimeter pertahanan:

1. **Mitigasi Path Traversal:**
   - Scanner wajib membersihkan seluruh path relatif menggunakan `filepath.Clean()`.
   - Mencegah evaluasi berkas di luar *root target directory* saat ada symlink berbahaya (*symlink loop & directory escape protection*).
2. **Proteksi Zip Slip / Symlink DoS:**
   - Membatasi kedalaman traversal direktori (maksimum 64 tingkat) untuk mencegah recursive symlink denial-of-service.
3. **Ukuran Berkas Maksimum:**
   - Membatasi pembacaan berkas sumber tunggal maksimal 10 MB per berkas guna mencegah memory exhaustion akibat berkas non-kode (misal: binary atau video yang salah taruh di folder frontend).

---

## 5. Gerbang Kualitas Pull Request (Quality Gate Checklist)

Sebelum sebuah PR di-merge ke branch `main`, pipeline otomatis CI akan memverifikasi:
- [ ] `golangci-lint run` keluar dengan status exit `0` (clean).
- [ ] `govulncheck ./...` tidak menemukan kerentanan aktif.
- [ ] `go test -race ./...` lulus tanpa mendeteksi data race.
- [ ] Seluruh golden test cocok dan benchmark alokasi memori tidak menunjukkan regresi.
