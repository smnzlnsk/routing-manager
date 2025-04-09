# Build stage
FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

ARG VERSION
ARG COMMIT
ARG DATE

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -a -installsuffix cgo -o routing-manager ./cmd/server

# Final stage
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user to run the application
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the binary from the builder stage
COPY --from=builder /app/routing-manager /app/
COPY --from=builder /app/config.yaml /app/

# Set ownership of the application files
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Set environment variables
ENV LOG_FORMAT=console \
    LOG_LEVEL=info

# Command to run the application
ENTRYPOINT ["/app/routing-manager"]