# theme.no-reduced-motion

> **Rule ID:** `theme.no-reduced-motion`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** WCAG 2.2 Success Criterion 2.3.3 (Animation from Interactions), W3C Media Queries Level 5 (prefers-reduced-motion), Accessible Web Animation & Vestibular Safety Guidelines

---

## 1. Overview & Core Invariant

Detects global theme transitions without prefers-reduced-motion media query wrapping

### Core Invariant:
> **"Global theme and color transitions must be scoped within prefers-reduced-motion: no-preference or mitigated with reduced-motion overrides."**

---
## 2. Technical Grounding & Engine Realities

Smooth CSS transitions applied to root or theme switching (such as * { transition: background-color 0.3s, color 0.3s; } or transition: all 0.2s) can cause dizziness, headaches, and nausea for users with vestibular disorders.

WCAG 2.2 Success Criterion 2.3.3 requires that non-essential animations triggered by user interaction can be turned off or respect system accessibility preferences.

Charites enforces wrapping theme transitions in @media (prefers-reduced-motion: no-preference) or providing an explicit @media (prefers-reduced-motion: reduce) { transition: none; } fallback.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Vestibular Distress** | MEDIUM | Rapid or uncontrolled surface transitions induce disorientation or motion sickness for sensitive users. |
| **WCAG 2.2 SC 2.3.3 Non-Compliance** | MEDIUM | Failure to honor OS-level accessibility preferences prevents compliance with regulatory accessibility standards. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Unmitigated global theme transition in Astro style):
```astro
<style>
  * {
    transition: background-color 0.3s ease, color 0.3s ease;
  }
</style>
```
### TSX (Broad transition all without motion preference in TSX):
```tsx
<style>{`
  body {
    transition: all 0.25s ease-in-out;
  }
`}</style>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Theme transition scoped to no-preference media query):
```astro
<style>
  @media (prefers-reduced-motion: no-preference) {
    * {
      transition: background-color 0.3s ease, color 0.3s ease;
    }
  }
</style>
```
### ASTRO (Explicit reduced-motion override):
```astro
<style>
  body {
    transition: background-color 0.3s ease;
  }
  @media (prefers-reduced-motion: reduce) {
    body {
      transition: none;
    }
  }
</style>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.no-reduced-motion intentional exception -->
```

```tsx
// charites:ignore theme.no-reduced-motion intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.no-reduced-motion:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


