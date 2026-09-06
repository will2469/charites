# cls.collapsible-height-jump

> **Rule ID:** `cls.collapsible-height-jump`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** CSS Grid Module Level 3 (grid-template-rows interpolation), Google Core Web Vitals (Interactive Animation CLS Invariants), Modern Zero-Shift Accordion Architectural Standards

---

## 1. Overview & Core Invariant

Collapsible accordion or drawer animates arbitrary max-height bounds instead of zero-shift CSS Grid

### Core Invariant:
> **"Collapsible content drawers and accordions must avoid animating arbitrary max-height bounds and instead adopt zero-shift CSS Grid (grid-template-rows: 0fr -> 1fr)."**

---
## 2. Technical Grounding & Engine Realities

A common legacy technique for animating collapsible elements involves transitioning 'max-height' from 0 to an arbitrarily large value (such as 'max-h-[1000px]').

Because CSS transitions interpolate over the declared boundary rather than actual content height, the animation duration becomes severely distorted: closing appears delayed and snapping occurs at the end of the transition, forcing layout reflow on surrounding elements.

The modern zero-shift solution utilizes CSS Grid: '<div class="grid transition-[grid-template-rows] duration-300 grid-rows-[0fr]"><div class="overflow-hidden">...</div></div>'. This allows CSS to smoothly interpolate intrinsic content height from 0fr to 1fr without any duration distortion or layout jumps.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Duration Distortion & Layout Snapping** | MEDIUM | Collapsing animations finish before the transition duration elapses, causing abrupt layout snaps and visual hitching. |
| **Continuous Main-Thread Reflow During Accordion Expansion** | MEDIUM | Transitioning max-height triggers continuous layout passes across all frames during accordion interactions. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Accordion drawer animating arbitrary max-height bounds):
```tsx
<div className="transition-all duration-300 overflow-hidden max-h-[1000px]">
  <AccordionBody />
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Modern zero-shift CSS Grid accordion interpolation):
```tsx
<div className="grid transition-[grid-template-rows] duration-300 grid-rows-[1fr]">
  <div className="overflow-hidden">
    <AccordionBody />
  </div>
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore cls.collapsible-height-jump intentional exception -->
```

```tsx
// charites:ignore cls.collapsible-height-jump intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.collapsible-height-jump:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [cls Category Guide](cls).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


