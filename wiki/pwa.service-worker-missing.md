# pwa.service-worker-missing

> **Rule ID:** `pwa.service-worker-missing`
> **Severity:** `WARN`
> **Category:** `pwa`
> **Target Standards:** W3C Service Workers (Document Integration), W3C Web App Manifest (Installation Requirements), Google Chrome PWA Criteria (Offline Capability)

---

## 1. Overview & Core Invariant

Warns when an HTML document head links to a Web App Manifest but lacks a Service Worker registration in the document

### Core Invariant:
> **"When an HTML document head links to a Web App Manifest, the document must register a Service Worker via navigator.serviceWorker.register or an external worker script."**

---
## 2. Technical Grounding & Engine Realities

A Progressive Web App requires a registered Service Worker to cache shell assets, intercept network requests during outages, and satisfy full mobile installability audits.

Linking a manifest file without registering a Service Worker leaves the application behaving like a conventional static website, incapable of offline execution or background updates.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Lack of Offline Capability** | MEDIUM | Users cannot open or use the application when device connectivity is intermittent or offline. |
| **Failed Installability Criteria** | LOW | Modern mobile browsers will not trigger automated installation banners without an active Service Worker. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Document head links to manifest but no Service Worker is registered):
```tsx
<head>
  <title>Layanan Desa</title>
  <link rel="manifest" href="/manifest.webmanifest" />
</head>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Document head links to manifest and registers a Service Worker):
```tsx
<head>
  <title>Layanan Desa</title>
  <link rel="manifest" href="/manifest.webmanifest" />
  <script>
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.register('/sw.js').catch(console.error);
    }
  </script>
</head>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore pwa.service-worker-missing intentional exception -->
```

```tsx
// charites:ignore pwa.service-worker-missing intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  pwa.service-worker-missing:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [pwa Category Guide](pwa).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


