# 03-TESTING: 06 - Golden Snapshot Regression & Fuzzing Verification Plan

> **Kode Dokumen:** `TEST-06-GOLDEN`
> **Tahapan:** Fase 6 - Validasi Penuh & Golden Snapshots (Milestone Selesai Pipa)
> **Peran Pilar:** TEST = PROOF (Harness Pengujian, Matriks Golden & Pembuktian Stabilitas)
> **Status:** Ready for Review
> **Standar Rujukan:** Snapshot Regression Testing & Native Fuzzing Protocols

Dokumen ini mendefinisikan strategi pengujian integrasi pipa compiler (*end-to-end*), matriks **Golden Master Snapshots**, verifikasi bertingkat L1-L3, serta benchmark latensi pipeline penuh (`BENCH-06-E2E-001`).

---

## 1. Strategi Pengujian Bertingkat (L1 - L3 Verification)

Untuk mencegah kegagalan diagnosis terselubung (*false-clean masked by reporter bug*):

1. **L1: Verifikasi Struktural AST & Parser:** Memvalidasi bahwa parser frontend menghasilkan pohon `*ir.Node` yang utuh, relasi parent-child terjaga, dan line offset akurat.
2. **L2: Verifikasi Logika Evaluasi Rule & Suppression:** Memvalidasi bahwa rule mendeteksi pelanggaran sesuai matriks batas 20 skenario dan filter inline ignore bekerja tepat pada baris/span target.
3. **L3: Verifikasi Golden Master Dokumen Penuh:** Memvalidasi bahwa serialisasi keluaran akhir (CLI JSON & Inline Text) menghasilkan byte yang identik dengan berkas kebenaran mutlak (*ground truth*).

---

## 2. Matriks Skenario Golden Snapshots (`tests/golden/projects/`)

| Nama Skenario | Path Fixture | Target Validasi | Golden File Terkait |
| :--- | :--- | :--- | :--- |
| `clean` | `tests/fixtures/projects/clean` | Repositori bersih tanpa pelanggaran (*Zero Noise Invariant*) | `clean.golden.json`<br>`clean.golden.txt` |
| `opacity_violations` | `tests/fixtures/projects/opacity_violations` | Pelanggaran opacity Tailwind pada Astro & React TSX | `opacity_violations.golden.json`<br>`opacity_violations.golden.txt` |
| `config_override` | `tests/fixtures/projects/config_override` | Penonaktifan rule via `charites.yaml` (`off`) & severity override (`warn`) | `config_override.golden.json`<br>`config_override.golden.txt` |
| `ignore_patterns` | `tests/fixtures/projects/ignore_patterns` | Pengabaian path via `.charitesignore` dan inline directives | `ignore_patterns.golden.json`<br>`ignore_patterns.golden.txt` |

### 2.1. Protokol Eksekusi & Pembaruan Snapshot
1. Runner membaca fixture dan mengeksekusi pipeline pemindaian.
2. Output JSON dinormalisasi (menghilangkan metadata non-deterministik `duration_ms`).
3. Membandingkan byte aktual dengan berkas golden referensi.
4. Pembaruan snapshot lokal:
   ```bash
   go test -v ./tests/... -run TestPipeline_GoldenSnapshots -update
   ```
5. CI memverifikasi integritas snapshot dan dilarang mengaktifkan flag `-update`.

---

## 3. Protokol Eksekusi Native Fuzzing (`tests/fuzz/`)

Fuzzing dijalankan untuk memastikan pipeline kebal terhadap input byte acak malformed:

```bash
# 1. Fuzzing Pipeline Astro Terintegrasi
go test -fuzz=FuzzAstroPipeline -fuzztime=60s ./tests/fuzz/...

# 2. Fuzzing Pipeline TSX Terintegrasi
go test -fuzz=FuzzTSXPipeline -fuzztime=60s ./tests/fuzz/...
```

### Kriteria Kelulusan Fuzzing:
- Fuzz target berjalan minimal **60 detik** per modul tanpa panic runtime (*unhandled panic*), segmentation fault, atau fatal error stack overflow.

---

## 4. Benchmark Latensi Pipeline Menyeluruh (`BENCH-06-E2E-001`)

Pengujian performa pipeline terpadu diukur menggunakan korpus monorepo mock terstandarisasi:

```go
func BenchmarkFullPipeline_Monorepo(b *testing.B) {
    fixtureDir := "tests/fixtures/projects/large_mock_repo" // 1.000 file campuran Astro & TSX

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

### Metodologi & Anggaran Desain (`BENCH-06-E2E-001`):
- **Korpus:** 1.000 berkas frontend (rata-rata ukuran $2\text{ KB}$ per berkas).
- **Environment:** Go 1.26 toolchain, `CGO_ENABLED=0`, runner CPU standar.
- **Target Desain (Performance Budget):** Pemindaian 1.000 berkas diharapkan selesai dalam waktu $\le 100\text{ ms}$ pada SSD modern.
