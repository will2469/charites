---
name: charites-wikis-builder
description: "Standard operational procedure (SOP) and automated transformation engine to generate clean, authoritative, public-facing open-source Rule Wikis (wiki/<category>.<slug>.md) and Home.md index conforming to Charites 8-Pillars documentation standards."
compatibility: "Markdown, GitHub Wiki, W3C/WCAG documentation standards"
metadata:
  version: "1.0.0"
  author: "Charites Team / Will"
  license: "MIT"
  citations:
    - "https://github.com/will2469/charites/wiki"
    - "W3C CSS Color Module Level 4"
    - "W3C Web Content Accessibility Guidelines (WCAG 2.2)"
    - "Core Web Vitals Performance Standards"
---

# Charites Public Wiki Builder Skill (`charites-wikis-builder`)

> **Mandat Utama:** Mentransformasikan spesifikasi aturan static analysis menjadi **Dokumentasi Wiki Publik Resmi** bertaraf industri di [`wiki/<category>.<slug>.md`](wiki/) dengan mematuhi **8-Pillars Documentation Matrix** dan memperbarui katalog [`wiki/Home.md`](wiki/Home.md).

---

## 1. Perimeter Dokumentasi Bersih (Clean Documentation Boundary)

Dokumentasi publik Charites (`github.com/will2469/charites`) wajib bersifat mandiri, modular, dan agnostik:

| Kategori | Format Dilarang | Format Standar Publik Terbuka |
| :--- | :--- | :--- |
| **Konteks Komponen** | Nama komponen privat aplikasi klien | Komponen generik standar: `Button`, `Card`, `Modal`, `Navbar`, `Header` |
| **Paket / Path Proyek** | Path direktori privat di luar repositori | Path standar web: `src/components/`, `src/layouts/`, `src/styles/` |
| **Tautan Referensi** | Tautan file lokal privat | Tautan standar publik: W3C, MDN, WCAG 2.2, web.dev, Tailwind CSS docs |
| **Identifier Aturan** | Kode numerik arbitrary | Semgrep-style canonical ID: `<category>.<slug>` (misal `theme.hardcode-color`) |

---

## 2. Struktur Standar Dokumen Wiki (8-Pillars Matrix)

Setiap berkas wiki di `wiki/<category>.<slug>.md` **wajib** mengikuti struktur 8 pilar berikut:

```markdown
# [category].[slug]: [Judul Aturan dalam Title Case]

> **Rule ID:** `<category>.<slug>`
> **Severity:** `CRITICAL` | `HIGH` | `MEDIUM` | `LOW`
> **Category:** `theme` | `a11y` | `responsive` | `perf` | `tailwind`
> **Target Standards:** [W3C Design Tokens, WCAG 2.2, Core Web Vitals, OKLCH]

---

## 1. Overview & Core Invariant
[Penjelasan singkat 1-2 paragraf mengenai prinsip dasar aturan dan apa yang ditegakkannya]

---

## 2. Technical Grounding & Engine Realities
[Penjelasan mendalam mengenai dampak runtime web: CSS cascade, layout recalculation, contrast ratio, atau touch target sizing]

---

## 3. Static Analysis Architecture & AST Detection
[Bagaimana Charites mendeteksi pelanggaran pada Leaf IR: Node types, attribute inspection, value normalization]

---

## 4. Vulnerability & Risk Taxonomy
[Tabel risiko jika aturan dilanggar: aesthetic fatigue, WCAG non-compliance, layout shift]

---

## 5. Non-Compliant Code Patterns (Bad Examples)
[Contoh kode melanggar pada Astro, TSX, atau CSS]

---

## 6. Compliant Implementation Patterns (Good Examples)
[Contoh kode yang diremediasi menggunakan design token atau best practice]

---

## 7. How to Suppress (Ignore Directives)
[Panduan sintaks inline dan block suppression: `<!-- charites:ignore -->`, `/* charites:ignore */`]

---

## 8. Configuration Reference (`charites.yaml`)
[Opsi konfigurasi yang relevan untuk aturan ini pada file charites.yaml]
```

---

## 3. Format Index Wiki (`wiki/Home.md`)

Setiap penambahan Rule Wiki baru wajib dicatatkan pada `wiki/Home.md` dengan format tabel kategori:

| Rule ID | Category | Severity | Description | Status |
| :--- | :--- | :--- | :--- | :--- |
| [`theme.hardcode-color`](theme.hardcode-color.md) | `theme` | `HIGH` | Detects raw un-tokenized hex/rgb color literals | `enabled` |
| [`a11y.alt-text`](a11y.alt-text.md) | `a11y` | `CRITICAL` | Requires descriptive alt attributes on images | `enabled` |

---

## 4. Standar Nada Tulisan (Voice & Tone)

- **Otoritatif & Berbasis Standar:** Gunakan terminologi resmi web platform (Fitts's Law, WCAG Contrast Minimum, Composite Layers, Cumulative Layout Shift).
- **Edukatif & Siap Pakai:** Setiap contoh bad pattern wajib dipasangkan dengan good pattern yang siap pakai (*copy-pasteable*).
- **Ringkas & To-the-Point:** Hindari cerita fiktif atau narasi berbunga-bunga; utamakan kejelasan teknis bagi software engineer.
