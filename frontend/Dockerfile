########################
# Build stage
########################
FROM node:20-alpine AS build

WORKDIR /app

# Build arguments для Vite переменных окружения
ARG VITE_API_BASE_URL=/api
ARG VITE_APP_VERSION=1.0.0
ARG VITE_ANALYTICS_ENABLED=true

# Копируем package files
COPY package*.json ./

# Устанавливаем ВСЕ зависимости (включая dev для сборки)
RUN npm ci

# Копируем исходники
COPY . .

# Сборка production build с переменными окружения
ENV NODE_ENV=production
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_APP_VERSION=${VITE_APP_VERSION}
ENV VITE_ANALYTICS_ENABLED=${VITE_ANALYTICS_ENABLED}
RUN npm run build

########################
# Production stage with nginx
########################
FROM nginx:1.25-alpine

# Копируем собранные файлы
COPY --from=build /app/dist /usr/share/nginx/html

# Копируем nginx конфигурацию
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Добавляем healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:80/ || exit 1

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
