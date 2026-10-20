package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// ===== B2: buat peminjaman =====

// POST /peminjaman
func (s *Server) buatPeminjaman(w http.ResponseWriter, r *http.Request) {
	var input PeminjamanInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		kirimError(w, http.StatusBadRequest, "body harus berupa JSON yang valid")
		return
	}

	baru, pesan := validasiPeminjaman(input)
	if pesan != "" {
		kirimError(w, http.StatusBadRequest, pesan)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	pesan = s.cekSumberDayaDanJumlah(baru)
	if pesan != "" {
		kirimError(w, http.StatusBadRequest, pesan)
		return
	}

	baru.ID = s.idPinjamBerikut
	s.idPinjamBerikut++
	baru.Status = "diajukan"
	baru.Alasan = ""
	s.daftarPinjam = append(s.daftarPinjam, baru)

	kirimJSON(w, http.StatusCreated, baru)
}

// ===== B3: lihat peminjaman =====

// GET /peminjaman?sumber_daya_id=&tanggal=&status=
func (s *Server) lihatSemuaPeminjaman(w http.ResponseWriter, r *http.Request) {
	filterSumber := strings.TrimSpace(r.URL.Query().Get("sumber_daya_id"))
	filterTanggal := strings.TrimSpace(r.URL.Query().Get("tanggal"))
	filterStatus := strings.TrimSpace(r.URL.Query().Get("status"))

	sumberID := 0
	if filterSumber != "" {
		angka, err := strconv.Atoi(filterSumber)
		if err != nil {
			kirimError(w, http.StatusBadRequest, "sumber_daya_id harus berupa angka")
			return
		}
		sumberID = angka
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	hasil := []Peminjaman{}
	for _, p := range s.daftarPinjam {
		if filterSumber != "" && p.SumberDayaID != sumberID {
			continue
		}
		if filterTanggal != "" && p.Tanggal != filterTanggal {
			continue
		}
		if filterStatus != "" && p.Status != filterStatus {
			continue
		}
		hasil = append(hasil, p)
	}

	kirimJSON(w, http.StatusOK, hasil)
}

// GET /peminjaman/{id}
func (s *Server) lihatSatuPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := ambilID(r)
	if !ok {
		kirimError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	index := s.cariPeminjaman(id)
	if index == -1 {
		kirimError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}

	kirimJSON(w, http.StatusOK, s.daftarPinjam[index])
}

// cariPeminjaman mengembalikan posisi peminjaman di slice, atau -1 kalau tidak ada.
// Harus dipanggil saat mutex sedang dikunci.
func (s *Server) cariPeminjaman(id int) int {
	for i := 0; i < len(s.daftarPinjam); i++ {
		if s.daftarPinjam[i].ID == id {
			return i
		}
	}
	return -1
}
