# Build stage
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first (better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application (correct path)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary and config
COPY --from=builder /app/main .
COPY --from=builder /app/files ./files

# Set ownership
RUN chown -R appuser:appgroup /app

# Use non-root user
USER appuser

# Set timezone
ENV TZ=Asia/Bangkok

EXPOSE 8080

CMD ["./main"]

