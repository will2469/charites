# theme.token-source-drift

> **Rule ID:** `theme.token-source-drift`
> **Severity:** `ERROR`
> **Category:** `theme`
> **Target Standards:** W3C Design Tokens Community Group (DTCG), Single Source of Truth (SSOT) Architecture

---

## 1. Overview & Core Invariant

Detects hardcoded color values bypassing the single source of truth design token pipeline

### Core Invariant:
> **"Custom properties representing theme tokens must not be assigned raw color literals in component scopes; they must resolve to SSOT token references."**

---
## 2. Technical Grounding & Engine Realities

Assigning raw hex/rgb color values directly to theme custom properties inside components or local stylesheets fractures the design token pipeline.

When developers write style="--primary: #2563eb" or declare local --color-brand: #3b82f6:
1. Drift from Global SSOT: The component diverges from centralized theme tokens (global.css), creating fragmented brand colors.
2. Theme Switching Failure: Dynamic theme changes (e.g. high-contrast, dark mode, multi-tenant branding) cannot override local hardcoded values.
3. Design System Audit Blind Spot: Design linters fail to track where rogue colors enter the application.

Charites enforces binding theme variables to global design tokens via var(--...) instead of raw literals.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Token SSOT Incoherence** | HIGH | Hardcoded local variable assignments decouple components from global design system updates. |
| **Theme Switch Blind Spot** | HIGH | Local variable assignments prevent dynamic color schemes and tenant styling from cascading. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Hardcoded hex assigned to theme token in inline style):
```astro
<div style="--primary: #2563eb; --background: #ffffff;">Drifting Tokens</div>
```
### TSX (Hardcoded rgb assigned to custom property in JSX style):
```tsx
export function Header() {
  return <header style={{ '--color-brand': 'rgb(37, 99, 235)' }}>Drifted Header</header>;
}
```
### ASTRO (Raw color assigned to theme custom property in style tag):
```astro
<style>
  .card {
    --card-bg: #1e293b;
  }
</style>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Theme token mapped via SSOT variable reference):
```astro
<div style="--primary: var(--color-blue-600);">SSOT Aligned</div>
```
### TSX (Non-color numeric custom property):
```tsx
export function Tabs() {
  return <div style={{ '--tab-index': '2' }}>Safe Property</div>;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.token-source-drift intentional exception -->
```

```tsx
// charites:ignore theme.token-source-drift intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.token-source-drift:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


