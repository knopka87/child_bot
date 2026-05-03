########################
# Build stage
########################
FROM node:20-bookworm-slim AS build

WORKDIR /app

# Build arguments для Vite переменных окружения
ARG VITE_API_BASE_URL=/api
ARG VITE_APP_VERSION=1.0.0
ARG VITE_ANALYTICS_ENABLED=true
ARG VITE_VK_APP_ID
ARG VITE_MAX_APP_ID
ARG VITE_TELEGRAM_BOT_USERNAME

# Копируем package files для кеширования слоя зависимостей
COPY package*.json ./

# Устанавливаем зависимости с использованием npm ci (чище и быстрее)
RUN npm ci --prefer-offline --no-audit

# Копируем исходники
COPY . .

# Сборка production build
ENV NODE_ENV=production
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_APP_VERSION=${VITE_APP_VERSION}
ENV VITE_ANALYTICS_ENABLED=${VITE_ANALYTICS_ENABLED}
ENV VITE_VK_APP_ID=${VITE_VK_APP_ID}
ENV VITE_MAX_APP_ID=${VITE_MAX_APP_ID}
ENV VITE_TELEGRAM_BOT_USERNAME=${VITE_TELEGRAM_BOT_USERNAME}

RUN npm run build

########################
# Production stage with nginx
########################
FROM nginx:1.25-alpine

# Копируем собранные файлы
COPY --from=build /app/dist /usr/share/nginx/html

# Копируем nginx конфигурацию
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:80/ || exit 1

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
