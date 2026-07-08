FROM golang:1.26-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

# Copy only dependency files first for caching
COPY go-executor/go.mod go-executor/go.sum ./
RUN go mod download

# Copy the source code
COPY go-executor/ .

# Build a static binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server .

FROM alpine:3.20
WORKDIR /app

# Only docker-cli — the Go binary needs this to talk to the host's docker.sock
RUN apk add --no-cache docker-cli

# Copy ONLY the compiled binary from the builder stage
COPY --from=builder /app/server .

EXPOSE 8050

CMD ["./server"]