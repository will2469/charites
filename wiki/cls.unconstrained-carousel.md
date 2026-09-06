# cls.unconstrained-carousel

> **Rule ID:** `cls.unconstrained-carousel`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** W3C Cumulative Layout Shift (CLS) Metric Specification, W3C CSS Scroll Snap Module Level 1, W3C CSS Box Sizing Module Level 4 (aspect-ratio)

---

## 1. Overview & Core Invariant

Warns when carousel or slider containers lack bounded height or slide aspect-ratio constraints

### Core Invariant:
> **"Carousel and slider viewport tracks must constrain container height or bind slide items to fixed aspect ratios to prevent vertical reflow during slide transitions."**

---
## 2. Technical Grounding & Engine Realities

Horizontal scrolling tracks and carousels render dynamic collections of cards, banners, or images.

When the carousel track lacks an explicit height (e.g. 'h-64' or 'min-h-[300px]') and slides do not have locked aspect ratios, incoming slides with varying image proportions or dynamic text will force the entire container to expand or collapse vertically.

Fixing the container height or assigning 'aspect-video' / 'aspect-square' to slide items ensures layout stability throughout horizontal panning.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Vertical Container Height Jitter** | MEDIUM | Slide transitions with varying content heights push subsequent page content up and down. |
| **Cumulative Layout Shift (CLS)** | HIGH | Carousel height adjustments contribute cumulative shift points during user scrolling. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Horizontal snap container without container height or slide aspect-ratio):
```tsx
<div className="flex overflow-x-auto snap-x">
  {slides.map(s => <img key={s.id} src={s.url} alt={s.title} />)}
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Carousel container with explicit height constraint):
```tsx
<div className="flex overflow-x-auto snap-x h-64 md:h-96 w-full">
  {slides.map(s => (
    <div key={s.id} className="snap-center shrink-0 w-full h-full">
      <img src={s.url} alt={s.title} className="w-full h-full object-cover" />
    </div>
  ))}
</div>
```
### TSX (Carousel slide items locked with aspect-video utility):
```tsx
<div className="flex overflow-x-auto snap-x w-full">
  {slides.map(s => (
    <div key={s.id} className="snap-center shrink-0 w-80 aspect-video">
      <img src={s.url} alt={s.title} className="w-full h-full object-cover" />
    </div>
  ))}
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: cls.unconstrained-carousel"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore cls.unconstrained-carousel` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/cls.unconstrained-carousel/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for cls.unconstrained-carousel"]
        subgraph P ["Positive Corpus (tests/correctness/cls.unconstrained-carousel/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/cls.unconstrained-carousel/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/cls.unconstrained-carousel/adversarial/)"]
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
<!-- charites:ignore cls.unconstrained-carousel intentional exception -->
```

```tsx
// charites:ignore cls.unconstrained-carousel intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.unconstrained-carousel:
    severity: warn # error | warn | info | off
```

