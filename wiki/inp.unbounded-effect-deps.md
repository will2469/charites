# inp.unbounded-effect-deps

> **Rule ID:** `inp.unbounded-effect-deps`
> **Severity:** `ERROR`
> **Category:** `inp`
> **Target Standards:** React Hooks Specification & Dependency Determinism, W3C Cooperative Scheduling & Frame Budget Invariants, Google Chrome Core Web Vitals (Input Presentation Delay)

---

## 1. Overview & Core Invariant

Lifecycle hook useEffect/useLayoutEffect is missing a dependency array, triggering unbounded re-executions on every render

### Core Invariant:
> **"React lifecycle hooks (useEffect, useLayoutEffect) must explicitly declare a dependency array as their second argument to prevent uncontrolled execution on every render cycle."**

---
## 2. Technical Grounding & Engine Realities

When 'useEffect' or 'useLayoutEffect' is invoked with only a callback and no second argument, React executes the effect after *every single render*.

Any state update, parent re-render, or user keystroke causes the entire effect callback to run again. If the effect queries DOM elements, reads layout properties, or synchronizes subscriptions, the main thread is constantly saturated by unnecessary computations.

Providing an explicit dependency array ('[]' for mount-only, or '[deps...]') restricts execution strictly to when dependencies change, protecting interaction frame rate.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Every-Render Effect Re-execution** | CRITICAL | Effects fire repeatedly on every keystroke or state change, causing severe CPU spikes and input lag. |
| **Infinite Render Loops** | HIGH | If an unbounded effect updates state, it causes an immediate infinite re-render loop that locks the browser tab. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (useEffect without a dependency array executes on every render):
```tsx
useEffect(() => {
  recomputeHeavyLayout();
});
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Explicit empty dependency array ensures execution only on mount):
```tsx
useEffect(() => {
  recomputeHeavyLayout();
}, []);
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore inp.unbounded-effect-deps intentional exception -->
```

```tsx
// charites:ignore inp.unbounded-effect-deps intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.unbounded-effect-deps:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [inp Category Guide](inp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


