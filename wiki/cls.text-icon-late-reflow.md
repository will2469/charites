# cls.text-icon-late-reflow

> **Rule ID:** `cls.text-icon-late-reflow`
> **Severity:** `INFO`
> **Category:** `cls`
> **Target Standards:** W3C Cumulative Layout Shift (CLS) Metric Specification, Google Material Icons & Material Symbols Integration Guide

---

## 1. Overview & Core Invariant

Requires locked bounding dimensions on text-ligature icon elements to prevent text reflow

### Core Invariant:
> **"Text-ligature icon elements must lock their bounding box via 'inline-block' (or block/flex), explicit width/height (or 'size-*'), and 'overflow-hidden' to prevent raw word ligature text from expanding the layout before the icon font is loaded."**

---
## 2. Technical Grounding & Engine Realities

Icon fonts like Material Icons or Material Symbols render icons by substituting raw text strings (such as 'shopping_cart', 'account_circle', or 'arrow_back') with icon glyphs via OpenType ligatures.

Before the web font finishes downloading, the browser displays the fallback word text ('shopping_cart') at full length (spanning 80-120px).

When the web font suddenly loads, the word shrinks into a 24x24px glyph, causing surrounding navigation bars, buttons, and text to collapse backward and triggering Cumulative Layout Shift (CLS).

Locking the container dimensions to 'inline-block size-6 overflow-hidden' ensures the element occupies exactly 24x24px regardless of font loading state.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Button and Header Layout Contraction** | MEDIUM | Long ligature strings expand buttons initially, then contract suddenly when font arrives. |
| **Cumulative Layout Shift (CLS)** | LOW | Shifts around interactive icons trigger layout recalculations in headers and toolbars. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Material icon with raw text ligature without locked box dimensions):
```tsx
<button className="flex items-center gap-2">
  <span className="material-icons">shopping_cart</span> Beli
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Locked icon bounding box with inline-block, size-6, and overflow-hidden):
```tsx
<button className="flex items-center gap-2">
  <span className="material-icons inline-block size-6 overflow-hidden">shopping_cart</span> Beli
</button>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore cls.text-icon-late-reflow intentional exception -->
```

```tsx
// charites:ignore cls.text-icon-late-reflow intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.text-icon-late-reflow:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [cls Category Guide](cls).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


