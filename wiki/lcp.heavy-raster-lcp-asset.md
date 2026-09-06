# lcp.heavy-raster-lcp-asset

> **Rule ID:** `lcp.heavy-raster-lcp-asset`
> **Severity:** `WARN`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Duration), W3C Web Performance Working Group Media Optimization Guidelines, IETF AVIF / WebP Image Compression Standards

---

## 1. Overview & Core Invariant

LCP candidate image uses legacy uncompressed raster format (.png, .bmp, .tiff, .gif); modern formats like WebP or AVIF should be served to reduce transfer size

### Core Invariant:
> **"Above-the-fold LCP candidate images must utilize next-generation compressed formats (WebP, AVIF) rather than legacy uncompressed raster formats (.png, .bmp, .tiff, .gif) to minimize byte transfer payload."**

---
## 2. Technical Grounding & Engine Realities

Serving high-resolution photographs or hero imagery in legacy raster formats such as PNG or uncompressed BMP results in massive byte payloads (often 2MB-5MB per image).

Next-generation formats such as WebP and AVIF provide superior lossy and lossless compression algorithms, reducing image file sizes by 30% to 70% compared to PNG and JPEG without perceptual visual degradation.

For above-the-fold hero images that dictate the LCP metric, reducing file transfer size directly accelerates the Resource Load Duration phase over bandwidth-constrained networks.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Resource Load Duration Delay** | HIGH | Downloading uncompressed 2MB-5MB PNG/BMP images over 4G/cellular connections introduces 800ms-3000ms delay to LCP. |
| **Memory Footprint & GPU Texture Pressure** | MEDIUM | Large uncompressed raster graphics consume excessive client RAM and GPU texture memory during decode and compositing. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Critical hero image served as an uncompressed 3MB PNG file):
```tsx
<section className="hero-section" data-perf-role="hero">
  <img
    src="/assets/hero-banner.png"
    alt="Hero Showcase"
    fetchpriority="high"
    className="w-full h-auto"
  />
</section>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Hero image converted to compressed modern WebP format):
```tsx
<section className="hero-section" data-perf-role="hero">
  <img
    src="/assets/hero-banner.webp"
    alt="Hero Showcase"
    fetchpriority="high"
    className="w-full h-auto"
  />
</section>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.heavy-raster-lcp-asset intentional exception -->
```

```tsx
// charites:ignore lcp.heavy-raster-lcp-asset intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.heavy-raster-lcp-asset:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


