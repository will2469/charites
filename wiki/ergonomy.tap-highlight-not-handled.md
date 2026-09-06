# ergonomy.tap-highlight-not-handled

> **Rule ID:** `ergonomy.tap-highlight-not-handled`
> **Severity:** `INFO`
> **Category:** `ergonomy`
> **Target Standards:** W3C Touch Events Community Group Guidelines, Chromium Android Tap Feedback UX Standards, Google Material Design (Tactile States & Surface Elevation)

---

## 1. Overview & Core Invariant

Flags clickable non-native custom elements lacking tactile tap feedback or tap-highlight management

### Core Invariant:
> **"Non-native clickable elements (<div onClick>, <span role="button">) must declare deliberate active feedback or suppress the default Android Chrome grey tap highlight box."**

---
## 2. Technical Grounding & Engine Realities

On Chromium Android, tapping an element without a native button role causes the browser to flash a rigid semi-transparent grey overlay box.

Without deliberate 'active:' micro-interactions (such as 'active:scale-[0.99]' or 'active:bg-muted') or setting '[-webkit-tap-highlight-color:transparent]', the application exhibits noticeable visual glitches and lacks native tactile responsiveness.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Visual Glitches on Android Chrome** | LOW | Rigid grey highlight rectangles flash abruptly over custom cards, badges, and list rows during touch. |
| **Poor Tactile Feedback** | LOW | Users cannot perceive if a touch tap registered, leading to repeated tapping and accidental duplicate submissions. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Non-native clickable div without active feedback):
```tsx
<div
  role="button"
  tabIndex={0}
  onClick={handleSelectCard}
  className="p-4 bg-card border rounded-2xl"
>
  <span>Pilihan Layanan</span>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Deliberate active feedback with suppressed grey highlight):
```tsx
<div
  role="button"
  tabIndex={0}
  onClick={handleSelectCard}
  className="p-4 bg-card border rounded-2xl active:bg-muted/60 active:scale-[0.99] transition-transform [-webkit-tap-highlight-color:transparent]"
>
  <span>Pilihan Layanan</span>
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ergonomy.tap-highlight-not-handled intentional exception -->
```

```tsx
// charites:ignore ergonomy.tap-highlight-not-handled intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ergonomy.tap-highlight-not-handled:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ergonomy Category Guide](ergonomy).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


