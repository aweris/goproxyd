# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/goproxyd

# Final stage
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary from builder
COPY --from=builder /go/bin/goproxyd /goproxyd

# Set proper metadata
LABEL org.opencontainers.image.source="https://github.com/aweris/goproxyd" \
      org.opencontainers.image.description="A proxy server for Go modules" \
      org.opencontainers.image.licenses="MIT"

# Expose the port
EXPOSE 8080

# Use nonroot user (provided by distroless)
USER nonroot:nonroot

# Set the entrypoint
ENTRYPOINT ["/goproxyd"]
