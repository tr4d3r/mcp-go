# Multi-platform build
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build for the target architecture
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o mcp-server main.go

# Final stage
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary and static files
COPY --from=builder /app/mcp-server .
COPY --from=builder /app/static ./static

# Expose port
EXPOSE 8080

# Run the application
CMD ["./mcp-server"]
