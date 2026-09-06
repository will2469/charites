# theme.inline-style-hardcode

> **Rule ID:** `theme.inline-style-hardcode`
> **Severity:** `ERROR`
> **Category:** `theme`
> **Target Standards:** W3C Cascading Style Sheets (CSS) Level 3, W3C Design Tokens Community Group (DTCG)

---

## 1. Overview & Core Invariant

Detects hardcoded color literals inside HTML/JSX style attributes that prevent theme cascade

### Core Invariant:
> **"Color properties must not be declared as raw literals inside inline style attributes; they must use semantic classes or CSS variables."**

---
## 2. Technical Grounding & Engine Realities

Inline style attributes have the highest specificity in CSS, superseding all class selectors and theme cascades.

When developers write style="color: #2563eb" or style={{ background: '#fff' }}:
1. Impossible Dark Mode: The inline declaration cannot be targeted or overridden by .dark or [data-theme='dark'] class rules.
2. Broken Theming Pipeline: Token transformations (such as high-contrast mode or tenant styling) fail completely.
3. Maintenance Pitfall: Colors hidden in inline style strings avoid static analysis tools unless specifically parsed.

Charites enforces moving inline colors into utility classes or CSS custom properties.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Theme Specificity Lockout** | HIGH | Inline style specificity completely disables dark mode and stylesheet theming. |
| **Accessibility Barrier** | HIGH | High-contrast mode and accessibility themes cannot override inline hardcoded styles. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Hardcoded hex in HTML inline style):
```astro
<div style="color: #2563eb; background: #ffffff;">Inline Color</div>
```
### TSX (Hardcoded rgb in JSX style object):
```tsx
export function Card() {
  return <div style={{ color: '#2563eb', backgroundColor: 'rgb(255, 0, 0)' }}>Bad Style</div>;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Semantic utility classes instead of inline style):
```astro
<div class="text-primary bg-background">Themed Color</div>
```
### TSX (CSS variable in inline style for dynamic calculations):
```tsx
export function Card() {
  return <div style={{ color: 'var(--primary)' }}>Safe Style</div>;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.inline-style-hardcode intentional exception -->
```

```tsx
// charites:ignore theme.inline-style-hardcode intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.inline-style-hardcode:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


