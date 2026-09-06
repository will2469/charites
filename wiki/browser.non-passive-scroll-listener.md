# browser.non-passive-scroll-listener

> **Rule ID:** `browser.non-passive-scroll-listener`
> **Severity:** `WARN`
> **Category:** `browser`
> **Target Standards:** W3C DOM Level 4 Events Specification (Passive Event Listeners), Chromium & WebKit Compositor Scrolling Pipeline Guidelines, Google Lighthouse Best Practices (Does not use passive listeners)

---

## 1. Overview & Core Invariant

Enforces { passive: true } option on touch and wheel event listeners to prevent main thread scroll blocking

### Core Invariant:
> **"Event listeners for 'touchstart', 'touchmove', 'wheel', or 'mousewheel' must declare '{ passive: true }' to ensure non-blocking compositor scrolling."**

---
## 2. Technical Grounding & Engine Realities

Browsers execute smooth scrolling on a dedicated compositor thread. Without '{ passive: true }', the compositor must block and wait for JavaScript execution on the main thread to see if 'preventDefault()' is called.

This introduces severe touch response latency and frame rate drops (scroll jank) on mobile devices.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Scroll Jank & Latency** | MEDIUM | Users experience jerky, lagging scrolling and delayed touch gestures on Safari iOS and Android Chrome. |
| **Lighthouse Performance Penalty** | LOW | Fails Lighthouse 'Does not use passive listeners to improve scrolling performance' audit. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### JAVASCRIPT (Adding touchmove listener without passive: true option):
```javascript
window.addEventListener("touchmove", (e) => {
  trackTouchPosition(e);
});
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### JAVASCRIPT (Specifying { passive: true } to unblock the compositor thread):
```javascript
window.addEventListener("touchmove", (e) => {
  trackTouchPosition(e);
}, { passive: true });
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore browser.non-passive-scroll-listener intentional exception -->
```

```tsx
// charites:ignore browser.non-passive-scroll-listener intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.non-passive-scroll-listener:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [browser Category Guide](browser).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


