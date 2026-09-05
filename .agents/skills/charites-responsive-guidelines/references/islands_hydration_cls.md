# Islands Hydration Discipline & CLS Optimization Reference

> Panduan kinerja arsitektur Astro SSR + React Islands: eliminasi JavaScript berlebih (Zero-JS default), pencegahan Flash of Unstyled Content (FOUC), stabilisasi layout ($CLS = 0$), dan integrasi komponen headless aksesibel.

---

## 1. Zero-JS by Default: Rationale Kinerja Mobile

Aplikasi berbasis Single Page Application (SPA) tradisional mengirimkan megabyte bundel JavaScript ke peramban klien. Pada ponsel murah (_low-tier mobile hardware_), proses parsing dan kompilasi script memakan CPU thread utama secara intensif, menyebabkan First Contentful Paint lambat (FCP > 2.5s) dan Total Blocking Time (TBT) yang tinggi.

Astro memecahkan masalah ini dengan **Islands Architecture**:

- Seluruh halaman dirender menjadi HTML statis dan CSS di sisi server.
- Secara default, **0 byte JavaScript klien yang dikirim ke browser**.
- Komponen interaktif diisolasi sebagai "pulau" (_island_) yang hanya dihidrasi ketika memenuhi direktif eksplisit.

### Matriks Direktif Hidrasi Astro

| Direktif Hidrasi             | Mekanisme Eksekusi                                              | Profil Beban CPU                                                   | Rekomendasi Penggunaan di Charites                                                        |
| :--------------------------- | :-------------------------------------------------------------- | :----------------------------------------------------------------- | :---------------------------------------------------------------------------------------- |
| **Default (Tanpa Direktif)** | Murni server-rendered HTML. **0 KB JS klien**.                  | Nol beban CPU klien; first paint instan.                           | Layout halaman, tabel data statis, kartu ringkasan, teks dokumentasi.                    |
| **`client:load`**            | Hidrasi langsung saat dokumen di-parse.                         | Prioritas tinggi; gunakan hemat agar tidak menghambat TBT.         | Interaksi kritis _above-the-fold_ (misal: bilah pencarian cepat, notifikasi penting).      |
| **`client:idle`**            | Hidrasi via `requestIdleCallback()` saat thread utama senggang. | Moderat; menunda JS sampai halaman selesai dirender.               | Panel filter data, sakelar tema (_theme toggle_), formulir feedback.                      |
| **`client:visible`**         | Hidrasi via `IntersectionObserver` saat elemen masuk viewport.  | Rendah; menunda eksekusi JS sampai pengguna menggulir ke komponen. | Grafik analitik, linimasa aktivitas, formulir transaksi di bawah lipatan layar.          |
| **`client:media="(query)"`** | Hidrasi hanya jika CSS media query bernilai benar.              | Nol beban di luar ukuran layar yang cocok.                         | Panel kontrol khusus desktop, multi-panel editor.                                         |
| **`client:only="react"`**    | Lewati SSR sepenuhnya; render murni di klien.                   | Tidak ada serialisasi server, butuh alokasi tinggi agar tidak CLS. | Autentikasi sesi privat, widget interaktif kompleks.                                     |

---

## 2. Pencegahan FOUC & Hydration Mismatch

_Hydration mismatch_ terjadi saat HTML hasil render server berbeda dengan pohon DOM awal di klien, biasanya karena pemanggilan API peramban non-deterministik (`localStorage`, `window.matchMedia`). Hal ini sering memicu _Flash of Unstyled Content_ (FOUC) saat tema gelap/terang berganti.

### A. Skrip Tema Blocking di `<head>`

Agar peramban menerapkan tema sebelum melukiskan frame pertama (_first paint_), sisipkan skrip inline di `<head>`:

```html
<!-- src/layouts/RootLayout.astro -->
<head>
	<meta charset="UTF-8" />
	<meta
		name="viewport"
		content="width=device-width, initial-scale=1.0, viewport-fit=cover, interactive-widget=resizes-content"
	/>

	<script is:inline>
		(function () {
			const storedTheme = localStorage.getItem("charites_theme_preference");
			const systemDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
			const theme = storedTheme || (systemDark ? "dark" : "light");
			document.documentElement.setAttribute("data-theme", theme);
		})();
	</script>
</head>
```

