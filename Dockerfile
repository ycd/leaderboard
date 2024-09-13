# Build stage
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/leaderboard-api cmd/api/main.go

# Final stage
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /bin/leaderboard-api /bin/leaderboard-api

# Set environment variables
ENV PORT=8080

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/bin/leaderboard-api"]
