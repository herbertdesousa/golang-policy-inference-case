package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type InferController struct {
	inferService *InferService
}

func NewInferController(inferService *InferService) *InferController {
	return &InferController{inferService: inferService}
}

func (c *InferController) HandleInfer(w http.ResponseWriter, r *http.Request) {
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

	resp, err := c.inferService.Evaluate(req)
	if err != nil {
		log.Printf("Failed to evaluate policy: %v", err)
		http.Error(w, "Failed to evaluate policy: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
