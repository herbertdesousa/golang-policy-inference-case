package test

import (
	"bytes"
	"golang-policy-inference-case/test/base"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InferTestSuite struct {
	base.BaseIntegrationTestSuite
}

func (s *InferTestSuite) Test_InferAPI_TableDriven() {
	type testCase struct {
		name      string
		policyDot string
		input     string
		expected  string
	}

	tests := []testCase{
		{
			name:      "basic approval",
			policyDot: `digraph { start [result=\"\"]; ok [result=\"approved=true\"]; no [result=\"approved=false\"]; start -> ok [cond=\"age>=18\"]; start -> no [cond=\"age<18\"]; }`,
			input:     `{"age": 20}`,
			expected:  `{"output":{"age":20,"approved":true}}`,
		},
		{
			name:      "multiple conditions - approved",
			policyDot: `digraph Policy { start [result=\"\"] check_income [result=\"\"] approved [result=\"tier=prime\"] rejected [result=\"approved=false\"] start -> check_income [cond=\"age>=18\"] start -> rejected [cond=\"age<18\"] check_income -> approved [cond=\"income >= 50000\"] check_income -> rejected [cond=\"income < 50000\"] }`,
			input:     `{"age": 20, "income": 50000}`,
			expected:  `{"output":{"age":20,"income":50000,"tier":"prime"}}`,
		},
		{
			name:      "multiple conditions - rejected - income < 50000",
			policyDot: `digraph Policy { start [result=\"\"] check_income [result=\"\"] approved [result=\"tier=prime\"] rejected [result=\"approved=false\"] start -> check_income [cond=\"age>=18\"] start -> rejected [cond=\"age<18\"] check_income -> approved [cond=\"income >= 50000\"] check_income -> rejected [cond=\"income < 50000\"] }`,
			input:     `{"age": 20, "income": 49999}`,
			expected:  `{"output":{"age":20,"approved":false,"income":49999}}`,
		},
		{
			name:      "multiple conditions - rejected - age < 18",
			policyDot: `digraph Policy { start [result=\"\"] check_income [result=\"\"] approved [result=\"tier=prime\"] rejected [result=\"approved=false\"] start -> check_income [cond=\"age>=18\"] start -> rejected [cond=\"age<18\"] check_income -> approved [cond=\"income >= 50000\"] check_income -> rejected [cond=\"income < 50000\"] }`,
			input:     `{"age": 17, "income": 50000}`,
			expected:  `{"output":{"age":17,"approved":false,"income":50000}}`,
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			payload := `
			{
				"policy_dot": "` + tc.policyDot + `",
				"input": ` + tc.input + `
			}
			`

			resp, err := s.HttpClient.Post("http://localhost:8080/infer", "application/json", bytes.NewBuffer([]byte(payload)))
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
