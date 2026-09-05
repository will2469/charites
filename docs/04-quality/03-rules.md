# 04-QUALITY: 03 - Rule Quality Constraints & Pure Function Invariants

> **Kode Dokumen:** `QUAL-03-RULES`
> **Tahapan:** Fase 3 - Rule Contract & Proving Ground Rule (`theme.hardcode-opacity-color`)
> **Status:** Ready for Review
> **Standar Rujukan:** OpenSSF Best Practices & Deterministic Static Analyzer Invariants

Dokumen ini mendefinisikan batasan kualitas ketat, penegakan fungsi evaluasi murni (*pure evaluation function*), keamanan konkurensi antar-goroutine, dan kepatuhan alokasi memori pada lapisan rule Charites.

---

## 1. Invarian Fungsi Evaluasi Murni (Pure Function Invariant)

Fungsi `rule.Evaluate(node *ir.Node) []ir.Diagnostic` merupakan inti dari static analysis engine. Modul rule **WAJIB** tunduk pada 4 prinsip kemurnian:

1. **Zero Disk / Network I/O:**
   Dilarang keras membuka berkas, membaca konfigurasi runtime, memanggil command shell, atau melakukan HTTP request di dalam `Evaluate()`. Semua konteks yang dibutuhkan rule harus sudah berada di dalam memori (`*ir.Node` dan token lookup table).
2. **Read-Only Invariant pada AST Node:**
   Dilarang memodifikasi field apapun pada pointer `*ir.Node` (misal: menambahkan class, mengubah span, atau memanipulasi pointer `Parent`/`Children`). AST diperlakukan sebagai struktur data *immutable*.
3. **Idempotensi Deterministik:**
   Untuk sembarang pointer `node`, pemanggilan berulang `rule.Evaluate(node)` di thread manapun **PASTI** menghasilkan diagnostic yang persis sama: urutan, pesan, line, dan column identik ($f(x) = f(x)$).
4. **Stateless Rule Instances:**
   Struct implementasi rule dilarang menyimpan state transien dari evaluasi sebelumnya (dilarang menyimpan counter, cache mutable tanpa sinkronisasi, atau pointer ke node aktif di dalam field struct rule).

---

## 2. Invarian Keamanan Konkurensi (Concurrency Safety)

Pada Fase 4, scanner Charites akan mendistribusikan puluhan ribu file dan ratusan ribu AST node ke worker pool goroutine:

- **Shared-Read Safe:** Satu instance singleton dari setiap rule (misal `theme.NewHardcodeOpacityColorRule()`) didaftarkan ke `Registry` sekali saja saat startup.
- **Concurrent Evaluators:** Puluhan goroutine akan memanggil `rule.Evaluate(node)` pada instance singleton yang sama secara paralel. Karena rule bersifat murni dan stateless, evaluasi dapat berjalan dengan skalabilitas paralel linear ($O(P)$) tanpa risiko data race atau contention mutex pada logika rule.
- **Thread-Safe Registry:** Operasi pendaftaran dan kueri registri dilindungi penuh oleh `sync.RWMutex`.

---

## 3. Ambang Batas Efisiensi Memori & Alokasi

Setiap pemindaian proyek frontend berskala enterprise dapat memproses jutaan AST node. Efisiensi alokasi memori pada fungsi evaluasi rule menjadi penentu performa:

| Kondisi Evaluasi | Batas Alokasi Memori | Batas Alokasi Objek | Batas Waktu Eksekusi |
| :--- | :---: | :---: | :---: |
| **Node Bersih (Tanpa Class)** | `0 B/op` | `0 allocs/op` | $\le 10\text{ ns/op}$ |
| **Node Bersih (Class Legal)** | `0 B/op` | `0 allocs/op` | $\le 50\text{ ns/op}$ |
| **Node Pelanggaran (1 Pelanggaran)** | $\le 128\text{ B/op}$ | $\le 2\text{ allocs/op}$ | $\le 250\text{ ns/op}$ |

### Aturan Alokasi Defensif:
- Selalu gunakan slice `nil` alih-alih `make([]ir.Diagnostic, 0)` ketika tidak ada diagnostic yang dihasilkan, sehingga tidak memicu alokasi heap kosong di Go runtime.

---

## 4. Standar Kualitas Kode & Metrik Kelulusan

| Metrik Kualitas | Ambang Batas Minimum | Cara Pengukuran |
| :--- | :---: | :--- |
| **Statement Coverage (`internal/rules/...`)** | $\ge 95\%$ | `go test -cover ./internal/rules/...` |
| **Branch Coverage Logika Rule** | $100\%$ | Pengujian seluruh cabang if/else pada `hardcode_opacity_color.go` |
| **Cyclomatic Complexity** | $\le 10$ per fungsi | `gocyclo -over 10 ./internal/rules` |
| **Data Race Verification** | $0$ data race detected | `go test -race ./internal/rules/...` |
| **Linter Compliance** | $0$ issues | `golangci-lint run ./internal/rules/...` |
