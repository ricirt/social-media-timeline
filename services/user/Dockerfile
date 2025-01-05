FROM golang:1.21-alpine AS builder

# Install build essentials
RUN apk add --no-cache git build-base

WORKDIR /build

# Copy the entire project for correct module resolution
COPY . .

# Build the application
WORKDIR /build/services/user
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and config from builder
COPY --from=builder /build/services/user/main ./
COPY --from=builder /build/services/user/config/appSettings.json ./config/

EXPOSE 8080

CMD ["./main"]
