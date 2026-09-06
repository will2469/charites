# inp.unbounded-effect-deps

> **Rule ID:** `inp.unbounded-effect-deps`
> **Severity:** `ERROR`
> **Category:** `inp`
> **Target Standards:** React Hooks Specification & Dependency Determinism, W3C Cooperative Scheduling & Frame Budget Invariants, Google Chrome Core Web Vitals (Input Presentation Delay)

---

## 1. Overview & Core Invariant

Lifecycle hook useEffect/useLayoutEffect is missing a dependency array, triggering unbounded re-executions on every render

### Core Invariant:
> **"React lifecycle hooks (useEffect, useLayoutEffect) must explicitly declare a dependency array as their second argument to prevent uncontrolled execution on every render cycle."**

---
## 2. Technical Grounding & Engine Realities

When 'useEffect' or 'useLayoutEffect' is invoked with only a callback and no second argument, React executes the effect after *every single render*.

Any state update, parent re-render, or user keystroke causes the entire effect callback to run again. If the effect queries DOM elements, reads layout properties, or synchronizes subscriptions, the main thread is constantly saturated by unnecessary computations.

Providing an explicit dependency array ('[]' for mount-only, or '[deps...]') restricts execution strictly to when dependencies change, protecting interaction frame rate.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Every-Render Effect Re-execution** | CRITICAL | Effects fire repeatedly on every keystroke or state change, causing severe CPU spikes and input lag. |
| **Infinite Render Loops** | HIGH | If an unbounded effect updates state, it causes an immediate infinite re-render loop that locks the browser tab. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (useEffect without a dependency array executes on every render):
```tsx
useEffect(() => {
  recomputeHeavyLayout();
});
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Explicit empty dependency array ensures execution only on mount):
```tsx
useEffect(() => {
  recomputeHeavyLayout();
}, []);
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.unbounded-effect-deps"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.unbounded-effect-deps` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.unbounded-effect-deps/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.unbounded-effect-deps"]
        subgraph P ["Positive Corpus (tests/correctness/inp.unbounded-effect-deps/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.unbounded-effect-deps/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.unbounded-effect-deps/adversarial/)"]
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
<!-- charites:ignore inp.unbounded-effect-deps intentional exception -->
```

```tsx
// charites:ignore inp.unbounded-effect-deps intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.unbounded-effect-deps:
    severity: error # error | warn | info | off
```

