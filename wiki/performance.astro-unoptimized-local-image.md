# performance.astro-unoptimized-local-image

> **Rule ID:** `performance.astro-unoptimized-local-image`
> **Severity:** `INFO`
> **Category:** `performance`
> **Target Standards:** Astro Asset Pipeline Best Practices ('astro:assets' Image & Picture), Core Web Vitals Largest Contentful Paint (LCP) Image Payload Optimization, W3C Next-Gen Responsive Image Delivery Standards (WebP/AVIF)

---

## 1. Overview & Core Invariant

Menganjurkan pemakaian komponen <Image /> dari astro:assets pada gambar lokal guna mengaktifkan konversi format modern dan kompresi build.

### Core Invariant:
> **"Local raster image assets in Astro templates should be rendered via '<Image />' from 'astro:assets' rather than raw '<img>' tags to leverage automated build-time format conversion and dimension inference."**

---
## 2. Technical Grounding & Engine Realities

Astro provides a native image optimization pipeline through the `astro:assets` module.

Using a raw HTML `<img>` tag pointing to a local file path (`src="../assets/banner.png"`) completely bypasses this pipeline, serving uncompressed, legacy formats (PNG/JPEG) with no automatic width/height dimension injection.

Migrating to `<Image />` allows Astro to automatically convert images to AVIF/WebP, generate responsive srcset attributes, and prevent Cumulative Layout Shift (CLS) by inferring exact dimensions at build time.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Inflated Asset Payload** | LOW | Serves unoptimized PNG/JPEG images that are 40-70% larger than modern WebP/AVIF equivalents. |
| **Missing Intrinsic Aspect Ratio** | LOW | Raw img tags without width and height attributes cause layout shifts during image load. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Tag img mentah melewatkan kompresi build-time Astro):
```astro
<!-- Advisory: Tag img mentah pada path lokal -->
<img src="../assets/product-hero.png" alt="Produk Baru" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Memanfaatkan komponen Image bawaan astro:assets):
```astro
---
import { Image } from 'astro:assets';
import productImg from '../assets/product-hero.png';
---
<Image src={productImg} alt="Produk Baru" />
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.astro-unoptimized-local-image intentional exception -->
```

```tsx
// charites:ignore performance.astro-unoptimized-local-image intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.astro-unoptimized-local-image:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [performance Category Guide](performance).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


