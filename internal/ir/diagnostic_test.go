package ir_test

import (
	"bytes"
	"encoding/json"
	"slices"
	"testing"

	"github.com/will2469/charites/internal/ir"
)

func TestDiagnostic_JSONDeterminism(t *testing.T) {
	diagWithHint := ir.Diagnostic{
		File:     "components/Button.tsx",
		Line:     12,
		Column:   5,
		Rule:     "theme.hardcode-opacity-color",
		Severity: ir.SeverityError,
		Message:  "Opacity color class found",
		Hint:     "Use design token instead",
	}

	data1, err := json.Marshal(diagWithHint)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	// Membuktikan byte-level determinism: marshaling ulang menghasilkan byte identik
	data2, err := json.Marshal(diagWithHint)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	if !bytes.Equal(data1, data2) {
		t.Fatalf("expected byte-identical JSON output across calls")
	}

	// Membuktikan flat JSON keys
	var rawMap map[string]any
	if err = json.Unmarshal(data1, &rawMap); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	expectedKeys := []string{"file", "line", "column", "rule", "severity", "message", "hint"}
	for _, key := range expectedKeys {
		if _, exists := rawMap[key]; !exists {
			t.Errorf("expected flat JSON key %q to exist in output", key)
		}
	}

	// Membuktikan omitempty pada Hint kosong
	diagWithoutHint := ir.Diagnostic{
		File:     "components/Card.astro",
		Line:     1,
		Column:   1,
		Rule:     "a11y.missing-alt",
		Severity: ir.SeverityWarn,
		Message:  "Missing alt text",
	}

	dataWithoutHint, err := json.Marshal(diagWithoutHint)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var rawMapNoHint map[string]any
	if err = json.Unmarshal(dataWithoutHint, &rawMapNoHint); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if _, exists := rawMapNoHint["hint"]; exists {
		t.Errorf("expected 'hint' to be omitted when empty string")
	}

	// Membuktikan serialisasi slice deterministik
	list := []ir.Diagnostic{diagWithHint, diagWithoutHint}
	listJSON1, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("failed to marshal list: %v", err)
	}
	listJSON2, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("failed to marshal list second time: %v", err)
	}
	if !bytes.Equal(listJSON1, listJSON2) {
		t.Fatalf("expected byte-identical slice JSON marshaling")
	}
}

