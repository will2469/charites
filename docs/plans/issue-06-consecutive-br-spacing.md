# Implementation Plan: Consecutive `<br>` Spacing Rule (Issue #6)

Resolves **[GitHub Issue #6](https://github.com/will2469/charites/issues/6)** (`semantic.consecutive-br-spacing`: feat(semantic): rule to detect consecutive `<br>` tags used for paragraph spacing violating WCAG 1.3.1 Technique F34).

---

## 1. Executive Summary & Problem Statement

In web and component development, developers occasionally use consecutive `<br>` tags (`<br /><br />`) to push content downward and simulate vertical paragraph separation:
```tsx
<>
  Apakah Anda yakin memilih kandidat ini sebagai Pimpinan Utama?
  <br />
  <br />
  <strong className="text-destructive">PERINGATAN:</strong> Tindakan ini tidak dapat dibatalkan.
</>
```

This pattern introduces serious accessibility and responsive design defects:
1. **Assistive Technology Degradation (WCAG 1.3.1 Technique F34):** Screen readers announce each `<br>` as a blank line or acoustic break sound, forcing screen reader users to navigate through phantom empty elements and disrupting document comprehension.
2. **Design System Token Bypassing:** Hardcoded `<br>` line breaks bypass design system spacing tokens (`space-y-*`, `gap-*`, container queries) and fail to scale with responsive typography.
3. **Semantic Markup Violation:** HTML Living Standard Section 4.5.27 specifies that `<br>` represents a line break within a single phrase or address block (e.g. postal addresses, poetry stanzas). Paragraphs must be wrapped in semantic container elements (`<p>`, `<div>`) and spaced with CSS.

---

## 2. The 4 Locked Architectural Invariants

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE 4 LOCKED ARCHITECTURAL INVARIANTS                    │
├─────────────────────────────────────────────────────────────────────────────┤
│ 1. Cardinality Invariant (Run Detector vs Pair Detector):                  │
│    A run of >= 2 consecutive <br> siblings represents ONE semantic defect. │
│    When runCount reaches exactly 2, emit exactly 1 diagnostic anchored on   │
│    the 2nd <br>. Subsequent <br> in the same run DO NOT emit duplicates:    │
│    - <br><br>          → 1 diagnostic                                       │
│    - <br><br><br>      → 1 diagnostic                                       │
│    - <br><br><br><br>  → 1 diagnostic                                       │
│                                                                             │
│ 2. Ignorable Sibling Invariant:                                             │
│    NodeComment ({/* ... */}, <!-- ... -->) and whitespace-only NodeText     │
│    are ignorable siblings: they do NOT break or reset the consecutive run.  │
│    Meaningful text or non-<br> elements strictly reset the run (runCount=0).│
│                                                                             │
│ 3. Structural Non-Flattening Invariant:                                     │
│    Inspect only direct children of NodeElement or NodeFragment.             │
│    Descendants in different child subtrees (e.g. <p><br/></p><p><br/></p>)   │
│    are NEVER flattened into the parent's sibling stream (0 diagnostics).    │
│                                                                             │
│ 4. AST Fact ≠ Human Intent (Defensible Diagnosis Message):                  │
│    Static AST analysis detects structural facts, not human developer intent.│
│    Message states the structural pattern fact without overclaiming intent:   │
│    "Consecutive '<br>' elements detected as a spacing pattern; use          │
│    structural elements and CSS for paragraph separation."                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Finite-State Run Detector State Machine

For each parent node (`NodeElement` or `NodeFragment`) independently:

```text
               ┌──────────┐
    Start ───> │ run = 0  │
               └────┬─────┘
                    │
         child == <br>
                    │
                    ▼
               ┌──────────┐   child == <br>    ┌──────────┐
               │ run = 1  │ ─────────────────> │ run = 2  │ ──> Emit Diagnostic (Line of 2nd <br>)
               └────┬─────┘                    └────┬─────┘
                    │                               │
        non-<br> / non-ignorable        child == <br>
                    │                               │
                    ▼                               ▼
               ┌──────────┐                    ┌──────────┐
               │ run = 0  │ <───────────────── │ run >= 3 │ (No additional diagnostic)
               └──────────┘    non-<br> /      └──────────┘
                              non-ignorable
```

- **Ignorable Siblings:** `child.Type == NodeComment` or `child.Type == NodeText && strings.TrimSpace(child.RawClasses) == ""` $\rightarrow$ `continue` (no state transition).
- **Run Reset:** Non-ignorable sibling $\rightarrow$ `run = 0`.

---

## 4. 1-SSOT Tri-Corpus Acceptance Matrix

| Corpus | Fixture | Expected Diagnostics | Rationale |
| :--- | :--- | :---: | :--- |
| **Positive** | `astrades_dialog.tsx` | 1 | Real-world specimen using double `<br />` inside fragment |
| **Positive** | `triple_br.astro` | 1 | Triple `<br />` in Astro template; asserts cardinality $\le 1$ |
| **Positive** | `comment_between_br.tsx` | 1 | Comments between `<br />` do not break the consecutive run |
| **Negative** | `address_single_br.tsx` | 0 | Address block with legitimate single line breaks |
| **Negative** | `poem_stanza.astro` | 0 | Poetry stanza with legitimate single line breaks |
| **Negative** | `semantic_paragraphs.tsx` | 0 | Paragraphs spaced via CSS `space-y-4` without `<br>` |
| **Adversarial** | `non_empty_text_between.tsx` | 0 | Meaningful text between line breaks resets run |
| **Adversarial** | `nested_different_parents.tsx`| 0 | `<br>` in different parent paragraphs (non-flattening) |
| **Adversarial** | `tag_containment_bait.tsx` | 0 | `<Breadcrumb />`, `<BrandedButton />` are not `<br>` |

---

## 5. Implementation Topology

- **Rule Definition:** `internal/rules/semantic/consecutive_br_spacing.go`
  - Canonical ID: `semantic.consecutive-br-spacing`
  - Category: `semantic`
  - Default Severity: `ir.SeverityWarn`
  - 8-Pillars Documentation with WCAG 2.2 Technique F34 and HTML Living Standard references
- **Unit & Contract Tests:** `internal/rules/semantic/consecutive_br_spacing_test.go` and `internal/rules/semantic/contract_test.go`
- **Registration:** `internal/rules/builtin.go` under `semantic Wave 1 (Document Structure & Content Semantics)`
- **Golden Tri-Corpus:** `tests/correctness/semantic/consecutive-br-spacing/`
- **Documentation:** Automated via `make wiki` $\rightarrow$ `wiki/semantic.consecutive-br-spacing.md`
