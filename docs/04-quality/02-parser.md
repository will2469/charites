# 04-QUALITY: 02 - Parser Robustness & Security Constraints

> **Kode Dokumen:** `QUAL-02-PARSER`
> **Tahapan:** Fase 2 - Parser Frontend & IR Builder
> **Peran Pilar:** QUALITY = QUALITY THRESHOLD (Ambang Batas Ketahanan, Security & Resource Budgets)
> **Status:** Ready for Execution
> **Standar Rujukan:** OpenSSF Best Practices & Resilient Parser Guidelines

Dokumen ini mendefinisikan ambang batas ketahanan (*robustness*), pencegahan eksploitasi rekursi (*stack overflow*), dan batas efisiensi anggaran sumber daya (*resource budgets*) pada modul parser frontend.

---

## 1. Invarian Ketahanan Tanpa Crash (Zero-Panic Invariant)

Parser frontend Charites bertindak sebagai garda terdepan sistem. Aturan ketahanan absolut:
1. **No Panic on Malformed Input:** Input sintaks apapun-mulai dari berkas kosong, tag HTML yang tidak ditutup, kutip tidak seimbang, hingga karakter biner acak-**DILARANG KERAS** memicu panic di Go runtime.
2. **Deterministic Recovery:** Kegagalan sintaks pada satu tag tidak boleh menghentikan pemindaian keseluruhan berkas. Parser memulihkan diri secara deterministik ke tag pembuka berikutnya tanpa menghasilkan node parsial yang rusak.
3. **No External Process Execution:** Dilarang memanggil Node.js, `esbuild`, atau shell eksternal saat proses parsing.

---

## 2. Pencegahan Serangan Stack Overflow & DoS

Untuk melindungi sistem dari berkas eksploitasi rekursi ekstrem (*deeply nested tags*):
- **Batas Kedalaman Elemen (Maximum Nesting Depth):** Parser membatasi kedalaman hierarki elemen hingga maksimal **256 tingkat nesting**. Jika melebihi 256 tingkat, elemen berikutnya diperlakukan datar (*flattened*) untuk mencegah *stack overflow*.
- **Kompleksitas Waktu Linear:** Algoritma pemindaian dan ekstraksi struktural **SHOULD** beroperasi dalam waktu linear $O(N)$ proporsional terhadap ukuran byte berkas pada grammar subset yang didukung dengan batas nesting 256. Dilarang keras menggunakan *regular expression* yang rentan terhadap *Catastrophic Backtracking* (ReDoS).

---

## 3. Ambang Batas Kualitas & Anggaran Sumber Daya (Resource Budgets)

| Indikator | Ambang Batas (*Threshold*) | Cara Evaluasi | Klasifikasi |
| :--- | :--- | :--- | :--- |
| **Batas Kedalaman Stack** | Maksimal 256 level hierarki | `TestNode_MaxNestingDepth` | Hard Gate |
| **Batas Memori per Berkas** | Peak heap live bytes $\le 4\times$ ukuran berkas sumber | `runtime.MemStats` pada benchmark corpus | Memory Budget |
| **Throughput Parsing** | Baseline terukur pada fixture standar | `go test -bench=BenchmarkParser -benchmem` | Performance Baseline |
| **Fuzzing Endurance** | 0 crash dalam minimal 60 detik mutasi | `go test -fuzz=. -fuzztime=60s` | Hard Gate |
| **Coverage Uji Paket Parser** | $\ge 85\%$ line coverage | `go test -cover ./internal/parser/...` | Hard Gate |

