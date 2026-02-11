# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /opt

# Copy the source code
COPY . .

# Download and install the dependencies
RUN go mod download

# Build the Go app
RUN go build -o aw-sync-agent .

# Final stage
FROM alpine:latest

WORKDIR /opt

# Copy the built binary from the builder stage
COPY --from=builder /opt/aw-sync-agent .

# Run the executable
CMD ["./aw-sync-agent"]