# inp.repeated-state-update

> **Rule ID:** `inp.repeated-state-update`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** React 18+ Automatic Batching Specification, W3C Web Performance Working Group (Interaction to Next Paint - INP), Concurrent React Scheduling & Reconciliation Cost

---

## 1. Overview & Core Invariant

Repeated state updater calls inside loops breaking automatic batching trigger cascading re-renders

### Core Invariant:
> **"React state setters must not be repeatedly invoked within loop iterations that break automatic batching (such as asynchronous loops containing 'await' or 'flushSync')."**

---
## 2. Technical Grounding & Engine Realities

While React 18 automatically batches multiple state updates within standard synchronous handlers, asynchronous loops (e.g. 'for ... of' with 'await' inside) or explicit 'flushSync' blocks break automatic batching.

Calling a state updater on every iteration of an asynchronous loop causes React to trigger a full re-render, VDOM diffing, and reconciliation cycle on every microtask tick.

This creates an enormous render queue backlog on the main thread, stalling user interactions and causing high Interaction to Next Paint (INP) latency.

Accumulating results locally into an array and issuing a single state update after the loop completes ensures a single, batched render pass.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Per-Iteration Re-render Cascades** | HIGH | Each iteration of an async loop schedules a separate render pass, saturating the React scheduler and freezing UI input. |
| **Presentation Delay Ballooning** | MEDIUM | Successive re-renders continuously postpone the browser paint phase, severely degrading INP. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (State updater called on each iteration of an async loop):
```tsx
for (const item of items) {
  const detail = await fetchDetail(item.id);
  setItems(prev => [...prev, detail]); // Memicu re-render pada setiap iterasi!
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Accumulating all results and updating state once after loop completion):
```tsx
const results = [];
for (const item of items) {
  results.push(await fetchDetail(item.id));
}
setItems(prev => [...prev, ...results]); // Hanya satu siklus render
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.repeated-state-update"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.repeated-state-update` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.repeated-state-update/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.repeated-state-update"]
        subgraph P ["Positive Corpus (tests/correctness/inp.repeated-state-update/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.repeated-state-update/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.repeated-state-update/adversarial/)"]
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
<!-- charites:ignore inp.repeated-state-update intentional exception -->
```

```tsx
// charites:ignore inp.repeated-state-update intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.repeated-state-update:
    severity: warn # error | warn | info | off
```

