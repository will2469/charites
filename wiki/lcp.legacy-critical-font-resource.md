# lcp.legacy-critical-font-resource

> **Rule ID:** `lcp.legacy-critical-font-resource`
> **Severity:** `WARN`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Duration), W3C WOFF File Format 2.0 (WOFF2) Recommendation, IETF Brotli Compressed Data Format Specification

---

## 1. Overview & Core Invariant

Custom '@font-face' declaration provides only legacy uncompressed font formats (.ttf, .otf, .eot) or deprioritizes WOFF2 in 'src:', inflating byte transfer payload

### Core Invariant:
> **"Custom @font-face declarations for web fonts must specify the modern WOFF2 format as the first item in the 'src' descriptor to maximize compression efficiency."**

---
## 2. Technical Grounding & Engine Realities

Legacy font formats such as raw TrueType (.ttf), OpenType (.otf), and Embedded OpenType (.eot) lack modern compression algorithms, resulting in file sizes ranging from 200KB to 800KB per font weight.

WOFF2 utilizes the Brotli compression algorithm, reducing font binary size by 50% to 80% compared to TTF/OTF and approximately 30% compared to WOFF 1.0 without loss of font hinting or OpenType layout features.

Browsers evaluate 'src' declarations in sequential order. Declaring WOFF2 first guarantees that modern browsers download the most compressed variant, accelerating the Resource Load Duration of LCP text elements.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Massive Font Transfer Payload** | HIGH | Downloading uncompressed 500KB+ TTF/OTF font files inflates Resource Load Duration on mobile networks. |
| **Bandwidth Competition with Hero Media** | MEDIUM | Bulky font files compete for socket bandwidth against hero images and critical CSS stylesheets during early page loading. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### CSS (Font declaration only provides raw uncompressed TTF format):
```css
@font-face {
  font-family: 'HeadingDisplay';
  src: url('/fonts/heading.ttf') format('truetype');
  font-display: swap;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### CSS (WOFF2 declared as primary format with progressive TTF fallback):
```css
@font-face {
  font-family: 'HeadingDisplay';
  src: url('/fonts/heading.woff2') format('woff2'),
       url('/fonts/heading.ttf') format('truetype');
  font-display: swap;
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: lcp.legacy-critical-font-resource"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore lcp.legacy-critical-font-resource` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/lcp.legacy-critical-font-resource/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for lcp.legacy-critical-font-resource"]
        subgraph P ["Positive Corpus (tests/correctness/lcp.legacy-critical-font-resource/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/lcp.legacy-critical-font-resource/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/lcp.legacy-critical-font-resource/adversarial/)"]
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
<!-- charites:ignore lcp.legacy-critical-font-resource intentional exception -->
```

```tsx
// charites:ignore lcp.legacy-critical-font-resource intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.legacy-critical-font-resource:
    severity: warn # error | warn | info | off
```

