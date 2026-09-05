# Mobile Responsive Validation Checklist

Use this checklist as a final review before merging any UI feature. Each item has a concrete test step.

## How to use

- Run through every item for new pages or major features.
- For small changes (a single button, a text fix), only the affected section is needed.
- The checklist is grouped by concern -- work top to bottom.

---

## A. Viewport & Layout Foundation

### A1. Viewport meta tag

- [ ] `<meta name="viewport">` includes `width=device-width, initial-scale=1.0, viewport-fit=cover, interactive-widget=resizes-content`
- [ ] **Test**: view page source, search for the viewport meta
- [ ] **Why**: missing `viewport-fit=cover` makes `env(safe-area-inset-*)` return 0 on iOS

### A2. No `100vh`

- [ ] Search codebase for `100vh` -- replace with `100dvh` (or `100svh`/`100lvh` per use case)
- [ ] **Test**: `grep -r "100vh" src/`
- [ ] **Why**: `100vh` on iOS Safari includes the URL bar area, content overflows

### A3. No horizontal overflow

- [ ] On a 320px-wide viewport (iPhone SE), drag horizontally -- page doesn't scroll sideways
- [ ] **Test**: DevTools → Toggle device → iPhone SE → swipe left/right
- [ ] **Why**: horizontal scroll is jarring and a sign of layout bugs

### A4. Safe area insets

- [ ] Fixed/sticky top elements have `pt-[env(safe-area-inset-top)]`
- [ ] Fixed/sticky bottom elements have `pb-[max(<default>, env(safe-area-inset-bottom))]`
- [ ] **Test**: on iPhone X+ in PWA standalone mode, content isn't under the notch or home indicator
- [ ] **Why**: without this, content disappears under hardware features

### A5. `overflow-x-hidden` on body/main

- [ ] Body and main wrapper have `overflow-x-hidden`
- [ ] **Test**: view source
- [ ] **Why**: prevents wide children from causing page-wide horizontal scroll

---

## B. Touch Targets & Input

### B1. Touch target size

- [ ] All interactive elements (buttons, links, icon buttons, nav items) are ≥ 44×44 CSS px
- [ ] **Test**: DevTools → Elements → pick a button → computed `min-width` and `min-height` ≥ 44px
- [ ] **Why**: WCAG 2.2 SC 2.5.8 (Level AA: 24px, Level AAA: 44px), Apple HIG

### B2. `touch-manipulation` on interactive elements

- [ ] Buttons and tap targets have `touch-manipulation` (Tailwind class: `touch-manipulation`)
- [ ] **Test**: view source
- [ ] **Why**: removes the legacy 300ms tap delay

### B3. Active state feedback

- [ ] Buttons have `active:scale-[0.97]` or `active:bg-slate-100` for tactile feedback
- [ ] **Test**: tap a button on mobile -- does it visibly change?
- [ ] **Why**: mobile users have no cursor; without feedback, taps feel unregistered

### B4. iOS tap highlight removed (with replacement)

- [ ] Global CSS has `-webkit-tap-highlight-color: transparent` on interactive elements
- [ ] **Test**: tap a button on iOS Safari -- no gray box appears
- [ ] **Why**: gray highlight looks dated once you have `:active` feedback

### B5. iOS auto-zoom prevention

- [ ] All `<input>`, `<select>`, `<textarea>` have `text-base` (16px) on mobile
- [ ] **Test**: tap an input on iOS Safari -- page doesn't zoom in
- [ ] **Why**: iOS Safari zooms on focus if font-size < 16px

### B6. `inputmode` + `autocomplete` + `enterkeyhint`

- [ ] Every input has appropriate `type`, `inputMode`, `autoComplete`, `enterKeyHint`
- [ ] **Test**: tap a phone field -- numeric keypad appears with tel layout, "Next" key
- [ ] **Why**: native-feel forms; without these, users get the wrong keyboard and no autofill

### B7. Form input keyboard types

| Field         | type     | inputMode | autoComplete   | enterKeyHint |
| :------------ | :------- | :-------- | :------------- | :----------- |
| NIK           | text     | numeric   | off            | next         |
| Phone         | tel      | tel       | tel            | next         |
| Email         | email    | email     | email          | next         |
| Date of birth | date     | (none)    | bday           | next         |
| Address       | text     | (none)    | street-address | next         |
| Postal code   | text     | numeric   | postal-code    | next         |
| Search        | search   | search    | off            | search       |
| Last field    | (varies) | (varies)  | (varies)       | go / done    |

---

## C. Astro Hydration

### C1. Minimal hydration directive

