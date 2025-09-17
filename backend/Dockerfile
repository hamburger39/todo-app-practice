# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates and sqlite
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create directory for database
RUN mkdir -p /var/lib/todo

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
