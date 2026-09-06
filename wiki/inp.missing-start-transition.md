# inp.missing-start-transition

> **Rule ID:** `inp.missing-start-transition`
> **Severity:** `INFO`
> **Category:** `inp`
> **Target Standards:** React 18/19 Concurrent Mode Architecture (startTransition & Transitions API), W3C User Timing & Cooperative Scheduling Invariants, Google Chrome Core Web Vitals (Input to Paint Responsiveness)

---

## 1. Overview & Core Invariant

Secondary non-urgent state update inside interactive handler should be wrapped in startTransition to prevent input lag

### Core Invariant:
> **"Secondary non-urgent state updates triggered alongside urgent user input must be scheduled as transitions via 'React.startTransition' to preserve typing responsiveness."**

---
## 2. Technical Grounding & Engine Realities

In modern user interfaces, an interactive event (such as typing in a search bar or clicking a filter tab) often triggers two types of updates: an urgent update (updating the input text cursor) and a non-urgent secondary update (filtering a large list or fetching preview cards).

When both updates are processed synchronously without transitions, React treats the expensive secondary re-render with the same high priority as the keystroke, blocking the main thread and causing noticeable keystroke stutter.

Wrapping secondary updates in 'React.startTransition' informs the scheduler that the secondary render is interruptible. React will immediately paint the user's keystroke, keeping INP low while deferring list rendering.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Keystroke Input Lag** | MEDIUM | Synchronous secondary re-renders block subsequent keystroke frames, creating sluggish typing feedback. |
| **Main Thread Presentation Delays** | LOW | The browser cannot acknowledge user interactions within the 200ms INP threshold. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Synchronously combining urgent input setter with heavy list filtering):
```tsx
function handleSearch(e: React.ChangeEvent<HTMLInputElement>) {
  setSearchQuery(e.target.value);
  setFilteredLargeList(expensiveFilter(e.target.value));
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Urgent input text is updated immediately; secondary list is wrapped in startTransition):
```tsx
function handleSearch(e: React.ChangeEvent<HTMLInputElement>) {
  setSearchQuery(e.target.value);
  React.startTransition(() => {
    setFilteredLargeList(expensiveFilter(e.target.value));
  });
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.missing-start-transition"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.missing-start-transition` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.missing-start-transition/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.missing-start-transition"]
        subgraph P ["Positive Corpus (tests/correctness/inp.missing-start-transition/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.missing-start-transition/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.missing-start-transition/adversarial/)"]
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
<!-- charites:ignore inp.missing-start-transition intentional exception -->
```

```tsx
// charites:ignore inp.missing-start-transition intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.missing-start-transition:
    severity: info # error | warn | info | off
```

