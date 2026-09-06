# cls.unsized-image

> **Rule ID:** `cls.unsized-image`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** W3C Cumulative Layout Shift (CLS) Metric Specification, Google Core Web Vitals Guidelines (Target CLS < 0.1), W3C CSS Box Sizing Module Level 4 (aspect-ratio), Astro Docs: Image Optimization (astro:assets)

---

## 1. Overview & Core Invariant

Warns when image elements lack explicit dimensions, aspect-ratio, or Tailwind box sizing

### Core Invariant:
> **"Image elements must establish a statically inferable reserved rendering box via explicit width/height attributes, CSS aspect-ratio, or Tailwind sizing utilities before the binary asset is downloaded."**

---
## 2. Technical Grounding & Engine Realities

When browsers parse an <img> tag without explicit dimensions or an aspect-ratio reservation, the layout engine initially allocates a 0x0 pixel box.

Once the remote image file is fetched and decoded, the browser performs a sudden reflow to accommodate the intrinsic image geometry, pushing surrounding content downward. This layout instability directly penalizes Cumulative Layout Shift (CLS) scores.

Specifying width and height attributes or utilizing Tailwind 'aspect-video' / 'aspect-square' allows modern browsers to compute the aspect ratio before network I/O completes, eliminating visual jank.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cumulative Layout Shift (CLS)** | HIGH | Unsized images push subsequent content down upon load, degrading Core Web Vitals and SEO rankings. |
| **Accidental User Mis-clicks** | MEDIUM | Users attempting to tap links or buttons near loading images accidentally trigger shifted elements. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Image with fluid width but missing height or aspect-ratio reservation):
```tsx
<img src={heroUrl} alt="Hero Banner" className="w-full h-auto" />
```
### ASTRO (Standard img tag lacking width and height attributes):
```astro
<img src="/pemandangan-desa.jpg" alt="Pemandangan Desa" class="rounded-lg" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Image with explicit numeric width and height attributes):
```tsx
<img src={heroUrl} alt="Hero Banner" width={1200} height={600} className="w-full h-auto" />
```
### TSX (Image with Tailwind v4 aspect-ratio utility):
```tsx
<img src={heroUrl} alt="Hero Banner" className="w-full aspect-video object-cover" />
```
### TSX (Avatar image with explicit width and height sizing utilities):
```tsx
<img src={avatarUrl} alt="Avatar" className="w-10 h-10 rounded-full" />
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore cls.unsized-image intentional exception -->
```

```tsx
// charites:ignore cls.unsized-image intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.unsized-image:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [cls Category Guide](cls).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


