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

Benchmarked at **50 RPS** against a 6–7 node deep graph with ~4-7 conditions (with IN, NOT, math) per edge:

| Configuration    | p(90)  | p(95)  | p(99)  | avg    | min      | max     |
| ---------------- | ------ | ------ | ------ | ------ | -------- | ------- |
| CPU 0.5 + 250 MB | 2.40ms | 2.51ms | 3.06ms | 1.96ms | 981.81µs | 10.38ms |

Take a little higher with **200RPS**

| Configuration                  | p(90)    | p(95)    | p(99) | avg      | min      | max   |
| ------------------------------ | -------- | -------- | ----- | -------- | -------- | ----- |
| CPU 0.5 + 256 MB + 2 instances | 339.22ms | 619.07ms | 1.12s | 117.92ms | 300.81µs | 1.43s |
