# lcp.lazy-loaded-lcp-image

> **Rule ID:** `lcp.lazy-loaded-lcp-image`
> **Severity:** `ERROR`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay), HTML Living Standard Lazy Loading Specification, W3C Web Performance Working Group Invariants

---

## 1. Overview & Core Invariant

Critical above-the-fold LCP candidate image has loading="lazy", delaying resource discovery and load initiation

### Core Invariant:
> **"Above-the-fold LCP candidate images must not be configured with loading='lazy'; lazy loading defers image download until layout completion, directly adding hundreds of milliseconds to LCP."**

---
## 2. Technical Grounding & Engine Realities

When a browser encounters an '<img>' with 'loading="lazy"', it deliberately pauses fetching the image resource until the page layout is calculated and the element is verified to be within or near the viewport.

For hero images and above-the-fold content that constitute the Largest Contentful Paint (LCP), this artificial pause wastes the initial network idle period. The browser speculative preload scanner is effectively blocked from fetching the hero asset early.

Removing 'loading="lazy"' or declaring 'loading="eager"' combined with 'fetchpriority="high"' allows the browser to initiate the network download immediately upon parsing the HTML token.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Resource Load Delay Inflation** | CRITICAL | Hero image download is postponed until stylesheet download, CSS parsing, and layout pass complete, adding 200ms-800ms to LCP. |
| **Speculative Preload Scanner Suppression** | HIGH | The browser's high-speed HTML lookahead parser skips downloading the hero asset during early stream processing. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Above-the-fold hero banner image configured with loading='lazy'):
```tsx
<section className="hero-section" data-perf-role="hero">
  <h1>Welcome to Our Platform</h1>
  <img src="/assets/hero.webp" alt="Hero Banner" loading="lazy" className="w-full h-auto" />
</section>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Hero image configured with loading='eager' and high fetch priority):
```tsx
<section className="hero-section" data-perf-role="hero">
  <h1>Welcome to Our Platform</h1>
  <img src="/assets/hero.webp" alt="Hero Banner" loading="eager" fetchpriority="high" className="w-full h-auto" />
</section>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.lazy-loaded-lcp-image intentional exception -->
```

```tsx
// charites:ignore lcp.lazy-loaded-lcp-image intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.lazy-loaded-lcp-image:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


