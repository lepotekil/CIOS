# Use a newer version of Go
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy the Go module files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o botcrawler ./cmd/botcrawler/main.go

# Create a minimal runtime image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built application from the builder
COPY --from=builder /app/botcrawler .

# Copy the configuration file
COPY config/config.yaml ./config/config.yaml
COPY sql/ ./sql/

# Command to run the application
CMD ["./botcrawler"]