- [ ] Below-the-fold islands use `client:visible` or `client:idle`, not `client:load`
- [ ] Above-the-fold critical UI uses `client:load` (justified)
- [ ] Browser-only components use `client:only="react"`
- [ ] **Test**: view page source -- count `client:` directives
- [ ] **Why**: every `client:load` downloads JS that runs immediately, slowing first interaction

### C2. No SSR-unsafe browser API access

- [ ] No `window`, `document`, `localStorage`, `navigator` access in component body
- [ ] All such access is inside `useEffect` or guarded with `typeof window !== "undefined"`
- [ ] **Test**: `grep -E "window\.|document\.|localStorage" src/components/`
- [ ] **Why**: causes hydration mismatch errors when SSR HTML differs from client render

### C3. View Transitions (if used)

- [ ] `ClientRouter` imported in layout
- [ ] Stateful islands wrapped with `transition:persist`
- [ ] Reduced-motion CSS disables view transition animations
- [ ] **Test**: navigate between pages -- does the transition work? Does it disappear when "Reduce Motion" is on?
- [ ] **Why**: view transitions are nice but can trigger motion sickness

---

## D. Navigation & Drawers

### D1. Drawer pattern

- [ ] Drawer panel is `fixed` with `transform` translate
- [ ] `max-w-[85vw]` so backdrop is tappable
- [ ] Backdrop has `bg-slate-900/60 backdrop-blur-sm`
- [ ] Backdrop tap closes the drawer
- [ ] ESC key closes the drawer
- [ ] Body scroll locked while open (`document.body.style.overflow = "hidden"`)
- [ ] Focus trapped inside the drawer
- [ ] Focus returns to the trigger on close
- [ ] **Test**: open drawer → Tab through → close → focus returns to trigger button

### D2. Mobile bottom navigation (if used)

- [ ] 3-5 items maximum
- [ ] Each item ≥ 44px tall (recommend `min-h-14` = 56px for icon + label)
- [ ] `pb-[env(safe-area-inset-bottom)]` on the nav container
- [ ] `md:hidden` (mobile-only)
- [ ] Page `<main>` has `pb-20` (or equivalent) to clear the nav
- [ ] **Test**: navigate using bottom nav on a phone -- does the active state update?

### D3. Back button handling (Android)

- [ ] System back button closes open modals/drawers (not navigate away)
- [ ] Pattern: `history.pushState` + `popstate` listener
- [ ] **Test**: open drawer on Android → press back → drawer closes (not page navigation)

---

## E. Modals & Bottom Sheets

### E1. Responsive modal pattern

- [ ] Mobile: bottom sheet (`items-end`, `rounded-t-2xl`)
- [ ] Desktop: centered dialog (`sm:items-center`, `sm:rounded-2xl`)
- [ ] `max-h-[90dvh]` (NOT `90vh`)
- [ ] Drag handle visible on mobile (`sm:hidden`)
- [ ] Body has `overscroll-contain` to prevent scroll chaining
- [ ] Footer has `flex-col-reverse sm:flex-row` (primary on bottom on mobile for thumb reach)
- [ ] Footer has `pb-[max(1rem, env(safe-area-inset-bottom))]`
- [ ] **Test**: open modal on iPhone → swipe down on body → background doesn't scroll

### E2. Focus management

- [ ] `role="dialog"` and `aria-modal="true"` on the modal container
- [ ] `aria-label` or `aria-labelledby` set
- [ ] Focus moves into the modal on open
- [ ] Focus trap cycles Tab/Shift+Tab within the modal
- [ ] Focus returns to trigger on close
- [ ] ESC closes the modal
- [ ] **Test**: open modal with keyboard → Tab cycles within → ESC closes → focus is back on trigger

### E3. Drag-to-dismiss (if implemented)

- [ ] Drag handle is the only drag initiator (not the body)
- [ ] Only downward drag moves the sheet (upward scrolls body)
- [ ] 100px threshold for dismiss
- [ ] Transition disabled during drag, enabled on release
- [ ] **Test**: drag down slowly → sheet follows → release at 50px (snaps back) → release at 120px (dismisses)

---

## F. Tables & Data Display

### F1. Card fallback on mobile

- [ ] Tables with > 4 columns have a card-based mobile view
- [ ] Mobile view: `md:hidden`
- [ ] Desktop table: `hidden md:block`
- [ ] Both views render from the same data source
- [ ] **Test**: resize browser from desktop to mobile -- view switches cleanly

### F2. Horizontal scroll container (if table must stay)

- [ ] Container has `overflow-x-auto`
- [ ] Table has `min-w-[<appropriate>]px` so columns don't squish
- [ ] Sticky header with `bg-white` and `z-10`
- [ ] `-mx-4 px-4` to bleed to edge on mobile
- [ ] Scroll indicator (gradient/shadow) on the right edge
- [ ] **Test**: scroll table horizontally on mobile -- header stays visible

