# 06-ROADMAP: Project Charter, Phasing & Risk Governance

> **Dokumen Status:** Active / Draft
> **Standar Rujukan:** CNCF Project Charter & Linux Foundation OSPO Framework
> **Domain:** Peta Jalan Proyek, Manajemen Fase Rilis, & Matriks Mitigasi Risiko Teknis

Dokumen ini mendefinisikan visi strategis, pentahapan rilis (*phased implementation*), dan manajemen risiko teknis untuk proyek **Charites**.

---

## 1. Visi & Misi Proyek

- **Visi:** Menjadikan Charites sebagai standar static analyzer tercepat, paling ringan, dan paling andal di ekosistem frontend modern berbasis Astro, TypeScript/JSX, dan Tailwind CSS.
- **Misi:** Menghilangkan kelemahan linter warisan berbasis Node.js yang lambat, boros memori, dan bergantung pada `node_modules` raksasa, dengan menghadirkan single binary Go native yang mengeksekusi analisis dalam hitungan milidetik.

---

## 2. Tahapan Rilis (Phased Implementation Roadmap)

```mermaid
timeline
    title Peta Jalan Rilis Charites
    section Phase 1 : 2026 Q3
      Vertical Slice MVP : go.mod 1.26 : IR Builder : Rule #1 hardcode_opacity_color : Golden Tests
    section Phase 2 : 2026 Q4
      Core Rules Expansion : A11y Rules : Web Vitals LCP/CLS/INP : React Hooks Rules : SEO Rules
    section Phase 3 : 2027 Q1
      AI & Knowledge Ecosystem : MCP Server Stdio : Subcommand charites wiki : Auto-docs wiki/*.md
    section Phase 4 : 2027 Q2
      IDE Integration : Language Server Protocol (LSP) : VS Code Extension : Real-time Diagnostics
```

### Rincian Tiap Fase:

### Phase 1: Vertical Slice MVP (Fokus Saat Ini)
- **Target:** Membuktikan stabilitas pipa compiler dari ujung ke ujung dengan 1 rule minimal.
- **Deliverables:**
  - Setup fondasi Go 1.26 (`cmd/charites/main.go`).
  - Parser token Tailwind CSS v4 (`@theme`), parser Astro frontmatter, dan parser TSX.
  - Normalisasi ke kontrak data `ir.Node`.
  - Implementasi Rule #1: `hardcode_opacity_color.go` (Semgrep ID: `theme.hardcode-opacity-color`).
  - Traversal engine dengan iterator Go 1.26 dan konkurensi worker pool.
  - Reporter terminal ANSI dan JSON.
  - Verifikasi golden snapshot test pass 100%.

### Phase 2: Core Rules Expansion (Porting dari Legacy)
- **Target:** Mem-porting seluruh 30+ rules legacy ke dalam arsitektur Go yang baru.
- **Deliverables:**
  - Subpackage A11y (`missing_htmlfor`, `button_accessible_name`, `form_controls`).
  - Subpackage Web Vitals (`lcp_image_priority`, `cls_dimensions`, `inp_listeners`).
  - Subpackage Hooks (`exhaustive_deps`, `hooks_rules`).
  - Subpackage SEO (`meta_tags`, `canonical_url`).

### Phase 3: AI & Knowledge Ecosystem
- **Target:** Menjadikan Charites ramah agen kecerdasan buatan (*AI pair-programming*).
- **Deliverables:**
  - Server MCP Stdio (`charites mcp`) mematuhi spek JSON-RPC 2.0 (2026-07-28).
  - Generator ensiklopedia rule otomatis (`charites wiki`) yang mengekspor metadata rule ke direktori `wiki/` (seperti format dokumentasi Argus).

### Phase 4: IDE Integration & Language Server Protocol (LSP)
- **Target:** Diagnostik real-time saat developer mengetik di text editor.
- **Deliverables:**
  - Subcommand `charites lsp` yang mendukung protokol LSP resmi.
  - Ekstensi resmi VS Code yang menampilkan squiggly lines langsung pada editor.

---

## 3. Matriks Manajemen Risiko Teknis (Failure Modes & Effects Analysis)

| Risiko Teknis | Tingkat Dampak | Probabilitas | Rencana Mitigasi Arsitektural |
| :--- | :--- | :--- | :--- |
| **Kompleksitas Parsing TSX & Astro di Go** (karena tidak ada standard parser bawaan Go) | **TINGGI** | Sedang | Menggunakan pemisahan dua tahap: Split frontmatter `---` secara streaming, lalu isolasi template HTML/JSX ke token parser ringan tanpa memaksakan full TypeScript compiler complexity. |
| **Monorepo File Walker I/O Bottleneck** (proyek dengan puluhan ribu file) | Sedang | Rendah | Menggunakan fast dirwalker paralel yang membaca direktori secara chunk, mematuhi `.ignore` sebelum membaca konten file, dan memanfaatkan thread pool Go. |
| **Memory Exhaustion pada Berkas Raksasa** | Sedang | Rendah | Menerapkan hard-limit ukuran berkas (maksimal 10 MB per berkas) dan daur ulang alokasi memory menggunakan `sync.Pool`. |
| **Kelelahan False Positive (False-Positive Fatigue)** | **TINGGI** | Sedang | Menyediakan inline ignore directive: `<!-- charites-ignore -->` (Astro) atau `{/* charites-ignore */}` (JSX), serta memisahkan severity `WARN` vs `ERROR`. |

---

## 4. Tata Kelola Kontribusi & Proposal Rule Baru (RFC Workflow)

Untuk menjaga kualitas dan integritas repositori:
1. **Proposal Rule Baru Wajib via RFC:**
   - Sebelum membuat rule baru, kontributor wajib mengajukan proposal spesifikasi di `docs/01-spec/` atau issue GitHub.
   - Proposal harus mencakup: nama rule, kategori, severity, contoh kode salah (*bad code*), contoh perbaikan (*good code*), dan penjelasan mengapa pelanggaran tersebut berdampak negatif bagi performa/UX.
2. **Review Minimal 1 Maintainer:**
   - Setiap PR wajib disetujui minimal 1 core maintainer dan lolos seluruh rangkaian CI/CD automated gates.
