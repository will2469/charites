# cls.layout-trigger-animation

> **Rule ID:** `cls.layout-trigger-animation`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** W3C CSS Animations Level 1 (@keyframes declaration blocks), Google Core Web Vitals (CLS Compositor Thread Guidelines), High Performance Mobile Web (GPU Compositing vs CPU Reflow)

---

## 1. Overview & Core Invariant

CSS @keyframes animation mutates layout-triggering geometry properties instead of GPU-composited transforms

### Core Invariant:
> **"CSS @keyframes animations must mutate GPU-composited layer properties ('transform', 'opacity') rather than layout-triggering geometry properties ('top', 'left', 'width', 'height', 'margin', 'padding')."**

---
## 2. Technical Grounding & Engine Realities

When CSS keyframes animate geometry properties (such as top, left, width, height, margin, or padding), the browser is forced to execute full layout recalculations (reflow) and repaint stages on the main CPU thread for every animation frame (typically 60-120 times per second).

This continuous geometry invalidation directly triggers Cumulative Layout Shift (CLS) for neighboring elements and causes noticeable frame jank (dropped frames) on mobile and low-power hardware.

Modern browser rendering pipelines offload 'transform' and 'opacity' mutations directly to the GPU compositor thread, executing smooth, 60fps animations that never invalidate document geometry or cause layout shifts.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Continuous CPU Layout Reflow** | HIGH | Animating geometry properties causes browser recalculation of surrounding elements on every frame, generating Cumulative Layout Shift. |
| **Rendering Pipeline Jank & Dropped Frames** | MEDIUM | High CPU load from continuous layout reflow stalls main thread execution, resulting in choppy animations and poor touch responsiveness. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### CSS (Keyframe animation mutating positional and margin geometry properties):
```css
@keyframes slideIn {
  from { top: -20px; margin-top: 10px; }
  to { top: 0; margin-top: 0; }
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### CSS (GPU-composited keyframe animation using transform and opacity):
```css
@keyframes slideIn {
  from { transform: translateY(-20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: cls.layout-trigger-animation"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore cls.layout-trigger-animation` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/cls.layout-trigger-animation/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for cls.layout-trigger-animation"]
        subgraph P ["Positive Corpus (tests/correctness/cls.layout-trigger-animation/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/cls.layout-trigger-animation/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/cls.layout-trigger-animation/adversarial/)"]
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
<!-- charites:ignore cls.layout-trigger-animation intentional exception -->
```

```tsx
// charites:ignore cls.layout-trigger-animation intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.layout-trigger-animation:
    severity: warn # error | warn | info | off
```

