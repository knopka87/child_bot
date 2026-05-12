# Security Headers и HTTPS

## Обзор

Приложение использует несколько слоев безопасности на уровне HTTP для защиты от распространенных атак.

## Security Middleware

### 1. HTTPS Redirect

**Файл:** `api/internal/api/middleware/security.go`

#### Что делает
- Автоматически редиректит все HTTP запросы на HTTPS в production режиме
- Использует 301 Permanent Redirect
- Проверяет заголовок `X-Forwarded-Proto` (для работы за reverse proxy)

#### Режимы работы

**Development (`ENV=development`):**
- Редирект отключен
- Можно работать по HTTP на localhost

**Production (`ENV=production`):**
- Все HTTP запросы редиректятся на HTTPS
- Обязательное использование HTTPS

#### Примеры

```bash
# Development - проходит без редиректа
curl http://localhost:8080/health
# → 200 OK

# Production - редиректится на HTTPS
curl http://api.obiasnyatel-dz.ru/health
# → 301 Moved Permanently
# → Location: https://api.obiasnyatel-dz.ru/health
```

### 2. Security Headers

Автоматически добавляются ко всем ответам API.

#### HSTS (HTTP Strict Transport Security)

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

**Защита от:**
- Man-in-the-Middle атак
- SSL Stripping атак
- Случайного использования HTTP

**Как работает:**
- Браузер запоминает что сайт должен работать только по HTTPS
- На 1 год (31536000 секунд)
- Применяется ко всем поддоменам

#### X-Frame-Options

```
X-Frame-Options: DENY
```

**Защита от:**
- Clickjacking атак
- UI Redressing

**Как работает:**
- Полностью запрещает загрузку страницы в iframe/frame
- Никто не может встроить ваше приложение в свой сайт

**Альтернативы:**
- `SAMEORIGIN` - разрешает только для того же домена
- `DENY` - запрещает полностью (рекомендуется)

#### X-Content-Type-Options

```
X-Content-Type-Options: nosniff
```

**Защита от:**
- MIME type sniffing атак
- Выполнение вредоносных скриптов

**Как работает:**
- Браузер доверяет только Content-Type заголовку
- Не пытается угадать тип файла по содержимому

#### X-XSS-Protection

```
X-XSS-Protection: 1; mode=block
```

**Защита от:**
- Reflected XSS атак (legacy браузеры)

**Как работает:**
- Включает встроенный XSS фильтр браузера
- При обнаружении XSS блокирует загрузку страницы

**Примечание:** Современные браузеры используют CSP вместо этого, но header полезен для старых браузеров.

#### Content-Security-Policy (CSP)

```
Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://vk.com https://*.vk.com https://*.vk.me; style-src 'self' 'unsafe-inline'; img-src 'self' data: https: blob:; font-src 'self' data:; connect-src 'self' https://vk.com https://*.vk.com https://*.vk.me wss://im-0*.vk.com; frame-src https://vk.com https://*.vk.com; frame-ancestors 'none'; form-action 'self'; base-uri 'self'; object-src 'none'
```

**Защита от:**
- XSS (Cross-Site Scripting)
- Data injection атак
- Clickjacking
- Загрузки вредоносного контента

**Директивы:**

| Директива | Значение | Что разрешает |
|-----------|----------|---------------|
| `default-src` | `'self'` | По умолчанию только с текущего домена |
| `script-src` | `'self' 'unsafe-inline' 'unsafe-eval' https://vk.com https://*.vk.com` | Скрипты с домена и VK (для VK Bridge) |
| `style-src` | `'self' 'unsafe-inline'` | Стили с домена и inline стили |
| `img-src` | `'self' data: https: blob:` | Картинки с любых HTTPS, data URLs, blob |
| `font-src` | `'self' data:` | Шрифты с домена и data URLs |
| `connect-src` | `'self' https://vk.com https://*.vk.com wss://im-0*.vk.com` | AJAX к домену и VK API |
| `frame-src` | `https://vk.com https://*.vk.com` | Iframe только от VK |
| `frame-ancestors` | `'none'` | Запрет на встраивание в iframe |
| `form-action` | `'self'` | Отправка форм только на свой домен |
| `base-uri` | `'self'` | `<base>` только для своего домена |
| `object-src` | `'none'` | Запрет `<object>`, `<embed>`, `<applet>` |

**Почему `unsafe-inline` и `unsafe-eval`?**
- Необходимо для работы VK Bridge и некоторых frontend библиотек
- В идеале нужно убрать, но это требует рефакторинга frontend
- TODO: Использовать nonce-based CSP в будущем

#### Referrer-Policy

```
Referrer-Policy: strict-origin-when-cross-origin
```

**Защита от:**
- Утечки чувствительной информации в referrer
- Privacy leaks

**Как работает:**
- Same-origin запросы: отправляется полный referrer
- Cross-origin запросы: отправляется только origin (без пути)
- HTTP → HTTPS: referrer не отправляется

#### Permissions-Policy

```
Permissions-Policy: geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=(), accelerometer=()
```

**Защита от:**
- Несанкционированного доступа к device APIs
- Privacy leaks

**Как работает:**
- Отключает доступ к геолокации, микрофону, камере и другим APIs
- `()` означает "никому не разрешено"

## Secure Cookies

**Файл:** `api/internal/api/middleware/security.go`

### SecureCookie Helper

Используйте эту функцию при установке cookies в handlers:

