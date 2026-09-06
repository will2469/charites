# performance.react-derived-state-in-effect

> **Rule ID:** `performance.react-derived-state-in-effect`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** React Official Documentation ('You Might Not Need an Effect'), React Reconciliation Lifecycle (Avoiding Cascading Secondary Renders), React Best Practices on Pure Render-Phase Computations

---

## 1. Overview & Core Invariant

Mencegah sinkronisasi derived state dari props atau state yang sudah ada melalui useEffect, yang memicu siklus perenderan sekunder ganda.

### Core Invariant:
> **"Derived values computed synchronously from props or existing state must be calculated directly in the component body during render; updating state inside 'useEffect' triggers redundant secondary render passes."**

---
## 2. Technical Grounding & Engine Realities

Updating state within a `useEffect` callback causes React to first render the component with the stale value, commit it to the DOM, and immediately schedule a second render pass to apply the updated state.

When derived state is calculated synchronously (e.g. concatenating names, filtering a list, or calculating a total), calculating it in an effect needlessly burns main thread time on layout calculations and DOM diffing twice per interaction.

Computing the value directly during the render pass completely eliminates the secondary render pass and keeps state management minimal.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cascading Secondary Render Cycles** | HIGH | Forces React to run duplicate render and diff cycles on every prop change, directly degrading interaction responsiveness (INP). |
| **Visual Stutter / Layout Shift** | MEDIUM | May momentarily display stale computed values before the effect updates, causing brief content flicker or layout shifts. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Sinkronisasi derived state via useEffect memicu render ganda):
```tsx
const [fullName, setFullName] = useState('');
useEffect(() => {
  setFullName(firstName + ' ' + lastName);
}, [firstName, lastName]);
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Dihitung secara sinkron dalam satu kali fase render):
```tsx
const fullName = firstName + ' ' + lastName;
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: performance.react-derived-state-in-effect"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore performance.react-derived-state-in-effect` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/performance.react-derived-state-in-effect/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for performance.react-derived-state-in-effect"]
        subgraph P ["Positive Corpus (tests/correctness/performance.react-derived-state-in-effect/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/performance.react-derived-state-in-effect/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/performance.react-derived-state-in-effect/adversarial/)"]
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
<!-- charites:ignore performance.react-derived-state-in-effect intentional exception -->
```

```tsx
// charites:ignore performance.react-derived-state-in-effect intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.react-derived-state-in-effect:
    severity: warn # error | warn | info | off
```

