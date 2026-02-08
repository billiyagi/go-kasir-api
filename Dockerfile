# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run Stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy environment file example if needed, but usually we rely on Env Vars in production
# COPY .env .env 

# Expose the port application runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
