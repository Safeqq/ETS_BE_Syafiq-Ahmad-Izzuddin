package main

import (
	"strings"
	"time"
	"unicode/utf8"
)

// ===== fungsi validasi =====
// Mengembalikan pesan error. Kalau pesan kosong ("") berarti data valid.

func validasiSumberDaya(input SumberDayaInput) (SumberDaya, string) {
	var hasil SumberDaya

	if input.Nama == nil || input.Jenis == nil || input.Stok == nil {
		return hasil, "nama, jenis, dan stok wajib diisi"
	}

	nama := strings.TrimSpace(*input.Nama)
	if nama == "" {
		return hasil, "nama tidak boleh kosong"
	}
	if utf8.RuneCountInString(nama) > 50 {
		return hasil, "nama maksimal 50 karakter"
	}

	jenis := strings.TrimSpace(*input.Jenis)
	if jenis != "ruang" && jenis != "alat" {
		return hasil, "jenis harus \"ruang\" atau \"alat\""
	}

	if *input.Stok < 1 {
		return hasil, "stok minimal 1"
	}

	hasil.Nama = nama
	hasil.Jenis = jenis
	hasil.Stok = *input.Stok
	return hasil, ""
}

// validasiPeminjaman memeriksa body POST dan PUT /peminjaman.
// Pengecekan yang butuh data sumber daya ada di cekSumberDayaDanJumlah.
func validasiPeminjaman(input PeminjamanInput) (Peminjaman, string) {
	var hasil Peminjaman

	if input.SumberDayaID == nil || input.Peminjam == nil || input.Tanggal == nil ||
		input.SlotMulai == nil || input.SlotSelesai == nil || input.Jumlah == nil {
		return hasil, "sumber_daya_id, peminjam, tanggal, slot_mulai, slot_selesai, dan jumlah wajib diisi"
	}

	peminjam := strings.TrimSpace(*input.Peminjam)
	if peminjam == "" {
		return hasil, "peminjam tidak boleh kosong"
	}
	if utf8.RuneCountInString(peminjam) > 50 {
		return hasil, "peminjam maksimal 50 karakter"
	}

	tanggal := strings.TrimSpace(*input.Tanggal)
	_, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return hasil, "tanggal harus berformat YYYY-MM-DD dan merupakan tanggal yang ada"
	}

	slotMulai := *input.SlotMulai
	slotSelesai := *input.SlotSelesai
	if slotMulai < 1 || slotMulai > 8 || slotSelesai < 1 || slotSelesai > 8 {
		return hasil, "slot harus antara 1 sampai 8"
	}
	if slotMulai > slotSelesai {
		return hasil, "slot_mulai tidak boleh lebih besar dari slot_selesai"
	}

	if *input.Jumlah < 1 {
		return hasil, "jumlah minimal 1"
	}

	hasil.SumberDayaID = *input.SumberDayaID
	hasil.Peminjam = peminjam
	hasil.Tanggal = tanggal
	hasil.SlotMulai = slotMulai
	hasil.SlotSelesai = slotSelesai
	hasil.Jumlah = *input.Jumlah
	return hasil, ""
}

// cekSumberDayaDanJumlah memeriksa sumber_daya_id menunjuk sumber daya yang ada
// dan jumlah tidak melebihi stoknya. Harus dipanggil saat mutex sedang dikunci.
func (s *Server) cekSumberDayaDanJumlah(p Peminjaman) string {
	index := s.cariSumberDaya(p.SumberDayaID)
	if index == -1 {
		return "sumber_daya_id tidak menunjuk sumber daya yang ada"
	}
	if p.Jumlah > s.daftarSumber[index].Stok {
		return "jumlah tidak boleh melebihi stok sumber daya"
	}
	return ""
}
