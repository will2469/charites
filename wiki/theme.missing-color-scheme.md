# theme.missing-color-scheme

> **Rule ID:** `theme.missing-color-scheme`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C CSS Color Adjustment Module Level 1, HTML Living Standard Section 4.2.5.5 (color-scheme)

---

## 1. Overview & Core Invariant

Detects dark theme definitions (.dark, [data-theme="dark"]) missing color-scheme property

### Core Invariant:
> **"Dark mode theme selectors must declare color-scheme: dark to synchronize native browser UI elements with the theme."**

---
## 2. Technical Grounding & Engine Realities

When developers configure dark mode using .dark or [data-theme='dark'] CSS rules without declaring color-scheme: dark:

1. White Form Controls: Native form elements (<select> dropdown popovers, <input type='date'> calendars, checkboxes) remain bright white.
2. Inverted Scrollbars: Operating system scrollbars fail to enter dark mode, glaring against dark page content.
3. Blinding Autofill Scrims: Browser credential autofill backgrounds turn blinding yellow-white.

Charites enforces declaring color-scheme: dark within all dark mode selector scopes.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Native UI Flash** | MEDIUM | Native controls, scrollbars, and dropdowns render in high-contrast light theme inside dark mode. |
| **Form Usability Degradation** | LOW | Datepickers and select menus become illegible due to mismatched browser chrome defaults. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Dark mode class without color-scheme declaration):
```astro
<style>
  .dark {
    --background: #09090b;
    --foreground: #fafafa;
  }
</style>
```
### TSX (Data-theme attribute without color-scheme):
```tsx
export function DarkTheme() {
  return (
    <style>{`
      [data-theme="dark"] {
        --bg-main: #121212;
      }
    `}</style>
  );
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Dark selector with explicit color-scheme):
```astro
<style>
  .dark {
    color-scheme: dark;
    --background: #09090b;
    --foreground: #fafafa;
  }
</style>
```
### ASTRO (Global color-scheme on root):
```astro
<style>
  :root {
    color-scheme: light dark;
  }
</style>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.missing-color-scheme intentional exception -->
```

```tsx
// charites:ignore theme.missing-color-scheme intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.missing-color-scheme:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


