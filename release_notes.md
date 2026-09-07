# Release Notes - Charites v1.0.0-beta.2 (2026-09-07)

Welcome to **Charites v1.0.0-beta.2**, the second beta release of the compile-time static analyzer and design token linter for **Astro**, **React TSX/JSX**, and **Tailwind CSS**.

This release introduces the **File-First Report Hierarchy** for JSON and Markdown outputs, the **Spacing Rhythm Drift** cadence analyzer (`ux.spacing-rhythm-drift`), several new UX and form safety rules, and strict SSOT canonical ordering.

---

## 1. Beta Evaluation Rationale & Multi-Project Calibration

> [!IMPORTANT]
> **Mengapa Rilis Ini Masih Berstatus Beta (`v1.0.0-beta.2`)?**
> * **Empirical Multi-Project Calibration:** Evaluasi performa deteksi, akurasi parsing AST, serta rasio *false positive* (FP) / *true positive* (TP) masih memerlukan kontribusi dan pengujian riil pada lebih dari 1 proyek eksternal. Kalibrasi lintas berbagai monorepo produksi (Astro, Next.js/React, Tailwind v3/v4) sangat esensial sebelum menetapkan rilis stabil `v1.0.0`.
> * **Spatial Rhythm Drift Fine-Tuning:** Rule baru `ux.spacing-rhythm-drift` memperkenalkan model relasional spasial ($N-1$ interval sibling dan observasi properti peer container). Model ini membutuhkan benchmarking performa dan variasi cadence layout dari komunitas developer.
> * **Community Feedback & Edge Cases:** Kami mengundang tim pengembang untuk menguji rilis ini pada proyek aktif mereka dan melaporkan temuan, performa deteksi, atau false positive melalui GitHub Issues atau MCP tool bawaan `charites_report_issue`.

---

## 2. Update & Upgrade Commands

Bagi pengguna yang sudah memasang Charites versi sebelumnya (`v1.0.0-beta.1` atau rilis development), gunakan perintah berikut untuk memperbarui:

### A. In-Place Self-Update (Rekomendasi Utama)
Charites dilengkapi mekanisme pembaruan otomatis di tempat tanpa perlu download manual:
```bash
charites update
```
*Alias:*
```bash
charites --update
charites -u
```

### B. Via Go Toolchain
```bash
go install github.com/will2469/charites/cmd/charites@v1.0.0-beta.2
```

### C. Linux & macOS (One-Line Updater / Installer)
```bash
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/install.sh | bash
```

### D. Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/will2469/charites/main/install.ps1 | iex
```

### E. Verifikasi Versi Terpasang
```bash
charites --version
```
*Output yang diharapkan:*
```text
charites version 1.0.0-beta.2 (go1.26.x)
```

---

## 3. Fresh Installation (Pemasangan Baru)

Untuk pemasangan baru pada mesin atau pipeline CI/CD:

```bash
# Go Developers:
go install github.com/will2469/charites/cmd/charites@v1.0.0-beta.2

# Linux / macOS:
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/install.sh | bash

# Windows PowerShell:
irm https://raw.githubusercontent.com/will2469/charites/main/install.ps1 | iex
```

---

## 4. Rangkuman Pembaruan & Fitur Baru (What's New)

###  File-First Output Contract (JSON & Markdown)
* **Remediation-First Hierarchy:** Mengubah hierarki pelaporan dari *rule-first* menjadi *file-first* (`files: [...]`). Setiap berkas merangkum jumlah error/warning serta daftar temuan dengan koordinat tepat `line:column` dan arahan supresi (`suppression`). Hal ini mempercepat alur kerja remediasi baik bagi pengembang maupun AI coding agent.
* **Backward Compatibility:** Root-level `diagnostics: [...]` tetap dipertahankan dengan jaminan invarian identik dengan gabungan pelanggaran berkas.
* **File-Grouped Markdown Audit:** Format `--format=markdown` / `-o audit-report.md` kini menyajikan kartu skor ringkasan eksekutif, tabel pelanggaran terkelompok per-berkas, dan tautan langsung ke wiki dokumentasi kanonikal.

###  Aturan Baru (New Quality & UX Rules)
1. **`ux.spacing-rhythm-drift`:** Mendeteksi inkonsistensi irama spasial (*spacing cadence*) menggunakan model relasional interval sibling $N-1$ serta observasi properti kontainer peer.
2. **`ux.destructive-action-unconfirmed` (Fixes #11):** Analisis context-aware state-gating yang mewajibkan dialog atau konfirmasi inline pada tombol aksi destruktif.
3. **`responsive.mobile-text-overflow` (Fixes #10):** Dukungan penuh utilitas `wrap-break-word` dan overflow handling untuk viewport mobile.
4. **`responsive.fractional-width-gap-drift` (Fixes #8):** Deteksi layout drift pada flexbox akibat pembagian lebar fraksional yang bertabrakan dengan gap.
5. **`ux.multiline-input-misuse` (Fixes #7):** Deteksi pencegahan penggunaan `<textarea>` atau multiline input pada field data baris tunggal.
6. **`semantic.consecutive-br-spacing` (Fixes #6):** Deteksi penggunaan tag `<br>` ganda berturut-turut untuk manipulasi jarak tata letak visual.
7. **`ux.number-input-wheel-hazard`, `ux.number-input-identity-misuse`, `ux.number-input-missing-bounds` (Fixes #5):** Rangkaian audit keamanan dan keandalan input formulir bertipe number.
8. **Component-Scoped Style Drift Analyzer (Fixes #4):** Deteksi deviasi token dan style drift antar instans komponen serupa.

###  Presisi Diagnostik, Performa & Determinisme
* **SSOT Canonical Ordering:** Pengurutan diagnostik 7-dimensi kanonikal kini berada di `internal/ir` (`ir.SortDiagnostics`), mengeliminasi multi-pass sorting di lapisan presenter.
* **Linear File Grouping ($O(N)$):** Partisi linier murni untuk pengelompokan berkas tanpa alokasi memori berlebih.
* **Invariant I10 Guarantee:** Metrik `ScanSummary` dijamin murni tanpa manipulasi di layer pelaporan.
* **Deterministic & Clockless:** Resolusi path murni berbasis string arithmetic (tanpa syscall I/O `os.Stat`) dan timestamp markdown deterministik untuk hasil yang 100% *byte-for-byte reproducible*.
* **Token Hardening:** Penutupan celah arbitrer pada `scale-[...]`, slash opacity (`bg-primary/10`), dan ekspansi token utilitas `shadow-*`.

---

## 5. Panduan Penggunaan Singkat

```bash
# Pemindaian standar direktori saat ini
charites scan .

# Pemindaian dengan laporan JSON terstruktur (file-first)
charites scan -f json .

# Pemindaian dan pembuatan laporan audit Markdown lengkap
charites scan -f markdown -o audit-report.md .

# Menjalankan server Model Context Protocol (MCP 2026-07-28)
charites mcp
```

---

## 6. Changelog Lengkap

Rincian commit lengkap dan perbandingan diff dapat dilihat di [CHANGELOG.md](CHANGELOG.md) serta [GitHub Releases](https://github.com/will2469/charites/releases/tag/v1.0.0-beta.2).
