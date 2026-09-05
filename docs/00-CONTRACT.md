# 00-CONTRACT: Charites Architectural Governance & Cross-Pillar Authority Contract

> **Kode Dokumen:** `GOV-00-CONTRACT`
> **Peran:** Single Source of Truth (SSOT) Tata Kelola & Hirarki Otoritas Dokumentasi
> **Status:** Active & Normative
> **Standar Rujukan:** ISO/IEC/IEEE 42010 Architecture Description & RFC 2119 Key Words

Dokumen ini adalah **kontrak tata kelola utama** yang mengikat seluruh tahapan (Fase 0 hingga 8) dalam repositori Charites (`github.com/will2469/charites`). Kontrak ini dibuat untuk mencegah pergeseran semantik (*semantic drift*), kontradiksi lintas pilar, dan kebocoran perilaku antar dokumen.

---

## 1. Hirarki Otoritas 6 Pilar (The 6-Pillar Authority Model)

Setiap tahapan pengembangan Charites didokumentasikan dalam 6 folder pilar dengan batas otoritas tegas:

```
01-SPEC (WHAT)
   │
   ▼
02-ARCH (HOW)
   │
   ▼
03-TEST (PROOF)
   │
   ▼
04-QUALITY (QUALITY THRESHOLD)
   │
   ▼
05-RELEASE (SEMVER & RELEASES)
   │
   ▼
06-ROADMAP (PHASE GATE ONLY)
```

### 1.1. Aturan Otoritas Normatif Lintas Pilar

1. **`01-SPEC` = WHAT (Otoritas Tertinggi Perilaku Produk)**
   - Mendefinisikan kebutuhan fungsional, antarmuka CLI, tata bahasa masukan, skema keluaran, kode keluar, dan batas semantik.
   - Seluruh batasan operasional dan perilaku yang tampak oleh pengguna (*user-visible behavior*) **WAJIB** bersumber dari SPEC.

2. **`02-ARCH` = HOW (Otoritas Realisasi & Struktur Internal)**
   - Mendefinisikan topologi komponen, struktur data Go, pola konkurensi, dan enkapsulasi algoritma.
   - **Aturan Tegas:** ARCH **DILARANG** memperkenalkan perilaku baru, fitur tersembunyi, atau fallback yang tidak terdefinisi di SPEC (*ARCH MUST NOT introduce behavior not defined by SPEC*).

3. **`03-TEST` = PROOF (Otoritas Pembuktian & Verifikasi)**
   - Mendefinisikan skenario pengujian, harness pembuktian batas (*boundary matrix*), korpus regresi, dan pengujian subprocess E2E.
   - **Aturan Tegas:** TEST **DILARANG** menciptakan persyaratan fungsional baru yang tidak ada di SPEC (*TEST MUST NOT invent requirements*).

4. **`04-QUALITY` = QUALITY THRESHOLD (Otoritas Ambang Batas Non-Fungsional)**
   - Mendefinisikan batas alokasi memori, anggaran performa (*Performance Budgets*), metodologi benchmark terisolasi, keamanan memori/thread, dan ambang batas cakupan pengujian (*coverage*).
   - **Aturan Tegas:** QUALITY **DILARANG** menyisipkan batasan fungsional (seperti batas ukuran berkas atau semantik error) kecuali batasan tersebut telah diangkat dan dipromosikan secara eksplisit ke SPEC (*QUALITY MUST NOT introduce functional behavior unless explicitly promoted to SPEC*).

5. **`05-RELEASE` = SEMVER & RELEASES (Otoritas Manajemen Rilis & Perubahan)**
   - Mendefinisikan penomoran SemVer 2.0.0, changelog, audit backward-compatibility, dan kebijakan breaking changes.

6. **`06-ROADMAP` = PHASE GATE ONLY (Otoritas Evaluasi Gerbang Transisi)**
   - Menetapkan daftar *deliverables* berkas dan evaluasi kelulusan gerbang (`ROAD-XX-GATE-001` s/d `GATE-004`).
   - **Aturan Tegas:** ROADMAP **DILARANG** menjelaskan ulang deskripsi fungsional atau menciptakan persentase metrik tandingan yang bertentangan dengan QUALITY (*ROADMAP MUST NOT redefine requirements or duplicate behaviors*).

---

## 2. Matriks Resolusi Konflik (Conflict Resolution Matrix)

Jika ditemukan inkonsistensi antar-dokumen, aturan resolusi berikut berlaku secara otomatis tanpa kompromi:

| Sifat Kontradiksi | Dokumen Menang | Dokumen Kalah | Tindakan Wajib |
| :--- | :--- | :--- | :--- |
| **Perilaku Fungsional / CLI / Skema** | `01-SPEC` | `02-ARCH` / `04-QUALITY` | ARCH & QUALITY wajib diselaraskan dengan SPEC. |
| **Struktur Internal / Algoritma** | `02-ARCH` | `03-TEST` | Harness TEST disesuaikan dengan arsitektur internal ARCH. |
| **Ambang Batas Kualitas / Coverage** | `04-QUALITY` | `06-ROADMAP` | ROADMAP wajib merujuk `QUAL-XX Compliance = PASS`. |
| **Nomenklatur Rule ID** | Charites Rule ID (`<category>.<slug>`) | Semgrep ID | Dilarang menyebut Semgrep sebagai otoritas Rule ID. |
| **Performa & Nanodetik** | Target Desain / Performance Budget | Hard Gate Mutlak | Ubah angka mutlak menjadi metodologi benchmark terukur. |

