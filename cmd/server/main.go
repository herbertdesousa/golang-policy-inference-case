package main

import (
	"fmt"
	"golang-policy-inference-case/internal/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/infer", api.HandleInfer)

	port := ":8080"
	fmt.Printf("HTTP server running on port %s\n", port)
	fmt.Printf("Send POST to http://localhost%s/infer\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
