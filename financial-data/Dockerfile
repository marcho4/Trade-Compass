# ==================== BUILDER STAGE ====================
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o financial-data cmd/main.go

# ==================== RUNTIME STAGE ====================
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata curl
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/financial-data .
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8082
ENTRYPOINT ["/app/financial-data"]
