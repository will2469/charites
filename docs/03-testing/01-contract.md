# 03-TESTING: 01 - Data Contract & Iterator Verification Plan

> **Kode Dokumen:** `TEST-01-CONTRACT`
> **Tahapan:** Fase 1 - Kunci Kontrak Data (IR & Diagnostic)
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Skenario Smoke & Asersi Pembuktian)
> **Status:** Graduated (All Phase Gates Passed)
> **Standar Rujukan:** Modern Testing Principles & Go Benchmarking Standards

Dokumen ini mendefinisikan skenario pembuktian empiris, rencana pengujian unit komprehensif, dan benchmark alokasi memori untuk membuktikan pemenuhan seluruh kontrak fungsional pada **`internal/ir`**.

---

## 1. Skenario Uji Coba Unit (`internal/ir/node_test.go`)

Paket `internal/ir` wajib memiliki pengujian unit deterministik yang membuktikan seluruh invarian normatif:

### Test Case 1: Traversal Pohon Pre-Order (`TestNode_Walk`)
- **Tujuan:** Memastikan iterator `Walk()` mengunjungi seluruh node anak secara teratur sesuai urutan *depth-first search*.
- **Struktur Pohon Uji:**
  ```text
  root (div)
  ├── child1 (span)
  │   └── grandchild1 (text)
  └── child2 (button)
  ```
- **Asersi Pembuktian:**
  - Node dikunjungi dengan urutan presisi: `root` -> `child1` -> `grandchild1` -> `child2`.
  - Total kunjungan berjumlah 4 node.

### Test Case 2: Terminasi Awal Iterator (*Early Exit*)
- **Tujuan:** Memastikan perulangan `for n := range root.Walk()` dapat dihentikan seketika di tengah jalan menggunakan `break` tanpa memicu panic atau traversal lanjutan.
- **Asersi Pembuktian:** Traversal berhenti seketika saat callback `yield` mengembalikan `false`.

### Test Case 3: Invarian Hubungan Induk-Anak (`TestNode_ParentChildInvariant`)
- **Tujuan:** Membuktikan invarian traversal dua arah (*bidirectional traversal*).
- **Asersi Pembuktian:**
  - Untuk setiap node anak, pointer `child.Parent` menunjuk secara identik ke node induknya (`child.Parent == parent`).
  - Irisan `parent.Children` memuat pointer `child`.
  - Node akar (*root*) memiliki `Parent == nil`.

### Test Case 4: Tokenisasi Kelas Semantik (`TestNode_ClassTokenization`)
- **Tujuan:** Membuktikan tokenisasi kelas CSS berbasis whitespace.
- **Input:** Atribut kelas dengan berbagai variasi spasi (`"p-4   bg-primary   text-sm\t\n"`).
- **Asersi Pembuktian:**
  - Menghasilkan slice `[]string{"p-4", "bg-primary", "text-sm"}` tanpa elemen string kosong atau whitespace tersisa.
  - Nilai `RawClasses` tetap mempertahankan string asli yang belum di-split.

### Test Case 5: Pengindeksan Rentang Posisi (`TestNode_SpanIndexing`)
- **Tujuan:** Membuktikan kepatuhan format posisi 1-indexed.
- **Asersi Pembuktian:**
  - `Span.Line >= 1` dan `Span.Column >= 1`.
  - `Span.EndLine >= Span.Line`.
  - Jika satu baris, `Span.EndColumn >= Span.Column`.

### Test Case 6: Verifikasi Ukuran Memori Struct (`TestNode_StructSize`)
- **Tujuan:** Membuktikan optimasi memory layout dan field alignment pada arsitektur 64-bit.
- **Asersi Pembuktian:** `unsafe.Sizeof(ir.Node{}) <= 136` bytes.

