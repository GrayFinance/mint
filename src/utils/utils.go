package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

type JSONError struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func RandomHex(nbytes int) string {
	b := make([]byte, nbytes)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func SendJSONError(w http.ResponseWriter, code int, message string) {
	data, _ := json.Marshal(JSONError{Ok: false, Message: message})
	http.Error(w, string(data), code)
}

func SendJSONResponse(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
