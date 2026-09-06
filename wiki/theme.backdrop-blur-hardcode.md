# theme.backdrop-blur-hardcode

> **Rule ID:** `theme.backdrop-blur-hardcode`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C Filter & Backdrop Filter Specification, Design System Glassmorphism Standards, Hardware-Accelerated Compositing Guidelines

---

## 1. Overview & Core Invariant

Detects hardcoded arbitrary blur and backdrop-blur scalars in Tailwind utility classes

### Core Invariant:
> **"Glassmorphism and surface blur effects must adhere to standardized blur tokens, never arbitrary scalar lengths."**

---
## 2. Technical Grounding & Engine Realities

Using arbitrary blur values (e.g. backdrop-blur-[5px] or blur-[12px]) produces inconsistent glassmorphism and performance bottlenecks:

1. GPU Overdraw Fragility: Arbitrary blur radii bypass optimized compositor layer pooling, causing unnecessary GPU rasterization penalties on mobile devices.
2. Glassmorphism Fragmentation: Slightly differing blur radii (e.g. 5px vs 8px vs 10px) across headers, dialogs, and drawer sheets ruin visual polish.
3. Inflexible Accessibility Adjustments: Standard tokens allow globally disabling or tuning blurs for users requesting reduced motion or low-power modes.

Charites enforces utilizing standard blur scale tokens (e.g. backdrop-blur-sm, backdrop-blur-md, backdrop-blur-lg) or CSS variables.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Glassmorphism Visual Discordance** | MEDIUM | Irregular blur intensity breaks cohesive layering and depth hierarchy across interface overlays. |
| **Mobile GPU Performance Stutter** | HIGH | Unstandardized backdrop-filter passes induce dropped frames during touch scrolling and bottom sheet gestures. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Arbitrary backdrop-blur on navigation header):
```tsx
<header className="backdrop-blur-[5px] bg-background/80">Sticky Nav</header>
```
### ASTRO (Arbitrary filter blur in Astro component):
```astro
<div class="blur-[12px] [backdrop-filter:blur(7px)]">Frosted Panel</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Using standard backdrop blur token):
```tsx
<header className="backdrop-blur-md bg-background/80">Sticky Nav</header>
```
### ASTRO (Standard filter blur token):
```astro
<div class="blur-md backdrop-blur-sm">Frosted Panel</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.backdrop-blur-hardcode intentional exception -->
```

```tsx
// charites:ignore theme.backdrop-blur-hardcode intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.backdrop-blur-hardcode:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


