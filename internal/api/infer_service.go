package api

import (
	"fmt"
	"golang-policy-inference-case/internal/cache"
	"golang-policy-inference-case/internal/policy"
	"log"
)

type InferService struct {
	policyCache *cache.Cache[policy.PolicyEngine]
}

func NewInferService(policyCache *cache.Cache[policy.PolicyEngine]) *InferService {
	return &InferService{policyCache: policyCache}
}

func (s *InferService) Evaluate(req InferRequestDto) (InferResponseDto, error) {
	// Could be a round-trip middleware http debugger
	log.Printf("Evaluating policy: %v with input: %v", req.PolicyDot, req.Input)

	cachedEngine, ok := s.policyCache.Get(req.PolicyDot)
	if ok {
		_, resultStr := cachedEngine.Evaluate("start", req.Input)
		return NewInferResponseDto(resultStr, req.Input), nil
	}

	log.Printf("Cache miss for policy, compiling %v", req.PolicyDot)

	engine, err := policy.NewPolicyEngine(req.PolicyDot, req.Input)
	if err != nil {
		return InferResponseDto{}, fmt.Errorf("failed to initialize policy: %w", err)
	}

	s.policyCache.Add(req.PolicyDot, engine)

	_, resultStr := engine.Evaluate("start", req.Input)

	return NewInferResponseDto(resultStr, req.Input), nil
}
