package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const defaultVersion = "v0.0.0"

func getVersion() string {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		return defaultVersion
	}
	return version
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World from K3s!\n")
}

func main() {
	version := getVersion()
	log.Printf("Starting hello-service version: %s", version)

	http.HandleFunc("/", helloHandler)

	port := ":3000"
	fmt.Printf("Server running at http://0.0.0.0%s/\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
