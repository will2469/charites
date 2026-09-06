# inp.heavy-event-handler

> **Rule ID:** `inp.heavy-event-handler`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** W3C Web Performance Working Group (Interaction to Next Paint - INP), Google Chrome Core Web Vitals (Input Delay & Processing Duration), Browser Cooperative Scheduling Guidelines (scheduler.yield)

---

## 1. Overview & Core Invariant

Interactive event handler executes heavy synchronous operations (JSON.parse, Array.sort) without cooperative yields

### Core Invariant:
> **"Interactive event handlers must avoid heavy synchronous computations on the main thread, adopting cooperative task yielding or Web Worker offloading."**

---
## 2. Technical Grounding & Engine Realities

When users tap, click, or type, the browser expects the main thread to quickly acknowledge the interaction and schedule the next paint frame (ideally within 50ms, with INP target <= 200ms).

Executing heavy synchronous operations (such as large JSON parsing, array sorting, or complex synchronous data manipulation) directly inside event handler callbacks blocks the main thread during the crucial input processing phase.

This delays the presentation of visual feedback (e.g. active button states, loading spinners) and directly inflates the INP metric.

Breaking long tasks with 'await scheduler.yield?.()' or offloading computation to a dedicated Web Worker allows the browser to present visual feedback immediately before executing intensive processing.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Processing Phase Thread Saturation** | HIGH | Synchronous algorithms in click/key handlers block the main thread, exceeding the 200ms INP threshold. |
| **Frozen Visual Feedback** | MEDIUM | Buttons appear unresponsive or stuck because UI rendering is starved by long synchronous handler execution. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Synchronous heavy data parsing and sorting directly inside onClick handler):
```tsx
<button onClick={() => {
  const data = JSON.parse(hugePayload);
  const sorted = data.sort((a, b) => b.score - a.score);
  setResults(sorted);
}}>
  Urutkan Data
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Cooperative yielding to acknowledge user input before intensive processing):
```tsx
<button onClick={async () => {
  setLoading(true);
  await (window.scheduler?.yield?.() ?? new Promise(r => setTimeout(r, 0)));
  const data = JSON.parse(hugePayload);
  setResults(data.sort((a, b) => b.score - a.score));
  setLoading(false);
}}>
  Urutkan Data
</button>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore inp.heavy-event-handler intentional exception -->
```

```tsx
// charites:ignore inp.heavy-event-handler intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.heavy-event-handler:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [inp Category Guide](inp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


