# theme.pseudo-hardcode-color

> **Rule ID:** `theme.pseudo-hardcode-color`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** W3C Design Tokens Community Group (DTCG), Tailwind CSS Pseudo-Element Architecture

---

## 1. Overview & Core Invariant

Detects hardcoded primitive, arbitrary hex, or monochrome colors inside pseudo-element and pseudo-class variants

### Core Invariant:
> **"Pseudo-elements (placeholder, selection, file, marker) must consume semantic tokens, never raw primitive or arbitrary colors."**

---
## 2. Technical Grounding & Engine Realities

Pseudo-element styling often slips past generic linters that only inspect top-level classes.

When developers specify placeholder:text-gray-400 or selection:bg-blue-200:
1. Input Readability Degradation: A placeholder styled with light gray-400 becomes completely invisible on light input surfaces or garish on dark inputs.
2. Selection Contrast Clashes: Static blue-200 selection background can fail WCAG contrast against the text color in dark mode.
3. Inconsistent State Branding: File inputs and list markers fail to reflect global theme tokens.

Charites enforces using semantic tokens (placeholder:text-muted-foreground, selection:bg-primary-light).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Form Accessibility Failure** | HIGH | Low-contrast placeholder text fails WCAG minimum ratio, making form inputs confusing for users. |
| **Selection Highlight Glitch** | MEDIUM | Hardcoded selection backgrounds obliterate text visibility under dark themes. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Hardcoded primitive colors in placeholder and selection):
```astro
<input class="placeholder:text-gray-400 selection:bg-blue-200" />
```
### TSX (Arbitrary hex in pseudo variants):
```tsx
export function Input() {
  return <input className="placeholder:text-[#94a3b8] file:bg-slate-100" />;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Semantic tokens for pseudo styling):
```astro
<input class="placeholder:text-muted-foreground selection:bg-primary-light" />
```
### TSX (Semantic tokens adapting to active theme):
```tsx
export function Input() {
  return <input className="placeholder:text-muted-foreground file:bg-secondary" />;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.pseudo-hardcode-color intentional exception -->
```

```tsx
// charites:ignore theme.pseudo-hardcode-color intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.pseudo-hardcode-color:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