### B. Sinkronisasi React via `useSyncExternalStore`

Hindari `useEffect` untuk mendeteksi tema agar tidak memicu re-render ganda dan layout shift:

```typescript
import React, { useSyncExternalStore } from "react";

function subscribeTheme(callback: () => void): () => void {
  const observer = new MutationObserver(callback);
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ["data-theme"] });
  return () => observer.disconnect();
}

export function ThemeToggle(): React.JSX.Element {
  const currentTheme = useSyncExternalStore(
    subscribeTheme,
    () => document.documentElement.getAttribute("data-theme") || "light",
    () => "light" // Fallback SSR
  );

  const toggleTheme = () => {
    const next = currentTheme === "dark" ? "light" : "dark";
    localStorage.setItem("charites_theme_preference", next);
    document.documentElement.setAttribute("data-theme", next);
  };

  return (
    <button
      type="button"
      onClick={toggleTheme}
      className="min-h-[44px] min-w-[44px] px-3 py-2 border border-border rounded-subtle bg-card hover:bg-card-active focus-visible:ring"
      aria-label={`Ganti tema. Status saat ini: ${currentTheme}`}
    >
      <span className="font-mono text-sm uppercase">{currentTheme}</span>
    </button>
  );
}
```

---

## 3. Stabilisasi Layout ($CLS = 0$)

Saat komponen dihidrasi secara asinkron via `client:visible` atau `client:idle`, ketiadaan alokasi tinggi kontainer dapat mendorong elemen di bawahnya, memicu Cumulative Layout Shift (CLS).

Gunakan `content-visibility: auto` dan `contain-intrinsic-size` pada pembungkus island:

```css
.island-reservation-box {
	content-visibility: auto;
	contain-intrinsic-size: auto 340px;
	min-height: 340px;
	width: 100%;
}
```

---

## 4. Komponen Headless Aksesibel (Radix / Ark UI)

Untuk menjamin kepatuhan penuh WAI-ARIA 1.2/1.3 tanpa menulis kode perangkap fokus (_focus trap_) manual yang rentan bug, bungkus primitif headless dengan token desain semantik:

```typescript
import * as Dialog from "@radix-ui/react-dialog";
import React from "react";

export function AccessibleActionDialog({ triggerLabel, title, description, confirmLabel, onConfirm }: any) {
  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <button type="button" className="min-h-[44px] px-4 py-2 bg-card border border-border rounded-subtle hover:bg-card-active focus-visible:ring">
          {triggerLabel}
        </button>
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 bg-black/40 backdrop-blur-sm" />
        <Dialog.Content className="fixed bottom-0 md:bottom-auto md:top-1/2 left-0 md:left-1/2 md:-translate-x-1/2 md:-translate-y-1/2 w-full md:max-w-lg bg-card border border-border p-6 rounded-t-default md:rounded-default focus:outline-none">
          <Dialog.Title className="text-lg font-semibold text-primary">{title}</Dialog.Title>
          <Dialog.Description className="mt-2 text-sm text-muted">{description}</Dialog.Description>
          <div className="mt-6 flex flex-col-reverse md:flex-row md:justify-end gap-3">
            <Dialog.Close asChild>
              <button type="button" className="min-h-[44px] px-4 py-2 border border-border rounded-subtle hover:bg-card-active">Batal</button>
            </Dialog.Close>
            <button type="button" onClick={onConfirm} className="min-h-[44px] px-4 py-2 bg-danger hover:bg-danger-hover text-white rounded-subtle focus-visible:ring">
              {confirmLabel}
            </button>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
```

---

## 5. Sitasi & Standar

- Astro Technology Corporation. (2024-2026). _Astro Islands Architecture & Hydration Reference._
- W3C Web Accessibility Initiative. (2023). _WAI-ARIA 1.2 Authoring Practices Guide (Modal Dialogs & Focus Traps)._
- Google Chrome Web Vitals Team. (2024). _Optimizing Cumulative Layout Shift (CLS) in Dynamic Client Hydration._
