# Stage 1: Modules caching
FROM golang:1.23 as modules
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Build
FROM golang:1.23 as builder
WORKDIR /app
COPY --from=modules /go/pkg /go/pkg
COPY . .

# Build your app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o pantan ./cmd/main.go


# Stage 3: Final
FROM ubuntu:24.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# Copy binaries
COPY --from=builder /app/pantan /app/pantan

WORKDIR /app
CMD ["/app/pantan"]