```go
cookie := middleware.SecureCookie("session_id", "value123", 3600)
http.SetCookie(w, cookie)
```

### Параметры

| Параметр | Development | Production | Назначение |
|----------|------------|------------|------------|
| `HttpOnly` | ✅ true | ✅ true | Защита от XSS - JavaScript не может прочитать cookie |
| `Secure` | ❌ false | ✅ true | Передача только по HTTPS |
| `SameSite` | Strict | Strict | CSRF защита - cookie отправляется только с same-site запросами |
| `Path` | `/` | `/` | Cookie доступна для всего сайта |
| `MaxAge` | custom | custom | Время жизни cookie в секундах |

## Проверка Security Headers

### Online Tools

1. **Security Headers:**
   - https://securityheaders.com/
   - Анализирует все security headers
   - Выставляет оценку A+ / A / B / C / D / F

2. **Mozilla Observatory:**
   - https://observatory.mozilla.org/
   - Комплексный security аудит
   - Рекомендации по улучшению

### Командная строка

```bash
# Проверка заголовков
curl -I https://api.obiasnyatel-dz.ru/health

# Проверка CSP
curl -I https://api.obiasnyatel-dz.ru/health | grep Content-Security-Policy

# Проверка HSTS
curl -I https://api.obiasnyatel-dz.ru/health | grep Strict-Transport-Security
```

### Browser DevTools

1. Откройте DevTools (F12)
2. Перейдите в **Network** tab
3. Выберите любой запрос
4. Откройте **Headers** tab
5. Проверьте **Response Headers**

Должны быть все headers из списка выше.

## Testing

### Unit Tests

```bash
cd api
go test -v ./internal/api/middleware/security_test.go ./internal/api/middleware/security.go
```

Тесты проверяют:
- ✅ Все security headers устанавливаются
- ✅ HTTPS redirect работает в production
- ✅ HTTPS redirect отключен в development
- ✅ SecureCookie создает правильные параметры

### Integration Testing

```bash
# 1. Запустите backend в production режиме
ENV=production go run cmd/api/main.go

# 2. Проверьте security headers
curl -I http://localhost:8080/health

# Должен быть редирект на HTTPS (если настроен reverse proxy)
```

## Deployment

### Nginx Reverse Proxy

Если используете Nginx перед API:

```nginx
server {
    listen 80;
    server_name api.obiasnyatel-dz.ru;

    # Редирект HTTP -> HTTPS (дополнительный слой)
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.obiasnyatel-dz.ru;

    # SSL сертификаты
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # Передаем протокол в API
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    location / {
        proxy_pass http://localhost:8080;
    }
}
```

### Docker

В docker-compose.production.yml:

```yaml
backend:
  environment:
    - ENV=production  # Включает HTTPS redirect
```

### Environment Variables

```bash
# Production
ENV=production

# Development
ENV=development
```

## Security Checklist

Перед deploy в production:

- [ ] `ENV=production` установлен
- [ ] SSL/TLS сертификаты настроены
- [ ] Nginx/reverse proxy настроен
- [ ] X-Forwarded-Proto передается в API
- [ ] Все запросы идут через HTTPS
- [ ] Security headers проверены через securityheaders.com
- [ ] HSTS preload list submission (опционально)

## HSTS Preload (Optional)

Для максимальной безопасности можно добавить домен в HSTS preload list:

1. Обновите HSTS header:
   ```
   Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
   ```

2. Подайте заявку:
   - https://hstspreload.org/

**Внимание:** Это необратимое действие! После добавления в preload list домен НИКОГДА не сможет работать по HTTP.

## Известные Issues

### CSP `unsafe-inline` и `unsafe-eval`

**Проблема:** Разрешены inline скрипты и eval

**Риск:** Средний (потенциальный XSS)

**План исправления:**
1. Провести аудит frontend кода
2. Заменить inline скрипты на external
3. Убрать eval() если используется
4. Использовать CSP nonces для разрешенных inline скриптов

**Трекинг:** TODO - создать issue в GitHub

## Мониторинг

### CSP Violation Reports

Можно настроить reporting endpoint для CSP violations:

```go
// В security.go добавить:
csp := "... report-uri /api/v1/csp-report"
```

Создать handler для приема репортов:

```go
func CSPReportHandler(w http.ResponseWriter, r *http.Request) {
    // Логировать нарушения CSP
    // Отправлять в Sentry/мониторинг
}
```

### Логи

Middleware логирует:

```
[Security] Redirecting HTTP -> HTTPS: /path
```

Проверяйте логи на частые редиректы - может указывать на misconfiguration.

## Best Practices

1. **Всегда используйте HTTPS в production**
   - Никогда не передавайте чувствительные данные по HTTP

2. **Регулярно проверяйте security headers**
   - Используйте securityheaders.com
   - Автоматизируйте проверки в CI/CD

3. **Обновляйте CSP по мере развития**
   - Убирайте unsafe-inline когда возможно
   - Используйте nonces для inline скриптов

4. **Тестируйте после изменений**
   - Проверяйте что приложение работает
   - CSP может сломать функциональность

5. **Документируйте исключения**
   - Объясняйте почему используется unsafe-inline
   - Планируйте исправление

## Ссылки

- [OWASP Secure Headers Project](https://owasp.org/www-project-secure-headers/)
- [MDN: Content Security Policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)
- [Security Headers Check](https://securityheaders.com/)
- [Mozilla Observatory](https://observatory.mozilla.org/)
- [HSTS Preload List](https://hstspreload.org/)
