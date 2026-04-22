########################
# Build stage
########################
# Используем стабильный Debian-based образ (меньше проблем с сетью, чаще кешируется)
FROM golang:1.22-bookworm AS build
WORKDIR /src

# Добавляем retry для apt (на случай сбоев сети)
RUN for i in 1 2 3 4 5; do \
      apt-get update && \
      apt-get install -y --no-install-recommends ca-certificates tzdata && \
      rm -rf /var/lib/apt/lists/* && \
      break || sleep 10; \
    done

# Сначала зависимости — кэшируется отдельно
COPY go.mod go.sum ./

# go mod download с retry
RUN for i in 1 2 3; do \
      go mod download && break || sleep 10; \
    done

# Код
COPY api ./api

# Сборка REST API server (новый entrypoint)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/server ./api/cmd/server

# golang-migrate (postgres)
RUN GOBIN=/out go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

# Миграции и entrypoint
COPY api/migrations /out/migrations
COPY api/docker/entrypoint.sh /out/entrypoint.sh
RUN chmod +x /out/entrypoint.sh

########################
# Runtime stage (app-only)
########################
# Используем Debian-slim для лучшей совместимости с Chromium
FROM debian:bookworm-slim
WORKDIR /app

# Устанавливаем зависимости с retry, включая chromium для генерации PDF
RUN for i in 1 2 3 4 5; do \
      apt-get update && \
      apt-get install -y --no-install-recommends \
        ca-certificates \
        tzdata \
        bash \
        wget \
        chromium \
        fonts-liberation \
        fonts-noto-emoji \
        fonts-dejavu && \
      rm -rf /var/lib/apt/lists/* && \
      break || sleep 10; \
    done

# бинарь и миграции
COPY --from=build /out/server /app/server
COPY --from=build /out/migrate /usr/local/bin/migrate
COPY --from=build /out/migrations /app/migrations
COPY --from=build /out/entrypoint.sh /app/entrypoint.sh

ENV PORT=8080 \
    MIGRATIONS_DIR=/app/migrations

EXPOSE 8080
ENTRYPOINT ["/app/entrypoint.sh"]