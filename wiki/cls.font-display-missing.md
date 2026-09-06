# cls.font-display-missing

> **Rule ID:** `cls.font-display-missing`
> **Severity:** `ERROR`
> **Category:** `cls`
> **Target Standards:** W3C CSS Fonts Module Level 4 (@font-face font-display descriptor), Google Core Web Vitals Guidelines (Cumulative Layout Shift & FOUT/FOIT), Web.dev Font Best Practices

---

## 1. Overview & Core Invariant

Requires font-display descriptor on custom @font-face declarations to prevent FOIT reflow

### Core Invariant:
> **"All custom @font-face declarations must declare an explicit, valid 'font-display' descriptor ('swap', 'optional', or 'fallback') to ensure continuous text visibility and prevent layout reflow."**

---
## 2. Technical Grounding & Engine Realities

When a browser encounters a custom @font-face without a 'font-display' descriptor, it defaults to 'font-display: auto' (often identical to 'block').

Under the 'block' period, the browser hides text completely (Flash of Invisible Text / FOIT) for up to 3 seconds while waiting for the web font file. Once the font arrives, the browser suddenly swaps the font and recalculates line wrapping, triggering Cumulative Layout Shift (CLS).

Using 'font-display: swap' renders system fallback fonts immediately and swaps gracefully when the font finishes loading, ensuring accessibility and predictable rendering.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Flash of Invisible Text (FOIT)** | HIGH | Users stare at blank spaces on slow networks while waiting for fonts to load. |
| **Cumulative Layout Shift (CLS)** | HIGH | Late font swaps cause text wrapping reflow that pushes subsequent content down. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Custom @font-face rule missing font-display descriptor):
```astro
<style>
  @font-face {
    font-family: 'GeistSans';
    src: url('/fonts/geist.woff2') format('woff2');
  }
</style>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Custom @font-face declaring font-display: swap):
```astro
<style>
  @font-face {
    font-family: 'GeistSans';
    src: url('/fonts/geist.woff2') format('woff2');
    font-display: swap;
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: cls.font-display-missing"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore cls.font-display-missing` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/cls.font-display-missing/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for cls.font-display-missing"]
        subgraph P ["Positive Corpus (tests/correctness/cls.font-display-missing/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/cls.font-display-missing/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/cls.font-display-missing/adversarial/)"]
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
<!-- charites:ignore cls.font-display-missing intentional exception -->
```

```tsx
// charites:ignore cls.font-display-missing intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.font-display-missing:
    severity: error # error | warn | info | off
```

