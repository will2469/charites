# browser.date-input-format-assumption

> **Rule ID:** `browser.date-input-format-assumption`
> **Severity:** `ERROR`
> **Category:** `browser`
> **Target Standards:** HTML Living Standard Section 4.10.5.1.7 (Date State - type=date), W3C RFC 3339 / ISO 8601 Normative Date Representation (YYYY-MM-DD)

---

## 1. Overview & Core Invariant

Prohibits localized string splitting assumptions on HTML5 date input values in favor of normative ISO 8601 parsing

### Core Invariant:
> **"Native <input type="date"> values are guaranteed by W3C specification to be serialized strictly as ISO 8601 (YYYY-MM-DD). Code must not split values by localized delimiters ('/' or '.')."**

---
## 2. Technical Grounding & Engine Realities

While the browser UI may render localized date pickers according to OS settings (e.g., DD/MM/YYYY in Indonesia/UK, MM/DD/YYYY in US), the programmatic 'element.value' is ALWAYS serialized in ISO 8601 format ('YYYY-MM-DD').

Splitting by '/' causes catastrophic silent failures because the delimiter never exists in 'element.value'.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Silent Date Parsing Failure** | HIGH | Splitting ISO 8601 date string by '/' returns an array of length 1, corrupting day, month, and year data sent to APIs. |
| **Cross-Locale Form Submission Corruption** | HIGH | User birth dates, appointment bookings, or legal document dates are stored as NaN, undefined, or incorrect epochs. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Splitting native date input value by '/' based on UI display assumption):
```tsx
<input
  type="date"
  onChange={(e) => {
    // BUG: e.target.value is '2026-09-06'. Splitting by '/' fails!
    const [day, month, year] = e.target.value.split('/');
    saveDate(day, month, year);
  }}
/>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Splitting native date input by normative ISO 8601 dash delimiter or using valueAsDate):
```tsx
<input
  type="date"
  onChange={(e) => {
    // Correct: ISO 8601 format YYYY-MM-DD
    const [year, month, day] = e.target.value.split('-');
    saveDate(day, month, year);
  }}
/>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore browser.date-input-format-assumption intentional exception -->
```

```tsx
// charites:ignore browser.date-input-format-assumption intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.date-input-format-assumption:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [browser Category Guide](browser).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


