# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO disabled - pure Go SQLite, no GCC needed
RUN CGO_ENABLED=0 GOOS=linux go build -o ticket-system ./cmd/main.go

# Stage 2: Minimal runtime
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ticket-system .

EXPOSE 8080

CMD ["./ticket-system"]