---

## 3. Presedensi Kebijakan Eksekusi (Unified Execution Precedence)

Untuk seluruh subsistem Charites (Rule, Config, Scanner, CLI), hierarki presedensi tunggal ditetapkan sebagai berikut:

$$\text{Registry (Base Candidates)} \longrightarrow \text{CLI Scope (--category, --rule)} \longrightarrow \text{Config Policy (charites.yaml)} \longrightarrow \text{Filesystem Ignore (.charitesignore)} \longrightarrow \text{Built-in Invariant (Hard Exclusion)}$$

1. **CLI Scope:** Menyaring kandidat rule atau path target.
2. **Config Policy:** Menetapkan status aktif/nonaktif dan severity override. Kebijakan `off` pada config mengalahkan flag `--rule` pada CLI.
3. **Filesystem Ignore:** Membatasi penelusuran berkas berdasarkan aturan pengguna.
4. **Built-in Invariant:** Pengecualian mutlak bawaan (`node_modules/`, `.git/`, direktori symlink) tidak dapat dinegasi oleh aturan pengguna atau target berkas langsung.

---

## 4. Invarian Bebas Residu & Kebersihan Siklus Hidup (Zero Residual Footprint Invariant)

Untuk melindungi integritas sistem operasi pengguna dan memastikan portabilitas murni, Charites memberlakukan **Invarian Bebas Residu & Kebersihan Siklus Hidup** (*Zero Residual Footprint Invariant*) yang mengikat seluruh tahapan rilis dan biner eksekutabel:

### 4.1. Filosofi Biner Tunggal Mandiri (Pure Standalone Executable)
1. **Zero External Runtime Dependencies:** Charites didistribusikan sebagai artefak biner tunggal mandiri (`CGO_ENABLED=0`) tanpa membutuhkan pustaka sistem dinamis (glibc/musl) khusus, driver, interpreter eksternal (Node.js, Python), atau daemon pendukung.
2. **Atomic Placement:** Instalasi Charites semata-mata adalah penempatan berkas biner tunggal ke dalam direktori eksekutabel pilihan pengguna (misalnya `/usr/local/bin/charites` atau `$HOME/.local/bin/charites`).

### 4.2. Invarian Bebas Polusi Host (Zero Host Pollution Invariant)
1. **Dilarang Menulis di Luar Workspace Target:** Selama proses eksekusi, pemindaian, analisis, pelaporan, atau parsing, Charites **DILARANG KERAS** membuat, memodifikasi, atau meninggalkan berkas/direktori di lingkungan sistem pengguna di luar target workspace yang ditentukan secara eksplisit.
2. **Dilarang Membuat Direktori Global Tersembunyi:** Charites tidak pernah membuat atau mengasumsikan keberadaan direktori konfigurasi atau cache global pengguna, termasuk namun tidak terbatas pada:
   - Linux/Unix: `~/.config/charites`, `~/.cache/charites`, `~/.charites`, `/var/cache/charites`, `/tmp/charites*`
   - macOS: `~/Library/Application Support/charites`, `~/Library/Caches/charites`
   - Windows: `%APPDATA%\charites`, `%LOCALAPPDATA%\charites`, `%TEMP%\charites*`
3. **Dilarang Memodifikasi Konfigurasi Shell / Sistem Host:** Charites tidak pernah menginjeksikan script atau alias ke `.bashrc`, `.zshrc`, `.profile`, Windows Registry, atau environment variable sistem secara tersembunyi.
4. **Dilarang Menjalankan Background Daemon / Sockets / PID Files:** Charites beroperasi sebagai CLI synchronous ephemeral run-to-completion. Charites **DILARANG** mendaftarkan background daemon, systemd unit, launchd plist, cron job, Windows service, telemetry emitter di background, Unix domain socket, atau meninggalkan PID file.

### 4.3. Jaminan Bersih Total Siklus Hidup: Update & Uninstall (Zero Leftover Guarantee)
1. **Invarian Pembaruan (Clean Update Invariant):**
   - Pembaruan versi Charites dilakukan secara atomik murni dengan mengganti biner lama dengan biner baru (`cp new_charites $(which charites)` atau `curl ... -o $(which charites)`).
   - Tidak ada skema database tersembunyi, state migrasi, atau cache versi lama yang tertinggal di mesin pengguna yang dapat memicu konflik antar-versi.
2. **Invarian Pencopotan (Clean Uninstall Invariant):**
   - Pencopotan Charites dilakukan secara tuntas dan bersih 100% hanya dengan **menghapus satu-satunya berkas biner `charites`** (`rm $(which charites)`).
   - Menghapus biner Charites secara deterministik menjamin tidak ada satupun artefak sisa (*zero residual leftover*), berkas konfigurasi yatim piatu, cache tersembunyi, entri registri, atau proses zombie di sistem operasi pengguna.
3. **In-Memory Ephemeral Execution:**
   - Seluruh pemrosesan memori (AST buffer, cache token Tailwind, dictionary IR, context) dialokasikan dalam memori proses (*heap*) dan dibebaskan seketika saat proses CLI berakhir (`os.Exit`). Tidak ada berkas sementara (*temporary swap files*) yang ditulis ke disk.

