# theme.hardcode-shadow-color

> **Rule ID:** `theme.hardcode-shadow-color`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C DTCG Elevation Tokens, Tailwind CSS Box Shadow Specification, Dark Mode Optical Physics & Contrast

---

## 1. Overview & Core Invariant

Detects hardcoded color literals embedded in box-shadow declarations

### Core Invariant:
> **"Elevation shadows must not embed raw hex or arbitrary color literals; shadow tints must adapt dynamically across light and dark modes via semantic tokens."**

---
## 2. Technical Grounding & Engine Realities

Embedding raw color literals inside arbitrary shadow brackets (e.g. shadow-[0_4px_10px_#00000040]) introduces major theme defects:

1. Dark Mode Disappearance: Dark shadows (black/gray with alpha) disappear completely when rendered over dark backgrounds (e.g. #09090b), leaving elevated cards looking flat.
2. Unadaptive Tints: Brand theme colors cannot tint shadows realistically when hardcoded hex codes are baked into individual classes.
3. Specificity Collisions: Overriding arbitrary shadow strings requires higher specificity or duplicate classes.

Charites enforces using standard shadow scale tokens (e.g. shadow-sm, shadow-md, shadow-lg) or semantic elevation tokens defined in global.css.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Dark Mode Elevation Invisibility** | HIGH | Hardcoded dark shadows become completely invisible against dark canvases, collapsing visual depth. |
| **Inconsistent Ambient Occlusion** | MEDIUM | Disparate shadow colors across components destroy uniform light-source perception in the design system. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Arbitrary shadow with embedded hex color):
```tsx
<div className="shadow-[0_4px_10px_#00000040] p-6">Floating Card</div>
```
### ASTRO (Arbitrary property box-shadow with rgb):
```astro
<section class="[box-shadow:0_10px_15px_rgba(0,0,0,0.1)]">Elevated Panel</section>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Using standard elevation shadow tokens):
```tsx
<div className="shadow-md p-6">Adaptive Floating Card</div>
```
### ASTRO (CSS variable shadow color):
```astro
<section class="shadow-[0_4px_6px_var(--shadow-color)]">Elevated Panel</section>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.hardcode-shadow-color intentional exception -->
```

```tsx
// charites:ignore theme.hardcode-shadow-color intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.hardcode-shadow-color:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


