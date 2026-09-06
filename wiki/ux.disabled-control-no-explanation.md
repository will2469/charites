# ux.disabled-control-no-explanation

> **Rule ID:** `ux.disabled-control-no-explanation`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Nielsen Heuristic #1: Visibility of System Status, Feedforward Principle & Gulf of Evaluation (Don Norman), WCAG 2.2 Success Criterion 3.3.2 (Labels or Instructions)

---

## 1. Overview & Core Invariant

Enforces feedforward explanation for disabled interactive controls to prevent user dead ends

### Core Invariant:
> **"Disabled interactive controls ('disabled', 'aria-disabled="true"') must provide a feedforward explanation via 'aria-describedby', tooltip, or visible helper text."**

---
## 2. Technical Grounding & Engine Realities

When users encounter a disabled button or locked form control without any visible reason, they reach a cognitive dead end: they cannot proceed and have no information on how to unlock the action.

Providing contextual feedforward (e.g. 'Minimum belanja Rp 50.000 untuk checkout' or an explanatory tooltip) clarifies system constraints and transforms frustration into actionable user guidance.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cognitive Dead Ends & Workflow Abandonment** | MEDIUM | Users are unable to diagnose why an action is blocked and abandon conversion or submission workflows. |
| **Accessibility Barrier (Screen Reader Confusion)** | MEDIUM | Assisted technology users hear 'button disabled' without auditory explanation of required missing inputs. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Disabled checkout button without explanation of minimum order constraint):
```tsx
<div className="mt-4">
  <button disabled={cartTotal < 50000} className="bg-primary text-white px-4 py-2 rounded">
    Checkout
  </button>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Disabled button accompanied by aria-describedby linking to explanatory hint text):
```tsx
<div className="mt-4">
  <button
    disabled={cartTotal < 50000}
    aria-describedby="min-order-hint"
    className="bg-primary text-white px-4 py-2 rounded disabled:opacity-50"
  >
    Checkout
  </button>
  {cartTotal < 50000 && (
    <p id="min-order-hint" className="text-xs text-muted-foreground mt-1">
      Minimum belanja Rp 50.000 untuk melanjutkan checkout.
    </p>
  )}
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.disabled-control-no-explanation intentional exception -->
```

```tsx
// charites:ignore ux.disabled-control-no-explanation intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.disabled-control-no-explanation:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


