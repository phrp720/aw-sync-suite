FROM golang:1.23.2-alpine

WORKDIR /opt

# Copy the source code
COPY .. .

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o aw-sync-agent .

# Run the executable
CMD ["./aw-sync-agent"]