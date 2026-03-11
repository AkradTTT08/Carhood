# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and static assets
COPY --from=builder /app/main .
RUN mkdir uploads

EXPOSE 8081

# Command to run the application
CMD ["./main"]
