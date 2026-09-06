# pwa.service-worker-registration

> **Rule ID:** `pwa.service-worker-registration`
> **Severity:** `WARN`
> **Category:** `pwa`
> **Target Standards:** W3C Service Workers (Registration Lifecycle & Error Handling), MDN Progressive Web App Guides (Registering a Service Worker Safely), Google Web Fundamentals (Service Worker Reliability)

---

## 1. Overview & Core Invariant

Warns when Service Worker registration lacks feature detection ('serviceWorker' in navigator) or error handling (.catch)

### Core Invariant:
> **"Calls to navigator.serviceWorker.register must be guarded by feature detection ('serviceWorker' in navigator) and handled with error callbacks (.catch or try/catch)."**

---
## 2. Technical Grounding & Engine Realities

Calling navigator.serviceWorker.register() without feature detection triggers fatal runtime TypeErrors in legacy browsers, restricted WebViews, or non-secure HTTP contexts.

Furthermore, failing to handle registration failure (.catch or try/catch) causes unhandled promise rejections that can disrupt analytics scripts and break client-side bootstrapping.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Runtime Script Crash on Older Browsers** | MEDIUM | Browsers or WebViews lacking Service Worker support crash with Uncaught TypeError: Cannot read properties of undefined. |
| **Silent Unhandled Promise Rejections** | LOW | Registration rejections (e.g. 404 or SSL errors) pollute telemetry logs and fail to notify diagnostic listeners. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Unsafe registration without feature detection and error handling):
```tsx
<script>
  navigator.serviceWorker.register('/sw.js');
</script>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Safe registration with feature detection and error handler):
```tsx
<script>
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/sw.js')
      .then((reg) => console.log('SW registered:', reg.scope))
      .catch((err) => console.error('SW registration failed:', err));
  }
</script>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore pwa.service-worker-registration intentional exception -->
```

```tsx
// charites:ignore pwa.service-worker-registration intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  pwa.service-worker-registration:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [pwa Category Guide](pwa).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


