package main

import (
	"encoding/json"
	"net/http"
)

// ===== B1: kelola sumber daya =====

// POST /sumber-daya
func (s *Server) buatSumberDaya(w http.ResponseWriter, r *http.Request) {
	var input SumberDayaInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		kirimError(w, http.StatusBadRequest, "body harus berupa JSON yang valid")
		return
	}

	baru, pesan := validasiSumberDaya(input)
	if pesan != "" {
		kirimError(w, http.StatusBadRequest, pesan)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	baru.ID = s.idSumberBerikut
	s.idSumberBerikut++
	s.daftarSumber = append(s.daftarSumber, baru)

	kirimJSON(w, http.StatusCreated, baru)
}

// GET /sumber-daya
func (s *Server) lihatSemuaSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	kirimJSON(w, http.StatusOK, s.daftarSumber)
}

// GET /sumber-daya/{id}
func (s *Server) lihatSatuSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, ok := ambilID(r)
	if !ok {
		kirimError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	index := s.cariSumberDaya(id)
	if index == -1 {
		kirimError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}

	kirimJSON(w, http.StatusOK, s.daftarSumber[index])
}

// cariSumberDaya mengembalikan posisi sumber daya di slice, atau -1 kalau tidak ada.
// Harus dipanggil saat mutex sedang dikunci.
func (s *Server) cariSumberDaya(id int) int {
	for i := 0; i < len(s.daftarSumber); i++ {
		if s.daftarSumber[i].ID == id {
			return i
		}
	}
	return -1
}
