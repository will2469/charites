# theme.chart-color-hardcode

> **Rule ID:** `theme.chart-color-hardcode`
> **Severity:** `ERROR`
> **Category:** `theme`
> **Target Standards:** WCAG 2.2 Success Criterion 1.4.3 (Contrast Minimum), WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast), Accessible Data Visualization Design Tokens

---

## 1. Overview & Core Invariant

Detects hardcoded color values on chart visualization components

### Core Invariant:
> **"Chart components must reference semantic theme tokens (e.g. var(--chart-1)) rather than hardcoded hex or color literals."**

---
## 2. Technical Grounding & Engine Realities

Data visualization libraries (such as Recharts, Chart.js, or Nivo) rely on SVG fill and stroke attributes to render bars, lines, and areas.

When developers hardcode hex colors onto chart elements (e.g. <Bar dataKey="sales" fill="#3b82f6" />):
1. Dark Mode Contrast Inversion: The hardcoded colors clash with dark card backgrounds, failing accessibility contrast minimums.
2. Theme Blindness: Visualizations fail to adapt when switching between light, dark, or high-contrast themes.
3. Fragmented Visual Identity: Brand colors drift between charts and surrounding interface tokens.

Charites enforces using CSS custom properties (fill="var(--chart-1)") or dynamic theme mappings.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Chart Contrast Invalidation** | HIGH | Chart bars and lines become illegible against inverted dark backgrounds, obscuring critical analytics. |
| **Theme Desynchronization** | MEDIUM | Data visualizations remain locked to legacy colors while the rest of the application adapts dynamically. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Hardcoded hex fill and stroke on Recharts Bar and Line):
```tsx
<>
  <Bar dataKey="revenue" fill="#3b82f6" />
  <Line dataKey="profit" stroke="#10b981" />
</>
```
### ASTRO (Hardcoded color on Area and Cell components):
```astro
<Area dataKey="uv" fill="#8884d8" stroke="#82ca9d" />
<Cell fill="#f43f5e" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Semantic chart tokens from design system):
```tsx
<>
  <Bar dataKey="revenue" fill="var(--chart-1)" />
  <Line dataKey="profit" stroke="var(--chart-2)" />
</>
```
### ASTRO (CSS variable references on Area and Cell):
```astro
<Area dataKey="uv" fill="var(--chart-1)" stroke="var(--chart-2)" />
<Cell fill="var(--chart-destructive)" />
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.chart-color-hardcode intentional exception -->
```

```tsx
// charites:ignore theme.chart-color-hardcode intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.chart-color-hardcode:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


