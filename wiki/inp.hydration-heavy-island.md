# inp.hydration-heavy-island

> **Rule ID:** `inp.hydration-heavy-island`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** Astro Zero-JS Server-Side Rendering (SSR) Architecture, React Virtual DOM Hydration Complexity & Reconciliation Budget, W3C Web Performance Working Group (Main Thread Input Delay)

---

## 1. Overview & Core Invariant

Client island wraps excessive static DOM subtree forcing heavy virtual DOM reconciliation on the client

### Core Invariant:
> **"Client-hydrated islands must remain compact and isolate only truly interactive elements; static text, articles, and decorative containers must remain in zero-JS Astro SSR."**

---
## 2. Technical Grounding & Engine Realities

Hydrating a monolithic React island forces the client browser to parse JavaScript, construct virtual DOM representations, and reconcile every DOM node against server HTML-even for completely static elements.

When developers wrap entire articles or document structures inside a single `<ArticleViewer client:load>`, hundreds of static paragraphs and headings are needlessly reconciled, consuming excessive main-thread CPU time.

By decomposing the UI and rendering static content through native Astro components (zero client JS), only individual interactive widgets (such as like buttons or comment inputs) are hydrated, keeping the main thread free for user interaction.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Virtual DOM Reconciliation Bloat** | HIGH | Large static subtrees force long synchronous VDOM tree reconciliation during island hydration. |
| **Excessive Client Bundle & Parse Overhead** | MEDIUM | Shipping static component trees to client bundles increases script evaluation time and input latency. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Static article text wrapped inside a client-hydrated island):
```astro
<ArticleViewerIsland client:load>
  <header><h1>Article Title</h1></header>
  <article>
    <p>Paragraph 1...</p>
    <p>Paragraph 2...</p>
    <p>Paragraph 3...</p>
  </article>
  <CommentButton />
</ArticleViewerIsland>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Static content rendered via zero-JS Astro SSR; only interactive button is an island):
```astro
<header><h1>Article Title</h1></header>
<article>
  <p>Paragraph 1...</p>
  <p>Paragraph 2...</p>
  <p>Paragraph 3...</p>
</article>
<CommentButton client:visible />
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.hydration-heavy-island"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.hydration-heavy-island` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.hydration-heavy-island/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.hydration-heavy-island"]
        subgraph P ["Positive Corpus (tests/correctness/inp.hydration-heavy-island/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.hydration-heavy-island/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.hydration-heavy-island/adversarial/)"]
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
<!-- charites:ignore inp.hydration-heavy-island intentional exception -->
```

```tsx
// charites:ignore inp.hydration-heavy-island intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.hydration-heavy-island:
    severity: warn # error | warn | info | off
```

