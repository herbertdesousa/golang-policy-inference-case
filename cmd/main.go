package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type CompiledEdge struct {
	Dst     string
	Program *vm.Program
}

type PolicyEngine struct {
	Nodes     map[string]*gographviz.Node
	Adjacency map[string][]CompiledEdge
}

type InferRequest struct {
	PolicyDot string                 `json:"policy_dot"`
	Input     map[string]interface{} `json:"input"`
}

type InferResponse struct {
	Output map[string]interface{} `json:"output"`
}

func NewPolicyEngine(dotString string, env interface{}) (*PolicyEngine, error) {
	dotString = strings.ReplaceAll(dotString, "result=", "comment=")
	dotString = strings.ReplaceAll(dotString, "cond=", "label=")

	graphAst, err := gographviz.ParseString(dotString)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return nil, fmt.Errorf("analysis error: %w", err)
	}

	engine := &PolicyEngine{
		Nodes:     graph.Nodes.Lookup,
		Adjacency: make(map[string][]CompiledEdge),
	}

	for _, edge := range graph.Edges.Edges {
		condStr := edge.Attrs["label"]
		if condStr == "" {
			continue
		}

		cond := strings.Trim(condStr, "\"")

		program, err := expr.Compile(cond, expr.Env(env))
		if err != nil {
			return nil, fmt.Errorf("failed to compile condition '%s': %w", cond, err)
		}

		engine.Adjacency[edge.Src] = append(engine.Adjacency[edge.Src], CompiledEdge{
			Dst:     edge.Dst,
			Program: program,
		})
	}

	if hasCycle(engine.Adjacency) {
		return nil, fmt.Errorf("invalid policy graph: an infinite loop (cycle) was detected")
	}

	return engine, nil
}

func hasCycle(adjacency map[string][]CompiledEdge) bool {
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recursionStack[node] = true

		for _, edge := range adjacency[node] {
			if !visited[edge.Dst] {
				if dfs(edge.Dst) {
					return true
				}
			} else if recursionStack[edge.Dst] {
				return true
			}
		}

		recursionStack[node] = false
		return false
	}

	for node := range adjacency {
		if !visited[node] {
			if dfs(node) {
				return true
			}
		}
	}
	return false
}

func (e *PolicyEngine) Evaluate(startNode string, payload map[string]interface{}) (string, string) {
	currentNode := startNode

	for {
		outgoingEdges, hasEdges := e.Adjacency[currentNode]

		if !hasEdges || len(outgoingEdges) == 0 {
			break
		}

		moved := false
		for _, edge := range outgoingEdges {
			res, err := expr.Run(edge.Program, payload)

			if err != nil {
				log.Printf("Eval error on %s: %v", currentNode, err)
				continue
			}

			if isMatch, ok := res.(bool); ok && isMatch {
				currentNode = edge.Dst
				moved = true
				break
			}
		}

		if !moved {
			return currentNode, "error: stuck in graph, no conditions matched"
		}
	}

	node := e.Nodes[currentNode]
	resultStr := ""
	if node != nil && node.Attrs["comment"] != "" {
		resultStr = strings.Trim(node.Attrs["comment"], "\"")
	}

	return currentNode, resultStr
}

func enrichResult(resultStr string, input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		output[k] = v
	}

	if resultStr == "" {
		return output
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

	return output
}

func handleInfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	engine, err := NewPolicyEngine(req.PolicyDot, req.Input)
	if err != nil {
		http.Error(w, "Failed to initialize policy: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, resultStr := engine.Evaluate("start", req.Input)

	finalOutput := enrichResult(resultStr, req.Input)

	resp := InferResponse{Output: finalOutput}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/infer", handleInfer)

	port := ":8080"
	fmt.Printf("HTTP server running on port %s\n", port)
	fmt.Printf("Send POST to http://localhost%s/infer\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
