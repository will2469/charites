# EXPANSION-BATCH-06: Progressive Web App & Offline Standards (`pwa.*`)
> **Kode Dokumen:** `SPEC-EXP-06-PWA`
> **Kategori:** `pwa`
> **Pilar:** `01-SPEC` (WHAT - Spesifikasi Perilaku & Kontrak Rule)
> **Status:** Active Expansion Specification (10 Aturan Terkurasi Bebas Redundansi Linter Standar)
> **Standar Rujukan:**
> - W3C Web App Manifest Specification
> - W3C Service Workers Nightly & Cache Storage API
> - Google Web App Installability Criteria & Maskable Icons Standard
> - Apple Safari Web Content Guide (Configuring Web Applications)
> - W3C Secure Contexts (Mixed Content Mitigation)

---

## 1. Ikhtisar Kategori `pwa` (10 Aturan Non-Redundan)

> **Prinsip Eliminasi Redundansi:** Linter kode konvensional (ESLint, `tsc`, Prettier, Stylelint) **buta total terhadap berkas konfigurasi PWA (`manifest.json`), siklus hidup Service Worker (`sw.js`), dan keterkaitan meta tag Apple**. Kategori `pwa` Charites berfokus pada validasi statis lintas-file (*cross-file static analysis*) yang menghubungkan dokumen HTML/JSX dengan file manifest dan worker script tanpa perlu menyalakan runtime browser.

```mermaid
flowchart TD
    subgraph W1 ["Wave 1: Web App Manifest & Branding (4 Rules)"]
        P1["pwa.manifest-required-fields-missing (Field wajib: name, start_url, display)"]
        P2["pwa.icon-maskable-missing (Ikon adaptif Android launcher)"]
        P3["pwa.manifest-missing (Keberadaan tag link rel manifest)"]
        P4["pwa.start-url-inconsistency (Validitas cakupan start_url terhadap base path)"]
    end

    subgraph W2 ["Wave 2: Apple Standalone & Security (2 Rules)"]
        P5["pwa.apple-meta-missing (Meta tag apple-mobile-web-app & touch icon)"]
        P6["pwa.insecure-context-resource (Pencegahan mixed content http:// di secure context)"]
    end

    subgraph W3 ["Wave 3: Service Worker Lifecycle & Offline Cache (4 Rules)"]
        P7["pwa.service-worker-no-offline-fallback (Fallback fetch handler saat offline)"]
        P8["pwa.service-worker-missing (Registrasi SW pada project PWA)"]
        P9["pwa.service-worker-registration (Registrasi aman & error handling)"]
        P10["pwa.pwa-cache-runtime-api-risk (Cegah akses DOM/window di Worker)"]
    end
```

---

## 2. Spesifikasi Detail Rule Wave 1: Web App Manifest & Branding (4 Rules)

---

### 2.1. `pwa.manifest-required-fields-missing`
- **Design Rationale:** W3C Web App Manifest Specification & Google Chrome Web App Installability Criteria.
- **Konteks Realitas Mobile:**
  Agar sebuah web app memenuhi syarat instalasi PWA (*Add to Home Screen*) pada Android dan iOS, berkas manifest harus mendefinisikan field wajib minimum: `name` (atau `short_name`), `start_url`, `display` (seperti `standalone` atau `fullscreen`), dan array `icons` yang memuat minimal satu ikon. Jika salah satu field hilang, browser menolak memunculkan prompt instalasi atau menampilkan aplikasi dengan nama placeholder kosong dan ikon default yang rusak.
- **Invariant (Predikat AST):**
  Untuk setiap elemen deklarasi manifest web app $M \in \text{ManifestDeclarations}$:
  $$\neg \text{hasName}(M) \lor \neg \text{hasStartUrl}(M) \lor \neg \text{hasDisplay}(M) \lor \neg \text{hasIcons}(M) \implies \text{Violation (Error)}$$
  di mana $\text{ManifestDeclarations}$ mencakup elemen `<script type="application/manifest+json">` atau komponen konfigurasi manifest.
- **Mengapa Lolos Linter Standar:**
  Berkas atau blok manifest valid secara sintaksis JSON. Linter bahasa umum tidak memvalidasi skema field wajib W3C Web App Manifest.
