# Lcp Rules (`lcp`)

The `lcp` category contains static analysis rules for code quality, architectural constraints, and design system governance.

---

## Category Rule Index

| Rule ID | Severity | Summary | Full Specification | Status |
| :--- | :---: | :--- | :--- | :---: |
| `lcp.client-only-lcp-content` | `WARN` | Above-the-fold hero island declared with 'client:only' without an SSR fallback slot, eliminating server HTML and delaying LCP until client-side bundle execution | [`lcp.client-only-lcp-content`](lcp.client-only-lcp-content) | `enabled` |
| `lcp.heavy-raster-lcp-asset` | `WARN` | LCP candidate image uses legacy uncompressed raster format (.png, .bmp, .tiff, .gif); modern formats like WebP or AVIF should be served to reduce transfer size | [`lcp.heavy-raster-lcp-asset`](lcp.heavy-raster-lcp-asset) | `enabled` |
| `lcp.image-source-density-mismatch` | `INFO` | Fixed-dimension LCP candidate image lacks aligned '1x, 2x' pixel density descriptors in 'srcset', risking blurry rendering or unoptimized asset delivery on high-DPI screens | [`lcp.image-source-density-mismatch`](lcp.image-source-density-mismatch) | `enabled` |
| `lcp.lazy-loaded-lcp-image` | `ERROR` | Critical above-the-fold LCP candidate image has loading="lazy", delaying resource discovery and load initiation | [`lcp.lazy-loaded-lcp-image`](lcp.lazy-loaded-lcp-image) | `enabled` |
| `lcp.missing-lcp-image-preload` | `INFO` | Delayed-discovery LCP image lacks <link rel="preload" as="image"> in document head to initiate early asset transfer | [`lcp.missing-lcp-image-preload`](lcp.missing-lcp-image-preload) | `enabled` |
| `lcp.oversized-lcp-resource-selection` | `WARN` | Fluid responsive LCP candidate image lacks responsive 'srcset' and 'sizes' attributes, forcing mobile viewports to download oversized desktop assets | [`lcp.oversized-lcp-resource-selection`](lcp.oversized-lcp-resource-selection) | `enabled` |
| `lcp.undiscoverable-lcp-image` | `WARN` | Above-the-fold hero container loads primary image via CSS background without <link rel="preload"> in document head | [`lcp.undiscoverable-lcp-image`](lcp.undiscoverable-lcp-image) | `enabled` |
| `lcp.unhinted-lcp-image-priority` | `WARN` | Above-the-fold LCP candidate image lacks fetchpriority="high", delaying bandwidth allocation in early network stream | [`lcp.unhinted-lcp-image-priority`](lcp.unhinted-lcp-image-priority) | `enabled` |

---
## How the Lcp Analysis Pipeline Works

The `lcp` engine applies static analysis checks against component source code:

```mermaid
flowchart LR
    TargetFiles["Target Files (*.astro, *.tsx)"] --> Parser["Leaf IR AST Parser"]
    Parser --> Engine["Rule Evaluator Engine"]
    Engine --> Check{"Evaluate Invariant"}
    Check -- "Compliant" --> Safe["Pass"]
    Check -- "Violation" --> Diag["Diagnostic: lcp.*"]
```

### Pipeline Flow:
1. **AST Node Traversal:** Scans target template files into normalized intermediate representation.
2. **Invariant Assertion:** Validates structural and semantic invariants.
3. **Diagnostic Reporting:** Emits structured diagnostics for non-compliant patterns.

---

## How Lcp Tests Work (Verification Harness)

All rules in `lcp` are verified using the canonical 1-SSOT Tri-Corpus (`tests/correctness/lcp.*/`) encompassing Positive (P1-P5), Negative (N1-N5), and Adversarial (A1-A7) fixture matrices.
