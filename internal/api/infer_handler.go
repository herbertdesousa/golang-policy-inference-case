package api

import (
	"encoding/json"
	"golang-policy-inference-case/internal/policy"
	"log"
	"net/http"
)

func HandleInfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InferRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode json: %v", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	engine, err := policy.NewPolicyEngine(req.PolicyDot, req.Input)
	if err != nil {
		log.Printf("Failed to initialize policy: %v", err)
		http.Error(w, "Failed to initialize policy: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Kind high log, but useful to track policy <> input
	log.Printf("Evaluating policy: %v with input: %v", req.PolicyDot, req.Input)

	_, resultStr := engine.Evaluate("start", req.Input)

	resp := NewInferResponseDto(resultStr, req.Input)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
