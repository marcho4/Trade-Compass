# ==================== BUILDER STAGE ====================
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN make build

# ==================== RUNTIME STAGE ====================
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata curl
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/ai-service .
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8082
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8082/health || exit 1
ENTRYPOINT ["/app/ai-service"]
