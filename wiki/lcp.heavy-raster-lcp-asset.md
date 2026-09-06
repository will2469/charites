# lcp.heavy-raster-lcp-asset

> **Rule ID:** `lcp.heavy-raster-lcp-asset`
> **Severity:** `WARN`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Duration), W3C Web Performance Working Group Media Optimization Guidelines, IETF AVIF / WebP Image Compression Standards

---

## 1. Overview & Core Invariant

LCP candidate image uses legacy uncompressed raster format (.png, .bmp, .tiff, .gif); modern formats like WebP or AVIF should be served to reduce transfer size

### Core Invariant:
> **"Above-the-fold LCP candidate images must utilize next-generation compressed formats (WebP, AVIF) rather than legacy uncompressed raster formats (.png, .bmp, .tiff, .gif) to minimize byte transfer payload."**

---
## 2. Technical Grounding & Engine Realities

Serving high-resolution photographs or hero imagery in legacy raster formats such as PNG or uncompressed BMP results in massive byte payloads (often 2MB-5MB per image).

Next-generation formats such as WebP and AVIF provide superior lossy and lossless compression algorithms, reducing image file sizes by 30% to 70% compared to PNG and JPEG without perceptual visual degradation.

For above-the-fold hero images that dictate the LCP metric, reducing file transfer size directly accelerates the Resource Load Duration phase over bandwidth-constrained networks.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Resource Load Duration Delay** | HIGH | Downloading uncompressed 2MB-5MB PNG/BMP images over 4G/cellular connections introduces 800ms-3000ms delay to LCP. |
| **Memory Footprint & GPU Texture Pressure** | MEDIUM | Large uncompressed raster graphics consume excessive client RAM and GPU texture memory during decode and compositing. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Critical hero image served as an uncompressed 3MB PNG file):
```tsx
<section className="hero-section" data-perf-role="hero">
  <img
    src="/assets/hero-banner.png"
    alt="Hero Showcase"
    fetchpriority="high"
    className="w-full h-auto"
  />
</section>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Hero image converted to compressed modern WebP format):
```tsx
<section className="hero-section" data-perf-role="hero">
  <img
    src="/assets/hero-banner.webp"
    alt="Hero Showcase"
    fetchpriority="high"
    className="w-full h-auto"
  />
</section>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: lcp.heavy-raster-lcp-asset"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore lcp.heavy-raster-lcp-asset` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/lcp.heavy-raster-lcp-asset/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for lcp.heavy-raster-lcp-asset"]
        subgraph P ["Positive Corpus (tests/correctness/lcp.heavy-raster-lcp-asset/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/lcp.heavy-raster-lcp-asset/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/lcp.heavy-raster-lcp-asset/adversarial/)"]
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
<!-- charites:ignore lcp.heavy-raster-lcp-asset intentional exception -->
```

```tsx
// charites:ignore lcp.heavy-raster-lcp-asset intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.heavy-raster-lcp-asset:
    severity: warn # error | warn | info | off
```

