.PHONY: build build-db build-server clean test fmt vet lint run-db run-server help

BINARY_DB=bin/db
BINARY_SERVER=bin/scrappah
CMD_DB=./cmd/db
CMD_SERVER=./cmd/scrappah

help:
	@echo "Available targets:"
	@echo "  build        - Build both executables"
	@echo "  build-db     - Build database helper tool"
	@echo "  build-server - Build main server"
	@echo "  run-db       - Run database helper (pass args with ARGS=...)"
	@echo "  run-server   - Run main server"
	@echo "  test-ip      - Test IP address through proxy"
	@echo "  test         - Run all tests"
	@echo "  fmt          - Format Go code"
	@echo "  vet          - Run go vet"
	@echo "  lint         - Run golangci-lint (if available)"
	@echo "  clean        - Remove built binaries"

build: build-db build-server

build-db:
	@mkdir -p bin
	go build -o $(BINARY_DB) $(CMD_DB)

build-server:
	@mkdir -p bin
	go build -o $(BINARY_SERVER) $(CMD_SERVER)

run-db: build-db
	$(BINARY_DB) $(ARGS)

run-server: build-server
	$(BINARY_SERVER)

test-ip: 
	@echo "Testing IP through proxy..."
	curl -x http://localhost:8002 https://api.ipify.org

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -cover ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet instead"; \
		go vet ./...; \
	fi

clean:
	rm -rf bin/

dev-setup:
	go mod tidy
	go mod download

.DEFAULT_GOAL := help