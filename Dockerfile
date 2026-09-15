# Build Stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git and certificates
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/cinema-api ./cmd/api

# Final Minimal Stage
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/bin/cinema-api /app/cinema-api
COPY --from=builder /app/database /app/database
COPY --from=builder /app/.env.example /app/.env

EXPOSE 8080

CMD ["/app/cinema-api"]
