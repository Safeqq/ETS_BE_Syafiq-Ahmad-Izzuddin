# PROSES

## 1. Persona dan masalahnya
Pengguna persona saya adalah admin himpunan. Untuk permasalahannya ruang himpunan tidak boleh langsung terpakai. Setiap permintaan harus disetujui admin dulu, dan permintaan yang ditolak harus disertai alasan. Lalu dengan fitur yang akan saya buat ini diharapkan seseorang yang ingin menggunakan ruang himpunan dapat melalui alur peminjaman yang jelas (diajukan, disetujui/ditolak (dengan alasan), selesai).

(2-3 kalimat dengan kata-katamu sendiri: siapa pengguna persona kamu, masalah apa yang diselesaikan fitur khasnya.)

## 2. Rancangan struktur data
Rancangan Awal:
Untuk saat ini rancangan akan terdapat struct
- `peminjaman` yang berisi {id, tanggal, slot_mulai, slot_selesai, jumlah, status, alasan} 
- `peminjam` {sumber_saya_id, nama}. 
Slice yang akan disimpan adalah sumber_daya_id, nama, status, dan alasan karena itu merupakan data terpisah yang akan digabungkan (sumber_daya_id dan nama) dan data yang akan berubah seiring waktu (status dan alasan).

Perubahan Rancangan:
- Struct peminjam tidak jadi dibuat terpisah. Di kontrak API, peminjam hanya berupa teks (nama), jadi cukup menjadi field `Peminjam` di struct `Peminjaman`. Struct `Peminjaman` sekarang berisi {id, sumber_daya_id, peminjam, tanggal, slot_mulai, slot_selesai, jumlah, status, alasan}.
- Ditambah struct `SumberDaya` {id, nama, jenis, stok} karena B1 meminta data sumber daya disimpan sendiri, dan `sumber_daya_id` di peminjaman menunjuk ke data ini.
- Yang disimpan di slice bukan per field (sumber_daya_id, nama, status, alasan), tetapi satu objek utuh: `[]SumberDaya` dan `[]Peminjaman`. Kalau per field, data satu peminjaman tersebar di banyak slice dan mudah tidak sinkron. Status dan alasan cukup diubah langsung di objek peminjamannya.
- Kedua slice disimpan di dalam struct `Server` bersama `sync.Mutex` (supaya aman untuk request bersamaan) dan penghitung id untuk masing-masing jenis data (supaya id tidak terpakai ulang setelah data dihapus).
- Ditambah struct input (`SumberDayaInput`, `PeminjamanInput`, `TolakInput`) yang field-nya pointer, untuk membaca body request dan mengecek apakah field wajib benar-benar dikirim.

(Tulis struct dan slice/map apa yang kamu simpan, dan kenapa dipilih begitu. Rancang ini sebelum menulis kode, lalu catat kalau berubah di tengah jalan.)

## 3. Satu keputusan desain
Keputusan: validasi peminjaman saya bagi menjadi dua tahap. Tahap pertama (`validasiPeminjaman`) memeriksa isi body saja, seperti field wajib, peminjam, tanggal, slot, dan jumlah minimal, dan dijalankan sebelum mutex dikunci. Tahap kedua (`cekSumberDayaDanJumlah`) memeriksa hal yang butuh data tersimpan, yaitu sumber daya ada dan jumlah tidak melebihi stok, dan dijalankan setelah mutex dikunci.

Pilihan lain yang saya pertimbangkan:
- Semua validasi di luar lock. Ini berbahaya karena antara "cek" dan "simpan" ada jeda, dan request lain bisa mengubah data di jeda itu.
- Semua validasi di dalam lock. Ini aman, tapi request lain harus ikut menunggu walaupun yang diperiksa hanya isi body yang tidak menyentuh data bersama.

Alasan memilih dua tahap: pemeriksaan yang tidak menyentuh data bersama tidak perlu membuat request lain menunggu, sedangkan pemeriksaan yang membaca data bersama harus berada dalam kunci yang sama dengan penyimpanan supaya tidak disela request lain. Pola ini juga akan dipakai untuk cek bentrok di B5 dan keamanan request bersamaan di B7.

(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
