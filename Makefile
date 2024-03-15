.PHONY: tests
tests:
	go test -v -race -cover ./...

.PHONY: mocks
mocks:
	@echo "=====> Installing mockgen"
	@go install github.com/golang/mock/mockgen@latest

	@echo "=====> Removing old mocks"
	@rm -rf mocks

	@echo "=====> Generating mocks"

	@for file in application/contract/*.go; do \
		filename=$$(basename $$file); \
		mockgen -package mocks -source=$$file -destination=mocks/$$filename; \
	done
	
	@echo "=====> Mocks generated"

.PHONY: start
start: build
	@echo "=====> Starting application"
	@./myapp 

.PHONY: build
build:
	@echo "=====> Building application"
	@go build -o myapp ./cmd/logparser/main.go
