# golang-policy-inference-case

HTTP service that evaluates decision policies defined as **DOT digraphs**. Each request provides a graph and an input payload; the engine traverses the graph evaluating edge conditions (via [expr-lang](https://github.com/expr-lang/expr)) and returns the result of the terminal node.

## Folder Structure

```
cmd/server/          # Entrypoint (main.go)
internal/
  api/               # HTTP handler + DTOs
  policy/            # Graph parsing & evaluation engine
test/
  digraphs/          # Sample .dot policy files
  integration/       # Integration tests (Go)
  stress/            # Stress tests (k6)
```

## Makefile

| Command                 | Description                            |
| ----------------------- | -------------------------------------- |
| `make run`              | Run the server locally                 |
| `make build`            | Build the Docker image                 |
| `make up`               | Start the container via docker-compose |
| `make build-and-up`     | Build image then start container       |
| `make test-integration` | Run integration tests                  |
| `make test-stress`      | Run k6 stress test                     |

## API

```
POST /infer
Content-Type: application/json

{
  "policy_dot": "<DOT digraph string>",
  "input": { "<key>": "<value>", ... }
}
```

## Performance

Benchmarked at **30 RPS** against a 6–7 node deep graph with ~2 conditions per edge:

| Resources          | p(90)   | p(95)   | p(99)   | avg     | min      | max      |
| ------------------ | ------- | ------- | ------- | ------- | -------- | -------- |
| 0.5 CPU / 250M mem | 1.60 ms | 1.70 ms | 2.00 ms | 1.23 ms | 666.27µs | 16.12 ms |
| 0.1 CPU / 125M mem | 1.60 ms | 1.72 ms | 1.98 ms | 1.25 ms | 664.85µs | 35.95 ms |
| 0.1 CPU / 250M mem | 1.58 ms | 1.66 ms | 1.94 ms | 1.23 ms | 674.66µs | 21.25 ms |
| 0.5 CPU / 125M mem | 1.61 ms | 1.73 ms | 2.03 ms | 1.24 ms | 681.71µs | 15.24 ms |

## Roadmap

- [x] Core engine (graph traversal + expression evaluation)
- [x] HTTP server
- [x] Integration tests
- [x] Stress tests
- [x] Performance improvements & code readability
