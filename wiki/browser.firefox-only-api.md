# browser.firefox-only-api

> **Rule ID:** `browser.firefox-only-api`
> **Severity:** `WARN`
> **Category:** `browser`
> **Target Standards:** W3C Fullscreen API Specification, W3C DOM Standards (Gecko Extension Deprecations), MDN Web Docs (Gecko-Specific DOM Interfaces)

---

## 1. Overview & Core Invariant

Flags usage of legacy Gecko/Firefox-exclusive DOM extensions and APIs without standard W3C equivalents

### Core Invariant:
> **"Gecko-prefixed DOM methods and proprietary APIs ('mozRequestFullScreen', 'InstallTrigger', etc.) must provide standard W3C equivalents for Blink and WebKit."**

---
## 2. Technical Grounding & Engine Realities

Mozilla Firefox historically exposed vendor-prefixed APIs such as 'mozRequestFullScreen' and browser-specific globals like 'InstallTrigger'.

Calling these directly without standard W3C methods causes instant crashes or undefined behavior in Blink (Chrome/Edge) and WebKit (Safari).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Crash in Chrome and Safari** | MEDIUM | Invoking 'element.mozRequestFullScreen()' throws 'TypeError: element.mozRequestFullScreen is not a function' in all non-Gecko browsers. |
| **Obsolete Browser Sniffing** | LOW | Relying on 'InstallTrigger' to detect Firefox is deprecated and breaks as Firefox modernizes its engine. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### JAVASCRIPT (Direct invocation of mozRequestFullScreen without standard check):
```javascript
function enterFullscreen(element) {
  element.mozRequestFullScreen();
}
```
### JAVASCRIPT (Direct access to Gecko-specific inner screen property):
```javascript
const screenX = window.mozInnerScreenX;
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### JAVASCRIPT (Prioritizing standard W3C fullscreen method):
```javascript
function enterFullscreen(element) {
  if (element.requestFullscreen) {
    element.requestFullscreen();
  } else if (element.mozRequestFullScreen) {
    element.mozRequestFullScreen();
  }
}
```
### JAVASCRIPT (Standard fullscreenElement fallback chain):
```javascript
const fsElement = document.fullscreenElement || document.mozFullScreenElement;
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore browser.firefox-only-api intentional exception -->
```

```tsx
// charites:ignore browser.firefox-only-api intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.firefox-only-api:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [browser Category Guide](browser).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


