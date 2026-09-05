package analyzer_test

import (
	"testing"

	"github.com/will2469/charites/internal/analyzer"
	"github.com/will2469/charites/internal/ir"
)

func TestSortDiagnostics_TotalOrdering(t *testing.T) {
	diags := []ir.Diagnostic{
		{File: "a.tsx", Line: 10, Column: 5, Rule: "theme.b", Severity: ir.SeverityWarn, Message: "msg B", Hint: "hint B"},
		{File: "a.tsx", Line: 10, Column: 5, Rule: "theme.a", Severity: ir.SeverityError, Message: "msg A", Hint: "hint A"},
		{File: "a.tsx", Line: 10, Column: 5, Rule: "theme.a", Severity: ir.SeverityError, Message: "msg A2", Hint: "hint A"},
		{File: "a.tsx", Line: 5, Column: 1, Rule: "theme.c", Severity: ir.SeverityInfo, Message: "msg C", Hint: "hint C"},
		{File: "b.tsx", Line: 1, Column: 1, Rule: "theme.a", Severity: ir.SeverityError, Message: "msg BFile", Hint: "hint BFile"},
		// Duplikat identik untuk menguji deduplikasi idempoten
		{File: "a.tsx", Line: 10, Column: 5, Rule: "theme.a", Severity: ir.SeverityError, Message: "msg A", Hint: "hint A"},
	}

	sorted := analyzer.SortDiagnostics(diags)

	// Verifikasi deduplikasi: 6 input dengan 1 duplikat -> 5 elemen unik
	if len(sorted) != 5 {
		t.Fatalf("expected 5 sorted diagnostics after deduplication, got %d", len(sorted))
	}

	// 1. File a.tsx line 5 harus berada paling depan
	if sorted[0].File != "a.tsx" || sorted[0].Line != 5 {
		t.Errorf("expected first element to be a.tsx:5, got %s:%d", sorted[0].File, sorted[0].Line)
	}

	// 2. Pada baris 10 kolom 5: theme.a (Error) harus mendahului theme.b (Warn)
	if sorted[1].Rule != "theme.a" || sorted[1].Message != "msg A" {
		t.Errorf("expected sorted[1] to be theme.a msg A, got %+v", sorted[1])
	}
	if sorted[2].Rule != "theme.a" || sorted[2].Message != "msg A2" {
		t.Errorf("expected sorted[2] to be theme.a msg A2, got %+v", sorted[2])
	}
	if sorted[3].Rule != "theme.b" {
		t.Errorf("expected sorted[3] to be theme.b, got %+v", sorted[3])
	}

	// 3. File b.tsx harus berada paling akhir
	if sorted[4].File != "b.tsx" {
		t.Errorf("expected last element to be b.tsx, got %s", sorted[4].File)
	}
}
