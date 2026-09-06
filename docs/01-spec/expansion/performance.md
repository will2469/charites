# EXPANSION-BATCH-11: Framework & Build Engine Optimization Standards (`performance.*`)
> **Kode Dokumen:** `SPEC-EXP-11-PERFORMANCE`
> **Kategori:** `performance` (Framework, Compiler & Build-Time Performance Optimization; Alias CLI: `perf`)
> **Pilar:** `01-SPEC` (WHAT - Spesifikasi Perilaku & Kontrak Rule)
> **Status:** Calibrated & Peer-Reviewed Expansion Specification (16 Aturan Terkurasi: 4 Wave × 4 Aturan)
> **Kalibrasi Desain:** Multi-Framework Graph Engine (React 18/19 & Compiler Awareness, Astro Islands Architecture, Tailwind CSS v4 Oxide)
> **Migrasi Sumber:** [`charites-legacy/perf-checker.ts`](file:///home/will/Monorepo/charites/charites-legacy/perf-checker.ts)
> **Standar Rujukan:**
> - React Official Architecture: Reconciliation Identity, Referential Stability & React Compiler Integration
> - React Lifecycle Standards: Effect Cleanup Symmetry & Derived State Elimination ("You Might Not Need an Effect")
> - Astro Island Architecture: Zero-JS Default Paradigm, Island Hydration Boundary & Asset Pipeline Integration
> - Tailwind CSS v4 CSS-First Architecture: Static Token Scanner (Oxide), `@source` Monorepo Discovery & `@utility` Hygiene
> **Pilar Terkait:** [01-SPEC: cls.md](cls.md), [01-SPEC: inp.md](inp.md), [01-SPEC: lcp.md](lcp.md), [01-SPEC: themes.md](themes.md), & [01-SPEC: responsive.md](responsive.md)

---

## 1. Epistemologi & Pemisahan Lapisan: Runtime Browser vs Framework & Build Engine

Untuk memahami posisi kategori `performance.*` (alias `perf.*`), Charites membagi arsitektur evaluasi performa ke dalam dua domain yang berbeda secara fundamental:

```text
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                        ARSIKTEKTUR DUA LAPISAN EVALUASI PERFORMA                       │
├────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                        │
│  [LAPISAN 1: BROWSER RUNTIME METRICS (Core Web Vitals - Symptom Layer)]                │
│  • cls.*  : Stabilitas geometri tata letak & eliminasi pergeseran tak terduga (CLS)    │
│  • inp.*  : Responsivitas event loop, pemrosesan tugas panjang & kelancaran input (INP)│
│  • lcp.*  : Waktu penemuan aset kritis di browser & siklus First Contentful Paint (LCP) │
│  → Karakteristik: Berbasis gejala runtime yang dirasakan pengguna (*user-perceived*). │
│  → Hakim Terakhir: Web browser riil (DevTools, LoAF API, CrUX, Lighthouse CI).         │
│                                                                                        │
│  ─────────────────────────────────── BATAS TEGAS ───────────────────────────────────  │
│                                                                                        │
│  [LAPISAN 2: FRAMEWORK & BUILD ENGINE OPTIMIZATION (performance.* - Root-Cause Layer)] │
│  • React  : Rekonsiliasi VDOM, kestabilan referensi memori, memory leak & compiler mode│
│  • Astro  : Penegakan Zero-JS bawaan, isolasi pulau hidrasi & utilisasi aset pipeline │
│  • Tailwind: Kompatibilitas scanner statis v4, resolusi @source & payload CSS minimal   │
│  → Karakteristik: Berbasis efisiensi abstraksi kode & kontrak kompilasi bundler.       │
│  → Hakim Terakhir: AST Static Analyzer & Bundler Compiler (0 B/op, 0 allocs/op).       │
│                                                                                        │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

### 1.1. Mengapa Kategori `performance.*` Unik dan Wajib Ada?
Jika kategori `cls.*`, `inp.*`, dan `lcp.*` berfokus pada **perilaku render dan interaksi yang dialami peramban di sisi pengguna**, kategori `performance.*` berfokus pada **efisiensi penulisan kode pada tingkat abstraksi framework dan compiler**:
1. **Pemanfaatan Fitur Framework Modern:** Memastikan pengembang memanfaatkan kemampuan native React 18/19 (misal: code splitting via `React.lazy`, stabilisasi referensi), Astro (misal: pemrosesan aset lokal via `astro:assets`, prinsip Zero-JS), dan Tailwind CSS v4 (CSS-first scanner, direktif `@source`).
2. **Pencegahan Anti-Pola Abstraksi:** Mencegah kode yang secara diam-diam memboroskan memori (misal: `useEffect` tanpa fungsi cleanup), merusak algoritma rekonsiliasi VDOM (misal: `key={index}` pada mutable list), atau menggelembungkan ukuran berkas produksi (misal: pengiriman JavaScript untuk komponen yang sebenarnya murni statis).

---

## 2. Garansi Zero Redundancy, Registry SSOT & Precedensi

### 2.1. Tabel Rekonsiliasi Aturan Warisan (`perf-checker.ts`)
Untuk mencegah redundansi, aturan dari skrip warisan [`charites-legacy/perf-checker.ts`](file:///home/will/Monorepo/charites/charites-legacy/perf-checker.ts) yang telah memiliki domain spesifik di Core Web Vitals atau Tema **TIDAK DIMASUKKAN LAGI** ke dalam `performance.*`:

| Aturan Legacy (`perf-checker.ts`) | Status Rekonsiliasi | Dialihkan / Diserap ke | Alasan Pengalihan & Batas Domain |
|---|:---:|---|---|
| `effect-no-deps` (R3) | **Diserap** | `inp.unbounded-effect-deps` | `useEffect` tanpa dependensi memicu infinite loop / long task yang menghantam main thread (domain INP). |
| `provider-inline-value` (R6) | **Diserap** | `inp.context-re-render-cascade` | Objek literal inline pada provider value memicu re-render massal saat interaksi (domain INP). |
| `too-many-client-load` (R10) | **Diserap** | `inp.hydration-contention` | Penumpukan hidrasi `client:load` awal memicu kompetisi CPU pada first load (domain INP). |
| `no-client-only` (R11) | **Diserap** | `cls.client-only-hydration-pop` & `lcp.client-only-lcp-content` | Ketiadaan fallback `client:only` menyebabkan layout jump (CLS) dan hilangnya elemen hero di HTML awal (LCP). |
| `tailwind-apply-bloat` (R16) | **Diserap** | `theme.apply-bloat` | Penumpukan utility classes di dalam `@apply` adalah masalah tata kelola token CSS desain (domain Themes). |
| `tailwind-content-missing` (R14) | **Dimodernisasi** | `performance.tailwind-untracked-package-source` | Di Tailwind v4, file config JS sudah usang. Digantikan oleh audit direktif `@source` untuk monorepo. |
| `tailwind-content-no-astro` (R15) | **Dimodernisasi** | `performance.tailwind-dynamic-class-concatenation` | Di Tailwind v4, file `.astro` discan otomatis; yang merusak scanner adalah penggabungan string dinamis. |

### 2.2. Matriks Ortogonalitas Lintas Kategori (Genuine Nearest Neighbors & Root-Cause Mapping)

Setiap aturan dalam `performance.*` dipasangkan dengan aturan terdekat yang nyata atau ditelusuri ke dampak metrik runtime browser (*Root-Cause $\to$ Symptom Mapping*):

| Rule `performance.*` (Root Cause) | Rule Terdekat Kategori Lain (Symptom / Partner) | Fokus Domain Kategori Lain | Fokus Domain Kategori `performance.*` | Garansi Batasan Ortogonal |
|---|---|---|---|---|
| `performance.react-inline-prop-memo` | `inp.context-re-render-cascade` | Mengaudit instansiasi objek baru pada `Context.Provider` yang merusak memoization konsumen context global. | Mengaudit pengiriman prop inline (objek, array, fungsi) pada pemanggilan komponen spesifik yang dibungkus `React.memo()`. | `inp` mengaudit pohon langganan Context; `performance` mengaudit kontrak referensial prop pada komponen ter-memoize. |
| `performance.react-index-as-key` | `inp.unbounded-collection-render` | Mengaudit volume elemen DOM masif yang dirender tanpa virtual windowing. | Mengaudit penggunaan indeks array (`key={index}`) sebagai kunci rekonsiliasi VDOM pada koleksi dinamis yang mengalami mutasi/reorder. | `inp` mengaudit kapasitas render DOM fisik; `performance` mengaudit stabilitas identitas rekonsiliasi VDOM saat mutasi data. |
| `performance.react-effect-missing-cleanup` | `inp.unbounded-effect-deps` | Mengaudit frekuensi eksekusi efek akibat ketiadaan array dependensi yang menyita main-thread. | Mengaudit ketiadaan fungsi cleanup simetris pada efek yang mengakuisisi resource persisten (listener, timer, observer, abort controller). | `inp` mengaudit frekuensi eksekusi loop efek; `performance` mengaudit pelepasan sumber daya siklus hidup memori (*memory leak*). |
| `performance.react-context-domain-coupling` | `inp.context-re-render-cascade` | Mengaudit pembuatan objek literal baru pada prop `value` di level instansiasi markup. | Mengaudit penggabungan state yang berbeda frekuensi mutasinya ke dalam satu Context terpadu dengan basis konsumen luas. | `inp` mengaudit instansiasi referensi objek runtime; `performance` mengaudit arsitektur penggabungan domain state konteks. |
| `performance.react-static-heavy-import` | `lcp.client-only-lcp-content` | Mengaudit ketiadaan konten visual hero di SSR awal akibat pemintasan hidrasi di atas pelipatan. | Mengaudit impor statis modul berbobot besar non-kritis (chart, editor, admin) yang memicu pembengkakan bundle JS awal. | `lcp` mengaudit waktu keterlihatan konten hero; `performance` mengaudit pemisahan kode (*code splitting*) bundler. |
| `performance.react-redundant-function-memoization` | `inp.expensive-render-computation` | Mengaudit komputasi algoritma berat yang dieksekusi sinkron di badan fungsi render. | Mengaudit pembungkusan fungsi sederhana dalam `useCallback` yang tidak memiliki konsumen sensitif identitas atau saat React Compiler aktif. | `inp` mengaudit komputasi CPU berat; `performance` mengaudit pemborosan alokasi overhead memori hook yang tidak berguna. |
| `performance.react-derived-state-in-effect` | `inp.repeated-state-update` | Mengaudit mutasi `setState` berulang-ulang di dalam loop sinkron. | Mengaudit pola anti-pattern sinkronisasi state turunan via `useEffect` yang memicu re-render sekunder ganda. | `inp` mengaudit batching mutasi perulangan; `performance` mengaudit arsitektur render pass (kalkulasi inline vs efek sekunder). |
| `performance.react-unstable-hook-reference` | `performance.react-inline-prop-memo` | Mengaudit literal inline langsung pada posisi atribut JSX komponen ter-memoize. | Mengaudit custom React hook yang mengembalikan objek/fungsi baru tanpa memoization dan diteruskan ke konsumen hilir sensitif identitas. | `inline-prop-memo` mengaudit call-site JSX langsung; `unstable-hook-reference` mengaudit kontrak return value custom hook. |
| `performance.astro-unnecessary-client-directive` | `inp.hydration-contention` | Mengaudit saturasi thread akibat kompetisi hidrasi serentak banyak pulau interaktif di awal load. | Mengaudit pemberian direktif hidrasi `client:*` pada komponen Astro yang terbukti murni statis tanpa interaktivitas. | `inp` mengaudit antrean hidrasi komponen aktif; `performance` mengaudit eliminasi pengiriman JS untuk UI statis (Zero-JS). |
| `performance.astro-island-boundary-overlap` | `cls.client-only-hydration-pop` | Mengaudit pergeseran tata letak saat pulau interaktif dihidrasi tanpa reservasi ruang slot. | Mengaudit tumpang-tindih boundary hidrasi independen atau runtime multi-framework yang memicu inkonsistensi state & hidrasi ganda. | `cls` mengaudit stabilitas geometri fisik; `performance` mengaudit integritas boundary hidrasi pulau modular Astro. |
| `performance.astro-unoptimized-local-image` | `lcp.oversized-lcp-resource-selection` | Mengaudit ketiadaan varian responsif `srcset`/`sizes` pada gambar hero kandidat LCP. | Mengaudit penggunaan tag `<img>` biasa untuk berkas gambar lokal yang melewatkan pipeline build kompresi bawaan `astro:assets`. | `lcp` mengaudit efisiensi bandwidth gambar LCP; `performance` mengaudit utilisasi pipeline build framework pada aset lokal. |
| `performance.astro-over-prefetching` | `lcp.missing-critical-origin-hint` | Mengaudit ketiadaan tag `preconnect` untuk origin CDN pihak ketiga tempat aset LCP di-host. | Mengaudit konfigurasi prefetch agresif (`viewport`/`load`) pada tautan navigasi internal yang probabilitas kunjungannya rendah. | `lcp` mengaudit koneksi awal origin primer; `performance` mengaudit efisiensi kuota jaringan pada prefetch navigasi internal. |
| `performance.tailwind-dynamic-class-concatenation` | `theme.unknown-token-usage` | Mengaudit pemakaian token utilitas CSS yang tidak terdaftar di konfigurasi tema. | Mengaudit pembuatan kelas via template literal dinamis yang tidak dapat diekstrak oleh static scanner Tailwind v4 (Oxide). | `theme` mengaudit kepatuhan token desain; `performance` mengaudit keteruraian kode terhadap build scanner compiler. |
| `performance.tailwind-duplicate-arbitrary-rules` | `theme.apply-bloat` | Mengaudit penggunaan `@apply` dengan jumlah kelas utilitas berlebih di dalam aturan CSS. | Mengaudit penggunaan nilai arbitrary ad-hoc berulang yang menghasilkan aturan CSS ekuivalen ganda di output stylesheet. | `theme` mengaudit kebersihan CSS custom; `performance` mengaudit redundansi payload aturan CSS terkompilasi. |
| `performance.tailwind-untracked-package-source` | `theme.token-source-drift` | Mengaudit inkonsistensi sumber kebenaran token desain visual antar berkas. | Mengaudit ketiadaan deklarasi `@source` pada berkas CSS untuk komponen yang diimpor dari paket monorepo eksternal. | `theme` mengaudit keselarasan nilai token; `performance` mengaudit pelacakan dependensi kompilasi bundler Tailwind v4. |
| `performance.tailwind-duplicate-utility-definition` | `theme.apply-bloat` | Mengaudit pemakaian `@apply` yang menumpuk utilitas dalam selector CSS. | Mengaudit deklarasi kustom `@utility` di CSS yang menduplikasi utilitas core bawaan Tailwind v4. | `theme` mengaudit struktur selector CSS; `performance` mengaudit redundansi aturan pada kamus utilitas CSS final. |

### 2.3. Matriks Redundansi & Precedensi Intra-Kategori ($N \times N$)

Tabel ini menjamin aturan di dalam `performance.*` tidak saling tumpang tindih saat mengevaluasi satu simpul AST yang kompleks:

```text
┌──────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────┐
│ Pasangan Aturan Intra-Kategori               │ Resolusi Precedensi & Kebijakan Dedup Pelaporan                                        │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────┤
│ performance.react-inline-prop-memo           │ SCOPE SEPARATION:                                                                      │
│ vs performance.react-unstable-hook-reference │ react-inline-prop-memo HANYA mengaudit ekspresi literal (Object/Array/ArrowFunction)   │
│                                              │ yang ditulis langsung pada posisi atribut prop JSX di call-site komponen.              │
│                                              │ react-unstable-hook-reference HANYA mengaudit return statement di dalam badan          │
│                                              │ deklarasi custom hook. Jika nilai dari hook tak-stabil dioper ke komponen memo, linter │
│                                              │ HANYA memicu 1 warning pada deklarasi custom hook (akar masalah referensial).          │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────┤
│ performance.astro-unnecessary-client-directive│ PRECEDENCE 1:                                                                          │
│ vs performance.astro-island-boundary-overlap │ Jika sebuah pulau terbukti 100% statis (tidak butuh JS sama sekali), aturan            │
│                                              │ unnecessary-client-directive mengambil prioritas utama (hapus direktif hidrasi).       │
│                                              │ Aturan island-boundary-overlap HANYA dievaluasi jika kedua komponen memang             │
│                                              │ terbukti membutuhkan interaktivitas client-side.                                       │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────┤
│ performance.react-redundant-function-memo    │ COMPILER-AWARE CONDITIONAL:                                                            │
│ vs React Compiler Mode                       │ Jika project mengaktifkan React Compiler (Babel/Vite plugin), aturan ini secara        │
│                                              │ otomatis di-suppress atau diturunkan ke level info, karena compiler telah mengotomatisasi│
│                                              │ memoization fungsi dan dependensi.                                                     │
└──────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Ringkasan Matriks 16 Rule `performance.*` (4 Wave Terkalibrasi)

| Wave | Rule ID | Legacy Ref | Domain Parser | Tier | Confidence | Severity | Kelayakan Autofix |
| :---: | :--- | :--- | :---: | :---: | :---: | :---: | :---: |
| **W1** | `performance.react-inline-prop-memo` | R1 | React JSX/TSX AST | T1 | `PROVEN` | `warning` | Menengah (sarankan stabilisasi referensi via `useCallback`/`useMemo`) |
| **W1** | `performance.react-index-as-key` | R2 | React JSX/TSX AST | T1 | `PROVEN` | `warning` | Menengah (sarankan pemakaian properti unik stabil misal `item.id`) |
| **W1** | `performance.react-effect-missing-cleanup` | R4 | React JSX/TSX AST | T1 | `PROVEN` | `error` | **Tinggi** (auto-inject return callback cleanup skeleton) |
| **W1** | `performance.react-context-domain-coupling` | R5 | React JSX/TSX AST | T2 | `LIKELY` | `warning` | Rendah (sarankan pemisahan konteks menjadi multi-context modular) |
| **W2** | `performance.react-static-heavy-import` | R7 | React JSX/TSX AST | T2 | `LIKELY` | `warning` | Menengah (sarankan konversi ke `React.lazy()` + `<Suspense>`) |
| **W2** | `performance.react-redundant-function-memoization` | Baru | React JSX/TSX AST | T1 | `PROVEN` | `info` | **Tinggi** (lepaskan pembungkus `useCallback` menjadi fungsi biasa) |
| **W2** | `performance.react-derived-state-in-effect` | Baru | React JSX/TSX AST | T1 | `PROVEN` | `warning` | Menengah (sarankan kalkulasi inline langsung di badan fungsi render) |
| **W2** | `performance.react-unstable-hook-reference` | Baru | React TS/TSX AST | T1 | `PROVEN` | `warning` | Menengah (sarankan `useCallback`/`useMemo` pada objek kembalian custom hook) |
| **W3** | `performance.astro-unnecessary-client-directive` | Baru | Astro Compiler AST | T1 | `PROVEN` | `error` | **Tinggi** (hapus atribut direktif `client:*` untuk memulihkan Zero-JS) |
| **W3** | `performance.astro-island-boundary-overlap` | Baru | Astro + React AST | T2 | `LIKELY` | `warning` | Rendah (sarankan restrukturisasi komposisi pulau via Astro slots) |
| **W3** | `performance.astro-unoptimized-local-image` | R9 | Astro Template AST | T1 | `PROVEN` | `info` | **Tinggi** (migrasikan tag `<img>` lokal ke komponen `<Image />`) |
| **W3** | `performance.astro-over-prefetching` | R13 | Astro Template AST | T1 | `PROVEN` | `warning` | **Tinggi** (turunkan strategi prefetch ke `hover` pada rute sekunder) |
| **W4** | `performance.tailwind-dynamic-class-concatenation`| Baru | JSX/Astro AST | T1 | `PROVEN` | `error` | Menengah (sarankan pemetaan kamus objek eksplisit statis) |
| **W4** | `performance.tailwind-duplicate-arbitrary-rules` | Baru | JSX/Astro + CSS | T2 | `LIKELY` | `warning` | Menengah (sarankan konsolidasi aturan arbitrary yang ekuivalen ke `@theme`) |
| **W4** | `performance.tailwind-untracked-package-source` | Baru | CSS AST (TW4) | T1 | `PROVEN` | `error` | **Tinggi** (sisipkan direktif `@source "../packages/..."` di CSS root) |
| **W4** | `performance.tailwind-duplicate-utility-definition`| Baru | CSS AST (TW4) | T1 | `PROVEN` | `warning` | **Tinggi** (hapus `@utility` kustom yang menduplikasi utilitas core) |

---

## 4. Spesifikasi Detail & Kontrak Formal 16 Rule `performance.*`

### 4.1. `performance.react-inline-prop-memo` (Wave 1 - Migrasi Legacy R1)
- **Sumber Legacy:** `charites-legacy/perf-checker.ts` (Rule R1: `inline-prop-breaks-memo`).
- **Domain Parser:** React JSX/TSX AST + Component Memoization Graph.
- **Tier / Severity:** Tier 1 (Memoization Integrity Invariant) / `warning`.
- **Formal Contract:**
  - **Subject:** Pemanggilan komponen JSX yang dibungkus dalam `React.memo` atau `memo(...)`.
  - **Evidence:** Komponen menerima atribut prop berupa ekspresi inline instan: `ObjectExpression` (`prop={{ ... }}`), `ArrayExpression` (`prop={[ ... ]}`), `ArrowFunctionExpression` (`prop={() => ...}`), atau `JSXSpreadAttribute` objek baru.
  - **Predicate:** Komponen yang di-memoize **wajib menerima referensi prop yang stabil**, karena pembuatan objek/fungsi inline pada saat parent render menghasilkan alamat memori baru di setiap siklus, yang secara total membatalkan manfaat perbandingan dangkal (*shallow comparison*) `React.memo`.
  - **Confidence:** `PROVEN` (tipe ekspresi inline pada komponen teridentifikasi memo terbukti statis).
  - **Exceptions:** Komponen standar native DOM (`<div>`, `<button>`) atau komponen biasa yang tidak dibungkus `memo`.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Linter konvensional seperti ESLint (`react/jsx-no-bind`) hanya memeriksa apakah ada fungsi bind atau arrow function secara naif tanpa mengetahui apakah komponen penerima adalah komponen memoized. Akibatnya, ESLint membanjiri developer dengan warning palsu pada elemen native `<button onClick={() => ...}>`, atau justru membiarkan objek literal `data={{ id: 1 }}` lolos. Charites menelusuri definisi komponen penerima melalui Graf Relasional L3; jika terbukti dibungkus `React.memo()`, Charites secara presisi menandai alokasi objek/array/fungsi inline di call-site yang menghancurkan referensial identity dan membatalkan memoization.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Aturan bawaan ESLint `react/jsx-key` hanya memvalidasi ada atau tidaknya atribut `key` secara literal. Menulis `key={index}` membuat ESLint puas 100% dan menganggapnya kode yang wajar. Charites memeriksa apakah koleksi yang di-map adalah koleksi mutable atau memiliki interaksi dinamis (penyortiran, penyaringan, penghapusan item). Charites menangkap bahwa indeks numerik merusak algoritma rekonsiliasi VDOM React, memicu re-render seluruh subtree dan mempertahankan state internal yang salah pada kontrol input form.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  ESLint `react-hooks/exhaustive-deps` hanya memeriksa kelengkapan array dependensi, sama sekali buta terhadap operasi di dalam tubuh fungsi efek. Menulis `useEffect(() => { window.addEventListener("resize", onResize); }, [])` tanpa return statement dianggap 100% normal dan legal oleh ESLint. Charites menganalisis simpul panggilan di dalam efek (*resource acquisition symmetry*). Ketika mendeteksi pendaftaran listener, timer, observer, atau websocket tanpa fungsi cleanup simetris, Charites menandainya sebagai cacat struktural yang memicu kebocoran memori (*memory leak*) persisten.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Menggabungkan berbagai domain state ke dalam satu `createContext()` adalah pola yang lumrah bagi linter JavaScript standar karena secara sintaksis hanyalah objek JS biasa. Charites menganalisis frekuensi mutasi properti di dalam `value` konteks dan memetakan graf dependensi konsumennya. Charites mendeteksi ketika state interaktif berfrekuensi mutasi tinggi (seperti koordinat mouse atau scroll progress) digabungkan dengan state statis (seperti auth profile) dalam satu provider tunggal, menangkap potensi *re-render cascade* masif yang membebani CPU di seluruh pohon konsumen.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Pernyataan impor statis tingkat atas `import { Editor } from "monaco-editor"` di berkas rute halaman utama lolos pemeriksaan bundler dan linter tanpa komplain. Charites memeriksa pustaka berbobot besar (*heavy module signatures*) pada rute masuk awal (*initial route*) dan mendeteksi ketiadaan pemisahan kode asinkron, mewajibkan pemakaian `React.lazy()` dipadukan dengan `<Suspense>` guna memangkas ukuran bundle JavaScript awal (*critical JS payload*) yang harus diunduh dan diparsing peramban.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Membungkus fungsi dengan `useCallback` sering kali dianggap sebagai "best practice" umum oleh developer, dan tidak ada linter yang melarangnya. Charites menginspeksi konsumen fungsi tersebut: jika fungsi hanya diteruskan ke elemen native HTML (`<button onClick={...}>`), atau bila proyek mengaktifkan React Compiler yang sudah mengotomatisasi memoization secara global, Charites menandai `useCallback` tersebut sebagai pemborosan alokasi overhead memori hook yang tidak memberikan keuntungan performa apa pun.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Menulis `useEffect(() => { setFullName(first + " " + last); }, [first, last])` adalah kode yang sepenuhnya valid bagi linter hooks standar. Charites mengenali anti-pattern sinkronisasi state turunan ini (sesuai dok resmi React "You Might Not Need an Effect"). Charites mendeteksi bahwa kalkulasi turunan dari props/state seharusnya dihitung secara inline saat render pass berlangsung, mengeliminasi siklus render sekunder ganda yang membuang waktu main-thread.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Custom hook yang mengembalikan objek literal baru `{ data, refetch }` di setiap pemanggilan dianggap fungsi JavaScript biasa oleh ESLint. Charites menelusuri aliran nilai balik custom hook ke komponen pemanggil (*call-site data-flow*). Jika nilai balik tersebut diteruskan ke array dependensi hook hilir (`useEffect`, `useMemo`), Charites mendeteksi bahaya loop efek tak berujung (*infinite effect re-execution*) yang dipicu oleh instansiasi referensi memori baru di setiap render.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Menulis `<Header client:load />` pada komponen antarmuka yang isinya hanya teks dan gambar statis adalah sintaksis Astro resmi yang 100% valid. Linter biasa tidak memiliki visibilitas ke dalam isi komponen pulau untuk mengetahui apakah komponen tersebut benar-benar interaktif. Charites menginspeksi AST internal komponen; jika tidak menemukan interaktivitas (tanpa hooks, tanpa state, tanpa listener event), Charites menandai pelanggaran arsitektur Zero-JS yang mengirimkan JavaScript sia-sia ke browser klien.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Menyarangkan komponen interaktif ber-client directive di dalam pulau framework lain adalah sintaksis JSX/Astro yang valid dan lolos compiler. Charites menganalisis batas-batas pulau hidrasi (*island boundaries*), mendeteksi penyarangan pulau multi-framework atau transfer props reaktif lintas batas hidrasi yang memicu desinkronisasi state, kegagalan hidrasi parsial, dan overhead ganda runtime framework di memori peramban.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Tag HTML native `<img src="/images/banner.png" />` adalah markup standar yang sah di semua linter HTML/JSX. Charites mengenali bahwa berkas gambar berada di direktori lokal repositori dan merekomendasikan migrasi ke komponen `<Image />` bawaan `astro:assets`. Hal ini memastikan aset lokal melewati pipeline kompresi otomatis, inferensi dimensi resolusi, dan konversi ke format modern WebP/AVIF saat proses build produksi.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Menambahkan `data-astro-prefetch="viewport"` pada seluruh tautan navigasi internal di footer atau menu samping dianggap wajar oleh linter markup. Charites mengaudit penempatan strategi prefetch agresif pada tautan navigasi sekunder yang probabilitas kunjungannya rendah, mencegah pemborosan kuota data pengguna ponsel dan menghindari saturasi antrean unduhan jaringan yang seharusnya diprioritaskan untuk aset halaman aktif.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Penggabungan string dinamis via template literal seperti `` className={`text-${size} bg-${color}-500`} `` adalah JavaScript modern yang 100% legal bagi ESLint dan TypeScript (`tsc`). Charites memahami cara kerja compiler scanner statis Tailwind CSS v4 (Oxide engine). Karena scanner statis tidak mengeksekusi runtime JS saat build, penggabungan dinamis ini membuat utilitas CSS tidak terdeteksi sama sekali, mengakibatkan styling hilang total di CSS produksi.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Menulis kelas arbitrary seperti `w-[327px]` di berbagai berkas komponen adalah sintaksis Tailwind yang sah dan tidak pernah dipermasalahkan oleh linter CSS biasa. Charites mengaudit penggunaan nilai arbitrary berulang di seluruh proyek, menganjurkan pendaftaran nilai ke skala utilitas `@theme` atau konfigurasi token resmi guna mencegah pembengkakan duplikasi aturan CSS ekuivalen pada berkas stylesheet terkompilasi.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Mengimpor komponen dari paket monorepo internal (`@workspace/ui`) dan memakainya di template Astro/TSX adalah impor modul ES standar yang valid. Charites memeriksa berkas stylesheet root Tailwind CSS v4, mendeteksi ketiadaan direktif `@source` yang menyebabkan static scanner Oxide melewatkan pemindaian pada paket eksternal tersebut, sehingga seluruh utilitas Tailwind komponen paket monorepo hilang di hasil build.
- **AST Parser Unfair Advantage vs Linter Konvensional:**
  Deklarasi kustom `@utility` di dalam berkas CSS adalah fitur resmi Tailwind v4 yang valid secara sintaksis. Charites membandingkan deklarasi `@utility` kustom terhadap kamus utilitas bawaan Tailwind v4, mendeteksi jika pengembang menduplikasi properti dan nilai yang sudah disediakan secara native, mengeliminasi redundansi definisi CSS dan memelihara kebersihan codebase.
- **Autofix Feasibility:** Menengah (sarankan pembungkusan via `useMemo` atau `useCallback` di komponen induk).
- **Suspicious:**
  ```tsx
  // Pelanggaran: Komponen memoized UserCard selalu re-render karena opsi baru tiap render
  const UserCard = React.memo(({ user, config }: { user: User; config: Config }) => { ... });

  function Parent() {
    return <UserCard user={currentUser} config={{ theme: 'dark', compact: true }} />;
  }
  ```
- **Compliant:**
  ```tsx
  // Patuh: Menstabilkan referensi objek konfigurasi di luar atau via useMemo
  const USER_CONFIG = { theme: 'dark', compact: true } as const;

  function Parent() {
    return <UserCard user={currentUser} config={USER_CONFIG} />;
  }
  ```

---

### 4.2. `performance.react-index-as-key` (Wave 1 - Migrasi Legacy R2)
- **Sumber Legacy:** `charites-legacy/perf-checker.ts` (Rule R2: `no-index-key`).
- **Domain Parser:** React JSX/TSX AST + Collection Mutability Analyzer.
- **Tier / Severity:** Tier 1 (Reconciliation Identity Invariant) / `warning`.
- **Formal Contract:**
  - **Subject:** Pemanggilan metode `.map()` pada koleksi data yang merender elemen JSX berakar kunci (`key`).
  - **Evidence:** Nilai atribut `key` merujuk langsung ke parameter indeks perulangan (`index`, `idx`, `i`) **DAN** koleksi tersebut merupakan data dinamis yang berpotensi mengalami penambahan, penghapusan, penyaringan (*filter*), atau pengurutan ulang (*sort/reorder*).
  - **Predicate:** Elemen dalam koleksi mutable **dilarang menggunakan indeks array sebagai `key`**, karena rekonsiliasi React akan keliru mempertahankan state internal simpul anak saat urutan elemen berubah, memicu bug visual dan perenderan ulang DOM yang tidak perlu.
  - **Confidence:** `PROVEN` (korelasi parameter index callback terhadap atribut key pada koleksi mutable terdeteksi jelas).
  - **Exceptions:**
    1. Perulangan elemen kerangka statis read-only: `Array.from({ length: 5 }).map((_, i) => <Skeleton key={i} />)`.
    2. Koleksi konstan statis yang tidak pernah diubah (nama variabel UPPER_SNAKE_CASE).
- **Autofix Feasibility:** Menengah (sarankan pemakaian ID entitas stabil `item.id`).
- **Suspicious:**
  ```tsx
  // Pelanggaran: Menggunakan index pada list transaksi yang bisa difilter/diurutkan
  {transactions.map((tx, index) => (
    <TransactionRow key={index} data={tx} />
  ))}
  ```
- **Compliant:**
  ```tsx
  // Patuh: Menggunakan identifier unik permanen dari data
  {transactions.map((tx) => (
    <TransactionRow key={tx.id} data={tx} />
  ))}
  ```

---

### 4.3. `performance.react-effect-missing-cleanup` (Wave 1 - Migrasi Legacy R4)
- **Sumber Legacy:** `charites-legacy/perf-checker.ts` (Rule R4: `effect-no-cleanup`).
- **Domain Parser:** React JSX/TSX AST + Resource Acquisition Symmetry Detector.
- **Tier / Severity:** Tier 1 (Memory Leak Prevention Invariant) / `error`.
- **Formal Contract:**
  - **Subject:** Deklarasi hook `useEffect` atau `useLayoutEffect`.
  - **Evidence:** Badan hook melakukan akuisisi sumber daya eksternal persisten:
    $$\text{Acquired Resources} \in \{\text{addEventListener}, \text{setInterval}, \text{setTimeout}, \text{ResizeObserver}, \text{IntersectionObserver}, \text{MutationObserver}, \text{WebSocket}, \text{AbortController}\}$$
    namun fungsi callback hook **tidak mengembalikan (*does not return*) fungsi pembersih simetris (*cleanup release function*)**.
  - **Predicate:** Efek samping yang mengalokasikan langganan peramban atau timer persisten **wajib mengembalikan fungsi cleanup**, guna mencegah kebocoran memori (*memory leak*) dan eksekusi callback pada komponen yang telah dilepas (*unmounted*).
  - **Confidence:** `PROVEN` (keberadaan pemanggilan subscription tanpa return statement terbukti deterministik).
  - **Exceptions:** Efek murni *fire-and-forget* tanpa langganan persisten (misal: penulisan `document.title = 'Judul'`).
- **Autofix Feasibility:** **Tinggi**. Menyisipkan kerangka fungsi pengembalian pembersih.
- **Suspicious:**
  ```tsx
  // Pelanggaran: Menambahkan listener window tanpa fungsi cleanup saat unmount
  useEffect(() => {
    const onResize = () => setWidth(window.innerWidth);
    window.addEventListener('resize', onResize);
  }, []);
  ```
- **Compliant:**
  ```tsx
  // Patuh: Membersihkan listener saat komponen unmount atau dependensi berubah
  useEffect(() => {
    const onResize = () => setWidth(window.innerWidth);
    window.addEventListener('resize', onResize);
    return () => window.removeEventListener('resize', onResize);
  }, []);
  ```

---

### 4.4. `performance.react-context-domain-coupling` (Wave 1 - Refined Legacy R5)
- **Sumber Legacy:** `charites-legacy/perf-checker.ts` (Rule R5: `god-context`).
- **Domain Parser:** React JSX/TSX AST + Context Consumer Dependency Graph.
- **Tier / Severity:** Tier 2 (Context Granularity Architecture) / `warning`.
- **Formal Contract:**
  - **Subject:** Deklarasi penyedia konteks `<Context.Provider>`.
  - **Evidence:** Objek nilai konteks (`value={{ ... }}`) menggabungkan domain state yang memiliki frekuensi pembaruan (*mutation frequency*) sangat berbeda (misal: state navigasi berfrekuensi tinggi digabung dengan data profil pengguna statis berfrekuensi rendah) dalam satu boundary provider yang dikonsumsi oleh puluhan komponen berbeda.
  - **Predicate:** State aplikasi yang memiliki frekuensi pembaruan berbeda **dilarang digabungkan ke dalam satu Context terpadu dengan basis konsumen luas**, karena setiap pembaruan pada state yang sering berubah akan memaksa seluruh komponen konsumen statis untuk melakukan re-render penuh.
  - **Confidence:** `LIKELY` (analisis graf mendeteksi korelasi dependensi konsumen terhadap domain nilai konteks).
  - **Exceptions:** Konteks konfigurasi statis murni yang nilainya tidak pernah bermutasi setelah inisialisasi aplikasi.
- **Autofix Feasibility:** Rendah (membutuhkan dekomposisi arsitektur konteks menjadi multi-provider terpisah).
- **Suspicious:**
  ```tsx
  // Pelanggaran: State kursor/hover (high-frequency) digabung dengan AuthUser (low-frequency)
  <GlobalContext.Provider value={{ currentUser, theme, mousePosition, scrollProgress }}>
    {children}
  </GlobalContext.Provider>
  ```
- **Compliant:**
  ```tsx
  // Patuh: Memisahkan konteks statis dari konteks interaksi dinamis
  <AuthContext.Provider value={currentUser}>
    <InteractionContext.Provider value={{ mousePosition, scrollProgress }}>
      {children}
    </InteractionContext.Provider>
  </AuthContext.Provider>
  ```

---

### 4.5. `performance.react-static-heavy-import` (Wave 2 - Migrasi Legacy R7)
- **Sumber Legacy:** `charites-legacy/perf-checker.ts` (Rule R7: `no-lazy-heavy-route`).
- **Domain Parser:** React JSX/TSX AST + Module Graph & Estimated Weight Analyzer.
- **Tier / Severity:** Tier 2 (Bundle Splitting Optimization) / `warning`.
- **Formal Contract:**
  - **Subject:** Pernyataan impor statis di tingkat atas (*top-level static imports*) pada berkas komponen halaman atau tata letak.
  - **Evidence:** Impor statis merujuk langsung ke pustaka eksternal atau modul komponen lokal yang diketahui berbobot muatan sangat besar:
    $$\text{Heavy Module Signatures} \in \{\text{monaco-editor}, \text{chart.js}, \text{echarts}, \text{quill}, \text{pdfjs}, \text{three}, \text{canvas}, \text{datatables}\}$$
    pada rute awal (*initial route*) tanpa memanfaatkan pemisahan kode (*code splitting*).
  - **Predicate:** Komponen atau pustaka non-kritis yang berukuran besar **wajib dimuat secara asinkron menggunakan `React.lazy()` dipadukan dengan `<Suspense>`**, guna mencegah pembengkakan bundle JavaScript awal (*initial bundle size bloat*).
  - **Confidence:** `LIKELY` (pola nama modul dan pustaka grafis berat terdeteksi presisi).
  - **Exceptions:** Modul yang diperlukan secara instan untuk merender antarmuka di atas pelipatan (*above-the-fold*).
- **Autofix Feasibility:** Menengah (sarankan transformasi ke `const HeavyComp = React.lazy(...)`).
- **Suspicious:**
  ```tsx
  // Pelanggaran: Mengimpor modul visualisasi grafik berat secara statis di halaman awal
  import { AnalyticalChart } from '@/components/AnalyticalChart';

  export function Dashboard() {
    return <AnalyticalChart data={stats} />;
  }
  ```
- **Compliant:**
  ```tsx
  // Patuh: Memisahkan bundel via React.lazy dan menampilkan loading fallback
  import { Suspense, lazy } from 'react';
  const AnalyticalChart = lazy(() => import('@/components/AnalyticalChart'));

  export function Dashboard() {
    return (
      <Suspense fallback={<ChartSkeleton />}>
        <AnalyticalChart data={stats} />
      </Suspense>
    );
  }
  ```

---

### 4.6. `performance.react-redundant-function-memoization` (Wave 2 - Rewritten & Compiler-Aware)
- **Sumber Legacy:** Konsep baru (Hook Overhead Elimination & React Compiler Readiness).
- **Domain Parser:** React JSX/TSX AST + Downstream Identity Consumer Inspector.
- **Tier / Severity:** Tier 1 (Hook Allocation Economy) / `info` (advisory).
- **Formal Contract:**
  - **Subject:** Penggunaan hook `useCallback` atau `useMemo`.
  - **Evidence:** Fungsi dibungkus dalam `useCallback` namun hanya diteruskan langsung ke elemen native HTML DOM (misal `<button onClick={cb}>`) dan tidak pernah digunakan sebagai array dependensi hook lain atau diteruskan ke komponen `React.memo` **DAN** project tidak menggunakan React Compiler (manual mode).
  - **Predicate:** Penggunaan `useCallback` **sebaiknya dihindari jika fungsi tidak memiliki konsumen yang sensitif terhadap identitas referensi**, karena alokasi closure internal hook dan array dependensi justru memakan lebih banyak memori dan siklus CPU dibandingkan pembuatan fungsi inline biasa.
  - **Confidence:** `PROVEN` (tujuan passing fungsi terbukti langsung ke simpul native JSX tanpa perantara memo).
  - **Exceptions:**
    1. Project yang mengaktifkan **React Compiler** (aturan otomatis di-suppress).
    2. Fungsi yang dijadikan dependensi pada `useEffect` atau diteruskan ke komponen memoized anak.
- **Autofix Feasibility:** **Tinggi**. Menghapus pembungkus `useCallback` menjadi fungsi reguler.
- **Suspicious:**
  ```tsx
  // Advisory: Membungkus handler tombol native dengan useCallback adalah overhead sia-sia
  const handleClick = useCallback(() => {
    setOpen(true);
  }, []);

  return <button onClick={handleClick}>Buka Modal</button>;
  ```
- **Compliant:**
  ```tsx
  // Patuh: Gunakan fungsi reguler untuk elemen DOM biasa
  const handleClick = () => {
    setOpen(true);
  };

  return <button onClick={handleClick}>Buka Modal</button>;
  ```

---

### 4.7. `performance.react-derived-state-in-effect` (Wave 2 - Baru)
- **Sumber Legacy:** Konsep baru (Elimination of Cascading Secondary Renders).
- **Domain Parser:** React JSX/TSX AST.
- **Tier / Severity:** Tier 1 (Render Pass Minimization) / `warning`.
- **Formal Contract:**
  - **Subject:** Deklarasi hook `useEffect`.
  - **Evidence:** Hook `useEffect` hanya berfungsi menghitung nilai turunan (*derived state*) dengan memanggil setter state lokal (`setState`) setiap kali props atau state lain berubah:
    $$\text{Pattern: } \text{useEffect}(() \Rightarrow \{\text{setFullName}(\text{first} + ' ' + \text{last})\}, [\text{first}, \text{last}])$$
  - **Predicate:** Nilai turunan dari props atau state yang sudah ada **dilarang disinkronkan melalui `useEffect`**, karena pola ini memaksa React menjalankan siklus render sekunder (*cascading render pass*); nilai tersebut wajib dihitung secara sinkron selama fase render.
  - **Confidence:** `PROVEN` (setter state murni dari variabel dependensi di dalam efek terdeteksi statis).
  - **Exceptions:** Operasi asinkron (misal: pengambilan data dari API luar saat prop ID berubah).
- **Autofix Feasibility:** Menengah (sarankan penghapusan `useEffect` dan pembuatan variabel kalkulasi langsung).
- **Suspicious:**
  ```tsx
  // Pelanggaran: Sinkronisasi nama lengkap via useEffect memicu render sekunder ganda
  const [fullName, setFullName] = useState('');
  useEffect(() => {
    setFullName(`${firstName} ${lastName}`);
  }, [firstName, lastName]);
  ```
- **Compliant:**
  ```tsx
  // Patuh: Dihitung secara sinkron dalam satu kali fase render
  const fullName = `${firstName} ${lastName}`;
  ```

---

### 4.8. `performance.react-unstable-hook-reference` (Wave 2 - Baru)
- **Sumber Legacy:** Konsep baru (Custom Hook Reference Stability Invariant).
- **Domain Parser:** React TS/TSX AST + Downstream Consumer Graph.
- **Tier / Severity:** Tier 1 (Contract Integrity) / `warning`.
- **Formal Contract:**
  - **Subject:** Deklarasi custom hook (fungsi berawalan `use...`).
  - **Evidence:** Custom hook mengembalikan objek atau fungsi helper baru di setiap pemanggilan (`return { data, refetch: () => ... }`) **DAN** nilai return tersebut dikonsumsi oleh komponen hilir sebagai dependensi `useEffect` atau prop pada komponen ter-memoize.
  - **Predicate:** Custom hook yang mengekspos fungsi atau objek kembalian **wajib menstabilkan referensi memori fungsi tersebut (via `useCallback`/`useMemo`)**, agar komponen konsumen dapat menyertakannya secara aman di dalam array dependensi tanpa memicu perulangan re-render tak terkendali.
  - **Confidence:** `PROVEN` (instansiasi objek baru di return statement custom hook terbukti statis).
  - **Exceptions:** Hook yang hanya mengembalikan nilai primitif (`boolean`, `string`, `number`).
- **Autofix Feasibility:** Menengah (sarankan pembungkusan fungsi return dalam `useCallback`).
- **Suspicious:**
  ```tsx
  // Pelanggaran: refetch adalah fungsi baru di setiap render komponen konsumen
  export function useProfile(userId: string) {
    const [data, setData] = useState(null);
    const refetch = () => { fetchProfile(userId).then(setData); };
    return { data, refetch };
  }
  ```
- **Compliant:**
  ```tsx
  // Patuh: Menstabilkan fungsi pembantu menggunakan useCallback
  export function useProfile(userId: string) {
    const [data, setData] = useState(null);
    const refetch = useCallback(() => {
      fetchProfile(userId).then(setData);
    }, [userId]);
    return { data, refetch };
  }
  ```

---

### 4.9. `performance.astro-unnecessary-client-directive` (Wave 3 - Baru)
- **Sumber Legacy:** Konsep baru (Astro Zero-JS Principle Enforcement).
- **Domain Parser:** Astro Compiler AST + Component Interactivity Inference Engine.
- **Tier / Severity:** Tier 1 (Zero-JS Paradigm Invariant) / `error`.
- **Formal Contract:**
  - **Subject:** Elemen pulau (*island*) komponen di dalam template Astro.
  - **Evidence:** Komponen diberi direktif hidrasi peramban (`client:load`, `client:idle`, `client:visible`), padahal analisis inferensi membuktikan bahwa komponen tersebut **murni statis**: tidak memiliki state (`useState`), tidak memiliki efek samping (`useEffect`), tidak memiliki event listener (`onClick`, `onChange`), dan tidak memanggil Browser Web APIs.
  - **Predicate:** Komponen antarmuka yang tidak interaktif **dilarang menyertakan direktif hidrasi `client:*`**, karena direktif tersebut memaksa bundler mengirimkan kode runtime framework dan bundel komponen ke peramban, melanggar arsitektur Zero-JS bawaan Astro.
  - **Confidence:** `PROVEN` (ketiadaan interaktivitas pada komponen terbukti deterministik).
  - **Exceptions:** Komponen yang memanfaatkan animasi CSS murni tetapi membutuhkan ref DOM dinamis khusus.
- **Autofix Feasibility:** **Tinggi**. Menghapus atribut direktif hidrasi dari markup Astro.
- **Suspicious:**
  ```astro
  ---
  // Pelanggaran: HeaderFooter statis dipaksa terhidrasi ke browser mengirim 40KB React runtime
  import HeaderStatic from '../components/HeaderStatic.tsx';
  ---
  <HeaderStatic client:load title="Selamat Datang" />
  ```
- **Compliant:**
  ```astro
  ---
  // Patuh: Dirender sebagai pure static HTML tanpa mengirim satu byte pun JavaScript ke client
  import HeaderStatic from '../components/HeaderStatic.tsx';
  ---
  <HeaderStatic title="Selamat Datang" />
  ```

---

### 4.10. `performance.astro-island-boundary-overlap` (Wave 3 - Rewritten)
- **Sumber Legacy:** Konsep baru (Island Boundary Overlap & Multi-Framework Runtime Hygiene).
- **Domain Parser:** Astro Compiler AST + Island Hydration Boundary Traversal.
- **Tier / Severity:** Tier 2 (Hydration Correctness & Runtime Overhead) / `warning`.
- **Formal Contract:**
  - **Subject:** Hierarki pulau (*islands*) interaktif di dalam proyek Astro.
  - **Evidence:** Sebuah komponen pulau interaktif (`client:*`) menyarangkan pulau client lain yang terhidrasi secara terpisah dengan framework berbeda (misal: React Island di dalam Vue Island) atau mengoper props reaktif lintas batas hidrasi independen yang memicu **state isolation breakdown** dan duplikasi runtime framework.
  - **Predicate:** Pulau hidrasi Astro **wajib menjaga isolasi boundary hidrasi yang jelas**; hindari nesting pulau interaktif multi-framework secara langsung. Komposisikan antarmuka menggunakan Astro Slots (`<slot />`) agar boundary hidrasi tetap terisolasi dan tidak memicu konflik listener atau duplikasi runtime framework ganda di klien.
  - **Confidence:** `LIKELY` (struktur penyarangan pulau dengan runtime terpisah terdeteksi dari AST graph).
  - **Exceptions:** Komponen murni satu framework yang berada dalam satu pohon virtual DOM internal pulau yang sama.
- **Autofix Feasibility:** Rendah (memerlukan refactoring komposisi antarmuka via Astro slot).
- **Suspicious:**
  ```astro
  <!-- Pelanggaran: Penyarangan pulau multi-framework memicu konflik hidrasi & runtime ganda -->
  <ReactDashboardContainer client:load>
    <VueChartWidget client:idle />
  </ReactDashboardContainer>
  ```
- **Compliant:**
  ```astro
  <!-- Patuh: Memanfaatkan Astro slot untuk menjaga boundary hidrasi tetap terisolasi -->
  <ReactDashboardContainer client:load>
    <div slot="chart-slot">
      <VueChartWidget client:idle />
    </div>
  </ReactDashboardContainer>
  ```

---

### 4.11. `performance.astro-unoptimized-local-image` (Wave 3 - Refined Legacy R9)
- **Sumber Legacy:** `charites-legacy/perf-checker.ts` (Rule R9: `no-astro-image`).
- **Domain Parser:** Astro Template AST.
- **Tier / Severity:** Tier 1 (Asset Pipeline Integration) / `info` (advisory).
- **Formal Contract:**
  - **Subject:** Elemen `<img>` di dalam template berkas `.astro`.
  - **Evidence:** Tag `<img>` merujuk langsung ke berkas gambar lokal di direktori `src/` atau `assets/` (misal: `src="../assets/banner.png"`) tanpa menggunakan komponen bawaan `<Image />` atau `<Picture />` dari modul `astro:assets`.
  - **Predicate:** Berkas gambar lokal **disarankan diproses melalui komponen `<Image />` dari `astro:assets`**, agar bundler Astro dapat mengoptimasi resolusi gambar, mengonversi otomatis ke format modern WebP/AVIF, dan menginjeksi atribut dimensi `width`/`height` secara otomatis saat proses build.
  - **Confidence:** `PROVEN` (penggunaan tag img mentah pada path lokal terbukti dari AST).
  - **Exceptions:** Gambar vektor SVG, gambar eksternal yang URL-nya baru diketahui saat runtime dinamis, atau gambar di dalam direktori `public/` yang sengaja disajikan apa adanya tanpa kompresi build.
- **Autofix Feasibility:** **Tinggi**. Mengimpor dan mengonversi tag `<img>` ke `<Image />`.
- **Suspicious:**
  ```astro
  <!-- Advisory: Tag img mentah melewatkan kompresi format dan ekstraksi dimensi otomatis -->
  <img src="../assets/product-hero.png" alt="Produk Baru" />
  ```
- **Compliant:**
  ```astro
  ---
  // Patuh: Memanfaatkan pipeline optimasi build bawaan Astro
  import { Image } from 'astro:assets';
  import productImg from '../assets/product-hero.png';
  ---
  <Image src={productImg} alt="Produk Baru" />
  ```

---

### 4.12. `performance.astro-over-prefetching` (Wave 3 - Migrasi Legacy R13)
- **Sumber Legacy:** `charites-legacy/perf-checker.ts` (Rule R13: `too-many-prefetch`).
- **Domain Parser:** Astro Template AST + Route Probability Graph.
- **Tier / Severity:** Tier 1 (Bandwidth Economy) / `warning`.
- **Formal Contract:**
  - **Subject:** Tautan navigasi (`<a>`) di dalam template halaman Astro.
  - **Evidence:** Halaman memuat sejumlah besar tautan dengan strategi prefetch agresif (`data-astro-prefetch="viewport"` atau `"load"`) pada tautan yang memiliki probabilitas navigasi rendah (seperti tautan footer, halaman syarat & ketentuan, kebijakan privasi).
  - **Predicate:** Prefetch agresif berbasis `viewport` atau `load` **hanya boleh diterapkan pada rute konversi prioritas tinggi**; tautan sekunder wajib menggunakan strategi prefetch pasif (`data-astro-prefetch="hover"` atau `tap`) untuk mencegah pemborosan kuota data seluler dan perebutan bandwidth pada jaringan lambat.
  - **Confidence:** `PROVEN` (atribut prefetch agresif pada tautan sekunder terdeteksi jelas).
  - **Exceptions:** Tautan paginasi berikutnya (*next page*) pada alur artikel multi-halaman.
- **Autofix Feasibility:** **Tinggi**. Mengubah nilai atribut `data-astro-prefetch` menjadi `"hover"`.
- **Suspicious:**
  ```astro
  <!-- Pelanggaran: Memasang prefetch agresif pada tautan footer non-kritis -->
  <footer>
    <a href="/terms" data-astro-prefetch="viewport">Syarat & Ketentuan</a>
    <a href="/privacy" data-astro-prefetch="viewport">Kebijakan Privasi</a>
  </footer>
  ```
- **Compliant:**
  ```astro
  <!-- Patuh: Prefetch pasif saat kursor menyentuh tautan (hover) -->
  <footer>
    <a href="/terms" data-astro-prefetch="hover">Syarat & Ketentuan</a>
    <a href="/privacy" data-astro-prefetch="hover">Kebijakan Privasi</a>
  </footer>
  ```

---

### 4.13. `performance.tailwind-dynamic-class-concatenation` (Wave 4 - Baru)
- **Sumber Legacy:** Dimodernisasi dari R14/R15 (Tailwind v4 Static Scanner Integrity).
- **Domain Parser:** React JSX/TSX AST + Astro Template AST.
- **Tier / Severity:** Tier 1 (Static Scanner Compatibility) / `error`.
- **Formal Contract:**
  - **Subject:** Penetapan atribut nama kelas (`className` atau `class`) pada komponen atau elemen markup.
  - **Evidence:** Nama kelas dikonstruksi melalui interpolasi string dinamis atau penggabungan parsial ekspresi runtime:
    $$\text{Pattern: } \text{className}=\{`\text{text-}\{\text{color}\}-600`\} \quad \text{atau} \quad \text{class}=\{`\text{bg-}\{\text{variant}\}`\}$$
  - **Predicate:** Kelas utilitas Tailwind CSS **wajib ditulis sebagai string literal utuh dan statis**, karena mesin kompilasi Tailwind CSS v4 (Oxide scanner) mengekstrak kelas melalui pemindaian statis cepat pada kode sumber; penggabungan string dinamis membuat kelas tersebut tidak terdeteksi sehingga aturan CSS tidak akan pernah dihasilkan di produksi.
  - **Confidence:** `PROVEN` (keberadaan template literal dinamis pada prefiks kelas Tailwind terbukti deterministik).
  - **Exceptions:** Pemilihan kelas utuh menggunakan conditional operator (ternary) atau library utilitas seperti `clsx` / `cn('bg-red-500', isLarge && 'p-8')`.
- **Autofix Feasibility:** Menengah (sarankan pemetaan kamus objek eksplisit).
- **Suspicious:**
  ```tsx
  // Pelanggaran: Scanner statis Tailwind v4 tidak bisa mengekstrak kelas ini
  function Badge({ color }: { color: 'red' | 'blue' }) {
    return <span className={`bg-${color}-100 text-${color}-800`}>Status</span>;
  }
  ```
- **Compliant:**
  ```tsx
  // Patuh: Menuliskan nama kelas secara utuh dalam kamus statis
  const COLOR_MAP = {
    red: 'bg-red-100 text-red-800',
    blue: 'bg-blue-100 text-blue-800',
  } as const;

  function Badge({ color }: { color: 'red' | 'blue' }) {
    return <span className={COLOR_MAP[color]}>Status</span>;
  }
  ```

---

### 4.14. `performance.tailwind-duplicate-arbitrary-rules` (Wave 4 - Rewritten)
- **Sumber Legacy:** Konsep baru (Tailwind v4 CSS Declaration Deduplication).
- **Domain Parser:** JSX/Astro AST + CSS Output Equivalence Resolver.
- **Tier / Severity:** Tier 2 (Compiled CSS Payload Hygiene) / `warning`.
- **Formal Contract:**
  - **Subject:** Penggunaan utilitas nilai sembarang (*arbitrary values*) dalam atribut kelas markup.
  - **Evidence:** Proyek mendefinisikan beberapa nilai arbitrary yang menghasilkan aturan CSS deklarasi identik (misal: `mt-[16px]` bercampur dengan `mt-4`, atau `p-[1rem]` bercampur dengan `p-4`) di berbagai berkas komponen.
  - **Predicate:** Pengembang **sebaiknya menggunakan utilitas skala bawaan Tailwind v4** daripada menuliskan nilai arbitrary ad-hoc yang ekuivalen, guna mencegah duplikasi deklarasi CSS di dalam stylesheet produksi.
  - **Confidence:** `LIKELY` (ekuivalensi nilai arbitrary terhadap utilitas core terbukti via resolver).
  - **Exceptions:** Nilai arbitrary mikro yang memang tidak memiliki ekuivalensi pada skala desain Tailwind v4 (misal `top-[13px]`).
- **Autofix Feasibility:** **Tinggi**. Mengganti nilai arbitrary dengan kelas skala standar Tailwind v4.
- **Suspicious:**
  ```tsx
  // Pelanggaran: Menggunakan arbitrary p-[16px] yang menduplikasi utilitas core p-4
  <div className="p-[16px] mt-[1rem]">...</div>
  ```
- **Compliant:**
  ```tsx
  // Patuh: Menggunakan kelas skala core yang sudah ada di runtime CSS
  <div className="p-4 mt-4">...</div>
  ```

---

### 4.15. `performance.tailwind-untracked-package-source` (Wave 4 - Dimodernisasi dari Legacy R14/R15)
- **Sumber Legacy:** Dimodernisasi dari `charites-legacy/perf-checker.ts` (Rule R14/R15: `tailwind-content-missing`).
- **Domain Parser:** CSS AST (Tailwind v4 Root Stylesheet) + Monorepo Package Graph.
- **Tier / Severity:** Tier 1 (Compiler Discovery Boundary) / `error`.
- **Formal Contract:**
  - **Subject:** Berkas stylesheet root Tailwind CSS v4 (`src/styles/global.css`).
  - **Evidence:** Proyek mengimpor komponen antarmuka dari paket monorepo eksternal (misal: `import { Button } from '@repo/ui'`) **namun berkas CSS root tidak menyertakan direktif pemindaian `@source`** untuk direktori paket tersebut:
    $$\text{Missing Directive: } @\text{source } "\text{../../packages/ui}";$$
  - **Predicate:** Dalam arsitektur Tailwind CSS v4, pustaka komponen lokal atau paket workspace monorepo **wajib didaftarkan secara eksplisit melalui direktif `@source` di dalam stylesheet utama**, agar compiler scanner dapat menemukan dan menyertakan seluruh kelas utilitas yang digunakan oleh paket tersebut ke dalam CSS produksi.
  - **Confidence:** `PROVEN` (korelasi package import terhadap ketiadaan direktif `@source` di CSS terbukti statis).
  - **Exceptions:** Komponen yang sudah didistribusikan dalam bentuk CSS pra-kompilasi (*pre-compiled standalone CSS*).
- **Autofix Feasibility:** **Tinggi**. Menyisipkan baris `@source "../../packages/..."` pada berkas CSS root.
- **Suspicious:**
  ```css
  /* Pelanggaran: global.css tidak memindai paket monorepo eksternal @repo/ui */
  @import "tailwindcss";
  ```
- **Compliant:**
  ```css
  /* Patuh: Mendaftarkan path paket eksternal agar discan oleh compiler v4 */
  @import "tailwindcss";
  @source "../../packages/ui";
  ```

---

### 4.16. `performance.tailwind-duplicate-utility-definition` (Wave 4 - Baru)
- **Sumber Legacy:** Konsep baru (Tailwind v4 `@utility` Hygiene).
- **Domain Parser:** CSS AST (PostCSS) + Native Utility Dictionary.
- **Tier / Severity:** Tier 1 (Utility Duplication Elimination) / `warning`.
- **Formal Contract:**
  - **Subject:** Deklarasi direktif `@utility` kustom di dalam berkas CSS.
  - **Evidence:** Blok `@utility <name>` mendeklarasikan aturan CSS yang properti dan nilainya sudah disediakan secara native oleh utilitas bawaan Tailwind v4 (misal: `@utility flex-center { display: flex; justify-content: center; align-items: center; }`).
  - **Predicate:** Deklarasi `@utility` kustom **sebaiknya tidak menduplikasi utilitas core bawaan**, melainkan dikomposisikan langsung di markup atau hanya digunakan untuk properti CSS mutakhir yang belum didukung core Tailwind, guna mencegah duplikasi byte output CSS yang tidak perlu.
  - **Confidence:** `PROVEN` (pencocokan properti utilitas terhadap kamus core utility Tailwind v4).
  - **Exceptions:** Utilitas kustom yang menggabungkan vendor-prefixes kompleks atau properti CSS experimental.
- **Autofix Feasibility:** **Tinggi**. Menghapus deklarasi utilitas duplikat dan menyarankan kombinasi kelas core.
- **Suspicious:**
  ```css
  /* Pelanggaran: Menduplikasi utilitas native flexbox bawaan Tailwind */
  @utility center-flex {
    display: flex;
    align-items: center;
  }
  ```
- **Compliant:**
  ```tsx
  /* Patuh: Gunakan kombinasi utilitas core native langsung di markup */
  <div className="flex items-center">...</div>
  ```

---

## 5. Rubrik Keparahan, Matriks Keyakinan & Kelayakan Autofix

### 5.1. Skala Keparahan (*Severity Scale*)
```text
┌──────────────┬───────────────────────────────┬──────────────────────────────────────────┐
│   Severity   │ Kriteria Penentuan            │ Contoh Kasus                             │
├──────────────┼───────────────────────────────┼──────────────────────────────────────────┤
│ error        │ Pelanggaran deterministik     │ • performance.react-effect-missing-clean │
│              │ yang memicu memory leak,      │ • performance.astro-unnecessary-client-d │
│              │ rusaknya scanner kompilasi,   │ • performance.tailwind-dynamic-class-con │
│              │ atau hilangnya styling build. │ • performance.tailwind-untracked-package │
├──────────────┼───────────────────────────────┼──────────────────────────────────────────┤
│ warning      │ Anti-pola yang memicu         │ • performance.react-inline-prop-memo     │
│              │ pembengkakan bundle,          │ • performance.react-index-as-key         │
│              │ re-render VDOM sia-sia, atau  │ • performance.astro-island-boundary-over │
│              │ pemborosan alokasi memori.    │ • performance.tailwind-duplicate-utility │
├──────────────┼───────────────────────────────┼──────────────────────────────────────────┤
│ info         │ Saran efisiensi kode mikro    │ • performance.react-redundant-function-m │
│              │ atau adopsi pipeline modern.  │ • performance.astro-unoptimized-local-im │
└──────────────┴───────────────────────────────┴──────────────────────────────────────────┘
```

### 5.2. Kebijakan Transaksional Autofix Aman
1. **Pencegahan Mutasi Semantik Berbahaya:**
   - Autofix `performance.astro-unnecessary-client-directive` hanya menghapus direktif hidrasi jika dan hanya jika komponen terbukti 100% statis (tidak ada hook, event, atau interaktivitas).
2. **Preservasi String Literal Tailwind v4:**
   - Autofix `performance.tailwind-untracked-package-source` secara otomatis menambahkan path direktif `@source` ke berkas `global.css` tanpa mengubah urutan impor `@import "tailwindcss"`.
3. **Penyisipan Skeleton Cleanup Effect:**
   - Autofix `performance.react-effect-missing-cleanup` menyisipkan blok return pembersih (`return () => { ... }`) dengan nama fungsi pembersih yang sesuai (`removeEventListener` atau `disconnect`).

---

## 6. Roadmap Implementasi 4 Wave

Penerapan engine static analyzer Go di `internal/rules/performance/` (dengan alias registry `perf`) dijadwalkan secara bertahap:

1. **Wave 1 (React Reconciliation & Memory Lifecycle):**
   - `performance.react-inline-prop-memo`
   - `performance.react-index-as-key`
   - `performance.react-effect-missing-cleanup`
   - `performance.react-context-domain-coupling`
2. **Wave 2 (React Code Splitting & Compilation Hygiene):**
   - `performance.react-static-heavy-import`
   - `performance.react-redundant-function-memoization`
   - `performance.react-derived-state-in-effect`
   - `performance.react-unstable-hook-reference`
3. **Wave 3 (Astro Architecture & Zero-JS Island Pipeline):**
   - `performance.astro-unnecessary-client-directive`
   - `performance.astro-island-boundary-overlap`
   - `performance.astro-unoptimized-local-image`
   - `performance.astro-over-prefetching`
4. **Wave 4 (Tailwind CSS v4 & Build-Time Stylesheet Hygiene):**
   - `performance.tailwind-dynamic-class-concatenation`
   - `performance.tailwind-duplicate-arbitrary-rules`
   - `performance.tailwind-untracked-package-source`
   - `performance.tailwind-duplicate-utility-definition`
