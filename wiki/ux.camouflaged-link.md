# ux.camouflaged-link

> **Rule ID:** `ux.camouflaged-link`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** W3C Web Content Accessibility Guidelines (WCAG 2.2 SC 1.4.1 Use of Color - Level A), Gestalt Law of Similarity & Visual Affordance Principles (Norman, 2013), Nielsen Norman Group (Hyperlink Affordance Guidelines: Persistent Underlines)

---

## 1. Overview & Core Invariant

Warns when inline prose links rely solely on color without persistent underline or non-color affordance

### Core Invariant:
> **"Inline prose hyperlinks embedded within body text must provide a persistent non-color visual cue ('underline' or 'border-b') in idle state rather than relying solely on text color or hover-only transitions."**

---
## 2. Technical Grounding & Engine Realities

WCAG 2.2 Success Criterion 1.4.1 mandates that color must not be used as the only visual means of conveying information, indicating an action, prompting a response, or distinguishing a visual element.

When an inline link inside body copy removes underlines ('no-underline') or only shows underlines on hover ('hover:underline') while displaying text in primary brand color, users with color vision deficiency (protanopia, deuteranopia, tritanopia) or those using monitors with non-calibrated contrast cannot perceive the text as an interactive link.

Furthermore, according to NN/g research, static text that lacks standard underline affordance forces users to hunt and peck with the cursor, drastically diminishing reading fluency and scan efficiency. Persistent underlines or bottom border decorations provide immediate visual affordance.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Inaccessible Color-Dependent Affordance** | HIGH | Color-blind users cannot distinguish clickable inline links from regular static prose text. |
| **Reduced Reading & Scanning Fluency** | MEDIUM | Users fail to discover important inline reference links, increasing task failure rates. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Inline prose link with no-underline relying only on text color and hover):
```tsx
<p className="text-base text-neutral-700">
  Untuk informasi lebih lengkap, silakan kunjungi
  <a href="/panduan" className="text-primary hover:underline"> buku panduan warga</a>.
</p>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Persistent underline in idle state ensures immediate non-color affordance):
```tsx
<p className="text-base text-neutral-700">
  Untuk informasi lebih lengkap, silakan kunjungi
  <a href="/panduan" className="text-primary underline decoration-primary/50 hover:decoration-primary"> buku panduan warga</a>.
</p>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.camouflaged-link intentional exception -->
```

```tsx
// charites:ignore ux.camouflaged-link intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.camouflaged-link:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


