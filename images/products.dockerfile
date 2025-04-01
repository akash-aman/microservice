# Build stage
FROM golang:1.24 AS builder

# Set destination for COPY
WORKDIR /
# Copy the source code
COPY ./internal ./internal

WORKDIR /internal/pkg
RUN go mod download
WORKDIR /internal/services/products
RUN go mod download


# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /products /internal/services/products/cmd/main.go

# Final stage
FROM alpine:latest

# Install necessary certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from the build stage
COPY --from=builder /products /products
