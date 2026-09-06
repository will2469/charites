# responsive.viewport-unit-leak

> **Rule ID:** `responsive.viewport-unit-leak`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** W3C CSS Values and Units Module Level 4 (Small, Large, and Dynamic Viewport Units), WebKit Safari iOS Dynamic Viewport Sizing Specification, Core Web Vitals Cumulative Layout Shift (CLS) Mitigation

---

## 1. Overview & Core Invariant

Warns when viewport height relies on static 100vh instead of modern dynamic dvh or svh units

### Core Invariant:
> **"Viewport height declarations should use CSS Level 4 dynamic units (dvh, svh) rather than static 100vh (h-screen, min-h-screen) to eliminate mobile layout shifts."**

---
## 2. Technical Grounding & Engine Realities

On mobile browsers (Safari iOS and Chrome Android), the browser address bar and bottom toolbar dynamically expand and collapse during vertical scrolling.

The classic 100vh unit uses the Large Viewport Height, which does not account for the visible URL bar. This causes bottom-anchored content to be covered by browser chrome and leads to disruptive layout jumps when the address bar toggles.

Utilizing dynamic viewport units (min-h-dvh or h-dvh) ensures the layout adapts smoothly to the actual visible viewport height.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Browser Layout Jumps (CLS)** | MEDIUM | Content suddenly shifts when mobile address bar hides or appears during scroll. |
| **Occluded Bottom CTA Buttons** | LOW | Bottom buttons in a 100vh container are partially covered beneath Safari's bottom navigation bar. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Static 100vh height causing layout jumps on mobile browsers):
```tsx
<main className="min-h-screen bg-background flex flex-col justify-between">
  <h1>Beranda Desa</h1>
  <button className="h-11 px-4 bg-primary text-primary-foreground">Lanjutkan</button>
</main>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Dynamic viewport height unit adapting smoothly to mobile address bar):
```tsx
<main className="min-h-dvh bg-background flex flex-col justify-between">
  <h1>Beranda Desa</h1>
  <button className="h-11 px-4 bg-primary text-primary-foreground">Lanjutkan</button>
</main>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore responsive.viewport-unit-leak intentional exception -->
```

```tsx
// charites:ignore responsive.viewport-unit-leak intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.viewport-unit-leak:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [responsive Category Guide](responsive).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