- **Suspicious (Manifest Kehilangan Field Wajib):**
  ```tsx
  {/* Hilang: start_url, display, dan icons */}
  <script type="application/manifest+json">
    {JSON.stringify({
      name: "Desa Digital"
    })}
  </script>
  ```
- **Compliant (Seluruh Field Wajib Didefinisikan):**
  ```tsx
  {/* Seluruh field wajib terpenuhi */}
  <script type="application/manifest+json">
    {JSON.stringify({
      name: "Desa Digital",
      short_name: "Desa",
      start_url: "/",
      display: "standalone",
      icons: [
        { src: "/icon-192.png", sizes: "192x192", type: "image/png" },
        { src: "/icon-512-maskable.png", sizes: "512x512", type: "image/png", purpose: "maskable" }
      ]
    })}
  </script>
  ```
- **Engine:** L1 Syntax + L2 JSON Manifest AST (`internal/rules/pwa/manifest_required_fields_missing.go`).
- **Severity:** `error`.
- **Autofix:** No.

---

### 2.2. `pwa.icon-maskable-missing`
- **Design Rationale:** W3C Web App Manifest (Adaptive Icon Masking) & Google Android Maskable Icons Specification.
- **Konteks Realitas Mobile:**
  Mulai Android 8.0 Oreo, peluncur (*launcher*) perangkat Android memotong ikon aplikasi mengikuti bentuk sistem adaptif (lingkaran, *squircle*, atau persegi membulat). Jika manifest PWA hanya menyediakan ikon standar (`purpose: "any"` atau tanpa atribut `purpose`), Android membungkus ikon di dalam kotak putih kaku berukuran kecil (*letterboxing*) yang merusak estetika antarmuka native. Menyediakan ikon dengan `purpose: "maskable"` menjamin ikon ditampilkan penuh secara mulus tanpa batas putih.
- **Invariant (Predikat AST):**
  Untuk setiap deklarasi manifest yang memiliki koleksi `icons`:
  $$\text{hasIcons}(M) \land \neg \text{hasMaskableIcon}(M) \implies \text{Violation (Warn)}$$
  di mana $\text{hasMaskableIcon}$ memeriksa keberadaan entri ikon dengan properti `purpose` yang memuat nilai `"maskable"`.
- **Mengapa Lolos Linter Standar:**
  Ikon dengan `purpose: "any"` adalah sah secara W3C manifest schema. Schema validator umum tidak mewajibkan entri maskable, sehingga masalah baru disadari saat aplikasi diinstal di ponsel Android.
- **Suspicious (Hanya Ikon Biasa Tanpa Maskable):**
  ```tsx
  {/* Menghasilkan kotak putih jelek di launcher Android adaptif */}
  <script type="application/manifest+json">
    {JSON.stringify({
      name: "Desa Digital",
      start_url: "/",
      display: "standalone",
      icons: [
        { src: "/icon-512.png", sizes: "512x512", type: "image/png" }
      ]
    })}
  </script>
  ```
- **Compliant (Ikon Maskable Adaptif Disediakan):**
  ```tsx
  {/* Ikon maskable adaptif menyesuaikan bentuk launcher Android */}
  <script type="application/manifest+json">
    {JSON.stringify({
      name: "Desa Digital",
      start_url: "/",
      display: "standalone",
      icons: [
        { src: "/icon-512.png", sizes: "512x512", type: "image/png" },
        { src: "/icon-512-maskable.png", sizes: "512x512", type: "image/png", purpose: "maskable" }
      ]
    })}
  </script>
  ```
- **Engine:** L1 Syntax + L2 JSON Manifest AST (`internal/rules/pwa/icon_maskable_missing.go`).
- **Severity:** `warning`.
- **Autofix:** No.

---

### 2.3. `pwa.manifest-missing`
- **Design Rationale:** W3C Web App Manifest Section 4 (Linking to a Manifest) & HTML Living Standard.
- **Konteks Realitas Mobile:**
  Browser mobile hanya dapat mengenali aplikasi sebagai Progressive Web App jika dokumen HTML dasar memuat tag `<link rel="manifest" href="...">` di dalam elemen `<head>`. Tanpa tag ini, browser memperlakukan web sebagai situs desktop biasa dan tidak pernah memicu prompt penginstalan web app.
- **Invariant (Predikat AST):**
  Untuk setiap elemen `<head>` atau root layout dokumen:
  $$\neg \text{hasManifestLink}(Head) \implies \text{Violation (Warn)}$$
  di mana $\text{hasManifestLink}$ memeriksa keberadaan tag `<link rel="manifest">` dengan atribut `href` yang tidak kosong.
