# lcp.missing-critical-origin-hint

> **Rule ID:** `lcp.missing-critical-origin-hint`
> **Severity:** `INFO`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay), W3C Resource Hints (preconnect and dns-prefetch specification), Web Performance Working Group Network Socket Pre-warming Guidelines

---

## 1. Overview & Core Invariant

Critical LCP visual asset loaded from third-party CDN origin without '<link rel="preconnect">' connection hint in '<head>'

### Core Invariant:
> **"Critical visual LCP assets hosted on external cross-origin domains should have an early '<link rel="preconnect">' hint in '<head>' to eliminate DNS, TCP, and TLS socket handshake round-trips."**

---
## 2. Technical Grounding & Engine Realities

When an above-the-fold hero image is loaded from a third-party domain (e.g. 'images.unsplash.com', 'cdn.shopify.com', or 'res.cloudinary.com'), the browser cannot initiate the network connection until the `<img>` tag is discovered and parsed.

Establishing a new HTTPS connection to an external origin requires three sequential round-trips: DNS resolution, TCP three-way handshake, and TLS cryptographic negotiation, adding 150ms to 400ms of idle latency on cellular connections.

Declaring `<link rel="preconnect" href="https://images.unsplash.com" />` in `<head>` instructs the browser to open the socket in the background during early HTML document streaming.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Connection Handshake Latency Spike** | MEDIUM | Adds 150ms-400ms of socket negotiation delay to LCP Resource Load Delay before image bytes start streaming. |
| **Cellular Round-Trip Time Penalty** | LOW | Multi-RTT connection setup noticeably degrades mobile performance metrics in emerging market networks. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### HTML (Critical hero image loaded from external CDN without preconnect hint in head):
```html
<head>
  <title>Product Showcase</title>
</head>
<body>
  <img src="https://images.unsplash.com/photo-hero" fetchpriority="high" data-perf-role="hero" />
</body>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### HTML (Preconnect hint declared in head to pre-warm CDN connection sockets early):
```html
<head>
  <title>Product Showcase</title>
  <link rel="preconnect" href="https://images.unsplash.com" />
</head>
<body>
  <img src="https://images.unsplash.com/photo-hero" fetchpriority="high" data-perf-role="hero" />
</body>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore lcp.missing-critical-origin-hint intentional exception -->
```

```tsx
// charites:ignore lcp.missing-critical-origin-hint intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.missing-critical-origin-hint:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [lcp Category Guide](lcp).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


