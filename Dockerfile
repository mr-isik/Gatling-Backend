# Builder
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o gatling-backend ./cmd/server/main.go

# Final Stage
FROM alpine:latest

# Install necessary runtime dependencies like tzdata and ca-certificates
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the built executable from the builder stage
COPY --from=builder /app/gatling-backend .

# Copy configuration and migration files
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Expose the API port
EXPOSE 8080

# Run the application
CMD ["./gatling-backend"]
