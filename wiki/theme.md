# Theme Rules (`theme`)

The `theme` category enforces design token compliance, color system integrity, OKLCH color spaces, and consistent visual styling across Astro components and React TSX templates.

---

## Category Rule Index

| Rule ID | Severity | Summary | Status |
| :--- | :---: | :--- | :---: |
| [`theme.hardcode-opacity-color`](#themehardcode-opacity-color) | `ERROR` | Detects utility classes with hardcoded slash opacity modifiers that have official semantic token replacements | `enabled` |

---

## `theme.hardcode-opacity-color`

> **Rule ID:** `theme.hardcode-opacity-color`
> **Severity:** `ERROR`
> **Category:** `theme`
> **Target Standards:** W3C Design Tokens Community Group (DTCG), Tailwind CSS Design Token Architecture, WCAG 2.2 Relative Contrast

### 1. Overview & Core Invariant

The `theme.hardcode-opacity-color` rule forbids ad-hoc slash opacity modifiers on base semantic colors (such as `bg-primary/10`, `text-primary/20`, or `border-destructive/20`) when an officially curated semantic replacement token exists in the design system.

#### Core Invariant:
> **"Every color opacity variation that represents a semantic state or visual elevation must use a centralized semantic design token rather than an arbitrary slash modifier."**

Direct slash opacity modifiers bypass CSS custom property definitions, break automated theme switching (e.g. Dark Mode contrast recalculation), and cause visual drift across components.

### 2. Technical Grounding & Engine Realities

In modern design token architecture (such as Tailwind CSS with CSS Variables or OKLCH color spaces), semantic colors like `primary` and `destructive` are calibrated for foreground/background contrast against explicit color stops.

When developers append arbitrary slash modifiers (e.g., `bg-primary/10`), the resulting alpha-blended color:
1. **Destroys WCAG 2.2 Contrast Predictability:** Transparent alpha layers depend on whatever background color sits underneath. In dark mode or high-contrast themes, 10% opacity can drop contrast ratios below the 4.5:1 WCAG AA minimum.
2. **Breaks Theme Export & Reusability:** When exporting design tokens to mobile apps, Figma, or print styles, runtime alpha calculations cannot be resolved statically.
3. **Creates Aesthetic Inconsistency:** Different developers use varying opacities (`/5`, `/10`, `/15`, `/20`) for the same intended visual state (such as subtle hover backgrounds or tinted badge pills).

Charites enforces pre-calibrated semantic tokens (e.g. `primary-light`, `primary-subtle`, `muted-light`, `destructive-light`) that are mathematically verified for contrast and consistent across themes.

### 3. Static Analysis Architecture & AST Detection

Charites inspects unified Leaf IR (`*ir.Node`) nodes across Astro components (`.astro`) and React JSX/TSX files (`.tsx`, `.jsx`):

1. **Class Tokenization:** Each element's `class` or `className` attribute is extracted into normalized tokens.
2. **Fast-Path Filter:** Tokens without a slash (`'/'`) are skipped immediately with $O(1)$ complexity, allocating zero heap memory.
3. **Lexical Normalization (`stripVariants`):** Single and chained variants (such as `hover:`, `dark:`, `md:hover:`) are stripped down to the root utility.
4. **Utility Family Match:** Only color utility families (`bg-`, `text-`, `border-`, `ring-`) are evaluated.
5. **Token Map Lookup:** The extracted color and opacity modifier are checked against `OpacityTokenMap`.

#### Official Semantic Replacement Matrix (`OPACITY_TOKEN_MAP`):

| Violating Pattern | Replacement Semantic Token | Visual Intent |
| :--- | :--- | :--- |
| `primary/10`, `primary/20` | `primary-light` | Light primary tint (hover, active card) |
| `primary/5` | `primary-subtle` | Subtle tinted background |
| `secondary/10`, `muted/10` | `muted-light` | Subtle muted surface |
| `secondary/5`, `muted/5` | `muted-subtle` | Very subtle container tint |
| `destructive/10`, `destructive/20` | `destructive-light` | Light error/alert background |
| `destructive/5` | `destructive-subtle` | Subtle destructive surface |
| `accent/10`, `accent/20` | `accent-light` | Accent highlight tint |
| `accent/5` | `accent-subtle` | Subtle accent container |
| `warning/10` | `warning-light` | Warning alert tint |
| `warning/5` | `warning-subtle` | Subtle warning surface |
| `amber/10` | `amber-light` | Amber pill/badge tint |
| `amber/5` | `amber-subtle` | Subtle amber surface |
| `emerald/10` | `emerald-light` | Success badge tint |
| `emerald/5` | `emerald-subtle` | Subtle success surface |

#### Excluded Utilities (Out-of-Scope):
The rule deliberately ignores non-color fractions and utilities that belong to other domains:
- Layout fractions: `w-1/2`, `h-1/3`, `max-w-1/2`
- Aspect ratios: `aspect-16/9`, `aspect-4/3`
- Grid template fractions: `grid-cols-2/3`
- Typography line-height modifiers: `text-sm/6`, `text-xs/relaxed`
- Arbitrary hex colors: `bg-[#123456]/10` (handled by `theme.hardcode-color`)
- Raw palette colors: `bg-red-500/10` (handled by `theme.hardcode-palette-color`)

### 4. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Accessibility Degradation** | HIGH | Contrast ratio drops below 4.5:1 under dark mode themes due to uncalibrated alpha blending. |
| **Visual Debt & Inconsistency** | MEDIUM | Proliferation of slightly different opacities (`/5`, `/10`, `/20`) degrades product polish. |
| **Theme Portability Failure** | MEDIUM | External design token exporters cannot map hardcoded alpha values to standalone color systems. |

### 5. Non-Compliant Code Patterns (Bad Examples)

#### Astro (`.astro`):
```astro
---
// BAD: Direct slash opacity modifiers on semantic colors
---
<div class="card p-6 rounded-xl bg-primary/10 border border-destructive/20">
  <h2 class="text-xl font-bold text-primary/20">Card Title</h2>
  <span class="badge ring-1 ring-warning/10 bg-primary/5">Warning</span>
</div>
```

#### React TSX (`.tsx`):
```tsx
// BAD: Chained and single variants with hardcoded opacity
export function ActionCard() {
  return (
    <div className="p-4 rounded-lg hover:bg-primary/10 dark:bg-primary/10 md:hover:bg-primary/10">
      <button className="px-3 py-2 text-sm dark:border-destructive/20 sm:dark:hover:border-destructive/20">
        Delete
      </button>
    </div>
  );
}
```

### 6. Compliant Implementation Patterns (Good Examples)

#### Astro (`.astro`):
```astro
---
// GOOD: Using official semantic tokens from global.css
---
<div class="card p-6 rounded-xl bg-primary-light border border-destructive-light">
  <h2 class="text-xl font-bold text-primary">{Astro.props.title}</h2>
  <span class="badge ring-1 ring-warning-light bg-primary-subtle">Warning</span>
</div>
```

#### React TSX (`.tsx`):
```tsx
// GOOD: Using semantic tokens with variants
export function ActionCard() {
  return (
    <div className="p-4 rounded-lg hover:bg-primary-light dark:bg-primary-light md:hover:bg-primary-light">
      <button className="px-3 py-2 text-sm dark:border-destructive-light">
        Delete
      </button>
    </div>
  );
}
```

### 7. How to Suppress (Ignore Directives)

If a third-party integration or specific one-off visual effect intentionally requires a hardcoded opacity modifier, suppress the diagnostic using the canonical Charites Rule ID:

#### Single-Line Ignore (Astro / HTML):
```astro
<!-- charites:ignore theme.hardcode-opacity-color intentional brand exception -->
<div class="bg-primary/10">Special Brand Banner</div>
```

#### Single-Line Ignore (React TSX / JSX):
```tsx
// charites:ignore theme.hardcode-opacity-color intentional brand exception
<div className="bg-primary/10">Special Brand Banner</div>
```

#### Block-Range Ignore (Astro / TSX):
```astro
<!-- charites:ignore-start theme.hardcode-opacity-color -->
<div class="bg-primary/10">Legacy Widget 1</div>
<div class="bg-primary/20">Legacy Widget 2</div>
<!-- charites:ignore-end -->
```

### 8. Configuration Reference (`charites.yaml`)

Configure rule severity or disable it in `charites.yaml`:

```yaml
rules:
  theme.hardcode-opacity-color:
    severity: error # error | warn | info | off
```
