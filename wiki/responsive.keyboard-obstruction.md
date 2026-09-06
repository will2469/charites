# responsive.keyboard-obstruction

> **Rule ID:** `responsive.keyboard-obstruction`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** WCAG 2.2 Guideline 2.1 (Keyboard Accessible), Material Design 3 Mobile Form Guidelines, iOS Human Interface Guidelines (Managing the Virtual Keyboard)

---

## 1. Overview & Core Invariant

Warns against fixed bottom action bars in forms lacking vertical scroll containers, which can be obstructed by the mobile virtual keyboard

### Core Invariant:
> **"Forms containing text inputs and fixed/sticky bottom action bars must provide a scrollable container ('overflow-y-auto') so inputs and actions are never obscured when the virtual keyboard expands."**

---
## 2. Technical Grounding & Engine Realities

When a user taps an input on a smartphone, the virtual software keyboard slides up from the bottom of the screen, consuming 40% to 50% of the visible viewport.

Elements styled with 'fixed bottom-0' or 'sticky bottom-0' remain pinned above the viewport bottom or above the keyboard. If the parent form is not wrapped in a vertical scroll container ('overflow-y-auto'), the active input field gets pushed behind the keyboard or under the fixed button, leaving the user unable to view their input or complete submission.

Charites enforces scrollable viewport resilience for mobile form layouts.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Virtual Keyboard Input Obstruction** | HIGH | Users cannot see text being typed into lower form inputs because fixed bottom bars pin directly over them. |
| **Form Abandonment & Submission Blockers** | MEDIUM | When keyboard expansion pushes inputs offscreen without scroll capabilities, conversion rates drop. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Fixed bottom submit bar in a rigid form lacking scrollable region):
```tsx
<form className="h-screen flex flex-col justify-between">
  <div className="p-4 space-y-4">
    <input type="text" placeholder="Nama Lengkap" />
    <input type="email" placeholder="Alamat Surel" />
    <textarea placeholder="Pesan Anda" />
  </div>
  <div className="fixed bottom-0 inset-x-0 p-4 bg-surface border-t">
    <button type="submit" className="w-full bg-primary text-white py-3 rounded">Kirim</button>
  </div>
</form>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Scrollable body with fixed bottom bar allowing smooth mobile keyboard reflow):
```tsx
<form className="h-screen flex flex-col">
  <div className="flex-1 overflow-y-auto p-4 space-y-4">
    <input type="text" placeholder="Nama Lengkap" />
    <input type="email" placeholder="Alamat Surel" />
    <textarea placeholder="Pesan Anda" />
  </div>
  <div className="p-4 bg-surface border-t">
    <button type="submit" className="w-full bg-primary text-white py-3 rounded">Kirim</button>
  </div>
</form>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore responsive.keyboard-obstruction intentional exception -->
```

```tsx
// charites:ignore responsive.keyboard-obstruction intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.keyboard-obstruction:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [responsive Category Guide](responsive).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


