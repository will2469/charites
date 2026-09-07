# Release Notes - Charites v1.0.0-beta.3 (2026-09-07)

Welcome to **Charites v1.0.0-beta.3**, the third beta release of the compile-time static analyzer and design token linter for **Astro**, **React TSX/JSX**, and **Tailwind CSS**.

This release introduces the **System & Environment Doctor (`charites doctor`)**, installs **robust concurrency and version guards** in `scripts/install.sh`, adds native **archive unpacking (.tar.gz and .zip)** for `charites update`, and documents real-world environment observations ("honest errors") surfaced during beta testing.

---

## 1. Beta Evaluation Rationale & Honest Field-Testing Retrospective

> [!NOTE]
> **Mengapa Rilis Ini Dinaikkan ke `v1.0.0-beta.3`?**
> Pengujian lapangan riil pada `v1.0.0-beta.2` mengungkap beberapa kondisi lingkungan pengembang nyata yang memerlukan mitigasi sistem tingkat pertama:
>
> 1. **PATH Pollution & Collision:** Lingkungan shell sering kali memuat duplikasi entri direktori di variabel `$PATH` (misalnya `~/.local/bin` yang terpanggil belasan kali dari `.bashrc` majemuk) atau adanya beberapa salinan biner `charites` di direktori berbeda (`~/.local/bin` vs `/usr/local/bin` vs `~/go/bin`) yang berpotensi membingungkan versi mana yang aktif.
> 2. **Installer Race Condition & State Guard:** Script installer lama tidak memiliki pengecekan status "already installed" dan proteksi kunci konkurensi (*lockfile*), sehingga jika proses download/install berjalan di latar belakang bersamaan dengan perintah pencopotan (*uninstall*), file yang baru dicopot dapat langsung tercipta kembali tanpa disengaja.
> 3. **In-Place Update Archive Unpacking:** Fitur pembaruan mandiri (`charites update`) kini dilengkapi mesin ekstraksi bawaan untuk mengekstrak biner langsung dari paket arsip resmi GitHub (`.tar.gz` dan `.zip`) dengan pertahanan terhadap *decompression bomb* (`io.LimitReader`).
> 4. **Root CLI Shorthand Aliases:** Mendukung `-u`, `--update`, serta `-doctor`, `--doctor` langsung dari level root CLI untuk mempermudah eksekusi tanpa keharusan mengingat struktur subcommand.

---

## 2. Fitur Baru Utama (Highlights)

###  System & Environment Diagnostic Tool (`charites doctor`)
Perintah baru untuk memvalidasi kesehatan instalasi dan lingkungan shell pengembang:
```bash
charites doctor
# atau via shortcut:
charites -doctor
charites --doctor
```
Fitur audit yang dilakukan:
* **$PATH Cleanliness:** Mendeteksi entri direktori duplikat di `$PATH` dan memberikan rekomendasi pembersihan.
* **Binary Discovery & Collision:** Menemukan semua lokasi biner `charites` di `$PATH`, menandai biner aktif, dan memperingatkan jika ada biner sekunder yang tertimpa (*shadowed*).
* **Write Permissions:** Menguji apakah direktori biner dapat ditulis secara langsung untuk operasi `charites update` tanpa sudo.
* **API Reachability:** Memvalidasi konektivitas ke GitHub Releases API (batas waktu 3 detik).
* **Workspace Health:** Mendeteksi konfigurasi `charites.yaml` lokal.

###  Installer Guards (`scripts/install.sh`)
* **Already-Installed Guard:** Jika versi yang sama sudah terpasang, installer akan memberi tahu dan keluar secara aman tanpa men-download ulang (dapat dipaksa dengan opsi `--force` atau `CHARITES_FORCE=1`).
* **Concurrency Lock Protection:** Menggunakan lockfile `/tmp/charites-installer.lock` berbasis PID untuk mencegah eksekusi installer ganda secara simultan.
* **Pencegahan Duplikasi Shell Profile:** Menghindari penambahan baris `export PATH` berulang ke `~/.bashrc`, `~/.zshrc`, atau `~/.profile`.

###  Archive Unpacking Engine pada `charites update`
* Mendukung pembaruan mandiri (*in-place self-update*) langsung dari arsip `.tar.gz` (Linux/macOS) dan `.zip` (Windows) tanpa memerlukan perkakas dekompresi eksternal.
* Menerapkan batas ukuran ekstraksi 100MB (`io.LimitReader`) untuk mencegah kerentanan dekompresi bom (CWE-409 / gosec G110).

---

## 3. Update & Upgrade Commands

Bagi pengguna yang sudah memasang Charites versi sebelumnya:

### A. In-Place Self-Update (Rekomendasi Utama)
```bash
charites update
# atau alias:
charites --update
charites -u
```

### B. Via Go Toolchain
```bash
go install github.com/will2469/charites/cmd/charites@v1.0.0-beta.3
```

### C. Linux & macOS (One-Line Updater / Installer)
```bash
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/scripts/install.sh | bash
```

### D. Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/will2469/charites/main/install.ps1 | iex
```

### E. Verifikasi Versi & Diagnosis Sistem
```bash
charites --version
charites doctor
```

---

## 4. Fresh Installation (Pemasangan Baru)

```bash
# Go Developers:
go install github.com/will2469/charites/cmd/charites@v1.0.0-beta.3

# Linux & macOS:
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/scripts/install.sh | bash

# Windows PowerShell:
irm https://raw.githubusercontent.com/will2469/charites/main/install.ps1 | iex
```

---

## 5. Changelog & Riwayat Rilis

* Rincian perubahan lengkap dapat dilihat di [CHANGELOG.md](CHANGELOG.md).
* Arsip versi historis (seperti `v1.0.0-beta.1`) tersimpan di direktori [docs/05-release/changelogs/](docs/05-release/changelogs/).
* Unduhan biner dan checksums SHA-256 tersedia langsung di [GitHub Releases](https://github.com/will2469/charites/releases/tag/v1.0.0-beta.3).
