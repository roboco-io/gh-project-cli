# Build stage
FROM golang:1.23.4-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (required for Go modules)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ghp ./cmd/ghp

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates git

# Create non-root user
RUN adduser -D -s /bin/sh ghp

# Set working directory
WORKDIR /home/ghp

# Copy binary from builder stage
COPY --from=builder /app/ghp /usr/local/bin/ghp

# Change ownership
RUN chown ghp:ghp /usr/local/bin/ghp

# Switch to non-root user
USER ghp

# Set entrypoint
ENTRYPOINT ["ghp"]

# Default command
CMD ["--help"]