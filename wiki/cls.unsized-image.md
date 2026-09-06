# cls.unsized-image

> **Rule ID:** `cls.unsized-image`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** W3C Cumulative Layout Shift (CLS) Metric Specification, Google Core Web Vitals Guidelines (Target CLS < 0.1), W3C CSS Box Sizing Module Level 4 (aspect-ratio), Astro Docs: Image Optimization (astro:assets)

---

## 1. Overview & Core Invariant

Warns when image elements lack explicit dimensions, aspect-ratio, or Tailwind box sizing

### Core Invariant:
> **"Image elements must establish a statically inferable reserved rendering box via explicit width/height attributes, CSS aspect-ratio, or Tailwind sizing utilities before the binary asset is downloaded."**

---
## 2. Technical Grounding & Engine Realities

When browsers parse an <img> tag without explicit dimensions or an aspect-ratio reservation, the layout engine initially allocates a 0x0 pixel box.

Once the remote image file is fetched and decoded, the browser performs a sudden reflow to accommodate the intrinsic image geometry, pushing surrounding content downward. This layout instability directly penalizes Cumulative Layout Shift (CLS) scores.

Specifying width and height attributes or utilizing Tailwind 'aspect-video' / 'aspect-square' allows modern browsers to compute the aspect ratio before network I/O completes, eliminating visual jank.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cumulative Layout Shift (CLS)** | HIGH | Unsized images push subsequent content down upon load, degrading Core Web Vitals and SEO rankings. |
| **Accidental User Mis-clicks** | MEDIUM | Users attempting to tap links or buttons near loading images accidentally trigger shifted elements. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Image with fluid width but missing height or aspect-ratio reservation):
```tsx
<img src={heroUrl} alt="Hero Banner" className="w-full h-auto" />
```
### ASTRO (Standard img tag lacking width and height attributes):
```astro
<img src="/pemandangan-desa.jpg" alt="Pemandangan Desa" class="rounded-lg" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Image with explicit numeric width and height attributes):
```tsx
<img src={heroUrl} alt="Hero Banner" width={1200} height={600} className="w-full h-auto" />
```
### TSX (Image with Tailwind v4 aspect-ratio utility):
```tsx
<img src={heroUrl} alt="Hero Banner" className="w-full aspect-video object-cover" />
```
### TSX (Avatar image with explicit width and height sizing utilities):
```tsx
<img src={avatarUrl} alt="Avatar" className="w-10 h-10 rounded-full" />
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: cls.unsized-image"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore cls.unsized-image` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/cls.unsized-image/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for cls.unsized-image"]
        subgraph P ["Positive Corpus (tests/correctness/cls.unsized-image/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/cls.unsized-image/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/cls.unsized-image/adversarial/)"]
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
<!-- charites:ignore cls.unsized-image intentional exception -->
```

```tsx
// charites:ignore cls.unsized-image intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.unsized-image:
    severity: warn # error | warn | info | off
```

