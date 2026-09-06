# inp.render-blocking-script

> **Rule ID:** `inp.render-blocking-script`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** HTML Living Standard (The script element & execution pipeline), W3C Web Performance & Navigation Timing Specification, Google Chrome Core Web Vitals (Eliminating Render-Blocking Resources)

---

## 1. Overview & Core Invariant

External script element without defer, async, or type="module" synchronously blocks rendering and input responsiveness

### Core Invariant:
> **"External script elements must declare 'defer', 'async', or 'type="module"' to avoid synchronously blocking HTML parsing and main-thread readiness."**

---
## 2. Technical Grounding & Engine Realities

When the browser encounters a synchronous `<script src="...">` tag, it must pause HTML parsing, establish a network connection, download the script, and execute it before resuming document rendering.

In Astro, standard `<script>` tags are automatically bundled into deferred ES modules. However, scripts marked with `is:inline` or raw external scripts in HTML document heads bypass bundling and execute synchronously.

Adding 'defer' or 'type="module"' ensures the script is downloaded in the background and executed without halting the parser, keeping the browser immediately receptive to early user taps and clicks.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Synchronous Parser Halting** | HIGH | HTML parsing and initial rendering are paused until external scripts download and execute. |
| **Delayed Main-Thread Input Availability** | MEDIUM | The browser input event loop is delayed, resulting in unacknowledged early user taps. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Synchronous external inline script blocking HTML parser):
```astro
<script is:inline src="https://analytics.example.com/heavy-bundle.js"></script>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (External inline script deferred to prevent parser blocking):
```astro
<script is:inline src="https://analytics.example.com/heavy-bundle.js" defer></script>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore inp.render-blocking-script intentional exception -->
```

```tsx
// charites:ignore inp.render-blocking-script intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.render-blocking-script:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [inp Category Guide](inp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


