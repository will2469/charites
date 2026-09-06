# pwa.start-url-inconsistency

> **Rule ID:** `pwa.start-url-inconsistency`
> **Severity:** `ERROR`
> **Category:** `pwa`
> **Target Standards:** W3C Web App Manifest Section 5.2 (The start_url member), W3C Secure Contexts (Mixed Content Mitigation), RFC 3986 Uniform Resource Identifier (URI): Generic Syntax

---

## 1. Overview & Core Invariant

Errors when a Web App Manifest start_url uses an insecure protocol (http://), script scheme (javascript:), or path traversal (../)

### Core Invariant:
> **"Web App Manifest 'start_url' must not use insecure HTTP protocols, script URI schemes (javascript:), or directory traversal ('../')."**

---
## 2. Technical Grounding & Engine Realities

The 'start_url' member defines the preferred URL that should be loaded when the user launches the web application from the mobile launcher.

According to W3C PWA specifications, PWAs must operate strictly within secure contexts (HTTPS). Setting 'start_url' to an insecure HTTP URL ('http://') causes mobile browsers to block launch execution. Setting 'start_url' to a path traversal sequence ('../') escapes the intended navigation scope and causes unpredictable routing failures.

Using a clean relative path (e.g. '/' or '/app') under the secure origin ensures consistent and secure PWA startup.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Insecure Context Launch Failure** | HIGH | Mobile browsers block launching PWAs whose start_url does not satisfy Secure Context requirements. |
| **Path Traversal Outside Navigation Scope** | HIGH | Using '../' breaks origin scope confinement, leading to broken navigation and failed manifest resolution. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Insecure HTTP protocol in start_url):
```tsx
<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    start_url: "http://desa.id/app",
    display: "standalone",
    icons: [{ src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }]
  })}
</script>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Valid relative path for start_url):
```tsx
<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    start_url: "/",
    display: "standalone",
    icons: [{ src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }]
  })}
</script>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore pwa.start-url-inconsistency intentional exception -->
```

```tsx
// charites:ignore pwa.start-url-inconsistency intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  pwa.start-url-inconsistency:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [pwa Category Guide](pwa).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


