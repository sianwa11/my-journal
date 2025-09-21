package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling json: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Fatal("Failed to write data to connection")
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {}