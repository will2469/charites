# lcp.missing-lcp-image-preload

> **Rule ID:** `lcp.missing-lcp-image-preload`
> **Severity:** `INFO`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay), W3C Preload Specification (<link rel="preload" as="image">), Document Layout Graph & Early Resource Discovery Optimization

---

## 1. Overview & Core Invariant

Delayed-discovery LCP image lacks <link rel="preload" as="image"> in document head to initiate early asset transfer

### Core Invariant:
> **"LCP candidate media elements that suffer delayed discovery (dynamic data attributes, client script hydration, or CSS backgrounds) should be preloaded in '<head>' with 'fetchpriority="high"'."**

---
## 2. Technical Grounding & Engine Realities

When an LCP image cannot be immediately parsed from a direct '<img>' element in the server-rendered HTML stream (for example, when its URL is stored in a dynamic data attribute 'data-bg-src', rendered by a client island, or defined via CSS background), its network fetch is delayed.

Injecting '<link rel="preload" as="image" href="..." fetchpriority="high">' inside '<head>' compensates for this discovery delay by instructing the browser lookahead scanner to initiate connection and download immediately during initial document streaming.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Delayed Resource Discovery** | MEDIUM | Image download is postponed until JavaScript hydration or style resolution finishes, inflating LCP. |
| **Initial Viewport Flash** | LOW | Late arrival of visual hero media causes prolonged empty or placeholder hero appearance. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Hero container with dynamic data-bg-src without preload in document head):
```astro
<head>
  <title>Product Gallery</title>
</head>
<body>
  <div id="hero-root" data-perf-role="hero" data-bg-src="https://cdn.example.com/promo.webp"></div>
</body>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Document head preloading the hero image asset with high fetch priority):
```astro
<head>
  <title>Product Gallery</title>
  <link rel="preload" as="image" href="https://cdn.example.com/promo.webp" fetchpriority="high" />
</head>
<body>
  <div id="hero-root" data-perf-role="hero" data-bg-src="https://cdn.example.com/promo.webp"></div>
</body>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.missing-lcp-image-preload intentional exception -->
```

```tsx
// charites:ignore lcp.missing-lcp-image-preload intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.missing-lcp-image-preload:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


