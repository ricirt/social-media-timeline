FROM golang:1.21-alpine AS builder

# Install build essentials
RUN apk add --no-cache git build-base

WORKDIR /build

# Copy only go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy only the necessary source code
COPY user/ ./user/

# Build the application
RUN cd user && CGO_ENABLED=0 GOOS=linux go build -o main cmd/user/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and config from builder
COPY --from=builder /build/user/main ./
COPY user/config/appSettings.json ./config/

EXPOSE 8080

CMD ["./main"]
