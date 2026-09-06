# theme.svg-hardcode-fill

> **Rule ID:** `theme.svg-hardcode-fill`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C SVG 2 Specification (Styling & currentColor), WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast), Design System Scalable Iconography Architecture

---

## 1. Overview & Core Invariant

Detects hardcoded color attributes on SVG markup preventing theme adaptation

### Core Invariant:
> **"SVG vector elements must inherit colors dynamically via currentColor or semantic CSS variables, never hardcoded hex or primitive colors."**

---
## 2. Technical Grounding & Engine Realities

Directly hardcoding raw colors onto SVG elements (such as <path fill="#000000"> or <stop stop-color="#3b82f6">) locks graphics to a static palette:

1. Theme Blindness: Dark icons with fill="#000" vanish when the user toggles dark mode.
2. Inverted Hover/Active States: Hardcoded stroke attributes prevent buttons and navigation links from changing icon color on hover or focus.
3. Reusability Breakdown: Components cannot share identical SVG glyphs across varying semantic surfaces without duplicating markup.

Charites enforces dynamic inheritance using fill="currentColor", stroke="currentColor", or semantic design tokens (var(--primary)).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Dark Mode Icon Invisibility** | HIGH | Vector icons hardcoded to black or dark shades become completely invisible against dark backgrounds. |
| **Broken State Affordance** | MEDIUM | Icons fail to inherit hover, focus, and disabled states from parent interactive components. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Hardcoded hex fill on SVG path in TSX):
```tsx
<path fill="#000000" d="M10 10 H 90 V 90 H 10 Z" />
```
### ASTRO (Primitive hex stop-color and stroke in Astro SVG):
```astro
<svg viewBox="0 0 100 100">
  <stop stop-color="#3b82f6" offset="100%" />
  <circle cx="50" cy="50" r="40" stroke="#ef4444" fill="none" />
</svg>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Adaptive currentColor fill in TSX):
```tsx
<path fill="currentColor" d="M10 10 H 90 V 90 H 10 Z" />
```
### ASTRO (Dynamic CSS variable in gradient stop and currentColor stroke):
```astro
<svg viewBox="0 0 100 100">
  <stop stop-color="var(--primary)" offset="100%" />
  <circle cx="50" cy="50" r="40" stroke="currentColor" fill="none" />
</svg>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.svg-hardcode-fill intentional exception -->
```

```tsx
// charites:ignore theme.svg-hardcode-fill intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.svg-hardcode-fill:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


