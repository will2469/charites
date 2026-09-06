# theme.split-theme-state

> **Rule ID:** `theme.split-theme-state`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** Design System Single Source of Truth (SSOT) Architecture, React State Management Best Practices (Context & Hooks), WCAG 2.2 Predictable Navigation & Consistency

---

## 1. Overview & Core Invariant

Detects ad-hoc direct access to theme state via localStorage outside ThemeProvider

### Core Invariant:
> **"Component UI state must consume theme through a unified ThemeProvider context or custom hook, never querying localStorage directly in component bodies or handlers."**

---
## 2. Technical Grounding & Engine Realities

When developers directly access or mutate localStorage.getItem('theme') or localStorage.theme in scattered components:

1. Fragmented State: Component A reads localStorage while Component B listens to React Context, causing disparate parts of the UI to display inconsistent themes.
2. Missing Reactivity: Updates directly to localStorage do not trigger React or framework re-renders across sibling components.
3. Testability Breakdown: Components cannot be unit tested or rendered in isolation without mocking global browser APIs.

Charites enforces routing all theme state access through a unified Theme Provider / useTheme hook, permitting direct localStorage access only in root <head> bootstrap scripts.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Theme Desynchronization Across UI** | MEDIUM | Different page regions display discordant color schemes due to uncoordinated local state reads. |
| **Broken Component Reactivity** | MEDIUM | Theme switches fail to re-render affected components without a full browser page refresh. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Direct localStorage mutation in button onClick handler):
```tsx
<button onClick={() => localStorage.setItem('theme', 'dark')}>Toggle</button>
```
### ASTRO (Direct localStorage inspection in Astro component body):
```astro
<div data-theme={localStorage.getItem('theme')}>Container</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Using unified useTheme hook from ThemeProvider):
```tsx
const { theme, setTheme } = useTheme();
<button onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}>Toggle</button>
```
### ASTRO (Permitted inline bootstrap script inside root <head>):
```astro
<head>
  <script is:inline>
    const theme = localStorage.getItem('theme') || 'light';
    document.documentElement.classList.toggle('dark', theme === 'dark');
  </script>
</head>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.split-theme-state intentional exception -->
```

```tsx
// charites:ignore theme.split-theme-state intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.split-theme-state:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


