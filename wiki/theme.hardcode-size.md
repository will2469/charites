# theme.hardcode-size

> **Rule ID:** `theme.hardcode-size`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C DTCG Spatial Scale Standard, 8pt Modular Grid Rhythm, Tailwind CSS Spacing Architecture

---

## 1. Overview & Core Invariant

Detects hardcoded arbitrary size, spacing, or typography scalars in Tailwind utility classes

### Core Invariant:
> **"Spatial dimensions, spacing intervals, and typography sizes must use standardized modular scale tokens or CSS variables, never arbitrary raw scalar values or non-standard fractional steps."**

---
## 2. Technical Grounding & Engine Realities

Embedding arbitrary scalar dimensions (e.g. p-[19px], w-[320px], or text-[15px]) or non-standard fractional scales (e.g. p-3.25, w-2.75) introduces severe UI design regressions:

1. Spatial Rhythm Degradation: Arbitrary pixel/rem values and fractional decimals shatter the visual harmony of 4px/8px modular grid systems.
2. Sub-pixel Anti-Aliasing Blur: Fractional step dimensions like p-3.25 (13px) or w-2.75 (11px) fail to align cleanly on various mobile device pixel ratios (DPR), causing fuzzy borders and sub-pixel rendering artifacts.
3. False Conformance: Tailwind IntelliSense suggests shorthand decimals (e.g. 'p-[13px] can be written as p-3.25') which pass Tailwind validation without warning, but violate design system consistency.
4. Maintenance Overhead: Dispersed magic numbers make global layout scaling and responsive adaptation cumbersome.

Charites enforces migrating arbitrary sizing and non-standard fractional utilities to standard token steps (e.g. p-3, p-3.5, p-4, w-80, text-base) or semantic token variables.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Visual Rhythm Breakdown** | MEDIUM | Inconsistent micro-spacing across components causes fragmented alignment and sloppy UI rendering. |
| **Typography Scale Drift** | HIGH | Unchecked font sizes degrade readability, leading calculation, and accessibility scaling. |
| **Sub-pixel Rendering Artifacts** | MEDIUM | Non-standard off-grid decimal scales (e.g. 11px, 13px) produce rounding errors and blurry borders across fractional display scales. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Arbitrary padding, width, and non-standard fractional scale in JSX):
```tsx
<div className="p-[19px] w-[320px] p-3.25 w-2.75 text-[15px]">Hardcoded container</div>
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


