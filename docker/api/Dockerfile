FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY api/go.mod api/go.sum ./

RUN go mod download

COPY api/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api