# browser.chrome-only-api

> **Rule ID:** `browser.chrome-only-api`
> **Severity:** `WARN`
> **Category:** `browser`
> **Target Standards:** W3C Device Memory Specification (Chromium-only draft), W3C Network Information API (Excluded from Safari and Firefox), W3C Web Platform Status (Mozilla & WebKit Positions on Chromium APIs)

---

## 1. Overview & Core Invariant

Flags reliance on Chromium-exclusive APIs without cross-browser fallbacks for Firefox and Safari

### Core Invariant:
> **"Reliance on Chromium-exclusive APIs ('navigator.deviceMemory', 'navigator.connection', Web Serial/USB, etc.) must include runtime guards and cross-browser fallbacks."**

---
## 2. Technical Grounding & Engine Realities

Chromium exposes several proprietary or non-consensus APIs that Mozilla Firefox and Apple WebKit have formally opposed due to privacy or architecture concerns.

Direct, unguarded access to 'navigator.deviceMemory' or 'navigator.connection' will fail or return 'undefined', throwing TypeErrors when accessing nested properties.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Runtime Crash on Safari and Firefox** | MEDIUM | Accessing 'navigator.connection.effectiveType' throws 'TypeError: Cannot read properties of undefined' on Safari and desktop Firefox. |
| **Feature Lockout for Non-Chromium Users** | MEDIUM | Users on Safari (iOS/macOS) and Firefox cannot complete critical workflows if File System Access or Web Serial is required without a standard fallback. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### JAVASCRIPT (Direct unguarded access to navigator.connection):
```javascript
const effectiveSpeed = navigator.connection.effectiveType;
```
### JAVASCRIPT (Direct access to navigator.deviceMemory):
```javascript
const memoryGiB = navigator.deviceMemory;
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### JAVASCRIPT (Guarded access with default fallback for non-Chromium browsers):
```javascript
const effectiveSpeed = (typeof navigator !== "undefined" && navigator.connection?.effectiveType) || "4g";
```
### JAVASCRIPT (Capability guard using in operator):
```javascript
const memoryGiB = (typeof navigator !== "undefined" && "deviceMemory" in navigator) ? navigator.deviceMemory : 4;
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore browser.chrome-only-api intentional exception -->
```

```tsx
// charites:ignore browser.chrome-only-api intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.chrome-only-api:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [browser Category Guide](browser).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


