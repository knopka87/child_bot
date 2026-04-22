# Сравнение VK Mini Apps и MAX Mini Apps

## VK Mini Apps

### Технические требования

| Параметр | Требование |
|----------|------------|
| Протокол | **Только HTTPS** |
| Хостинг | Внешний сервер с публичным URL |
| SSL | Обязательный сертификат |
| Entry point | `index.html` |

### SDK и библиотеки

```bash
npm install @vkontakte/vk-bridge @vkontakte/vkui
```

| Библиотека | Назначение |
|------------|------------|
| **VK Bridge** | Мост к API ВК и функциям устройства |
| **VKUI** | 100+ адаптивных UI компонентов |

**VK Bridge методы:**
- `VKWebAppGetUserInfo` — информация о пользователе
- `VKWebAppStorageGet/Set` — хранилище данных
- Доступ к камере, вибрации, списку друзей
- **100+ методов и событий**

### Модерация и публикация

- Бета-тестирование: **минимум 3 дня**
- Модерация командой VK
- Инструмент **VK Tunnel** для локальной разработки

### Документация

- [VK Mini Apps: обзор возможностей](https://habr.com/ru/companies/vk/articles/961286/)
- [Создание VK Mini App: туториал](https://selectel.ru/blog/tutorials/vk-mini-apps/)
- [VK для разработчиков](https://dev.vk.com/ru/mini-apps/overview)

---

## MAX Mini Apps

### Технические требования

| Параметр | Требование |
|----------|------------|
| Протокол | **Только HTTPS** |
| URL длина | До 1024 символов |
| Символы в URL | Латиница, цифры, точка, дефис |
| Payload (диплинки) | До 512 символов |

### SDK и библиотеки

```html
<script src="https://st.max.ru/js/max-web-app.js"></script>
```

| Библиотека | Назначение |
|------------|------------|
| **MAX Bridge** | Мост к API MAX и функциям устройства |
| **MAX UI** | React-компоненты в стиле MAX |

**MAX Bridge API (`window.WebApp`):**
- `ready()` — уведомление о готовности
- `close()` — закрытие приложения
- `openLink()`, `openMaxLink()` — переходы
- `shareContent()` — шеринг
- `requestContact()` — запрос телефона
- `BackButton` — управление навигацией
- `HapticFeedback` — вибрация (soft, light, medium, heavy)
- `DeviceStorage`, `SecureStorage` — хранилище
- `BiometricManager` — биометрия

### Модерация и публикация

**КРИТИЧНО**: Публикация доступна только для:
- Юридические лица РФ
- ИП — резиденты РФ
- ~~Физлица~~ — **НЕТ**
- ~~Самозанятые~~ — **НЕТ**
- ~~Нерезиденты РФ~~ — **НЕТ**

**Требования к карточке:**

| Параметр | Требование |
|----------|------------|
| Название | 1–59 символов, без эмодзи |
| Ник бота | 11–60 символов, окончание `_bot` или `bot` |
| Сайт | Обязателен, HTTPS, до 1024 символов |
| Логотип | 500×500 px, JPG/PNG, до 5 МБ |
| Описание | До 200 символов |

**Комплаенс обязателен:**
- Пользовательское соглашение
- Политика конфиденциальности
- Сведения о правообладателе
- Канал поддержки

**Верификация:** до 48 часов (рабочие дни)

### Документация

- [MAX для разработчиков](https://dev.max.ru/)
- [Подключение мини-приложения](https://dev.max.ru/docs/webapps/introduction)
- [MAX Bridge SDK](https://dev.max.ru/docs/webapps/bridge)
- [Правила публикации (Хабр)](https://habr.com/ru/articles/951326/)

---

## Сравнительная таблица

| Критерий | VK Mini Apps | MAX Mini Apps |
|----------|--------------|---------------|
| **Протокол** | HTTPS | HTTPS |
| **SDK** | VK Bridge | MAX Bridge |
| **UI Kit** | VKUI (100+ компонентов) | MAX UI (React) |
| **Методы API** | 100+ | ~20-30 |
| **Физлица** | Да | Нет |
| **ИП** | Да | Да |
| **Юрлица** | Да | Да |
| **Нерезиденты** | Да | Нет |
| **Модерация** | 3+ дней | до 48 часов |
| **Документация** | Полная | Базовая |
| **Зрелость** | Высокая | Развивается |

---

## Рекомендации для проекта

### 1. Приоритет — VK Mini Apps

- Более зрелая платформа
- Богаче API (100+ методов)
- Меньше юридических ограничений
- Лучшая документация

### 2. MAX Mini Apps — позже

- Требует юрлицо/ИП
- API похож на Telegram (легко портировать)
- Платформа молодая, API может меняться

### 3. Общий код — абстракция платформы

Оба используют похожий паттерн Bridge:

```tsx
// VK
import bridge from '@vkontakte/vk-bridge';
bridge.send('VKWebAppGetUserInfo');

// MAX
window.WebApp.ready();
window.WebApp.requestContact();
```

Можно создать абстракцию:

```tsx
// src/platform/bridge.ts
import bridge from '@vkontakte/vk-bridge';

type Platform = 'vk' | 'max' | 'web';

function detectPlatform(): Platform {
  if (window.WebApp) return 'max';
  if (window.vkBridge) return 'vk';
  return 'web';
}

export const platform = detectPlatform();

export const platformBridge = {
  init: async () => {
    if (platform === 'vk') {
      await bridge.send('VKWebAppInit');
    } else if (platform === 'max') {
      window.WebApp.ready();
    }
  },

  getUser: async () => {
    if (platform === 'vk') {
      return bridge.send('VKWebAppGetUserInfo');
    } else if (platform === 'max') {
      // MAX не имеет прямого метода получения пользователя
      // данные приходят через initData
      return window.WebApp.initDataUnsafe?.user;
    }
    return null;
  },

  hapticFeedback: (type: 'success' | 'error' | 'warning') => {
    if (platform === 'vk') {
      bridge.send('VKWebAppTapticNotificationOccurred', { type });
    } else if (platform === 'max') {
      window.WebApp.HapticFeedback.notificationOccurred(type);
    }
  },

  close: () => {
    if (platform === 'vk') {
      bridge.send('VKWebAppClose', { status: 'success' });
    } else if (platform === 'max') {
      window.WebApp.close();
    }
  },

  openLink: (url: string) => {
    if (platform === 'vk') {
      bridge.send('VKWebAppOpenLink', { url });
    } else if (platform === 'max') {
      window.WebApp.openLink(url);
    } else {
      window.open(url, '_blank');
    }
  },
};
```

### 4. Структура проекта для мультиплатформы

```
src/
├── platform/
│   ├── bridge.ts           # Абстракция платформы
│   ├── vk.ts               # VK-специфичный код
│   ├── max.ts              # MAX-специфичный код
│   └── types.ts            # Общие типы
├── components/
│   ├── ui/                 # Кроссплатформенные компоненты
│   └── platform/
│       ├── VKComponents/   # VKUI обёртки
│       └── MAXComponents/  # MAX UI обёртки
├── hooks/
│   ├── usePlatform.ts      # Хук определения платформы
│   └── useUser.ts          # Хук получения пользователя
└── app/
    └── App.tsx
```

---

## Чеклист перед публикацией

### VK Mini Apps

- [ ] Приложение работает по HTTPS
- [ ] Настроен SSL сертификат
- [ ] Подключён VK Bridge
- [ ] Используется VKUI для UI
- [ ] Протестировано через VK Tunnel
- [ ] Заполнена карточка приложения
- [ ] Пройдено бета-тестирование (3+ дня)

### MAX Mini Apps

- [ ] Есть юрлицо/ИП в РФ
- [ ] Пройдена верификация организации
- [ ] Приложение работает по HTTPS
- [ ] Подключён MAX Bridge (`window.WebApp.ready()`)
- [ ] Работает BackButton
- [ ] Логотип 500x500 px
- [ ] Есть пользовательское соглашение
- [ ] Есть политика конфиденциальности
- [ ] Указаны сведения о правообладателе
- [ ] Настроен канал поддержки