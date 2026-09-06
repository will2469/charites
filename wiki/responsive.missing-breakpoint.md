# responsive.missing-breakpoint

> **Rule ID:** `responsive.missing-breakpoint`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** Mobile-First Responsive Web Design Principles, W3C CSS Grid Layout Module Level 2, Tailwind CSS Responsive Design Specification

---

## 1. Overview & Core Invariant

Warns when multi-column grids or giant font sizes are declared on mobile baseline without responsive breakpoint modifiers

### Core Invariant:
> **"Multi-column grids (grid-cols-[3-12]) and giant font sizes (text-[5-9]xl) must not be defined on mobile baseline without responsive breakpoint prefixes (sm:, md:, lg:)."**

---
## 2. Technical Grounding & Engine Realities

On compact smartphone screens (360px-390px), defining 3 or more columns directly on mobile baseline squeezes individual columns below 100px width, causing severe card distortion and text wrapping.

Similarly, declaring giant typography (e.g. text-6xl) on mobile baseline causes single words to span multiple vertical lines, breaking header visual balance.

Adhering to mobile-first progression requires starting from single-column baselines (grid-cols-1) and scaling up via responsive modifiers (sm:grid-cols-2 md:grid-cols-4).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Severe Column Squeeze on Mobile** | MEDIUM | Multi-column cards become unreadable and distorted when squeezed into 360px phone screens. |
| **Typography Layout Blowout** | LOW | Giant font headings wrap awkwardly into 4-5 lines on narrow mobile viewports. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Multi-column grid on mobile baseline without responsive modifier):
```tsx
<div className="grid grid-cols-4 gap-4">
  <div className="bg-card p-4">Item 1</div>
  <div className="bg-card p-4">Item 2</div>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Mobile-first progression starting from 1 column to multi-column on desktop):
```tsx
<div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
  <div className="bg-card p-4">Item 1</div>
  <div className="bg-card p-4">Item 2</div>
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore responsive.missing-breakpoint intentional exception -->
```

```tsx
// charites:ignore responsive.missing-breakpoint intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.missing-breakpoint:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [responsive Category Guide](responsive).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


