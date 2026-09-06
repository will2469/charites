# theme.unlayered-token-definition

> **Rule ID:** `theme.unlayered-token-definition`
> **Severity:** `ERROR`
> **Category:** `theme`
> **Target Standards:** W3C Cascading Style Sheets (CSS) Level 5 (Cascade Layers), W3C CSS Custom Properties for Cascading Variables Module Level 1

---

## 1. Overview & Core Invariant

Detects CSS custom property definitions declared outside @layer theme or @layer base

### Core Invariant:
> **"CSS custom properties representing theme tokens must be declared within @layer theme or @layer base to ensure deterministic cascade resolution."**

---
## 2. Technical Grounding & Engine Realities

In modern frontend architectures and Tailwind CSS v4, unlayered CSS custom properties automatically take precedence over all layered styles regardless of specificity.

When developers declare :root { --primary: #... } without @layer theme or @layer base:
1. Cascade Inversion: Unlayered rules override framework layers and variant cascades unexpectedly.
2. Dark Mode Clashes: Nested dark mode themes defined within layers cannot reliably override unlayered root variables.
3. Specificity Pollution: Subsequent theme overrides require !important or higher specificity hacks to function.

Charites enforces encapsulating theme custom property definitions inside @layer theme or @layer base.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cascade Priority Inversion** | HIGH | Unlayered properties override all cascade layers, preventing dark mode and variant styles from taking effect. |
| **Theme Specificity Escalation** | MEDIUM | Teams resort to !important declarations to override unlayered variables, causing style degradation. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Unlayered :root custom property definition in style tag):
```astro
<style>
  :root {
    --primary: #2563eb;
    --background: #ffffff;
  }
</style>
```
### TSX (Unlayered [data-theme] custom property definition):
```tsx
export function GlobalStyles() {
  return (
    <style>{`
      :root {
        --brand-color: #3b82f6;
      }
    `}</style>
  );
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Enclosed within @layer theme):
```astro
<style>
  @layer theme {
    :root {
      --primary: #2563eb;
      --background: #ffffff;
    }
  }
</style>
```
### ASTRO (Enclosed within @layer base):
```astro
<style>
  @layer base {
    :root {
      --primary: #2563eb;
    }
  }
</style>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.unlayered-token-definition intentional exception -->
```

```tsx
// charites:ignore theme.unlayered-token-definition intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.unlayered-token-definition:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


