# iOS Safari & PWA Quirks -- Deep Dive

Open this file when debugging iOS-specific issues: status bar overlap, Dynamic Island, home indicator, overscroll bounce, PWA standalone behavior, biometric prompts, or Apple-specific meta tags.

## Table of Contents

1. [Safe Area Insets -- the full picture](#1-safe-area-insets)
2. [Status Bar & Theme Color](#2-status-bar--theme-color)
3. [Dynamic Island (iPhone 14 Pro+)](#3-dynamic-island)
4. [Overscroll Bounce & Scroll Chaining](#4-overscroll-bounce--scroll-chaining)
5. [iOS PWA Standalone Mode](#5-ios-pwa-standalone-mode)
6. [iOS Input & Keyboard Quirks](#6-ios-input--keyboard-quirks)
7. [iOS-specific CSS & API support](#7-ios-specific-css--api-support)

---

## 1. Safe Area Insets

iOS exposes four environment variables for the safe area (the region not covered by notch, Dynamic Island, home indicator, or rounded corners):

```css
padding-top: env(safe-area-inset-top); /* status bar / notch */
padding-bottom: env(safe-area-inset-bottom); /* home indicator */
padding-left: env(safe-area-inset-left); /* landscape notch */
padding-right: env(safe-area-inset-right); /* landscape notch */
```

### Requirements (all three):

1. **Viewport meta must have `viewport-fit=cover`** -- without it, all `env()` values are `0`.
   ```astro
   <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
   ```
2. **Use `env()` with a fallback** -- older browsers (and some WebViews) don't support it:
   ```css
   padding-bottom: max(1rem, env(safe-area-inset-bottom));
   ```
3. **`max()` for combined padding** -- so you always have at least your design padding even when the safe-area inset is `0`:
   ```css
   padding-bottom: max(1rem, env(safe-area-inset-bottom));
   ```

### Where to apply safe areas

| Element                           | Inset to apply                      |
| :-------------------------------- | :---------------------------------- |
| Fixed top header (standalone PWA) | `top`                               |
| Fixed top header (browser)        | None -- Safari's chrome handles it  |
| Fixed/sticky bottom bar (always)  | `bottom`                            |
| Bottom navigation tab bar         | `bottom`                            |
| FAB (floating action button)      | `bottom` (and `right` in landscape) |
| Side rails in landscape           | `left`/`right`                      |

### Common bug: safe-area returns 0

If `env(safe-area-inset-*)` is `0` on a real iPhone X+:

- Missing `viewport-fit=cover` in the viewport meta (most common).
- Inside an iframe (safe areas don't propagate into iframes).
- The element is inside a `transform: translateZ(0)` parent that creates a new stacking context (rare).
- The page is loaded inside an app's WebView that doesn't forward the insets (Capacitor/Cordova require plugin configuration).

---

## 2. Status Bar & Theme Color

iOS Safari 15+ supports `theme-color` for the area above the URL bar:

```astro
<meta name="theme-color" media="(prefers-color-scheme: light)" content="#ffffff" />
<meta name="theme-color" media="(prefers-color-scheme: dark)" content="#0f172a" />
```

In **standalone PWA mode** (added to home screen), `theme-color` tints the entire status bar background. Pick a color that has enough contrast with the iOS status-bar icons (which are white in dark mode, black in light mode).

`color-scheme` tells Safari which color palettes the page supports, affecting form controls, scrollbars, and the overscroll background:

```astro
<meta name="color-scheme" content="light dark" />
```

Without this, native form controls (date picker, dropdowns) will render in light mode even when the user has dark mode on.

---

## 3. Dynamic Island

iPhone 14 Pro and later have the Dynamic Island -- a taller, active area at the top. For web content, it's handled identically to the notch via `safe-area-inset-top`. There's no special API to interact with the Island itself.

### Things to verify

- **Landscape**: the Island is on the left side. `safe-area-inset-left` becomes non-zero. Test landscape explicitly.
- **Media playback**: if your page plays video/audio, iOS may show a Now Playing widget in the Island. There's nothing to code for -- but design your header so it doesn't break if the user opens another app and comes back.
- **Live Activities**: web content can't directly post Live Activities. That requires a native app companion.

### Common bug: header content clipped

If you have a fixed top header with absolute-positioned children, the Island can clip them:

```css
/*  Child escapes the safe area */
.header {
	position: fixed;
	top: 0;
	padding-top: env(safe-area-inset-top);
}
.header .logo {
	position: absolute;
	top: 8px;
	left: 0;
} /* clipped! */

/*  Use flex, not absolute, for header content */
.header {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	padding-top: env(safe-area-inset-top);
	display: flex;
	align-items: center;
	gap: 0.75rem;
	padding-inline: 1rem;
}
```

---

## 4. Overscroll Bounce & Scroll Chaining

iOS Safari has the famous "rubber band" scroll bounce when you scroll past the top/bottom of a scrollable area. It also "chains" scroll events: when an inner scrollable element reaches its end, the parent takes over.

### Disable bounce globally

```css
:root {
	overscroll-behavior: none;
}
```

This disables both the bounce _and_ the pull-to-refresh gesture. Use this on app-like pages where the bounce feels broken (modals, dashboards, admin tools). Keep it enabled on content pages where users expect native scroll.

### Contain scroll inside modals/drawers

```css
.modal-body {
	overflow-y: auto;
	overscroll-behavior: contain; /* scroll stops at the edge -- doesn't chain to body */
}
```

**Critical for modals**: without `overscroll-behavior: contain`, scrolling inside a modal on iOS will eventually scroll the background page -- a very common bug.

### Two-finger scroll inside nested scroll containers

Older iOS Safari (15 and below) sometimes won't scroll a nested `overflow: auto` element with one finger. The legacy fix is `-webkit-overflow-scrolling: touch`. Modern Safari (16+) does this by default -- the property is now no-op, but doesn't hurt to include:

```css
.scrollable {
	overflow-y: auto;
	-webkit-overflow-scrolling: touch;
}
```

---

## 5. iOS PWA Standalone Mode

When a user adds your site to their home screen and launches it from there, Safari runs in "standalone" mode -- no URL bar, no bottom toolbar, full screen. Detect and adapt:

### Detection

```ts
// Inside a React island (client:only)
const isStandalone =
	typeof window !== "undefined" &&
	(window.matchMedia("(display-mode: standalone)").matches ||
		// iOS Safari (pre-15.4) uses a non-standard property
		(window.navigator as any).standalone === true);
```

### Required meta tags

```astro
<!-- iPhone PWA support -->
<meta name="apple-mobile-web-app-capable" content="yes" />
<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
<meta name="apple-mobile-web-app-title" content="Charites" />
<link rel="apple-touch-icon" href="/icons/ios/192.png" />
```

**`black-translucent`** makes the status bar transparent with white text -- your content extends under it. You MUST then add safe-area-top padding to your header. **`default`** gives a solid black bar with white text and no safe-area overlap.

### Common standalone bugs

1. **Status bar overlap** -- header text disappears under the status bar. Fix: `padding-top: env(safe-area-inset-top)` on the fixed header.
2. **No back button** -- iOS PWA has no browser chrome. Add an in-app back button on every page that isn't the home screen.
3. **External links open new context** -- `<a target="_blank">` opens Safari and leaves the PWA. Use `<a target="_blank" rel="noopener">` for external links, and accept this behavior.
4. **`100vh` is the full screen** -- but `100dvh` still works correctly. Stick with `dvh`.
5. **Audio must be user-initiated** -- autoplay is blocked. Always start audio from a tap.

---

## 6. iOS Input & Keyboard Quirks

### Auto-zoom on focus (16px rule)

Reiterated here because it's the most common iOS form bug: any `<input>`, `<select>`, or `<textarea>` with computed font-size < 16px triggers page zoom on focus. Fix: `text-base` (16px) on mobile for all form controls. You can shrink to `text-sm` on `sm:` and up.

### Fixed-position elements jump when keyboard opens

iOS Safari (15-) doesn't resize the layout viewport when the keyboard opens -- it pushes the visual viewport up. `position: fixed` elements anchored to the bottom (like a sticky submit button) may appear above the keyboard or behind it.

**Fix**: anchor bottom-fixed elements to `100dvh` instead of `100vh`, and add `interactive-widget=resizes-content` to the viewport meta (forces the layout to reflow).

### Date picker native UI

iOS Safari renders `<input type="date">` as a native wheel picker. The displayed format follows the device locale, but `input.value` is always `YYYY-MM-DD`. Don't try to parse the displayed text.

### File upload & camera

```tsx
<input
	type="file"
	accept="image/*"
	capture="environment" // "user" for front camera, "environment" for rear
/>
```

On iOS, `capture="environment"` opens the rear camera directly (skipping the photo library). Without `capture`, the user gets the standard picker with camera/photo-library/files options.

### `position: sticky` inside `overflow: auto` parent

Works in iOS 13+. Older versions had bugs. If you have a sticky header inside a scrollable table, verify on iOS 13+.

---

## 7. iOS-specific CSS & API support

| Feature                  | Safari iOS                        | Notes                                      |
| :----------------------- | :-------------------------------- | :----------------------------------------- |
| `dvh`/`svh`/`lvh`        | 15.4+                             | Use `@supports` to fallback for older iOS  |
| `overscroll-behavior`    | 16+                               | Below 16, you can't disable bounce via CSS |
| `env(safe-area-inset-*)` | 11.1+                             | Requires `viewport-fit=cover`              |
| `backdrop-filter`        | 9+ (prefixed `-webkit-` until 18) | Use both prefixes for safety               |
| `gap` in flexbox         | 14.5+                             | Below that, use margins                    |
| `:has()`                 | 15.4+                             | Below that, fallback to JS                 |
| `aspect-ratio`           | 15+                               | Below that, use padding-bottom hack        |
| `container` queries      | 16+                               | Below that, use media queries              |
| View Transitions API     | 18+                               | Polyfill needed for older iOS              |
| `font-display: optional` | 13.1+                             | Use to prevent invisible text              |

### `@supports` fallback for viewport units

```css
.full-height {
	height: 100vh; /* fallback */
}
@supports (height: 100dvh) {
	.full-height {
		height: 100dvh;
	}
}
```

### `backdrop-filter` prefix

```css
.glass {
	-webkit-backdrop-filter: blur(8px);
	backdrop-filter: blur(8px);
}
```
