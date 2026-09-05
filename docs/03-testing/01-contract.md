# 03-TESTING: 01 - Data Contract & Iterator Verification Plan

> **Kode Dokumen:** `TEST-01-CONTRACT`
> **Tahapan:** Fase 1 - Kunci Kontrak Data (IR & Diagnostic)
> **Status:** Ready for Review
> **Standar Rujukan:** Modern Testing Principles & Go Benchmarking Standards

Dokumen ini mendefinisikan rencana pengujian unit dan pengukuran performa alokasi memori untuk paket kontrak data **`internal/ir`**.

---

## 1. Skenario Uji Coba Unit (`internal/ir/node_test.go`)

### Test Case 1: Traversal Pohon Pre-Order (`TestNode_Walk`)
- **Tujuan:** Memastikan iterator `Walk()` mengunjungi seluruh node anak secara teratur sesuai urutan *depth-first search*.
- **Struktur Pohon Uji:**
  ```text
  root (div)
  ├── child1 (span)
  │   └── grandchild1 (text)
  └── child2 (button)
  ```
- **Ekspektasi:**
  - Node dikunjungi dengan urutan tepat: `root` -> `child1` -> `grandchild1` -> `child2`.
  - Total kunjungan berjumlah 4 node.

### Test Case 2: Terminasi Awal Iterator (*Early Exit*)
- **Tujuan:** Memastikan perulangan `for n := range root.Walk()` dapat dihentikan di tengah jalan menggunakan `break` tanpa memicu memory leak atau panic.
- **Ekspektasi:** Traversal berhenti seketika saat `yield` mengembalikan `false`.

### Test Case 3: Helper Metode (`TestNode_Helpers`)
- Menguji `HasClass("bg-primary")` menghasilkan `true` dan `HasClass("missing")` menghasilkan `false`.
- Menguji `GetAttr("role")` mengembalikan nilai dan boolean exist yang valid.
- Menguji `IsElement("div", "span")` berfungsi variadic dengan benar.

---

## 2. Skenario Uji Coba Serialisasi Diagnostic (`internal/ir/diagnostic_test.go`)

### Test Case 1: JSON Marshaling Determinism
- **Input:** Struct `Diagnostic` terisi lengkap.
- **Ekspektasi:**
  - Hasil `json.Marshal(diag)` menghasilkan flat JSON dengan kunci tepat: `file`, `line`, `column`, `rule`, `severity`, `message`, `hint`.
  - Tidak ada field kosong yang tidak perlu.

---

## 3. Benchmark Alokasi Memori (Zero-Alloc Gate)

Mekanisme traversal pohon IR adalah *hot-path* yang dieksekusi jutaan kali. Benchmark alokasi wajib dijalankan:

```go
func BenchmarkNode_Walk(b *testing.B) {
    root := buildMockSubtree(100) // Membuat 100 node bersarang
    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        count := 0
        for range root.Walk() {
            count++
        }
        if count != 100 {
            b.Fatalf("expected 100 nodes, got %d", count)
        }
    }
}
```

### Kriteria Penerimaan Benchmark:
- **Alokasi Heap:** **`0 B/op`**
- **Jumlah Alokasi:** **`0 allocs/op`**
- Jika ditemukan alokasi $> 0$, Pull Request **DITOLAK**.
