# browser.hover-only-interaction

> **Rule ID:** `browser.hover-only-interaction`
> **Severity:** `ERROR`
> **Category:** `browser`
> **Target Standards:** W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 2.1.1 (Keyboard), WICG / WHATWG Touch Events & Pointer Events Level 3 (Touch vs Hover Ergonomics), Apple Human Interface Guidelines for iOS Touch Interactions

---

## 1. Overview & Core Invariant

Ensures interactive actions and state reveals have keyboard and touch counterparts instead of relying solely on hover

### Core Invariant:
> **"Interactive controls and revealed elements must not rely exclusively on ':hover' or 'group-hover:' without keyboard/touch counterparts ('focus-visible:', 'group-focus-within:')."**

---
## 2. Technical Grounding & Engine Realities

Touchscreen devices (the majority of web traffic on Safari iOS and Chrome Android) have no physical cursor and cannot perform genuine mouse hover.

When critical action buttons (e.g. delete, edit, copy) are hidden by default with 'opacity-0' and only revealed via 'group-hover:opacity-100':
1. Total Mobile Inaccessibility: Touchscreen users cannot see or activate the controls because hovering does not exist on mobile.
2. iOS Sticky Hover Bug: Tapping an element on Safari iOS triggers an inconsistent 'sticky hover' state, requiring multiple confusing taps.
3. Keyboard Navigation Failure: Users navigating with the 'Tab' key cannot discover or focus hidden controls unless focus-within or focus-visible is bound.

Charites enforces that any hover-revealed element provides accessible keyboard and touch parity.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Touch Exclusion** | HIGH | Critical action controls are completely invisible and unreachable on smartphones and tablets. |
| **Keyboard Accessibility Barrier** | MEDIUM | Fails WCAG 2.2 Level A keyboard navigation audits when hidden controls cannot be focused with the Tab key. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Delete button hidden by default and only revealed on group-hover (invisible on touchscreens)):
```tsx
<div className="group flex items-center justify-between p-3 border rounded-xl">
  <span>Berkas_KTP.pdf</span>
  <button
    type="button"
    onClick={handleDelete}
    className="opacity-0 group-hover:opacity-100 text-destructive text-sm"
  >
    Hapus
  </button>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Button accessible via hover, keyboard Tab navigation, and touch focus):
```tsx
<div className="group flex items-center justify-between p-3 border rounded-xl">
  <span>Berkas_KTP.pdf</span>
  <button
    type="button"
    onClick={handleDelete}
    className="opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 focus-visible:opacity-100 text-destructive text-sm"
  >
    Hapus
  </button>
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore browser.hover-only-interaction intentional exception -->
```

```tsx
// charites:ignore browser.hover-only-interaction intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.hover-only-interaction:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [browser Category Guide](browser).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


