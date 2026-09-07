# ergonomy.number-input-wheel-hazard

> **Rule ID:** `ergonomy.number-input-wheel-hazard`
> **Severity:** `WARN`
> **Category:** `ergonomy`
> **Target Standards:** Nielsen Norman Group (Input Steppers and Number Inputs), GOV.UK Design System (Numbers and Form Inputs), W3C HTML 5.2 Section 4.10.5.1.12 (Number State)

---

## 1. Overview & Core Invariant

Flags number inputs without wheel prevention to prevent unintended value mutations during scrolling (Scroll Hijacking)

### Core Invariant:
> **"Number inputs (<input type="number">) must declare an explicit wheel event handler (e.g. 'onWheel={(e) => e.currentTarget.blur()}') to claim ownership and prevent accidental value mutations during page scrolling."**

---
## 2. Technical Grounding & Engine Realities

When a user focuses on an input with type="number" and scrolls the document using a mouse wheel or trackpad, desktop browsers intercept the scroll delta to increment or decrement the numeric value rather than moving the viewport.

This phenomenon (Scroll Hijacking) silently mutates quantitative user data (such as item quantities or monetary amounts) without the user's conscious intent. Disabled and readOnly inputs are excluded from this rule because they are not intended for direct user mutation.

Note: Static analysis checks whether the developer explicitly declared a wheel handler (such as 'onWheel' or 'onWheelCapture') as a static protection claim. It does not perform runtime interprocedural flow analysis on the callback body.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Accidental Value Mutation (Scroll Hijacking)** | HIGH | Users scrolling past focused inputs unintentionally increment or decrement orders, quantities, or amounts. |
| **Silent Submission Corruption** | HIGH | Mutated values submit without visual warning because the user assumes their gesture only shifted the viewport. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Unprotected number input vulnerable to mouse wheel mutation during page scroll):
```tsx
<input
  type="number"
  name="quantity"
  min="1"
/>
```
### TSX (Unprotected component-based number input):
```tsx
<Input
  type="number"
  name="item_count"
  min="0"
/>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Blur on wheel prevents mouse wheel value increment while preserving page scrolling):
```tsx
<input
  type="number"
  name="quantity"
  min="1"
  onWheel={(e) => e.currentTarget.blur()}
/>
```
### TSX (Component-based input with designated wheel handler callback):
```tsx
<Input
  type="number"
  name="item_count"
  min="0"
  onWheel={handleWheel}
/>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ergonomy.number-input-wheel-hazard intentional exception -->
```

```tsx
// charites:ignore ergonomy.number-input-wheel-hazard intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ergonomy.number-input-wheel-hazard:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ergonomy Category Guide](ergonomy).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


