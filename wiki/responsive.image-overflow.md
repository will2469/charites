# responsive.image-overflow

> **Rule ID:** `responsive.image-overflow`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** HTML Living Standard (Embedded Content: img, video, svg), Web.dev Responsive Media & Core Web Vitals (CLS Prevention)

---

## 1. Overview & Core Invariant

Warns when media elements with large fixed dimensions lack responsive max-w-full scaling

### Core Invariant:
> **"Media elements with explicit width dimensions exceeding 320px must declare 'max-w-full' or 'w-full' to prevent horizontal viewport tearing on mobile screens."**

---
## 2. Technical Grounding & Engine Realities

Specifying explicit 'width' and 'height' attributes on media elements is recommended for Core Web Vitals to reserve aspect ratio boxes and prevent Cumulative Layout Shift (CLS).

However, when large static dimensions (e.g. width={1200}) lack responsive CSS scaling ('max-w-full h-auto'), mobile browsers render the media at full absolute physical pixel width, breaking outside narrow 360px viewports and causing severe horizontal scrolling.

Applying 'max-w-full h-auto' preserves CLS aspect ratio benefits while ensuring the media smoothly downsizes to fit compact screens.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Viewport Tearing** | MEDIUM | Large unconstrained images expand outside mobile viewport boundaries, forcing horizontal scrollbars. |
| **Aspect Ratio Distortion** | LOW | Images constrained by height but not width stretch disproportionately on narrow viewports. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Media element with large width attribute lacking max-w-full):
```tsx
<img src="/hero-desa.jpg" width={1200} height={800} alt="Pemandangan Desa" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Responsive media element with max-w-full and h-auto):
```tsx
<img className="max-w-full h-auto" src="/hero-desa.jpg" width={1200} height={800} alt="Pemandangan Desa" />
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: responsive.image-overflow"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore responsive.image-overflow` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/responsive.image-overflow/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for responsive.image-overflow"]
        subgraph P ["Positive Corpus (tests/correctness/responsive.image-overflow/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/responsive.image-overflow/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/responsive.image-overflow/adversarial/)"]
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
<!-- charites:ignore responsive.image-overflow intentional exception -->
```

```tsx
// charites:ignore responsive.image-overflow intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.image-overflow:
    severity: warn # error | warn | info | off
```

