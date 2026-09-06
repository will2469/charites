# cls.layout-trigger-transition

> **Rule ID:** `cls.layout-trigger-transition`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** W3C CSS Transitions Level 1 (transition-property), Google Core Web Vitals (CLS & Animation Frame Stability), Tailwind CSS v4 Transition Best Practices

---

## 1. Overview & Core Invariant

CSS transition targets layout-triggering geometry properties instead of GPU-composited transforms

### Core Invariant:
> **"CSS transitions must avoid animating layout-triggering geometry properties ('width', 'height', 'margin', 'padding', 'top', 'left') and instead utilize GPU-composited 'transform' or 'opacity'."**

---
## 2. Technical Grounding & Engine Realities

Transitioning geometry properties (such as width, height, padding, or positional offsets) triggers continuous CPU layout recalculations and repaints throughout the transition duration.

When geometry transitions execute on interactive elements (e.g., hover expansion or focus state enlargement), neighboring elements are continuously pushed and shifted, generating Cumulative Layout Shift (CLS).

Transitioning 'transform' (e.g. scale, translate) or 'opacity' executes entirely on the GPU compositor thread without triggering layout reflow.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Interactive Reflow Stalling** | MEDIUM | Hovering or focusing interactive elements triggers continuous layout passes, causing adjacent content to jitter and shift. |
| **Dropped Frames During Hover Transitions** | MEDIUM | Main-thread CPU layout calculations during mousemove or hover states cause micro-stutters and responsiveness lag. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### CSS (CSS declaration transitioning width directly):
```css
.sidebar {
  transition: width 300ms ease-in-out;
}
```
### TSX (Tailwind transition-all combined with hover geometry mutation):
```tsx
<div className="w-32 transition-all hover:w-64">
  Sidebar
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### CSS (CSS transition utilizing GPU-composited transform):
```css
.sidebar {
  transition: transform 300ms ease-in-out;
}
```
### TSX (Tailwind transition-transform with scale):
```tsx
<div className="w-32 transition-transform hover:scale-110">
  Sidebar
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore cls.layout-trigger-transition intentional exception -->
```

```tsx
// charites:ignore cls.layout-trigger-transition intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.layout-trigger-transition:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [cls Category Guide](cls).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


