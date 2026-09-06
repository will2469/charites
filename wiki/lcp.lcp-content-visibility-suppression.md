# lcp.lcp-content-visibility-suppression

> **Rule ID:** `lcp.lcp-content-visibility-suppression`
> **Severity:** `ERROR`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Element Render Delay), W3C CSS Containment Module Level 2 (content-visibility property specification), Chromium Blink Rendering Pipeline Initial Viewport Invariants

---

## 1. Overview & Core Invariant

Above-the-fold hero container specifies 'content-visibility: auto' or 'content-auto', suppressing initial paint and severely inflating LCP

### Core Invariant:
> **"Initial viewport hero containers are strictly forbidden from specifying 'content-visibility: auto' or 'content-auto' as it instructs the browser engine to suppress early layout and paint passes."**

---
## 2. Technical Grounding & Engine Realities

The CSS property 'content-visibility: auto' is a powerful rendering performance feature for below-the-fold content, allowing browsers to skip layout and painting for elements until they approach the viewport.

However, when applied to above-the-fold hero elements (such as `<header>`, hero sections, or containers with `data-perf-role="hero"`), the browser initially skips rendering the element entirely during the first layout pass.

Only after subsequent scrolling or intersection calculations does the engine lay out and paint the content, resulting in massive Element Render Delay and catastrophically failing LCP benchmarks.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Initial Layout Pass Suppression** | CRITICAL | Blink skips the initial layout and paint of the primary hero container, delaying LCP registration until post-hydration intersection observer checks. |
| **Blank Initial Viewport Flash** | HIGH | Users experience an empty screen or blank space in the initial viewport on fast network connections. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Hero section using content-visibility: auto in the initial viewport):
```tsx
<section className="hero-section content-auto" data-perf-role="hero">
  <h1>Solusi Cloud Enterprise</h1>
  <img src="/hero.webp" fetchpriority="high" />
</section>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Hero section rendered immediately without initial paint suppression):
```tsx
<section className="hero-section" data-perf-role="hero">
  <h1>Solusi Cloud Enterprise</h1>
  <img src="/hero.webp" fetchpriority="high" />
</section>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.lcp-content-visibility-suppression intentional exception -->
```

```tsx
// charites:ignore lcp.lcp-content-visibility-suppression intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.lcp-content-visibility-suppression:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


