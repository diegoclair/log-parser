logpath ?= "./qgames.log"

.PHONY: start
start: build
	@echo "=====> Starting application"
	@./myapp --logpath=$(logpath)

.PHONY: build
build:
	@echo "=====> Building application"
	@go build -o myapp ./cmd/logparser/main.go

.PHONY: tests
tests:
	go test -v -race -cover ./...
