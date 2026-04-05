# Настройка VK Mini App для тестирования

Это руководство основано на официальной документации VK для разработчиков.

## Введение

VK Mini Apps — это веб-приложения, которые встраиваются в экосистему ВКонтакте через:
- **Мобильное приложение** (WebView)
- **Десктопная версия** (iFrame)
- **Мобильные браузеры** (MVK)

Приложение размещается на собственном сервере и открывается пользователям через платформу VK.

## 1. Создание приложения в VK

### Шаг 1: Создайте новое приложение
1. Перейдите в [панель администрирования VK Developer](https://dev.vk.com/ru/admin/create-app)
2. Выберите тип **"mini-app"** (Мини-приложение)
3. Укажите название приложения
4. Выберите подходящую категорию
5. Нажмите **"Создать приложение"**
6. Подтвердите создание выбранным способом

**Ваше приложение создано с ID**: `54517931`

### Шаг 2: Получите параметры приложения

После создания вам понадобятся:
- **App ID** (ID приложения): `54517931`
- **Защищенный ключ** (Secret key) — для проверки подписи параметров запуска
- **Сервисный токен** (Service token) — для работы с VK API на стороне сервера

Эти параметры доступны в разделе **"Настройки"** → **"Ключи доступа"**.

## 2. Настройка приложения

### Откройте настройки приложения
Перейдите по ссылке: https://vk.com/editapp?id=54517931

### Раздел "Размещение"

VK Mini Apps требует указать URL для трёх вариантов размещения:

#### 1. Мобильное приложение (VK App)
- **Состояние**: Включено (в режиме разработки — только для администраторов)
- **URL**: `https://wet-poets-open.loca.lt`
- Приложение открывается в WebView внутри мобильного приложения VK

#### 2. Мобильный браузер (MVK)
- **Состояние**: Включено
- **URL**: `https://wet-poets-open.loca.lt`
- Приложение открывается в мобильном браузере m.vk.com

#### 3. Десктопная версия (Web)
- **Состояние**: Включено
- **URL**: `https://wet-poets-open.loca.lt`
- Приложение открывается в iFrame на vk.com

⚠️ **Важно**: Для всех трёх вариантов укажите один и тот же URL — ваше приложение адаптируется автоматически.

### Раздел "Настройки"

#### Основные параметры:
- **Тип приложения**: Mini Apps (VK Apps)
- **Доверенный redirect URI**: `https://wet-poets-open.loca.lt`
- **Категория**: Выберите подходящую (Образование, Игры и т.д.)

#### Безопасность:
- ✅ Используйте HTTPS (обязательно для продакшена)
- ✅ Включите проверку подписи параметров запуска
- ✅ Храните Secret Key в безопасности (только на сервере)

### Раздел "Доступность"

Для тестирования:
- **Статус приложения**: Тестовое приложение
- **Видимость**: Приложение видно всем (для открытого тестирования)
- **Режим разработки**: Включён — позволяет тестировать без модерации

Для продакшена:
- **Статус приложения**: Открытое приложение
- Потребуется пройти **модерацию** (до 7 дней)
- Релизы новых приложений — еженедельно по четвергам

### Сохраните изменения
Нажмите **"Сохранить изменения"** внизу страницы.

## 3. Параметры запуска VK Mini Apps

При открытии вашего приложения, VK автоматически добавляет параметры в URL:

```
https://wet-poets-open.loca.lt/?vk_user_id=123&vk_app_id=54517931&vk_platform=mobile_web&sign=abc123...
```

### Основные параметры:

| Параметр | Описание |
|----------|----------|
| `vk_user_id` | ID пользователя, открывшего приложение |
| `vk_app_id` | ID вашего приложения (54517931) |
| `vk_platform` | Платформа запуска: `mobile_android`, `mobile_iphone`, `desktop_web`, `mobile_web` |
| `vk_language` | Язык интерфейса: `ru`, `en` и т.д. |
| `vk_are_notifications_enabled` | Включены ли уведомления: `0` или `1` |
| `vk_is_favorite` | Приложение в избранном: `0` или `1` |
| `vk_group_id` | ID сообщества (если запущено из сообщества) |
| `vk_viewer_group_role` | Роль в сообществе: `admin`, `editor`, `member`, `none` |
| `vk_ref` | Источник запуска (для аналитики) |
| `sign` | **Подпись параметров** (SHA-256 HMAC с Secret Key) |

### Проверка подписи (обязательно!)

Параметр `sign` защищает от подделки данных. Проверяйте подпись на **сервере**:

```javascript
// Пример проверки подписи на Node.js
const crypto = require('crypto');

function checkSign(params, secretKey) {
  const ordered = {};
  Object.keys(params).filter(key => key.startsWith('vk_'))
    .sort()
    .forEach(key => ordered[key] = params[key]);

  const stringParams = Object.entries(ordered)
    .map(([key, value]) => `${key}=${value}`)
    .join('&');

  const hash = crypto
    .createHmac('sha256', secretKey)
    .update(stringParams)
    .digest('base64')
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '');

  return hash === params.sign;
}
```

⚠️ **Никогда** не полагайтесь на параметры без проверки подписи!

## 4. Интеграция VK Bridge

VK Bridge (ранее VK Connect) — библиотека для взаимодействия с платформой VK.

### Установка:

```bash
npm install @vkontakte/vk-bridge
```

### Базовое использование:

```javascript
import bridge from '@vkontakte/vk-bridge';

// Инициализация
bridge.send('VKWebAppInit');

// Получение информации о пользователе
bridge.send('VKWebAppGetUserInfo')
  .then(data => {
    console.log('User:', data);
  })
  .catch(error => {
    console.error('Error:', error);
  });

// Проверка поддержки метода
if (bridge.supports('VKWebAppGetUserInfo')) {
  // Метод поддерживается на текущей платформе
}
```

### Популярные методы VK Bridge:

- `VKWebAppGetUserInfo` — информация о пользователе
- `VKWebAppStorageGet` / `VKWebAppStorageSet` — локальное хранилище
- `VKWebAppShare` — поделиться приложением
- `VKWebAppShowImages` — просмотр изображений
- `VKWebAppOpenPayForm` — платежи
- `VKWebAppCallAPIMethod` — вызов методов VK API

📖 Полный список методов: https://dev.vk.com/bridge/overview

## 5. Использование VKUI

VKUI — набор React-компонентов для создания интерфейсов в стиле VK.

### Установка:

```bash
npm install @vkontakte/vkui @vkontakte/icons
```

### Обязательный метатег:

Добавьте в `<head>` вашего HTML:

```html
<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no, user-scalable=no, viewport-fit=cover">
```

### Базовая структура:

```jsx
import { AdaptivityProvider, ConfigProvider, AppRoot } from '@vkontakte/vkui';
import '@vkontakte/vkui/dist/vkui.css';

function App() {
  return (
    <ConfigProvider>
      <AdaptivityProvider>
        <AppRoot>
          {/* Ваше приложение */}
        </AppRoot>
      </AdaptivityProvider>
    </ConfigProvider>
  );
}
```

📖 Документация VKUI: https://vkcom.github.io/VKUI/

## 6. Настройка локального туннеля (для разработки)

### Текущий статус:
✅ **Локальный туннель запущен**: `https://wet-poets-open.loca.lt`
✅ **Процесс работает**: PID 15726, 15058

### Запуск localtunnel:

```bash
# Установка
npm install -g localtunnel

# Запуск на порт 80
lt --port 80
```

Вы получите URL вида: `https://random-name.loca.lt`

### ⚠️ Важно про localtunnel:

1. **URL меняется** при каждом перезапуске
2. При первом открытии показывает предупреждение — нажмите **"Continue"**
3. Для продакшена используйте постоянный домен с HTTPS

### Обновление URL после перезапуска:

Если получили новый URL от localtunnel:

1. **Обновите настройки VK**: https://vk.com/editapp?id=54517931
   - Укажите новый URL во всех трёх вариантах размещения

2. **Обновите ALLOWED_ORIGINS** в `.env`:
   ```bash
   ALLOWED_ORIGINS=http://localhost,https://NEW-URL.loca.lt
   ```

3. **Перезапустите backend**:
   ```bash
   docker compose restart backend
   ```

4. **Обновите nginx** (если нужно):
   ```bash
   docker compose up -d --build frontend
   ```

## 7. Открытие приложения для тестирования

### Вариант 1: Прямая ссылка (Web)
Откройте в браузере: **https://vk.com/app54517931**

### Вариант 2: Через каталог сервисов
1. Откройте VK в браузере
2. Перейдите в раздел **"Сервисы"**
3. Найдите своё приложение или откройте: https://vk.com/services?w=app54517931

### Вариант 3: В мобильном приложении VK
1. Откройте приложение VK на телефоне (iOS/Android)
2. Перейдите в раздел **"Сервисы"**
3. Найдите приложение в списке или откройте по deep link: `vk://vk.com/app54517931`

### Вариант 4: В мобильном браузере
Откройте в мобильном браузере: **https://m.vk.com/app54517931**

## 8. Тестирование и отладка

### Проверка параметров запуска

Откройте DevTools → Console и проверьте URL:

```javascript
console.log('Location:', window.location.href);
console.log('Search params:', window.location.search);

// Парсинг параметров
const params = new URLSearchParams(window.location.search);
console.log('vk_user_id:', params.get('vk_user_id'));
console.log('vk_platform:', params.get('vk_platform'));
console.log('sign:', params.get('sign'));
```

### Проверка VK Bridge

```javascript
import bridge from '@vkontakte/vk-bridge';

// Проверка инициализации
bridge.send('VKWebAppInit')
  .then(() => console.log('VK Bridge initialized'))
  .catch(err => console.error('VK Bridge error:', err));

// Получение информации о платформе
bridge.send('VKWebAppGetClientVersion')
  .then(data => console.log('Client version:', data))
  .catch(err => console.error('Error:', err));
```

### Логи в консоли

Должны появиться логи:
```
[Platform] Platform detected: vk
[VKStorage] Using VK Bridge storage
[APIClient] Set X-Platform-ID header: vk
```

### Проверка CORS

Откройте DevTools → Network и убедитесь:
- ✅ Запросы к `/api/*` не падают с ошибками CORS
- ✅ В заголовках есть `Access-Control-Allow-Origin`
- ✅ Preflight запросы (OPTIONS) возвращают 204

### Проверка запросов к API

Проверьте в Network tab:
- ✅ Заголовок `X-Platform-ID: vk` присутствует
- ✅ Заголовок `X-Child-Profile-ID` есть (после онбординга)
- ✅ Нет ошибок 401/403

## 9. Требования и лучшие практики

### Технические требования:

✅ **HTTPS обязателен** (для продакшена)
✅ **Responsive дизайн** — приложение должно работать на всех разрешениях
✅ **Поддержка iOS и Android** — тестируйте на обеих платформах
✅ **Быстрая загрузка** — оптимизируйте размер бандла
✅ **Офлайн-режим** — используйте Service Workers (опционально)

### Безопасность:

🔒 **Всегда проверяйте подпись** параметров запуска на сервере
🔒 **Не храните секреты** в frontend-коде
🔒 **Используйте HTTPS** для всех запросов
🔒 **Валидируйте** все входные данные
🔒 **Не передавайте токены** через Referer или query параметры

### UX/UI рекомендации:

📱 **Используйте VKUI** — приложение будет выглядеть нативно
📱 **Поддержка свайпов** — для навигации назад
📱 **Аппаратные кнопки** — обрабатывайте "Назад" на Android
📱 **Адаптивность** — десктоп, планшет, мобильный
📱 **Тёмная тема** — поддержите через VKUI

## 10. Модерация и публикация

### Требования для публикации в каталог:

1. ✅ Соответствие [Правилам размещения](https://dev.vk.com/ru/rules)
2. ✅ Заполненное описание и скриншоты
3. ✅ Иконки 72x72, 128x128, 256x256 px
4. ✅ Обложка для каталога
5. ✅ Работающий функционал без критических багов

### Процесс модерации:

- **Срок**: До 7 дней
- **Релизы**: Еженедельно по четвергам
- **Тестирование**: Через Testpool (опционально, ~2 недели)

Для участия в Testpool сообщите модератору о готовности в понедельник до 18:00 МСК.

### Чек-лист перед отправкой:

- [ ] Приложение работает на всех платформах (iOS, Android, Web)
- [ ] Нет критических ошибок и багов
- [ ] Проверена подпись параметров запуска
- [ ] HTTPS настроен (для продакшена)
- [ ] Загружены иконки и обложки
- [ ] Заполнено описание и категория
- [ ] Протестировано на реальных устройствах

## 11. Монетизация (для будущего)

### Варианты монетизации:

1. **Покупки внутри приложения** — через VK Pay
2. **Реклама** — баннеры и прероллы через AppsCentrum
3. **Подписки** — регулярные платежи

💡 Можно комбинировать несколько моделей одновременно.

### Преимущества рекламы VK:

- Отдельный аукцион без конкуренции
- Стоимость показов в 2-3 раза дешевле
- Нативная интеграция в платформу

## 12. Полезные ресурсы

### Официальная документация:

- 📖 [VK Mini Apps документация](https://dev.vk.com/ru/mini-apps/overview)
- 📖 [VK Bridge методы](https://dev.vk.com/bridge/overview)
- 📖 [VK API методы](https://dev.vk.com/ru/method)
- 📖 [VKUI компоненты](https://vkcom.github.io/VKUI/)
- 📖 [VK Icons](https://vkcom.github.io/icons/)

### Инструменты:

- 🛠 [create-vk-mini-app](https://github.com/VKCOM/create-vk-mini-app) — стартовый шаблон
- 🛠 [VK Bridge Sandbox](https://github.com/VKCOM/vk-bridge) — тестирование VK Bridge
- 🛠 [Примеры приложений](https://github.com/VKCOM/vk-mini-apps-examples)
- 🛠 [Параметры запуска (примеры)](https://github.com/VKCOM/vk-apps-launch-params)

### Обучающие материалы:

- 📚 [База знаний VK Mini Apps (Habr)](https://habr.com/ru/company/vk/blog/521192/)
- 📚 [Разработка VK Mini Apps (Habr)](https://habr.com/ru/articles/480974/)
- 📚 [Туториал Selectel](https://selectel.ru/blog/tutorials/vk-mini-apps/)
- 📚 [Академия ВКонтакте](https://vk.com/academy) — мастер-классы и лекции

### Поддержка:

- 💬 [Сообщество VK Mini Apps](https://vk.com/vkappsdev)
- 💬 [Техподдержка для разработчиков](https://vk.com/support?act=home_cat&cat=35)

---

## Текущий статус проекта

✅ **Локальный туннель**: `https://wet-poets-open.loca.lt` (работает)
✅ **Backend**: CORS настроен, API работает
✅ **Frontend**: React + VKUI, онбординг реализован
✅ **База данных**: PostgreSQL с миграциями
✅ **VK App ID**: `54517931`

### Следующие шаги:

1. ✅ Настроить приложение в VK (указать URL)
2. ✅ Открыть приложение через VK
3. ⏳ Протестировать онбординг в VK
4. ⏳ Проверить VK Bridge интеграцию
5. ⏳ Подготовить к модерации

---

## Troubleshooting

### Проблема: Localtunnel показывает предупреждение

**Решение**: Это нормально при первом открытии. Нажмите кнопку **"Continue"**.

### Проблема: Ошибки CORS

**Решение**:
1. Проверьте `.env`: URL должен быть в `ALLOWED_ORIGINS`
2. Перезапустите backend: `docker compose restart backend`
3. Проверьте логи: `docker logs homework_backend`

### Проблема: VK Bridge не работает

**Решение**:
1. Убедитесь, что приложение открыто через VK (а не напрямую)
2. Проверьте инициализацию: `bridge.send('VKWebAppInit')`
3. Проверьте поддержку метода: `bridge.supports('VKWebAppMethodName')`

### Проблема: Параметры запуска не приходят

**Решение**:
1. Откройте приложение через VK (не напрямую по URL)
2. Проверьте URL в консоли: `console.log(window.location.search)`
3. Убедитесь, что режим разработки включен в настройках VK

### Проблема: Не работает на мобильном

**Решение**:
1. Проверьте метатег viewport в HTML
2. Тестируйте в реальном приложении VK (не в браузере)
3. Проверьте логи в Safari/Chrome Remote Debugging

---

**Источники:**

- [VK Mini Apps: база знаний](https://habr.com/ru/company/vk/blog/521192/)
- [Разработка VK Mini Apps](https://habr.com/ru/articles/480974/)
- [Туториал Selectel](https://selectel.ru/blog/tutorials/vk-mini-apps/)
- [VK Apps Launch Params](https://github.com/VKCOM/vk-apps-launch-params)
- [Create VK Mini App](https://github.com/VKCOM/create-vk-mini-app)

---

Если возникнут проблемы, проверьте:
- Frontend логи: DevTools → Console
- Backend логи: `docker logs homework_backend -f`
- Localtunnel статус: процессы с PID 15726, 15058 должны работать
