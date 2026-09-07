# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

> **Retention Policy:** Berkas ini hanya menyimpan rilis terbaru dan rilis sebelumnya ($N$ dan $N-1$). Riwayat rilis yang lebih lama diarsipkan di direktori [`docs/05-release/changelogs/`](docs/05-release/changelogs/).

---

## [v1.0.0-beta.3] - 2026-09-07

### Beta Evaluation Rationale & Honest Field-Testing Retrospective
* **Honest Error & Real-World Diagnostic Discovery:** Pengujian lapangan langsung pada `v1.0.0-beta.2` mengungkap tiga tantangan nyata dalam lingkungan sistem pengembang:
  1. **PATH Pollution & Collision:** Akumulasi entri path duplikat di shell environment pengembang (seperti `~/.local/bin` yang terduplikasi berulang kali di shell profile) serta potensi tabrakan multiple biner di direktori berbeda.
  2. **Installer Concurrency & State Guard:** Tidak adanya guard status "already installed" dan lockfile proteksi konkurensi pada script instalasi, yang dapat memicu race condition saat installer dijalankan di latar belakang bersamaan dengan proses uninstall.
  3. **In-Place Update Archive Unpacking & CLI Flag Shorthand:** Perintah pembaruan bawaan memerlukan kemampuan ekstraksi otomatis langsung dari arsip rilis (`.tar.gz` dan `.zip`) serta pengenalan alias flag (`-u`, `--update`, `-doctor`, `--doctor`) di root CLI.
* **Continued Multi-Project Calibration:** Pengujian performa deteksi dan kalibrasi rasio false-positive (FP) tetap dilanjutkan pada proyek-proyek eksternal sebelum deklarasi rilis stabil `v1.0.0`.

### Added
* **Diagnostic Doctor Command (`charites doctor`, `-doctor`, `--doctor`):**
  - Memeriksa kebersihan `$PATH`, mendeteksi entri direktori duplikat, dan mengaudit tabrakan biner (*multiple binary discovery* & *shadowing*).
  - Menguji izin penulisan direktori biner untuk menjamin kelancaran fitur pembaruan mandiri (*in-place self-update*).
  - Memverifikasi konektivitas jaringan ke GitHub Release API dengan batas waktu 3 detik.
  - Menghasilkan ringkasan rekomendasi remediasi yang dapat langsung dieksekusi pengguna.
* **Installer Concurrency Lock & Already-Installed Guard (`scripts/install.sh`):**
  - Lockfile berbasis PID (`/tmp/charites-installer.lock`) untuk mencegah race condition instalasi paralel.
  - Deteksi otomatis versi yang telah terpasang; jika versi yang sama telah terpasang, instalasi dilewati secara aman tanpa download ulang (dapat dipaksa via `--force` atau `CHARITES_FORCE=1`).
  - Pencegahan penambahan entri duplikat pada profil shell (`~/.bashrc`, `~/.zshrc`, `~/.profile`).
* **Archive Unpacking Engine for `charites update`:**
  - Parser internal berbasis `archive/tar`, `archive/zip`, dan `compress/gzip` untuk mengekstrak biner secara langsung dari arsip rilis resmi GitHub Releases.
  - Proteksi batas pembacaan (`io.LimitReader` 100MB) untuk mencegah ancaman eksploitasi *decompression bomb* (CWE-409 / gosec G110).
* **Root CLI Shorthand Aliases:**
  - Penambahan routing alias `-u`, `--update`, `-update` untuk pembaruan mandiri.
  - Penambahan routing alias `-doctor`, `--doctor` untuk diagnosis sistem.

### Changed
* Mengubah retensi berkas `CHANGELOG.md` utama agar fokus pada 2 versi teratas ($N$ dan $N-1$), dengan pengarsipan otomatis versi historis ke `docs/05-release/changelogs/`.

---

## [v1.0.0-beta.2] - 2026-09-07

### Beta Evaluation Rationale & Field-Testing Focus
* **Multi-Project Empirical Calibration:** Rilis ini dipertahankan pada status **Beta** karena masih memerlukan kontribusi dan pengujian empiris dalam performa deteksi serta kalibrasi false positive (FP) / true positive (TP) pada lebih dari 1 proyek nyata (lintas repositori produksi Astro, React TSX/JSX, dan Tailwind CSS).
* **Cadence & Spatial Analysis Field Validation:** Validasi model relasional spasial baru (`ux.spacing-rhythm-drift`) membutuhkan benchmarking pada berbagai ragam arsitektur layout komponen untuk memastikan akurasi deteksi ritme tanpa noise.
* **Community & Agent Feedback:** Kontribusi hasil pengujian riil, benchmark performa deteksi, dan pelaporan edge-case dikumpulkan melalui GitHub Issues dan tool MCP bawaan `charites_report_issue`.

