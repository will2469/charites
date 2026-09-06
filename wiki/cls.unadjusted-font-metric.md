# cls.unadjusted-font-metric

> **Rule ID:** `cls.unadjusted-font-metric`
> **Severity:** `INFO`
> **Category:** `cls`
> **Target Standards:** W3C CSS Fonts Module Level 4 (size-adjust, ascent-override, descent-override), Google Chrome Font Metric Override Guidelines

---

## 1. Overview & Core Invariant

Recommends font metric overrides on fallback font declarations to reduce swap CLS

### Core Invariant:
> **"Local fallback @font-face declarations (using 'src: local(...)') should specify metric adjustment descriptors ('size-adjust', 'ascent-override', or 'descent-override') to align bounding boxes with the principal web font."**

---
## 2. Technical Grounding & Engine Realities

When a web font downloads and replaces a system fallback font, variations in glyph x-height, ascent, and descent alter the computed bounding boxes of every text line.

This disparity causes sudden vertical expansion or contraction of paragraphs and headers, contributing to Cumulative Layout Shift.

By declaring 'size-adjust', 'ascent-override', and 'descent-override' on the fallback @font-face, developers can calibrate the system font's metrics to match the web font, creating a near zero-shift swap.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Font Swap Layout Jitter** | LOW | Paragraphs and navigation bars visibly shift lines when the web font swaps in. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Local fallback font-face without metric override descriptors):
```astro
<style>
  @font-face {
    font-family: 'InterFallback';
    src: local('Arial');
  }
</style>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Local fallback font-face with size-adjust and ascent-override):
```astro
<style>
  @font-face {
    font-family: 'InterFallback';
    src: local('Arial');
    ascent-override: 90%;
    descent-override: 22%;
    size-adjust: 107%;
  }
</style>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: cls.unadjusted-font-metric"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore cls.unadjusted-font-metric` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/cls.unadjusted-font-metric/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for cls.unadjusted-font-metric"]
        subgraph P ["Positive Corpus (tests/correctness/cls.unadjusted-font-metric/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/cls.unadjusted-font-metric/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/cls.unadjusted-font-metric/adversarial/)"]
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
<!-- charites:ignore cls.unadjusted-font-metric intentional exception -->
```

```tsx
// charites:ignore cls.unadjusted-font-metric intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.unadjusted-font-metric:
    severity: info # error | warn | info | off
```

