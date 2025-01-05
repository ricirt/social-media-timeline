FROM golang:1.21-alpine AS builder

# Debug için gerekli araçları yükle
RUN apk add --no-cache git gcc libc-dev

# Build aşaması için /app dizini
WORKDIR /app

# Delve debugger'ı yükle
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Önce go.mod ve go.sum dosyalarını kopyala
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY user/ ./user/

# Debug sembolleriyle birlikte build et
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o main ./user/cmd/user/main.go

# Yeni aşama - çalışma zamanı
FROM alpine:latest

# Debug için gerekli araçları yükle
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Builder aşamasından gerekli dosyaları kopyala
COPY --from=builder /app/main .
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

# Config dizinini oluştur
RUN mkdir -p /app/config

EXPOSE 8080 40000

# Debug modunda başlat
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]
