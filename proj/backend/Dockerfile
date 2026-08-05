# backend/Dockerfile

# Estágio 1: compilação
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Estágio 2: execução (imagem mínima)
FROM alpine:latest
RUN apk --no-cache add ca-certificates wget

# roda como usuário sem privilégios em vez de root
RUN adduser -D -H -u 10001 app

WORKDIR /app
COPY --from=builder /app/main .

# /data é o volume persistente do SQLite; precisa pertencer ao usuário app
RUN mkdir -p /data && chown -R app /data /app
USER app

EXPOSE 8080

# healthcheck bate no /healthz (sem auth) pra orquestrador saber se está vivo
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/healthz || exit 1

CMD ["./main"]
