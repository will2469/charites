# Specification: File-First Output Contract for JSON and Markdown Reporters

- **Status:** APPROVED / READY-TO-CODE
- **Document ID:** `SPEC-EXP-CLI-FILE-FIRST-OUTPUT`
- **Related Documents:**
  - Canonical Spec: [`docs/01-spec/05-cli.md`](file:///home/will/Monorepo/charites/docs/01-spec/05-cli.md)
  - Core Architecture: [`docs/02-architecture/05-cli.md`](file:///home/will/Monorepo/charites/docs/02-architecture/05-cli.md)
  - Core Contract: [`docs/00-CONTRACT.md`](https://github.com/will2469/charites/blob/main/docs/00-CONTRACT.md)

---

## 1. Technical Rationale & Invariant Principle

### 1.1. Invariant Formulation
> **Primary File-First Invariant:**
> A static analysis audit report MUST provide a file-primary projection in which all diagnostics belonging to the same normalized file path are contiguous and deterministically ordered.

### 1.2. Operational Rationale
In earlier iterations, findings were grouped by **Rule ID** (`Rule -> [FileA, FileB]`). As Charites scales to 35+ static rules, grouping by rule causes severe operational friction:
1. **Context Fragmentation:** An agent or developer remediating a specific component (`fileA.tsx`) must parse scattered sections throughout the audit report.
2. **Repetitive Edit Cycles:** Fixing one rule at a time forces repeated file touches, multiple parser passes, and excessive build/lint/test iterations.
3. **Atomic File Remediation:** Under the file-first hierarchy, each file is a self-contained, atomic remediation block. Developers and agents open a file once, address all diagnostics chronologically from top to bottom, and proceed to the next file without hopping between directories.

---

## 2. Shared Canonical Pipeline Architecture

To prevent semantic divergence between reporters, neither the JSON nor Markdown reporter performs ad-hoc sorting or grouping. Both reporters consume a **Single Canonical Diagnostic View**:

```text
                    ScanResult
                        │
                        ▼
              Canonical Diagnostic Set
                        │
              ┌─────────┴─────────┐
              │                   │
              ▼                   ▼
       Canonical Ordering    Aggregate Facts
              │                   │
              └─────────┬─────────┘
                        ▼
                File Group Builder
                        │
             ┌──────────┴──────────┐
             ▼                     ▼
       JSON Projection       Markdown Projection
             │                     │
             ▼                     ▼
        files[] +             Results by File
        diagnostics[]             │
                                  ├─ file header (count)
                                  ├─ line anchor
                                  ├─ rule/wiki link
                                  ├─ message
                                  ├─ hint
                                  └─ suppression
```

### Invariants Matrix

| Invariant | Definition | Verification Method |
| :--- | :--- | :--- |
| **I1: File Lexicographical Order** | Files in `files[]` and Markdown sections MUST be sorted in ascending lexicographical order by normalized POSIX relative path. | `TestReporter_Invariants/I1_FileOrdering` |
| **I2: Canonical Diagnostic Order** | Within each file, diagnostics MUST follow strict total ordering: $\text{Line ASC} \to \text{Col ASC} \to \text{Rule ASC} \to \text{Severity ASC} \to \text{Message ASC} \to \text{Hint ASC}$. | `TestReporter_Invariants/I2_ViolationOrdering` |
| **I3: Flatten Consistency** | $\text{Flatten}(\text{files}[].\text{violations}) \equiv \text{diagnostics}$ under canonical ordering. No diagnostic may be omitted, altered, or duplicated. | `TestReporter_Invariants/I3_FlattenConsistency` |
| **I4: File Total Issues** | $\text{file}.\text{total\_issues} == \text{len}(\text{file}.\text{violations})$. | `TestReporter_Invariants/I4_FileTotalIssues` |
| **I5: File Severity Counts** | $\text{file}.\text{severity\_count} == \sum \text{violations with that severity}$. | `TestReporter_Invariants/I5_FileSeverityCounts` |
| **I6: Global Severity Counts** | $\text{summary}.\text{severity\_count} == \sum \text{file}.\text{severity\_count}$. | `TestReporter_Invariants/I6_GlobalSeverityCounts` |
| **I7: Files With Issues** | $\text{summary}.\text{files\_with\_issues} == \text{len}(\text{files})$. | `TestReporter_Invariants/I7_FilesWithIssuesCount` |
| **I8: Clean Files** | $\text{summary}.\text{clean\_files} == \text{summary}.\text{scanned\_files} - \text{summary}.\text{files\_with\_issues}$. | `TestReporter_Invariants/I8_CleanFilesCount` |
| **I9: Determinism** | Consecutive rendering passes over the same `ScanResult` MUST produce byte-for-byte identical output. | `TestReporter_Invariants/I9_Determinism` |
| **I10: Semantic Non-Alteration** | Presentation layer MUST NOT alter, filter, or fabricate diagnostic values. | `TestReporter_Invariants/I10_SemanticIntegrity` |

---

## 3. Data Transfer Models & Single Source of Truth

### 3.1. SSOT Report Version
```go
const DefaultReportVersion = "1.1.0"
```
The report version is anchored at the package level and exposed consistently across JSON and Markdown outputs.

### 3.2. Extended `ScanSummary`
```go
type ScanSummary struct {
	ScannedFiles    int   `json:"scanned_files"`
	FilesWithIssues int   `json:"files_with_issues"`
	CleanFiles      int   `json:"clean_files"`
	DurationMS      int64 `json:"duration_ms"`
	ErrorCount      int   `json:"error_count"`
	WarningCount    int   `json:"warning_count"`
	InfoCount       int   `json:"info_count"`
	TotalIssues     int   `json:"total_issues"`
	Passed          bool  `json:"passed"`
}
```

### 3.3. Centralized Suppression Directive Contract
Suppression comment syntax is governed by a shared resolver rather than formatting heuristics:
```go
// SuppressionDirective resolves the valid inline suppression comment syntax based on target file extension.
func SuppressionDirective(filePath string, ruleID string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".astro", ".html":
		return fmt.Sprintf("<!-- charites:ignore %s <reason> -->", ruleID)
	case ".css":
		return fmt.Sprintf("/* charites:ignore %s <reason> */", ruleID)
	default: // .tsx, .jsx, .ts, .js
		return fmt.Sprintf("// charites:ignore %s <reason>", ruleID)
	}
}
```

---

## 4. JSON Output Contract (`--format=json`)

### 4.1. Concrete JSON Document
```json
{
  "version": "1.1.0",
  "summary": {
    "scanned_files": 28,
    "files_with_issues": 2,
    "clean_files": 26,
    "duration_ms": 18,
    "error_count": 1,
    "warning_count": 1,
    "info_count": 0,
    "total_issues": 2,
    "passed": false
  },
  "files": [
    {
      "file": "src/components/Card.tsx",
      "error_count": 0,
      "warning_count": 1,
      "info_count": 0,
      "total_issues": 1,
      "violations": [
        {
          "line": 42,
          "column": 12,
          "rule": "theme.hardcode-color",
          "category": "theme",
          "severity": "warn",
          "message": "Hardcode hex color: \"#2563eb\"",
          "hint": "Use semantic token \"bg-primary\".",
          "doc_url": "https://github.com/will2469/charites/wiki/theme.hardcode-color",
          "suppression": "// charites:ignore theme.hardcode-color <reason>"
        }
      ]
    },
    {
      "file": "src/pages/index.astro",
      "error_count": 1,
      "warning_count": 0,
      "info_count": 0,
      "total_issues": 1,
      "violations": [
        {
          "line": 14,
          "column": 8,
          "rule": "theme.hardcode-opacity-color",
          "category": "theme",
          "severity": "error",
          "message": "Hardcode opacity color: \"bg-primary/10\"",
          "hint": "Use semantic token \"primary-light\".",
          "doc_url": "https://github.com/will2469/charites/wiki/theme.hardcode-opacity-color",
          "suppression": "<!-- charites:ignore theme.hardcode-opacity-color <reason> -->"
        }
      ]
    }
  ],
  "diagnostics": [
    {
      "file": "src/components/Card.tsx",
      "line": 42,
      "column": 12,
      "rule": "theme.hardcode-color",
      "category": "theme",
      "severity": "warn",
      "message": "Hardcode hex color: \"#2563eb\"",
      "hint": "Use semantic token \"bg-primary\".",
      "doc_url": "https://github.com/will2469/charites/wiki/theme.hardcode-color"
    },
    {
      "file": "src/pages/index.astro",
      "line": 14,
      "column": 8,
      "rule": "theme.hardcode-opacity-color",
      "category": "theme",
      "severity": "error",
      "message": "Hardcode opacity color: \"bg-primary/10\"",
      "hint": "Use semantic token \"primary-light\".",
      "doc_url": "https://github.com/will2469/charites/wiki/theme.hardcode-opacity-color"
    }
  ]
}
```

---

## 5. Markdown Output Contract (`--format=md`, `--format=markdown`)

### 5.1. Exact Contract Structure
The audit report section heading is strictly `## Results by File` (no alias).

```markdown
# Charites Frontend Static Analysis & UI Ergonomics Audit Report

**Timestamp:** 2026-09-07T12:00:00.000Z
**Status:** FAILED (Violations Found)

## Summary

| Metric | Jumlah |
| :--- | :--- |
| Total Berkas | 28 |
| Berkas Bermasalah | 2 |
| Berkas Bersih | 26 |
| Durasi Pemindaian | 18ms |
| Rules Attached | 35 |
| Total Issues | 2 |
| Errors | 1 |
| Warnings | 1 |
| Info | 0 |

## Detailed Info

| ID | Category | Description | Issues Found | Status |
| :--- | :--- | :--- | :---: | :---: |
| [theme.hardcode-color](https://github.com/will2469/charites/wiki/theme.hardcode-color) | theme | Detects hardcoded hex/rgb colors... | 1 | FAILED |
| [theme.hardcode-opacity-color](https://github.com/will2469/charites/wiki/theme.hardcode-opacity-color) | theme | Detects hardcoded slash opacity... | 1 | FAILED |

## Results by File

Found 2 violations across 2 files:

### [src/components/Card.tsx](file:///workspace/project/src/components/Card.tsx#L42) (1 issue: 0 errors, 1 warning)

- **[L42:C12](file:///workspace/project/src/components/Card.tsx#L42)** • `[WARN]` • [`theme.hardcode-color`](https://github.com/will2469/charites/wiki/theme.hardcode-color)
  - **Message:** Hardcode hex color: "#2563eb"
  - **Hint:** Use semantic token "bg-primary".
  - **Suppression:** `// charites:ignore theme.hardcode-color <reason>`

---

### [src/pages/index.astro](file:///workspace/project/src/pages/index.astro#L14) (1 issue: 1 error, 0 warnings)

- **[L14:C8](file:///workspace/project/src/pages/index.astro#L14)** • `[ERROR]` • [`theme.hardcode-opacity-color`](https://github.com/will2469/charites/wiki/theme.hardcode-opacity-color)
  - **Message:** Hardcode opacity color: "bg-primary/10"
  - **Hint:** Use semantic token "primary-light".
  - **Suppression:** `<!-- charites:ignore theme.hardcode-opacity-color <reason> -->`
```
