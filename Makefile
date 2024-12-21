.PHONY: setup lint test format proto docker

# Environment setup
GO_BIN := $(shell go env GOPATH)/bin
PYTHON_BIN := $(HOME)/Library/Python/3.11/bin
export PATH := $(GO_BIN):$(PYTHON_BIN):$(PATH)

# Go related variables
GOFILES = $(shell find . -type f -name '*.go' -not -path "./api/types/*")

# Install all dependencies and set up development environment
setup: setup-env setup-go setup-python setup-hooks setup-buf setup-proto-tools

# Set up environment
setup-env:
	@echo "Setting up environment..."
	@mkdir -p $(GO_BIN)
	@echo "export PATH=$(GO_BIN):$(PYTHON_BIN):\$$PATH" > .env.local
	@echo "Please run 'source .env.local' after setup completes"

# Set up Go dependencies
setup-go:
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing Go tools..."
	GOBIN=$(GO_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0
	GOBIN=$(GO_BIN) go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installed Go tools to $(GO_BIN)"

# Set up Python dependencies
setup-python:
	@echo "Installing/upgrading pip and pre-commit..."
	python3 -m pip install --user --upgrade pip pre-commit

# Set up git hooks
setup-hooks:
	@echo "Setting up git hooks..."
	@$(PYTHON_BIN)/pre-commit install
	@mkdir -p .git/hooks
	@cp .git-hooks/pre-push .git/hooks/
	@chmod +x .git/hooks/pre-push
	@echo "Git hooks installed successfully"

# Set up buf
setup-buf:
	@echo "Installing buf..."
	@if ! command -v buf >/dev/null 2>&1; then \
		brew install bufbuild/buf/buf; \
	fi

# Set up protoc plugins
setup-proto-tools:
	@echo "Installing protoc plugins..."
	GOBIN=$(GO_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	GOBIN=$(GO_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	GOBIN=$(GO_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1
	GOBIN=$(GO_BIN) go install github.com/shengyjs/protoc-gen-go-grpc-mock@latest

# Lint the code
lint:
	@echo "Running golangci-lint..."
	@$(GO_BIN)/golangci-lint run

# Run tests
test:
	go test -v ./...

# Format code
format:
	@echo "Formatting Go files..."
	gofmt -w $(GOFILES)
	$(GO_BIN)/goimports -w $(GOFILES)

# Generate proto files
proto:
	@echo "Generating proto files..."
	@if ! command -v buf >/dev/null 2>&1; then \
		echo "buf is not installed. Running setup-buf..."; \
		$(MAKE) setup-buf; \
	fi
	./scripts/proto-gen.sh

# Build docker image
docker:
	docker-compose build

# Run locally
run:
	docker-compose up

# Clean up
clean:
	docker-compose down
	rm -f tables.db
