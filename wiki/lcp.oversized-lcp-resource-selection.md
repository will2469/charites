# lcp.oversized-lcp-resource-selection

> **Rule ID:** `lcp.oversized-lcp-resource-selection`
> **Severity:** `WARN`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Duration), Responsive Images Community Group (RICG) Responsive Images Specification, HTML Living Standard srcset and sizes Attributes Specification

---

## 1. Overview & Core Invariant

Fluid responsive LCP candidate image lacks responsive 'srcset' and 'sizes' attributes, forcing mobile viewports to download oversized desktop assets

### Core Invariant:
> **"Fluid responsive LCP candidate images must provide responsive 'srcset' with width descriptors and a 'sizes' attribute to prevent mobile devices from downloading oversized desktop assets."**

---
## 2. Technical Grounding & Engine Realities

When a fluid image (such as a full-width hero banner) only specifies a single large 'src' attribute, mobile devices with small viewports are forced to download the same high-resolution asset designed for 4K desktop screens.

This unnecessary byte payload directly prolongs the Resource Load Duration component of LCP over cellular networks.

By providing a 'srcset' attribute with width descriptors (e.g. '400w, 800w, 1200w') alongside a matching 'sizes' attribute (or using the '<Image />' component from 'astro:assets'), the browser can accurately select the optimal image variant for the user's viewport and device pixel ratio.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Resource Load Duration Bloat** | HIGH | Mobile devices download 2MB-5MB desktop-resolution assets over cellular connections, adding 500ms-2500ms to LCP. |
| **Excess Mobile Data Consumption** | MEDIUM | Users on metered cellular data plans consume excessive bandwidth downloading unneeded pixels. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Fluid hero image only provides a single massive desktop asset without srcset and sizes):
```tsx
<section className="hero-section" data-perf-role="hero">
  <img
    src="/images/hero-3840x2160.jpg"
    alt="Hero Banner"
    className="w-full h-auto"
    fetchpriority="high"
  />
</section>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Fluid hero image configured with responsive srcset width descriptors and sizes attribute):
```tsx
<section className="hero-section" data-perf-role="hero">
  <img
    src="/images/hero-1200.webp"
    srcset="/images/hero-400.webp 400w, /images/hero-800.webp 800w, /images/hero-1200.webp 1200w"
    sizes="(max-width: 768px) 100vw, 1200px"
    alt="Hero Banner"
    className="w-full h-auto"
    fetchpriority="high"
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: lcp.oversized-lcp-resource-selection"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore lcp.oversized-lcp-resource-selection` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/lcp.oversized-lcp-resource-selection/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for lcp.oversized-lcp-resource-selection"]
        subgraph P ["Positive Corpus (tests/correctness/lcp.oversized-lcp-resource-selection/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/lcp.oversized-lcp-resource-selection/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/lcp.oversized-lcp-resource-selection/adversarial/)"]
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
<!-- charites:ignore lcp.oversized-lcp-resource-selection intentional exception -->
```

```tsx
// charites:ignore lcp.oversized-lcp-resource-selection intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.oversized-lcp-resource-selection:
    severity: warn # error | warn | info | off
```

