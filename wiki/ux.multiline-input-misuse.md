# ux.multiline-input-misuse

> **Rule ID:** `ux.multiline-input-misuse`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** GOV.UK Design System (Text input vs Textarea for multi-line user answers), Nielsen Norman Group (Textareas for Multi-Line Text Input), W3C WCAG 2.2 Supporting Context (Success Criterion 3.3.2 Labels or Instructions)

---

## 1. Overview & Core Invariant

Flags single-line <Input> elements used for open-ended commentary, notes, or descriptions instead of <Textarea>

### Core Invariant:
> **"Form fields intended for open-ended commentary, descriptions, or notes must use multiline <Textarea> controls, never single-line <input>."**

---
## 2. Technical Grounding & Engine Realities

Form fields intended for open-ended commentary, descriptions, reasons, or notes (identified by keywords such as 'keterangan', 'catatan', 'notes', 'deskripsi', 'description', 'alasan', 'komentar') frequently misuse single-line '<input>' controls.

When users type explanatory sentences into a single-line text box:
1. Text Overflow: The text scrolls horizontally out of view, preventing users from reviewing or verifying their responses without manual cursor scrolling.
2. Line Break Inability: Users cannot press 'Enter' to structure their thoughts into paragraphs or distinct points.
3. Usability Failure: According to GOV.UK Design System and Nielsen Norman Group form design guidelines, text inputs must only be used for short single-line entries, whereas open-ended commentary and notes require multiline '<Textarea rows={3}>' controls.

Charites inspects single-line input controls across attribute channels ('name', 'id', 'placeholder', 'aria-label') while respecting single-line qualifiers (such as 'note_title', 'reason_code', 'short_description') to ensure zero false positives on legitimate single-line fields.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Input Review Obstruction** | HIGH | Users cannot read their full input at a glance, increasing submission errors on critical notes and reasons. |
| **Cognitive Load & Text Formatting Frustration** | MEDIUM | Users are unable to insert line breaks or structure thoughts into paragraphs in single-line controls. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Single-line input for open-ended notes/keterangan (Astrades specimen)):
```tsx
<Label htmlFor="keterangan-input">Keterangan (Opsional)</Label>
<Input id="keterangan-input" placeholder="Catatan tambahan (bila ada)" />
```
### ASTRO (Single-line text input for rejection reason in Astro template):
```astro
<label for="rejection-reason">Alasan Penolakan</label>
<input id="rejection-reason" name="alasan_penolakan" type="text" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Multiline Textarea for readable notes):
```tsx
<Label htmlFor="keterangan-input">Keterangan (Opsional)</Label>
<Textarea id="keterangan-input" placeholder="Catatan tambahan (bila ada)" rows={3} />
```
### ASTRO (Multiline textarea for cancellation reason):
```astro
<label for="rejection-reason">Alasan Penolakan</label>
<textarea id="rejection-reason" name="alasan_penolakan" rows={3}></textarea>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.multiline-input-misuse intentional exception -->
```

```tsx
// charites:ignore ux.multiline-input-misuse intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.multiline-input-misuse:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


