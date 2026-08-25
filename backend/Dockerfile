# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
ENV GOTOOLCHAIN=auto

# Install git and certs
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o adfluence-backend ./cmd/main.go

# Production runtime stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy compiled binary and migration files
COPY --from=builder /app/adfluence-backend ./adfluence-backend
COPY --from=builder /app/migration ./migration

# Default port
ENV PORT=8080
EXPOSE 8080

CMD ["./adfluence-backend"]