func TestDiagnostic_CollectionOrdering(t *testing.T) {
	// 1. Uji 7-Level Strict Tie-Breaking
	t.Run("7-Level Tie Breaking", func(t *testing.T) {
		cases := []struct {
			name     string
			a        ir.Diagnostic
			b        ir.Diagnostic
			expected int // -1 jika a < b, 1 jika a > b, 0 jika a == b
		}{
			{
				name:     "Level 1: Different File",
				a:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 1, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				b:        ir.Diagnostic{File: "b.tsx", Line: 10, Column: 1, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				expected: -1,
			},
			{
				name:     "Level 2: Same File, Different Line",
				a:        ir.Diagnostic{File: "a.tsx", Line: 5, Column: 10, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				b:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 1, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				expected: -1,
			},
			{
				name:     "Level 3: Same Line, Different Column",
				a:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				b:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 8, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				expected: -1,
			},
			{
				name:     "Level 4: Same Location, Different Rule",
				a:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "a11y.foo", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				b:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "theme.bar", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				expected: -1,
			},
			{
				name:     "Level 5: Same Rule, Different Severity (error < warn < info)",
				a:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				b:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityWarn, Message: "m", Hint: "h"},
				expected: -1,
			},
			{
				name:     "Level 5b: Warn < Info",
				a:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityWarn, Message: "m", Hint: "h"},
				b:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityInfo, Message: "m", Hint: "h"},
				expected: -1,
			},
			{
				name:     "Level 5c: Info < Unknown Severity",
				a:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityInfo, Message: "m", Hint: "h"},
				b:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.Severity("unknown"), Message: "m", Hint: "h"},
				expected: -1,
			},
			{
				name:     "Level 6: Same Severity, Different Message",
				a:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "Message A", Hint: "h"},
				b:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "Message B", Hint: "h"},
				expected: -1,
			},
			{
				name:     "Level 7: Same Message, Different Hint",
				a:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "Fix A"},
				b:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "Fix B"},
				expected: -1,
			},
			{
				name:     "Absolute Identical",
				a:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				b:        ir.Diagnostic{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "m", Hint: "h"},
				expected: 0,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				res := ir.CompareDiagnostics(tc.a, tc.b)
				if res != tc.expected {
					t.Errorf("CompareDiagnostics(a, b) = %d; expected %d", res, tc.expected)
				}
				// Antisimetri: Compare(b, a) == -res
				resRev := ir.CompareDiagnostics(tc.b, tc.a)
				if resRev != -tc.expected {
					t.Errorf("antisymmetry violated: Compare(b, a) = %d; expected %d", resRev, -tc.expected)
				}
			})
		}
	})

	// 2. Uji Permutation Invariance (Shuffled Order Convergence)
	t.Run("Permutation Invariance", func(t *testing.T) {
		canonicalOrder := []ir.Diagnostic{
			{File: "a.tsx", Line: 1, Column: 1, Rule: "r1", Severity: ir.SeverityError, Message: "A", Hint: "H1"},
			{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "A", Hint: "H1"},
			{File: "a.tsx", Line: 10, Column: 5, Rule: "r1", Severity: ir.SeverityError, Message: "A", Hint: "H1"},
			{File: "a.tsx", Line: 10, Column: 5, Rule: "r2", Severity: ir.SeverityError, Message: "A", Hint: "H1"},
			{File: "a.tsx", Line: 10, Column: 5, Rule: "r2", Severity: ir.SeverityWarn, Message: "A", Hint: "H1"},
			{File: "a.tsx", Line: 10, Column: 5, Rule: "r2", Severity: ir.SeverityWarn, Message: "B", Hint: "H1"},
			{File: "a.tsx", Line: 10, Column: 5, Rule: "r2", Severity: ir.SeverityWarn, Message: "B", Hint: "H2"},
			{File: "b.astro", Line: 2, Column: 1, Rule: "r1", Severity: ir.SeverityError, Message: "A", Hint: "H1"},
		}

		// Marshalling dari urutan kanonikal
		expectedBytes, err := json.Marshal(canonicalOrder)
		if err != nil {
			t.Fatalf("failed to marshal canonical: %v", err)
		}

		// Variasi permutasi 1: Reverse order
		perm1 := slices.Clone(canonicalOrder)
		slices.Reverse(perm1)
		sorted1 := ir.SortDiagnostics(perm1)
		bytes1, err := json.Marshal(sorted1)
		if err != nil {
			t.Fatalf("failed to marshal sorted1: %v", err)
		}
		if !bytes.Equal(bytes1, expectedBytes) {
			t.Fatalf("permutation 1 failed: output bytes differ from canonical")
		}

		// Variasi permutasi 2: Interleaved swap
		perm2 := slices.Clone(canonicalOrder)
		perm2[0], perm2[4] = perm2[4], perm2[0]
		perm2[1], perm2[6] = perm2[6], perm2[1]
		perm2[2], perm2[7] = perm2[7], perm2[2]
		sorted2 := ir.SortDiagnostics(perm2)
		bytes2, err := json.Marshal(sorted2)
		if err != nil {
			t.Fatalf("failed to marshal sorted2: %v", err)
		}
		if !bytes.Equal(bytes2, expectedBytes) {
			t.Fatalf("permutation 2 failed: output bytes differ from canonical")
		}
	})

	// 3. Uji Idempotent Deduplication
	t.Run("Idempotent Deduplication", func(t *testing.T) {
		duplicated := []ir.Diagnostic{
			{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "M", Hint: "H"},
			{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "M", Hint: "H"},
			{File: "a.tsx", Line: 10, Column: 2, Rule: "r1", Severity: ir.SeverityError, Message: "M", Hint: "H"},
			{File: "b.tsx", Line: 5, Column: 1, Rule: "r2", Severity: ir.SeverityWarn, Message: "M2", Hint: "H2"},
			{File: "b.tsx", Line: 5, Column: 1, Rule: "r2", Severity: ir.SeverityWarn, Message: "M2", Hint: "H2"},
		}

		deduped := ir.SortDiagnostics(duplicated)
		if len(deduped) != 2 {
			t.Fatalf("expected 2 unique diagnostics after deduplication, got %d", len(deduped))
		}

		if deduped[0].File != "a.tsx" || deduped[1].File != "b.tsx" {
			t.Errorf("unexpected deduplication ordering: %+v", deduped)
		}
	})
}
