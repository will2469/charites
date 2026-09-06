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

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.oversized-lcp-resource-selection intentional exception -->
```

```tsx
// charites:ignore lcp.oversized-lcp-resource-selection intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.oversized-lcp-resource-selection:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


