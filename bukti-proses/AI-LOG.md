# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Memahami requirement, Setup
- Prompt: "Pahami requirement B1-B7 ini, kira-kira implementasi apa yang saya butuhkan dulu sebelum mulai bikin endpoint?"
- Dipakai: rancangan dua entitas (`SumberDaya` dan `Peminjaman`) yang disimpan di memori dalam slice, plus struct `Server` yang berisi `sync.Mutex`, `daftarSumber`, `daftarPinjam`, `idSumberBerikut`, dan `idPinjamBerikut` (`main.go`).
- Penjelasan: belum ada database, jadi data disimpan di slice. Mutex disiapkan dari awal karena B7 meminta aman untuk permintaan bersamaan.

## Penghitung ID, Setup
- Prompt: "Kenapa ID-nya pakai penghitung sendiri, bukan `len(slice) + 1`?"
- Dipakai: tidak ada kode baru (penjelasan saja).
- Penjelasan: nanti di B6 ada DELETE. Setelah data dihapus `len` mengecil, sehingga ID baru bisa sama dengan ID lama. Penghitung yang selalu naik menjamin ID unik.

## Struct input dengan pointer, Setup
- Prompt: "Kenapa ada `SumberDayaInput` yang field-nya pointer semua? Bukannya sudah ada `SumberDaya`?"
- Dipakai: `SumberDayaInput`, `PeminjamanInput`, `TolakInput` di `model.go`.
- Penjelasan: pointer membedakan "field tidak dikirim" (`nil`) dengan "field dikirim bernilai nol" (misalnya `"stok": 0`). Dengan begitu kita bisa memberi pesan "wajib diisi".

## Helper response JSON, Setup
- Prompt: "Buatkan helper untuk kirim response JSON, karena B4 bilang error harus selalu JSON."
- Dipakai: `kirimJSON` dan `kirimError` di `helper.go`.
- Penjelasan: `Header().Set` harus dipanggil sebelum `WriteHeader`, karena setelah status ditulis header sudah terkirim. `kirimError` selalu membungkus pesan menjadi `{"error": "..."}`.

## Handler sumber daya, B1
- Prompt: "Coba generate-kan handler untuk B1: membuat sumber daya, melihat daftar, dan melihat satu sumber daya."
- Dipakai: `buatSumberDaya`, `lihatSemuaSumberDaya`, `lihatSatuSumberDaya` di `sumber_daya.go`; route `POST /sumber-daya`, `GET /sumber-daya`, `GET /sumber-daya/{id}` di `NewHandler` (`main.go`).
- Penjelasan: alur POST adalah decode JSON, validasi, kunci mutex, beri ID, simpan, lalu balas 201. Go 1.22 sudah mendukung method dan `{id}` langsung di `ServeMux`.

## Validasi di luar lock, B1
- Prompt: "Kenapa validasinya dilakukan sebelum `s.mu.Lock()`?"
- Dipakai: tidak ada kode baru (penjelasan saja).
- Penjelasan: validasi hanya memeriksa body dan tidak menyentuh data bersama. Kalau dikerjakan di luar lock, request lain tidak perlu menunggu. Lock cukup dipakai saat membaca atau mengubah slice dan penghitung ID.

## Validasi sumber daya, B1
- Prompt: "Jelaskan `validasiSumberDaya`, kenapa return-nya `(SumberDaya, string)`?"
- Dipakai: `validasiSumberDaya` di `validasi.go`.
- Penjelasan: string kedua adalah pesan error; kalau `""` berarti valid dan struct pertama sudah berisi data bersih (sudah di-`TrimSpace`). Aturannya: ketiga field wajib ada, nama tidak kosong dan maksimal 50 karakter (`utf8.RuneCountInString` agar huruf non-ASCII dihitung satu), jenis hanya `"ruang"` atau `"alat"`, stok minimal 1.

## Ambil ID dari URL, B1
- Prompt: "Untuk GET `/sumber-daya/{id}`, gimana kalau id-nya bukan angka?"
- Dipakai: `ambilID` di `helper.go`.
- Penjelasan: `ambilID` mengubah `r.PathValue("id")` ke `int` dan mengembalikan `ok`. Kalau bukan angka hasilnya 400, kalau angka tapi tidak ditemukan hasilnya 404.

## Cari berdasarkan index, B1
- Prompt: "Kenapa `cariSumberDaya` return index, bukan struct-nya langsung?"
- Dipakai: `cariSumberDaya` di `sumber_daya.go`.
- Penjelasan: B6 (PUT/DELETE) butuh posisi elemen di slice untuk mengganti atau menghapusnya; mengubah salinan struct tidak mengubah slice. Fungsi ini tidak mengunci sendiri, jadi harus dipanggil saat mutex sedang dikunci.

