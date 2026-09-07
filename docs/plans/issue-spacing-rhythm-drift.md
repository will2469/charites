# Implementation Plan: `ux.spacing-rhythm-drift` (Local Spacing Consistency & Rhythm Invariant Engine)

Implements the formal expansion specification **[`docs/01-spec/expansion/ux-spacing-rhythm-drift.md`](../01-spec/expansion/ux-spacing-rhythm-drift.md)** (`SPEC-EXP-UX-SPACING-RHYTHM-DRIFT`).

---

## 1. Executive Summary & Problem Formulation

In modern design systems and frontend component authoring, spacing sequences establish the visual cadence and rhythm of user interfaces. A common subtle defect occurs when a layout container contains repeated sibling relationships with homogeneous spacing, but one relationship unexpectedly deviates (e.g. `4, 4, 7, 4` or `mb-4, mb-4, mb-6, mb-4`).

Existing linters (Stylelint, ESLint) only check scalar tokens or global modulus. In Charites:
- `ux.spacing-inversion` checks parent-child hierarchical spacing (intra-child spacing $\ge$ parent gap).
- `theme.hardcode-size` checks raw pixel/arbitrary token validity.
- `ux.spacing-rhythm-drift` detects **local relational rhythm drift** within a single layout group without imposing arbitrary universal multiples.

### Core Invariant
> Spacing values serving the same spatial role within the same layout group should normally come from a coherent local spacing family. When a dominant local rhythm is established ($\ge 3$ occurrences, dominance ratio $\ge 60\%$), isolated unexplained outliers are reported as `warning`.

### Architectural Invariant Triad
```text
Classification MUST precede comparison.
Comparison MUST precede diagnostic.
Diagnostic MUST NOT perform classification.
```

---

## 2. Ownership Matrix (Separation of Concerns)

| Scenario | `theme.hardcode-size` | `ux.spacing-inversion` | `ux.spacing-rhythm-drift` | Architectural Rationale |
| :--- | :---: | :---: | :---: | :--- |
| `gap-[17px]` |  | 0 | 0 | Token validator detects non-standard token; rhythm analysis isolates non-system tokens. |
| child `gap-8` > parent `space-y-3` | 0 |  | 0 | Hierarchical inversion detects child spacing exceeding parent boundary. |
| `4, 4, 7, 4` same role & layout owner | 0 | 0 |  | Rhythm drift detects isolated outlier in a dominant local family. |
| `8` section vs `4` component | 0 | 0 | 0 | Different spatial roles partition into separate rhythm pools; no drift flagged. |
| base `gap-4`, md `gap-6` | 0 | 0 | 0 | Different responsive states partition into separate rhythm pools; no drift flagged. |
| `4, 4, 6, 6` (tie) | 0 | 0 | 0 | No clear dominant rhythm established ($50\% < 60\%$ threshold, tie policy). |
| `margin` + parent `gap-4` | 0 | 0 | 0 | Competing parent layout invalidates margin ownership (`LayoutOwnerNone`). |
| `gap-4` + `gap-7` peer containers | 0 | 0 |  | Peer container layout detects outlier gap among equivalent structural peers. |

---

## 3. Locked Architectural & Semantic Contracts

1. **Typed `LayoutOwnerKind`:**
   Layout ownership is explicitly classified into `LayoutOwnerContainerGap`, `LayoutOwnerContainerSpace`, `LayoutOwnerMarginSequence`, `LayoutOwnerPeerContainers`, or `LayoutOwnerNone`.
2. **N-1 Sibling Relationships Model:**
   A spacing sequence represents relationships between consecutive peers. For $N$ items, there are $N-1$ sibling intervals:
   - For container `gap-*`/`space-*`, $N$ children yield $N-1$ spatial relationships.
   - For margin sequences, `child[i]` only represents a relationship if `child[i+1]` exists; trailing margins on the last child (`nextChild == nil`) are strictly excluded.
3. **Structural Comparability First:**
   Sibling peers under a layout owner default to `SpacingRoleGroup`. Semantic classification acts as an intentional boundary refinement (`SpacingRoleComponent` for label/input pairs, `SpacingRoleSection` for distinct macro landmarks), preventing lexical lockout on custom design system components (`<Widget />`, `<Row />`, etc.).
4. **Canonical Variant Infrastructure:**
   Reuses Charites `StripVariants` engine and general responsive variant matching (`sm`, `md`, `lg`, `xl`, `2xl`, `max-*`, `@*`) without hardcoded or redundant breakpoint tier parsers.
5. **Peer Container Gap Evaluation:**
   When homogeneous structural peers under a common parent declare gaps on the same axis (e.g. 4 peer rows with `gap-4, gap-4, gap-7, gap-4`), the parent evaluates the peer rhythm and emits exactly 1 diagnostic on the deviant peer container.
6. **Pure KISS Majority Outlier Evaluation:**
   `MinOccurrences = 3`, `MinDominance = 0.60`. Invariant: Dominant established ($\text{count} \ge 2$, $\text{ratio} \ge 0.60$) $+$ minority value $\ne$ dominant $\implies$ outlier.
