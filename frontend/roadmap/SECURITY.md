# Security Guidelines — Безопасность приложения

**Проект:** Объяснятель ДЗ MiniApp
**Дата:** 2026-03-29

---

## ⚠️ Критически важно!

Приложение работает с персональными данными детей и родителей. **Безопасность — приоритет №1.**

---

## 🔐 Аутентификация и авторизация

### 1. Получение токена

При запуске приложения через VK/Max/Telegram, платформа передает параметры:

```typescript
interface PlatformLaunchParams {
  vk_user_id: string          // ID пользователя в VK/Max
  vk_access_token_settings: string
  vk_app_id: string
  vk_are_notifications_enabled: string
  vk_is_app_user: string
  vk_is_favorite: string
  vk_language: string
  vk_platform: string
  vk_ref: string
  sign: string                // Подпись для верификации
}
```

**⚠️ КРИТИЧНО:**
1. **НИКОГДА не доверяй `vk_user_id` напрямую** - его можно подделать
2. **ВСЕГДА валидируй `sign`** на backend'е перед выдачей токена
3. **Используй backend как единственный источник истины** для идентификации

### Правильный flow аутентификации

```typescript
// ❌ НЕПРАВИЛЬНО - уязвимость!
async function authenticate() {
  const userId = bridge.getUserId() // Можно подделать!
  const profile = await api.getProfile(userId) // ❌ Любой может получить любой профиль
  return profile
}

// ✅ ПРАВИЛЬНО
async function authenticate() {
  // 1. Получаем launch params от платформы
  const launchParams = await bridge.getLaunchParams()

  // 2. Отправляем все параметры на backend
  const authResponse = await api.post('/api/v1/auth/platform', {
    platform_type: 'vk',
    launch_params: launchParams // включая sign
  })

  // 3. Backend валидирует sign, проверяет подлинность
  // 4. Backend возвращает JWT токен с зашифрованным user_id
  const { access_token, user } = authResponse

  // 5. Сохраняем токен для всех последующих запросов
  storage.setToken(access_token)

  return user
}
```

### Валидация на Backend

```go
// Backend должен валидировать sign
func ValidateVKSign(params map[string]string, secretKey string) bool {
    // 1. Извлечь sign
    sign := params["sign"]
    delete(params, "sign")

    // 2. Отсортировать параметры
    keys := make([]string, 0, len(params))
    for k := range params {
        if strings.HasPrefix(k, "vk_") {
            keys = append(keys, k)
        }
    }
    sort.Strings(keys)

    // 3. Построить query string
    var queryString string
    for _, k := range keys {
        queryString += k + "=" + params[k] + "&"
    }
    queryString = strings.TrimSuffix(queryString, "&")

    // 4. Вычислить HMAC
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(queryString))
    expectedSign := base64.URLEncoding.EncodeToString(h.Sum(nil))

    // 5. Сравнить
    return sign == expectedSign
}
```

---

## 🛡️ API Security

### 1. Авторизация каждого запроса

**КАЖДЫЙ** API запрос должен содержать JWT токен:

```typescript
// api/client.ts
class APIClient {
  private getHeaders(): Record<string, string> {
    const token = storage.getToken()

    if (!token) {
      throw new Error('Not authenticated')
    }

    return {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      'X-Platform-Type': getPlatformType(), // vk/max/telegram
      'X-App-Version': APP_VERSION
    }
  }

  async get(url: string) {
    const response = await fetch(API_BASE_URL + url, {
      headers: this.getHeaders()
    })

    if (response.status === 401) {
      // Token expired, re-authenticate
      await this.reAuthenticate()
      throw new Error('Authentication required')
    }

    return response.json()
  }
}
```

### 2. Запрет прямого доступа к чужим данным

**❌ НИКОГДА не позволяй фронтенду указывать ID другого пользователя:**

```typescript
// ❌ УЯЗВИМОСТЬ!
async function getProfile(childProfileId: string) {
  // Любой может подставить чужой ID!
  return api.get(`/api/v1/child-profile/${childProfileId}`)
}

// ❌ УЯЗВИМОСТЬ!
async function getAttempts(childProfileId: string) {
  // Можно получить историю другого ребенка!
  return api.get(`/api/v1/attempts/history?child_profile_id=${childProfileId}`)
}
```

**✅ ПРАВИЛЬНО - Backend сам определяет user_id из токена:**

