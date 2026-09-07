# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

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
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/install.sh | bash
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

## [v1.0.0-beta.1] - 2026-09-06

### Beta Evaluation Rationale & Field-Testing Focus
* **Real-World Empirical Validation:** Over 50% of advanced rule combinations, nested component hierarchies, and bespoke Astro/React architectural patterns have only been tested internally. Testing on real-world projects is required to evaluate practical effectiveness, diagnostic precision, and false positive (FP) / false negative (FN) rates under production conditions.
* **Deferred Scopes:** Certain expansion domains-notably `seo.*` (`SPEC-EXP-12-SEO`)-have been intentionally deferred to avoid superficial overlap with generic linters and keep the compiler focused on design token integrity, accessibility, and Core Web Vitals.
* **Community & Agent Feedback:** Field feedback and real-world false positive reports are gathered via GitHub Issues and the built-in MCP two-phase HITL tool (`charites_report_issue`).

### Added
* **Ultra-Fast Zero-CGO Static Analysis Compiler:**
  - High-performance Go 1.26 AST parsing engine for `.astro`, `.tsx`, `.jsx`, and `.css` files with sub-millisecond per-file traversal without Node.js runtime or CGO overhead.
  - Unified Intermediate Representation (Leaf IR) streaming Astro and TSX/JSX ASTs into a normalized node graph.
  - SSOT Multi-Format Design Token Engine parsing `global.css`, `index.css`, `@theme`, and `tokens.json` (W3C DTCG format) with directed graph cycle detection (`ErrCycleDetected`) and recursion budget limits.
  - Evidence-based token verification ("The Banana Test") guaranteeing zero false positives for untokenized custom utility classes when no semantic token is declared.
* **90 Canonical Quality Rules across 8 Domains:**
  - `theme.*` (32 rules): Design token enforcement, slash opacity elimination (`bg-primary/10` $\rightarrow$ `bg-primary-light`), dark mode elevation preservation, CSS Cascade Layer boundaries.
  - `a11y.*` (16 rules): WCAG 2.2 AA standards, mathematical relative color contrast calculation, form control label bindings, modal keyboard trap prevention, iOS Safari auto-zoom prevention.
  - `responsive.*` (17 rules): Apple HIG/WCAG touch target ergonomics ($\ge 44 \times 44\text{px}$), container queries (`@container`), dynamic viewport units (`dvh`/`svh`), responsive table overflow wrapping, mobile keyboard safe areas.
  - `lcp.*`, `cls.*`, `inp.*`, `performance.*` (25 rules): Core Web Vitals optimization including hero image priority (`fetchpriority="high"`), preload links, explicit media aspect-ratio dimensions, layout shift prevention, long task unyielding detection, and Astro island deferred hydration.
* **Pure Stateless Model Context Protocol (MCP) Server (2026-07-28 Standard):**
  - `charites_scan`: Workspace component static analysis with rich structured diagnostics and online wiki links.
  - `charites_explain_rule`: Returns complete 8-Pillars architectural rationale, risk taxonomy, non-compliant examples, and remediation guidance.
  - `charites_list_rules`: Dynamic discovery of all 90 registered rules, domains, and default severities.
  - `charites_report_issue`: Two-Phase Human-in-the-Loop (HITL) reporting tool with SHA-256 draft signatures (Phase 1) and user-verified submission (Phase 2).
* **Rich Multi-Format CLI Reporters:**
  - Inline ANSI terminal reporter with colorized source snippets and actionable remediation hints.
  - Machine-readable JSON streaming format (`--format=json`) with rule metadata and online doc URLs.
  - Markdown audit reporter (`--format=markdown` / `-o report.md`) with executive scorecards, category violation breakdowns, and direct links to online wiki documentation.
* **1-SSOT Tri-Corpus Testing Harness:**
  - 17-pattern adversarial test matrix (P1-P5 positive, N1-N5 negative, A1-A7 adversarial) across all rules.
  - Continuous Go 1.26 fuzzing suite with 14,000+ synthetic mutations verifying zero crashes, memory leaks, or panic hazards.
* **Self-Management & Installation:**
  - Automated in-place self-update (`charites update`) and uninstaller (`charites uninstall`).
  - Cross-platform installation via `go install`, curl installer (`install.sh`), and PowerShell (`install.ps1`).

---
