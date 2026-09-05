# Android Chrome & WebView Quirks -- Deep Dive

Open this file when debugging Android-specific issues: font boosting (text inflation), Visual Viewport API, on-screen keyboard (OSK) resize, theme-color, system back button, or WebView (Capacitor/Tauri/WebView-based apps).

## Table of Contents

1. [Font Boosting (Text Inflation)](#1-font-boosting-text-inflation)
2. [Visual Viewport API & OSK Resize](#2-visual-viewport-api--osk-resize)
3. [System Back Button](#3-system-back-button)
4. [Theme Color & Chrome Custom Tabs](#4-theme-color--chrome-custom-tabs)
5. [Android PWA / TWA](#5-android-pwa--twa)
6. [WebView App Quirks (Capacitor / Tauri)](#6-webview-app-quirks)
7. [Text Selection & Long-Press Menu](#7-text-selection--long-press-menu)

---

## 1. Font Boosting (Text Inflation)

Android Chrome automatically **enlarges text** on small screens to make it readable without zooming. The algorithm uses heuristics: it inflates text in elements that look like body copy (paragraphs, list items, divs with text) but skips form controls, headings, and elements with explicit `line-height`.

### When it breaks your layout

- A button label suddenly becomes 1.5× larger than expected, breaking the layout.
- Text in a flex row wraps unexpectedly.
- A `<span>` inside a `<p>` inflates at a different rate than the parent.

### Disabling font boosting

```css
/* Globally -- stops Android from inflating any text */
html {
	-webkit-text-size-adjust: none; /* Old WebKit */
	text-size-adjust: none; /* Standard, but Chrome ignores `none` */
	text-size-adjust: 100%; /* Chrome accepts 100% */
}

/* On specific elements that look like body copy but shouldn't inflate */
.button-label,
.nav-item {
	text-size-adjust: 100%;
}
```

**Trade-off**: Disabling font boosting entirely hurts accessibility -- users who rely on the inflation for readability lose it. The better approach is to design with **rem** units and let the user's font-size preference scale the whole layout, then explicitly opt out only on UI chrome elements (buttons, nav).

### Best practice: use rem, respect user font scale

```css
html {
	font-size: 100%; /* 1rem = 16px at default scale, scales with Android font setting */
}
body {
	font-size: 1rem;
}
```

Test your layout with Android's **Settings → Accessibility → Font size → Largest**. If anything overflows, fix it -- don't disable inflation.

---

## 2. Visual Viewport API & OSK Resize

Since Chrome 108, the Android on-screen keyboard resizes the **visual viewport** only, not the layout viewport. This means:

- `window.innerHeight` stays the same when the keyboard opens.
- `100vh` / `100dvh` containers stay the same height.
- `position: fixed; bottom: 0` elements are pushed up by the keyboard if `interactive-widget=resizes-content` is set, otherwise they're hidden behind the keyboard.

### Detecting keyboard open

```ts
function useKeyboardOpen() {
	const [open, setOpen] = useState(false);

	useEffect(() => {
		if (!window.visualViewport) return;

		const onResize = () => {
			// Visual viewport shrinks when keyboard opens
			const keyboardHeight = window.innerHeight - window.visualViewport!.height;
			setOpen(keyboardHeight > 100); // 100px threshold to avoid URL bar noise
		};

		window.visualViewport!.addEventListener("resize", onResize);
		return () => window.visualViewport!.removeEventListener("resize", onResize);
	}, []);

	return open;
}
```

### Scrolling focused input into view

When a user taps an input near the bottom of the screen, the keyboard may cover it. Use `scrollIntoView` on focus:

```tsx
function Input({ ...props }) {
	const ref = useRef<HTMLInputElement>(null);

	return (
		<input
			ref={ref}
			onFocus={() => {
				// Wait for keyboard to start animating
				setTimeout(() => {
					ref.current?.scrollIntoView({
						behavior: "smooth",
						block: "center",
					});
				}, 300);
			}}
			{...props}
		/>
	);
}
```

**Why the 300ms delay**: Chrome animates the keyboard in over ~250ms. Calling `scrollIntoView` immediately scrolls before the layout settles, then the keyboard covers the field anyway.

### Sticky bottom bar that should hide above keyboard

```tsx
function StickySubmitBar() {
	const keyboardOpen = useKeyboardOpen();

	return (
		<footer
			className="sticky bottom-0 transition-transform duration-200"
			style={{
				// Push the bar up by the keyboard height
				paddingBottom:
					keyboardOpen && window.visualViewport
						? `${window.innerHeight - window.visualViewport.height}px`
						: "max(0.75rem, env(safe-area-inset-bottom))",
			}}
		>
			<Button>Simpan</Button>
		</footer>
	);
}
```

Or simpler: with `interactive-widget=resizes-content` in the viewport meta, the layout reflows and the sticky bar naturally sits above the keyboard. Prefer this approach when possible.

---

## 3. System Back Button

Android has a hardware/system back button. By default, it navigates browser history. In an SPA-like Astro app (with View Transitions / ClientRouter), you need to handle it explicitly for things like:

- Closing an open modal/drawer when back is pressed (instead of navigating away)
- Closing an expanded menu
- Stepping back through an in-app wizard

### Pattern: virtual history entries

```ts
function useBackHandler(onBack: () => void) {
	useEffect(() => {
		// Push a virtual state so the back button closes the modal instead of leaving
		history.pushState({ modal: true }, "");

		const onPop = (e: PopStateEvent) => {
			onBack();
			// Don't push again -- let this pop happen
		};

		window.addEventListener("popstate", onPop);
		return () => {
			window.removeEventListener("popstate", onPop);
			// If we unmount without the back button being pressed, clean up the virtual entry
			if (history.state?.modal) history.back();
		};
	}, [onBack]);
}
```

This pattern is what native mobile apps do (and what iOS gives you for free with the swipe-back gesture). Use it for any modal/drawer that should close on back.

---

## 4. Theme Color & Chrome Custom Tabs

Android Chrome 93+ respects `theme-color` for the browser chrome (URL bar + status bar area):

```astro
<meta name="theme-color" media="(prefers-color-scheme: light)" content="#ffffff" />
<meta name="theme-color" media="(prefers-color-scheme: dark)" content="#0f172a" />
```

When your site is opened inside a **Chrome Custom Tab** (a common pattern when an Android app links to your site), `theme-color` tints the top bar -- a cheap way to look "integrated" without an SDK.

### Status bar (Android system) color

The system status bar color is set by the launcher / OS theme, not by your site -- `theme-color` only affects the browser/WebView chrome. To control the system status bar, you'd need a PWA manifest's `theme_color` and `display: standalone`.

---

## 5. Android PWA / TWA

For a true native-feeling Android PWA:

### Manifest requirements

```json
{
	"name": "Charites",
	"short_name": "Charites",
	"display": "standalone",
	"theme_color": "#0f172a",
	"background_color": "#ffffff",
	"orientation": "portrait-primary",
	"icons": [
		{ "src": "/icons/android/192.png", "sizes": "192x192", "type": "image/png" },
		{ "src": "/icons/android/512.png", "sizes": "512x512", "type": "image/png" },
		{
			"src": "/icons/android/maskable-512.png",
			"sizes": "512x512",
			"type": "image/png",
			"purpose": "maskable"
		}
	]
}
```

**`purpose: "maskable"`** is critical on Android -- without a maskable icon, the launcher may render your icon inside a white circle. Maskable icons have a "safe zone" in the center 80%; design your icon to be legible when cropped to a circle, squircle, or rounded square.

**`orientation: "portrait-primary"`** locks the PWA to portrait. Use only when landscape breaks your UI (e.g. an admin app designed specifically for portrait). Allow free rotation otherwise.

### TWA (Trusted Web Activity)

For shipping as a real Android APK, wrap your PWA in a TWA using Bubblewrap or PWABuilder. No code changes needed -- just a `assetlinks.json` on your domain to prove ownership. The TWA uses Chrome under the hood but renders with no browser chrome.

---

## 6. WebView App Quirks (Capacitor / Tauri)

When your Astro app is wrapped in a native shell (Capacitor, Tauri, Cordova, React Native WebView), several things change:

### Safe area insets

WebViews don't expose safe-area insets by default. In Capacitor, install `@capacitor/status-bar` and `@capacitor/keyboard` plugins and forward the values:

```ts
import { StatusBar } from "@capacitor/status-bar";

// In a top-level useEffect
StatusBar.setStyle({ style: Style.Default });
// Safe area insets are then forwarded to CSS env() automatically
```

### Keyboard

Capacitor's Keyboard plugin overrides the browser's keyboard behavior. Configure it in `capacitor.config.json`:

```json
{
	"plugins": {
		"Keyboard": {
			"resize": "body",
			"resizeOnFullScreen": true
		}
	}
}
```

`resize: "body"` resizes the layout (like `interactive-widget=resizes-content`). `resize: "native"` keeps the layout fixed and pushes the visual viewport (Android default). Choose based on whether your app uses `position: fixed; bottom: 0` elements.

### File system & camera

Inside a WebView, `input[type="file"]` requires native permission handlers. Capacitor handles this automatically with `@capacitor/camera` and `@capacitor/filesystem` plugins -- don't try to use the raw `<input type="file">` for anything beyond simple uploads.

### JavaScript bridge timing

If your Astro app calls a native bridge (e.g. `window.AndroidBridge.something()`), the bridge may not be ready when your component mounts. Wait for a deviceready-style event:

```ts
await new Promise<void>((resolve) => {
	if (window.AndroidBridge) return resolve();
	document.addEventListener("bridgeReady", () => resolve(), { once: true });
});
```

---

## 7. Text Selection & Long-Press Menu

Android Chrome shows a context menu on long-press of images (offering "Save image", "Share image") and selects text on long-press of text. This is normal browser behavior -- don't disable it globally.

### When to disable

- **UI chrome elements** (icons, buttons, drag handles): prevent the long-press menu from interfering with drag-and-drop or context actions.
- **Card layouts where long-press opens a custom menu**: same.

```css
.no-select {
	-webkit-user-select: none;
	user-select: none;
	-webkit-touch-callout: none; /* iOS too */
}
```

Don't apply `user-select: none` to `<input>`, `<textarea>`, or content the user needs to copy (NIK, phone numbers, error messages).

### Tap-and-hold to confirm destructive actions

A native-feel pattern: long-press a "Delete" button for 1.5s to confirm. Better than a `confirm()` dialog on mobile. Implement with `onTouchStart` / `onPointerDown` + a timeout.
