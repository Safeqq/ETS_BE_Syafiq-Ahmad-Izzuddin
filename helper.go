package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ===== fungsi bantu untuk kirim response =====

func kirimJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func kirimError(w http.ResponseWriter, status int, pesan string) {
	kirimJSON(w, status, map[string]string{"error": pesan})
}

// ambilID mengubah {id} di URL menjadi angka. ok bernilai false kalau id bukan angka.
func ambilID(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return 0, false
	}
	return id, true
}
