# inp.heavy-event-handler

> **Rule ID:** `inp.heavy-event-handler`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** W3C Web Performance Working Group (Interaction to Next Paint - INP), Google Chrome Core Web Vitals (Input Delay & Processing Duration), Browser Cooperative Scheduling Guidelines (scheduler.yield)

---

## 1. Overview & Core Invariant

Interactive event handler executes heavy synchronous operations (JSON.parse, Array.sort) without cooperative yields

### Core Invariant:
> **"Interactive event handlers must avoid heavy synchronous computations on the main thread, adopting cooperative task yielding or Web Worker offloading."**

---
## 2. Technical Grounding & Engine Realities

When users tap, click, or type, the browser expects the main thread to quickly acknowledge the interaction and schedule the next paint frame (ideally within 50ms, with INP target <= 200ms).

Executing heavy synchronous operations (such as large JSON parsing, array sorting, or complex synchronous data manipulation) directly inside event handler callbacks blocks the main thread during the crucial input processing phase.

This delays the presentation of visual feedback (e.g. active button states, loading spinners) and directly inflates the INP metric.

Breaking long tasks with 'await scheduler.yield?.()' or offloading computation to a dedicated Web Worker allows the browser to present visual feedback immediately before executing intensive processing.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Processing Phase Thread Saturation** | HIGH | Synchronous algorithms in click/key handlers block the main thread, exceeding the 200ms INP threshold. |
| **Frozen Visual Feedback** | MEDIUM | Buttons appear unresponsive or stuck because UI rendering is starved by long synchronous handler execution. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Synchronous heavy data parsing and sorting directly inside onClick handler):
```tsx
<button onClick={() => {
  const data = JSON.parse(hugePayload);
  const sorted = data.sort((a, b) => b.score - a.score);
  setResults(sorted);
}}>
  Urutkan Data
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Cooperative yielding to acknowledge user input before intensive processing):
```tsx
<button onClick={async () => {
  setLoading(true);
  await (window.scheduler?.yield?.() ?? new Promise(r => setTimeout(r, 0)));
  const data = JSON.parse(hugePayload);
  setResults(data.sort((a, b) => b.score - a.score));
  setLoading(false);
}}>
  Urutkan Data
</button>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.heavy-event-handler"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.heavy-event-handler` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.heavy-event-handler/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.heavy-event-handler"]
        subgraph P ["Positive Corpus (tests/correctness/inp.heavy-event-handler/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.heavy-event-handler/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.heavy-event-handler/adversarial/)"]
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
<!-- charites:ignore inp.heavy-event-handler intentional exception -->
```

```tsx
// charites:ignore inp.heavy-event-handler intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.heavy-event-handler:
    severity: warn # error | warn | info | off
```

