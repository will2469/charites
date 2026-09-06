# responsive.desktop-only-content

> **Rule ID:** `responsive.desktop-only-content`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** Mobile-First Responsive Web Design (Luke Wroblewski), WCAG 2.2 Guideline 3.2 (Predictable - Consistent Navigation), Google Mobile-Friendly & Core Web Vitals Guidance

---

## 1. Overview & Core Invariant

Warns when primary action buttons or form submit controls are hidden on mobile viewports without mobile alternatives

### Core Invariant:
> **"Primary call-to-action (CTA) controls, checkout buttons, and form submissions must not be hidden on the mobile baseline ('hidden md:...') without accessible mobile parity."**

---
## 2. Technical Grounding & Engine Realities

In mobile-first design, hiding ancillary content (such as secondary marketing badges or decorative sidebars) on narrow viewports is standard practice.

However, hiding vital action triggers (e.g. 'Checkout', 'Bayar Sekarang', 'Kirim Berkas', or form submit buttons) via 'hidden md:flex' or 'hidden lg:block' leaves smartphone users stranded with incomplete user flows and broken core functionality.

Charites enforces that essential primary actions remain discoverable across all breakpoints, whether inline, within a bottom action sheet, or through a responsive floating bar.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile User Conversion Disruption** | HIGH | Smartphone users cannot complete essential transactions, checkouts, or form submissions when action buttons are hidden on mobile. |
| **Inconsistent Navigation Experience** | MEDIUM | Users switching between devices experience confusing functional disparity between desktop and mobile layouts. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Primary checkout button hidden on mobile baseline):
```tsx
<button type="submit" className="hidden md:flex items-center px-4 py-2 bg-primary text-primary-foreground">
  Bayar Sekarang
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Primary checkout button visible across all viewports with responsive sizing):
```tsx
<button type="submit" className="flex items-center justify-center w-full md:w-auto px-4 py-2 bg-primary text-primary-foreground">
  Bayar Sekarang
</button>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore responsive.desktop-only-content intentional exception -->
```

```tsx
// charites:ignore responsive.desktop-only-content intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.desktop-only-content:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [responsive Category Guide](responsive).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