### F3. Long content wrapping

- [ ] NIK / NIP / NPWP fields use `font-mono` + `break-all`
- [ ] Long URLs/email use `break-all` or `break-words`
- [ ] **Test**: view a list with long identifiers -- no horizontal overflow

---

## G. Forms

### G1. iOS auto-zoom (re-checked)

- [ ] All inputs `text-base` on mobile (16px)
- [ ] **Test**: focus an input on iOS Safari -- page doesn't zoom

### G2. Sticky submit button (long forms)

- [ ] Submit button `sticky bottom-0` for forms longer than one screen
- [ ] Background covers content scrolling under it
- [ ] `pb-[max(0.75rem, env(safe-area-inset-bottom))]`
- [ ] **Test**: scroll a long form -- submit button stays visible

### G3. Field grouping

- [ ] Related fields grouped (e.g. RT and RW side by side)
- [ ] Labels above inputs (not left-aligned) on mobile
- [ ] **Test**: scan a form top-to-bottom on mobile -- is the structure clear?

### G4. Validation

- [ ] Inline validation on blur (not on every keystroke)
- [ ] Errors have `role="alert"` and `aria-describedby` linking to the field
- [ ] Submit-time errors scroll the first error into view and focus it
- [ ] `aria-invalid="true"` on fields with errors
- [ ] **Test**: submit empty form → errors appear → first error is focused

### G5. WCAG 3.3.7 Redundant Entry

- [ ] Multi-step forms don't ask for the same info twice without autofill
- [ ] **Test**: complete step 1 with email → step 3 needs email → it's pre-filled

---

## H. Images & Media

### H1. Responsive images

- [ ] Astro `<Image>` component used for local images
- [ ] `width` and `height` (or `aspect-ratio`) set to prevent CLS
- [ ] `sizes` attribute set on every responsive image
- [ ] `loading="lazy"` on below-the-fold images
- [ ] `loading="eager"` on the LCP image (above the fold)
- [ ] `decoding="async"` on all images
- [ ] **Test**: DevTools → Performance → reload → no layout shifts

### H2. Alt text

- [ ] Content images have descriptive `alt`
- [ ] Decorative images have `alt=""` (empty, not omitted)
- [ ] Icon-only buttons have `aria-label`
- [ ] **Test**: browse with images disabled -- can you still understand the page?

### H3. Video / audio

- [ ] Autoplay disabled (or muted + playsinline)
- [ ] `playsinline` attribute on `<video>` (prevents iOS forced fullscreen)
- [ ] User-gesture-initiated playback
- [ ] Captions / transcript provided
- [ ] **Test**: open a video page on iOS Safari → it doesn't autoplay fullscreen

---

## I. Motion & User Preferences

### I1. Reduced motion

- [ ] Global CSS has `@media (prefers-reduced-motion: reduce)` block
- [ ] All animations/transitions disabled or shortened in that block
- [ ] JS-driven animations (Framer Motion, GSAP) check the media query
- [ ] View Transitions animations disabled in reduced motion
- [ ] **Test**: macOS Settings → Accessibility → Display → Reduce Motion → reload page

### I2. Dark mode

- [ ] `theme-color` meta for both light and dark
- [ ] `color-scheme: light dark` meta
- [ ] Tailwind `dark:` variants applied to all custom colors
- [ ] Images/logos have dark-mode variants if needed
- [ ] **Test**: System → Dark mode → page is readable, no white flashes

### I3. Contrast

- [ ] Body text: ≥ 4.5:1 contrast ratio
- [ ] Large text (≥ 18pt): ≥ 3:1
- [ ] UI components (borders, icons): ≥ 3:1
- [ ] **Test**: Chrome DevTools → CSS Overview → Colors → contrast issues

### I4. Other preference queries

- [ ] `prefers-contrast: more` increases border weights (optional)
- [ ] `forced-colors: active` ensures borders are explicit (Windows High Contrast)
- [ ] `inverted-colors: inverted` doesn't break (iOS Smart Invert)
- [ ] **Test**: Windows → High Contrast Mode → page is usable

---

## J. Accessibility

### J1. Semantic structure

- [ ] Page has exactly one `<h1>`
- [ ] Heading hierarchy is sequential (h1 → h2 → h3, no skips)
- [ ] Landmarks: `<header>`, `<nav>`, `<main>`, `<footer>`, `<aside>` present
- [ ] Skip-to-content link is the first focusable element

### J2. Keyboard navigation

