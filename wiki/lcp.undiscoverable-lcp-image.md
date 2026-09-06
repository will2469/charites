# lcp.undiscoverable-lcp-image

> **Rule ID:** `lcp.undiscoverable-lcp-image`
> **Severity:** `WARN`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay), Chromium Speculative Preload Scanner Discovery Architecture, W3C Preload Specification (<link rel="preload" as="image">)

---

## 1. Overview & Core Invariant

Above-the-fold hero container loads primary image via CSS background without <link rel="preload"> in document head

### Core Invariant:
> **"Hero visual assets must not be embedded exclusively via CSS background-image without a corresponding '<link rel="preload">' in '<head>'; CSS-based images are undiscoverable by the HTML preload scanner."**

---
## 2. Technical Grounding & Engine Realities

When an image is referenced inside CSS ('background-image: url(...)' or Tailwind 'bg-[url(...)]'), the browser's speculative preload scanner cannot discover the asset URL while streaming the HTML.

The browser must first download all render-blocking CSS, parse the cascade, and run the style computation step before it even learns that the image URL exists. This creates massive Resource Load Delay for LCP.

Migrating the visual background to a native '<img>' element (e.g. with 'absolute inset-0 w-full h-full object-cover -z-10') makes it immediately discoverable in HTML. If CSS background is necessary, injecting '<link rel="preload" as="image" href="..." fetchpriority="high">' into '<head>' bridges the discovery gap.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **CSS Cascade Dependency Block** | HIGH | Hero image discovery is delayed until external CSS stylesheets are downloaded and parsed, adding 300ms-1000ms to LCP. |
| **Speculative Scanner Blindness** | HIGH | Lookahead scanner cannot prefetch the LCP asset during initial document TCP streaming. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Hero container using CSS background-image without head preload):
```tsx
<header className="w-full h-[480px] bg-[url('/hero.webp')] bg-cover" data-perf-role="hero">
  <h1 className="text-white">Galactic Exploration</h1>
</header>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Native <img> with object-cover immediately discoverable by preload scanner):
```tsx
<header className="relative w-full h-[480px] overflow-hidden" data-perf-role="hero">
  <img src="/hero.webp" alt="Hero Background" fetchpriority="high" className="absolute inset-0 w-full h-full object-cover -z-10" />
  <h1 className="relative z-10 text-white p-8">Galactic Exploration</h1>
</header>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.undiscoverable-lcp-image intentional exception -->
```

```tsx
// charites:ignore lcp.undiscoverable-lcp-image intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.undiscoverable-lcp-image:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


