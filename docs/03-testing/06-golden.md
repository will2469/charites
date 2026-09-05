# 03-TESTING: 06 - Golden Snapshot Regression & Fuzzing Verification Plan

> **Kode Dokumen:** `TEST-06-GOLDEN`
> **Tahapan:** Fase 6 - Validasi Penuh & Golden Snapshots (Milestone Selesai Pipa)
> **Status:** Ready for Review
> **Standar Rujukan:** Snapshot Regression Testing & Native Fuzzing Protocols

Dokumen ini mendefinisikan rencana pengujian regresi menyeluruh (*regression testing*) berbasis **Golden Snapshots**, protokol eksekusi *native fuzzing*, serta metrik validasi kestabilan pipa compiler secara menyeluruh.

---

## 1. Skenario Pengujian Golden Snapshot (`tests/golden_test.go`)

Seluruh skenario pengujian diuji secara otomatis dari lapisan terluar CLI hingga lapisan terdalam parser:

### 1.1. Matriks Skenario Snapshot
| Nama Skenario | Path Fixture | Target Validasi | Golden File Terkait |
| :--- | :--- | :--- | :--- |
| `astro_opacity` | `tests/fixtures/astro_opacity` | Pelanggaran opacity di Astro, verifikasi line offset frontmatter presisi | `astro_opacity.golden.json`<br>`astro_opacity.golden.txt` |
| `tsx_opacity` | `tests/fixtures/tsx_opacity` | Pelanggaran di file React TSX, template literal, JSX attributes | `tsx_opacity.golden.json`<br>`tsx_opacity.golden.txt` |
| `clean_project` | `tests/fixtures/clean_project` | Repositori bersih tanpa pelanggaran (*Zero Noise Invariant*) | `clean_project.golden.json`<br>`clean_project.golden.txt` |
| `config_override` | `tests/fixtures/config_override` | Validasi penonaktifan rule via `charites.yaml` | `config_override.golden.json`<br>`config_override.golden.txt` |
| `ignore_patterns` | `tests/fixtures/ignore_patterns` | Validasi `.charitesignore` dan inline comment ignore | `ignore_patterns.golden.json`<br>`ignore_patterns.golden.txt` |

### 1.2. Prosedur Uji Coba Snapshot
1. Runner mengeksekusi pipeline pemindaian dan menangkap output `stdout` ke buffer memori.
2. Membandingkan byte aktual dengan berkas `.golden.*`.
3. Jika terdapat diskrepansi:
   - Pengujian gagal (*FAIL*).
   - Menampilkan visual perbandingan diff baris per baris (*unified diff*).
4. Pembaruan snapshot hanya diizinkan via flag `-update`:
   ```bash
   go test -v ./tests/... -run TestPipeline_GoldenSnapshots -update
   ```

---

## 2. Protokol Eksekusi Native Fuzzing (`tests/fuzz/`)

Fuzzing bertujuan menemukan cacat tersembunyi pada logika parser dan IR builder sebelum binary dilepas ke publik:

### 2.1. Target Fuzzing
```bash
# 1. Fuzzing Pipeline Astro
go test -fuzz=FuzzAstroPipeline -fuzztime=60s ./tests/fuzz/...

# 2. Fuzzing Pipeline TSX
go test -fuzz=FuzzTSXPipeline -fuzztime=60s ./tests/fuzz/...
```

### 2.2. Kriteria Kelulusan Fuzzing
- Tidak terjadi *panic: runtime error* (seperti slice out of range atau nil pointer dereference).
- Tidak terjadi *fatal error: stack overflow* akibat kedalaman nesting HTML/JSX yang ekstrem.
- Memori Go runtime tidak mengalami pembengkakan liar (*memory leak / OOM*).

---

## 3. Benchmark Pipa Penuh (End-to-End Latency)

Untuk membuktikan klaim performa sub-100ms pada monorepo:

```go
func BenchmarkFullPipeline_Monorepo(b *testing.B) {
    fixtureDir := "tests/fixtures/large_mock_repo" // Memuat 1.000 file

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var buf bytes.Buffer
        code := cli.ExecuteWithBuffer([]string{"scan", fixtureDir, "-f", "json"}, &buf)
        if code != 0 && code != 1 {
            b.Fatalf("unexpected exit code: %d", code)
        }
    }
}
```

### Target Ambang Batas:
- Memindai 1.000 berkas frontend di SSD selesai dalam waktu $\le 80\text{ milidetik}$.