```typescript
// ✅ БЕЗОПАСНО
async function getProfile() {
  // Backend извлекает user_id из JWT токена
  // и возвращает только данные текущего пользователя
  return api.get('/api/v1/profile/me')
}

// ✅ БЕЗОПАСНО
async function getAttempts() {
  // Backend сам определяет child_profile_id из токена
  return api.get('/api/v1/attempts/history')
}
```

### Backend реализация

```go
// middleware для извлечения user_id из JWT
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Извлечь токен из заголовка
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Unauthorized", 401)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")

        // 2. Валидировать JWT
        claims, err := jwt.Verify(token, jwtSecret)
        if err != nil {
            http.Error(w, "Invalid token", 401)
            return
        }

        // 3. Добавить user_id в context
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "child_profile_id", claims.ChildProfileID)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// В handlers используем только данные из context
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
    // ✅ Извлекаем только из токена
    childProfileID := r.Context().Value("child_profile_id").(string)

    // ✅ Пользователь может получить только свой профиль
    profile, err := db.GetChildProfile(childProfileID)
    if err != nil {
        http.Error(w, "Not found", 404)
        return
    }

    json.NewEncoder(w).WriteJson(profile)
}

// ❌ НИКОГДА не принимай ID из query параметров!
func GetProfileHandlerVULNERABLE(w http.ResponseWriter, r *http.Request) {
    // ❌ УЯЗВИМОСТЬ!
    childProfileID := r.URL.Query().Get("child_profile_id")

    // Любой может получить любой профиль!
    profile, _ := db.GetChildProfile(childProfileID)
    json.NewEncoder(w).WriteJson(profile)
}
```

---

## 🔒 Доступ к данным профиля

### Child Profile

```typescript
// ❌ НЕПРАВИЛЬНО
interface GetChildProfileRequest {
  child_profile_id: string // ❌ Позволяет получить любой профиль
}

// ✅ ПРАВИЛЬНО
// Backend определяет child_profile_id из JWT токена
GET /api/v1/profile/me
// Возвращает только текущий профиль
```

### Попытки (Attempts)

```typescript
// ❌ НЕПРАВИЛЬНО
GET /api/v1/attempts/history?child_profile_id=xxx // ❌ Можно указать чужой ID

// ✅ ПРАВИЛЬНО
GET /api/v1/attempts/history // Backend сам определяет владельца
```

### Достижения

```typescript
// ❌ НЕПРАВИЛЬНО
GET /api/v1/achievements?child_profile_id=xxx // ❌ Чужие достижения

// ✅ ПРАВИЛЬНО
GET /api/v1/achievements // Backend вернет только для текущего пользователя
```

---

## 👨‍👩‍👧 Parent Gate - Доступ родителей

Родитель может иметь несколько детских профилей. Нужна дополнительная проверка:

```typescript
// 1. Родитель выбирает профиль ребенка
async function switchChildProfile(childProfileId: string) {
  // Backend проверяет, что этот ребенок принадлежит родителю
  const response = await api.post('/api/v1/profile/switch', {
    child_profile_id: childProfileId
  })

  // Backend возвращает новый JWT токен с child_profile_id
  const { access_token } = response
  storage.setToken(access_token)
}
```

Backend проверка:

```go
func SwitchChildProfileHandler(w http.ResponseWriter, r *http.Request) {
    parentUserID := r.Context().Value("parent_user_id").(string)

    var req struct {
        ChildProfileID string `json:"child_profile_id"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    // ✅ КРИТИЧНО: Проверить, что профиль принадлежит родителю
    childProfile, err := db.GetChildProfile(req.ChildProfileID)
    if err != nil {
        http.Error(w, "Not found", 404)
        return
    }

    if childProfile.ParentUserID != parentUserID {
        // ❌ Попытка доступа к чужому профилю!
        http.Error(w, "Forbidden", 403)
        return
    }

    // ✅ Всё ок, выдаем новый токен
    token := jwt.Create(map[string]interface{}{
        "parent_user_id": parentUserID,
        "child_profile_id": req.ChildProfileID,
    })

    json.NewEncoder(w).WriteJson(map[string]string{
        "access_token": token,
    })
}
```

---

## 🖼️ Безопасность изображений

### Upload

```typescript
// При загрузке изображения
async function uploadImage(attemptId: string, image: File) {
  // Backend должен проверить, что attempt принадлежит пользователю
  const formData = new FormData()
  formData.append('image', image)

  const response = await api.post(`/api/v1/attempts/${attemptId}/images`, formData)
  return response
}
```

Backend проверка:

```go
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
    childProfileID := r.Context().Value("child_profile_id").(string)
    attemptID := mux.Vars(r)["attempt_id"]

    // ✅ Проверить, что attempt принадлежит пользователю
    attempt, err := db.GetAttempt(attemptID)
    if err != nil || attempt.ChildProfileID != childProfileID {
        http.Error(w, "Forbidden", 403)
        return
    }

    // Продолжаем загрузку
    // ...
}
```

### Доступ к изображениям

```typescript
// Backend должен генерировать signed URLs
interface AttemptImage {
  attempt_image_id: string
  image_url: string // https://storage.com/signed-url?token=xxx&expires=xxx
  thumbnail_url: string
}
```

Signed URLs истекают через определенное время и содержат подпись, которую нельзя подделать.

---

## 📧 Email и персональные данные

### Parent User Email

```typescript
// ❌ НЕПРАВИЛЬНО - email в открытом виде
interface ParentUser {
  email: "user@example.com" // ❌ Полный email
}

