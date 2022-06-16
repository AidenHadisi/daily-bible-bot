# Create a builder for building initial binary
FROM golang:1.18-alpine as builder

WORKDIR /app

# Download deps
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy rest of the files
COPY . .

# Build the binary
RUN go build -o twitterApp ./cmd/twitter

# Now use a small unix image
FROM alpine:latest

WORKDIR /app/binary

# now copy from above container to this one
COPY --from=builder /app/twitterApp /app/binary

# Start it
CMD ["/app/binary/twitterApp"]