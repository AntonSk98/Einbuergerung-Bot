# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install gcc for cgo
RUN apk add --no-cache gcc musl-dev

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and .env file
COPY . .

# Build the application with cgo enabled for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o /bot ./cmd/bot

# Final runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /bot /app/bot

# Run the bot
CMD ["/app/bot"]