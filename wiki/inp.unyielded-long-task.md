# inp.unyielded-long-task

> **Rule ID:** `inp.unyielded-long-task`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** W3C Cooperative Scheduling Controller (scheduler.yield), Google Chrome Core Web Vitals (Long Tasks & Input Responsiveness), Main-Thread Cooperative Concurrency Invariants

---

## 1. Overview & Core Invariant

Long task processing large arrays without cooperative scheduling yields stalls main-thread responsiveness

### Core Invariant:
> **"Long execution tasks triggered by or affecting user interactions must periodically yield control to the browser event loop via cooperative scheduling boundaries."**

---
## 2. Technical Grounding & Engine Realities

Long tasks running uninterrupted on the main thread (> 50ms) prevent the browser from acknowledging new user inputs (clicks, keypresses, taps) or rendering visual updates.

When user actions initiate extensive batch computations, running the entire process synchronously locks the page until completion, producing high Interaction to Next Paint (INP) latency.

By periodically pausing execution using modern cooperative scheduling: 'await (window.scheduler?.yield?.() ?? new Promise(r => setTimeout(r, 0)))', the browser is given immediate opportunities to handle pending user inputs and paint frames before continuing task work.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Main-Thread Input Starvation** | HIGH | Long uninterrupted execution loops starve the browser input queue, leaving pages unresponsive during batch processing. |
| **High INP & Dropped Frames** | MEDIUM | Presentation of user feedback is blocked for hundreds of milliseconds, breaching the 200ms INP threshold. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TS (Long calculation loop over large dataset without cooperative yield):
```ts
function processLargeArray(items: string[]) {
  for (let i = 0; i < items.length; i++) {
    heavyCalculation(items[i]);
  }
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TS (Periodic cooperative yielding to maintain input responsiveness):
```ts
async function processLargeArray(items: string[]) {
  for (let i = 0; i < items.length; i++) {
    heavyCalculation(items[i]);
    if (i % 50 === 0) {
      await (window.scheduler?.yield?.() ?? new Promise(r => setTimeout(r, 0)));
    }
  }
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.unyielded-long-task"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.unyielded-long-task` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.unyielded-long-task/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.unyielded-long-task"]
        subgraph P ["Positive Corpus (tests/correctness/inp.unyielded-long-task/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.unyielded-long-task/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.unyielded-long-task/adversarial/)"]
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
<!-- charites:ignore inp.unyielded-long-task intentional exception -->
```

```tsx
// charites:ignore inp.unyielded-long-task intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.unyielded-long-task:
    severity: warn # error | warn | info | off
```

