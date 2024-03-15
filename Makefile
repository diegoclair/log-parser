.PHONY: tests
tests:
	go test -v -race -cover ./...

.PHONY: start
start: build
	@echo "=====> Starting application"
	@./myapp 

.PHONY: build
build:
	@echo "=====> Building application"
	@go build -o myapp ./cmd/logparser/main.go
