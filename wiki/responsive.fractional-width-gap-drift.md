# responsive.fractional-width-gap-drift

> **Rule ID:** `responsive.fractional-width-gap-drift`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** W3C CSS Flexible Box Layout Module Level 1 (Section 9.3: Line Breaking & Section 9.7: Resolving Flexible Lengths), W3C CSS Grid Layout Module Level 2 (Section 7.2: Flexible Track Sizing), Tailwind CSS Responsive Design & Spacing Geometry

---

## 1. Overview & Core Invariant

Warns when flex container declares a horizontal gap alongside fractional-width children totaling >= 100% with flex-wrap or disabled shrinking, causing unexpected wrapping or container blowout

### Core Invariant:
> **"Flex containers declaring horizontal gaps alongside fractional-width children totaling 100% or more must not enable 'flex-wrap' or disable flex shrinking ('shrink-0' / 'flex-none'), which causes unexpected line wrapping or container blowout."**

---
## 2. Technical Grounding & Engine Realities

In CSS Flexbox, percentage/fractional widths (e.g. 'w-1/2', 'w-2/3', 'w-1/3') are computed against the container's inner content box, but they do NOT deduct the horizontal gap space declared on the container.

When fractional children sum to 100% (e.g. 50% + 50% = 100%) and a positive horizontal gap is declared (e.g. 'gap-4' = 16px), the total required row width becomes 100% + 16px.

Under default 'flex-wrap: nowrap' and default 'flex-shrink: 1', browser Flexbox automatically shrinks items to absorb the 16px negative free space, making it visually safe. However, when 'flex-wrap' is active, the second item is forced onto a new line (unexpected line wrapping). Similarly, when flex shrinking is disabled on fractional items ('shrink-0' or 'flex-none'), the row refuses to shrink and overflows the container.

The modern, defect-free alternative is CSS Grid (e.g. 'grid grid-cols-2 gap-4' or 'grid grid-cols-1 md:grid-cols-3 gap-4') where fraction units ('fr') automatically distribute remaining track space after deducting gap sizes.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Unexpected Line Wrapping** | HIGH | When 'flex-wrap' is active, columns intended to appear side-by-side break awkwardly onto new lines because fractional widths + gap exceed 100%. |
| **Container Blowout & Horizontal Leak** | MEDIUM | When 'shrink-0' is applied to fractional flex items, the row cannot absorb gap spacing and blows out the container horizontally. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Flex container with flex-wrap and fractional children summing to 100%):
```tsx
<div className="flex flex-wrap gap-4">
  <div className="w-1/2">Kolom 1</div>
  <div className="w-1/2">Kolom 2 (terdorong ke baris baru karena gap-4)</div>
</div>
```
### TSX (Flex container with shrink-disabled fractional items causing horizontal blowout):
```tsx
<div className="flex gap-4">
  <div className="w-2/3 shrink-0">Kolom Utama</div>
  <div className="w-1/3 shrink-0">Sidebar</div>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (CSS Grid allocation where fraction units (fr) automatically deduct gap width):
```tsx
<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
  <div className="md:col-span-2">Kolom Utama</div>
  <div>Sidebar</div>
</div>
```
### TSX (Fluid flex distribution using flex-1 / basis-0):
```tsx
<div className="flex gap-4">
  <div className="flex-1 basis-0">Kolom 1</div>
  <div className="flex-1 basis-0">Kolom 2</div>
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore responsive.fractional-width-gap-drift intentional exception -->
```

```tsx
// charites:ignore responsive.fractional-width-gap-drift intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.fractional-width-gap-drift:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [responsive Category Guide](responsive).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


