# Platform Service

## Prerequisites
- Go 1.23.2 or later
- Docker & Docker Compose
- Python 3.8 or later (for pre-commit hooks)
- Git
- `buf` CLI tool (for proto generation)

## Quick Start

1. Clone the repository: git clone https://github.com/skip-mev/platform-take-home.git
2. Navigate to the project directory: cd platform-take-home


# Build and Run Guide

## Initial Setup
Run the following command to set up your development environment: `make setup`
This will:
1. Set up environment paths (`setup-env`)
   - Creates necessary directories
   - Sets up PATH variables

2. Install Go dependencies (`setup-go`)
   - Downloads project dependencies
   - Installs golangci-lint
   - Installs goimports

3. Install Python tools (`setup-python`)
   - Upgrades pip
   - Installs pre-commit

4. Configure Git hooks (`setup-hooks`)
   - Installs pre-commit hooks
   - Sets up pre-push hook for proto generation

5. Install Protocol Buffer tools
   - Installs `buf` CLI (`setup-buf`)
   - Installs protoc plugins (`setup-proto-tools`)

## Building the Application
To build the Docker image:
`make docker`
This command:
- Runs `docker-compose build`
- Creates the application image with:
  - gRPC server (port 9008)
  - REST gateway (port 8080)
  - Metrics endpoint (port 8081)

## Running the Services
To start all services:
`make run`
This command:
1. Starts the Postgres database
2. Launches the application container
3. Exposes all service endpoints:
   - `localhost:8080` - REST API
   - `localhost:9008` - gRPC Server
   - `localhost:8081` - Metrics

To stop all services:
`make clean`
