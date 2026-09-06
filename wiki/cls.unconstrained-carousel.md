# cls.unconstrained-carousel

> **Rule ID:** `cls.unconstrained-carousel`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** W3C Cumulative Layout Shift (CLS) Metric Specification, W3C CSS Scroll Snap Module Level 1, W3C CSS Box Sizing Module Level 4 (aspect-ratio)

---

## 1. Overview & Core Invariant

Warns when carousel or slider containers lack bounded height or slide aspect-ratio constraints

### Core Invariant:
> **"Carousel and slider viewport tracks must constrain container height or bind slide items to fixed aspect ratios to prevent vertical reflow during slide transitions."**

---
## 2. Technical Grounding & Engine Realities

Horizontal scrolling tracks and carousels render dynamic collections of cards, banners, or images.

When the carousel track lacks an explicit height (e.g. 'h-64' or 'min-h-[300px]') and slides do not have locked aspect ratios, incoming slides with varying image proportions or dynamic text will force the entire container to expand or collapse vertically.

Fixing the container height or assigning 'aspect-video' / 'aspect-square' to slide items ensures layout stability throughout horizontal panning.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Vertical Container Height Jitter** | MEDIUM | Slide transitions with varying content heights push subsequent page content up and down. |
| **Cumulative Layout Shift (CLS)** | HIGH | Carousel height adjustments contribute cumulative shift points during user scrolling. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Horizontal snap container without container height or slide aspect-ratio):
```tsx
<div className="flex overflow-x-auto snap-x">
  {slides.map(s => <img key={s.id} src={s.url} alt={s.title} />)}
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Carousel container with explicit height constraint):
```tsx
<div className="flex overflow-x-auto snap-x h-64 md:h-96 w-full">
  {slides.map(s => (
    <div key={s.id} className="snap-center shrink-0 w-full h-full">
      <img src={s.url} alt={s.title} className="w-full h-full object-cover" />
    </div>
  ))}
</div>
```
### TSX (Carousel slide items locked with aspect-video utility):
```tsx
<div className="flex overflow-x-auto snap-x w-full">
  {slides.map(s => (
    <div key={s.id} className="snap-center shrink-0 w-80 aspect-video">
      <img src={s.url} alt={s.title} className="w-full h-full object-cover" />
    </div>
  ))}
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore cls.unconstrained-carousel intentional exception -->
```

```tsx
// charites:ignore cls.unconstrained-carousel intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.unconstrained-carousel:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [cls Category Guide](cls).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


