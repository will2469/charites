# Mobile Touch Ergonomics & Viewport Dynamics Reference

> Pemodelan matematis ergonomi sentuhan jempol (Hukum Fitts), standar ukuran target sentuh WCAG 2.2, unit viewport modern (`svh`/`dvh`), dan penanganan keyboard virtual.

---

## 1. Hukum Fitts & Ergonomi Jempol (The Thumb Zone)

### A. Pemodelan Matematis Hukum Fitts (Paul Fitts, 1954)

Waktu tempuh jari pengguna menuju target interaktif ($MT$) dimodelkan sebagai:
$$MT = a + b \cdot ID = a + b \log_2 \left( \frac{2D}{W} \right)$$

- $D$ (_Distance_): Jarak fisik dari posisi istirahat jempol ke target.
- $W$ (_Width_): Lebar bidang sentuh interaktif target.
- $ID$ (_Index of Difficulty_): Indeks kesulitan dalam bit.
- $a, b$: Konstanta empiris biomekanik manusia.

> **Hukum Rekayasa:** Untuk mempercepat eksekusi aksi pengguna pada ponsel cerdas, kita harus **memperpendek jarak $D$** (menempatkan aksi di bagian bawah layar yang dijangkau jempol) atau **memperbesar lebar $W$** (memperluas hit area minimal 44×44px).

### B. Pemetaan Tiga Zona Jempol (Steven Hoober & Luke Wroblewski)

1. **Natural Zone (Bawah & Tengah):** Area yang dijangkau jempol secara alami tanpa meregangkan sendi ibu jari. Wajib untuk: tombol aksi primer, form submit, bottom navigation, dan drawer trigger.
2. **Reach Zone (Tengah Atas):** Area yang dapat dijangkau dengan sedikit peregangan. Cocok untuk: konten kartu, tabel pembacaan data, filter sekunder.
3. **Stretch Zone (Pojok Atas):** Area sulit yang membutuhkan perpindahan genggaman tangan (_grip repositioning_). Hanya untuk aksi destruktif atau navigasi mundur (_back button_). Dilarang keras menaruh tombol submit utama di area ini!

---

## 2. Standar Target Sentuh (WCAG 2.2 SC 2.5.8)

- **Ketentuan Minimum WCAG 2.2 Level AA:** Area sentuh minimal **24×24 CSS piksel**, kecuali terdapat jarak pemisah (_spacing clearance_) yang mencukupi.
- **Standar Rekomendasi Charites (Apple HIG & Material 3):**
  - Tombol primer/sekunder: Minimal tinggi **44px** atau **48px**.
  - Ikon interaktif (tanpa teks): Padding tak kasat mata diperluas agar hit-area mencapai minimal 44×44px (`min-h-[44px] min-w-[44px]`).

---

## 3. Unit Viewport Modern: `svh`, `lvh`, `dvh` vs `100vh`

Penggunaan `100vh` pada peramban seluler (iOS Safari, Android Chrome) adalah sumber utama bug tata letak, karena `100vh` mengabaikan ekspansi dan kolaps bilah alamat URL peramban.

| Unit CSS  | Kepanjangan             | Perilaku Bilah Alamat (Browser Chrome)                                                                       | Skenario Penggunaan                                                                                  |
| :-------- | :---------------------- | :----------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------- |
| **`svh`** | Small Viewport Height   | Dihitung saat bilah alamat **sedang terbuka penuh** (area layar paling sempit). Nilai stabil saat di-scroll. | **Wajib untuk:** Bottom action bar, modal dialog, formulir input penting yang tidak boleh terpotong. |
| **`lvh`** | Large Viewport Height   | Dihitung saat bilah alamat **sedang tertutup/kolaps** (area layar paling luas).                              | Latar belakang hero full-bleed, kanvas promosi scrollable.                                           |
| **`dvh`** | Dynamic Viewport Height | Menyesuaikan ukuran secara dinamis _real-time_ mengikuti buka-tutup bilah alamat saat di-scroll.             | Pembungkus layout utama (gunakan dengan hati-hati karena dapat memicu reflow saat scroll aktif).     |
| **`vh`**  | Legacy Height           | Tidak stabil dan meluap di balik bilah URL.                                                                  | **Dilarang keras pada mobile.**                                                                      |

---

## 4. Penanganan On-Screen Keyboard (OSK) Occlusion

Saat input teks difokuskan, keyboard virtual perangkat lunak muncul dan sering kali menutupi form serta tombol submit.

### A. Meta Viewport Standar W3C

```html
<meta
	name="viewport"
	content="width=device-width, initial-scale=1.0, viewport-fit=cover, interactive-widget=resizes-content"
/>
```

- `interactive-widget=resizes-content`: Memaksa keyboard virtual mengubah ukuran _layout viewport_, sehingga elemen `position: sticky` atau fixed footer otomatis terangkat di atas keyboard tanpa bantuan JavaScript.

### B. Fallback `window.visualViewport` (iOS Safari Compatibility)

Karena dukungan iOS Safari terhadap `interactive-widget` bervariasi, sertakan penyesuaian visual viewport:

```typescript
export function initVirtualKeyboardAdjuster(): () => void {
	const vv = window.visualViewport;
	if (!vv) return () => {};

	const handleResize = () => {
		const layoutHeight = document.documentElement.clientHeight;
		const visualHeight = vv.height;
		const offset = Math.max(0, layoutHeight - visualHeight);
		document.documentElement.style.setProperty("--virtual-keyboard-offset", `${offset}px`);
	};

	vv.addEventListener("resize", handleResize);
	vv.addEventListener("scroll", handleResize);
	handleResize();

	return () => {
		vv.removeEventListener("resize", handleResize);
		vv.removeEventListener("scroll", handleResize);
	};
}
```

```css
/* Bottom bar yang selalu mengambang aman di atas keyboard */
.mobile-bottom-bar {
	position: fixed;
	left: 0;
	right: 0;
	bottom: 0;
	padding-bottom: calc(env(safe-area-inset-bottom, 16px) + var(--virtual-keyboard-offset, 0px));
}
```

---

## 5. Safe-Area Insets (`env(safe-area-inset-*)`)

Untuk perangkat berponi, Dynamic Island, dan gesture pill Android:

1. Wajib mendeklarasikan `viewport-fit=cover` pada meta viewport.
2. Seluruh kontainer tepi layar wajib memperhitungkan safe-area:

```css
padding-top: max(var(--primitive-space-4), env(safe-area-inset-top));
padding-bottom: max(var(--primitive-space-4), env(safe-area-inset-bottom));
padding-left: max(var(--primitive-space-4), env(safe-area-inset-left));
padding-right: max(var(--primitive-space-4), env(safe-area-inset-right));
```

---

## 6. Sitasi & Standar

- Fitts, P. M. (1954). _The information capacity of the human motor system in controlling the amplitude of movement._ Journal of Experimental Psychology.
- Hoober, S. (2013). _How Do Users Really Hold Flagship Smartphones?_ UXmatters.
- Wroblewski, L. (2011). _Mobile First._ A Book Apart.
- W3C CSS Working Group. (2024). _CSS Values and Units Module Level 4 (Viewport Units)._
- W3C Viewport Working Group. (2024). _Virtual Keyboard and Viewport Extension Specification._
