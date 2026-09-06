# ux.spacing-inversion

> **Rule ID:** `ux.spacing-inversion`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Gestalt Law of Proximity (Visual Perceptual Hierarchy), Tailwind CSS v3 Space-Between Sibling Selector Specificity Quirks, W3C Design Tokens Community Group (DTCG v2025.10 - Spatial Scale)

---

## 1. Overview & Core Invariant

Warns when child element intra-spacing exceeds parent gap or when space-y conflicts with child mt margin in Tailwind v3

### Core Invariant:
> **"Child element intra-spacing must be strictly tighter than the inter-element gap separating parent sibling groups, and parent 'space-y-*' must not conflict with child 'mt-*' overrides."**

---
## 2. Technical Grounding & Engine Realities

According to the Gestalt Law of Proximity, elements that belong to the same logical group must have smaller internal spacing than the boundary gap between distinct sibling groups.

When a child card or section applies an internal margin or gap that is larger than or equal to the parent container's gap (e.g., parent has 'space-y-4' while child has 'mb-8'), the visual cohesion dissolves, confusing users about which headings, texts, or actions belong together.

Furthermore, in Tailwind CSS v3, 'space-y-*' generates a complex sibling selector '> :not([hidden]) ~ :not([hidden])' with specificity (0, 3, 0), which silently overrides any child 'mt-*' utility (0, 1, 0) without a compiler error. Switching the parent to 'flex flex-col gap-*' restores deterministic CSS cascade and spatial clarity.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cognitive Grouping Disruption** | MEDIUM | Users misattribute subheadings and actions to unrelated neighbouring cards due to broken proximity cues. |
| **Silent CSS Specificity Override** | MEDIUM | Tailwind v3 sibling selectors override child margins without error, leading to unintended layout shifts. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Parent uses space-y-4 while child card specifies mb-8, causing Gestalt proximity inversion and v3 specificity clash):
```tsx
<section className="space-y-4">
  <div className="mb-8">
    <h3 className="text-sm font-semibold">Grup A</h3>
  </div>
  <div className="mb-8">
    <h3 className="text-sm font-semibold">Grup B</h3>
  </div>
</section>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Parent sets wider gap-8 separating groups, while child maintains tighter mb-3 intra-spacing):
```tsx
<section className="flex flex-col gap-8">
  <div className="mb-3">
    <h3 className="text-sm font-semibold">Grup A</h3>
  </div>
  <div className="mb-3">
    <h3 className="text-sm font-semibold">Grup B</h3>
  </div>
</section>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.spacing-inversion intentional exception -->
```

```tsx
// charites:ignore ux.spacing-inversion intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.spacing-inversion:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


