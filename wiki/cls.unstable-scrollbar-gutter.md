# cls.unstable-scrollbar-gutter

> **Rule ID:** `cls.unstable-scrollbar-gutter`
> **Severity:** `INFO`
> **Category:** `cls`
> **Target Standards:** W3C CSS Box Model Module Level 4 (scrollbar-gutter), Google Core Web Vitals (Horizontal Layout Shift Prevention), Desktop & Mobile Multi-Platform Viewport Consistency

---

## 1. Overview & Core Invariant

Root document scroller declares overflow-y: auto without scrollbar-gutter: stable, risking horizontal layout shifts

### Core Invariant:
> **"Root document scrollers (html, body, :root) with dynamic overflow should specify 'scrollbar-gutter: stable' to permanently reserve viewport space for scrollbars."**

---
## 2. Technical Grounding & Engine Realities

When a webpage renders with 'overflow-y: auto' at the document root level, the operating system initially renders content spanning the full viewport width.

As dynamic content streams in, hydration completes, or the user navigates to a longer page, the vertical scrollbar suddenly appears. On non-overlay desktop platforms (Windows, Linux, non-overlay macOS), this vertical scrollbar consumes 15-17px of width.

This sudden shrinkage of available client width causes all centered layouts, responsive grids, and full-bleed headers to snap and shift horizontally, registering an instant Cumulative Layout Shift.

Adding 'scrollbar-gutter: stable;' reserves the scrollbar space permanently, ensuring viewport dimensions remain completely invariant regardless of page height.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Document-Wide Horizontal Snapping** | LOW | Sudden appearance of vertical scrollbars causes all centered containers and flex items to jump 15-17px horizontally. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### CSS (Root html scroller with auto overflow but no reserved scrollbar gutter):
```css
html {
  overflow-y: auto;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### CSS (Root html scroller with stable scrollbar gutter reservation):
```css
html {
  overflow-y: auto;
  scrollbar-gutter: stable;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore cls.unstable-scrollbar-gutter intentional exception -->
```

```tsx
// charites:ignore cls.unstable-scrollbar-gutter intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.unstable-scrollbar-gutter:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [cls Category Guide](cls).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


