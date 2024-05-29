package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	JsonError(w, err.Error(), http.StatusBadRequest)
}

func ReadJSON(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("invalid JSON input")
	}
	return nil
}

func JsonOk(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	WriteJSON(w, v)
}

func WriteJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, fmt.Sprintf("json encoding error: %v", err), http.StatusInternalServerError)
		return
	}
	WriteBytes(w, b)
}

func WriteBytes(w http.ResponseWriter, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		http.Error(w, fmt.Sprintf("write error: %v", err), http.StatusInternalServerError)
		return
	}
}

func JsonError(w http.ResponseWriter, errStr string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	WriteJSON(w, &jsonErr{Err: errStr})
}

type jsonErr struct {
	Err string `json:"err"`
}