package test

import (
	"bytes"
	"encoding/json"
	"golang-policy-inference-case/test/integration/base"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InferTestSuite struct {
	base.BaseIntegrationTestSuite
}

func (s *InferTestSuite) readPolicyDotFile(policyDotFile string) string {
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFile)
	dotFilePath := filepath.Join(testDir, "..", "digraphs", policyDotFile)

	dotContent, err := os.ReadFile(dotFilePath)
	if err != nil {
		s.T().Fatalf("failed to read policy DOT file %s: %v", dotFilePath, err)
	}

	return string(dotContent)
}

func (s *InferTestSuite) Test_InferAPI_TableDriven() {
	type testCase struct {
		name          string
		policyDotFile string
		input         string
		expected      string
	}

	tests := []testCase{
		{
			name:          "basic approval",
			policyDotFile: "age.dot",
			input:         `{"age": 20}`,
			expected:      `{"output":{"age":20,"approved":true}}`,
		},
		{
			name:          "multiple conditions - approved",
			policyDotFile: "age_and_income.dot",
			input:         `{"age": 20, "income": 50000}`,
			expected:      `{"output":{"age":20,"income":50000,"tier":"prime"}}`,
		},
		{
			name:          "multiple conditions - rejected - income < 50000",
			policyDotFile: "age_and_income.dot",
			input:         `{"age": 20, "income": 49999}`,
			expected:      `{"output":{"age":20,"approved":false,"income":49999}}`,
		},
		{
			name:          "multiple conditions - rejected - age < 18",
			policyDotFile: "age_and_income.dot",
			input:         `{"age": 17, "income": 50000}`,
			expected:      `{"output":{"age":17,"approved":false,"income":50000}}`,
		},
		{
			name:          "complex policy - reject_minor",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 17, "score": 0, "income": 0, "country": "BR", "tenure_months": 0, "risk_flag": false, "has_default": false, "product_type": "credit"}`,
			expected:      `{"output":{"age":17,"approved":false,"country":"BR","has_default":false,"income":0,"product_type":"credit","reason":"underage","risk_flag":false,"score":0,"tenure_months":0}}`,
		},
		{
			name:          "complex policy - reject_low_score",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 400, "income": 0, "country": "BR", "tenure_months": 0, "risk_flag": false, "has_default": false, "product_type": "credit"}`,
			expected:      `{"output":{"age":30,"approved":false,"country":"BR","has_default":false,"income":0,"product_type":"credit","reason":"low_score","risk_flag":false,"score":400,"tenure_months":0}}`,
		},
		{
			name:          "complex policy - reject_low_income",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 550, "income": 3000, "country": "BR", "tenure_months": 0, "risk_flag": false, "has_default": false, "product_type": "credit"}`,
			expected:      `{"output":{"age":30,"approved":false,"country":"BR","has_default":false,"income":3000,"product_type":"credit","reason":"low_income","risk_flag":false,"score":550,"tenure_months":0}}`,
		},
		{
			name:          "complex policy - review_foreign",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 800, "income": 20000, "country": "US", "risk_flag": true, "tenure_months": 24, "has_default": false, "product_type": "credit"}`,
			expected:      `{"output":{"age":30,"approved":false,"country":"US","has_default":false,"income":20000,"product_type":"credit","reason":"foreign_risk","risk_flag":true,"score":800,"tenure_months":24}}`,
		},
		{
			name:          "complex policy - approve_prime",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 800, "income": 20000, "country": "BR", "tenure_months": 48, "risk_flag": false, "has_default": false, "product_type": "credit"}`,
			expected:      `{"output":{"age":30,"approved":true,"country":"BR","has_default":false,"income":20000,"limit_multiplier":"3","product_type":"credit","risk_flag":false,"score":800,"segment":"prime","tenure_months":48}}`,
		},
		{
			name:          "complex policy - approve_standard",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 700, "income": 10000, "country": "BR", "tenure_months": 24, "risk_flag": false, "has_default": false, "product_type": "credit"}`,
			expected:      `{"output":{"age":30,"approved":true,"country":"BR","has_default":false,"income":10000,"limit_multiplier":"2","product_type":"credit","risk_flag":false,"score":700,"segment":"standard","tenure_months":24}}`,
		},
		{
			name:          "complex policy - approve_basic via tenure_mid",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 700, "income": 10000, "country": "BR", "tenure_months": 24, "risk_flag": false, "has_default": false, "product_type": "loan"}`,
			expected:      `{"output":{"age":30,"approved":true,"country":"BR","has_default":false,"income":10000,"limit_multiplier":"1","product_type":"loan","risk_flag":false,"score":700,"segment":"basic","tenure_months":24}}`,
		},
		{
			name:          "complex policy - approve_basic via tenure_low",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 550, "income": 5000, "country": "BR", "tenure_months": 6, "risk_flag": false, "has_default": false, "product_type": "credit"}`,
			expected:      `{"output":{"age":30,"approved":true,"country":"BR","has_default":false,"income":5000,"limit_multiplier":"1","product_type":"credit","risk_flag":false,"score":550,"segment":"basic","tenure_months":6}}`,
		},
		{
			name:          "complex policy - approve_basic via foreign tenure_mid",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 800, "income": 20000, "country": "US", "risk_flag": false, "tenure_months": 24, "has_default": false, "product_type": "loan"}`,
			expected:      `{"output":{"age":30,"approved":true,"country":"US","has_default":false,"income":20000,"limit_multiplier":"1","product_type":"loan","risk_flag":false,"score":800,"segment":"basic","tenure_months":24}}`,
		},
		{
			name:          "complex policy - manual_review via tenure_high",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 800, "income": 20000, "country": "BR", "tenure_months": 48, "risk_flag": false, "has_default": true, "product_type": "credit"}`,
			expected:      `{"output":{"age":30,"approved":false,"country":"BR","has_default":true,"income":20000,"product_type":"credit","reason":"manual_review","risk_flag":false,"score":800,"tenure_months":48}}`,
		},
		{
			name:          "complex policy - manual_review via tenure_mid",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 700, "income": 10000, "country": "BR", "tenure_months": 24, "risk_flag": false, "has_default": true, "product_type": "credit"}`,
			expected:      `{"output":{"age":30,"approved":false,"country":"BR","has_default":true,"income":10000,"product_type":"credit","reason":"manual_review","risk_flag":false,"score":700,"tenure_months":24}}`,
		},
		{
			name:          "complex policy - manual_review via tenure_low",
			policyDotFile: "complex_policy.dot",
			input:         `{"age": 30, "score": 550, "income": 5000, "country": "BR", "tenure_months": 6, "risk_flag": false, "has_default": true, "product_type": "credit"}`,
			expected:      `{"output":{"age":30,"approved":false,"country":"BR","has_default":true,"income":5000,"product_type":"credit","reason":"manual_review","risk_flag":false,"score":550,"tenure_months":6}}`,
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			dotContent := s.readPolicyDotFile(tc.policyDotFile)

			var inputJSON map[string]interface{}
			if err := json.Unmarshal([]byte(tc.input), &inputJSON); err != nil {
				t.Fatalf("failed to parse input JSON: %v", err)
			}

			requestPayload := map[string]interface{}{
				"policy_dot": string(dotContent),
				"input":      inputJSON,
			}

			payloadBytes, err := json.Marshal(requestPayload)
			if err != nil {
				t.Fatalf("failed to marshal request payload: %v", err)
			}

			resp, err := s.HttpClient.Post("http://localhost:8080/infer", "application/json", bytes.NewBuffer(payloadBytes))
			if err != nil {
				t.Fatalf("failed to POST /infer: %v", err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			assert.Equal(t, tc.expected+"\n", string(body))
		})
	}
}

func TestInferIntegration(t *testing.T) {
	suite.Run(t, new(InferTestSuite))
}
