# theme.hardcode-size

> **Rule ID:** `theme.hardcode-size`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C DTCG Spatial Scale Standard, 8pt Modular Grid Rhythm, Tailwind CSS Spacing Architecture, Tailwind CSS Transform Scale Architecture

---

## 1. Overview & Core Invariant

Detects hardcoded arbitrary size, spacing, or typography scalars in Tailwind utility classes

### Core Invariant:
> **"Spatial dimensions, spacing intervals, typography sizes, and transform scale modifiers must use standardized modular scale tokens or CSS variables, never arbitrary raw scalar values, non-standard fractional steps, or inline arbitrary scale modifiers."**

---
## 2. Technical Grounding & Engine Realities

Embedding arbitrary scalar dimensions (e.g. p-[19px], w-[320px], or text-[15px]), non-standard fractional scales (e.g. p-3.25, w-2.75), or arbitrary inline scale modifiers (e.g. active:scale-[0.99], hover:scale-[1.02]) introduces severe UI design regressions:

1. Spatial Rhythm Degradation: Arbitrary pixel/rem values and fractional decimals shatter the visual harmony of 4px/8px modular grid systems.
2. Sub-pixel Anti-Aliasing Blur: Fractional step dimensions like p-3.25 (13px) or w-2.75 (11px) fail to align cleanly on various mobile device pixel ratios (DPR), causing fuzzy borders and sub-pixel rendering artifacts.
3. False Conformance: Tailwind IntelliSense suggests shorthand decimals (e.g. 'p-[13px] can be written as p-3.25') which pass Tailwind validation without warning, but violate design system consistency.
4. Maintenance Overhead: Dispersed magic numbers make global layout scaling and responsive adaptation cumbersome.
5. Transform Scale Drift & Subpixel Blurring: Arbitrary inline scale modifiers (e.g. active:scale-[0.99], hover:scale-[1.02], -scale-x-[0.95], [scale:0.98]) fragment tactile physics and micro-interactions across pages, trigger blurry text and 1px borders due to fractional rasterization, and violate centralized token governance.

Charites enforces migrating arbitrary sizing, non-standard fractional utilities, and arbitrary inline scales to standard token steps (e.g. p-3, p-3.5, p-4, w-80, text-base, scale-95, scale-105) or token-backed CSS variables (e.g. scale-[var(--scale-press)]).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Visual Rhythm Breakdown** | MEDIUM | Inconsistent micro-spacing across components causes fragmented alignment and sloppy UI rendering. |
| **Typography Scale Drift** | HIGH | Unchecked font sizes degrade readability, leading calculation, and accessibility scaling. |
| **Sub-pixel Rendering Artifacts** | MEDIUM | Non-standard off-grid decimal scales (e.g. 11px, 13px) produce rounding errors and blurry borders across fractional display scales. |
| **Physics & Motion Fragmentation** | MEDIUM | Arbitrary inline scale modifiers produce erratic tactile micro-interactions and blurry subpixel text rendering on standard DPI displays. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Arbitrary padding, width, and non-standard fractional scale in JSX):
```tsx
<div className="p-[19px] w-[320px] p-3.25 w-2.75 text-[15px]">Hardcoded container</div>
```
### TSX (Arbitrary inline transform scale modifiers in JSX):
```tsx
<div className="active:scale-[0.99] hover:scale-[1.02] -scale-x-[0.95] [scale:0.98]">Arbitrary scale</div>
```
### ASTRO (Arbitrary spacing and non-standard decimal in Astro component):
```astro
<section class="gap-1.25 mt-[27px] [padding:19px]">Arbitrary layout</section>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Standard modular scale tokens (integers and canonical .5 half-steps)):
```tsx
<div className="p-3.5 p-4 w-80 text-base">Standard modular container</div>
```
### TSX (Standard Tailwind scale steps and token-backed CSS variable):
```tsx
<div className="active:scale-95 hover:scale-105 active:scale-[var(--scale-press)]">Standard scale</div>
```
### ASTRO (System tokens and CSS variables):
```astro
<section class="gap-3 mt-6 p-4">Standard layout</section>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.hardcode-size intentional exception -->
```

```tsx
// charites:ignore theme.hardcode-size intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.hardcode-size:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


