# Mobile Component Patterns -- Deep Dive

Open this file when implementing drawers, bottom sheets, bottom navigation bars, FABs, swipeable cards, search interfaces, sticky headers, or any complex mobile UI pattern.

## Table of Contents

1. [Drawer / Sidebar (slide-in)](#1-drawer--sidebar-slide-in)
2. [Bottom Sheet Modal](#2-bottom-sheet-modal)
3. [Bottom Navigation Bar](#3-bottom-navigation-bar)
4. [Floating Action Button (FAB)](#4-floating-action-button-fab)
5. [Sticky Header with Search](#5-sticky-header-with-search)
6. [Swipeable List Items](#6-swipeable-list-items)
7. [Pull-to-Refresh](#7-pull-to-refresh)
8. [Search Experience](#8-search-experience)
9. [Skeleton Loading States](#9-skeleton-loading-states)
10. [Infinite Scroll](#10-infinite-scroll)

---

## 1. Drawer / Sidebar (slide-in)

Expands on the pattern in SKILL.md §4. Use for navigation menus, filter panels, and detail panels.

### Full implementation with focus trap

```tsx
import { useEffect, useRef } from "react";

interface DrawerProps {
	open: boolean;
	onClose: () => void;
	side?: "left" | "right";
	title?: string;
	children: React.ReactNode;
}

export function Drawer({ open, onClose, side = "left", title, children }: DrawerProps) {
	const panelRef = useRef<HTMLDivElement>(null);
	const triggerRef = useRef<Element | null>(null);

	// Body scroll lock
	useEffect(() => {
		if (!open) return;
		const original = document.body.style.overflow;
		document.body.style.overflow = "hidden";
		return () => {
			document.body.style.overflow = original;
		};
	}, [open]);

	// Focus trap + ESC
	useEffect(() => {
		if (!open) return;
		triggerRef.current = document.activeElement;

		const panel = panelRef.current;
		if (!panel) return;

		const focusable = panel.querySelectorAll<HTMLElement>(
			'a[href], button:not([disabled]), textarea, input, select, [tabindex]:not([tabindex="-1"])'
		);
		focusable[0]?.focus();

		const onKey = (e: KeyboardEvent) => {
			if (e.key === "Escape") {
				onClose();
				return;
			}
			if (e.key !== "Tab" || focusable.length === 0) return;
			const first = focusable[0];
			const last = focusable[focusable.length - 1];
			if (e.shiftKey && document.activeElement === first) {
				e.preventDefault();
				last.focus();
			} else if (!e.shiftKey && document.activeElement === last) {
				e.preventDefault();
				first.focus();
			}
		};

		document.addEventListener("keydown", onKey);
		return () => {
			document.removeEventListener("keydown", onKey);
			(triggerRef.current as HTMLElement)?.focus();
		};
	}, [open, onClose]);

	const sideClasses =
		side === "left"
			? "left-0 " + (open ? "translate-x-0" : "-translate-x-full")
			: "right-0 " + (open ? "translate-x-0" : "translate-x-full");

	return (
		<div aria-hidden={!open}>
			{/* Backdrop */}
			<div
				className={`fixed inset-0 z-40 bg-slate-900/60 backdrop-blur-sm transition-opacity duration-300
          ${open ? "opacity-100" : "opacity-0 pointer-events-none"}`}
				onClick={onClose}
				aria-hidden="true"
			/>

			{/* Panel */}
			<aside
				ref={panelRef}
				role="dialog"
				aria-modal="true"
				aria-label={title}
				className={`fixed inset-y-0 ${sideClasses} z-50 w-80 max-w-[85vw] bg-white shadow-2xl
          transform transition-transform duration-300 ease-out flex flex-col
          md:hidden`}
			>
				{/* Status bar space */}
				<div className="pt-[env(safe-area-inset-top)]">
					{title && (
						<header className="flex items-center justify-between px-4 py-3 border-b">
							<h2 className="text-base font-semibold">{title}</h2>
							<button
								onClick={onClose}
								aria-label="Tutup"
								className="min-h-11 min-w-11 flex items-center justify-center rounded-lg active:bg-slate-100"
							>
								<XIcon className="w-5 h-5" />
							</button>
						</header>
					)}
					<div className="flex-1 overflow-y-auto overscroll-contain">{children}</div>
				</div>
			</aside>
		</div>
	);
}
```

### Critical details

1. **`max-w-[85vw]`** -- leaves a sliver of backdrop for tap-to-close.
2. **Body scroll lock + cleanup** -- restores original overflow, doesn't blindly set `""`.
3. **Focus trap returns focus to the trigger** on close -- critical for keyboard users.
4. **`overscroll-contain`** on the scrollable body -- prevents scroll chaining to the page behind.
5. **`md:hidden`** -- drawer only shows on mobile. On desktop, use a static sidebar (separate component or `md:static md:translate-x-0`).
6. **`transition-transform` only** -- transitioning `opacity` on the panel causes a flash on iOS Safari. Transform-only is GPU-accelerated and smooth.

---

## 2. Bottom Sheet Modal

The mobile-preferred modal pattern. Snaps to the bottom, can be dragged down to dismiss.

### Basic version (no drag)

See SKILL.md §7.1 for the basic version.

### With drag-to-dismiss

```tsx
import { useRef, useState, useEffect } from "react";

export function BottomSheet({ open, onClose, children, title }: Props) {
	const [dragY, setDragY] = useState(0);
	const startY = useRef<number | null>(null);

	const onTouchStart = (e: React.TouchEvent) => {
		startY.current = e.touches[0].clientY;
	};

	const onTouchMove = (e: React.TouchEvent) => {
		if (startY.current === null) return;
		const delta = e.touches[0].clientY - startY.current;
		if (delta > 0) setDragY(delta); // only allow downward drag
	};

	const onTouchEnd = () => {
		if (dragY > 100) {
			onClose();
		}
		setDragY(0);
		startY.current = null;
	};

	if (!open) return null;

	return (
		<div className="fixed inset-0 z-50 flex items-end justify-center">
			<div className="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" onClick={onClose} />
			<div
				role="dialog"
				aria-modal="true"
				aria-label={title}
				className="relative w-full bg-white rounded-t-2xl shadow-2xl max-h-[90dvh] flex flex-col"
				style={{
					transform: `translateY(${dragY}px)`,
					transition: dragY === 0 ? "transform 0.3s ease-out" : "none",
				}}
			>
				{/* Drag handle -- also the drag initiator */}
				<div
					onTouchStart={onTouchStart}
					onTouchMove={onTouchMove}
					onTouchEnd={onTouchEnd}
					className="flex justify-center pt-3 pb-2 cursor-grab active:cursor-grabbing"
				>
					<div className="h-1.5 w-10 rounded-full bg-slate-300" />
				</div>

				<header className="px-6 py-3 border-b">
					<h2 className="text-lg font-semibold">{title}</h2>
				</header>

				<div className="flex-1 overflow-y-auto overscroll-contain px-6 py-4">{children}</div>
			</div>
		</div>
	);
}
```

### Implementation notes

- **Only allow downward drag** (delta > 0) -- upward drag should scroll content, not move the sheet.
- **100px threshold** for dismiss -- too small and accidental drags dismiss; too large and it feels unresponsive.
- **`transition: none` during drag** -- applying a transition while dragging causes lag. Re-enable the transition on release so the sheet animates back smoothly.
- **Don't allow drag from the body** -- only from the handle. Otherwise users can't scroll the body content (every scroll attempt would drag the sheet).
- **Close on backdrop tap** AND on ESC AND on drag-down -- multiple paths to close is a feature, not a bug.

### Snap points (optional)

For sheets with multiple heights (half-open, full-open):

```tsx
const [snapIndex, setSnapIndex] = useState(0);
const snapPoints = [0.5, 0.9]; // 50%, 90% of viewport

<div
	style={{
		height: `${snapPoints[snapIndex] * 100}dvh`,
		transition: "height 0.3s ease-out",
	}}
>
	{/* content with snap controls */}
</div>;
```

Use sparingly -- snap points add complexity. Most sheets should be either full-height or content-height.

---

## 3. Bottom Navigation Bar

The mobile-app pattern: 3-5 destinations pinned to the bottom. Native to iOS/Android, increasingly common in PWAs.

```tsx
import { NavLink } from "react-router-dom"; // or your router

const items = [
	{ to: "/", label: "Beranda", icon: HomeIcon },
	{ to: "/search", label: "Cari", icon: SearchIcon },
	{ to: "/notifications", label: "Notifikasi", icon: BellIcon, badge: 3 },
	{ to: "/profile", label: "Profil", icon: UserIcon },
];

export function BottomNav() {
	return (
		<nav
			aria-label="Navigasi utama"
			className="fixed bottom-0 inset-x-0 z-30 bg-white border-t border-slate-200
                 md:hidden
                 pb-[env(safe-area-inset-bottom)]"
		>
			<ul className="flex">
				{items.map((item) => (
					<li key={item.to} className="flex-1">
						<NavLink
							to={item.to}
							className={({ isActive }) =>
								`flex flex-col items-center justify-center gap-0.5 py-2 px-1 min-h-14
                 ${isActive ? "text-blue-600" : "text-slate-600"}`
							}
						>
							<span className="relative">
								<item.icon className="w-6 h-6" />
								{item.badge && (
									<span
										className="absolute -top-1 -right-2 min-w-4 h-4 px-1
                               bg-red-500 text-white text-[10px] rounded-full
                               flex items-center justify-center"
										aria-label={`${item.badge} notifikasi baru`}
									>
										{item.badge}
									</span>
								)}
							</span>
							<span className="text-xs">{item.label}</span>
						</NavLink>
					</li>
				))}
			</ul>
		</nav>
	);
}
```

### Critical details

1. **3-5 items max** -- beyond 5, labels truncate and tap targets shrink. Use a "More" menu for the overflow.
2. **`min-h-14` (56px)** -- taller than the 44px minimum because the bar contains both icon and label.
3. **`pb-[env(safe-area-inset-bottom)]`** -- extends the bar's background below the home indicator.
4. **`md:hidden`** -- bottom nav is mobile-only. On tablet/desktop, use a sidebar.
5. **Badges with `aria-label`** -- don't just show "3"; tell screen readers what it means.
6. **Active state** uses color, not just bold -- color-blind users need multiple cues.

### Page body padding

Pages using `<BottomNav>` must add bottom padding so content isn't hidden behind it:

```tsx
<main className="pb-20 md:pb-0">{/* page content */}</main>
```

`pb-20` = 80px (56px nav + safe area + breathing room). `md:pb-0` because the nav is hidden on desktop.

---

## 4. Floating Action Button (FAB)

A circular button pinned to the bottom-right, above the bottom nav. Primary action only.

```tsx
export function FAB({ onClick, icon: Icon, label }: Props) {
	return (
		<button
			onClick={onClick}
			aria-label={label}
			className="fixed bottom-4 right-4 z-30
                 mb-[env(safe-area-inset-bottom)]
                 w-14 h-14 rounded-full bg-blue-600 text-white shadow-lg
                 flex items-center justify-center
                 active:scale-95 active:bg-blue-700
                 transition-transform
                 md:hidden"
		>
			<Icon className="w-6 h-6" />
		</button>
	);
}
```

### Variants

- **Regular FAB** (56px / `w-14 h-14`) -- most common.
- **Mini FAB** (40px / `w-10 h-10`) -- when below a bottom nav that already has a primary action.
- **Extended FAB** (pill with icon + label) -- when the action isn't obvious from the icon alone:
  ```tsx
  <button className="fixed bottom-4 right-4 ... h-14 px-6 rounded-full ... gap-2">
  	<PlusIcon /> <span>Tambah</span>
  </button>
  ```

### Position relative to bottom nav

If you have both a bottom nav (56px) and a FAB:

```tsx
// FAB sits ABOVE the nav
<button className="fixed bottom-20 right-4 ... md:hidden">
```

`bottom-20` = 80px = 56px nav + 16px gap + safe area.

### Don't use a FAB when

- The primary action is already obvious (e.g. a "Send" button in a chat input bar).
- There are multiple primary actions (use an "extended FAB" that expands into a stack on tap, or just put them in the page).
- The page has nothing to create -- FABs imply "add new thing".

---

## 5. Sticky Header with Search

Admin pages typically have a title + optional search field + filter button at the top. Make it sticky so users can search from anywhere on the page.

```tsx
export function PageHeader({ title, onSearch }: Props) {
	return (
		<header
			className="sticky top-0 z-20 bg-white/95 backdrop-blur-sm border-b
                 pt-[env(safe-area-inset-top)]"
		>
			<div className="px-4 py-3">
				<h1 className="text-xl font-semibold mb-3">{title}</h1>

				<div className="flex gap-2">
					<label className="relative flex-1">
						<span className="sr-only">Cari</span>
						<SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
						<input
							type="search"
							inputMode="search"
							enterKeyHint="search"
							placeholder="Cari..."
							onChange={(e) => onSearch(e.target.value)}
							className="w-full pl-10 pr-3 py-2.5 text-base rounded-lg border border-slate-300
                         focus:outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</label>
					<button
						aria-label="Filter"
						className="min-h-11 min-w-11 px-3 flex items-center justify-center rounded-lg border border-slate-300"
					>
						<FilterIcon className="w-5 h-5" />
					</button>
				</div>
			</div>
		</header>
	);
}
```

### Critical details

1. **`pt-[env(safe-area-inset-top)]`** -- in PWA standalone mode, the header extends under the status bar.
2. **`backdrop-blur-sm` + `bg-white/95`** -- content scrolling underneath is visible but not distracting.
3. **`text-base` on input** -- prevents iOS zoom.
4. **`enterKeyHint="search"`** -- keyboard shows "Search" instead of "Return".
5. **Sticky header adds to the scroll container's effective height** -- make sure the page's main scroll container accounts for it. If you have a sticky footer too, your scroll area is `100dvh - header - footer`.

---

## 6. Swipeable List Items

Swipe-left to reveal "Delete", swipe-right to reveal "Archive" -- a native-feel pattern.

### Implementation with Pointer Events

```tsx
function SwipeableItem({ item, onDelete, onArchive, children }: Props) {
	const [offset, setOffset] = useState(0);
	const startX = useRef<number | null>(null);

	const onPointerDown = (e: React.PointerEvent) => {
		startX.current = e.clientX;
		(e.target as HTMLElement).setPointerCapture(e.pointerId);
	};

	const onPointerMove = (e: React.PointerEvent) => {
		if (startX.current === null) return;
		const delta = e.clientX - startX.current;
		// Resist swipe beyond limits
		setOffset(Math.max(-120, Math.min(120, delta)));
	};

	const onPointerUp = () => {
		if (offset < -60) {
			// Swiped left → reveal delete
			setOffset(-80);
		} else if (offset > 60) {
			// Swiped right → reveal archive
			setOffset(80);
		} else {
			setOffset(0);
		}
		startX.current = null;
	};

	return (
		<div className="relative overflow-hidden">
			{/* Action backgrounds */}
			<div className="absolute inset-0 flex">
				<div className="bg-green-500 flex items-center justify-start pl-4 flex-1">
					<ArchiveIcon className="w-5 h-5 text-white" />
				</div>
				<div className="bg-red-500 flex items-center justify-end pr-4 flex-1">
					<TrashIcon className="w-5 h-5 text-white" />
				</div>
			</div>

			{/* Foreground item */}
			<div
				onPointerDown={onPointerDown}
				onPointerMove={onPointerMove}
				onPointerUp={onPointerUp}
				onPointerCancel={onPointerUp}
				style={{
					transform: `translateX(${offset}px)`,
					transition: startX.current === null ? "transform 0.2s ease-out" : "none",
				}}
				className="relative bg-white"
			>
				{children}
			</div>
		</div>
	);
}
```

### Critical details

1. **Pointer Events, not Touch Events** -- Pointer Events work for mouse, touch, and pen uniformly.
2. **`setPointerCapture`** -- keeps receiving move events even if the finger leaves the element.
3. **Resist swipe beyond limits** -- `Math.max(-120, Math.min(120, delta))` adds a "rubber band" feel.
4. **60px threshold** for action reveal -- feels responsive but not accidental.
5. **`transition: none` during drag** -- see bottom sheet, same reason.
6. **Accessibility** -- swipeable items are hard to use with a keyboard. Provide a context menu button as well:
   ```tsx
   <button aria-label="Opsi lainnya" onClick={openMenu}>
   	<MoreIcon />
   </button>
   ```
7. **Mobile-only** -- on desktop, hide the swipe affordance and use a hover-revealed action bar instead.

---

## 7. Pull-to-Refresh

A native-feel pattern for refreshing list data.

### Implementation

```tsx
function usePullToRefresh(onRefresh: () => Promise<void>) {
	const [pullDistance, setPullDistance] = useState(0);
	const [refreshing, setRefreshing] = useState(false);
	const startY = useRef<number | null>(null);

	useEffect(() => {
		const onTouchStart = (e: TouchEvent) => {
			// Only when scrolled to top
			if (window.scrollY > 0) return;
			startY.current = e.touches[0].clientY;
		};

		const onTouchMove = (e: TouchEvent) => {
			if (startY.current === null) return;
			const delta = e.touches[0].clientY - startY.current;
			if (delta > 0) {
				// Resistance: pulls get harder beyond 80px
				setPullDistance(Math.min(80, delta * 0.5));
			}
		};

		const onTouchEnd = async () => {
			if (pullDistance > 60 && !refreshing) {
				setRefreshing(true);
				try {
					await onRefresh();
				} finally {
					setRefreshing(false);
				}
			}
			setPullDistance(0);
			startY.current = null;
		};

		window.addEventListener("touchstart", onTouchStart, { passive: true });
		window.addEventListener("touchmove", onTouchMove, { passive: true });
		window.addEventListener("touchend", onTouchEnd);

		return () => {
			window.removeEventListener("touchstart", onTouchStart);
			window.removeEventListener("touchmove", onTouchMove);
			window.removeEventListener("touchend", onTouchEnd);
		};
	}, [pullDistance, refreshing, onRefresh]);

	return { pullDistance, refreshing };
}
```

### Critical details

1. **Only trigger when `scrollY === 0`** -- otherwise the gesture conflicts with normal scrolling.
2. **Apply resistance** (`delta * 0.5`) -- without resistance, a 200px drag would translate to a 200px UI shift, which feels wrong.
3. **60px threshold** -- anything less is too easy to trigger accidentally.
4. **Cap at 80px** -- the spinner area shouldn't grow infinitely.
5. **`{ passive: true }`** on touch listeners -- non-passive listeners block scrolling and hurt performance.
6. **Don't compete with `overscroll-behavior`** -- if you set `overscroll-behavior-y: contain` on the scroll container, the native pull-to-refresh is disabled, which is what you want when implementing your own.

### When NOT to implement pull-to-refresh

- The page doesn't have a list (no point refreshing a static article).
- The data refreshes automatically (websockets, polling) -- users will be confused by the manual trigger.
- The list is paginated and a refresh resets to page 1 -- destructive.

---

## 8. Search Experience

### Search input with results dropdown

```tsx
function SearchWithResults() {
	const [query, setQuery] = useState("");
	const [results, setResults] = useState<Result[]>([]);
	const [open, setOpen] = useState(false);
	const [highlightIndex, setHighlightIndex] = useState(-1);

	// Debounced search
	useEffect(() => {
		if (!query.trim()) {
			setResults([]);
			return;
		}
		const timer = setTimeout(async () => {
			const data = await searchAPI(query);
			setResults(data);
			setOpen(true);
		}, 300);
		return () => clearTimeout(timer);
	}, [query]);

	const onKeyDown = (e: React.KeyboardEvent) => {
		if (e.key === "ArrowDown") {
			e.preventDefault();
			setHighlightIndex((i) => Math.min(results.length - 1, i + 1));
		} else if (e.key === "ArrowUp") {
			e.preventDefault();
			setHighlightIndex((i) => Math.max(-1, i - 1));
		} else if (e.key === "Enter" && highlightIndex >= 0) {
			e.preventDefault();
			selectResult(results[highlightIndex]);
		} else if (e.key === "Escape") {
			setOpen(false);
		}
	};

	return (
		<div className="relative">
			<input
				type="search"
				inputMode="search"
				enterKeyHint="search"
				role="combobox"
				aria-expanded={open}
				aria-controls="search-results"
				aria-activedescendant={highlightIndex >= 0 ? `result-${highlightIndex}` : undefined}
				value={query}
				onChange={(e) => {
					setQuery(e.target.value);
					setHighlightIndex(-1);
				}}
				onKeyDown={onKeyDown}
				onFocus={() => results.length && setOpen(true)}
				onBlur={() => setTimeout(() => setOpen(false), 150)}
				className="w-full px-4 py-2.5 text-base rounded-lg border"
				placeholder="Cari penduduk..."
			/>

			{open && results.length > 0 && (
				<ul
					id="search-results"
					role="listbox"
					className="absolute top-full inset-x-0 mt-1 bg-white border rounded-lg shadow-lg max-h-72 overflow-y-auto"
				>
					{results.map((result, i) => (
						<li
							key={result.id}
							id={`result-${i}`}
							role="option"
							aria-selected={i === highlightIndex}
							onMouseDown={() => selectResult(result)} // mousedown, not click -- fires before blur
							className={`px-4 py-3 cursor-pointer ${i === highlightIndex ? "bg-blue-50" : ""}`}
						>
							{result.nama}
						</li>
					))}
				</ul>
			)}
		</div>
	);
}
```

### Critical details

1. **Debounce 300ms** -- searches on every keystroke hammer the server.
2. **`onMouseDown` not `onClick`** -- click fires after blur, which closes the dropdown before the click registers.
3. **`role="combobox"` + `role="listbox"` + `role="option"`** -- full ARIA combobox pattern for screen readers.
4. **`aria-activedescendant`** -- tells screen readers which option is highlighted as the user arrows through.
5. **Keyboard nav** -- ArrowUp/Down to highlight, Enter to select, Escape to close.
6. **`enterKeyHint="search"`** -- submit key says "Search" on mobile keyboard.

---

## 9. Skeleton Loading States

Don't show a blank page or spinner while data loads -- show a skeleton that mimics the layout.

```tsx
function PendudukCardSkeleton() {
	return (
		<li className="rounded-xl border border-slate-200 bg-white p-4 animate-pulse">
			<div className="flex justify-between mb-3">
				<div className="h-5 w-32 bg-slate-200 rounded" />
				<div className="h-5 w-16 bg-slate-200 rounded-full" />
			</div>
			<div className="space-y-2">
				<div className="h-4 w-full bg-slate-200 rounded" />
				<div className="h-4 w-2/3 bg-slate-200 rounded" />
			</div>
			<div className="mt-4 flex gap-2">
				<div className="h-9 flex-1 bg-slate-200 rounded" />
				<div className="h-9 flex-1 bg-slate-200 rounded" />
			</div>
		</li>
	);
}

function PendudukList({ loading, items }: Props) {
	if (loading) {
		return (
			<ul className="grid grid-cols-1 gap-3 md:hidden">
				{Array.from({ length: 5 }).map((_, i) => (
					<PendudukCardSkeleton key={i} />
				))}
			</ul>
		);
	}
	// render actual list
}
```

### Critical details

1. **`animate-pulse`** -- Tailwind's built-in pulse animation. Smooth and consistent.
2. **Match the layout exactly** -- same heights, widths, gaps as the real content. Otherwise the page jumps when data loads (layout shift).
3. **5 items** for list skeletons -- enough to fill the visible area.
4. **For images, use `aspect-ratio`** so the skeleton has the right shape:
   ```tsx
   <div className="aspect-square bg-slate-200 rounded" />
   ```
5. **Respect reduced motion**:
   ```css
   @media (prefers-reduced-motion: reduce) {
   	.animate-pulse {
   		animation: none;
   		opacity: 0.6;
   	}
   }
   ```

---

## 10. Infinite Scroll

For paginated lists -- load more as the user scrolls near the bottom.

### Implementation with IntersectionObserver

```tsx
function useInfiniteScroll onLoadMore: () => void, hasMore: boolean) {
  const sentinelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!hasMore) return;
    const sentinel = sentinelRef.current;
    if (!sentinel) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          onLoadMore();
        }
      },
      { rootMargin: "200px" }  // trigger 200px before reaching the bottom
    );

    observer.observe(sentinel);
    return () => observer.disconnect();
  }, [onLoadMore, hasMore]);

  return sentinelRef;
}

// Usage
function PendudukList() {
  const [items, setItems] = useState<PendudukItem[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  const loadMore = useCallback(async () => {
    if (loading) return;
    setLoading(true);
    const newItems = await fetchPage(page + 1);
    if (newItems.length === 0) {
      setHasMore(false);
    } else {
      setItems((prev) => [...prev, ...newItems]);
      setPage(page + 1);
    }
    setLoading(false);
  }, [page, loading]);

  const sentinelRef = useInfiniteScroll(loadMore, hasMore);

  return (
    <div>
      <ul>
        {items.map((item) => (
          <PendudukCard key={item.id} item={item} />
        ))}
      </ul>
      <div ref={sentinelRef} className="h-4" />
      {loading && <PendudukCardSkeleton />}
      {!hasMore && (
        <p className="text-center text-sm text-slate-500 py-4">
          Tidak ada data lagi
        </p>
      )}
    </div>
  );
}
```

### Critical details

1. **`rootMargin: "200px"`** -- trigger loading before the user reaches the bottom, so they never see the loading state (if the network is fast).
2. **`if (loading) return`** -- prevent double-loading.
3. **Show a sentinel + skeleton** -- the sentinel is invisible, the skeleton provides visual feedback.
4. **Show "no more data"** when `hasMore` is false -- without this, users keep scrolling and think the page is broken.
5. **Don't reset to page 1 on filter change** -- invalidate `hasMore`, clear items, reset page, re-fetch.
6. **For accessibility**, also include a "Load more" button as a fallback -- keyboard users can't trigger scroll.

### Don't use infinite scroll when

- Users need to reach the footer (use "Load more" button instead).
- The list is ordered and users need to jump to a specific position.
- The list is short enough to load all at once.
