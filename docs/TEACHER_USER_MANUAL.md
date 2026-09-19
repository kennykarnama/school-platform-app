# Buku Panduan Pengguna (User Manual): Guru (Teacher)
*School Administration & Attendance Platform*

---

## Daftar Isi
1. [Pengantar](#1-pengantar)
2. [Login, Profil, dan Akses](#2-login-profil-dan-akses)
3. [Manajemen Siswa (`#students`)](#3-manajemen-siswa-students)
4. [Pencatatan Absensi Harian (`#absensi/input`)](#4-pencatatan-absensi-harian-absensiinput)
5. [Rekapitulasi Absensi Individu (`#absensi/stats`)](#5-rekapitulasi-absensi-individu-absensistats)
6. [Rekapitulasi Klasikal (`#absensi/stats-classical`)](#6-rekapitulasi-klasikal-absensistats-classical)
7. [Tips & Pemecahan Masalah Singkat](#7-tips--pemecahan-masalah-singkat)

---

## 1. Pengantar

Panduan ini ditujukan bagi **Guru** untuk mengelola data kehadiran siswa, pemindahan kelas siswa, serta memantau dan mengunduh laporan rekapitulasi absensi pada aplikasi Administrasi Sekolah.

Sebagai pengguna dengan hak akses **Guru**, menu yang tersedia untuk Anda meliputi:
- **Home**: Halaman beranda sambutan.
- **Siswa**: Direktori pencarian siswa dan pemindahan rombongan belajar.
- **Absensi**:
  - **Input**: Pengisian kehadiran harian kelas.
  - **Rekap**: Rekapitulasi presensi individu dan tren statistik.
  - **Rekap Klasikal**: Rekapitulasi agregat kehadiran harian kelas.
- **About**: Informasi sistem aplikasi.

---

## 2. Login, Profil, dan Akses

### 2.1 Masuk ke Sistem (Login)
1. Buka tautan portal aplikasi di web browser.
2. Masukkan **Username / Email** dan **Password** Anda.
3. Klik tombol **Login**.

### 2.2 Memeriksa Informasi Akun
Pada panel navigasi samping (*sidebar*):
- Avatar inisial nama Anda.
- **Nama Lengkap Guru**.
- **NIP / Nomor Identitas Pegawai**.
- **Nama Sekolah**.
- Label peran: `Guru`.

### 2.3 Mengubah Password
1. Buka tautan menu ganti password (`/#password/change`).
2. Masukkan password lama Anda.
3. Masukkan dan konfirmasi password baru.
4. Klik **Simpan**.

### 2.4 Keluar dari Aplikasi (Logout)
- Klik tombol **Keluar** di bagian bawah menu navigasi setelah selesai menggunakan aplikasi.

---

## 3. Manajemen Siswa (`#students`)

Halaman ini digunakan untuk melihat daftar siswa, status aktif, penempatan kelas, serta melakukan pemindahan kelas secara massal (*bulk transfer*).

### 3.1 Mencari dan Memfilter Siswa
1. Masuk ke menu **Siswa** di menu samping.
2. Gunakan filter pencarian:
   - **Cari nama atau ID**: Masukkan kata kunci pencarian, lalu tekan `Enter` atau klik ikon pencarian 🔍.
   - **Dropdown Kelas**: Memfilter daftar siswa pada rombongan belajar tertentu.
   - **Dropdown Status**: Memfilter siswa berdasarkan status `Aktif` atau `Tidak Aktif`.
3. Klik tombol **Reset** untuk menghapus filter dan memuat kembali seluruh siswa.

### 3.2 Memindahkan Siswa Antar Kelas (*Transfer*)
1. Beri tanda centang pada kotak di samping nama siswa yang ingin dipindahkan (bisa memilih satu atau lebih siswa).
2. Klik tombol **Pindahkan Siswa Terpilih**.
3. Pada kotak dialog:
   - Pilih **Tahun Akademik Tujuan**.
   - Pilih **Kelas Tujuan**.
4. Klik tombol **Pindahkan**.

> **Catatan:** Fitur penambahan siswa baru dan penonaktifan siswa secara permanen dikelola oleh Administrator Sekolah. Guru memiliki akses untuk melihat data dan melakukan pemindahan kelas.

---

## 4. Pencatatan Absensi Harian (`#absensi/input`)

Fitur utama guru untuk mencatat presensi kehadiran siswa di kelas yang diajar.

### 4.1 Langkah Pengisian Absensi Harian
1. Buka menu **Absensi** > **Input**.
2. Pilih kriteria pada bilah filter:
   - **Tahun Akademik - Semester** (Contoh: *2023/2024 - Ganjil*).
   - **Pilih Kelas** (Contoh: *7A*).
   - **Tanggal Kehadiran**: Klik pada bidang tanggal untuk memilih tanggal dari kalender pop-up (*datepicker*).
3. Klik tombol **Cari** (ikon kaca pembesar 🔍).
4. Daftar siswa kelas tersebut akan ditampilkan pada tabel.
5. Pada kolom **Kehadiran**, tentukan status kehadiran tiap siswa:
   - `Hadir`
   - `Sakit`
   - `Izin`
   - `Alpha` (Tanpa keterangan)
6. Setelah seluruh data siswa selesai diisi, gulir ke bagian bawah tabel dan klik tombol **Submit Absensi**.
7. Tunggu hingga muncul pesan notifikasi berhasil di pojok layar.

### 4.2 Menambahkan Siswa Langsung ke Kelas
1. Jika ada siswa yang belum masuk ke dalam daftar kelas, klik tombol **Tambah Siswa** di bagian bawah tabel absensi.
2. Isi nama siswa, pilih tahun akademik serta kelas tujuan, lalu simpan.

### 4.3 Menghapus / Menonaktifkan Siswa dari Absensi Kelas
1. Klik tombol merah **Hapus** pada baris siswa yang bersangkutan.
2. Pada kotak dialog konfirmasi, masukkan alasan (contoh: *Pindah sekolah* atau *Salah kelas*).
3. Klik tombol **Simpan**.

### 4.4 Memindahkan Siswa dari Tabel Absensi
1. Centang siswa yang bersangkutan pada tabel absensi.
2. Klik tombol **Pindahkan (X)** di bawah tabel.
3. Tentukan tahun ajaran dan kelas tujuan, lalu klik **Pindahkan**.

---

## 5. Rekapitulasi Absensi Individu (`#absensi/stats`)

Halaman ini digunakan untuk melihat total dan persentase kehadiran setiap siswa secara detail dalam periode tertentu (misal: bulanan atau semesteran).

### 5.1 Menampilkan Data Rekap
1. Buka menu **Absensi** > **Rekap**.
2. Tentukan filter:
   - **Tahun Akademik - Semester**
   - **Pilih Kelas**
   - **Dari tanggal**: Tanggal awal periode.
   - **Sampai tanggal**: Tanggal akhir periode.
3. Klik tombol **Cari** (🔍).

### 5.2 Membaca Tabel dan Navigasi
- **Kolom Tabel**: Menampilkan jumlah kehadiran per jenis (`Hadir`, `Sakit`, `Izin`, `Alpha`), total hari pertemuan, serta persentase kehadiran masing-masing kategori.
- **Navigasi Layar Lebar**: Gunakan tombol panah kiri (`←`) dan panah kanan (`→`) di atas tabel untuk menggeser tampilan tabel jika kolom terpotong.

### 5.3 Grafik Tren Kehadiran
- Di bawah tabel rekap, grafik visual interaktif akan otomatis terbentuk untuk melihat tren kehadiran siswa pada periode yang dipilih.

### 5.4 Mengunduh File Excel (`.xls`)
- Klik tombol **Download** (ikon unduh 📥) berwarna biru di atas tabel.
- File spreadsheet `rekap.xls` akan terunduh secara otomatis.

---

## 6. Rekapitulasi Klasikal (`#absensi/stats-classical`)

Halaman ini menyajikan rekap agregat kelas secara keseluruhan per tanggal untuk melihat tingkat partisipasi kelas secara makro.

### 6.1 Menampilkan Rekap Klasikal
1. Buka menu **Absensi** > **Rekap Klasikal**.
2. Tentukan parameter:
   - **Tahun Akademik - Semester**
   - **Pilih Kelas**
   - **Dari tanggal** dan **Sampai tanggal**
3. Klik tombol **Cari** (🔍).

### 6.2 Membaca Laporan Klasikal
- Data ditampilkan per tanggal pelaksanaan pembelajaran.
- Setiap tanggal memuat:
  - Jumlah total siswa kelas.
  - Jumlah siswa per status (`Hadir`, `Sakit`, `Izin`, `Alpha`).
  - Persentase klasikal kelas pada hari tersebut.

### 6.3 Mengunduh Laporan Klasikal
- Klik tombol download Excel yang tersedia untuk mengunduh rekapitulasi klasikal.

---

## 7. Tips & Pemecahan Masalah Singkat

| Situasi | Solusi |
| :--- | :--- |
| **Data absensi tidak tersimpan** | Pastikan Anda menekan tombol **Submit Absensi** di bawah tabel sebelum berpindah halaman atau menutup browser. |
| **Siswa tidak muncul di daftar kelas** | Pastikan filter **Tahun Akademik** dan **Kelas** yang dipilih sudah tepat. Jika siswa baru bergabung, gunakan tombol **Tambah Siswa** atau koordinasikan dengan Admin Sekolah. |
| **Tabel terpotong di layar laptop/tablet** | Gunakan tombol navigasi geser panah kiri/kanan di atas tabel, atau geser tabel ke samping (*horizontal scroll*). |
| **Lupa password** | Hubungi Administrator Sekolah untuk melakukan reset password akun Anda. |
