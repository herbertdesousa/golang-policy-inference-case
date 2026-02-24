run:
	go run cmd/server/main.go

test-integration:
	go test -v -count=1 ./test/integration

test-stress:
	k6 run test/stress/infer_stress.ts --summary-trend-stats="p(90),p(95),p(99),avg,min,max"

build:
	docker build -t golang-policy-inference-case .

up:
	@if [ -n "$$(docker ps -q -f name=golang-policy-inference-case)" ]; then \
		docker-compose down; \
	fi
	docker-compose up -d;

build-and-up:
	make build;
	make up;

ci:
	@echo "up env"
	@echo "build and run app"
	@echo "run int tests"
