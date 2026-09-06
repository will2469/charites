# inp.hydration-contention

> **Rule ID:** `inp.hydration-contention`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** Astro Islands Architecture & Partial Hydration Specification, W3C Cooperative Scheduling & Main-Thread Budget Invariants, Google Core Web Vitals (INP Input Delay & Hydration Contention)

---

## 1. Overview & Core Invariant

Concurrently hydrating multiple Astro client:load islands saturates the main thread and spikes input delay

### Core Invariant:
> **"Astro templates must avoid declaring multiple eager 'client:load' island directives simultaneously; non-critical islands must use deferred hydration directives ('client:idle' or 'client:visible')."**

---
## 2. Technical Grounding & Engine Realities

The 'client:load' directive instructs the browser to immediately fetch and execute island JavaScript upon page load, before user interaction or idle periods.

When multiple islands (3 or more) declare 'client:load' on the same page, their hydration phases execute in parallel or rapid succession on the main thread. This contention monopolizes CPU resources during initial user interactions, generating severe Long Tasks and inflating Input Delay.

By reserving 'client:load' strictly for critical interactive UI (such as primary navigation) and deferring secondary components to 'client:idle' or 'client:visible', the main thread remains responsive to user taps, clicks, and keystrokes.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Initial Hydration CPU Saturation** | HIGH | Multiple islands running concurrent React hydration lock the main thread during the window when users attempt first interaction. |
| **Severe Input Delay Spikes** | MEDIUM | User clicks or keystrokes are queued behind synchronous island hydration tasks, resulting in INP > 200ms. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Multiple non-critical islands concurrently hydrated with client:load):
```astro
---
import HeaderNav from '../components/HeaderNav.tsx';
import SearchBar from '../components/SearchBar.tsx';
import PromoBanner from '../components/PromoBanner.tsx';
---
<HeaderNav client:load />
<SearchBar client:load />
<PromoBanner client:load />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Only critical navigation uses client:load; secondary islands use deferred hydration):
```astro
---
import HeaderNav from '../components/HeaderNav.tsx';
import SearchBar from '../components/SearchBar.tsx';
import PromoBanner from '../components/PromoBanner.tsx';
---
<HeaderNav client:load />
<SearchBar client:idle />
<PromoBanner client:visible />
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.hydration-contention"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.hydration-contention` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.hydration-contention/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.hydration-contention"]
        subgraph P ["Positive Corpus (tests/correctness/inp.hydration-contention/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.hydration-contention/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.hydration-contention/adversarial/)"]
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
<!-- charites:ignore inp.hydration-contention intentional exception -->
```

```tsx
// charites:ignore inp.hydration-contention intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.hydration-contention:
    severity: warn # error | warn | info | off
```

