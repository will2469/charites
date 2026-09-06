# pwa.icon-maskable-missing

> **Rule ID:** `pwa.icon-maskable-missing`
> **Severity:** `WARN`
> **Category:** `pwa`
> **Target Standards:** W3C Web App Manifest Specification (Adaptive Icon Masking), Google Android Maskable Icons Specification, Android Oreo+ Adaptive Launcher Icon Architecture

---

## 1. Overview & Core Invariant

Warns when a Web App Manifest defines icons but none has purpose: 'maskable' for Android adaptive launcher icons

### Core Invariant:
> **"When a Web App Manifest defines icons, at least one icon must declare 'purpose: "maskable"' to prevent Android launcher letterboxing."**

---
## 2. Technical Grounding & Engine Realities

Starting in Android 8.0 Oreo, native device launchers crop application icons according to user-selected device masks (circles, squircles, rounded rectangles).

When a PWA provides only standard icons (purpose: 'any' or omitted purpose), modern Android launchers place the icon inside a small white square box (letterboxing) to fit the mask. This disrupts the visual consistency of native mobile app trays.

Providing at least one icon with 'purpose: "maskable"' (with an appropriate safe zone margin) ensures the launcher can scale and mask the icon seamlessly to fill the full shape.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Letterboxed Android Launcher Icons** | MEDIUM | PWA icon appears inside an awkward white square box on Android device home screens. |
| **Degraded Native Visual Immersion** | LOW | Breaks aesthetic parity with native Android apps installed from Google Play. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Manifest defines icons without any maskable purpose):
```tsx
<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    start_url: "/",
    display: "standalone",
    icons: [
      { src: "/icon-512.png", sizes: "512x512", type: "image/png" }
    ]
  })}
</script>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Manifest includes an adaptive icon with purpose: maskable):
```tsx
<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    start_url: "/",
    display: "standalone",
    icons: [
      { src: "/icon-512.png", sizes: "512x512", type: "image/png" },
      { src: "/icon-512-maskable.png", sizes: "512x512", type: "image/png", purpose: "maskable" }
    ]
  })}
</script>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore pwa.icon-maskable-missing intentional exception -->
```

```tsx
// charites:ignore pwa.icon-maskable-missing intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  pwa.icon-maskable-missing:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [pwa Category Guide](pwa).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


