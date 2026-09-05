# theme.chart-color-hardcode

> **Rule ID:** `theme.chart-color-hardcode`
> **Severity:** `ERROR`
> **Category:** `theme`
> **Target Standards:** WCAG 2.2 Success Criterion 1.4.3 (Contrast Minimum), WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast), Accessible Data Visualization Design Tokens

---

## 1. Overview & Core Invariant

Detects hardcoded color values on chart visualization components

### Core Invariant:
> **"Chart components must reference semantic theme tokens (e.g. var(--chart-1)) rather than hardcoded hex or color literals."**

---
## 2. Technical Grounding & Engine Realities

Data visualization libraries (such as Recharts, Chart.js, or Nivo) rely on SVG fill and stroke attributes to render bars, lines, and areas.

When developers hardcode hex colors onto chart elements (e.g. <Bar dataKey="sales" fill="#3b82f6" />):
1. Dark Mode Contrast Inversion: The hardcoded colors clash with dark card backgrounds, failing accessibility contrast minimums.
2. Theme Blindness: Visualizations fail to adapt when switching between light, dark, or high-contrast themes.
3. Fragmented Visual Identity: Brand colors drift between charts and surrounding interface tokens.

Charites enforces using CSS custom properties (fill="var(--chart-1)") or dynamic theme mappings.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Chart Contrast Invalidation** | HIGH | Chart bars and lines become illegible against inverted dark backgrounds, obscuring critical analytics. |
| **Theme Desynchronization** | MEDIUM | Data visualizations remain locked to legacy colors while the rest of the application adapts dynamically. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Hardcoded hex fill and stroke on Recharts Bar and Line):
```tsx
<>
  <Bar dataKey="revenue" fill="#3b82f6" />
  <Line dataKey="profit" stroke="#10b981" />
</>
```
### ASTRO (Hardcoded color on Area and Cell components):
```astro
<Area dataKey="uv" fill="#8884d8" stroke="#82ca9d" />
<Cell fill="#f43f5e" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Semantic chart tokens from design system):
```tsx
<>
  <Bar dataKey="revenue" fill="var(--chart-1)" />
  <Line dataKey="profit" stroke="var(--chart-2)" />
</>
```
### ASTRO (CSS variable references on Area and Cell):
```astro
<Area dataKey="uv" fill="var(--chart-1)" stroke="var(--chart-2)" />
<Cell fill="var(--chart-destructive)" />
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
    IgnoreCheck -- "Not Ignored" --> Diag["7. Emit Diagnostic: theme.chart-color-hardcode with Replacement Suggestion"]
```

### Step-by-Step Evaluation:
1. **AST Node Traversal:** `internal/analyzer` streams JSX/Astro AST elements to the rule's `Evaluate` visitor.
2. **Variant Normalization:** Strips responsive (`sm:`, `md:`), interaction state (`hover:`, `focus:`), and theme (`dark:`) prefixes to isolate the core utility class.
3. **Modifier Extraction:** Parses utility segments and extracts slash opacity modifiers.
4. **Token Convention Resolution:** Consults the `TokenConvention` adapter to determine the official semantic design token replacement candidate.
5. **Token Graph Verification (Banana Test):** Queries `token.Context` to verify that the candidate token is declared in `global.css` or `tokens.json` within the element's scope. If not declared, the custom value is permitted without a false-positive diagnostic.
6. **Directive Suppression Check:** Inspects preceding AST comments for `charites:ignore theme.chart-color-hardcode`.
7. **Diagnostic Emission:** Produces a structured diagnostic with line number, column span, and actionable replacement suggestion.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/theme.chart-color-hardcode/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for theme.chart-color-hardcode"]
        subgraph P ["Positive Corpus (tests/correctness/theme.chart-color-hardcode/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/theme.chart-color-hardcode/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/theme.chart-color-hardcode/adversarial/)"]
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
<!-- charites:ignore theme.chart-color-hardcode intentional exception -->
```

```tsx
// charites:ignore theme.chart-color-hardcode intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.chart-color-hardcode:
    severity: error # error | warn | info | off
```

