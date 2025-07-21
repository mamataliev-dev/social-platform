# syntax=docker/dockerfile:1

############################
# 1. Build stage
############################
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 1.1 Pull in your module definitions (from repo root)
COPY go.mod go.sum ./
RUN go mod download

# 1.2 Copy the entire repo into the image
COPY . .

# 1.3 Build your service binary
WORKDIR /app/services/user-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/bin/user-service ./cmd/server

############################
# 2. Runtime stage
############################
FROM alpine:latest

# 2.1 TLS CA certs + gettext for envsubst
RUN apk add --no-cache ca-certificates gettext

# 2.2 Non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# 2.3 Copy in the binary and the YAML template from the builder
COPY --from=builder /app/bin/user-service    /app/user-service
COPY --from=builder /app/services/user-service/config.yaml /app/config.yaml.tpl

# 2.4 Chown & switch
RUN chown -R appuser:appgroup /app
USER appuser

# 2.5 Expose your gRPC & HTTP ports
EXPOSE 50000 10

# 2.6 Render the real config.yaml and launch
ENTRYPOINT ["sh", "-c", "\
    envsubst < /app/config.yaml.tpl > /app/config.yaml && \
    exec /app/user-service --config /app/config.yaml \
"]
