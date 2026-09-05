# 06-ROADMAP: 01 - Phase 1 Milestone & Transition Gate

> **Kode Dokumen:** `ROAD-01-CONTRACT`
> **Tahapan:** Fase 1 - Kunci Kontrak Data (IR & Diagnostic)
> **Status:** Ready for Review

Dokumen ini menetapkan kriteria kelulusan (*exit criteria*) dan gerbang transisi (*transition gate*) untuk **Fase 1 (Kunci Kontrak Data)** sebelum tim melangkah ke **Fase 2 (Parser Frontend & IR Builder)**.

---

## 1. Deliverables Wajib Fase 1

1. **`internal/ir/node.go`**: Definisi struct `Node`, `Span`, enum `NodeType`, dan iterator `Walk()` berbasis `iter.Seq`.
2. **`internal/ir/diagnostic.go`**: Definisi struct `Diagnostic`, enum `Severity`, dan JSON serialization tags.
3. **`internal/ir/node_test.go`**: Pengujian unit traversal `Walk()`, terminasi awal, dan helper methods (`HasClass`, `GetAttr`, `IsElement`).
4. **`internal/ir/diagnostic_test.go`**: Pengujian determinisme serialisasi JSON `Diagnostic`.
5. **`internal/rules/rule.go`**: Definisi interface baku `rules.Rule` yang mengonsumsi `*ir.Node` dan menghasilkan `[]ir.Diagnostic`.

---

## 2. Checklist Definition of Done (DoD) Fase 1

- [ ] Seluruh berkas deliverable Fase 1 terimplementasi tanpa paket eksternal pihak ketiga.
- [ ] `go test -v -cover ./internal/ir/...` menghasilkan coverage $\ge 90\%$ tanpa kegagalan.
- [ ] Benchmark `BenchmarkNode_Walk` membuktikan **0 B/op** dan **0 allocs/op**.
- [ ] Kompilasi membuktikan bebas dari dependensi sirkular (`import cycle`).
- [ ] `golangci-lint run ./internal/ir/...` lolos 100% tanpa warning.

---

## 3. Gerbang Transisi ke Fase 2 (Parser & IR Builder)

Begitu seluruh checklist di atas terpenuhi:
1. Buat git commit: `feat(ir): define unified node, diagnostic contracts and zero-alloc walk iterator`.
2. Buka dokumen `docs/01-spec/02-parser.md` untuk memulai spesifikasi ekstraksi AST Tailwind CSS, Astro, dan TSX.
