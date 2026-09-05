# 04-QUALITY: 01 - Data Contract Hygiene & Memory Constraints

> **Kode Dokumen:** `QUAL-01-CONTRACT`
> **Tahapan:** Fase 1 - Kunci Kontrak Data (IR & Diagnostic)
> **Peran Pilar:** QUALITY = QUALITY THRESHOLD (Ambang Batas Kualitas Linter, Zero Dep & Hygiene)
> **Status:** Ready for Execution
> **Standar Rujukan:** OpenSSF Best Practices & Go Staticcheck Rules

Dokumen ini menetapkan ambang batas kualitas kode statis, kebersihan tipe data, dan batas anggaran performa (*performance budget*) untuk paket kontrak inti **`internal/ir`**.

---

## 1. Invarian Kebersihan Paket `internal/ir`

Sebagai fondasi *Single Source of Truth* seluruh arsitektur, paket `internal/ir` harus mematuhi aturan ketat:
1. **Zero External Dependency:** Dilarang menambahkan dependensi pihak ketiga (`third-party dependencies`).
2. **Zero Internal Dependency:** Paket `internal/ir` dilarang mengimpor paket `internal/*` mana pun untuk mencegah keterikatan sirkular (*circular coupling*).
3. **No Unsafe Pointer Manipulation:** Seluruh operasi pointer di dalam `internal/ir` wajib menggunakan pointer Go standar yang aman (*safe Go*), dilarang menggunakan `unsafe.Pointer` pada logika produksi. Penggunaan `unsafe.Sizeof` hanya diizinkan di dalam berkas pengujian (`_test.go`) untuk memverifikasi batas memori struct.

---

## 2. Standar Kualitas Analisis Statis

Pemeriksaan linter wajib lolos 100% pada kode `internal/ir`:
- **`govet`:** Memastikan kebenaran semantik Go standar (tanpa menyandarkan optimasi field alignment ke govet, karena field alignment diverifikasi langsung di unit test).
- **`staticcheck`:** Memastikan tidak ada perbandingan tipe yang tidak aman pada enum `NodeType` dan `Severity`.
- **`unused`:** Seluruh field dan metode helper publik wajib memiliki peruntukan fungsional yang jelas.
- **`errcheck`:** JSON serialization helper tidak boleh mengabaikan error secara diam-diam.

---

## 3. Ambang Batas Kualitas & Anggaran Performa (Quality Thresholds)

| Indikator | Ambang Batas (*Threshold*) | Metode Evaluasi | Klasifikasi |
| :--- | :--- | :--- | :--- |
| **Ukuran Struct `ir.Node`** | $\le 136$ bytes pada target 64-bit | `unsafe.Sizeof(ir.Node{})` di unit test | Hard Gate |
| **Determinisme Total Order** | 100% byte-identical under permutations | `TestDiagnostic_CollectionOrdering` | Hard Gate (`QUAL-01-INVAR-001`) |
| **Alokasi Iterator `Walk()`** | Target: **0 B/op** & **0 allocs/op** | `go test -bench=BenchmarkNode_Walk -benchmem` | Performance Budget |
| **Coverage Uji Paket `ir`** | $\ge 90\%$ line coverage | `go test -cover ./internal/ir/...` | Hard Gate |
| **Sirkularitas Dependensi** | **0 cycles** | `go vet ./internal/...` | Hard Gate |

> [!NOTE]
> - **QUAL-01-INVAR-001 (Deterministic Total Ordering):** Seluruh diagnosis yang dikumpulkan dari berbagai goroutine wajib melewati `ir.SortDiagnostics()` dengan 7-kunci `DiagnosticOrderKey`. Permutasi urutan kedatangan tidak boleh mengubah hasil akhir byte output JSON.
> - Alokasi iterator `Walk()` dikelola sebagai **Performance Budget**. Deviasi alokasi ($> 0$) pada compiler environment tertentu memicu audit regresi performa, tanpa membatalkan keabsahan kontrak fungsional `ir.Node`.