### Update & Upgrade Commands
Untuk pengguna yang telah menginstal Charites, jalankan salah satu perintah pembaruan berikut:

#### In-Place Self-Update (Rilis Terpasang)
```bash
charites update
# atau alias:
charites --update
charites -u
```

#### Via Go Toolchain
```bash
go install github.com/will2469/charites/cmd/charites@v1.0.0-beta.2
```

#### Linux & macOS (One-Line Script)
```bash
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/scripts/install.sh | bash
```

#### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/will2469/charites/main/install.ps1 | iex
```

### Added
* **File-First Output Contract for JSON & Markdown Reporters:**
  - Mengubah hierarki pelaporan dari *rule-first* menjadi *file-first* (`files: [...]`), mengelompokkan temuan berdasarkan berkas target dengan koordinat `line:column`, jumlah pelanggaran, dan arahan supresi untuk mempermudah alur kerja perbaikan berkas-demi-berkas oleh pengembang maupun AI agent.
  - Mempertahankan backward compatibility skema JSON melalui root-level `diagnostics: [...]` yang terhubung secara kanonikal ke pelanggaran berkas.
  - Menambahkan laporan audit Markdown terkelompok per-berkas (`--format=markdown` / `-o report.md`) lengkap dengan scorecard ringkasan eksekutif, rincian kategori, dan tautan langsung ke dokumentasi wiki online.
* **New Quality & UX Invariant Rules:**
  - `ux.spacing-rhythm-drift`: Analisis ritme spasial relasional yang mengevaluasi variasi interval $N-1$ antar sibling dan observasi properti kontainer peer.
  - `ux.destructive-action-unconfirmed`: Analisis context-aware state-gating untuk aksi destruktif tanpa konfirmasi dialog/inline (Fixes #11).
  - `responsive.mobile-text-overflow`: Dukungan utilitas modern `wrap-break-word` dan penanganan overflow teks mobile (Fixes #10).
  - `responsive.fractional-width-gap-drift`: Deteksi drift layout akibat flexbox fractional width gap (Fixes #8).
  - `ux.multiline-input-misuse`: Pencegahan penggunaan input multiline pada kolom data teks baris tunggal (Fixes #7).
  - `semantic.consecutive-br-spacing`: Penegakan semantik struktur HTML terhadap penggunaan tag `<br>` berulang untuk manipulasi spacing (Fixes #6).
  - `ux.number-input-wheel-hazard`, `ux.number-input-identity-misuse`, `ux.number-input-missing-bounds`: Rangkaian aturan keandalan kontrol formulir numerik (Fixes #5).
  - Component-scoped style drift analyzer (Fixes #4).

### Changed
* **SSOT Canonical Diagnostic Ordering:** Menyatukan sorting kanonikal 7-dimensi langsung di lapisan `internal/ir` (`ir.SortDiagnostics`), mengeliminasi duplikasi logika sorting pada layer presenter reporter.
* **Linear File Grouping ($O(N)$):** Mengubah `GroupByFile()` menjadi partisi linier murni satu lintasan tanpa pemanggilan sort redundan.
* **Invariant I10 Preservation:** Menghapus `NormalizeSummary()` pada presenter agar metrik agregasi scanner (`ScanSummary`) tetap murni tanpa modifikasi di lapisan presentasi.
* **Deterministic String Path Arithmetic:** Menghilangkan syscall `os.Stat()` dan `filepath.Abs()` pada presenter, menjamin resolusi path relatif deterministik lintas platform tanpa ketergantungan disk I/O.
* **Clockless Markdown Reproducibility:** Menghilangkan dependensi clock `time.Now()` pada Markdown reporter dengan fallback epoch deterministik (`time.Unix(0, 0).UTC()`) untuk pengujian golden snapshot yang stabil dan *byte-for-byte reproducible*.

### Fixed
* **Peer Container Rhythm Fingerprint:** Penegakan *structural fingerprinting* dan observasi properti kontainer peer pada `ux.spacing-rhythm-drift`.
* **Carousel Control Slider Disambiguation:** Menghindari false positive pada slider kontrol di `cls.unconstrained-carousel` (Fixes #9).
* **Arbitrary Inline Scale Tokens:** Mendeteksi `scale-[...]` arbitrer dan menegakkan variabel token `global.css` pada rule `theme.*`.
* **Slash Opacity & Shadow Utility Closure:** Menutup celah bypass arbitrary slash opacity dan memperluas dukungan token `shadow-*` pada rule `theme.*`.

---

> **Historical Releases:** Changelogs for older versions (e.g. `v1.0.0-beta.1` and earlier) are archived under [`docs/05-release/changelogs/`](docs/05-release/changelogs/).
