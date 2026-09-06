# performance.astro-unnecessary-client-directive

> **Rule ID:** `performance.astro-unnecessary-client-directive`
> **Severity:** `ERROR`
> **Category:** `performance`
> **Target Standards:** Astro Islands Architecture Specification (Zero-JS Baseline Principle), W3C Web Performance Client-Side Script Minimization Invariants, Astro Official Documentation ('Template Directives: client:*')

---

## 1. Overview & Core Invariant

Menegakkan prinsip Zero-JS Astro dengan melarang penambahan direktif hidrasi (client:*) pada komponen antarmuka yang murni statis.

### Core Invariant:
> **"Static UI components must not include 'client:*' hydration directives; adding hydration directives to non-interactive components forces the framework runtime and component bundle to be downloaded, violating Astro's Zero-JS guarantee."**

---
## 2. Technical Grounding & Engine Realities

Astro by default renders all components to pure, static HTML at build time with zero client-side JavaScript overhead.

When a developer unnecessarily adds a `client:*` directive (`client:load`, `client:idle`, `client:visible`) to a purely presentational component, Astro treats it as an interactive island.

This forces the bundler to extract the component into a separate client bundle and ship the framework runtime (such as React or Vue, weighing 30-50KB+) to the browser, needlessly delaying page interactivity and squandering network bandwidth.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Zero-JS Guarantee Violation** | HIGH | Transmits unnecessary framework runtimes and component code to the client, increasing page weight and parse time. |
| **Main Thread Hydration Lag** | MEDIUM | Wastes browser CPU cycles hydrating static DOM trees that have no event listeners or interactive state. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Header statis dipaksa terhidrasi ke peramban):
```astro
---
import HeaderStatic from '../components/HeaderStatic.tsx';
---
<HeaderStatic client:load title="Selamat Datang" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Dirender sebagai pure static HTML tanpa JavaScript):
```astro
---
import HeaderStatic from '../components/HeaderStatic.tsx';
---
<HeaderStatic title="Selamat Datang" />
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: performance.astro-unnecessary-client-directive"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore performance.astro-unnecessary-client-directive` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/performance.astro-unnecessary-client-directive/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for performance.astro-unnecessary-client-directive"]
        subgraph P ["Positive Corpus (tests/correctness/performance.astro-unnecessary-client-directive/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/performance.astro-unnecessary-client-directive/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/performance.astro-unnecessary-client-directive/adversarial/)"]
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
<!-- charites:ignore performance.astro-unnecessary-client-directive intentional exception -->
```

```tsx
// charites:ignore performance.astro-unnecessary-client-directive intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.astro-unnecessary-client-directive:
    severity: error # error | warn | info | off
```

