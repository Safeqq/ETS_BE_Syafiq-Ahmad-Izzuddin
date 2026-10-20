package main

import (
	"net/http"
	"sync"
)

// ===== tempat simpan data =====

type Server struct {
	mu              sync.Mutex
	daftarSumber    []SumberDaya
	daftarPinjam    []Peminjaman
	idSumberBerikut int
	idPinjamBerikut int
}

func NewHandler() http.Handler {
	s := &Server{}
	s.daftarSumber = []SumberDaya{}
	s.daftarPinjam = []Peminjaman{}
	s.idSumberBerikut = 1
	s.idPinjamBerikut = 1

	mux := http.NewServeMux()
	// B1
	mux.HandleFunc("POST /sumber-daya", s.buatSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.lihatSemuaSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.lihatSatuSumberDaya)

	// B2 dan B3
	mux.HandleFunc("POST /peminjaman", s.buatPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.lihatSemuaPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.lihatSatuPeminjaman)

	return mux
}

func main() {
	http.ListenAndServe(":8080", NewHandler())
}
