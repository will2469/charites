# theme.apply-bloat

> **Rule ID:** `theme.apply-bloat`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** Tailwind CSS v3/v4 Architectural Best Practices, W3C Web Performance & CSS Bundle Size Guidelines

---

## 1. Overview & Core Invariant

Detects excessive use of @apply with more than 8 utility classes in CSS or style blocks

### Core Invariant:
> **"The @apply directive must not aggregate more than 8 utility classes in a single declaration to prevent CSS bloat and abstraction decay."**

---
## 2. Technical Grounding & Engine Realities

The @apply directive in Tailwind CSS was designed for small semantic abstractions (such as buttons or form inputs). Overusing @apply by stacking dozens of utility classes recreates the worst aspects of monolithic CSS.

When developers write @apply flex items-center justify-between p-4 bg-white rounded-lg shadow-md border border-gray-200 text-sm font-medium:
1. Bundle Size Inflation: Utility classes are duplicated into individual CSS selectors, negating Tailwind's atomic deduplication benefits.
2. Loss of Utility Ergonomics: Developers lose the ability to override individual styles via props or conditional classes.
3. Maintenance Decay: Giant @apply strings become unreadable 'css-in-css' dumping grounds.

Charites enforces a maximum threshold of 8 utility classes per @apply directive.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **CSS Bundle Bloat** | MEDIUM | Overloaded @apply directives balloon production stylesheet size and defeat atomic CSS compression. |
| **Component Maintainability Decay** | LOW | Massive CSS helper blocks reduce readability and make conditional variant overrides difficult. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Overloaded @apply declaration with 11 utility classes):
```astro
<style>
  .card {
    @apply flex items-center justify-between p-4 bg-white rounded-lg shadow-md border border-gray-200 text-sm font-medium;
  }
</style>
```
### TSX (Bloated @apply inside TSX style tag):
```tsx
export function Widget() {
  return (
    <style>{`
      .btn-primary {
        @apply inline-flex items-center justify-center px-4 py-2 text-sm font-semibold rounded-md shadow-sm text-white bg-primary hover:bg-primary/90;
      }
    `}</style>
  );
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Concise @apply declaration with 4 utility classes):
```astro
<style>
  .card {
    @apply flex items-center justify-between p-4;
  }
</style>
```
### TSX (Utilities applied directly to JSX markup):
```tsx
export function Card() {
  return <div className="flex items-center justify-between p-4 bg-white rounded-lg shadow-md">Markup Utilities</div>;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.apply-bloat intentional exception -->
```

```tsx
// charites:ignore theme.apply-bloat intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.apply-bloat:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


