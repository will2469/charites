# theme.hardcode-border-color

> **Rule ID:** `theme.hardcode-border-color`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C Design Tokens Community Group (DTCG), Tailwind CSS Border Token Architecture

---

## 1. Overview & Core Invariant

Detects hardcoded border and divider colors using primitive palettes, raw hex literals, or static monochrome

### Core Invariant:
> **"Component borders and dividers must use semantic tokens (border-border, border-input), never primitive palette or arbitrary hex colors."**

---
## 2. Technical Grounding & Engine Realities

Border lines define container elevation, separation, and affordance. When border colors are hardcoded (e.g. border-gray-200, border-[#e5e5e5]):

1. Invisibility in Dark Mode: A light gray border (#e5e5e5) provides zero contrast or turns into an inverted stark line in dark themes.
2. Theme Disconnect: When the primary or brand palette changes, borders remain pinned to legacy gray scales.
3. Inconsistent Boundaries: Disparate components end up using gray-200, slate-200, zinc-300 arbitrarily for identical UI dividers.

Charites enforces using centralized border tokens (border-border, border-input, divide-border).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Dark Mode Invisibility** | HIGH | Hardcoded light borders vanish or glow unnaturally on dark theme backgrounds. |
| **Visual Fragmentation** | MEDIUM | Different shades of gray borders destroy cohesive surface elevation hierarchy. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Hardcoded border primitive and arbitrary hex):
```astro
<div class="border border-gray-200 divide-y divide-[#e5e5e5]">List</div>
```
### TSX (Primitive directional border in JSX):
```tsx
export function Card() {
  return <div className="border-t-slate-300 border-x-[#cccccc]">Content</div>;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Using semantic border and divider tokens):
```astro
<div class="border border-border divide-y divide-border">List</div>
```
### TSX (Semantic border tokens with dark mode adaptability):
```tsx
export function Card() {
  return <div className="border-t border-border">Content</div>;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.hardcode-border-color intentional exception -->
```

```tsx
// charites:ignore theme.hardcode-border-color intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.hardcode-border-color:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


