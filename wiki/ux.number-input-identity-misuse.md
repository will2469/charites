# ux.number-input-identity-misuse

> **Rule ID:** `ux.number-input-identity-misuse`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** GOV.UK Design System (Numbers vs Numeric Text), Nielsen Norman Group (Number Inputs vs. Numeric Text), W3C HTML 5.2 Section 4.10.5.1.12 (Number State vs Text State)

---

## 1. Overview & Core Invariant

Flags type="number" on discrete identity codes and serial identifiers where mathematical operations are meaningless

### Core Invariant:
> **"Inputs capturing discrete identity, serial, or cryptographic tokens (e.g. postal codes, account numbers, phone numbers, security PINs) must use 'type="text"' with 'inputMode="numeric"' rather than 'type="number"'."**

---
## 2. Technical Grounding & Engine Realities

Browsers treat '<input type="number">' as IEEE-754 floating-point mathematical numbers, which causes severe UX defects and data corruption when applied to identity codes:

1. Leading zeros are silently stripped by the browser parser (e.g. postal code '01234' becomes '1234').
2. Numbers exceeding 16 digits (such as bank accounts or payment cards) lose precision or collapse into exponential scientific notation ('1e+16').
3. Unhelpful stepper spinbuttons appear in the UI, prompting users to 'increment' or 'decrement' an identity code where arithmetic is nonsensical.
4. Screen readers announce spinbox increments that confuse assistive technology users.

Remediation requires switching to '<input type="text" inputMode="numeric" pattern="[0-9]*" />' (or 'type="tel"' for phone numbers), which provides mobile numeric keypads without stripping zeros or enabling stepper arrows.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Loss of Leading Zeros** | CRITICAL | Postal codes, account numbers, and identity codes starting with 0 are corrupted during form serialization. |
| **IEEE-754 Precision Truncation** | CRITICAL | Long card numbers or identifiers over 15 digits lose trailing digits or format as exponential notation. |
| **Spinbox Steppers on Identity Fields** | LOW | Confuses users and screen readers by implying identity codes can be arithmetically adjusted. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Using type='number' on postal code strips leading zeros):
```tsx
<input
  type="number"
  name="postal_code"
  placeholder="01234"
/>
```
### TSX (Using type='number' on account number risks precision truncation):
```tsx
<Input
  type="number"
  name="account_number"
/>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Text input with numeric inputMode preserves leading zeros and launches mobile number keypad):
```tsx
<input
  type="text"
  inputMode="numeric"
  pattern="[0-9]*"
  name="postal_code"
  autoComplete="postal-code"
/>
```
### TSX (Telephone field uses dedicated type='tel' with autocomplete):
```tsx
<input
  type="tel"
  name="phone_number"
  autoComplete="tel"
/>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.number-input-identity-misuse intentional exception -->
```

```tsx
// charites:ignore ux.number-input-identity-misuse intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.number-input-identity-misuse:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


