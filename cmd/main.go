package main

import (
	"fmt"
	"log"
	"strings"
	"time"

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

func NewPolicyEngine(dotString string, env interface{}) (*PolicyEngine, error) {
	// Bypass gographviz validation errors
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

func main() {
	dotString := `digraph Policy {
		start [result=""]
		check_income [result=""]
		approved [result="tier=prime"]
		rejected [result="approved=false"]
		
		start -> check_income [cond="age>=18"]
		start -> rejected [cond="age<18"]
		
		check_income -> approved [cond="income >= 50000"]
		check_income -> rejected [cond="income < 50000"]
	}`

	dummyEnv := map[string]interface{}{
		"age": 0, "income": 0,
	}

	fmt.Println("Initializing Engine...")
	engine, err := NewPolicyEngine(dotString, dummyEnv)
	if err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	payload := map[string]interface{}{
		"age":    25,
		"income": 60000,
	}

	fmt.Println("Evaluating Payload...")

	start := time.Now()
	finalNode, result := engine.Evaluate("start", payload)
	duration := time.Since(start)

	fmt.Printf("\nPayload:      %+v\n", payload)
	fmt.Printf("Landed Node:  %s\n", finalNode)
	fmt.Printf("Final Result: %s\n", result)
	fmt.Printf("Eval Time:    %s\n", duration)
}
