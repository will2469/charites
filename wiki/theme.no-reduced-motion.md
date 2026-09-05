# theme.no-reduced-motion

> **Rule ID:** `theme.no-reduced-motion`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** WCAG 2.2 Success Criterion 2.3.3 (Animation from Interactions), W3C Media Queries Level 5 (prefers-reduced-motion), Accessible Web Animation & Vestibular Safety Guidelines

---

## 1. Overview & Core Invariant

Detects global theme transitions without prefers-reduced-motion media query wrapping

### Core Invariant:
> **"Global theme and color transitions must be scoped within prefers-reduced-motion: no-preference or mitigated with reduced-motion overrides."**

---
## 2. Technical Grounding & Engine Realities

Smooth CSS transitions applied to root or theme switching (such as * { transition: background-color 0.3s, color 0.3s; } or transition: all 0.2s) can cause dizziness, headaches, and nausea for users with vestibular disorders.

WCAG 2.2 Success Criterion 2.3.3 requires that non-essential animations triggered by user interaction can be turned off or respect system accessibility preferences.

Charites enforces wrapping theme transitions in @media (prefers-reduced-motion: no-preference) or providing an explicit @media (prefers-reduced-motion: reduce) { transition: none; } fallback.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Vestibular Distress** | MEDIUM | Rapid or uncontrolled surface transitions induce disorientation or motion sickness for sensitive users. |
| **WCAG 2.2 SC 2.3.3 Non-Compliance** | MEDIUM | Failure to honor OS-level accessibility preferences prevents compliance with regulatory accessibility standards. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Unmitigated global theme transition in Astro style):
```astro
<style>
  * {
    transition: background-color 0.3s ease, color 0.3s ease;
  }
</style>
```
### TSX (Broad transition all without motion preference in TSX):
```tsx
<style>{`
  body {
    transition: all 0.25s ease-in-out;
  }
`}</style>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Theme transition scoped to no-preference media query):
```astro
<style>
  @media (prefers-reduced-motion: no-preference) {
    * {
      transition: background-color 0.3s ease, color 0.3s ease;
    }
  }
</style>
```
### ASTRO (Explicit reduced-motion override):
```astro
<style>
  body {
    transition: background-color 0.3s ease;
  }
  @media (prefers-reduced-motion: reduce) {
    body {
      transition: none;
    }
  }
</style>
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
    IgnoreCheck -- "Not Ignored" --> Diag["7. Emit Diagnostic: theme.no-reduced-motion with Replacement Suggestion"]
```

### Step-by-Step Evaluation:
1. **AST Node Traversal:** `internal/analyzer` streams JSX/Astro AST elements to the rule's `Evaluate` visitor.
2. **Variant Normalization:** Strips responsive (`sm:`, `md:`), interaction state (`hover:`, `focus:`), and theme (`dark:`) prefixes to isolate the core utility class.
3. **Modifier Extraction:** Parses utility segments and extracts slash opacity modifiers.
4. **Token Convention Resolution:** Consults the `TokenConvention` adapter to determine the official semantic design token replacement candidate.
5. **Token Graph Verification (Banana Test):** Queries `token.Context` to verify that the candidate token is declared in `global.css` or `tokens.json` within the element's scope. If not declared, the custom value is permitted without a false-positive diagnostic.
6. **Directive Suppression Check:** Inspects preceding AST comments for `charites:ignore theme.no-reduced-motion`.
7. **Diagnostic Emission:** Produces a structured diagnostic with line number, column span, and actionable replacement suggestion.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/theme.no-reduced-motion/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for theme.no-reduced-motion"]
        subgraph P ["Positive Corpus (tests/correctness/theme.no-reduced-motion/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/theme.no-reduced-motion/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/theme.no-reduced-motion/adversarial/)"]
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
<!-- charites:ignore theme.no-reduced-motion intentional exception -->
```

```tsx
// charites:ignore theme.no-reduced-motion intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.no-reduced-motion:
    severity: warn # error | warn | info | off
```

