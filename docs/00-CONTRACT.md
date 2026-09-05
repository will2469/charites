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
