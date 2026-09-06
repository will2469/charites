# pwa.insecure-context-resource

> **Rule ID:** `pwa.insecure-context-resource`
> **Severity:** `ERROR`
> **Category:** `pwa`
> **Target Standards:** W3C Secure Contexts Specification, W3C Mixed Content Level 2 Specification, RFC 7258 Pervasive Monitoring Is an Attack

---

## 1. Overview & Core Invariant

Errors when a resource element loads assets over an insecure HTTP protocol (http://) in violation of W3C Secure Contexts

### Core Invariant:
> **"Resource elements (<script>, <link>, <img>, <iframe>, <video>, <audio>) must not load external assets over insecure 'http://' (except localhost loopback)."**

---
## 2. Technical Grounding & Engine Realities

Progressive Web Apps strictly require a Secure Context (HTTPS) to enable service workers, cache storage, and device hardware APIs.

Loading external assets (scripts, stylesheets, images, media, iframes) via an unencrypted 'http://' connection triggers Mixed Content blocking. Active mixed content (scripts, stylesheets) is blocked immediately by modern mobile browsers, breaking application functionality. Passive mixed content (images, audio) generates security warnings and can be intercepted or tampered with on public Wi-Fi networks.

All asset references must use HTTPS ('https://'), protocol-relative URLs ('//'), or local origin paths. Localhost addresses ('http://localhost' and 'http://127.0.0.1') are excepted for development purposes.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Active Mixed Content Blocking** | HIGH | Browsers completely block insecure scripts and stylesheets, breaking UI styling and interactive application logic. |
| **Man-in-the-Middle Asset Tampering** | HIGH | Unencrypted HTTP traffic can be intercepted, inspected, or modified by malicious actors on untrusted networks. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Insecure HTTP external script and stylesheet links):
```tsx
<div>
  <script src="http://cdn.example.org/tracker.js" />
  <link rel="stylesheet" href="http://assets.desa.id/styles.css" />
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Secure HTTPS asset loading conforming to Secure Contexts):
```tsx
<div>
  <script src="https://cdn.example.org/tracker.js" />
  <link rel="stylesheet" href="https://assets.desa.id/styles.css" />
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore pwa.insecure-context-resource intentional exception -->
```

```tsx
// charites:ignore pwa.insecure-context-resource intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  pwa.insecure-context-resource:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [pwa Category Guide](pwa).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


