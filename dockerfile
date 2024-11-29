# Stage 1: Build environment
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Create a non-root user for building
RUN addgroup -S gogroup && adduser -S gouser -G gogroup

# Set working directory
WORKDIR /build

# Copy go mod files first to leverage layer caching
COPY go.mod go.sum ./

# Download dependencies
# Using go mod download instead of go get for better reproducibility
RUN go mod download && \
    go mod verify

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static linking
# -ldflags="-w -s" strips debug information and symbol tables
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /build/app ./main.go

# Stage 2: Production environment
FROM scratch AS production

# Import from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the binary
COPY --from=builder /build/app /app

# Copy any additional required files (like config files, static assets)
# COPY --from=builder /build/config /config
# COPY --from=builder /build/static /static

# Use non-root user
USER gouser

# Set environment variables
ENV APP_ENV=production

# Document the port
EXPOSE 3000

# Run the binary
ENTRYPOINT ["/app"]