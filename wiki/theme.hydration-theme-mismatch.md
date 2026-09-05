# theme.hydration-theme-mismatch

> **Rule ID:** `theme.hydration-theme-mismatch`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C Web Performance Working Group (Core Web Vitals FOUC Prevention), React 18/19 Hydration Boundary Specification, Astro SSR Zero-JS Script Tag Standards

---

## 1. Overview & Core Invariant

Detects SSR root layouts lacking blocking inline script for theme initialization

### Core Invariant:
> **"Root SSR document layouts (<head>) must include a render-blocking inline theme script to resolve theme state before first paint and prevent theme FOUC."**

---
## 2. Technical Grounding & Engine Realities

In Server-Side Rendered (SSR) architectures (such as Astro, Next.js, or Remix):

1. Flash of Unstyled Theme (FOUC): If theme detection runs only after deferred client hydration (e.g. inside useEffect), the browser paints a blinding white default page before snapping jarringly to dark mode.
2. React Hydration Mismatch: Inconsistent theme attributes between server-rendered HTML and client hydration trigger React warning cascades and forced DOM re-mounts.
3. Cumulative Layout Shift (CLS): Font, border, or icon shifts caused by late theme flipping harm Core Web Vitals.

Charites enforces placing an inline render-blocking theme initialization script directly in the SSR root <head>.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Theme FOUC Glare** | HIGH | Users in dark environments experience a painful full-screen white flash on every page navigation. |
| **Hydration Error Cascade** | MEDIUM | React discards server-rendered DOM nodes upon encountering mismatched class attributes, increasing TTI. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (SSR root head in Astro without inline theme script):
```astro
<html>
  <head>
    <meta charset="utf-8" />
    <title>Application</title>
  </head>
  <body>
    <slot />
  </body>
</html>
```
### TSX (Root head in TSX missing blocking theme initializer):
```tsx
<head>
  <meta charSet="utf-8" />
  <title>Dashboard</title>
</head>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Blocking inline theme script in Astro head):
```astro
<head>
  <meta charset="utf-8" />
  <script is:inline>
    const theme = localStorage.getItem('theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    document.documentElement.classList.toggle('dark', theme === 'dark');
  </script>
</head>
```
### TSX (Blocking dangerouslySetInnerHTML theme script in TSX head):
```tsx
<head>
  <script
    dangerouslySetInnerHTML={{
      __html: "document.documentElement.classList.add(localStorage.getItem('theme') || 'light');",
    }}
  />
</head>
```

---

## 6. Detection & Verification Pipeline (How The Rule Evaluates Code)
This `theme` rule evaluates source templates against the project's design token graph:

```mermaid
flowchart TD
    Node["AST Node (Astro / TSX element)"] --> Extract["1. Extract Class Names (e.g. 'hover:bg-primary/10')"]
    Extract --> Strip["2. Strip Variants (hover:, dark:, sm:) -> 'bg-primary/10'"]
    Strip --> Split["3. Split Opacity Modifier -> Utility: 'bg-primary', Opacity: '/10'"]
    Split --> Convention["4. Query TokenConvention (Candidate: '--color-primary-light')"]
    Convention --> GraphQuery{"5. Check Token Graph (Does token exist in active scope?)"}
    GraphQuery -- "No (Banana Test)" --> Safe["Pass (Valid Custom / Untokenized Color)"]
    GraphQuery -- "Yes (Official Token Exists)" --> IgnoreCheck{"6. Check charites:ignore directive"}
    IgnoreCheck -- "Ignored" --> Safe
    IgnoreCheck -- "Not Ignored" --> Diag["7. Emit Diagnostic: theme.hydration-theme-mismatch with Replacement Suggestion"]
```

### Step-by-Step Evaluation:
1. **AST Node Traversal:** `internal/analyzer` streams JSX/Astro AST elements to the rule's `Evaluate` visitor.
2. **Variant Normalization:** Strips responsive (`sm:`, `md:`), interaction state (`hover:`, `focus:`), and theme (`dark:`) prefixes to isolate the core utility class.
3. **Modifier Extraction:** Parses utility segments and extracts slash opacity modifiers.
4. **Token Convention Resolution:** Consults the `TokenConvention` adapter to determine the official semantic design token replacement candidate.
5. **Token Graph Verification (Banana Test):** Queries `token.Context` to verify that the candidate token is declared in `global.css` or `tokens.json` within the element's scope. If not declared, the custom value is permitted without a false-positive diagnostic.
6. **Directive Suppression Check:** Inspects preceding AST comments for `charites:ignore theme.hydration-theme-mismatch`.
7. **Diagnostic Emission:** Produces a structured diagnostic with line number, column span, and actionable replacement suggestion.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/theme.hydration-theme-mismatch/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for theme.hydration-theme-mismatch"]
        subgraph P ["Positive Corpus (tests/correctness/theme.hydration-theme-mismatch/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/theme.hydration-theme-mismatch/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/theme.hydration-theme-mismatch/adversarial/)"]
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
<!-- charites:ignore theme.hydration-theme-mismatch intentional exception -->
```

```tsx
// charites:ignore theme.hydration-theme-mismatch intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.hydration-theme-mismatch:
    severity: warn # error | warn | info | off
```

