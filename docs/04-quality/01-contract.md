# 04-QUALITY: 01 - Data Contract Hygiene & Memory Constraints

> **Kode Dokumen:** `QUAL-01-CONTRACT`
> **Tahapan:** Fase 1 - Kunci Kontrak Data (IR & Diagnostic)
> **Status:** Ready for Review
> **Standar Rujukan:** OpenSSF Best Practices & Go Staticcheck Rules

Dokumen ini menetapkan batasan kualitas kode statis, kebersihan tipe data, dan batas konsumsi memori heap untuk paket kontrak inti **`internal/ir`**.

---

## 1. Invarian Kebersihan Paket `internal/ir`

Sebagai fondasi *Single Source of Truth* seluruh arsitektur, paket `internal/ir` harus mematuhi aturan ketat:
1. **Zero External Dependency:** Dilarang menambahkan dependensi pihak ketiga (`third-party dependencies`).
2. **Zero Internal Dependency:** Paket `internal/ir` dilarang mengimpor paket `internal/*` mana pun untuk mencegah keterikatan sirkular (*circular coupling*).
3. **No Unsafe Pointer Manipulation:** Seluruh operasi pointer di dalam `internal/ir` wajib menggunakan pointer Go standar yang aman (*safe Go*), dilarang menggunakan `unsafe.Pointer` kecuali untuk pengujian pengukuran ukuran memori.

---

## 2. Standar Kualitas Analisis Statis

Pemeriksaan linter wajib lolos 100% pada `internal/ir`:
- **`govet` (Struct Alignment):** Memastikan urutan field struct tidak menimbulkan fragmentasi memori berlebih.
- **`staticcheck`:** Memastikan tidak ada perbandingan tipe yang tidak aman pada enum `NodeType` dan `Severity`.
- **`unused`:** Seluruh field dan metode helper publik wajib memiliki peruntukan fungsional yang jelas.
- **`errcheck`:** JSON serialization helper tidak boleh mengabaikan error secara diam-diam.

---

## 3. Ambang Batas Kualitas Memori (Quality Gate)

| Indikator | Ambang Batas (*Threshold*) | Evaluasi |
| :--- | :--- | :--- |
| **Ukuran Struct `ir.Node`** | $\le 136$ bytes pada OS 64-bit | `unsafe.Sizeof(ir.Node{})` |
| **Alokasi Iterator `Walk()`** | **0 B/op** dan **0 allocs/op** | `go test -bench=BenchmarkNode_Walk -benchmem` |
| **Coverage Uji Paket `ir`** | $\ge 90\%$ line coverage | `go test -cover ./internal/ir/...` |
| **Sirkularitas Dependensi** | **0 cycles** | `go vet ./internal/...` |
