# =============================================================================
# DeFi Yield Aggregator - Multi-Stage Dockerfile
# =============================================================================
# This Dockerfile uses multi-stage builds to create optimized images for
# both development (with hot reload) and production (minimal binary).

# -----------------------------------------------------------------------------
# Stage 1: Base - Common Go setup
# -----------------------------------------------------------------------------
FROM golang:1.21-alpine AS base

# Install essential packages
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go module files first for better caching
COPY go.mod go.sum* ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download

# -----------------------------------------------------------------------------
# Stage 2: Development - With hot reload using Air
# -----------------------------------------------------------------------------
FROM base AS development

# Install Air for hot reload (https://github.com/cosmtrek/air)
RUN go install github.com/cosmtrek/air@latest

# Install additional development tools
RUN apk add --no-cache curl

# Source code will be mounted as volume, not copied
# This allows hot reload to work properly

# Default command (can be overridden in docker-compose)
CMD ["air", "-c", ".air.toml"]

# -----------------------------------------------------------------------------
# Stage 3: Builder - Compile the application
# -----------------------------------------------------------------------------
FROM base AS builder

# Copy all source code
COPY . .

# Build arguments for versioning
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the API server binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /bin/api-server ./cmd/server

# Build the worker binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /bin/worker ./cmd/worker

# -----------------------------------------------------------------------------
# Stage 4: Production API Server - Minimal runtime image
# -----------------------------------------------------------------------------
FROM alpine:3.19 AS production-api

# Install CA certificates for HTTPS requests and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /bin/api-server /app/api-server

# Copy any required static files or configs if needed
# COPY --from=builder /app/docs /app/docs

# Use non-root user
USER appuser

# Expose API port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/v1/health || exit 1

# Run the server
ENTRYPOINT ["/app/api-server"]

# -----------------------------------------------------------------------------
# Stage 5: Production Worker - Minimal runtime image
# -----------------------------------------------------------------------------
FROM alpine:3.19 AS production-worker

# Install CA certificates for HTTPS requests and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /bin/worker /app/worker

# Use non-root user
USER appuser

# Run the worker
ENTRYPOINT ["/app/worker"]
