# responsive.mobile-density-overload

> **Rule ID:** `responsive.mobile-density-overload`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** Steven Hoober (Designing for Touch - Touch Target Interference), WCAG 2.2 SC 2.5.8 (Target Size - Minimum & Spacing), Material Design 3 Mobile App Bar & Toolbar Guidelines

---

## 1. Overview & Core Invariant

Warns when toolbars or action rows cram more than 4 interactive buttons in a single unscrollable row on mobile viewports

### Core Invariant:
> **"Horizontal action toolbars on mobile viewports must not cram more than 4 interactive buttons in a single rigid row without 'overflow-x-auto', 'flex-wrap', or an overflow menu."**

---
## 2. Technical Grounding & Engine Realities

On compact smartphone screens (360px viewport width), accommodating 5 or more buttons in a single unscrollable flex row forces button widths below 48px or induces layout squishing.

This severe spatial compression leads to:
1. High Error Rate / Mis-taps: Users inadvertently trigger adjacent destructive or unwanted actions due to finger pad overlap.
2. Text/Icon Clipping: Labels are aggressively truncated, and icon hitboxes overlap.

Best practice dictates limiting direct actions to 3-4 primary controls, wrapping the toolbar in a horizontal scroll container ('overflow-x-auto'), or collapsing secondary actions into a 'More (...)' dropdown menu.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Touch Target Mis-tap Interference** | MEDIUM | Users frequently tap the wrong button because adjacent targets are compressed below safe physical spacing limits. |
| **Mobile Visual Clutter & Overflow** | LOW | Rigid toolbars cause horizontal viewport tearing or text clipping on narrow mobile devices. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Five buttons packed tightly into an unscrollable horizontal flex row):
```tsx
<div className="flex items-center gap-2 p-2">
  <button type="button">Edit</button>
  <button type="button">Salin</button>
  <button type="button">Cetak</button>
  <button type="button">Bagikan</button>
  <button type="button">Hapus</button>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Scrollable horizontal action bar accommodating many actions comfortably):
```tsx
<div className="flex items-center gap-2 p-2 overflow-x-auto">
  <button type="button">Edit</button>
  <button type="button">Salin</button>
  <button type="button">Cetak</button>
  <button type="button">Bagikan</button>
  <button type="button">Hapus</button>
</div>
```

---

## 6. Detection & Verification Pipeline (How The Rule Evaluates Code)
This rule evaluates source code through the standard AST inspection pipeline:

```mermaid
flowchart TD
    Node["AST Node (Astro / TSX element)"] --> Inspect["1. Inspect Element & Attributes"]
    Inspect --> Invariant{"2. Evaluate Rule Invariant"}
    Invariant -- "Compliant" --> Safe["Pass"]
    Invariant -- "Non-Compliant" --> IgnoreCheck{"3. Check charites:ignore directive"}
    IgnoreCheck -- "Ignored" --> Safe
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: responsive.mobile-density-overload"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore responsive.mobile-density-overload` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/responsive.mobile-density-overload/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for responsive.mobile-density-overload"]
        subgraph P ["Positive Corpus (tests/correctness/responsive.mobile-density-overload/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/responsive.mobile-density-overload/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/responsive.mobile-density-overload/adversarial/)"]
            A1["A1: Template Literal Interpolations"]
            A2["A2: Ternary Conditional Expressions"]
            A3["A3: Spread Properties & Dynamic Overrides"]
            A4["A4: Dynamic Object Class Syntax"]
            A5["A5: Shadowed Variable Identifiers"]
            A6["A6: Nested Closures & HOC Wrappers"]
            A7["A7: Obfuscated Classes & Cyclic Tokens"]
        end
    end

    P --> TestRunner["Automated Runner (rule_test.go)"]
    N --> TestRunner
    A --> TestRunner
    TestRunner --> Gates["Quality Gates: Zero Panic, Zero False-Positive, Zero Bypass"]
```

- **Positive Fixtures (P1-P5):** Verified to trigger diagnostics at exact lines and column spans.
- **Negative Fixtures (N1-N5):** Verified to produce zero diagnostics on valid tokens and legitimate exemptions.
- **Adversarial Fixtures (A1-A7):** Verified to prevent evasion across dynamic expressions, string interpolations, and cyclic references.

---

## 8. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore responsive.mobile-density-overload intentional exception -->
```

```tsx
// charites:ignore responsive.mobile-density-overload intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.mobile-density-overload:
    severity: warn # error | warn | info | off
```

