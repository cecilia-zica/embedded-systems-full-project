# backend/Dockerfile

# Stage 1: build
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: runtime (minimal image)
# pinned image (not 'latest') for reproducible builds
FROM alpine:3.20
RUN apk --no-cache add ca-certificates wget

# run as an unprivileged user instead of root
RUN adduser -D -H -u 10001 app

WORKDIR /app
COPY --from=builder /app/main .

# /data is the persistent SQLite volume; it must be owned by the app user
RUN mkdir -p /data && chown -R app /data /app
USER app

EXPOSE 8080

# healthcheck hits /healthz (no auth) so the orchestrator knows it's alive
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/healthz || exit 1

CMD ["./main"]
