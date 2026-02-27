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

| Configuration                                    | p(90)    | p(95)    | p(99)  | avg      | min      | max     |
| ------------------------------------------------ | -------- | -------- | ------ | -------- | -------- | ------- |
| CPU 0.5 + 250 MB                                 | 2.40ms   | 2.51ms   | 3.06ms | 1.96ms   | 981.81µs | 10.38ms |
| CPU 0.25 + 125 MB + 2 instances + instance cache | 900.58µs | 977.96µs | 1.63ms | 751.01µs | 357.13µs | 15.87ms |

Question: Option 2 make it worth?

- Infra cost
- Code base complexitity increase

Take a little higher with **200RPS**

| Configuration                                    | p(90)    | p(95)   | p(99)  | avg      | min      | max     |
| ------------------------------------------------ | -------- | ------- | ------ | -------- | -------- | ------- |
| CPU 0.25 + 125 MB + 2 instances                  | 3.88s    | 4.74s   | 6.55s  | 1.90s    | 1.16ms   | 11.85s  |
| CPU 0.25 + 125 MB + 2 instances + instance cache | 717.51µs | 789.5µs | 1.14ms | 560.69µs | 257.47µs | 17.18ms |

Take a little more higher... **1000 RPS**

| Configuration                                             | p(90)    | p(95)    | p(99)   | avg      | min      | max      |
| --------------------------------------------------------- | -------- | -------- | ------- | -------- | -------- | -------- |
| CPU 0.125 + 62.5 MB + 4 instances (nginx CPU 0.1 + 64 MB) | 264.73ms | 361.75ms | 568.2ms | 78.76ms  | 235.02µs | 983.16ms |
| CPU 0.125 + 62.5 MB + 4 instances (nginx CPU 0.2 + 96 MB) | 617.97µs | 773.14µs | 1.92ms  | 940.37µs | 235.58µs | 286.42ms |
