# lcp.render-blocking-head-script

> **Rule ID:** `lcp.render-blocking-head-script`
> **Severity:** `WARN`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Element Render Delay), HTML Living Standard (The script element & parser-blocking execution pipeline), W3C Web Performance Working Group Critical Path Minimization Guidelines

---

## 1. Overview & Core Invariant

External script in '<head>' without 'defer', 'async', or 'type="module"' synchronously blocks HTML parser and delays LCP candidate paint

### Core Invariant:
> **"External scripts declared in the document '<head>' must specify 'defer', 'async', or 'type="module"' to prevent halting HTML tokenization and per-frame rendering before the LCP candidate is displayed."**

---
## 2. Technical Grounding & Engine Realities

When the browser HTML parser encounters a synchronous `<script src="...">` tag in the `<head>`, it must halt DOM construction, initiate a TCP/TLS connection to the script origin, download the JavaScript payload, and execute it before resuming document rendering.

In Astro, standard `<script>` tags are automatically processed by the bundler into deferred ES modules. However, external scripts tagged with `is:inline` or raw `<script src>` tags in document layouts bypass bundling and execute synchronously.

Adding 'defer' or 'type="module"' allows the browser to download the script in parallel in the background while continuing HTML parsing and per-frame rendering of the hero LCP element.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Synchronous HTML Parser Halting** | HIGH | Halts DOM construction and suppresses layout passes, directly inflating LCP Element Render Delay by the full script network latency. |
| **Head-of-Line Network Contention** | MEDIUM | Competes with critical hero media and external stylesheets for initial HTTP connection bandwidth. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Synchronous external inline script blocking HTML parsing in document head):
```astro
<head>
  <script is:inline src="https://analytics.example.com/tracker.js"></script>
</head>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (External inline script deferred to allow non-blocking HTML parsing):
```astro
<head>
  <script is:inline src="https://analytics.example.com/tracker.js" defer></script>
</head>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.render-blocking-head-script intentional exception -->
```

```tsx
// charites:ignore lcp.render-blocking-head-script intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.render-blocking-head-script:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


