# ==================== BUILDER STAGE ====================
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN make build-prod

# ==================== RUNTIME STAGE ====================
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
WORKDIR /ap
COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrations ./migrations
RUN chown -R appuser:appgroup /app
USER appuser
ENTRYPOINT ["/app/auth-service"]

