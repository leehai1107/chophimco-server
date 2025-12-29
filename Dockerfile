# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /usr/src/app

# Install build dependencies
RUN apk add --no-cache git

# Install swag for swagger documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate swagger docs
RUN swag init --parseDependency --parseInternal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/docs ./docs

EXPOSE 8081

CMD ["./main", "api"]