- **Mengapa Lolos Linter Standar:**
  Dokumen HTML tanpa link manifest adalah valid secara sintaksis. Tidak ada linter HTML default yang memeriksa ketersediaan manifest PWA.
- **Suspicious (Dokumen Root Tanpa Link Manifest):**
  ```tsx
  <head>
    <title>Layanan Surat Desa</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  ```
- **Compliant (Link Manifest Dideklarasikan di Head):**
  ```tsx
  <head>
    <title>Layanan Surat Desa</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="manifest" href="/manifest.webmanifest" />
  </head>
  ```
- **Engine:** L1 Syntax + L2 HTML Head AST (`internal/rules/pwa/manifest_missing.go`).
- **Severity:** `warning`.
- **Autofix:** No.

---

### 2.4. `pwa.start-url-inconsistency`
- **Design Rationale:** W3C Web App Manifest Section 5.2 (The start_url member) & W3C Secure Contexts.
- **Konteks Realitas Mobile:**
  Field `start_url` menentukan halaman awal yang dimuat ketika pengguna meluncurkan PWA dari homescreen ponsel. Spesifikasi W3C mewajibkan `start_url` berada dalam cakupan (*navigation scope*) aplikasi. Menetapkan `start_url` ke protokol tidak aman (`http://`), domain luar (*cross-origin*), atau jalur traversal (`../`) menyebabkan peluncuran PWA gagal dan dibatalkan oleh browser mobile.
- **Invariant (Predikat AST):**
  Untuk setiap deklarasi manifest dengan field `start_url`:
  $$\text{isInsecureProtocol}(\text{start\_url}) \lor \text{isPathTraversal}(\text{start\_url}) \lor \text{isScriptURL}(\text{start\_url}) \implies \text{Violation (Error)}$$
- **Mengapa Lolos Linter Standar:**
  Nilai `start_url` adalah string biasa di dalam JSON. Linter standar tidak memvalidasi keamanan protokol dan cakupan navigasi URL.
- **Suspicious (start_url Tidak Aman atau Melampaui Scope):**
  ```tsx
  {/* Protokol http tidak aman melanggar standar Secure Contexts PWA */}
  <script type="application/manifest+json">
    {JSON.stringify({
      name: "Desa Digital",
      start_url: "http://desa-insecure.org/app",
      display: "standalone",
      icons: [{ src: "/icon.png", sizes: "192x192", type: "image/png", purpose: "maskable" }]
    })}
  </script>
  ```
- **Compliant (start_url Valid Terhadap Root Scope):**
  ```tsx
  {/* start_url relatif terhadap root scope yang aman */}
  <script type="application/manifest+json">
    {JSON.stringify({
      name: "Desa Digital",
      start_url: "/",
      display: "standalone",
      icons: [{ src: "/icon.png", sizes: "192x192", type: "image/png", purpose: "maskable" }]
    })}
  </script>
  ```
- **Engine:** L1 Syntax + L2 JSON Manifest AST (`internal/rules/pwa/start_url_inconsistency.go`).
- **Severity:** `error`.
- **Autofix:** No.

---

## 3. Spesifikasi Detail Rule Wave 2: Apple Standalone & Security (2 Rules)

### 3.1. `pwa.apple-meta-missing`
- **Tujuan:** Safari iOS memerlukan meta tag WebKit khusus agar tampil tanpa bilah alamat dan memakai ikon kustom saat ditambahkan ke Home Screen.
- **In-Scope:** Root layout (`<head>`) yang memasang manifest tapi tidak menyertakan `<meta name="apple-mobile-web-app-capable">` dan `<link rel="apple-touch-icon">`.
- **Severity:** Warning.

### 3.2. `pwa.insecure-context-resource`
- **Tujuan:** Memastikan seluruh resource dalam PWA tidak menggunakan protokol HTTP yang tidak aman (`http://`), sesuai standar W3C Secure Contexts.
- **In-Scope:** URL literal berskema `http://` (bukan localhost) pada atribut atau pemanggilan resource.
- **Severity:** Error.

---

## 4. Spesifikasi Detail Rule Wave 3: Service Worker Lifecycle & Offline Cache (4 Rules)

