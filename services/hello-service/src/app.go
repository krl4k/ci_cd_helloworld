package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World from K3s!\n")
}

func main() {
	http.HandleFunc("/", helloHandler)

	port := ":3000"
	fmt.Printf("Server running at http://0.0.0.0%s/\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
