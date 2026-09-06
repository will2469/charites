# responsive.image-overflow

> **Rule ID:** `responsive.image-overflow`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** HTML Living Standard (Embedded Content: img, video, svg), Web.dev Responsive Media & Core Web Vitals (CLS Prevention)

---

## 1. Overview & Core Invariant

Warns when media elements with large fixed dimensions lack responsive max-w-full scaling

### Core Invariant:
> **"Media elements with explicit width dimensions exceeding 320px must declare 'max-w-full' or 'w-full' to prevent horizontal viewport tearing on mobile screens."**

---
## 2. Technical Grounding & Engine Realities

Specifying explicit 'width' and 'height' attributes on media elements is recommended for Core Web Vitals to reserve aspect ratio boxes and prevent Cumulative Layout Shift (CLS).

However, when large static dimensions (e.g. width={1200}) lack responsive CSS scaling ('max-w-full h-auto'), mobile browsers render the media at full absolute physical pixel width, breaking outside narrow 360px viewports and causing severe horizontal scrolling.

Applying 'max-w-full h-auto' preserves CLS aspect ratio benefits while ensuring the media smoothly downsizes to fit compact screens.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Viewport Tearing** | MEDIUM | Large unconstrained images expand outside mobile viewport boundaries, forcing horizontal scrollbars. |
| **Aspect Ratio Distortion** | LOW | Images constrained by height but not width stretch disproportionately on narrow viewports. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Media element with large width attribute lacking max-w-full):
```tsx
<img src="/hero-desa.jpg" width={1200} height={800} alt="Pemandangan Desa" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Responsive media element with max-w-full and h-auto):
```tsx
<img className="max-w-full h-auto" src="/hero-desa.jpg" width={1200} height={800} alt="Pemandangan Desa" />
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore responsive.image-overflow intentional exception -->
```

```tsx
// charites:ignore responsive.image-overflow intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.image-overflow:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [responsive Category Guide](responsive).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


