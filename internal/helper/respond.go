package helper

import (
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload any) error {
	// response, err := json.Marshal(payload)
	// if err != nil {
	// 	return err
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.WriteHeader(code)
	// w.Write(response)
	err := Encode(w, statusCode, payload)
	return err
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}
