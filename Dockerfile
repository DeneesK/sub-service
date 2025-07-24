# Stage 1: build
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/main.go -o api/docs

RUN go build -o app ./cmd/main.go

# Stage 2: runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY ./migrations .
COPY --from=builder /app/app .
COPY --from=builder /app/api/docs ./api/docs