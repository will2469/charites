---
name: charites-versioning
description: "MANDATORY SEMVER 2.0.0 RELEASE & DIFF GUARDIAN FOR CHARITES: Deterministic version increment evaluator and breaking change inspector. Auto-triggers when evaluating version bumps, checking git diff from latest tag to HEAD, auditing API backward compatibility, reviewing breaking changes in CLI/rules/config/MCP, or preparing release tags for Charites ('charites versioning', 'semver', 'version bump', 'major minor patch', 'tag diff', 'breaking change check', 'release check')."
compatibility: "Go 1.26+, git, bash"
metadata:
  version: "1.0.0"
  author: "Charites Team / Will"
  license: "MIT"
  citations:
    - "Semantic Versioning 2.0.0 Specification (https://semver.org/spec/v2.0.0.html) by Tom Preston-Werner"
    - "Conventional Commits 1.0.0 (https://www.conventionalcommits.org/)"
    - "IETF RFC 2119: Key words for use in RFCs to Indicate Requirement Levels"
---

# Charites Semantic Versioning & Diff Guardian (`charites-versioning`)

> **Core Thesis:** Penentuan nomor rilis pada **Charites** (`github.com/will2469/charites`) tunduk sepenuhnya pada **Semantic Versioning 2.0.0 (SemVer)**. Evaluasi kenaikan versi (`MAJOR.MINOR.PATCH`) wajib dihitung secara deterministik dari perubahan nyata (*ground-truth diff*) antara **tag rilis terakhir (`LATEST_TAG`)** dan **state saat ini (`HEAD`)**, menjaga integritas kontrak API publik tanpa asumsi spekulatif.

---

## 1. Fondasi Spesifikasi Resmi SemVer 2.0.0

Format versi stabil Charites adalah triplet numerik **`MAJOR.MINOR.PATCH`** (contoh: `1.0.0`, `1.1.2`):

```text
      X   .   Y   .   Z
    MAJOR   MINOR   PATCH
      │       │       │
      │       │       └─► Clause 5: Backward-compatible bug/security fix
      │       └─────────► Clause 6: Backward-compatible new features / deprecations
      └─────────────────► Clause 7: Incompatible API breaking changes
```

### Tabel Sitasi Klausul Otoritatif (semver.org)

| Komponen | Klausul | Kutipan Dokumen Resmi Spesifikasi SemVer 2.0.0 | Dampak Penomoran |
| :--- | :--- | :--- | :--- |
| **MAJOR (X)** | **Clause 7** | *"Major version X (X.y.z \| X > 0) MUST be incremented if any backwards incompatible changes are introduced to the public API. It MAY also include minor and patch level changes. Patch and minor version MUST be reset to 0 when major version is incremented."* | `1.1.2` → **`2.0.0`** *(Reset Minor & Patch ke 0)* |
| **MINOR (Y)** | **Clause 6** | *"Minor version Y (x.Y.z \| x > 0) MUST be incremented if new, backwards compatible functionality is introduced to the public API. It MUST be incremented if any public API functionality is marked as deprecated... Patch version MUST be reset to 0 when minor version is incremented."* | `1.1.2` → **`1.2.0`** *(Reset Patch ke 0)* |
| **PATCH (Z)** | **Clause 5** | *"Patch version Z (x.y.Z \| x > 0) MUST be incremented if only backwards compatible bug fixes are introduced. A bug fix is defined as an internal change that fixes incorrect behavior."* | `1.1.2` → **`1.1.3`** *(Hanya Patch naik)* |

---

## 2. Batas Kontrak API Publik Charites (*Public API Boundary*)

Berdasarkan **Clause 4**, penentuan breaking change di Charites dievaluasi terhadap **5 Permukaan Kontrak Publik**:

