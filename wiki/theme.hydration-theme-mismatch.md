# theme.hydration-theme-mismatch

> **Rule ID:** `theme.hydration-theme-mismatch`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C Web Performance Working Group (Core Web Vitals FOUC Prevention), React 18/19 Hydration Boundary Specification, Astro SSR Zero-JS Script Tag Standards

---

## 1. Overview & Core Invariant

Detects SSR root layouts lacking blocking inline script for theme initialization

### Core Invariant:
> **"Root SSR document layouts (<head>) must include a render-blocking inline theme script to resolve theme state before first paint and prevent theme FOUC."**

---
## 2. Technical Grounding & Engine Realities

In Server-Side Rendered (SSR) architectures (such as Astro, Next.js, or Remix):

1. Flash of Unstyled Theme (FOUC): If theme detection runs only after deferred client hydration (e.g. inside useEffect), the browser paints a blinding white default page before snapping jarringly to dark mode.
2. React Hydration Mismatch: Inconsistent theme attributes between server-rendered HTML and client hydration trigger React warning cascades and forced DOM re-mounts.
3. Cumulative Layout Shift (CLS): Font, border, or icon shifts caused by late theme flipping harm Core Web Vitals.

Charites enforces placing an inline render-blocking theme initialization script directly in the SSR root <head>.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Theme FOUC Glare** | HIGH | Users in dark environments experience a painful full-screen white flash on every page navigation. |
| **Hydration Error Cascade** | MEDIUM | React discards server-rendered DOM nodes upon encountering mismatched class attributes, increasing TTI. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (SSR root head in Astro without inline theme script):
```astro
<html>
  <head>
    <meta charset="utf-8" />
    <title>Application</title>
  </head>
  <body>
    <slot />
  </body>
</html>
```
### TSX (Root head in TSX missing blocking theme initializer):
```tsx
<head>
  <meta charSet="utf-8" />
  <title>Dashboard</title>
</head>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Blocking inline theme script in Astro head):
```astro
<head>
  <meta charset="utf-8" />
  <script is:inline>
    const theme = localStorage.getItem('theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    document.documentElement.classList.toggle('dark', theme === 'dark');
  </script>
</head>
```
### TSX (Blocking dangerouslySetInnerHTML theme script in TSX head):
```tsx
<head>
  <script
    dangerouslySetInnerHTML={{
      __html: "document.documentElement.classList.add(localStorage.getItem('theme') || 'light');",
    }}
  />
</head>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.hydration-theme-mismatch intentional exception -->
```

```tsx
// charites:ignore theme.hydration-theme-mismatch intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.hydration-theme-mismatch:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


