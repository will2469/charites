# responsive.dynamic-viewport-inconsistency

> **Rule ID:** `responsive.dynamic-viewport-inconsistency`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** W3C CSS Values and Units Module Level 4 (Small, Large, and Dynamic Viewport Units), WebKit Dynamic Viewport Sizing Specification, Chrome for Android URL Bar Scroll Resize Guidelines

---

## 1. Overview & Core Invariant

Warns when static viewport units (100vh, h-screen) and modern dynamic units (dvh, svh) are mixed inconsistently across layout hierarchies

### Core Invariant:
> **"Components nested within a dynamic viewport container ('dvh', 'svh') must not use static viewport units ('100vh', 'h-screen'), and conflicting viewport dimensions must not be declared on the same element."**

---
## 2. Technical Grounding & Engine Realities

Modern mobile browsers (Safari iOS and Chrome Android) feature dynamic interface chrome (URL address bar and bottom navigation toolbar) that expand and collapse during user scrolling.

The dynamic viewport unit 'dvh' continuously tracks the active visible viewport height. In contrast, classical '100vh' and 'h-screen' are fixed to the Large Viewport (the maximum screen height assuming all browser chrome is collapsed).

When an outer wrapper uses 'min-h-dvh' while an inner component specifies 'h-screen' or 'h-[100vh]', the child height exceeds the visible parent area whenever the address bar is visible, causing unexpected double scrollbars, layout clipping, and jarring viewport jitter.

Charites enforces consistent viewport unit adoption across component hierarchies.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Double Scrollbar & Viewport Jitter** | MEDIUM | Inner components sized with 100vh exceed the dvh container, causing double scrollbars and layout jerking during scroll. |
| **Content Clipping Behind Browser Chrome** | LOW | Bottom actions and footers are pushed offscreen beneath mobile browser toolbars. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Inner child with h-screen nested inside an outer min-h-dvh container):
```tsx
<main className="min-h-dvh flex flex-col">
  <div className="h-screen bg-surface">
    <h2>Konten Terpotong di Mobile</h2>
  </div>
</main>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Consistent dynamic viewport units across parent and child):
```tsx
<main className="min-h-dvh flex flex-col">
  <div className="h-full bg-surface">
    <h2>Konten Selaras Mengikuti Viewport</h2>
  </div>
</main>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore responsive.dynamic-viewport-inconsistency intentional exception -->
```

```tsx
// charites:ignore responsive.dynamic-viewport-inconsistency intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.dynamic-viewport-inconsistency:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [responsive Category Guide](responsive).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


