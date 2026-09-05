# 04-QUALITY: 02 - Parser Robustness & Security Constraints

> **Kode Dokumen:** `QUAL-02-PARSER`
> **Tahapan:** Fase 2 - Parser Frontend & IR Builder
> **Status:** Ready for Review
> **Standar Rujukan:** OpenSSF Best Practices & Resilient Parser Guidelines

Dokumen ini mendefinisikan batasan ketahanan (*robustness*), pencegahan eksploitasi rekursi (*stack overflow*), dan batas efisiensi waktu pemrosesan pada modul parser frontend.

---

## 1. Invarian Ketahanan Tanpa Crash (Zero-Panic Invariant)

Parser frontend Charites bertindak sebagai garda terdepan sistem. Aturan ketahanan absolut:
1. **No Panic on Malformed Input:** Input sintaks apapun-mulai dari berkas kosong, tag HTML yang tidak ditutup, kutip tidak seimbang, hingga karakter biner acak-**DILARANG KERAS** memicu panic di Go runtime.
2. **Error Recovery:** Kegagalan sintaks pada satu tag tidak boleh menghentikan pemindaian keseluruhan berkas. Parser wajib memulihkan diri (*graceful recovery*) ke tag valid berikutnya.
3. **No External Process Execution:** Dilarang memanggil Node.js, `esbuild`, atau shell eksternal saat proses parsing.

---

## 2. Pencegahan Serangan Stack Overflow & DoS

Untuk melindungi sistem dari berkas jahat atau rekursi ekstrem (*deeply nested tags*):
- **Batas Kedalaman Tag (Max Nesting Depth):** Parser membatasi kedalaman hirarki elemen hingga maksimal **256 tingkat nesting**. Jika melebihi 256 tingkat, elemen berikutnya diperlakukan datar (*flattened*) untuk mencegah *stack overflow*.
- **Kompleksitas Waktu Linear ($O(N)$):** Algoritma tokenisasi dan parsing wajib berjalan dalam waktu linear $O(N)$ proporsional terhadap ukuran byte berkas. Dilarang menggunakan *regular expression* yang rentan terhadap *Catastrophic Backtracking* (ReDoS).

---

## 3. Ambang Batas Kualitas Memori & Throughput

| Indikator | Ambang Batas (*Threshold*) | Cara Evaluasi |
| :--- | :--- | :--- |
| **Kecepatan Parsing** | $\ge 20.000$ lines of code per detik per CPU core | `go test -bench=BenchmarkParser -benchmem` |
| **Batas Memori per Berkas** | Alokasi heap $\le 3\times$ ukuran berkas sumber | Pengukuran heap profiler pprof |
| **Fuzzing Endurance** | 0 crash dalam 60 detik mutasi | `go test -fuzz=FuzzAstroParser -fuzztime=60s` |
| **Coverage Uji Paket Parser** | $\ge 85\%$ line coverage | `go test -cover ./internal/parser/...` |
