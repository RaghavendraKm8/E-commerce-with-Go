# syntax=docker/dockerfile:1
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o service .

# final image
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/service /app/service

EXPOSE 8080
CMD ["./service"]
