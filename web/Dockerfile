# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web-app .

# Final stage
FROM alpine:3.18

# Install necessary packages
RUN apk add --no-cache ca-certificates tzdata curl

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/web-app .

# Copy templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Create a non-root user and switch to it
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app
USER appuser

# Expose port
EXPOSE 3000

# Create health check script
USER root
RUN echo '#!/bin/sh' > /usr/local/bin/healthcheck.sh && \
    echo 'curl -s -f http://localhost:3000/health || exit 1' >> /usr/local/bin/healthcheck.sh && \
    chmod +x /usr/local/bin/healthcheck.sh
USER appuser

# Add health check
HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 CMD /usr/local/bin/healthcheck.sh

# Set environment variables
ENV PORT=3000

# Run the application
CMD ["./web-app"] 