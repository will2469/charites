# Design Tokens SSOT & Automated Governance Reference

> Spesifikasi arsitektur 3-Tier Design Tokens W3C DTCG (v2025.10), implementasi CSS Custom Properties di `global.css`, pemetaan Tailwind v4, dan penegakan tata kelola otomatis.

---

## 1. Arsitektur 3-Tier W3C DTCG Standard

Untuk membasmi fragmentasi desain dan nilai sembarangan, setiap keputusan visual dikodifikasi ke dalam tiga lapisan terpisah:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       THE 3-TIER DESIGN TOKEN HIERARCHY                     │
│                                                                             │
│  [TIER 1: GLOBAL PRIMITIVES] ──► Nilai mentah independen                    │
│                                  (slate-50, crimson-500, space-4, radius-8) │
│                                         │                                   │
│  [TIER 2: SEMANTIC TOKENS]   ──► Niat fungsional & Pemetaan Tema            │
│                                  (--surface-canvas, --text-primary)         │
│                                         │                                   │
│  [TIER 3: COMPONENT TOKENS]  ──► Aturan scoped spesifik UI                  │
│                                  (--btn-primary-bg, --dialog-radius)        │
└─────────────────────────────────────────────────────────────────────────────┘
```

1. **Tier 1: Global / Primitive Tokens:**
   - Mewakili nilai absolut yang tidak memiliki konteks bisnis.
   - Dilarang keras dipanggil langsung di dalam komponen antarmuka pengguna (`components/` atau `pages/`).
2. **Tier 2: Semantic / Intent-Driven Tokens:**
   - Memetakan nilai primitif ke dalam fungsi operasional aplikasi.
   - Contoh: `--surface-canvas` (latar belakang root), `--surface-card` (kontainer kartu), `--text-primary` (konten baca utama), `--border-subtle` (garis batas pemisah).
   - Mode gelap (_Dark Mode_) dieksekusi **hanya** pada lapisan ini melalui pergantian alias data attribute `[data-theme="dark"]`, tanpa perlu mengubah kode komponen di hilir.
3. **Tier 3: Component-Scoped Tokens:**
   - Nilai styling khusus untuk komponen atomik tertentu (misal `--badge-status-radius`, `--table-cell-padding`).

---

## 2. Implementasi CSS Custom Properties (`global.css` SSOT)

Kompilasi token wajib tinggal di satu file terpusat (`global.css`):

```css
:root {
	/* --- Tier 1: Primitives Scale --- */
	--primitive-slate-50: #f8fafc;
	--primitive-slate-100: #f1f5f9;
	--primitive-slate-200: #e2e8f0;
	--primitive-slate-700: #334155;
	--primitive-slate-800: #1e293b;
	--primitive-slate-900: #0f172a;
	--primitive-crimson-500: #ef4444;
	--primitive-crimson-600: #dc2626;

	/* Spacing: Modular 4px/8px Base */
	--primitive-space-1: 0.25rem; /* 4px */
	--primitive-space-2: 0.5rem; /* 8px */
	--primitive-space-3: 0.75rem; /* 12px */
	--primitive-space-4: 1rem; /* 16px */
	--primitive-space-6: 1.5rem; /* 24px */
	--primitive-space-8: 2rem; /* 32px */

	/* Radius: Subdued Industrial */
	--primitive-radius-none: 0px;
	--primitive-radius-subtle: 4px;
	--primitive-radius-default: 8px;

	/* Typography */
	--font-family-sans:
		"Inter Variable", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
	--font-family-mono: "JetBrains Mono Variable", monospace;
	--font-size-sm: 0.875rem;
	--font-size-base: 1rem;
	--font-size-lg: 1.125rem;
	--font-size-xl: 1.25rem;
	--line-height-snug: 1.25;
	--line-height-base: 1.5;

	/* --- Tier 2: Semantic Intent (Light Mode Default) --- */
	--surface-canvas: var(--primitive-slate-50);
	--surface-card: #ffffff;
	--surface-active: var(--primitive-slate-100);
	--text-primary: var(--primitive-slate-900);
	--text-muted: var(--primitive-slate-700);
	--border-subtle: var(--primitive-slate-200);
	--interactive-focus-ring: #0284c7;
	--status-danger: var(--primitive-crimson-500);
	--status-danger-hover: var(--primitive-crimson-600);
}

