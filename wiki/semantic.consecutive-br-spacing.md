# semantic.consecutive-br-spacing

> **Rule ID:** `semantic.consecutive-br-spacing`
> **Severity:** `WARN`
> **Category:** `semantic`
> **Target Standards:** W3C WCAG 2.2 Technique F34 (Failure of Success Criterion 1.3.1 due to using characters to make space or using <br> tags instead of structural markup like <p>), HTML Living Standard Section 4.5.27 (The br Element), WebAIM Screen Reader Usability Survey (Navigation and Empty Breaks)

---

## 1. Overview & Core Invariant

Flags consecutive <br> tags used for paragraph spacing violating WCAG 1.3.1 (Technique F34)

### Core Invariant:
> **"Paragraph spacing and vertical separation must use structural container elements (<p>, <div>) with CSS spacing, never consecutive <br> tags."**

---
## 2. Technical Grounding & Engine Realities

The '<br>' element represents a thematic line break within a single phrase or address block (such as in postal addresses or poetry stanzas).

Charites treats consecutive '<br>' siblings as a static indicator of spacing-by-markup. When consecutive '<br>' tags are stacked to push content downward, screen readers announce each break as a blank line or acoustic break sound, forcing assistive technology users to navigate through phantom empty elements.

Furthermore, hardcoded '<br>' line heights bypass design system spacing tokens ('space-y-*', 'gap-*') and fail to scale with responsive typography or container queries. Paragraphs should instead be wrapped in semantic '<p>' tags and spaced using CSS or Tailwind vertical rhythm classes.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Assistive Technology Reading Degradation** | HIGH | Screen reader users encounter repeated 'blank' announcements that disrupt document reading comprehension. |
| **Responsive Spacing Rigidity** | MEDIUM | Hardcoded line breaks bypass design system spacing tokens, preventing dynamic vertical rhythm scaling. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Consecutive <br> tags used as paragraph spacing (Astrades dialog specimen)):
```tsx
<>
  Apakah Anda yakin memilih kandidat ini sebagai Pimpinan Utama?
  <br />
  <br />
  <strong className="text-destructive">PERINGATAN:</strong> Tindakan ini tidak dapat dibatalkan.
</>
```
### ASTRO (Multiple <br> tags simulating paragraph gaps in Astro template):
```astro
<div>
  <p>Pengumuman penting untuk seluruh staf.</p>
  <br />
  <br />
  <br />
  <p>Harap berkumpul di aula utama pukul 09:00 WIB.</p>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Semantic paragraphs wrapped in a container with Tailwind vertical spacing):
```tsx
<div className="space-y-3">
  <p>Apakah Anda yakin memilih kandidat ini sebagai Pimpinan Utama?</p>
  <p>
    <strong className="text-destructive">PERINGATAN:</strong> Tindakan ini tidak dapat dibatalkan.
  </p>
</div>
```
### TSX (Permitted single line break within address semantics):
```tsx
<address>
  Jl. Jenderal Sudirman No. 42<br />
  Jakarta Pusat, 10220<br />
  Indonesia
</address>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore semantic.consecutive-br-spacing intentional exception -->
```

```tsx
// charites:ignore semantic.consecutive-br-spacing intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  semantic.consecutive-br-spacing:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [semantic Category Guide](semantic).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


