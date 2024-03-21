# STAGE 1
FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /go-blog ./cmd/web

# STAGE 2
FROM alpine:latest

WORKDIR /app
COPY --from=builder /go-blog ./go-blog
ENTRYPOINT ["/app/go-blog"]