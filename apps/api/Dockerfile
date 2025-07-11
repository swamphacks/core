
### STAGE 1: Build the Go binary ###
FROM golang:1.24-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache build-base ca-certificates

# Set working directory
WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/api

### STAGE 2: Minimal runtime container ###
FROM alpine:latest

# Install CA certs for HTTPS support
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/server .

# Run the binary
CMD ["./server"]
