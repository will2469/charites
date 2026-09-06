# theme.dynamic-class

> **Rule ID:** `theme.dynamic-class`
> **Severity:** `ERROR`
> **Category:** `theme`
> **Target Standards:** Tailwind CSS JIT Static Analysis & Extraction Guidelines, Build-Time CSS Zero-Runtime Architecture, W3C Web Performance & Production Reliability

---

## 1. Overview & Core Invariant

Detects unpadded dynamic template strings breaking Tailwind JIT class generation

### Core Invariant:
> **"Utility classes must be written as complete static string literals so the Tailwind build compiler can reliably extract and generate them."**

---
## 2. Technical Grounding & Engine Realities

Tailwind CSS searches source files using regular expressions looking for complete class strings at build time. It does not evaluate JavaScript at runtime.

When developers dynamically construct utility classes using template literal slicing (e.g. className={`text-${color}-500`} or `bg-${variant}`):
1. Missing Production CSS: The Tailwind compiler never matches the interpolated string, leaving the utility completely absent from the production stylesheet.
2. Silent Visual Degradation: The component appears broken or unstyled in production while appearing to work intermittently in dev if another component imported that class.
3. Inscrutable Debugging: Developers struggle to trace why specific color variants intermittently fail to render.

Charites enforces using static class maps or complete utility strings within conditional expressions.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Missing Production Styles** | CRITICAL | Tailwind JIT engine strips un-scanned utility classes from production bundles, breaking layout and colors. |
| **Heisenbug UI Regressions** | HIGH | Styles intermittently vanish or break depending on which other files are compiled in the same build chunk. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Dynamic class string splicing in JSX className):
```tsx
<div className={`text-${color}-500 font-bold`}>Status</div>
```
### ASTRO (Dynamic background variant splicing in Astro):
```astro
<button class={`px-4 py-2 bg-${variant}`}>Action</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Static class lookup map for dynamic variants):
```tsx
const colorMap: Record<string, string> = {
  red: "text-red-500",
  blue: "text-blue-500",
  green: "text-green-500",
};
<div className={`${colorMap[color]} font-bold`}>Status</div>
```
### TSX (Complete utility class strings in ternary expression):
```tsx
<button className={`px-4 py-2 ${isActive ? "bg-primary text-primary-foreground" : "bg-muted"}`}>Action</button>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.dynamic-class intentional exception -->
```

```tsx
// charites:ignore theme.dynamic-class intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.dynamic-class:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


