# inp.layout-thrashing

> **Rule ID:** `inp.layout-thrashing`
> **Severity:** `ERROR`
> **Category:** `inp`
> **Target Standards:** W3C Web Performance Working Group (Interaction to Next Paint - INP), Google Chrome Rendering Engine Pipeline (Forced Synchronous Layout), Browser Main-Thread Event Loop Scheduling

---

## 1. Overview & Core Invariant

Sequential DOM style mutation followed by layout geometry reading triggers forced synchronous reflow

### Core Invariant:
> **"Imperative JavaScript execution must separate layout queries from style mutations, avoiding forced synchronous reflow passes (read-then-write batching)."**

---
## 2. Technical Grounding & Engine Realities

When JavaScript mutates DOM styles or class names (e.g. 'el.style.width = ...') and subsequently reads a layout geometry property (e.g. 'el.offsetHeight' or 'getBoundingClientRect()') within the same synchronous execution block, the browser is forced to flush pending style changes and perform an immediate, blocking layout recalculation.

This phenomenon, known as 'Layout Thrashing' or 'Forced Synchronous Reflow', locks the browser main thread, preventing user interaction processing and drastically inflating Interaction to Next Paint (INP) latency.

Batching all layout reads before performing style writes, or deferring updates via 'requestAnimationFrame', prevents synchronous recalculations and keeps the main thread responsive.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Forced Synchronous Reflow Stalling** | HIGH | Synchronous layout computation blocks the main thread during interaction handling, causing dropped frames and severe INP degradation. |
| **Cascading Rendering Bottleneck** | HIGH | Interleaved write-read loops exponentially degrade interaction responsiveness on complex DOM trees. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Style mutation immediately followed by geometry reading (forced reflow)):
```tsx
function resizeBox(el: HTMLElement) {
  el.style.width = '200px';
  const height = el.offsetHeight; // Memaksa kalkulasi layout sinkron!
  el.style.height = (height * 2) + 'px';
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Read-then-write batching to prevent forced layout calculation):
```tsx
function resizeBox(el: HTMLElement) {
  const currentHeight = el.offsetHeight; // Baca di awal
  el.style.width = '200px';              // Tulis serentak
  el.style.height = (currentHeight * 2) + 'px';
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore inp.layout-thrashing intentional exception -->
```

```tsx
// charites:ignore inp.layout-thrashing intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.layout-thrashing:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [inp Category Guide](inp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


