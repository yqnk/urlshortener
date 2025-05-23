# syntax=docker/dockerfile:1
# Layer 1
FROM golang:1.24.3-alpine AS build-layer

WORKDIR /app

ENV CGO_ENABLED=1
ENV GOOS=linux

# required for go-sqlite3
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o urlshortener ./cmd/server

# Layer 2
FROM alpine:latest

WORKDIR /root/
COPY --from=build-layer /app/urlshortener .

EXPOSE 3434

CMD ["./urlshortener"]
