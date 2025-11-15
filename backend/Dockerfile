# Placeholder Dockerfile for Backend
# TODO: Replace with actual backend implementation

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Placeholder files - replace with actual backend code
RUN echo 'package main\n\nimport (\n\t"fmt"\n\t"net/http"\n)\n\nfunc main() {\n\thttp.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {\n\t\tw.WriteHeader(http.StatusOK)\n\t\tfmt.Fprint(w, "Backend is running")\n\t})\n\n\tfmt.Println("Backend server starting on :8080")\n\thttp.ListenAndServe(":8080", nil)\n}' > main.go

RUN go build -o backend main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/backend .

EXPOSE 8080

CMD ["./backend"]
