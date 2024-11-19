# Build stage
FROM golang:1.23.2-alpine AS builder

WORKDIR /opt

# Copy the source code
COPY .. .

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o aw-sync-agent .

# Final stage
FROM alpine:latest

WORKDIR /opt

# Copy the built binary from the builder stage
COPY --from=builder /opt/aw-sync-agent .

# Run the executable
CMD ["./aw-sync-agent"]