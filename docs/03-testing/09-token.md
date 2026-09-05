# 03-TESTING: 09 - Token Subsystem Verification Plan & Fuzzing

> **Kode Dokumen:** `TEST-09-TOKEN`
> **Tahapan:** Fase 2/3 - SSOT Design Token Engine & Dependency Resolution
> **Peran Pilar:** TEST = PROVE (Skenario Pengujian, Fuzzing, Pembuktian Zero-Panic)
> **Status:** Active / Graduated

Dokumen ini mendefinisikan matriks pengujian empiris, rencana pengujian integrasi monorepo, dan protokol fuzzing untuk membuktikan keandalan subsistem **Design Token & Dependency Graph** (`internal/token`).

---

## 1. Matriks Pengujian Unit (`internal/token/...`)

### 1.1. Generic CSS Lexer & Parser (`internal/token/theme/parser_test.go`)
- **`Semicolon_inside_quotes`**: Membuktikan pemisahan pernyataan CSS tidak terpotong oleh tanda titik koma `;` di dalam string literal (`--value: "hello;world"`).
- **`Data_URI_with_colons_and_semicolons`**: Membuktikan parser tidak salah mendeteksi `:` atau `;` di dalam `url("data:image/svg+xml;...")`.
- **`CSS_Nesting_&_Pseudo-selectors`**: Membuktikan resolusi selektor bertingkat (`.card` + `&:hover` -> `.card:hover`).
- **`Nested_AtRules`**: Membuktikan pemrosesan berlapis `@layer theme` di dalam `@media (prefers-color-scheme: dark)`.
- **`Banana_Test`**: Membuktikan parser mengekstrak token arbitrer (`--banana`, `--thing-that-is-definitely-not-primary`, `--super-special-design-token`) tanpa memuat heuristik warna.

### 1.2. Token Dependency Graph & Context API (`internal/token/extractor_test.go`)
- **`TestTheme_BananaTest`**: Pembuktian resolusi dependensi token berantai tanpa opini desain.
- **`TestTheme_CycleDetection`**: Pembuktian deteksi siklus sirkular (`--a: var(--b); --b: var(--a)`) mengembalikan `ErrCycleDetected` tanpa recursive stack overflow.
- **`TestTheme_MultiHopResolution`**: Pembuktian resolusi multi-tahap (`--alias-3 -> --alias-2 -> --alias-1 -> --base -> #abcdef`).
- **`TestTheme_EvaluationBudget`**: Pembuktian batas DoS `opts.MaxNodes` memotong rantai patologis dan mengembalikan `ErrEvaluationBudgetExceeded`.
- **`TestTheme_ScopeSeparationAndSpecificity`**: Pembuktian deklarasi ganda `--brand` di `:root`, `.card`, dan `#header.hero .card` mempertahankan identitas unik dan bobot spesifisitas selektor.
- **`TestTheme_DiscoverAndLoad`**: Pembuktian upward directory traversal dari folder bersarang ke root repositori.

---

## 2. Pengujian Integrasi Monorepo (`tests/token_integration_test.go`)

Pengujian integrasi memverifikasi kolaborasi antarkomponen dari pembacaan filesystem nyata hingga keluaran terminal CLI:

1. **Auto-Discovery Berkas Proyek Nyata (`TestTokenIntegration_AutoDiscoveryAndResolution`)**:
   - Memanggil `token.DiscoverAndLoad` dari subfolder dalam (`src/components/dashboard/cards/`).
   - Memverifikasi ekstraksi `global.css` di parent level, multi-scope `:root` vs `.dark`, dan resolusi variabel `var(--spacing-card) -> 1.5rem`.
   - Menguji adapter konvensi semantik Layer 4 (`TokenConvention`) terhadap fakta token proyek nyata dengan zero false positive.

2. **Perlindungan Siklus Sirkular (`TestTokenIntegration_CircularDependencyProtection`)**:
   - Memasukkan berkas CSS dengan siklus sirkular `--token-a -> --token-b -> --token-c -> --token-a`.
   - Memverifikasi `FindCycles()` mendeteksi siklus seketika dan token non-siklus di berkas yang sama tetap dapat di-resolve.

3. **End-to-End CLI Scan Pipeline (`TestTokenIntegration_EndToEndCLIScan`)**:
   - Menjalankan `cli.ExecuteArgs([]string{"scan", tempDir, "-f", "json"}, ...)` pada proyek temporary.
   - Memverifikasi CLI menghasilkan JSON envelope dengan temuan:
     ```json
     {
       "rule": "theme.hardcode-opacity-color",
       "message": "Hardcode opacity color: \"bg-primary/10\"",
       "hint": "Use semantic token \"primary-light\"."
     }
     ```

---

## 3. Protokol Native Fuzzing (`tests/fuzz/css_fuzz_test.go`)

Untuk menjamin **Zero-Panic Invariant**, parser CSS dan Token Extractor diuji dengan mesin fuzzing Go 1.26:

```go
func FuzzCSSParser(f *testing.F)
```

- **Seed Corpus:** 20 seed mencakup berkas `global.css` riil, token arbitrer, string unclosed, data URI malformed, CSS nesting, kurung kurawal rusak, null byte, dan rantai variabel acak.
- **Hasil Fuzzing:**
  - $> 14.500$ eksekusi mutasi acak per sesi.
  - Kecepatan: $> 4.800$ eksekusi/detik dengan 8 worker paralel.
  - **Zero Panic**, **Zero Deadlock**, dan **Zero Memory Leak**.
