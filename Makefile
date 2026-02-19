run:
	go run cmd/main.go

test-integration:
	go test -v -count=1 ./test

ci:
	@echo "up env"
	@echo "build and run app"
	@echo "run int tests"
