# responsive.missing-breakpoint

> **Rule ID:** `responsive.missing-breakpoint`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** Mobile-First Responsive Web Design Principles, W3C CSS Grid Layout Module Level 2, Tailwind CSS Responsive Design Specification

---

## 1. Overview & Core Invariant

Warns when multi-column grids or giant font sizes are declared on mobile baseline without responsive breakpoint modifiers

### Core Invariant:
> **"Multi-column grids (grid-cols-[3-12]) and giant font sizes (text-[5-9]xl) must not be defined on mobile baseline without responsive breakpoint prefixes (sm:, md:, lg:)."**

---
## 2. Technical Grounding & Engine Realities

On compact smartphone screens (360px-390px), defining 3 or more columns directly on mobile baseline squeezes individual columns below 100px width, causing severe card distortion and text wrapping.

Similarly, declaring giant typography (e.g. text-6xl) on mobile baseline causes single words to span multiple vertical lines, breaking header visual balance.

Adhering to mobile-first progression requires starting from single-column baselines (grid-cols-1) and scaling up via responsive modifiers (sm:grid-cols-2 md:grid-cols-4).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Severe Column Squeeze on Mobile** | MEDIUM | Multi-column cards become unreadable and distorted when squeezed into 360px phone screens. |
| **Typography Layout Blowout** | LOW | Giant font headings wrap awkwardly into 4-5 lines on narrow mobile viewports. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Multi-column grid on mobile baseline without responsive modifier):
```tsx
<div className="grid grid-cols-4 gap-4">
  <div className="bg-card p-4">Item 1</div>
  <div className="bg-card p-4">Item 2</div>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Mobile-first progression starting from 1 column to multi-column on desktop):
```tsx
<div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
  <div className="bg-card p-4">Item 1</div>
  <div className="bg-card p-4">Item 2</div>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: responsive.missing-breakpoint"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore responsive.missing-breakpoint` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/responsive.missing-breakpoint/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for responsive.missing-breakpoint"]
        subgraph P ["Positive Corpus (tests/correctness/responsive.missing-breakpoint/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/responsive.missing-breakpoint/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/responsive.missing-breakpoint/adversarial/)"]
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
<!-- charites:ignore responsive.missing-breakpoint intentional exception -->
```

```tsx
// charites:ignore responsive.missing-breakpoint intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.missing-breakpoint:
    severity: warn # error | warn | info | off
```

