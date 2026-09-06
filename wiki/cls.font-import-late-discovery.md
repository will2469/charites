# cls.font-import-late-discovery

> **Rule ID:** `cls.font-import-late-discovery`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** W3C Cumulative Layout Shift (CLS) Metric Specification, Google Core Web Vitals Guidelines (Render-Blocking Resources & Critical Path), Tailwind CSS v4 Import Specifications

---

## 1. Overview & Core Invariant

Warns when CSS @import is used for external font loading, delaying discovery and risking layout shift

### Core Invariant:
> **"External web fonts must be discovered and preconnected in HTML/Astro <head> rather than imported via cascading CSS @import directives, while whitelisting Tailwind CSS and local stylesheets."**

---
## 2. Technical Grounding & Engine Realities

Placing @import rules referencing external fonts (such as Google Fonts or Typekit) inside CSS creates a cascading waterfall of render-blocking requests.

The browser must download the HTML, parse the stylesheet, discover the nested @import, download the font CSS, and only then start downloading the binary font file. This profound delay forces long periods of fallback rendering and dramatic late layout shifts.

Declaring '<link rel="preconnect">' alongside '<link rel="stylesheet">' in the HTML layout starts DNS preconnection and font loading at the earliest possible phase of page load.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cascading Network Waterfall** | HIGH | Nested font CSS imports delay font delivery by several hundred milliseconds on mobile networks. |
| **Severe Late Layout Shift** | MEDIUM | Delayed font swapping abruptly reorganizes the text geometry long after initial paint. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Late discovery font import inside style block):
```astro
<style>
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap');
</style>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Fonts loaded via preconnect and stylesheet link in head):
```astro
<head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap" />
</head>
```
### ASTRO (Whitelisted Tailwind CSS and local file imports):
```astro
<style>
  @import "tailwindcss";
  @import "./local-tokens.css";
</style>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore cls.font-import-late-discovery intentional exception -->
```

```tsx
// charites:ignore cls.font-import-late-discovery intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.font-import-late-discovery:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [cls Category Guide](cls).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


