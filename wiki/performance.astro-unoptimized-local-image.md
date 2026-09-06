# performance.astro-unoptimized-local-image

> **Rule ID:** `performance.astro-unoptimized-local-image`
> **Severity:** `INFO`
> **Category:** `performance`
> **Target Standards:** Astro Asset Pipeline Best Practices ('astro:assets' Image & Picture), Core Web Vitals Largest Contentful Paint (LCP) Image Payload Optimization, W3C Next-Gen Responsive Image Delivery Standards (WebP/AVIF)

---

## 1. Overview & Core Invariant

Menganjurkan pemakaian komponen <Image /> dari astro:assets pada gambar lokal guna mengaktifkan konversi format modern dan kompresi build.

### Core Invariant:
> **"Local raster image assets in Astro templates should be rendered via '<Image />' from 'astro:assets' rather than raw '<img>' tags to leverage automated build-time format conversion and dimension inference."**

---
## 2. Technical Grounding & Engine Realities

Astro provides a native image optimization pipeline through the `astro:assets` module.

Using a raw HTML `<img>` tag pointing to a local file path (`src="../assets/banner.png"`) completely bypasses this pipeline, serving uncompressed, legacy formats (PNG/JPEG) with no automatic width/height dimension injection.

Migrating to `<Image />` allows Astro to automatically convert images to AVIF/WebP, generate responsive srcset attributes, and prevent Cumulative Layout Shift (CLS) by inferring exact dimensions at build time.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Inflated Asset Payload** | LOW | Serves unoptimized PNG/JPEG images that are 40-70% larger than modern WebP/AVIF equivalents. |
| **Missing Intrinsic Aspect Ratio** | LOW | Raw img tags without width and height attributes cause layout shifts during image load. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Tag img mentah melewatkan kompresi build-time Astro):
```astro
<!-- Advisory: Tag img mentah pada path lokal -->
<img src="../assets/product-hero.png" alt="Produk Baru" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Memanfaatkan komponen Image bawaan astro:assets):
```astro
---
import { Image } from 'astro:assets';
import productImg from '../assets/product-hero.png';
---
<Image src={productImg} alt="Produk Baru" />
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: performance.astro-unoptimized-local-image"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore performance.astro-unoptimized-local-image` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/performance.astro-unoptimized-local-image/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for performance.astro-unoptimized-local-image"]
        subgraph P ["Positive Corpus (tests/correctness/performance.astro-unoptimized-local-image/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/performance.astro-unoptimized-local-image/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/performance.astro-unoptimized-local-image/adversarial/)"]
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
<!-- charites:ignore performance.astro-unoptimized-local-image intentional exception -->
```

```tsx
// charites:ignore performance.astro-unoptimized-local-image intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.astro-unoptimized-local-image:
    severity: info # error | warn | info | off
```

