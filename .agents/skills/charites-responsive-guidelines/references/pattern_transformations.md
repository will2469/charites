# Mobile Pattern Transformations & Container Queries Reference

> Transformasi struktural paradigma antarmuka dari desktop ke ponsel pintar: Centered Modal $\rightarrow$ Bottom Sheet dan Wide Data Table $\rightarrow$ Progressive-Disclosure Cards via CSS Container Queries.

---

## 1. Transformasi Modal Tengah Menjadi Bottom Sheet

Pada layar desktop yang lebar, dialog modal di tengah layar (_centered modal_) bekerja optimal karena dekat dengan kursor mouse. Namun pada ponsel pintar:

- Dialog di tengah jatuh pada _Stretch Zone_ yang sulit dijangkau jempol.
- Dialog mudah terpotong atau terdorong keluar saat keyboard virtual muncul.

### Solusi Responsif: Hybrid Bottom Sheet

Komponen dialog bertransformasi menjadi _Bottom Sheet_ di layar kecil (< 768px) dan menjadi _Centered Dialog_ di layar lebar:

```css
/* Mobile: Bottom Sheet di zona alami jempol */
.responsive-dialog {
	position: fixed;
	inset: auto 0 0 0;
	width: 100%;
	max-height: 85svh; /* Gunakan svh agar tidak terpotong saat URL bar terlihat */
	background-color: var(--surface-card);
	border-top-left-radius: var(--primitive-radius-default);
	border-top-right-radius: var(--primitive-radius-default);
	padding: var(--primitive-space-4);
	padding-bottom: calc(var(--primitive-space-4) + env(safe-area-inset-bottom));
	touch-action: pan-y;
	box-shadow: 0 -4px 16px rgba(0, 0, 0, 0.08);
}

/* Desktop (Tablet ke atas): Centered Modal */
@media (min-width: 768px) {
	.responsive-dialog {
		inset: 50% auto auto 50%;
		transform: translate(-50%, -50%);
		width: 100%;
		max-width: 32rem;
		max-height: 80vh;
		border-radius: var(--primitive-radius-default);
		padding-bottom: var(--primitive-space-4);
		box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
	}
}
```

---

## 2. Transformasi Tabel Data Lebar Menjadi Progressive-Disclosure Cards

Membungkus tabel banyak kolom dengan scroll horizontal (`overflow-x: auto`) pada layar ponsel adalah sumber frustrasi:

- Konteks kolom tersembunyi di luar layar.
- Sulit membandingkan baris.
- Gesekan horizontal berbenturan dengan gesture navigasi bawaan peramban/OS.

### Solusi Modern: CSS Container Queries (`@container`)

Daripada mengandalkan media query viewport, tabel dibungkus dalam _container query_ berukuran lokal (`inline-size`). Saat kontainer menyempit (< 640px), baris tabel secara otomatis bertransformasi menjadi **kartu modular bertumpuk**:

```html
<div class="table-boundary">
	<table class="adaptive-table">
		<thead class="adaptive-thead">
			<tr>
				<th scope="col">No. Berkas</th>
				<th scope="col">Nama Warga</th>
				<th scope="col">Layanan</th>
				<th scope="col">Status</th>
			</tr>
		</thead>
		<tbody class="adaptive-tbody">
			<tr class="adaptive-row">
				<td data-label="No. Berkas" class="font-mono">REG-2026-0814</td>
				<td data-label="Nama Warga">Bambang Wijaya</td>
				<td data-label="Layanan">Surat Keterangan Usaha</td>
				<td data-label="Status"><span class="badge-success">Selesai</span></td>
			</tr>
		</tbody>
	</table>
</div>
```

```css
/* 1. Definisikan Batas Kontainer Lokal */
.table-boundary {
	container-type: inline-size;
	container-name: data-container;
	width: 100%;
}

.adaptive-table {
	width: 100%;
	border-collapse: collapse;
	font-size: var(--font-size-sm);
	text-align: left;
}

.adaptive-table th,
.adaptive-table td {
	padding: var(--primitive-space-3);
	border-bottom: 1px solid var(--border-subtle);
}

/* 2. Transformasi Kartu Otomatis Saat Kontainer < 640px */
@container data-container (max-width: 640px) {
	.adaptive-thead {
		display: none; /* Sembunyikan header tabel horizontal */
	}

	.adaptive-table,
	.adaptive-tbody,
	.adaptive-row {
		display: block;
		width: 100%;
	}

	.adaptive-row {
		margin-bottom: var(--primitive-space-3);
		padding: var(--primitive-space-3);
		background-color: var(--surface-card);
		border: 1px solid var(--border-subtle);
		border-radius: var(--primitive-radius-subtle);
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
	}

	.adaptive-table td {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--primitive-space-2) 0;
		border-bottom: 1px dashed var(--border-subtle);
		text-align: right;
	}

	.adaptive-table td:last-child {
		border-bottom: none;
	}

	/* Munculkan label kolom dari atribut data-label */
	.adaptive-table td::before {
		content: attr(data-label);
		font-weight: 600;
		color: var(--text-muted);
		text-align: left;
		padding-right: var(--primitive-space-2);
	}
}
```

---

## 3. Keunggulan Intrinsik Container Queries

1. **Decoupled from Viewport:** Komponen tabel di atas dapat ditempatkan di mana saja: di halaman penuh, di dalam panel samping (_sidebar_), di dalam dialog modal, atau di layout split multi-kolom.
2. **Zero JavaScript Reflow:** Transformasi dieksekusi 100% oleh CSS engine peramban, menghasilkan performa rendering 60 FPS tanpa beban CPU thread utama.

---

## 4. Sitasi & Standar

- W3C CSS Working Group. (2024-2026). _CSS Containment Module Level 3 (Container Queries & Container Units)._
- Wroblewski, L. (2018). _Responsive Navigation & Table Patterns._ LukeW Ideation.
