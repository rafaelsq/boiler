package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Fail writes the error message if debug is set.
func Fail(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if len(r.URL.Query()["debug"]) != 0 {
		fmt.Fprintf(w, message)
	}
}

// JSON writes the content of the param data as JSON.
func JSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
		Fail(w, r, http.StatusInternalServerError, "could not encode response")
	}
}
