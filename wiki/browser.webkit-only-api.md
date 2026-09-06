# browser.webkit-only-api

> **Rule ID:** `browser.webkit-only-api`
> **Severity:** `WARN`
> **Category:** `browser`
> **Target Standards:** W3C Fullscreen API Specification, W3C Web Audio API Specification, W3C Web Speech API Specification

---

## 1. Overview & Core Invariant

Flags direct invocation of WebKit-prefixed legacy APIs without standard W3C equivalents or graceful fallbacks

### Core Invariant:
> **"Direct invocation of WebKit-prefixed methods ('webkitRequestFullscreen', 'webkitAudioContext', etc.) must provide standard W3C fallbacks for non-WebKit browsers."**

---
## 2. Technical Grounding & Engine Realities

WebKit vendor-prefixed APIs were transitional features during early HTML5 standardization.

Calling them directly without checking W3C standard methods causes instant runtime crashes in Mozilla Firefox desktop and modern Chromium environments.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Crash in Firefox and Non-WebKit Browsers** | MEDIUM | Methods like 'element.webkitRequestFullscreen()' throw 'TypeError: not a function' in Firefox because Gecko does not implement WebKit prefixes. |
| **Missed W3C Standard Improvements** | LOW | Legacy WebKit methods lack updated return types (such as Promises) supported by modern standard W3C APIs. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### JAVASCRIPT (Direct invocation of WebKit fullscreen without standard check):
```javascript
function enterFullscreen(element) {
  element.webkitRequestFullscreen();
}
```
### JAVASCRIPT (Direct instantiation of webkitAudioContext):
```javascript
const audioCtx = new webkitAudioContext();
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### JAVASCRIPT (Prioritizing standard W3C method with WebKit fallback):
```javascript
function enterFullscreen(element) {
  if (element.requestFullscreen) {
    element.requestFullscreen();
  } else if (element.webkitRequestFullscreen) {
    element.webkitRequestFullscreen();
  }
}
```
### JAVASCRIPT (Standard AudioContext fallback chain):
```javascript
const AudioContextClass = window.AudioContext || window.webkitAudioContext;
const audioCtx = new AudioContextClass();
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore browser.webkit-only-api intentional exception -->
```

```tsx
// charites:ignore browser.webkit-only-api intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.webkit-only-api:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [browser Category Guide](browser).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


