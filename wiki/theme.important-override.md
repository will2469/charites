# theme.important-override

> **Rule ID:** `theme.important-override`
> **Severity:** `ERROR`
> **Category:** `theme`
> **Target Standards:** W3C Cascading Style Sheets (CSS) Specificity Level 4, Tailwind CSS Design Token Architecture

---

## 1. Overview & Core Invariant

Detects !important modifiers on color utility classes that break theme cascade and specificity hierarchy

### Core Invariant:
> **"Color utility classes must never use the !important modifier (!bg-*, !text-*); specificity must be managed via CSS Cascade Layers."**

---
## 2. Technical Grounding & Engine Realities

Using the ! modifier (e.g. !bg-red-500 or !text-white) forcefully escalates CSS declaration priority above normal cascade layers.

1. Destroys Theme Inversion: Dark mode variants (.dark bg-card) cannot override !bg-white without also adding !dark:bg-card, sparking an !important arms race.
2. Compromises Component Reusability: Reusable components with !important color classes cannot be customized or themed by parent containers.
3. Unpredictable State Styling: Hover, focus, and disabled state colors fail to trigger reliably when base colors are marked !important.

Charites enforces relying on Cascade Layers (@layer components, utilities) and semantic token definitions.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cascade Arms Race** | HIGH | Forces downstream theme overrides to duplicate !important, breaking modular CSS encapsulation. |
| **Dark Mode Override Failure** | HIGH | Dark mode variants fail to override base !important styles, resulting in inverted visual glitches. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (!important modifier on background and text color):
```astro
<button class="!bg-red-500 !text-white">Delete</button>
```
### TSX (!important on semantic and hover colors in JSX):
```tsx
export function Action() {
  return <div className="hover:!bg-primary !border-border">Action</div>;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Proper layer-based specificity without !important):
```astro
<button class="bg-destructive text-destructive-foreground">Delete</button>
```
### TSX (Clean semantic classes with natural CSS cascade):
```tsx
export function Action() {
  return <div className="hover:bg-primary border-border">Action</div>;
}
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
    IgnoreCheck -- "Not Ignored" --> Diag["7. Emit Diagnostic: theme.important-override with Replacement Suggestion"]
```

### Step-by-Step Evaluation:
1. **AST Node Traversal:** `internal/analyzer` streams JSX/Astro AST elements to the rule's `Evaluate` visitor.
2. **Variant Normalization:** Strips responsive (`sm:`, `md:`), interaction state (`hover:`, `focus:`), and theme (`dark:`) prefixes to isolate the core utility class.
3. **Modifier Extraction:** Parses utility segments and extracts slash opacity modifiers.
4. **Token Convention Resolution:** Consults the `TokenConvention` adapter to determine the official semantic design token replacement candidate.
5. **Token Graph Verification (Banana Test):** Queries `token.Context` to verify that the candidate token is declared in `global.css` or `tokens.json` within the element's scope. If not declared, the custom value is permitted without a false-positive diagnostic.
6. **Directive Suppression Check:** Inspects preceding AST comments for `charites:ignore theme.important-override`.
7. **Diagnostic Emission:** Produces a structured diagnostic with line number, column span, and actionable replacement suggestion.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/theme.important-override/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for theme.important-override"]
        subgraph P ["Positive Corpus (tests/correctness/theme.important-override/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/theme.important-override/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/theme.important-override/adversarial/)"]
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
<!-- charites:ignore theme.important-override intentional exception -->
```

```tsx
// charites:ignore theme.important-override intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.important-override:
    severity: error # error | warn | info | off
```

