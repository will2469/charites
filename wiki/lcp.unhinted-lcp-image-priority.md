# lcp.unhinted-lcp-image-priority

> **Rule ID:** `lcp.unhinted-lcp-image-priority`
> **Severity:** `WARN`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay), W3C Priority Hints Specification (fetchpriority attribute), Chrome Preload Scanner Network Bandwidth Scheduling

---

## 1. Overview & Core Invariant

Above-the-fold LCP candidate image lacks fetchpriority="high", delaying bandwidth allocation in early network stream

### Core Invariant:
> **"Above-the-fold LCP candidate images must declare 'fetchpriority="high"' to prioritize early network bandwidth ahead of non-critical stylesheets and scripts."**

---
## 2. Technical Grounding & Engine Realities

By default, browsers assign an initial fetch priority of 'Low' to image resources discovered in the HTML stream.

For the primary hero image (the LCP element), this default low priority forces the image download to compete with or yield to lower-priority scripts, stylesheets, and fonts.

Declaring 'fetchpriority="high"' instructs the speculative preload scanner to immediately elevate the resource to the highest network tier, initiating the TCP/TLS transfer with maximum allocated bandwidth.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Bandwidth Starvation by Non-Critical Assets** | HIGH | Hero image bytes are delayed behind non-critical deferred scripts and fonts, inflating LCP by 150ms-400ms. |
| **Sub-optimal Browser Network Scheduling** | MEDIUM | Browsers under HTTP/2 or HTTP/3 multiplexing prioritize other resources unless explicitly hinted. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Above-the-fold hero banner image lacking priority hint):
```tsx
<header className="hero-banner" data-perf-role="hero">
  <img src="/hero.webp" alt="Primary Banner" className="w-full aspect-video" />
</header>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Hero image explicitly prioritized with fetchpriority='high'):
```tsx
<header className="hero-banner" data-perf-role="hero">
  <img src="/hero.webp" alt="Primary Banner" fetchpriority="high" className="w-full aspect-video" />
</header>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.unhinted-lcp-image-priority intentional exception -->
```

```tsx
// charites:ignore lcp.unhinted-lcp-image-priority intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.unhinted-lcp-image-priority:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