- [ ] All interactive elements reachable via Tab
- [ ] Tab order matches visual order
- [ ] Visible focus indicators on all interactive elements
- [ ] No `outline: none` without a replacement (`:focus-visible:ring-*`)
- [ ] **Test**: unplug mouse → navigate page with Tab/Shift+Tab/Enter/Escape

### J3. Screen reader basics

- [ ] Icon-only buttons have `aria-label`
- [ ] Form inputs have associated `<label>` (or `aria-label`)
- [ ] Decorative images have `alt=""`
- [ ] Dynamic content uses `aria-live` regions
- [ ] Modals have `role="dialog"` and `aria-modal="true"`
- [ ] **Test**: iOS VoiceOver → swipe right through page → structure makes sense

### J4. WCAG 2.2 mobile criteria

- [ ] 2.5.7 Dragging Movements: drag-only interactions have tap alternatives
- [ ] 2.5.8 Target Size: ≥ 24×24 CSS px (AA), ≥ 44×44 (AAA)
- [ ] 3.3.7 Redundant Entry: previously entered data is autofilled
- [ ] 3.3.8 Accessible Authentication: no cognitive-only tests (CAPTCHA) for login
- [ ] 4.1.3 Status Messages: toasts/errors use `role="status"` or `aria-live`

---

## K. Performance (mobile-specific)

### K1. Bundle size

- [ ] No `client:load` on below-the-fold islands
- [ ] Third-party JS lazy-loaded (analytics, chat widgets)
- [ ] No render-blocking external scripts
- [ ] **Test**: Chrome DevTools → Network → JS → total < 200KB gzipped on first load

### K2. Image weight

- [ ] AVIF or WebP format
- [ ] Appropriate sizes per breakpoint (don't ship a 1920px image to a 360px screen)
- [ ] **Test**: DevTools → Network → Img → total < 300KB on first load

### K3. Font loading

- [ ] `font-display: swap` (or `optional`) to prevent invisible text
- [ ] Preload critical font files
- [ ] Subset to needed character sets (e.g. Latin + Latin Extended, not full CJK if not needed)
- [ ] **Test**: reload page on slow 3G → no invisible text flashes

### K4. CLS (Cumulative Layout Shift)

- [ ] All images have `width`/`height`/`aspect-ratio`
- [ ] No injected ads/banners that push content
- [ ] Fonts loaded with `font-display: swap` + size-adjust to minimize reflow
- [ ] **Test**: DevTools → Performance Insights → CLS < 0.1

### K5. LCP (Largest Contentful Paint)

- [ ] LCP element is text or has eager-loaded image
- [ ] No render-blocking resources above the fold
- [ ] **Test**: DevTools → Performance Insights → LCP < 2.5s on slow 3G

---

## L. Cross-Platform Test Matrix

### L1. iOS Safari

- [ ] iPhone SE (375×667) -- smallest common iPhone
- [ ] iPhone 15 (393×852) -- standard
- [ ] iPhone 15 Pro Max (430×932) -- large
- [ ] Landscape orientation (test safe-area side insets)
- [ ] With VoiceOver enabled
- [ ] In PWA standalone mode (Add to Home Screen)
- [ ] With iOS Dynamic Type set to largest

### L2. Android Chrome

- [ ] Galaxy A series (360×800) -- small / low-end
- [ ] Pixel 7 (412×915) -- standard
- [ ] Galaxy S23 Ultra (412×915 @ 2x) -- large
- [ ] Landscape orientation
- [ ] With TalkBack enabled
- [ ] With Android font size set to Largest
- [ ] With "Remove animations" enabled (Developer Options)

### L3. Special conditions

- [ ] Slow 3G network (DevTools → Network → Slow 3G)
- [ ] Offline mode (after first load -- does the page show cached content?)
- [ ] High Contrast Mode (Windows)
- [ ] Browser zoom at 200% (Ctrl/Cmd + =)
- [ ] Pinch-zoom to 200% on mobile (don't disable this!)

### L4. PWA / Hybrid

- [ ] iOS PWA standalone (Add to Home Screen) -- status bar, safe areas
- [ ] Android PWA standalone -- status bar, navigation
- [ ] Capacitor/Tauri WebView (if applicable) -- bridge, permissions

---

## M. Final Sanity Check

- [ ] Show the page to a non-developer on their own phone -- watch them use it without instruction
- [ ] Ask: "What is this page for?" -- they should be able to answer in < 5 seconds
- [ ] Ask them to complete the primary task -- note where they hesitate or fail
- [ ] Compare load time and feel against a top-tier mobile site (Twitter, Notion, Linear) -- your page should feel similarly responsive

If the page fails any of these, it's not ready to ship. Fix before merge.