### 4.1. `pwa.service-worker-no-offline-fallback`
- **Tujuan:** Memastikan Service Worker memiliki strategi offline fallback (bukan pass-through fetch kosong) agar tidak memunculkan layar dinosaurus putus koneksi.
- **Severity:** Warning.

### 4.2. `pwa.service-worker-missing`
- **Tujuan:** Memastikan proyek PWA mendaftarkan service worker untuk mengelola offline caching.
- **Severity:** Warning.

### 4.3. `pwa.service-worker-registration`
- **Tujuan:** Mendeteksi registrasi Service Worker tanpa feature detect `'serviceWorker' in navigator` atau tanpa penanganan error `.catch()`.
- **Severity:** Warning.

### 4.4. `pwa.pwa-cache-runtime-api-risk`
- **Tujuan:** Mencegah akses API DOM (`window`, `document`, `localStorage`) di dalam thread Service Worker yang memicu crash seketika.
- **Severity:** Error.

---

## 5. Ringkasan Matriks Rule `pwa.*` (10 Aturan Non-Redundan)

| Rule ID | Fokus Invarian | Wave | Severity | Engine Target |
|---|---|:---:|---|---|
| `pwa.manifest-required-fields-missing` | Field wajib manifest PWA (name, start_url, display, icons) | **W1** | error | JSON/Manifest AST |
| `pwa.icon-maskable-missing` | Ikon launcher adaptif Android (purpose: maskable) | **W1** | warning | JSON/Manifest AST |
| `pwa.manifest-missing` | Tag link manifest pada root layout (<head>) | **W1** | warning | JSX/HTML Head AST |
| `pwa.start-url-inconsistency` | Validitas URL awal aplikasi terhadap navigation scope | **W1** | error | JSON/Manifest AST |
| `pwa.apple-meta-missing` | Standalone mode Safari iOS (apple-touch-icon & meta) | **W2** | warning | JSX/HTML Head AST |
| `pwa.insecure-context-resource` | Pencegahan mixed-content HTTP di secure context | **W2** | error | String Literal AST |
| `pwa.service-worker-no-offline-fallback` | Fallback fetch handler saat koneksi terputus | **W3** | warning | Call Expression AST |
| `pwa.service-worker-missing` | Ketersediaan service worker pada proyek PWA | **W3** | warning | Project Scanner |
| `pwa.service-worker-registration` | Feature detect dan penanganan error registrasi SW | **W3** | warning | Call Expression AST |
| `pwa.pwa-cache-runtime-api-risk` | Pencegahan akses DOM/window di thread Web Worker | **W3** | error | Worker Scope AST |

---

## 6. Rule Classification & Execution Boundary

1. **Deterministic AST Rules (< 50ms pre-commit gate):**
   - `pwa.manifest-required-fields-missing`, `pwa.icon-maskable-missing`, `pwa.manifest-missing`, `pwa.start-url-inconsistency`, `pwa.apple-meta-missing`, `pwa.insecure-context-resource`, `pwa.service-worker-no-offline-fallback`, `pwa.service-worker-registration`, `pwa.pwa-cache-runtime-api-risk`.
2. **Project Scanner Layer:**
   - `pwa.service-worker-missing`.
3. **Runtime Validation Layer:**
   - Audit Lighthouse PWA dan pengujian prompt instalasi browser mobile sesungguhnya.

---

## 7. Struktur Modul Kode & Roadmap Eksekusi Wave 1

Implementasi aturan `pwa.*` ditempatkan secara modular di `internal/rules/pwa/`:

```text
internal/rules/pwa/
├── util.go                                 # Parser manifest JSON & AST helpers
├── manifest_required_fields_missing.go     # Wave 1: Required manifest fields validator
├── icon_maskable_missing.go                # Wave 1: Adaptive Android maskable icon check
├── manifest_missing.go                     # Wave 1: Head manifest link presence check
├── start_url_inconsistency.go              # Wave 1: Navigation scope & protocol safety
├── contract_test.go                        # 8-Pillars Canonical Contract Validator
└── benchmark_test.go                       # QUAL-03 Zero Allocation Benchmarks
```

Setiap rule divalidasi dengan **1-SSOT Golden Tri-Corpus** di `tests/correctness/pwa/<slug>/` yang mencakup kasus uji positif, negatif, dan adversarial untuk menjamin **nol regresi dan nol false-positive**.

