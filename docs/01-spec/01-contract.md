# 01-SPEC: 01 - Core Data Contracts (IR & Diagnostic)

> **Kode Dokumen:** `SPEC-01-CONTRACT`
> **Tahapan:** Fase 1 - Kunci Kontrak Data (IR & Diagnostic)
> **Status:** Ready for Review
> **Standar Rujukan:** IETF RFC 2119 / RFC 8174

Dokumen ini mendefinisikan spesifikasi formal kebutuhan fungsional dan representasi data untuk **Intermediate Representation (`ir.Node`)**, **Temuan Diagnosis (`ir.Diagnostic`)**, serta **Interface Evaluasi Rule (`rules.Rule`)** pada mesin Charites.

---

## 1. Filosofi Desain Kontrak Data

Kontrak data pada Fase 1 dirancang dengan prinsip **Leaf-Level SSOT (Single Source of Truth)**:
1. **Zero Circular Dependency:** Seluruh tipe data esensial (`Node`, `Span`, `Diagnostic`, `Severity`) didefinisikan di dalam paket daun `internal/ir`. Paket ini **DILARANG** mengimpor paket `analyzer`, `rules`, `parser`, atau `reporter`.
2. **Framework Agnostic:** Kontrak `ir.Node` menormalkan sintaks heterogen (Astro template, React TSX, HTML) menjadi representasi pohon terpadu tanpa membawa metadata spesifik AST parser pihak ketiga.
3. **Immutability Post-Construction:** Begitu pohon IR selesai dibangun oleh IR Builder, node **MUST** diperlakukan sebagai *read-only* selama traversal evaluasi rule.

---

## 2. Spesifikasi Tipe Data IR (`internal/ir/node.go`)

### 2.1. Enumerasi Tipe Node (`NodeType`)
```go
type NodeType uint8

const (
    NodeElement NodeType = iota // Elemen HTML/JSX (contoh: <div>, <Card>, <button>)
    NodeText                    // Teks literal di dalam elemen
    NodeComment                 // Komentar HTML/JSX (<!-- ... --> atau {/* ... */})
    NodeFragment                // Fragment container (<> ... </>)
)
```

### 2.2. Struktur Rentang Posisi Sumber (`Span`)
Untuk menjamin presisi pelaporan baris dan kolom:
```go
type Span struct {
    Line      int // Posisi baris awal (1-indexed)
    Column    int // Posisi kolom awal (1-indexed)
    EndLine   int // Posisi baris akhir (1-indexed)
    EndColumn int // Posisi kolom akhir (1-indexed)
}
```

### 2.3. Struktur Utama `ir.Node`
```go
type Node struct {
    Type       NodeType          // Jenis node
    Tag        string            // Nama tag elemen (contoh: "div", "button", "Card")
    Attributes map[string]string // Peta atribut mentah (contoh: {"id": "hero", "role": "button"})
    Classes    []string          // Daftar token class yang telah di-split & di-trim
    RawClasses string            // String class asli sebelum di-split
    Span       Span              // Posisi sumber kode asli
    Parent     *Node             // Pointer ke node induk (nil untuk root)
    Children   []*Node           // Irisan node anak
}
```

#### Persyaratan Invarian `ir.Node`:
- **`Classes` Tokenization:** Token class **MUST** dipisahkan berdasarkan spasi (*whitespace-delimited*) dan dibersihkan dari spasi kosong. Contoh atribut `class="p-4  bg-primary "` menghasilkan slice `[]string{"p-4", "bg-primary"}`.
- **1-Indexed Position:** Posisi `Span.Line` dan `Span.Column` **MUST** dimulai dari angka `1` sesuai dengan standar editor dan terminal IDE.
- **Bi-Directional Traversal:** Setiap `Node` yang memiliki anak **MUST** mengisi pointer `Parent` pada setiap anak yang menunjuk kembali ke node tersebut.

---

## 3. Spesifikasi Kontrak Diagnosis (`internal/ir/diagnostic.go`)

