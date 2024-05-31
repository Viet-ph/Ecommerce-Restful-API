package helper

import (
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload any) error {
	err := Encode(w, statusCode, payload)
	return err
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}
