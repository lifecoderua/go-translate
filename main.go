package main

import (
	"engine/google"
	"fmt"
	"net/http"
	"strings"
)

func getURLParams(requestPath string) (string, string, string) {
	chunks := strings.SplitN(requestPath, "/", 5)
	return chunks[1], chunks[2], chunks[3] + "://" + chunks[4]
}

func handler(w http.ResponseWriter, r *http.Request) {
	sourceLang, targetLang, targetURL := getURLParams(r.URL.Path)

	// TODO: use a top-level generic  fetch/cleanup with handlers provied by an imported engine
	doc := google.Fetch(sourceLang, targetLang, targetURL)
	clean := google.Cleanup(doc)

	fmt.Fprintf(w, clean)
}

func main() {
	port := ":8888"
	fmt.Println("Started on", port)
	fmt.Println("Usage: http://localhost:8888/en/de/https/formkeep.com/")
	http.HandleFunc("/", handler)
	http.ListenAndServe(port, nil)
}
