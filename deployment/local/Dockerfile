# Stage 1: Build the Go application
FROM golang:1.23.1-alpine as builder
WORKDIR /app

COPY . .

ARG BUILD_TARGET
RUN go build -o /app/bin/app $BUILD_TARGET

# Stage 2: Create a smaller image for running the binary
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/bin/app .

ENTRYPOINT ["./app"]