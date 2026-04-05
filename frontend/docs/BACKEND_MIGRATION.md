# Backend API Migration Guide

## Обзор

Этот документ описывает процесс миграции фронтенда с моковых данных на реальный бэкенд API.

---

## Текущее состояние

### API Client
Файл: `src/api/client.ts`

```typescript
export const apiClient = {
  async get<T>(url: string): Promise<T> {
    // Currently returns mock data
  },
  async post<T>(url: string, data: any): Promise<T> {
    // Currently returns mock data
  },
  // ... other methods
}
```

### Endpoints реализованы:
- ✅ `/profile/*` - Profile API
- ✅ `/achievements/*` - Achievements API
- ✅ `/referrals/*` - Referral API
- ✅ `/villains/*` - Villain API
- ✅ `/analytics/*` - Analytics API
- ✅ `/subscriptions/*` - Monetization API

---

## Backend Architecture

### Tech Stack (Рекомендуемый)
- **Language**: Go 1.21+
- **Framework**: Chi Router / Fiber
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Storage**: S3-compatible (MinIO / AWS S3)
- **Queue**: RabbitMQ / Redis Streams
- **Monitoring**: Prometheus + Grafana

### Services

```
Backend
├── API Gateway (port 8080)
├── Auth Service (port 8081)
├── Profile Service (port 8082)
├── Task Service (port 8083)
├── AI Service (port 8084)
├── Analytics Service (port 8085)
└── Payment Service (port 8086)
```

---

## Migration Steps

### Step 1: Environment Setup

Создать `.env.production` файл:

```env
VITE_API_BASE_URL=https://api.homework-app.ru
VITE_API_TIMEOUT=30000
VITE_APP_VERSION=1.0.0
VITE_ENABLE_ANALYTICS=true
VITE_ENABLE_LOGGING=true
```

### Step 2: API Client Update

Обновить `src/api/client.ts`:

```typescript
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const API_TIMEOUT = import.meta.env.VITE_API_TIMEOUT || 30000;

const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: API_TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor для добавления JWT токена
axiosInstance.interceptors.request.use(
  (config) => {
    const token = VKStorage.getItem('jwt_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor для обработки ошибок
axiosInstance.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      // Refresh token или redirect на login
    }
    return Promise.reject(error);
  }
);

export const apiClient = {
  async get<T>(url: string): Promise<T> {
    return axiosInstance.get(url);
  },
  async post<T>(url: string, data: any): Promise<T> {
    return axiosInstance.post(url, data);
  },
  async put<T>(url: string, data: any): Promise<T> {
    return axiosInstance.put(url, data);
  },
  async delete<T>(url: string): Promise<T> {
    return axiosInstance.delete(url);
  },
};
```

### Step 3: Authentication Flow

#### JWT Token Flow
```
1. User opens app (VK Bridge)
2. Frontend получает vk_user_id
3. POST /auth/login { vk_user_id, launch_params }
4. Backend проверяет подпись VK
5. Backend создаёт/находит user
6. Backend возвращает { jwt_token, user }
7. Frontend сохраняет JWT в VKStorage
8. Все запросы используют JWT в Authorization header
```

#### Обновить `src/api/auth.ts`:

```typescript
export const authAPI = {
  async login(vkUserId: string, launchParams: string): Promise<AuthResponse> {
    return apiClient.post('/auth/login', {
      vk_user_id: vkUserId,
      launch_params: launchParams,
    });
  },

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    return apiClient.post('/auth/refresh', {
      refresh_token: refreshToken,
    });
  },

  async logout(): Promise<void> {
    return apiClient.post('/auth/logout', {});
  },
};
```

### Step 4: Image Upload

#### Backend Endpoint
```
POST /tasks/upload-image
Content-Type: multipart/form-data

Response:
{
  "image_id": "uuid",
  "url": "https://cdn.homework-app.ru/images/...",
  "thumbnail_url": "https://cdn.homework-app.ru/thumbs/..."
}
```

