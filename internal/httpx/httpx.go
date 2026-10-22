package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

func Decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func PathString(r *http.Request, name string) string {
	return mux.Vars(r)[name]
}

func PathInt(r *http.Request, name string) (int64, bool) {
	n, err := strconv.ParseInt(mux.Vars(r)[name], 10, 64)
	return n, err == nil
}
