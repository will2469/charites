# cls.unadjusted-font-metric

> **Rule ID:** `cls.unadjusted-font-metric`
> **Severity:** `INFO`
> **Category:** `cls`
> **Target Standards:** W3C CSS Fonts Module Level 4 (size-adjust, ascent-override, descent-override), Google Chrome Font Metric Override Guidelines

---

## 1. Overview & Core Invariant

Recommends font metric overrides on fallback font declarations to reduce swap CLS

### Core Invariant:
> **"Local fallback @font-face declarations (using 'src: local(...)') should specify metric adjustment descriptors ('size-adjust', 'ascent-override', or 'descent-override') to align bounding boxes with the principal web font."**

---
## 2. Technical Grounding & Engine Realities

When a web font downloads and replaces a system fallback font, variations in glyph x-height, ascent, and descent alter the computed bounding boxes of every text line.

This disparity causes sudden vertical expansion or contraction of paragraphs and headers, contributing to Cumulative Layout Shift.

By declaring 'size-adjust', 'ascent-override', and 'descent-override' on the fallback @font-face, developers can calibrate the system font's metrics to match the web font, creating a near zero-shift swap.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Font Swap Layout Jitter** | LOW | Paragraphs and navigation bars visibly shift lines when the web font swaps in. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Local fallback font-face without metric override descriptors):
```astro
<style>
  @font-face {
    font-family: 'InterFallback';
    src: local('Arial');
  }
</style>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Local fallback font-face with size-adjust and ascent-override):
```astro
<style>
  @font-face {
    font-family: 'InterFallback';
    src: local('Arial');
    ascent-override: 90%;
    descent-override: 22%;
    size-adjust: 107%;
  }
</style>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore cls.unadjusted-font-metric intentional exception -->
```

```tsx
// charites:ignore cls.unadjusted-font-metric intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.unadjusted-font-metric:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [cls Category Guide](cls).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


