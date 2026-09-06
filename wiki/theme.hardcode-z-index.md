# theme.hardcode-z-index

> **Rule ID:** `theme.hardcode-z-index`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** CSS Stacking Context Specification, Design System Elevation Hierarchy, Modal & Overlay Governance Standards

---

## 1. Overview & Core Invariant

Detects hardcoded arbitrary z-index scalars that trigger stacking context wars

### Core Invariant:
> **"Element stacking context elevation must be declared using semantic elevation tokens or CSS variables, never arbitrary numerical z-index scalars."**

---
## 2. Technical Grounding & Engine Realities

Using arbitrary z-index values (e.g. z-[9999] or [z-index:1000]) triggers destructive 'z-index wars':

1. Stacking Context Escalation: When engineers pick arbitrary large numbers (999, 9999, 99999) to force elements to the top, other elements inevitably get occluded.
2. Overlay Clashes: Modals, tooltips, dropdown menus, toast notifications, and sticky navigation headers collide unpredictably.
3. Unmaintainable Layering: Without a centralized hierarchy, debugging stacking context bugs requires inspecting the entire DOM tree.

Charites enforces utilizing structured elevation tokens (e.g. z-dropdown, z-modal, z-toast) or CSS custom properties (e.g. z-[var(--z-modal)]).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Z-Index Escalation Wars** | HIGH | Engineers continually increase z-index numbers, eventually breaking native select popovers and dialogs. |
| **Overlay Occlusion** | HIGH | Tooltips and toasts become permanently trapped behind sticky navigations or dropdown menus. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Arbitrary runaway z-index in fixed modal):
```tsx
<div className="fixed top-0 z-[9999]">Escalated Modal</div>
```
### ASTRO (Arbitrary property z-index):
```astro
<nav class="sticky top-0 [z-index:1000]">Sticky Header</nav>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Semantic elevation token or standard scale):
```tsx
<div className="fixed top-0 z-50">Controlled Modal</div>
```
### ASTRO (Token variable elevation):
```astro
<nav class="sticky top-0 z-[var(--z-sticky)]">Sticky Header</nav>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.hardcode-z-index intentional exception -->
```

```tsx
// charites:ignore theme.hardcode-z-index intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.hardcode-z-index:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


