# theme.hardcode-color

> **Rule ID:** `theme.hardcode-color`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C Design Tokens Community Group (DTCG), Tailwind CSS Design Token Architecture, WCAG 2.2 Contrast Predictability

---

## 1. Overview & Core Invariant

Detects hardcoded arbitrary hex or rgb color literals in Tailwind utility classes and arbitrary properties

### Core Invariant:
> **"Color declarations in markup must use centralized semantic design tokens or CSS variables, never arbitrary raw hex or color function literals."**

---
## 2. Technical Grounding & Engine Realities

Directly embedding raw hex or rgb colors (e.g. bg-[#2563eb] or [color:#fff]) inside UI components creates serious maintenance barriers:

1. Theme Blindness: Arbitrary color values cannot respond to dark mode, high-contrast, or tenant theme switching.
2. Design Drift: Slight variations in hex codes (e.g. #2563eb vs #2564ea) fracture visual consistency.
3. Inflexible Rebranding: Global style updates require searching and replacing thousands of isolated class strings.

Charites enforces migrating arbitrary color literals to semantic tokens defined in global.css (e.g. bg-primary, text-card-foreground).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Theme Inflexibility** | HIGH | Hardcoded hex values remain static during dark mode toggle, causing illegible text and broken contrast. |
| **Maintenance Bloat** | MEDIUM | Scattered arbitrary colors prevent centralized palette changes and design system updates. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Arbitrary hex color in class attribute):
```astro
<div class="bg-[#1e293b] text-[#f8fafc] [color:#fff]">Un-tokenized Card</div>
```
### TSX (Arbitrary rgb and hex literals in JSX):
```tsx
export function Badge() {
  return <span className="hover:bg-[#2563eb] text-[rgb(255,0,0)]">Status</span>;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Using semantic tokens and CSS variables):
```astro
<div class="bg-card text-card-foreground">Tokenized Card</div>
```
### TSX (Semantic token utility with dark mode support):
```tsx
export function Badge() {
  return <span className="hover:bg-primary text-destructive">Status</span>;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.hardcode-color intentional exception -->
```

```tsx
// charites:ignore theme.hardcode-color intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.hardcode-color:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