// ✅ ПРАВИЛЬНО - маскированный email
interface ParentUser {
  email_masked: "u***@example.com" // ✅ Частично скрыт
  email_verified: boolean
  email_domain: "example.com" // Для аналитики
}
```

### Для изменения email

```typescript
// Родитель вводит новый email
async function updateEmail(newEmail: string) {
  // Backend отправляет код верификации
  await api.post('/api/v1/parent/email/update', {
    new_email: newEmail
  })

  // Просим ввести код
  const code = await promptVerificationCode()

  // Подтверждаем
  await api.post('/api/v1/parent/email/confirm', {
    verification_code: code
  })
}
```

Backend никогда не возвращает полный email в ответах API (только маскированный).

---

## 🔐 Sensitive Data Protection

### Что НЕ должно попадать в frontend

❌ **НИКОГДА не передавай на frontend:**
- Полные email адреса (только маскированные)
- Пароли или их хеши
- Токены доступа к внешним сервисам
- API keys
- Внутренние ID других пользователей
- PII (Personally Identifiable Information) других пользователей

### Логирование

```typescript
// ❌ ОПАСНО
console.log('User data:', {
  email: user.email, // ❌ email в логах
  token: storage.getToken() // ❌ токен в логах
})

// ✅ БЕЗОПАСНО
console.log('User data:', {
  user_id: user.child_profile_id.substring(0, 8) + '***', // Частично скрыт
  email_domain: user.email_domain // Только домен
})

// Для отладки используй специальный dev режим
if (process.env.NODE_ENV === 'development') {
  // Только в dev mode
  console.debug('Full user data:', user)
}
```

---

## 🛡️ XSS Protection

### Рендеринг пользовательского контента

```typescript
// Отображение имени профиля
function ProfileCard({ profile }: { profile: ChildProfile }) {
  // ✅ React автоматически экранирует
  return (
    <div>
      <h2>{profile.display_name}</h2> {/* ✅ Безопасно */}
    </div>
  )
}

// ❌ ОПАСНО - dangerouslySetInnerHTML
function ProfileCardUNSAFE({ profile }: { profile: ChildProfile }) {
  return (
    <div dangerouslySetInnerHTML={{ __html: profile.display_name }} />
    {/* ❌ XSS уязвимость! */}
  )
}

// ✅ Если нужен HTML - используй санитайзер
import DOMPurify from 'dompurify'

function SafeHTMLRender({ html }: { html: string }) {
  const cleanHTML = DOMPurify.sanitize(html)
  return <div dangerouslySetInnerHTML={{ __html: cleanHTML }} />
}
```

---

## 🔒 HTTPS Only

**КРИТИЧНО:** Приложение ДОЛЖНО работать ТОЛЬКО через HTTPS.

```typescript
// В index.html
<meta http-equiv="Content-Security-Policy"
      content="upgrade-insecure-requests">

// Проверка в runtime
if (window.location.protocol !== 'https:' &&
    window.location.hostname !== 'localhost') {
  window.location.href = 'https:' + window.location.href.substring(
    window.location.protocol.length
  )
}
```

---

## 🚫 CORS и CSP

### Content Security Policy

```html
<meta http-equiv="Content-Security-Policy"
      content="
        default-src 'self';
        script-src 'self' 'unsafe-inline' vk.com *.vk.com;
        style-src 'self' 'unsafe-inline';
        img-src 'self' data: https: blob:;
        connect-src 'self' https://api.homework.app wss://api.homework.app;
        font-src 'self';
        frame-src 'none';
      ">
