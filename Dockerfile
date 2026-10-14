# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache make git

WORKDIR /src

# Cache dependencies
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /usr/local/bin/gitforge ./cmd/gitforge

# Final stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates git && \
    addgroup -g 1000 gitforge && \
    adduser -u 1000 -G gitforge -s /bin/sh -D gitforge

COPY --from=builder /usr/local/bin/gitforge /usr/local/bin/gitforge

USER gitforge
WORKDIR /home/gitforge

ENTRYPOINT ["gitforge"]
