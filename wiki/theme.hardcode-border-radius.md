# theme.hardcode-border-radius

> **Rule ID:** `theme.hardcode-border-radius`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C DTCG Shape & Radius Tokens, Design System Shape Hierarchy, Nested Curvature Optics Standard

---

## 1. Overview & Core Invariant

Detects hardcoded arbitrary border-radius scalars in Tailwind utility classes

### Core Invariant:
> **"Corner rounding and shape curvature must use standardized shape tokens or CSS variables, never arbitrary raw radius scalars."**

---
## 2. Technical Grounding & Engine Realities

Specifying arbitrary border-radius values (e.g. rounded-[7px] or rounded-t-[14px]) harms UI consistency:

1. Geometric Discordance: Components with off-scale radii look disjointed when nested or placed side-by-side.
2. Outer/Inner Radius Mismatch: Nested cards require deliberate radius proportion calculations (R_inner = R_outer - padding) defined by the shape system.
3. Rebranding Vulnerability: Global shape system updates (e.g. switching from square to rounded theme) cannot adapt arbitrary bracket classes.

Charites enforces using standard shape tokens (e.g. rounded-sm, rounded-md, rounded-xl) or token variables (e.g. rounded-[var(--radius)]).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Geometric Incoherence** | MEDIUM | Arbitrary corner radii make cards, buttons, and inputs clash visually across user interfaces. |
| **Theme Rigidity** | HIGH | Hardcoded radius prevents sweeping design system theme modernizations or brand shape adjustments. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Arbitrary border radius on button and card):
```tsx
<button className="rounded-[7px] p-3">Submit</button>
```
### ASTRO (Directional arbitrary radius in Astro component):
```astro
<div class="rounded-t-[14px] [border-radius:9px]">Modal Header</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Using design system shape tokens):
```tsx
<button className="rounded-md p-3">Submit</button>
```
### ASTRO (Standard directional rounded tokens):
```astro
<div class="rounded-t-xl rounded-b-none">Modal Header</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.hardcode-border-radius intentional exception -->
```

```tsx
// charites:ignore theme.hardcode-border-radius intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.hardcode-border-radius:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


