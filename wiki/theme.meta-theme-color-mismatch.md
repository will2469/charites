# theme.meta-theme-color-mismatch

> **Rule ID:** `theme.meta-theme-color-mismatch`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** HTML Living Standard Section 4.2.5 (The meta element), Web App Manifest & Mobile OS Theme Integration, WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast)

---

## 1. Overview & Core Invariant

Detects static meta theme-color tags lacking media prefers-color-scheme queries

### Core Invariant:
> **"Meta theme-color elements must provide media query pairs (prefers-color-scheme: light/dark) to synchronize mobile browser chrome."**

---
## 2. Technical Grounding & Engine Realities

Modern mobile browsers (Safari on iOS, Chrome on Android) color the operating system status bar and address bar based on the <meta name="theme-color"> element in the document <head>.

When developers specify a single static theme-color without media queries (e.g. <meta name="theme-color" content="#ffffff">):
1. Blinding Address Bar: When the user toggles dark mode, the mobile address bar and status bar remain stark white, causing harsh visual glare.
2. Inverted Chrome: When the page switches to dark background, white status bar text collapses against white browser chrome.

Charites enforces declaring media="(prefers-color-scheme: light)" and media="(prefers-color-scheme: dark)" pairs on all meta theme-color definitions.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Chrome Glare** | MEDIUM | Mobile Safari and Chrome address bars blast high-brightness white chrome when the application is viewed in dark mode. |
| **Status Bar Text Invisibility** | LOW | OS status bar text (time, battery, Wi-Fi) becomes invisible due to poor contrast against unadapted address bar backgrounds. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Static theme-color meta tag in Astro layout):
```astro
<meta name="theme-color" content="#ffffff" />
```
### TSX (Static theme-color in TSX Document head):
```tsx
<meta name="theme-color" content="#09090b" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Adaptive light/dark meta theme-color pair in Astro):
```astro
<>
  <meta name="theme-color" media="(prefers-color-scheme: light)" content="#ffffff" />
  <meta name="theme-color" media="(prefers-color-scheme: dark)" content="#09090b" />
</>
```
### TSX (Adaptive meta theme-color pair in TSX):
```tsx
<head>
  <meta name="theme-color" media="(prefers-color-scheme: light)" content="#ffffff" />
  <meta name="theme-color" media="(prefers-color-scheme: dark)" content="#09090b" />
</head>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.meta-theme-color-mismatch intentional exception -->
```

```tsx
// charites:ignore theme.meta-theme-color-mismatch intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.meta-theme-color-mismatch:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


