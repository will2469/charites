# inp.sync-layout-effect

> **Rule ID:** `inp.sync-layout-effect`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** React useLayoutEffect Pre-Paint Execution Model, W3C Presentation Timing & Frame Pipeline Invariants, Google Chrome Core Web Vitals (INP Presentation Delay Optimization)

---

## 1. Overview & Core Invariant

Synchronous non-geometrical computation in useLayoutEffect blocks browser paint and inflates presentation delay

### Core Invariant:
> **"The 'useLayoutEffect' hook must be reserved strictly for synchronous DOM measurements; data fetching and non-geometrical state updates must reside in 'useEffect'."**

---
## 2. Technical Grounding & Engine Realities

Unlike 'useEffect' which runs asynchronously after the browser paints the screen, 'useLayoutEffect' fires synchronously immediately after React commits DOM mutations, *before* the browser renders pixels to the screen.

Executing non-geometrical operations (such as data fetching, localStorage I/O, or secondary state cascades) within 'useLayoutEffect' delays the browser paint phase directly, locking the main thread and dramatically increasing Presentation Delay.

Developers should restrict 'useLayoutEffect' exclusively to reading layout properties (e.g. 'getBoundingClientRect') to position popovers or tooltips without flicker, moving all other logic to 'useEffect'.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Browser Paint Phase Halting** | HIGH | Frame rendering is synchronously blocked while non-layout logic executes in useLayoutEffect. |
| **Presentation Delay Spikes** | HIGH | Visual acknowledgment of user interactions is delayed by hundreds of milliseconds, breaching the 200ms INP threshold. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Data fetching inside useLayoutEffect blocks the browser paint phase):
```tsx
useLayoutEffect(() => {
  fetchUserData(userId).then(setData);
}, [userId]);
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Data fetching moved to useEffect; browser paints pixels without delay):
```tsx
useEffect(() => {
  fetchUserData(userId).then(setData);
}, [userId]);
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore inp.sync-layout-effect intentional exception -->
```

```tsx
// charites:ignore inp.sync-layout-effect intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.sync-layout-effect:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [inp Category Guide](inp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


