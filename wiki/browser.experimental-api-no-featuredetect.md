# browser.experimental-api-no-featuredetect

> **Rule ID:** `browser.experimental-api-no-featuredetect`
> **Severity:** `ERROR`
> **Category:** `browser`
> **Target Standards:** WICG Web Share API Specification (navigator.share), WICG File System Access API (showOpenFilePicker), W3C CSS View Transitions Module Level 1 (startViewTransition), ECMA-262 Feature Detection & Defensive JavaScript Guidelines

---

## 1. Overview & Core Invariant

Detects invocation of experimental Web APIs without runtime feature detection guards

### Core Invariant:
> **"Experimental or non-universal Web APIs must be guarded with runtime capability checks ('prop' in obj, if (obj.prop), optional chaining, or try/catch) before invocation."**

---
## 2. Technical Grounding & Engine Realities

Modern Web APIs are adopted unevenly across browser vendors. For instance, 'showOpenFilePicker' is exclusive to Chromium and crashes instantly on Firefox and Safari.

'navigator.share' throws an uncaught TypeError on desktop Firefox or non-secure contexts.

Directly calling these APIs without feature guards results in severe runtime exceptions that crash SPAs.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Runtime Exception Crash** | HIGH | Uncaught TypeError: undefined is not a function immediately terminates JavaScript execution in non-supporting browsers. |
| **Broken Core User Action** | HIGH | Primary user actions like document sharing or file importing fail completely without informative fallback UX. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Direct invocation of navigator.share without runtime feature guard (crashes on desktop Firefox)):
```tsx
<button onClick={() => {
  navigator.share({ title: "Surat", url: window.location.href });
}}>
  Bagikan
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Defensive feature detection with fallback to clipboard copy):
```tsx
<button onClick={async () => {
  if (typeof navigator !== "undefined" && navigator.share) {
    await navigator.share({ title: "Surat", url: window.location.href });
  } else {
    await navigator.clipboard?.writeText(window.location.href);
  }
}}>
  Bagikan
</button>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore browser.experimental-api-no-featuredetect intentional exception -->
```

```tsx
// charites:ignore browser.experimental-api-no-featuredetect intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.experimental-api-no-featuredetect:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [browser Category Guide](browser).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


