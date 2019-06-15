package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Fail(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if len(r.URL.Query()["debug"]) != 0 {
		fmt.Fprintf(w, message)
	}
}

func JSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
		Fail(w, r, http.StatusInternalServerError, "could not encode response")
	}
}
