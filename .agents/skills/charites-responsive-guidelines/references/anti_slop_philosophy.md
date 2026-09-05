# Anti-AI Slop Philosophy & Industrial Functionalism

> Landasan psikologi kognitif, persepsi visual Gestalt, dan prinsip fungsionalisme industri untuk mengeliminasi degradasi antarmuka akibat _generative AI slop_.

---

## 1. Anatomi & Taksonomi "AI UI Slop"

Fenomena _AI Slop_ pada antarmuka frontend lahir karena model bahasa besar (LLM) dilatih atas distribusi kode publik yang didominasi oleh _marketing landing pages_, portofolio instan, dan tren visual dangkal. Model memperlakukan UI sebagai artefak grafis statis, bukan mesin status reaktif yang berpusat pada manusia.

Empat dimensi kegagalan utama AI Slop:

| Dimensi Kerusakan                      | Manifestasi Visual & Kode                                                                                                                             | Dampak pada Pengguna                                                                                                                       |
| :------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------- |
| **Visual Clutter (Kebisingan Grafis)** | Gradasi ungu/indigo (`from-purple-500 to-indigo-600`), efek _glow_ neon gelap, `backdrop-filter: blur()`, dan `rounded-3xl` seragam pada semua kartu. | Mengaburkan batas struktural, merusak rasio kontras, dan membebani GPU peramban ponsel murah.                                              |
| **Token & Hierarchy Drift**            | Penggunaan nilai arbitrer sewenang-wenang (`p-[19px]`, `m-[7px]`, `w-[342px]`, `#6366F1`) yang melanggar skala modular.                               | Menghancurkan ritme vertikal visual, menghasilkan kepadatan konten yang timpang, dan merusak keselarasan tema.                             |
| **Affordance & State Erasure**         | Tombol dan input form dibuat tanpa `:hover`, `:active`, `:focus-visible`, atau `:disabled`. Pembersihan outline via `outline: none`.                  | Membutakan pengguna navigasi keyboard, membingungkan pembaca layar (_screen reader_), dan menghilangkan konfirmasi taktil sentuhan.        |
| **Template Amnesia**                   | Setiap halaman baru terasa seperti template yang berbeda dari situs lain. Inkonsistensi tipografi, modal, dan deklarasi CSS antar-halaman.            | Menghancurkan model mental pengguna (_cognitive mental model_), membuat sistem terasa seperti aplikasi tambal sulam yang tidak terpercaya. |

---

## 2. Teori Beban Kognitif (John Sweller) & Psikologi Gestalt

### A. Cognitive Load Theory (Sweller, 1988)

Arsitektur kognitif manusia dibatasi oleh memori kerja (_working memory_) yang sangat terbatas. Total beban kognitif terbagi menjadi:
$$L_{\text{total}} = L_{\text{intrinsic}} + L_{\text{germane}} + L_{\text{extraneous}}$$

1. **Intrinsic Load ($L_{\text{intrinsic}}$):** Beban esensial untuk memproses tugas operasional nyata (misal: memverifikasi NIK warga desa, mengajukan surat bansos, atau membaca rekonsiliasi kas).
2. **Germane Load ($L_{\text{germane}}$):** Beban mental produktif yang dipakai untuk memahami integrasi data dan membangun pemahaman alur.
3. **Extraneous Load ($L_{\text{extraneous}}$):** Beban mental sia-sia akibat antarmuka yang buruk, gradasi mencolok, batas kartu yang buram, dan hierarki visual yang membingungkan.

> **Hukum Rekayasa:** _AI Slop_ melipatgandakan $L_{\text{extraneous}}$. Korteks visual pengguna dipaksa terus-menerus memfilter kebisingan visual yang tidak fungsional, memicu **Aesthetic Fatigue** (kelelahan visual) dan **Banner Blindness** (pengguna mengabaikan tombol utama karena mengiranya sebagai iklan/interupsi visual).

### B. Prinsip Gestalt dalam Antarmuka Tenang

- **Principle of Common Region & Proximity:** Elemen-elemen yang berada dalam batas ruang yang jelas akan dikelompokkan secara otomatis oleh otak. Kartu dengan border tegas (`1px solid var(--border-subtle)`) jauh lebih mudah dipindai daripada kartu melayang tanpa batas tepi.
- **Figure-Ground Separation:** Penggunaan _glassmorphism_ berlebihan merusak deteksi tepi visual. Konten wajib memiliki kontras latar belakang yang solid dan terdefinisi agar mata tidak lelah membedakan teks dari tekstur latar belakang.

