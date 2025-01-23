package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
	w.Write([]byte(msg))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snippetCreate"))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the server "))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/snippet/view/{id}", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	log.Println("starting ther server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
