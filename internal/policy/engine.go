package policy

import (
	"fmt"
	"log"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type CompiledEdge struct {
	dst     string
	program *vm.Program
}

type PolicyEngine struct {
	nodes     map[string]*gographviz.Node
	adjacency map[string][]CompiledEdge
}

func NewPolicyEngine(dotString string, env interface{}) (*PolicyEngine, error) {
	dotString = strings.ReplaceAll(dotString, "result=", "comment=")
	dotString = strings.ReplaceAll(dotString, "cond=", "label=")
	dotString = strings.ReplaceAll(dotString, " NOT ", " not ")
	dotString = strings.ReplaceAll(dotString, " IN ", " in ")

	graphAst, err := gographviz.ParseString(dotString)
	if err != nil {
		log.Printf("Parse error: %v", err)
		return nil, fmt.Errorf("parse error: %w", err)
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		log.Printf("Analysis error: %v", err)
		return nil, fmt.Errorf("analysis error: %w", err)
	}

	engine := &PolicyEngine{
		nodes:     graph.Nodes.Lookup,
		adjacency: make(map[string][]CompiledEdge),
	}

	for _, edge := range graph.Edges.Edges {
		condStr := edge.Attrs["label"]
		if condStr == "" {
			continue
		}

		cond := strings.Trim(condStr, "\"")

		cond = strings.ReplaceAll(cond, `\"`, `"`)
		cond = strings.ReplaceAll(cond, `\`, `"`)

		program, err := expr.Compile(cond, expr.Env(env))
		if err != nil {
			log.Printf("Failed to compile condition '%s': %v", cond, err)
			return nil, fmt.Errorf("failed to compile condition '%s': %w", cond, err)
		}

		engine.adjacency[edge.Src] = append(engine.adjacency[edge.Src], CompiledEdge{
			dst:     edge.Dst,
			program: program,
		})
	}

	if hasCycle(engine.adjacency) {
		log.Printf("Invalid policy graph: an infinite loop (cycle) was detected")
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
			if !visited[edge.dst] {
				if dfs(edge.dst) {
					return true
				}
			} else if recursionStack[edge.dst] {
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
		outgoingEdges, hasEdges := e.adjacency[currentNode]

		if !hasEdges || len(outgoingEdges) == 0 {
			break
		}

		moved := false
		for _, edge := range outgoingEdges {
			res, err := expr.Run(edge.program, payload)

			if err != nil {
				log.Printf("Eval error on %s: %v", currentNode, err)
				continue
			}

			if isMatch, ok := res.(bool); ok && isMatch {
				currentNode = edge.dst
				moved = true
				break
			}
		}

		if !moved {
			return currentNode, "error: stuck in graph, no conditions matched"
		}
	}

	node := e.nodes[currentNode]
	resultStr := ""
	if node != nil && node.Attrs["comment"] != "" {
		resultStr = strings.Trim(node.Attrs["comment"], "\"")
	}

	return currentNode, resultStr
}
