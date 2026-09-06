# lcp.preload-font-cors-mismatch

> **Rule ID:** `lcp.preload-font-cors-mismatch`
> **Severity:** `ERROR`
> **Category:** `lcp`
> **Target Standards:** W3C Preload Specification (Font Preload CORS Requirements), HTML Living Standard Crossorigin Attribute Specification, Google Chrome Core Web Vitals (Largest Contentful Paint Resource Optimization)

---

## 1. Overview & Core Invariant

Font preload '<link rel="preload" as="font">' lacks 'crossorigin' attribute, triggering browser cache discard and double network downloads

### Core Invariant:
> **"All '<link rel="preload" as="font">' tags must specify the 'crossorigin' attribute to ensure the preloaded font binary is accepted by the browser font cache."**

---
## 2. Technical Grounding & Engine Realities

The W3C CSS Fonts specification mandates that web fonts must be fetched using anonymous Cross-Origin Resource Sharing (CORS) mode, even when the font is hosted on the same origin as the page.

When a '<link rel="preload" as="font">' tag omits the 'crossorigin' attribute, the preload scanner fetches the font using standard non-CORS mode.

When the CSS parser subsequently requests the font in anonymous CORS mode, the browser detects a CORS mode mismatch, discards the preloaded resource from cache, and executes a second, redundant network download.

Adding 'crossorigin' (or 'crossorigin="anonymous"') ensures the preloaded font matches the CSS font engine request key, avoiding double downloads.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Redundant Network Double Download** | CRITICAL | The font asset is downloaded twice over the network, completely defeating the purpose of preloading and inflating bandwidth usage. |
| **LCP Text Rendering Delay** | HIGH | The second fetch is queued after CSS parsing, adding hundreds of milliseconds to text block paint times. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### HTML (Font preloaded without crossorigin attribute triggering cache discard):
```html
<head>
  <link rel="preload" href="/fonts/inter.woff2" as="font" type="font/woff2" />
</head>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### HTML (Font preload declared with crossorigin attribute matching W3C anonymous CORS requirement):
```html
<head>
  <link rel="preload" href="/fonts/inter.woff2" as="font" type="font/woff2" crossorigin />
</head>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.preload-font-cors-mismatch intentional exception -->
```

```tsx
// charites:ignore lcp.preload-font-cors-mismatch intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.preload-font-cors-mismatch:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


