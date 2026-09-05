# Mobile Accessibility -- WCAG 2.2 Deep Dive

Open this file when auditing a page for accessibility, building a focus-trap, supporting VoiceOver/TalkBack, or implementing the new WCAG 2.2 mobile success criteria.

## Table of Contents

1. [WCAG 2.2 -- Mobile-specific Success Criteria](#1-wcag-22-mobile-specific-success-criteria)
2. [VoiceOver (iOS) -- testing methodology](#2-voiceover-ios)
3. [TalkBack (Android) -- testing methodology](#3-talkback-android)
4. [Focus Management & Focus Traps](#4-focus-management--focus-traps)
5. [ARIA Patterns for Common Mobile Components](#5-aria-patterns)
6. [Announcing Dynamic Content](#6-announcing-dynamic-content)
7. [Keyboard & Switch Access](#7-keyboard--switch-access)

---

## 1. WCAG 2.2 -- Mobile-specific Success Criteria

WCAG 2.2 (Oct 2023) added several criteria that specifically affect mobile design:

### 2.5.8 Target Size (Minimum) -- Level AA

> Targets are at least **24×24 CSS px**, OR there's adequate spacing, OR the target is inline / in a user-agent-default style.

AAA: 44×44 CSS px (same as Apple HIG).

**Practical**: aim for 44×44 (`min-h-11 min-w-11` in Tailwind). It satisfies both AA and AAA.

### 2.5.7 Dragging Movements -- Level AA

> Anything that can be done with a drag must also be doable with a single tap.

If you have a drag-to-reorder list or a slider, you MUST also provide tap-based controls (up/down buttons, +/- buttons). Pure drag-only interactions fail this criterion.

### 3.3.7 Redundant Entry -- Level A

> Information previously entered by the user is auto-populated or available for selection, not re-entered.

If a multi-step form asks for the user's email on step 1 and again on step 3, you must auto-fill or offer a dropdown. Don't make users type the same thing twice.

### 3.3.8 Accessible Authentication (Minimum) -- Level AA

> No cognitive function test (memory test, puzzle) is required for login, OR an alternative is provided.

CAPTCHAs that require "select all traffic lights" fail this. Provide WebAuthn / passkey / biometric options. Don't force password re-entry with masked characters.

### 4.1.3 Status Messages -- Level AA

> Status messages can be programmatically determined without receiving focus.

Use `role="status"` or `aria-live="polite"` for toasts, "Saved successfully" messages, and form-validation errors that appear after submission.

---

## 2. VoiceOver (iOS)

VoiceOver is iOS's built-in screen reader. Test it on every major feature.

### How to enable

**Settings → Accessibility → VoiceOver → On.**

Practice the gestures first:

- **Single tap**: announce the element under your finger
- **Double tap**: activate the focused element
- **Swipe right/left**: move to next/previous element
- **Three-finger swipe up/down**: scroll page
- **Two-finger scrub (Z shape)**: dismiss modal/dialog
- **Three-finger triple tap**: toggle speech on/off

### Common VoiceOver bugs in React apps

1. **`onClick` on a `<div>`** -- VoiceOver won't announce it as interactive. Use `<button>` or `<a>`, or add `role="button" tabIndex={0}` and handle `onKeyDown` for Enter/Space.
2. **Icon-only buttons without `aria-label`** -- VoiceOver reads "button" with no context. Add `aria-label="Tutup menu"`.
3. **Modal without `aria-modal="true"`** -- VoiceOver reads background content through the modal. Add `aria-modal="true"` and `role="dialog"`.
4. **Dynamic content that appears** without an aria-live region -- VoiceOver doesn't announce it. Wrap in `aria-live="polite"`.
5. **Visual-only feedback** (" Saved" as a green checkmark with no text) -- VoiceOver is silent. Add visually-hidden text or `aria-label`.
6. **Decorative images without `alt=""`** -- VoiceOver reads the filename. Always set `alt=""` for decorative images.

### rotor navigation

VoiceOver users use the "rotor" (a two-finger twist gesture) to navigate by element type: headings, links, form controls, landmarks. Use semantic HTML (`<h1>`-`<h6>`, `<nav>`, `<main>`, `<header>`, `<footer>`) so rotor navigation works.

---

## 3. TalkBack (Android)

TalkBack is Android's built-in screen reader. Largely similar to VoiceOver but with different gestures.

### How to enable

**Settings → Accessibility → TalkBack → On.**

Gestures:

- **Single tap**: announce element
- **Double tap**: activate
- **Swipe right/left**: next/previous element
- **Three-finger swipe**: scroll (note: VoiceOver uses two fingers)
- **Swipe up then down**: navigate by element type

### Android-specific considerations

1. **`accessibilityLiveRegion`** -- Android's equivalent of `aria-live`. The HTML `aria-live` attribute works in Chrome → TalkBack, but make sure to test.
2. **Touch exploration mode** is always on when TalkBack is active -- your `:hover` styles never fire. Don't rely on hover-only interactions.
3. **`<select>` rendering** -- Android Chrome renders dropdowns as a bottom sheet by default, which TalkBack navigates well. Don't replace native `<select>` with custom dropdowns unless you implement full ARIA combobox semantics.
4. **WebView apps**: TalkBack works inside WebViews, but you need to enable it in the native app's `AndroidManifest.xml` (`android:importantForAccessibility`).

---

## 4. Focus Management & Focus Traps

### Skip-to-content link

The first focusable element on every page should be a skip link:

```tsx
<a
  href="#main"
  className="sr-only focus:not-sr-only focus:absolute focus:top-2 focus:left-2 focus:z-50
             focus:px-4 focus:py-2 focus:bg-white focus:shadow-lg focus:rounded"
>
  Lewati ke konten utama
</a>
<main id="main">{/* ... */}</main>
```

### Focus trap for modals/drawers

When a modal opens, focus must:

1. Move into the modal (typically to the first focusable element or the title).
2. Stay inside the modal (Tab/Shift+Tab cycle within it).
3. Return to the trigger element when closed.

**Use `focus-trap-react` or `@radix-ui/react-dialog`** -- both handle this correctly. If you must implement manually:

```tsx
function useFocusTrap(containerRef: Ref<HTMLElement>, open: boolean) {
	useEffect(() => {
		if (!open || !containerRef.current) return;

		const container = containerRef.current;
		const previouslyFocused = document.activeElement as HTMLElement | null;

		// Query focusable elements
		const selector =
			'a[href], button:not([disabled]), textarea, input, select, [tabindex]:not([tabindex="-1"])';
		const focusable = () => Array.from(container.querySelectorAll<HTMLElement>(selector));

		// Move focus into modal
		const elements = focusable();
		if (elements.length) elements[0].focus();

		const onKey = (e: KeyboardEvent) => {
			if (e.key !== "Tab") return;
			const list = focusable();
			if (list.length === 0) return;

			const first = list[0];
			const last = list[list.length - 1];
			const active = document.activeElement;

			if (e.shiftKey && active === first) {
				e.preventDefault();
				last.focus();
			} else if (!e.shiftKey && active === last) {
				e.preventDefault();
				first.focus();
			}
		};

		document.addEventListener("keydown", onKey);
		return () => {
			document.removeEventListener("keydown", onKey);
			previouslyFocused?.focus();
		};
	}, [open, containerRef]);
}
```

### Focus on route change

In Astro with View Transitions, focus doesn't automatically move to the new page content. After navigation, move focus to the `<main>` heading:

```ts
// In a top-level script
document.addEventListener("astro:after-swap", () => {
	const main = document.querySelector("main");
	const heading = main?.querySelector("h1");
	if (heading) {
		heading.setAttribute("tabindex", "-1");
		(heading as HTMLElement).focus();
	}
});
```

---

## 5. ARIA Patterns

### Accordion

```tsx
<div>
	<h3>
		<button
			aria-expanded={isOpen}
			aria-controls="panel-1"
			id="accordion-1"
			onClick={() => setIsOpen(!isOpen)}
		>
			Section title
		</button>
	</h3>
	<div id="panel-1" role="region" aria-labelledby="accordion-1" hidden={!isOpen}>
		{/* content */}
	</div>
</div>
```

### Tab interface

```tsx
<div role="tablist">
  <button role="tab" id="tab-1" aria-selected="true" aria-controls="panel-1">Tab 1</button>
  <button role="tab" id="tab-2" aria-selected="false" aria-controls="panel-2" tabIndex={-1}>Tab 2</button>
</div>
<div role="tabpanel" id="panel-1" aria-labelledby="tab-1">...</div>
<div role="tabpanel" id="panel-2" aria-labelledby="tab-2" hidden>...</div>
```

Keyboard: ArrowLeft/Right to move between tabs, Home/End for first/last, Tab to enter the panel.

### Toast / status message

```tsx
<div role="status" aria-live="polite" className="sr-only">
	{toastMessage}
</div>
```

- `role="status"` = `aria-live="polite"` + `aria-atomic="true"`. Use for non-urgent status.
- `aria-live="assertive"` interrupts the user -- use only for errors that need immediate attention (form validation on submit).
- `aria-live="off"` (default) -- content is not announced.

### Loading state

```tsx
<button disabled={loading}>
	{loading && (
		<span role="status" aria-label="Memuat">
			<Spinner className="animate-spin" />
		</span>
	)}
	{loading ? "Menyimpan..." : "Simpan"}
</button>
```

### Menu (dropdown)

```tsx
<button aria-haspopup="true" aria-expanded={open} aria-controls="menu-1">
  Menu
</button>
<ul id="menu-1" role="menu" hidden={!open}>
  <li role="menuitem"><button>Item 1</button></li>
  <li role="menuitem"><button>Item 2</button></li>
</ul>
```

---

## 6. Announcing Dynamic Content

### When content appears after an action

- Form validation errors after submit → `role="alert"` or `aria-live="assertive"`.
- "Saved successfully" toast → `role="status"`.
- Live data updates (e.g. a counter) → `aria-live="polite"`.
- Loading spinner that appears → `role="status"` on the spinner wrapper.

### Visually-hidden text (for screen reader-only content)

```tsx
function VisuallyHidden({ children }: { children: React.ReactNode }) {
	return <span className="sr-only">{children}</span>;
}
```

Tailwind's `sr-only` class. Use for:

- Additional context on icon buttons: `<button><Trash /><span className="sr-only">Hapus item</span></button>`.
- Form field hints: `<input /> <span className="sr-only">Format: DD/MM/YYYY</span>`.
- Live region announcements.

### Don't over-announce

Every aria-live region fires when its text content changes. If you update a counter every 100ms, the screen reader will talk over itself. Throttle updates to ~1s minimum, or use `aria-live="off"` and toggle it on only when the user should hear the update.

---

## 7. Keyboard & Switch Access

Some users navigate with a keyboard, switch device, or voice control -- none of which use a mouse. Test your page with **Tab**, **Shift+Tab**, **Enter**, **Space**, **Escape**, and arrow keys.

### Required keyboard behaviors

| Element          | Key                    | Behavior           |
| :--------------- | :--------------------- | :----------------- |
| Button           | Enter / Space          | Activate           |
| Link             | Enter                  | Navigate           |
| Dialog           | Escape                 | Close              |
| Tab (in tablist) | ArrowLeft / ArrowRight | Move between tabs  |
| Tab (in tablist) | Tab                    | Move to tabpanel   |
| Menu item        | ArrowUp / ArrowDown    | Move between items |
| Combobox         | ArrowDown              | Open dropdown      |
| Combobox         | Escape                 | Close dropdown     |

### Don't remove outline without replacement

```css
/*  Bad -- removes keyboard focus indicator */
button {
	outline: none;
}

/*  Use focus-visible -- only shows for keyboard, not mouse */
button:focus {
	outline: none;
}
button:focus-visible {
	outline: 2px solid theme(colors.blue.500);
	outline-offset: 2px;
}
```

Tailwind: `focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2`.

### Switch Access (Android)

Android's Switch Access lets users with motor impairments navigate by cycling through actionable items with a single switch. It only stops on elements that are focusable. Ensure every interactive element has `tabindex="0"` (or is a native button/link) -- `div onClick` without `tabindex` is invisible to Switch Access.

---

## Quick Mobile A11y Audit (15-minute check)

1. **VoiceOver**: Turn on. Navigate the page top-to-bottom with swipe-right. Can you understand the page structure? Are icon buttons labeled?
2. **TalkBack**: Same audit on Android.
3. **Keyboard**: Unplug mouse. Can you complete the primary task with Tab/Enter/Escape?
4. **Tap targets**: Open Chrome DevTools → Elements → toggle "Highlighting → Show layout shift regions". Use a ruler -- are all interactive elements ≥ 44×44px?
5. **Color contrast**: Use the WCAG plugin or Chrome DevTools → CSS Overview. 4.5:1 for body text, 3:1 for large text and UI components.
6. **Reduced motion**: macOS System Settings → Accessibility → Display → Reduce Motion. Does the page still work? Are animations gone or shortened?
7. **Dark mode**: System → Dark. Does the page remain readable?
8. **Largest text**: Android Settings → Font size → Largest. Does anything overflow?
