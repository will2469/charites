# performance.react-unstable-hook-reference

> **Rule ID:** `performance.react-unstable-hook-reference`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** React Official Documentation (Building Your Own Custom Hooks), React Hooks Referential Integrity & Stable Function Contracts, React Hooks Exhaustive Dependencies Safety Guidelines

---

## 1. Overview & Core Invariant

Mengaudit custom hook yang mengembalikan referensi fungsi tidak stabil tanpa dibungkus useCallback, yang memicu re-render loop pada komponen konsumen.

### Core Invariant:
> **"Custom React hooks exposing helper functions must stabilize them with 'useCallback'; returning fresh function instances causes downstream consumers using them in effect dependencies to trigger infinite render loops."**

---
## 2. Technical Grounding & Engine Realities

Custom hooks frequently return an object containing state and mutation functions (e.g. `{ data, refetch, reset }`).

If these functions are defined as regular arrow functions without `useCallback`, a brand-new function reference is created in memory on every render pass of the consuming component.

When the consuming component passes this function into the dependency array of `useEffect` or `useMemo`, or passes it down to a memoized child component, the newly allocated reference violates referential equality, defeating memoization and in many cases causing uncontrollable infinite re-render loops.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Downstream Infinite Render Loops** | HIGH | Triggers continuous re-execution of downstream useEffect hooks that list the unmemoized helper function in their dependency arrays. |
| **Bypassed Child Memoization** | MEDIUM | Breaks shallow prop comparison (React.memo) across all child components consuming functions returned from the custom hook. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (refetch dialokasikan sebagai fungsi baru di setiap pemanggilan hook):
```tsx
export function useProfile(userId: string) {
  const [data, setData] = useState(null);
  const refetch = () => { fetchProfile(userId).then(setData); };
  return { data, refetch };
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Menstabilkan referensi fungsi dengan useCallback):
```tsx
export function useProfile(userId: string) {
  const [data, setData] = useState(null);
  const refetch = useCallback(() => {
    fetchProfile(userId).then(setData);
  }, [userId]);
  return { data, refetch };
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: performance.react-unstable-hook-reference"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore performance.react-unstable-hook-reference` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/performance.react-unstable-hook-reference/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for performance.react-unstable-hook-reference"]
        subgraph P ["Positive Corpus (tests/correctness/performance.react-unstable-hook-reference/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/performance.react-unstable-hook-reference/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/performance.react-unstable-hook-reference/adversarial/)"]
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
<!-- charites:ignore performance.react-unstable-hook-reference intentional exception -->
```

```tsx
// charites:ignore performance.react-unstable-hook-reference intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.react-unstable-hook-reference:
    severity: warn # error | warn | info | off
```