1. **CLI Commands & Flags (`cmd/charites/` & `internal/cli/`):**
   Perintah `charites scan`, alias `check`/`run`, subcommand `wiki`, `mcp`, `update` (alias `-u` / `--update`), flags (`-f/--format`, `--ext`, `--category`, `--rule`, `--config`, `--ignore`, `--no-color`, `--fail-on-warn`), kontrak exit code (`0`, `1`, `2`), dan skema JSON stream output.
2. **Skema Konfigurasi & Ignore (`charites.yaml` & `.charitesignore`):**
   Struktur blok `rules:`, `ignore:`, penanganan sintaks glob, serta invarian *Default: YES*.
3. **Aturan Analyzer (`internal/rules/`):**
   Semgrep ID aturan (`<category>.<slug>`, misal `theme.hardcode-opacity-color`), tipe keparahan default (`DefaultSeverity`), dan skema payload `ir.Diagnostic`.
4. **Server Model Context Protocol (`internal/mcp/`):**
   Nama tool resmi (`charites_scan`, `charites_explain_rule`, `charites_list_rules`), parameter input JSON schema, dan struktur response content block.
5. **Kontrak Data IR Publik (`internal/ir/`):**
   Struktur `*ir.Node`, iterator `Walk() iter.Seq[*ir.Node]`, struct `Diagnostic`, dan enum `Severity`.

---

## 3. Pohon Keputusan Evaluasi Diff (*Decision Tree*)

```text
              [ EVALUASI GIT DIFF: LATEST_TAG...HEAD ]
                                 │
     Apakah terdapat perubahan API Publik yang TIDAK KOMPATIBEL?
     (Hapus flag/subcommand CLI, ubah skema charites.yaml, ubah nama tool MCP,
      ubah Semgrep ID rule, ubah kontrak publik ir.Node, 'feat!:' / 'BREAKING CHANGE:')
                                 │
                 ┌───────────────┴───────────────┐
                YA                              TIDAK
                 │                               │
                 ▼                               ▼
         ┌───────────────┐        Apakah terdapat FITUR BARU kompatibel
         │  BUMP MAJOR   │        atau DEPREKASI baru?
         │ 1.1.2 → 2.0.0 │        (Tambah rule baru, tambah flag opsional CLI,
         └───────────────┘        tambah tool MCP, '// Deprecated:', commit 'feat:')
                                                 │
                                 ┌───────────────┴───────────────┐
                                YA                              TIDAK
                                 │                               │
                                 ▼                               ▼
                         ┌───────────────┐        Apakah terdapat BUG FIX,
                         │  BUMP MINOR   │        FP/FN FIX, atau OPTIMASI?
                         │ 1.1.2 → 1.2.0 │        (Fix AST tokenizer, tambah Tri-Corpus,
                         └───────────────┘        optimasi memory pool, doc/test)
                                                                 │
                                                 ┌───────────────┴───────────────┐
                                                YA                              TIDAK
                                                 │                               │
                                                 ▼                               ▼
                                         ┌───────────────┐               ┌───────────────┐
                                         │  BUMP PATCH   │               │    NO BUMP    │
                                         │ 1.1.2 → 1.1.3 │               │   (Unchanged) │
                                         └───────────────┘               └───────────────┘
```

---

## 4. Kriteria Penentuan Berdasarkan Diff Tag Terakhir

### A. Kriteria MAJOR (`1.x.x` → `2.0.0`)
*(Klausul 7: Incompatible API Changes)*
1. **CLI Breaking Changes:**
   - Menghapus subcommand (`scan`, `wiki`, `mcp`, `update`) atau flag yang ada (`--format`, `--ext`, `--category`, `--rule`).
   - Mengubah skema output JSON (`--format=json`).
   - Mengubah arti kode keluar (exit code 0/1/2).
2. **Konfigurasi (`charites.yaml`):**
   - Menghapus opsi konfigurasi atau mengubah aturan parsing boolean/severity.
3. **Analyzer & Rules:**
   - Menghapus rule terdaftar atau mengganti Semgrep ID (misal mengubah `theme.hardcode-opacity-color` menjadi ID lain tanpa backwards alias).
