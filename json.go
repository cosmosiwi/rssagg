package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseErr(w http.ResponseWriter, code int, msg string){
	if code > 499 {
		log.Println("Responding with 5XX error:" + msg)
	}
	type errResponse struct{
		Error string `json:"error"`
	}
	responseWithJSON(w, code, errResponse{Error: msg})
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}){
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshal JSON: %s\n", err)
	}
	w.WriteHeader(code)
	w.Write(dat)
}