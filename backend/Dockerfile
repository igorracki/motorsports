FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8081
CMD ["./server"]
