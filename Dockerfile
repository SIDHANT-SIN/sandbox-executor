# ---------- Builder ----------
FROM golang:1.26-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

# Copy only dependency files first for caching
COPY go-executor/go.mod go-executor/go.sum ./
RUN go mod download

# Copy the source code
COPY go-executor/ .

# BUILD the binary here. This creates a file named 'server'
RUN go build -o server .

# ---------- Runtime ----------
FROM golang:1.26-alpine
WORKDIR /app

# Install docker-cli so the Go binary can talk to the host socket
RUN apk add --no-cache git docker-cli

# Copy ONLY the compiled binary from the builder stage
COPY --from=builder /app/server .

EXPOSE 8080

# Run the binary directly (No more 'go run'!)
CMD ["./server"]