# performance.react-derived-state-in-effect

> **Rule ID:** `performance.react-derived-state-in-effect`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** React Official Documentation ('You Might Not Need an Effect'), React Reconciliation Lifecycle (Avoiding Cascading Secondary Renders), React Best Practices on Pure Render-Phase Computations

---

## 1. Overview & Core Invariant

Mencegah sinkronisasi derived state dari props atau state yang sudah ada melalui useEffect, yang memicu siklus perenderan sekunder ganda.

### Core Invariant:
> **"Derived values computed synchronously from props or existing state must be calculated directly in the component body during render; updating state inside 'useEffect' triggers redundant secondary render passes."**

---
## 2. Technical Grounding & Engine Realities

Updating state within a `useEffect` callback causes React to first render the component with the stale value, commit it to the DOM, and immediately schedule a second render pass to apply the updated state.

When derived state is calculated synchronously (e.g. concatenating names, filtering a list, or calculating a total), calculating it in an effect needlessly burns main thread time on layout calculations and DOM diffing twice per interaction.

Computing the value directly during the render pass completely eliminates the secondary render pass and keeps state management minimal.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cascading Secondary Render Cycles** | HIGH | Forces React to run duplicate render and diff cycles on every prop change, directly degrading interaction responsiveness (INP). |
| **Visual Stutter / Layout Shift** | MEDIUM | May momentarily display stale computed values before the effect updates, causing brief content flicker or layout shifts. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Sinkronisasi derived state via useEffect memicu render ganda):
```tsx
const [fullName, setFullName] = useState('');
useEffect(() => {
  setFullName(firstName + ' ' + lastName);
}, [firstName, lastName]);
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Dihitung secara sinkron dalam satu kali fase render):
```tsx
const fullName = firstName + ' ' + lastName;
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.react-derived-state-in-effect intentional exception -->
```

```tsx
// charites:ignore performance.react-derived-state-in-effect intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.react-derived-state-in-effect:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [performance Category Guide](performance).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


