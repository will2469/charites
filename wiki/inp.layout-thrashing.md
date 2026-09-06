# inp.layout-thrashing

> **Rule ID:** `inp.layout-thrashing`
> **Severity:** `ERROR`
> **Category:** `inp`
> **Target Standards:** W3C Web Performance Working Group (Interaction to Next Paint - INP), Google Chrome Rendering Engine Pipeline (Forced Synchronous Layout), Browser Main-Thread Event Loop Scheduling

---

## 1. Overview & Core Invariant

Sequential DOM style mutation followed by layout geometry reading triggers forced synchronous reflow

### Core Invariant:
> **"Imperative JavaScript execution must separate layout queries from style mutations, avoiding forced synchronous reflow passes (read-then-write batching)."**

---
## 2. Technical Grounding & Engine Realities

When JavaScript mutates DOM styles or class names (e.g. 'el.style.width = ...') and subsequently reads a layout geometry property (e.g. 'el.offsetHeight' or 'getBoundingClientRect()') within the same synchronous execution block, the browser is forced to flush pending style changes and perform an immediate, blocking layout recalculation.

This phenomenon, known as 'Layout Thrashing' or 'Forced Synchronous Reflow', locks the browser main thread, preventing user interaction processing and drastically inflating Interaction to Next Paint (INP) latency.

Batching all layout reads before performing style writes, or deferring updates via 'requestAnimationFrame', prevents synchronous recalculations and keeps the main thread responsive.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Forced Synchronous Reflow Stalling** | HIGH | Synchronous layout computation blocks the main thread during interaction handling, causing dropped frames and severe INP degradation. |
| **Cascading Rendering Bottleneck** | HIGH | Interleaved write-read loops exponentially degrade interaction responsiveness on complex DOM trees. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Style mutation immediately followed by geometry reading (forced reflow)):
```tsx
function resizeBox(el: HTMLElement) {
  el.style.width = '200px';
  const height = el.offsetHeight; // Memaksa kalkulasi layout sinkron!
  el.style.height = (height * 2) + 'px';
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Read-then-write batching to prevent forced layout calculation):
```tsx
function resizeBox(el: HTMLElement) {
  const currentHeight = el.offsetHeight; // Baca di awal
  el.style.width = '200px';              // Tulis serentak
  el.style.height = (currentHeight * 2) + 'px';
}
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.layout-thrashing"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.layout-thrashing` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.layout-thrashing/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.layout-thrashing"]
        subgraph P ["Positive Corpus (tests/correctness/inp.layout-thrashing/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.layout-thrashing/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.layout-thrashing/adversarial/)"]
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
<!-- charites:ignore inp.layout-thrashing intentional exception -->
```

```tsx
// charites:ignore inp.layout-thrashing intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.layout-thrashing:
    severity: error # error | warn | info | off
```

