# Build stage
FROM golang:1.23.2-alpine AS builder

WORKDIR /app

# Install required build tools
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose ports
EXPOSE 8080 9008 8081

CMD ["./server"]