4. **MCP Server:**
   - Menghapus tool MCP atau mengubah nama parameter wajib pada tool schema.
5. **Runtime Toolchain:**
   - Menaikkan versi minimum Go toolchain (misal mewajibkan Go 1.27 dari sebelumnya Go 1.26).
6. **Commit Tag:**
   - Commit memuat `!` (`feat!:`, `fix!:`) atau body memuat `BREAKING CHANGE:`.

---

### B. Kriteria MINOR (`1.1.0` → `1.2.0`)
*(Klausul 6: Backward-Compatible New Features & Deprecations)*
1. **Aturan Baru:**
   - Menambahkan rule baru ke dalam registry (misal implementasi batch `theme.hardcode-color`, `a11y.*`, `perf.*`).
2. **CLI & Konfigurasi:**
   - Menambahkan flag CLI baru yang bersifat opsional dengan nilai default yang aman.
   - Menambahkan presenter reporter baru (misal format markdown/HTML).
3. **Fitur MCP:**
   - Menambahkan MCP tool baru atau property opsional baru pada tool schema.
4. **Deprekasi Terencana (*Clause 6 Mandate*):**
   - Menandai rule atau flag sebagai deprecated via komentar `// Deprecated:`.
5. **Commit Tag:**
   - Commit diawali dengan `feat:` atau `feat(...)`: tanpa tanda seru `!`.

---

### C. Kriteria PATCH (`1.1.1` → `1.1.2`)
*(Klausul 5: Backward-Compatible Bug Fixes)*
1. **Perbaikan False Positive (FP) / False Negative (FN):**
   - Memperbaiki parsing Astro frontmatter offset baris.
   - Menambahkan bait handling di sub-korpus adversarial.
2. **Optimasi & Hygiene:**
   - Mengurangi alokasi heap traversal AST via `sync.Pool`.
   - Menutup memory leak atau goroutine leak pada worker pool.
3. **Pengujian & Dokumentasi:**
   - Menambahkan fixture uji coba atau memperbarui `wiki/*.md`.
4. **Commit Tag:**
   - Commit diawali dengan `fix:`, `perf:`, `refactor:`, `docs:`, `test:`, `chore:`.

---

## 5. Prosedur Eksekusi Otomatis

Jalankan skrip inspeksi deterministik dari root repositori:

```bash
# 1. Analisis diff dari tag terakhir ke HEAD (Rekomendasi Bump)
./.agents/skills/charites-versioning/scripts/charites-diff-inspector.sh

# 2. Mode JSON untuk integrasi pipeline CI/CD
./.agents/skills/charites-versioning/scripts/charites-diff-inspector.sh --json

# 3. Generate Release Notes terkurasi (Dry-run pratinjau)
./.agents/skills/charites-versioning/scripts/generate-release-notes.sh

# 4. Tulis ke release_notes.md dan masukkan ke CHANGELOG.md
./.agents/skills/charites-versioning/scripts/generate-release-notes.sh --write
```

---

## 6. Checklist Verifikasi Rilis Charites

- [ ] **1. Identifikasi Tag Basis:** `git describe --tags --abbrev=0` berhasil membaca tag acuan.
- [ ] **2. Three-Dot Diff:** Evaluasi diff menggunakan `git diff ${LATEST_TAG}...HEAD`.
- [ ] **3. Audit 5 Permukaan Publik:** Verifikasi tidak ada breaking change pada CLI, `charites.yaml`, Rules, MCP, dan `internal/ir/`.
- [ ] **4. Audit Deprekasi:** Fitur yang diberi status `Deprecated` mendapatkan bump MINOR (bukan Major).
- [ ] **5. Invariant Reset:** Nilai minor/patch di-reset ke 0 sesuai hierarki SemVer.
- [ ] **6. Imutabilitas Tag:** Tag lama tidak di-overwrite (mematuhi Clause 2).
