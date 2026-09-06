# theme.image-theme-hardcode

> **Rule ID:** `theme.image-theme-hardcode`
> **Severity:** `WARN`
> **Category:** `theme`
> **Target Standards:** WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast), W3C Responsive Images & Art Direction Specification, Tailwind CSS Dark Mode Graphic Switching Guidelines

---

## 1. Overview & Core Invariant

Detects graphic assets and logos in img tags lacking dark mode theme adaptation

### Core Invariant:
> **"Graphic assets, logos, and diagrams in img tags must provide theme-adaptive variants via picture, dark: utility classes, or invert filters."**

---
## 2. Technical Grounding & Engine Realities

Embedding graphical assets (such as brand logos, SVG diagrams, and charts) via static <img> tags without dark mode adaptation triggers severe visual breakage:

1. Asset Invisibility: A dark or black logo rendered against a dark mode background becomes completely invisible.
2. Inverted Eye-Strain: High-glare white background diagrams blast excessive light on dark UI themes.
3. Inflexible Art Direction: Projects without responsive image pairing cannot tailor vector artwork to dark contrast requirements.

Charites enforces theme-aware graphic handling using dark:hidden / dark:block class pairs, dark:invert filters, or responsive <picture> elements.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Asset Disappearance in Dark Themes** | MEDIUM | Brand logos, technical diagrams, and icon artwork become illegible on dark surfaces. |
| **Non-text Contrast Failure (WCAG 1.4.11)** | MEDIUM | Visual cues necessary for interface understanding fail accessibility contrast requirements. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Static logo in img tag without dark mode variant in TSX):
```tsx
<img src="/images/logo-black.svg" alt="Company Logo" />
```
### ASTRO (Vector architecture diagram without theme switching in Astro):
```astro
<img src="/assets/diagram.svg" alt="Architecture Flow" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Theme-paired image switching using Tailwind dark utilities):
```tsx
<>
  <img src="/images/logo-light.svg" className="dark:hidden" alt="Logo" />
  <img src="/images/logo-dark.svg" className="hidden dark:block" alt="Logo" />
</>
```
### ASTRO (Using dark:invert filter for vector diagrams):
```astro
<img src="/assets/diagram.svg" class="dark:invert" alt="Architecture Flow" />
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore theme.image-theme-hardcode intentional exception -->
```

```tsx
// charites:ignore theme.image-theme-hardcode intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  theme.image-theme-hardcode:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [theme Category Guide](theme).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