#### Frontend Implementation

```typescript
export const taskAPI = {
  async uploadImage(file: File): Promise<ImageUploadResponse> {
    const formData = new FormData();
    formData.append('image', file);

    // Compress image before upload
    const compressedFile = await compressImage(file);

    return apiClient.post('/tasks/upload-image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
};
```

### Step 5: Websockets для Real-time Updates

#### Backend
```go
// WebSocket endpoint
GET /ws?token=jwt_token

// Events:
{
  "type": "task_status_update",
  "data": {
    "task_id": "uuid",
    "status": "processing" | "completed" | "error"
  }
}
```

#### Frontend

```typescript
class WebSocketService {
  private ws: WebSocket | null = null;

  connect(token: string) {
    this.ws = new WebSocket(`wss://api.homework-app.ru/ws?token=${token}`);

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };

    this.ws.onerror = (error) => {
      console.error('[WS] Error:', error);
      this.reconnect();
    };
  }

  private handleMessage(message: any) {
    switch (message.type) {
      case 'task_status_update':
        // Update UI based on task status
        break;
      case 'hint_received':
        // Show new hint
        break;
    }
  }
}
```

---

## API Endpoints Documentation

### Authentication

```
POST /auth/login
Body: { vk_user_id, launch_params }
Response: { jwt_token, refresh_token, user }

POST /auth/refresh
Body: { refresh_token }
Response: { jwt_token, refresh_token }

POST /auth/logout
Response: { success: true }
```

### Profile

```
GET /profile/:child_profile_id
Response: ProfileData

GET /profile/:child_profile_id/history
Query: ?mode=help|check&limit=20&offset=0
Response: HistoryAttempt[]

POST /profile/:child_profile_id/avatar
Body: { avatar_id }
Response: { avatar_url }

PUT /profile/:child_profile_id/grade
Body: { grade: number }
Response: { success: true }
```

### Tasks

```
POST /tasks/help
Body: { child_profile_id, image_id, scenario }
Response: { task_id, status }

GET /tasks/:task_id
Response: TaskDetails

POST /tasks/:task_id/answer
Body: { answer: string }
Response: { is_correct, hints }

POST /tasks/:task_id/hint
Response: { hint_text, remaining_hints }

POST /tasks/check
Body: { child_profile_id, image_ids[], scenario }
Response: { task_id, status }
```

### Achievements

```
GET /achievements/:child_profile_id
Response: Achievement[]

GET /achievements/:child_profile_id/stats
Response: AchievementsStats

POST /achievements/:child_profile_id/claim/:achievement_id
Response: { success, rewards }
```

### Referrals

```
GET /referrals/:child_profile_id
Response: ReferralData

POST /referrals/:child_profile_id/invite
Body: { channel: 'vk'|'telegram' }
Response: { success }

POST /referrals/:child_profile_id/claim/:goal_id
Response: { success, rewards }
```

### Villain

```
GET /villains/:child_profile_id/active
Response: Villain | null

GET /villains/:child_profile_id/battle/:villain_id
Response: VillainBattle

GET /villains/:child_profile_id/victory/:villain_id
Response: VillainVictory
```

### Analytics

```
POST /analytics/events
Body: { events: AnalyticsEvent[] }
Response: { success }

POST /analytics/properties
Body: { properties: UserProperties }
Response: { success }
```

### Subscriptions

```
GET /subscriptions/plans
Response: SubscriptionPlan[]

GET /subscriptions/:user_id/current
Response: Subscription | null

POST /subscriptions/:user_id/subscribe
Body: { plan_id, payment_method }
Response: PaymentIntent

