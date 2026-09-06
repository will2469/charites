# performance.react-static-heavy-import

> **Rule ID:** `performance.react-static-heavy-import`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** W3C Web Performance & Initial Route Code Splitting Best Practices, React Official Documentation (Code-Splitting with React.lazy and Suspense), Chrome DevTools Lighthouse JavaScript Bundle Size Optimization Guidelines

---

## 1. Overview & Core Invariant

Mengaudit pernyataan impor statis modul berukuran besar di tingkat atas yang membengkakkan bundel JavaScript awal dan mewajibkan pemisahan kode via React.lazy() dan <Suspense>.

### Core Invariant:
> **"Heavy visualization, editing, and utility libraries must be asynchronously loaded via 'React.lazy()' and wrapped in '<Suspense>'; top-level static imports bloat critical initial bundles and degrade FCP/TBT."**

---
## 2. Technical Grounding & Engine Realities

Top-level static import statements (`import { Chart } from 'chart.js'`) force modern JavaScript bundlers (Webpack, Vite, Rollup, esbuild) to include the imported module directly in the entry-point chunk.

Heavy third-party libraries such as `monaco-editor`, `chart.js`, `echarts`, `quill`, `pdfjs-dist`, `three`, or `xlsx` often weigh several hundreds of kilobytes compressed.

Because these libraries are typically secondary to initial above-the-fold content, loading them synchronously forces mobile devices to spend hundreds of milliseconds downloading, parsing, and compiling JavaScript before the user can interact with the page.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Initial Bundle Size Bloat** | HIGH | Drastically increases First Contentful Paint (FCP) and Total Blocking Time (TBT) due to massive synchronous script downloads. |
| **Mobile CPU Decompression Bottlenecks** | MEDIUM | Consumes limited mobile device CPU and memory parsing large scripts that are not immediately rendered on initial screen view. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Mengimpor modul grafik berat secara statis di tingkat atas):
```tsx
import { Chart } from 'chart.js';

export function Dashboard() {
  return <Chart data={stats} />;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Memisahkan modul grafik menggunakan React.lazy() dan Suspense):
```tsx
import { Suspense, lazy } from 'react';
const Chart = lazy(() => import('chart.js'));

export function Dashboard() {
  return (
    <Suspense fallback={<ChartSkeleton />}>
      <Chart data={stats} />
    </Suspense>
  );
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.react-static-heavy-import intentional exception -->
```

```tsx
// charites:ignore performance.react-static-heavy-import intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.react-static-heavy-import:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [performance Category Guide](performance).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


