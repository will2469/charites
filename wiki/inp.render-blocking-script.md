# inp.render-blocking-script

> **Rule ID:** `inp.render-blocking-script`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** HTML Living Standard (The script element & execution pipeline), W3C Web Performance & Navigation Timing Specification, Google Chrome Core Web Vitals (Eliminating Render-Blocking Resources)

---

## 1. Overview & Core Invariant

External script element without defer, async, or type="module" synchronously blocks rendering and input responsiveness

### Core Invariant:
> **"External script elements must declare 'defer', 'async', or 'type="module"' to avoid synchronously blocking HTML parsing and main-thread readiness."**

---
## 2. Technical Grounding & Engine Realities

When the browser encounters a synchronous `<script src="...">` tag, it must pause HTML parsing, establish a network connection, download the script, and execute it before resuming document rendering.

In Astro, standard `<script>` tags are automatically bundled into deferred ES modules. However, scripts marked with `is:inline` or raw external scripts in HTML document heads bypass bundling and execute synchronously.

Adding 'defer' or 'type="module"' ensures the script is downloaded in the background and executed without halting the parser, keeping the browser immediately receptive to early user taps and clicks.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Synchronous Parser Halting** | HIGH | HTML parsing and initial rendering are paused until external scripts download and execute. |
| **Delayed Main-Thread Input Availability** | MEDIUM | The browser input event loop is delayed, resulting in unacknowledged early user taps. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Synchronous external inline script blocking HTML parser):
```astro
<script is:inline src="https://analytics.example.com/heavy-bundle.js"></script>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (External inline script deferred to prevent parser blocking):
```astro
<script is:inline src="https://analytics.example.com/heavy-bundle.js" defer></script>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.render-blocking-script"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.render-blocking-script` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.render-blocking-script/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.render-blocking-script"]
        subgraph P ["Positive Corpus (tests/correctness/inp.render-blocking-script/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.render-blocking-script/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.render-blocking-script/adversarial/)"]
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
<!-- charites:ignore inp.render-blocking-script intentional exception -->
```

```tsx
// charites:ignore inp.render-blocking-script intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.render-blocking-script:
    severity: warn # error | warn | info | off
```