---

## 3. Fungsionalisme Industri & Calm Technology

### A. Prinsip Dieter Rams (1976) & Saint-Exupéry

- **Prinsip ke-10 Dieter Rams:** _"Good design is as little design as possible"_ (Desain yang baik adalah desain yang sesedikit mungkin). Desain tidak boleh bersaing mencari perhatian; ia harus mundur agar fungsi utilitas muncul tanpa halangan.
- **Aforisme Antoine de Saint-Exupéry:** _"Kesempurnaan tercapai bukan saat tidak ada lagi yang bisa ditambahkan, melainkan saat tidak ada lagi yang bisa dikurangkan."_

### B. Calm Technology (Mark Weiser & John Seely Brown, 1995)

Sistem yang baik beroperasi dengan tenang di tepi kesadaran pengguna (_periphery_), bergerak ke pusat perhatian hanya saat diperlukan, dan segera mundur kembali setelah tugas selesai. Antarmuka layanan publik desa tidak boleh berteriak; ia harus melayani dengan tenang, presisi, dan dapat diandalkan.

---

## 4. Matriks Komparasi: AI Slop vs Grounded Industrial Design

| Dimensi Desain        | AI-Generated Slop (Pola Kegagalan)                                                                     | Grounded Industrial Design (Standar Charites)                                                                          |
| :-------------------- | :----------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------- |
| **Tipografi**         | Font generik tanpa kalibrasi, hierarki datar, _line-height_ tidak sinkron dengan ukuran teks.          | Skala tipografi matematis ketat (_Major Third 1.25_) terikat pada baseline vertikal modular (Inter & JetBrains Mono).  |
| **Warna & Kontras**   | Gradasi ungu, violet, dan cyan berlebihan. Teks abu-abu terang di atas putih (gagal kontras).          | Palet tenang: kanvas netral (`--surface-canvas`), kartu solid (`--surface-card`), rasio kontras $\ge 4.5:1$ (WCAG AA). |
| **Ritme Spasi**       | Nilai arbitrer kurung siku (`p-[19px]`, `m-[7px]`). Elemen tersebar tanpa gravitasi visual.            | Skala modular kelipatan 4px/8px (`var(--primitive-space-*)`). Pengelompokan Gestalt ketat.                             |
| **Status Interaktif** | Tidak ada `:hover`, `:active`, `:focus-visible`. Menghapus ring fokus via `outline: none`.             | Kontrak status lengkap: transisi hover halus, tactile active depth, dan ring fokus `:focus-visible` kontras tinggi.    |
| **Adaptasi Mobile**   | Tampilan desktop hanya dikecilkan. Tabel di-overflow horizontal; modal desktop di tengah layar ponsel. | Mobile-first sejati: tabel diubah jadi _progressive-disclosure cards_, modal diubah jadi _bottom sheet_ di jempol.     |
| **Responsivitas**     | `@media` query global kaku yang pecah saat komponen ditaruh di dalam sidebar atau grid.                | Container Queries (`@container`) sehingga komponen adaptif secara intrinsik terhadap ruang lokalnya.                   |
| **Runtime JS**        | SPA monolitik yang mengirim megabyte JS, mengunci CPU ponsel murah, dan memicu TBT tinggi.             | Islands Architecture (Astro SSR): 0 KB client JS default, hidrasi selektif hanya saat diperlukan.                      |

---

## 5. Sitasi & Bibliografi Akademis

- Sweller, J. (1988). _Cognitive load during problem solving: Effects on learning._ Cognitive Science, 12(2), 257-285.
- Rams, D. (1976). _Ten Principles for Good Design._ Vitsoe Archive.
- Weiser, M., & Brown, J. S. (1995). _Designing Calm Technology._ Xerox PARC.
- Nielsen Norman Group (2024-2026). _Aesthetic Fatigue and the Homogenization of Generative Interfaces._
- W3C Web Accessibility Initiative (2023-2026). _WCAG 2.2 Understanding SC 2.4.13 Focus Appearance & SC 2.5.8 Target Size._
