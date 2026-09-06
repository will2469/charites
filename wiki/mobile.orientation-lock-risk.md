# mobile.orientation-lock-risk

> **Rule ID:** `mobile.orientation-lock-risk`
> **Severity:** `INFO`
> **Category:** `mobile`
> **Target Standards:** W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 1.3.4 (Orientation - Level AA), W3C Screen Orientation API (ScreenOrientation.lock), Google Web Accessibility (Orientation Invariants)

---

## 1. Overview & Core Invariant

Advises against rigid screen orientation locking which restricts accessibility for mounted or assistive mobile setups (WCAG 2.2 SC 1.3.4)

### Core Invariant:
> **"Applications must not rigidly lock display orientation to portrait or landscape unless essential to the core functionality (e.g. bank check capture or piano keyboard)."**

---
## 2. Technical Grounding & Engine Realities

Locking mobile orientation via 'screen.orientation.lock("portrait")' prevents users with assistive needs from accessing content.

Users who have smartphones mounted horizontally on wheelchairs, bed frames, or vehicle dashboards cannot rotate their devices.

Web interfaces should adapt fluidly using responsive CSS (e.g. 'landscape:flex-row') rather than programmatically forbidding device rotation.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Assistive Technology Exclusion** | LOW | Users with fixed horizontal device mounts are unable to view or operate the application naturally. |
| **Unintended Script Errors on Unsupported Browsers** | LOW | Calling orientation lock on Safari iOS or unsupported browsers triggers unhandled promise rejections. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Programmatic orientation lock forces portrait mode):
```tsx
useEffect(() => {
  if (screen.orientation && screen.orientation.lock) {
    screen.orientation.lock("portrait").catch(() => {});
  }
}, []);
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Fluid responsive layout adapting naturally to landscape orientation):
```tsx
<div className="flex flex-col landscape:flex-row gap-4 p-4">
  <aside className="w-full landscape:w-64">Navigasi</aside>
  <main className="flex-1">Konten Utama</main>
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore mobile.orientation-lock-risk intentional exception -->
```

```tsx
// charites:ignore mobile.orientation-lock-risk intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  mobile.orientation-lock-risk:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [mobile Category Guide](mobile).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


