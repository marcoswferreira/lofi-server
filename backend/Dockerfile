# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

# Install swag for documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Generate swagger documentation
RUN /go/bin/swag init

RUN go build -o lofi-server main.go

# Run stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/lofi-server .
# Copy generated docs to the runtime image
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8081

CMD ["./lofi-server"]
