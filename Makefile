run:
	go run cmd/main.go

test-integration:
	go test -v -count=1 ./test/integration

test-stress:
	k6 run test/stress/infer_stress.ts --summary-trend-stats="p(90),p(95),p(99),avg,min,max"

ci:
	@echo "up env"
	@echo "build and run app"
	@echo "run int tests"
