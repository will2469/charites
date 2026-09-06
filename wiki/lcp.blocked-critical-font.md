# lcp.blocked-critical-font

> **Rule ID:** `lcp.blocked-critical-font`
> **Severity:** `WARN`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Text Block Paint), W3C CSS Fonts Module Level 4 (font-display descriptor specification), Web Performance Working Group FOIT Minimization Guidelines

---

## 1. Overview & Core Invariant

Custom '@font-face' declaration lacks 'font-display: swap' or 'font-display: optional', risking FOIT (Flash of Invisible Text) and delaying LCP text paint

### Core Invariant:
> **"Custom @font-face declarations for web fonts must specify 'font-display: swap' or 'font-display: optional' to prevent Flash of Invisible Text (FOIT) on LCP text blocks."**

---
## 2. Technical Grounding & Engine Realities

When a browser discovers text styled with a custom web font, it evaluates the @font-face 'font-display' descriptor.

By default ('font-display: auto' or 'font-display: block'), modern browsers enter a 3-second 'block period' during which text is rendered with invisible transparent glyphs while the font binary is fetched from the network.

If the primary heading (<h1> or hero banner text) is the LCP candidate, this block period directly delays the Largest Contentful Paint until font download completes.

Specifying 'font-display: swap' enables immediate text rendering with a system fallback font followed by an in-place swap once the custom font arrives.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Flash of Invisible Text (FOIT)** | HIGH | LCP candidate text remains completely invisible for up to 3000ms on cellular or high-latency networks. |
| **Element Render Delay Inflation** | MEDIUM | Directly inflates LCP duration by coupling text paint to third-party or remote font network latency. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### CSS (Custom @font-face without font-display causing FOIT on hero headings):
```css
@font-face {
  font-family: 'CabinetGrotesk';
  src: url('/fonts/cabinet.woff2') format('woff2');
}
h1 {
  font-family: 'CabinetGrotesk', sans-serif;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### CSS (Custom @font-face configured with font-display: swap to ensure immediate text rendering):
```css
@font-face {
  font-family: 'CabinetGrotesk';
  src: url('/fonts/cabinet.woff2') format('woff2');
  font-display: swap;
}
h1 {
  font-family: 'CabinetGrotesk', sans-serif;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.blocked-critical-font intentional exception -->
```

```tsx
// charites:ignore lcp.blocked-critical-font intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.blocked-critical-font:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


