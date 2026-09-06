# responsive.grid-min-column

> **Rule ID:** `responsive.grid-min-column`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** W3C CSS Grid Layout Module Level 1 (The minmax() Function), WCAG 2.2 SC 1.4.10 (Reflow - Level AA), Mobile Web Best Practices: Preventing Horizontal Viewport Blowout

---

## 1. Overview & Core Invariant

Warns against CSS grid minmax column definitions with rigid minimum sizes (> 320px) that cause horizontal overflow on mobile viewports

### Core Invariant:
> **"CSS grid column minmax tracks on mobile baseline must not enforce rigid minimum widths greater than 320px without dynamic clamping ('min(100%, <size>)') or desktop breakpoint scoping ('md:grid-cols-...')."**

---
## 2. Technical Grounding & Engine Realities

A common CSS grid pattern for auto-fit cards is 'repeat(auto-fit, minmax(350px, 1fr))' or 'repeat(auto-fill, minmax(400px, 1fr))'. While this looks great on desktop and tablet monitors, the minimum column track width of 350px or 400px exceeds the 360px physical width of most smartphones and the 320px minimum WCAG reflow baseline.

Because CSS grid does not shrink tracks below their minimum minmax threshold, the grid blows out horizontally, introducing an unintended horizontal scrollbar across the entire mobile page.

Charites detects rigid minmax tracks and suggests using the standard modern CSS clamp idiom: 'minmax(min(100%, 20rem), 1fr)' or scoping multi-column grids behind 'md:' breakpoints.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Viewport Horizontal Scrollbar Blowout** | HIGH | The entire website scrolls sideways on mobile phones because a single card grid enforces 350px+ minimum track size. |
| **Broken Touch Gestures & Visual Glitches** | MEDIUM | Accidental horizontal swiping triggers page drift instead of vertical scroll. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Grid specifying 400px minimum column width on mobile baseline):
```tsx
<div className="grid grid-cols-[repeat(auto-fit,minmax(400px,1fr))] gap-4">
  <div className="card">Kartu Layanan 1</div>
  <div className="card">Kartu Layanan 2</div>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Clamped minmax ensuring column never exceeds 100% on narrow screens):
```tsx
<div className="grid grid-cols-[repeat(auto-fit,minmax(min(100%,20rem),1fr))] gap-4">
  <div className="card">Kartu Layanan 1</div>
  <div className="card">Kartu Layanan 2</div>
</div>
```
### TSX (Mobile single-column with desktop-scoped multi-column minmax):
```tsx
<div className="grid grid-cols-1 md:grid-cols-[repeat(auto-fit,minmax(350px,1fr))] gap-4">
  <div className="card">Kartu Layanan 1</div>
  <div className="card">Kartu Layanan 2</div>
</div>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: responsive.grid-min-column"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore responsive.grid-min-column` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/responsive.grid-min-column/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for responsive.grid-min-column"]
        subgraph P ["Positive Corpus (tests/correctness/responsive.grid-min-column/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/responsive.grid-min-column/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/responsive.grid-min-column/adversarial/)"]
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
<!-- charites:ignore responsive.grid-min-column intentional exception -->
```

```tsx
// charites:ignore responsive.grid-min-column intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.grid-min-column:
    severity: warn # error | warn | info | off
```

