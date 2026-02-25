package main

import (
	"fmt"
	"golang-policy-inference-case/internal/api"
	"golang-policy-inference-case/internal/cache"
	"golang-policy-inference-case/internal/policy"
	"log"
	"net/http"
)

func main() {
	policyCache, err := cache.NewCache[policy.PolicyEngine](128)

	if err != nil {
		log.Fatalf("Failed to create policy cache: %v", err)
	}

	inferService := api.NewInferService(policyCache)
	inferController := api.NewInferController(inferService)

	http.HandleFunc("/infer", inferController.HandleInfer)

	port := ":8080"
	fmt.Printf("HTTP server running on port %s\n", port)
	fmt.Printf("Send POST to http://localhost%s/infer\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
