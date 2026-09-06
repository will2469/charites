# theme.gradient-hardcode

> **Rule ID:** `theme.gradient-hardcode`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C Design Tokens Community Group (DTCG), Tailwind CSS Gradient Token Architecture

---

## 1. Overview & Core Invariant

Detects hardcoded primitive, arbitrary hex, or monochrome colors in gradient stops

### Core Invariant:
> **"Gradient color stops must use semantic tokens (from-primary, to-accent), never primitive palette or arbitrary hex stops."**

---
## 2. Technical Grounding & Engine Realities

Gradients often span large hero sections or callout backgrounds. When color stops use primitive or arbitrary values (e.g. from-[#3b82f6] to-blue-500):

1. Inverted Muddy Colors: Light mode gradients rendered in dark mode produce muddy, low-contrast, or unreadable backgrounds behind text.
2. Theme Decoupling: Rebranding or dynamic tenant themes cannot adjust the stops without manually updating every gradient class.
3. Accessibility Violations: Static gradient stops cannot guarantee compliance with WCAG 2.2 text contrast across all screen areas.

Charites enforces gradient stops constructed from semantic tokens (from-primary, to-secondary, via-accent, from-transparent).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Dark Mode Breakage** | HIGH | Hardcoded gradient stops destroy text legibility and brand alignment in dark themes. |
| **Design Token Fragmentation** | MEDIUM | Gradients drift out of sync with established design system tokens. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Gradient stops using arbitrary hex and primitive colors):
```astro
<div class="bg-gradient-to-r from-[#3b82f6] to-blue-500">Banner</div>
```
### TSX (Gradient stops using monochrome white and primitive red):
```tsx
export function Hero() {
  return <div className="bg-gradient-to-b from-white via-rose-500 to-black">Hero</div>;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Semantic tokens for gradient stops):
```astro
<div class="bg-gradient-to-r from-primary to-accent">Banner</div>
```
### TSX (Semantic tokens adapting cleanly to dark mode):
```tsx
export function Hero() {
  return <div className="bg-gradient-to-b from-card via-primary to-background">Hero</div>;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.gradient-hardcode intentional exception -->
```

```tsx
// charites:ignore theme.gradient-hardcode intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.gradient-hardcode:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