### 3.1. Tingkat Keparahan (`Severity`)
```go
type Severity string

const (
    SeverityError Severity = "error" // Pelanggaran fatal (exit code 1)
    SeverityWarn  Severity = "warn"  // Peringatan perbaikan (exit code 1)
    SeverityInfo  Severity = "info"  // Informasi/advis (exit code 0 jika hanya ada info)
)
```

### 3.2. Struktur Temuan (`Diagnostic`)
```go
type Diagnostic struct {
    File     string   `json:"file"`               // Path relatif berkas sumber
    Line     int      `json:"line"`               // Baris lokasi temuan (1-indexed)
    Column   int      `json:"column"`             // Kolom lokasi temuan (1-indexed)
    Rule     string   `json:"rule"`               // Kode unik Semgrep rule (contoh: "theme.hardcode-opacity-color")
    Severity Severity `json:"severity"`           // "error" | "warn" | "info"
    Message  string   `json:"message"`            // Deskripsi pelanggaran ringkas
    Hint     string   `json:"hint,omitempty"`     // Rekomendasi tindakan remediasi
}
```

#### Persyaratan Invarian `Diagnostic`:
- **Serialisasi JSON Flat:** Tag JSON **MUST** menghasilkan struktur datar (*flat array*) untuk mempermudah integrasi dengan `jq`, editor IDE, dan payload Model Context Protocol (MCP).
- **Deterministik:** Diagnosis dengan nilai `File`, `Line`, `Column`, dan `Rule` yang sama **MUST** menghasilkan representasi JSON byte-per-byte yang identik.

---

## 4. Spesifikasi Interface Baku Evaluasi (`rules.Rule`)

Interface evaluasi rule didefinisikan pada `internal/rules/rule.go` dengan mengonsumsi kontrak `internal/ir`:

```go
package rules

import (
    "github.com/will2469/charites/internal/ir"
)

type Rule interface {
    ID() string                                   // Identifier Semgrep tunggal (contoh: "theme.hardcode-opacity-color")
    Description() string                          // Penjelasan maksud dan tujuan rule
    Category() string                             // Kategori (theme, a11y, perf, layout, seo)
    DefaultSeverity() ir.Severity                 // Default severity jika tidak dioverride config
    Evaluate(node *ir.Node) []ir.Diagnostic       // Fungsi evaluasi murni tanpa side-effect I/O
}
```

#### Aturan Evaluasi Murni (*Pure Function Requirement*):
1. **No I/O Operations:** Fungsi `Evaluate()` dilarang melakukan operasi disk/jaringan atau memodifikasi tree `ir.Node`.
2. **Deterministic Output:** Pemanggilan `Evaluate(node)` berulang kali pada node yang sama **MUST** menghasilkan slice diagnosis yang identik tanpa bergantung pada state global.

---

## 5. Spesifikasi Direktif Pengabaian Inline (*Inline Ignore Directives*)

Sistem scanner dan analyzer **MUST** mendukung komentar pengabaian inline (*inline ignore*) untuk mengecualikan node tertentu dari evaluasi rule:

1. **Format Single & Multi-Rule:**
   - Sintaks: `charites:ignore <rule1>[, <rule2>...]`
   - Pola pencocokan: token rule dipisahkan oleh tanda koma `,` atau spasi.
   - Contoh di berkas TypeScript / JSX:
     ```tsx
     // charites:ignore theme.hardcode-opacity-color
     <div className="bg-primary/10" />

     // Multi-rule ignore (dipisahkan koma atau spasi):
     // charites:ignore theme.hardcode-opacity-color, theme.hardcode-color
     <button className="bg-primary/10 text-slate-500" />
     ```
   - Contoh di berkas template Astro / HTML:
     ```astro
     <!-- charites:ignore theme.hardcode-opacity-color -->
     <div class="bg-primary/10"></div>
     ```
2. **Aturan Evaluasi Pengabaian:**
   - Jika sebuah node didahului atau memuat komentar direktif `charites:ignore`, diagnostic yang memiliki `Rule` yang terdaftar dalam komentar tersebut **MUST** dibatalkan (*suppressed*) dan tidak dilaporkan ke output akhir.