### Test Case 7: Helper Metode (`TestNode_Helpers`)
- Menguji `HasClass("bg-primary")` bernilai `true` dan `HasClass("missing")` bernilai `false`.
- Menguji `GetAttr("role")` mengembalikan nilai dan flag exist yang valid.
- Menguji `IsElement("div", "span")` berfungsi secara variadic dengan benar.

---

## 2. Skenario Uji Coba Serialisasi & Pengurutan Diagnostic (`internal/ir/diagnostic_test.go`)

### Test Case 1: Struktur JSON Flat & Determinisme Byte-Level (`TestDiagnostic_JSONDeterminism`)
- **Tujuan:** Membuktikan bahwa `Diagnostic` menghasilkan serialisasi flat JSON deterministik tanpa objek bersarang.
- **Asersi Pembuktian:**
  - Hasil `json.Marshal(diag)` menghasilkan flat JSON object dengan kunci tepat: `file`, `line`, `column`, `rule`, `severity`, `message`, `hint` (dengan `hint` di-omit jika kosong).
  - Dua instance `Diagnostic` dengan seluruh nilai field identik **MUST** menghasilkan byte array JSON yang identik 100%.
  - Serialisasi array `[]Diagnostic` menghasilkan JSON array flat yang deterministik.

### Test Case 2: Determinisme Total Pengurutan Koleksi (`TestDiagnostic_CollectionOrdering`)
- **Tujuan:** Membuktikan stabilitas total urutan temuan diagnosis menggunakan 7-tingkat `DiagnosticOrderKey` untuk mengatasi ambiguitas arrival order dari concurrent worker pool.
- **Asersi Pembuktian:**
  - **7-Level Strict Tie-Breaking:** Menguji pemecah seri (*tie-breaking*) secara bertingkat:
    1. Berkas berbeda (`a.tsx` vs `b.tsx`) $\rightarrow$ terurut menurut nama berkas.
    2. Berkas sama, baris berbeda (Line 5 vs Line 10) $\rightarrow$ terurut menurut nomor baris.
    3. Baris sama, kolom berbeda (Col 2 vs Col 8) $\rightarrow$ terurut menurut nomor kolom.
    4. Lokasi sama, rule berbeda (`a11y.foo` vs `theme.bar`) $\rightarrow$ terurut menurut Rule ID.
    5. Rule sama, severity berbeda (`error` vs `warn`) $\rightarrow$ `error` mendahului `warn`.
    6. Severity sama, message berbeda (`"A"` vs `"B"`) $\rightarrow$ terurut leksikal menurut `Message`.
    7. Message sama, hint berbeda (`"Fix A"` vs `"Fix B"`) $\rightarrow$ terurut leksikal menurut `Hint`.
  - **Permutation Invariance:** Slice diagnosis yang diacak urutannya (*shuffled order*) setelah diproses dengan `ir.SortDiagnostics(diags)` **MUST** menghasilkan urutan elemen dan byte JSON yang 100% identik.
  - **Idempotent Deduplication:** Dua item diagnosis dengan seluruh 7 field yang identik otomatis dipangkas menjadi 1 elemen unik.

---

## 3. Benchmark Pengukuran Alokasi Memori (`BenchmarkNode_Walk`)

Benchmark dijalankan untuk mengukur performa alokasi pada *hot-path* traversal:

```go
func BenchmarkNode_Walk(b *testing.B) {
    root := buildMockSubtree(100) // Membuat 100 node bersarang
    b.ReportAllocs()

    for b.Loop() {
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
### Evaluasi Alokasi (Performance Budget Target):
- **Target Desain:** **`0 B/op`** dan **`0 allocs/op`** di bawah optimasi compiler inlining Go 1.26.
- **Kebijakan Ambang Batas:** Jika hasil benchmark terukur mengalami deviasi alokasi ($> 0$), hasil tersebut dicatat sebagai **Performance Regression Warning** untuk diinvestigasi lebih lanjut, tanpa otomatis membatalkan validitas fungsional kontrak `ir.Node`.

