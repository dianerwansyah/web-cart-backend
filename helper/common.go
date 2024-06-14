package helper

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func ParseJSONBody(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func GenerateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	WriteJSONResponse(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error converting result to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
