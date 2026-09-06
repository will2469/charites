# responsive.horizontal-overflow

> **Rule ID:** `responsive.horizontal-overflow`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** W3C CSS Overflow Module Level 3, Google Web Vitals (Preventing Unintended Layout Shifts), iOS WebKit Touch Gesture Chaining Guidelines

---

## 1. Overview & Core Invariant

Warns when unconstrained overflow-x-scroll is declared without fluid width boundary or dynamic auto-scrolling

### Core Invariant:
> **"Horizontal scroll containers on mobile baseline must declare fluid boundaries ('w-full', 'max-w-full') and use dynamic scrolling ('overflow-x-auto') rather than forced permanent scrollbar rails ('overflow-x-scroll')."**

---
## 2. Technical Grounding & Engine Realities

Declaring 'overflow-x-scroll' directly on mobile baseline forces WebKit and Chromium browsers to render a permanent, unyielding scrollbar rail even when content fits within the viewport.

Furthermore, when horizontal scrolling lacks explicit boundary containment ('w-full' or 'max-w-full'), touch drag events can bleed into root document scrolling, causing disorienting horizontal page wobble.

Using 'overflow-x-auto w-full' ensures content only scrolls when overflowing, preserves natural gesture chaining, and prevents horizontal page wobble.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Forced Permanent Scrollbar Rails** | MEDIUM | Rigid horizontal scrollbar tracks clutter compact mobile viewports even when content fits completely. |
| **Broken Vertical Touch Gesture Chaining** | LOW | Users attempting vertical swipe navigation get trapped inside unconstrained horizontal scroll containers. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Unbounded overflow-x-scroll forcing persistent scrollbar rail):
```tsx
<div className="overflow-x-scroll">
  <div className="flex gap-4">
    <div className="p-4 bg-card">Item 1</div>
  </div>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Fluid container using dynamic overflow-x-auto):
```tsx
<div className="w-full overflow-x-auto">
  <div className="flex gap-4 min-w-max">
    <div className="p-4 bg-card">Item 1</div>
  </div>
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore responsive.horizontal-overflow intentional exception -->
```

```tsx
// charites:ignore responsive.horizontal-overflow intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.horizontal-overflow:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [responsive Category Guide](responsive).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


