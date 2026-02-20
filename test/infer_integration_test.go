package test

import (
	"bytes"
	"encoding/json"
	"golang-policy-inference-case/test/base"
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
	dotFilePath := filepath.Join(testDir, "digraphs", policyDotFile)

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
