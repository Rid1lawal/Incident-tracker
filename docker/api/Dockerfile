FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY api/go.mod api/go.sum ./

RUN go mod download

COPY api/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/server ./cmd/api

# ---------- Runtime Stage ----------

FROM alpine:3.23

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /build/server .

COPY --from=builder /app/migrations ./migrations

RUN chown appuser:appgroup /app/server

USER appuser

EXPOSE 8080

CMD ["./server"]