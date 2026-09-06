# ergonomy.missing-inputmode-keyboard

> **Rule ID:** `ergonomy.missing-inputmode-keyboard`
> **Severity:** `INFO`
> **Category:** `ergonomy`
> **Target Standards:** HTML Living Standard Section 4.10.5.3 (The inputmode attribute), Tesler's Law (Conservation of Complexity in Virtual Keyboards), Apple iOS & Android Mobile Virtual Keyboard Guidelines

---

## 1. Overview & Core Invariant

Enforces contextual virtual keyboard inputmode and type attributes on mobile form inputs (Tesler's Law)

### Core Invariant:
> **"Form text inputs collecting numeric, phone, or email values must declare contextual 'inputmode' or specialized 'type' attributes to directly open the optimized mobile virtual keypad."**

---
## 2. Technical Grounding & Engine Realities

On mobile devices, focusing an input without specialized type or inputmode opens the full standard QWERTY keyboard.

For numeric, telephone, or OTP fields, this forces the user to repeatedly toggle keyboard layers to find digits. According to Tesler's Law, complexity must be absorbed by software rather than offloaded to user manual effort. Declaring 'inputmode="numeric"' or 'type="tel"' instantly summons large, thumb-friendly numeric keypads.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Cognitive Friction** | LOW | Users are forced to manually switch keyboard layers on small touchscreens to enter digits or phone numbers. |
| **High Form Abandonment & Typing Errors** | LOW | Entering OTP or financial amounts on dense QWERTY keys leads to frequent miss-taps and delayed checkout flows. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Phone input missing type or inputMode):
```tsx
<input
  name="nomor_hp"
  placeholder="08123456789"
  className="h-11 px-3.5 py-2.5 border rounded-xl"
/>
```
### ASTRO (OTP field defaulting to QWERTY keyboard):
```astro
<input
  id="otp_code"
  placeholder="123456"
  class="h-11 px-3.5 py-2.5 border rounded-xl"
/>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Explicit telephone keypad and autocomplete):
```tsx
<input
  name="nomor_hp"
  type="tel"
  inputMode="tel"
  autoComplete="tel"
  placeholder="08123456789"
  className="h-11 px-3.5 py-2.5 border rounded-xl"
/>
```
### ASTRO (Numeric keypad for OTP verification):
```astro
<input
  id="otp_code"
  type="text"
  inputmode="numeric"
  pattern="[0-9]*"
  placeholder="123456"
  class="h-11 px-3.5 py-2.5 border rounded-xl"
/>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ergonomy.missing-inputmode-keyboard intentional exception -->
```

```tsx
// charites:ignore ergonomy.missing-inputmode-keyboard intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ergonomy.missing-inputmode-keyboard:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ergonomy Category Guide](ergonomy).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


