package api

import "strings"

type InferRequestDto struct {
	PolicyDot string                 `json:"policy_dot"`
	Input     map[string]interface{} `json:"input"`
}

type InferResponseDto struct {
	Output map[string]interface{} `json:"output"`
}

func NewInferResponseDto(resultStr string, input map[string]interface{}) InferResponseDto {
	output := make(map[string]interface{})

	for k, v := range input {
		output[k] = v
	}

	if resultStr == "" {
		return InferResponseDto{Output: output}
	}

	pairs := strings.Split(resultStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])

			if val == "true" {
				output[key] = true
			} else if val == "false" {
				output[key] = false
			} else {
				output[key] = val
			}
		}
	}

	return InferResponseDto{Output: output}
}
