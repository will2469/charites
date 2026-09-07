# ux.number-input-missing-bounds

> **Rule ID:** `ux.number-input-missing-bounds`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Charites Design System & Form Ergonomy Policy, WCAG 2.2 Success Criterion 3.3.4 (Error Prevention - Contextual Usability Guidance), ISO 9241-110:2020 (Interaction Ergonomics - Error Tolerance)

---

## 1. Overview & Core Invariant

Flags mathematical number inputs lacking explicit domain lower bounds (min) to prevent accidental out-of-range submissions

### Core Invariant:
> **"Number inputs representing domain quantities should declare an explicit domain lower bound ('min') to prevent erroneous or accidental out-of-range submissions."**

---
## 2. Technical Grounding & Engine Realities

While the HTML5 specification permits '<input type="number">' without a 'min' attribute, Charites UX design system policy requires developers to declare an explicit domain lower bound for quantitative fields.

Without a declared bound, users can inadvertently scroll or step into nonsensical negative numbers (e.g. quantity: -5, tickets: -1), or submit values below business thresholds.

Cross-Rule Yield Hierarchy: When an input field is classified as a discrete identity code (such as postal code, account number, or NIK), this rule intentionally yields and suppresses itself. Recommending 'min="0"' on an identity field is erroneous pseudo-remediation; the field must instead be converted to text via 'ux.number-input-identity-misuse'.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Accidental Negative Value Submission** | MEDIUM | Users submit nonsensical negative quantities or counts when stepper controls decrement below zero. |
| **Unbounded Stepper Traversal** | LOW | Spinbox controls allow decrementing past business logic boundaries without visual constraint. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Quantity input lacking explicit lower bound can be decremented into negative numbers):
```tsx
<input
  type="number"
  name="quantity"
  onWheel={(e) => e.currentTarget.blur()}
/>
```
### TSX (Discount percentage input without lower bound):
```tsx
<Input
  type="number"
  name="discount_percent"
  onWheel={handleWheel}
/>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Explicit lower bound min='0' prevents stepping into negative numbers):
```tsx
<input
  type="number"
  name="quantity"
  min="0"
  onWheel={(e) => e.currentTarget.blur()}
/>
```
### TSX (Domain temperature input with valid negative bound):
```tsx
<input
  type="number"
  name="temperature"
  min="-50"
  onWheel={(e) => e.currentTarget.blur()}
/>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.number-input-missing-bounds intentional exception -->
```

```tsx
// charites:ignore ux.number-input-missing-bounds intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.number-input-missing-bounds:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