```

### CORS на Backend

Backend должен разрешать запросы только с известных доменов:

```go
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")

            // Проверить, что origin в whitelist
            allowed := false
            for _, allowedOrigin := range allowedOrigins {
                if origin == allowedOrigin {
                    allowed = true
                    break
                }
            }

            if allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 📊 Аналитика - Privacy

### Что можно отправлять

✅ **Безопасно отправлять:**
- `child_profile_id` (это ID, а не имя)
- `parent_user_id` (это ID, а не email)
- `email_domain` (только домен, не полный email)
- Агрегированные метрики (счетчики, прогресс)

❌ **НЕ отправлять:**
- `display_name` (имя ребенка)
- `email` (полный email родителя)
- Содержимое заданий (текст с фото)
- Ответы пользователя

### Пример

```typescript
// ✅ ПРАВИЛЬНО
analytics.sendEvent({
  event_name: 'grade_selected',
  child_profile_id: profile.child_profile_id, // ✅ UUID, не PII
  grade: 5, // ✅ Класс - не персональные данные
  platform_type: 'vk'
})

// ❌ НЕПРАВИЛЬНО
analytics.sendEvent({
  event_name: 'profile_created',
  child_name: profile.display_name, // ❌ Имя ребенка!
  parent_email: user.email // ❌ Email родителя!
})
```

---

## 🔐 Storage Security

### LocalStorage / SessionStorage

```typescript
// ✅ Храни только JWT токен
storage.setItem('access_token', token)

// ❌ НЕ храни чувствительные данные
storage.setItem('user_password', password) // ❌ НИКОГДА!
storage.setItem('user_email', email) // ❌ Не нужно
storage.setItem('credit_card', card) // ❌ НИКОГДА!
```

### Очистка при logout

```typescript
async function logout() {
  // Удалить токен
  storage.removeItem('access_token')

  // Очистить все данные
  storage.clear()

  // Очистить state
  profileStore.reset()

  // Redirect на онбординг
  router.push('/onboarding')
}
```

---

## 🚨 Security Checklist

Перед релизом проверь:

### Authentication & Authorization
- [ ] Все API запросы требуют JWT токен
- [ ] Backend валидирует sign от VK/Max/Telegram
- [ ] Токены истекают через разумное время (30 дней)
- [ ] Refresh token механизм реализован
- [ ] Logout полностью очищает токены

### Data Access
- [ ] Пользователь не может получить данные другого пользователя
- [ ] Backend извлекает user_id ТОЛЬКО из JWT токена
- [ ] Parent может переключаться только между своими детьми
- [ ] Изображения доступны по signed URLs
- [ ] Signed URLs истекают

### Input Validation
- [ ] Все пользовательские данные валидируются на backend
- [ ] XSS protection (React автоматически экранирует)
- [ ] SQL injection protection (используй prepared statements)
- [ ] Limit на размер файлов (изображения)
- [ ] Limit на rate (rate limiting на API)

### Sensitive Data
- [ ] Email маскированы на frontend
- [ ] Пароли НЕ хранятся на frontend
- [ ] API keys НЕ в коде (только на backend)
- [ ] Логи НЕ содержат чувствительные данные
- [ ] Analytics НЕ содержит PII

### Network Security
- [ ] HTTPS only (HTTP redirect на HTTPS)
- [ ] CORS настроен правильно (whitelist origins)
- [ ] CSP заголовки настроены
- [ ] Cookies с httpOnly, secure, sameSite flags

### Platform Security
- [ ] VK Bridge sign validation
- [ ] Telegram WebApp initData validation
- [ ] Max SDK signature validation

---

## 📚 Дополнительные ресурсы

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [JWT Best Practices](https://datatracker.ietf.org/doc/html/rfc8725)
- [VK Security Best Practices](https://dev.vk.com/mini-apps/development/security)
- [React Security](https://react.dev/learn/keeping-components-pure)

---

## 🔥 В случае инцидента

1. **Немедленно** отзови все токены
2. Заблокируй скомпрометированные аккаунты
3. Уведоми пользователей
4. Проведи аудит безопасности
5. Исправь уязвимость
6. Обнови документацию

**Контакты для сообщений о уязвимостях:**
security@homework.app
