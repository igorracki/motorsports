FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache wget
# Ensure the /data directory exists for the volume
RUN mkdir -p /data
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD sh -c "wget --no-verbose --tries=1 --spider http://localhost:${SERVER_PORT:-8080}/health || exit 1"
CMD ["./server"]