/* Deterministic Dark Mode Mutation (Zero JS Reflow) */
[data-theme="dark"] {
	--surface-canvas: var(--primitive-slate-900);
	--surface-card: var(--primitive-slate-800);
	--surface-active: var(--primitive-slate-700);
	--text-primary: var(--primitive-slate-50);
	--text-muted: var(--primitive-slate-200);
	--border-subtle: #334155;
	--interactive-focus-ring: #38bdf8;
	--status-danger: #f87171;
	--status-danger-hover: #ef4444;
}

/* Mechanical Focus Affordance */
:focus-visible {
	outline: 2px solid var(--interactive-focus-ring);
	outline-offset: 2px;
}
```

---

## 3. Pemetaan Tailwind CSS v4

Untuk mencegah agen menyisipkan utilitas arbitrer (`bg-[#6366F1]` atau `p-[17px]`), Tailwind dibatasi hanya pada token semantik resmi:

```typescript
// tailwind.config.ts
import type { Config } from "tailwindcss";

const config: Config = {
	darkMode: ["class", '[data-theme="dark"]'],
	content: ["./src/**/*.{astro,html,js,jsx,ts,tsx}"],
	theme: {
		colors: {
			transparent: "transparent",
			current: "currentColor",
			canvas: "var(--surface-canvas)",
			card: "var(--surface-card)",
			"card-active": "var(--surface-active)",
			primary: "var(--text-primary)",
			muted: "var(--text-muted)",
			border: "var(--border-subtle)",
			ring: "var(--interactive-focus-ring)",
			danger: {
				DEFAULT: "var(--status-danger)",
				hover: "var(--status-danger-hover)",
			},
		},
		spacing: {
			0: "0px",
			1: "var(--primitive-space-1)",
			2: "var(--primitive-space-2)",
			3: "var(--primitive-space-3)",
			4: "var(--primitive-space-4)",
			6: "var(--primitive-space-6)",
			8: "var(--primitive-space-8)",
		},
		borderRadius: {
			none: "var(--primitive-radius-none)",
			subtle: "var(--primitive-radius-subtle)",
			default: "var(--primitive-radius-default)",
		},
	},
};

export default config;
```

---

## 4. Tata Kelola Otomatis (CI Lint Gates)

Pengendalian mutu di CI menggunakan plugin Stylelint `stylelint-declaration-strict-value`:

```javascript
// .stylelintrc.cjs
module.exports = {
	plugins: ["stylelint-declaration-strict-value"],
	rules: {
		"scale-unlimited/declaration-strict-value": [
			["/color$/", "background-color", "border-color", "fill", "stroke"],
			{
				ignoreValues: ["transparent", "currentColor", "inherit", "initial"],
				message:
					"Visual Defect: Nilai warna hardcoded melanggar kontrak token. Gunakan var(--surface-*), var(--text-*), atau var(--border-*).",
			},
		],
		"scale-unlimited/declaration-strict-value": [
			["margin", "padding", "gap", "top", "right", "bottom", "left"],
			{
				ignoreValues: ["0", "auto", "unset", "inherit"],
				message: "Spatial Drift: Spasi tata letak wajib merujuk ke langkah token desain resmi.",
			},
		],
	},
};
```

---

## 5. Sitasi & Standar

- W3C Design Tokens Community Group. (2025). _Design Tokens Format Module v2025.10._
- Frost, B. (2016). _Atomic Design._ Brad Frost Web.
- CSS Working Group. (2024-2026). _CSS Custom Properties for Cascading Variables Module Level 1._
