package main

import (
	"engine/google"
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	baseURL := "https://formkeep.com/"
	targetURL := baseURL + r.URL.Path[1:]

	// TODO: use a top-level generic  fetch/cleanup with handlers provied by an imported engine
	doc := google.Fetch(targetURL)
	clean := google.Cleanup(doc)

	fmt.Fprintf(w, clean)
}

func main() {
	port := ":8888"
	fmt.Println("Started on", port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(port, nil)
}
