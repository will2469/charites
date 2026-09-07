package ux_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules/ux"
)

func TestMultilineInputMisuseRule(t *testing.T) {
	rule := ux.NewMultilineInputMisuseRule()

	cases := []struct {
		name      string
		code      string
		wantDiags int
	}{
		// Positive Cases
		{
			name:      "Astrades specimen: id=keterangan-input and placeholder=Catatan tambahan",
			code:      `<Input id="keterangan-input" placeholder="Catatan tambahan (bila ada)" />`,
			wantDiags: 1,
		},
		{
			name:      "Rejection reason on native input",
			code:      `<input name="rejection_reason" />`,
			wantDiags: 1,
		},
		{
			name:      "Alasan pembatalan on Input component",
			code:      `<Input name="alasan_pembatalan" />`,
			wantDiags: 1,
		},
		{
			name:      "Notes on Input component",
			code:      `<Input name="notes" />`,
			wantDiags: 1,
		},
		{
			name:      "Description on Input component",
			code:      `<Input name="description" />`,
			wantDiags: 1,
		},
		{
			name:      "Multi-channel resolution: placeholder is multiline while name is unknown",
			code:      `<Input name="user_name" placeholder="Catatan tambahan" />`,
			wantDiags: 1,
		},
		{
			name:      "Aria-label channel is multiline while name is generic",
			code:      `<Input name="foo" aria-label="Catatan tambahan" />`,
			wantDiags: 1,
		},

		// Negative Cases
		{
			name:      "Already Textarea component",
			code:      `<Textarea name="notes" rows={3} />`,
			wantDiags: 0,
		},
		{
			name:      "Already native textarea element",
			code:      `<textarea name="keterangan" rows={3}></textarea>`,
			wantDiags: 0,
		},
		{
			name:      "Single-line safeguard: note_title",
			code:      `<Input name="note_title" />`,
			wantDiags: 0,
		},
		{
			name:      "Single-line safeguard: reason_code",
			code:      `<Input name="reason_code" />`,
			wantDiags: 0,
		},
		{
			name:      "Single-line safeguard: short_description",
			code:      `<Input name="short_description" />`,
			wantDiags: 0,
		},
		{
			name:      "Qualifier override: description name with Short description placeholder",
			code:      `<Input name="description" placeholder="Short description" />`,
			wantDiags: 0,
		},
		{
			name:      "Quantity safeguard: notes_count",
			code:      `<Input name="notes_count" />`,
			wantDiags: 0,
		},
		{
			name:      "Date safeguard: catatan_tanggal",
			code:      `<Input name="catatan_tanggal" />`,
			wantDiags: 0,
		},

		// Adversarial Cases
		{
			name:      "Non-text type: password with name=notes",
			code:      `<input type="password" name="notes" />`,
			wantDiags: 0,
		},
		{
			name:      "Non-text type: number with name=notes",
			code:      `<input type="number" name="notes" />`,
			wantDiags: 0,
		},
		{
			name:      "Dynamic type expression: type={dynamicType}",
			code:      `<input type={dynamicType} name="notes" />`,
			wantDiags: 0,
		},
		{
			name:      "Abbreviated desc (omitted from strong tokens): sort_desc",
			code:      `<input name="sort_desc" />`,
			wantDiags: 0,
		},
		{
			name:      "Data-testid excluded from v1: data-testid=notes-field with generic name",
			code:      `<input data-testid="notes-field" name="user_name" />`,
			wantDiags: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `export function TestComponent() { return (` + tc.code + `); }`
			root, err := tsx.Extract([]byte(src))
			if err != nil {
				t.Fatalf("failed to parse TSX: %v", err)
			}

			var diagsCount int
			for node := range root.Walk() {
				findings := rule.Evaluate(node)
				diagsCount += len(findings)
			}

			if diagsCount != tc.wantDiags {
				t.Errorf("[%s] got %d diagnostics, want %d", tc.name, diagsCount, tc.wantDiags)
			}
		})
	}
}
