# theme.backdrop-blur-hardcode

> **Rule ID:** `theme.backdrop-blur-hardcode`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C Filter & Backdrop Filter Specification, Design System Glassmorphism Standards, Hardware-Accelerated Compositing Guidelines

---

## 1. Overview & Core Invariant

Detects hardcoded arbitrary blur and backdrop-blur scalars in Tailwind utility classes

### Core Invariant:
> **"Glassmorphism and surface blur effects must adhere to standardized blur tokens, never arbitrary scalar lengths."**

---
## 2. Technical Grounding & Engine Realities

Using arbitrary blur values (e.g. backdrop-blur-[5px] or blur-[12px]) produces inconsistent glassmorphism and performance bottlenecks:

1. GPU Overdraw Fragility: Arbitrary blur radii bypass optimized compositor layer pooling, causing unnecessary GPU rasterization penalties on mobile devices.
2. Glassmorphism Fragmentation: Slightly differing blur radii (e.g. 5px vs 8px vs 10px) across headers, dialogs, and drawer sheets ruin visual polish.
3. Inflexible Accessibility Adjustments: Standard tokens allow globally disabling or tuning blurs for users requesting reduced motion or low-power modes.

Charites enforces utilizing standard blur scale tokens (e.g. backdrop-blur-sm, backdrop-blur-md, backdrop-blur-lg) or CSS variables.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Glassmorphism Visual Discordance** | MEDIUM | Irregular blur intensity breaks cohesive layering and depth hierarchy across interface overlays. |
| **Mobile GPU Performance Stutter** | HIGH | Unstandardized backdrop-filter passes induce dropped frames during touch scrolling and bottom sheet gestures. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Arbitrary backdrop-blur on navigation header):
```tsx
<header className="backdrop-blur-[5px] bg-background/80">Sticky Nav</header>
```
### ASTRO (Arbitrary filter blur in Astro component):
```astro
<div class="blur-[12px] [backdrop-filter:blur(7px)]">Frosted Panel</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Using standard backdrop blur token):
```tsx
<header className="backdrop-blur-md bg-background/80">Sticky Nav</header>
```
### ASTRO (Standard filter blur token):
```astro
<div class="blur-md backdrop-blur-sm">Frosted Panel</div>
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
    IgnoreCheck -- "Not Ignored" --> Diag["7. Emit Diagnostic: theme.backdrop-blur-hardcode with Replacement Suggestion"]
```

### Step-by-Step Evaluation:
1. **AST Node Traversal:** `internal/analyzer` streams JSX/Astro AST elements to the rule's `Evaluate` visitor.
2. **Variant Normalization:** Strips responsive (`sm:`, `md:`), interaction state (`hover:`, `focus:`), and theme (`dark:`) prefixes to isolate the core utility class.
3. **Modifier Extraction:** Parses utility segments and extracts slash opacity modifiers.
4. **Token Convention Resolution:** Consults the `TokenConvention` adapter to determine the official semantic design token replacement candidate.
5. **Token Graph Verification (Banana Test):** Queries `token.Context` to verify that the candidate token is declared in `global.css` or `tokens.json` within the element's scope. If not declared, the custom value is permitted without a false-positive diagnostic.
6. **Directive Suppression Check:** Inspects preceding AST comments for `charites:ignore theme.backdrop-blur-hardcode`.
7. **Diagnostic Emission:** Produces a structured diagnostic with line number, column span, and actionable replacement suggestion.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/theme.backdrop-blur-hardcode/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for theme.backdrop-blur-hardcode"]
        subgraph P ["Positive Corpus (tests/correctness/theme.backdrop-blur-hardcode/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/theme.backdrop-blur-hardcode/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/theme.backdrop-blur-hardcode/adversarial/)"]
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
<!-- charites:ignore theme.backdrop-blur-hardcode intentional exception -->
```

```tsx
// charites:ignore theme.backdrop-blur-hardcode intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.backdrop-blur-hardcode:
    severity: warn # error | warn | info | off
```