POST /subscriptions/:user_id/cancel
Body: { subscription_id }
Response: { success }
```

---

## Error Handling

### Error Response Format

```typescript
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field": "email",
      "reason": "invalid_format"
    }
  }
}
```

### Error Codes

| Code | Status | Description |
|------|--------|-------------|
| `VALIDATION_ERROR` | 400 | Неверные входные данные |
| `UNAUTHORIZED` | 401 | Не авторизован |
| `FORBIDDEN` | 403 | Доступ запрещён |
| `NOT_FOUND` | 404 | Ресурс не найден |
| `CONFLICT` | 409 | Конфликт данных |
| `RATE_LIMIT` | 429 | Превышен лимит запросов |
| `INTERNAL_ERROR` | 500 | Внутренняя ошибка сервера |
| `SERVICE_UNAVAILABLE` | 503 | Сервис недоступен |

### Frontend Error Handling

```typescript
try {
  const data = await profileAPI.getProfile(childProfileId);
} catch (error) {
  if (axios.isAxiosError(error)) {
    const errorCode = error.response?.data?.error?.code;

    switch (errorCode) {
      case 'UNAUTHORIZED':
        // Redirect to login
        break;
      case 'RATE_LIMIT':
        // Show rate limit message
        break;
      default:
        // Show generic error
        showToast('Произошла ошибка. Попробуйте позже.');
    }
  }
}
```

---

## Testing Strategy

### Local Development

```bash
# Run backend locally
docker-compose up -d postgres redis rabbitmq
go run cmd/api/main.go

# Run frontend with local API
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```

### Staging Environment

```
Frontend: https://staging.homework-app.ru
Backend: https://api-staging.homework-app.ru
```

### Production Environment

```
Frontend: https://vk.com/app123456
Backend: https://api.homework-app.ru
CDN: https://cdn.homework-app.ru
```

---

## Monitoring & Logging

### Frontend Logging

```typescript
// Log errors to backend
window.addEventListener('error', (event) => {
  apiClient.post('/logs/errors', {
    message: event.message,
    stack: event.error?.stack,
    url: window.location.href,
    user_agent: navigator.userAgent,
  });
});
```

### Backend Metrics

```
- Request rate (req/s)
- Response time (p50, p95, p99)
- Error rate (%)
- Database query time
- Cache hit rate
- Queue depth
- Active WebSocket connections
```

---

## Security Checklist

- [ ] HTTPS для всех запросов
- [ ] JWT токены с коротким TTL (15 min)
- [ ] Refresh tokens с длинным TTL (30 days)
- [ ] Rate limiting (100 req/min per user)
- [ ] Input validation на backend
- [ ] SQL injection protection (prepared statements)
- [ ] XSS protection (sanitize user input)
- [ ] CSRF tokens для sensitive операций
- [ ] CORS настроен правильно
- [ ] Secrets в environment variables
- [ ] Логи не содержат sensitive data
- [ ] Шифрование персональных данных в БД

---

## Deployment Pipeline

### CI/CD Flow

```
1. Push to main branch
2. Run tests (unit + integration)
3. Build frontend (npm run build)
4. Build backend (go build)
5. Run security scan
6. Deploy to staging
7. Run E2E tests
8. Manual approval
9. Deploy to production
10. Health check
11. Rollback if failed
```

### Deployment Commands

```bash
# Build
npm run build
npm run test

# Deploy frontend (Vercel/Netlify)
vercel deploy --prod

# Deploy backend (Docker)
docker build -t homework-api .
docker push homework-api:latest
kubectl apply -f k8s/deployment.yaml
```

---

## Rollback Strategy

### Frontend Rollback
```bash
# Vercel
vercel rollback <deployment-url>

# Netlify
netlify deploy --prod --dir=dist-backup
```

### Backend Rollback
```bash
# Kubernetes
kubectl rollout undo deployment/homework-api

# Docker
docker tag homework-api:previous homework-api:latest
docker push homework-api:latest
```

---

## Support Contacts

- **Backend Team Lead**: backend@homework-app.ru
- **DevOps**: devops@homework-app.ru
- **API Documentation**: https://docs.homework-app.ru
- **Status Page**: https://status.homework-app.ru