## Buat peminjaman, B2
- Prompt: "Sekarang B2. Coba generate-kan fungsi untuk membuat peminjaman: satu sumber daya, satu tanggal, rentang slot, dan jumlah unit."
- Dipakai: `buatPeminjaman` di `peminjaman.go`; `validasiPeminjaman` dan `cekSumberDayaDanJumlah` di `validasi.go`; route `POST /peminjaman` di `NewHandler` (`main.go`).
- Penjelasan: validasi dibagi dua. `validasiPeminjaman` hanya memeriksa body (di luar lock), sedangkan `cekSumberDayaDanJumlah` butuh `daftarSumber` (di dalam lock).

## Cek sumber daya di dalam lock, B2
- Prompt: "Kenapa pengecekan sumber daya dipisah dan harus di dalam lock? Kan bisa dicek di awal saja."
- Dipakai: tidak ada kode baru (penjelasan saja).
- Penjelasan: kalau dicek di luar lock, ada jeda antara "cek" dan "simpan", dan request lain bisa mengubah data di jeda itu. Mengecek dan menyimpan dalam satu lock membuat keduanya tidak bisa disela. Prinsip yang sama nanti dipakai di B5 dan B7.

## Validasi tanggal dan slot, B2
- Prompt: "Jelaskan validasi tanggal dan slot-nya."
- Dipakai: bagian tanggal, slot, dan jumlah di `validasiPeminjaman` (`validasi.go`).
- Penjelasan: `time.Parse("2006-01-02", ...)` memakai layout referensi Go untuk `YYYY-MM-DD` dan juga menolak tanggal yang tidak ada (misalnya `2025-02-30`). Slot harus 1-8 dan `slot_mulai` tidak boleh lebih besar dari `slot_selesai`; kalau sama berarti satu slot. Jumlah minimal 1.

## Status awal peminjaman, B2
- Prompt: "Status awalnya kenapa langsung di-set `\"diajukan\"`?"
- Dipakai: `baru.Status = "diajukan"` dan `baru.Alasan = ""` di `buatPeminjaman`.
- Penjelasan: status tidak boleh dikirim client (tidak ada di `PeminjamanInput`), jadi server yang menentukan nilai awalnya. `Alasan` baru terisi kalau peminjaman ditolak (`TolakInput`).

## Menentukan filter, B3
- Prompt: "Lanjut B3. Saya butuh daftar peminjaman yang bisa difilter, plus detail satu peminjaman. Filternya sebaiknya apa saja?"
- Dipakai: tidak ada kode baru (rancangan saja).
- Penjelasan: filter `sumber_daya_id`, `tanggal`, dan `status`, semuanya lewat query string dan opsional, misalnya `GET /peminjaman?sumber_daya_id=1&tanggal=2025-10-20&status=diajukan`.

## Daftar peminjaman dengan filter, B3
- Prompt: "Oke, generate-kan `lihatSemuaPeminjaman`."
- Dipakai: `lihatSemuaPeminjaman` di `peminjaman.go`; route `GET /peminjaman` di `NewHandler` (`main.go`).
- Penjelasan: `sumber_daya_id` di-parse ke angka di luar lock (salah format langsung 400). Setelah itu kunci mutex, loop `daftarPinjam`, dan kumpulkan yang lolos filter.

## Loop dan `continue`, B3
- Prompt: "Jelaskan bagian loop-nya, kenapa pakai `continue` di tiap filter?"
- Dipakai: tidak ada kode baru (penjelasan saja).
- Penjelasan: tiap `if` artinya "kalau filter ini diisi dan datanya tidak cocok, lewati". Kondisi `filterX != ""` membuat filter yang kosong diabaikan. Hasilnya semua filter aktif digabung dengan logika AND.

## Slice kosong vs nil, B3
- Prompt: "Kenapa `hasil := []Peminjaman{}` dan bukan `var hasil []Peminjaman`?"
- Dipakai: tidak ada kode baru (penjelasan saja).
- Penjelasan: slice `nil` di-encode menjadi `null`, sedangkan `[]Peminjaman{}` menjadi `[]`. Array kosong lebih mudah diproses client. Ini juga alasan `NewHandler` mengisi `daftarSumber` dan `daftarPinjam` dengan slice kosong.

## Detail satu peminjaman, B3
- Prompt: "Terus untuk detail satu peminjaman gimana?"
- Dipakai: `lihatSatuPeminjaman` dan `cariPeminjaman` di `peminjaman.go`; route `GET /peminjaman/{id}` di `NewHandler` (`main.go`).
- Penjelasan: polanya sama dengan `lihatSatuSumberDaya`: `ambilID` untuk 400, `cariPeminjaman` untuk 404. `cariPeminjaman` mengembalikan index supaya bisa dipakai ulang di B6.
