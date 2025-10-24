package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respWithError(w http.ResponseWriter, status int, msg string, err error) {
	if err != nil {
		log.Printf("Error: %s, Details: %v", msg, err)
	}

	if status > 499 {
		log.Printf("Internal Server Error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	respWithJSON(w, status, errorResponse{Error: msg})

}

func respWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(dat)
}
