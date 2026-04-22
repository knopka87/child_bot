# VK Mini Apps - Best Practices и Рекомендации (2026)

## Содержание

1. [VK Bridge API - Интеграция](#1-vk-bridge-api---интеграция)
2. [VKUI - Компоненты и паттерны](#2-vkui---компоненты-и-паттерны)
3. [Оптимизация производительности](#3-оптимизация-производительности)
4. [VK-специфичные требования](#4-vk-специфичные-требования)
5. [Монетизация и VK Pay](#5-монетизация-и-vk-pay)
6. [Аналитика и метрики](#6-аналитика-и-метрики)
7. [PWA и Offline Support](#7-pwa-и-offline-support)
8. [Реклама (Ads SDK)](#8-реклама-ads-sdk)
9. [Требования безопасности](#9-требования-безопасности)
10. [Процесс модерации](#10-процесс-модерации)

---

## 1. VK Bridge API - Интеграция

### Описание

VK Bridge — это пакет для интеграции VK Mini Apps с официальными клиентами VK для iOS, Android и Web. Он предоставляет API для взаимодействия вашего веб-приложения с нативными функциями платформы.

### Установка

```bash
npm install @vkontakte/vk-bridge
```

### Базовая инициализация

```typescript
import bridge from '@vkontakte/vk-bridge';

// Инициализация моста
bridge.send('VKWebAppInit');
```

### Основные методы

#### 1. Получение данных пользователя

```typescript
// Получение базовой информации о пользователе
const getUserInfo = async () => {
  try {
    const user = await bridge.send('VKWebAppGetUserInfo');
    console.log(user);
  } catch (error) {
    console.error('Failed to get user info', error);
  }
};
```

#### 2. Запрос токена доступа

```typescript
// Запрос прав доступа
const getAccessToken = async () => {
  try {
    const data = await bridge.send('VKWebAppGetAuthToken', {
      app_id: YOUR_APP_ID,
      scope: 'friends,photos'
    });
    // data.access_token - никогда не храните в localStorage!
    // Отправьте на бэкенд для проверки
    return data.access_token;
  } catch (error) {
    console.error('Auth failed', error);
  }
};
```

#### 3. Тактильная обратная связь (Haptics)

```typescript
// Легкая вибрация
bridge.send('VKWebAppTapticImpactOccurred', { style: 'light' });

// Средняя вибрация
bridge.send('VKWebAppTapticImpactOccurred', { style: 'medium' });

// Сильная вибрация
bridge.send('VKWebAppTapticImpactOccurred', { style: 'heavy' });

// Уведомление об успехе
bridge.send('VKWebAppTapticNotificationOccurred', { type: 'success' });

// Уведомление об ошибке
bridge.send('VKWebAppTapticNotificationOccurred', { type: 'error' });

// Изменение выбора
bridge.send('VKWebAppTapticSelectionChanged');
```

#### 4. Шаринг и Stories

```typescript
// Поделиться ссылкой
const share = async () => {
  try {
    await bridge.send('VKWebAppShare', {
      link: 'https://vk.com/app123456'
    });
  } catch (error) {
    console.error('Share failed', error);
  }
};

// Подписка на обновления в истории
const subscribeToStory = async () => {
  try {
    await bridge.send('VKWebAppSubscribeStoryApp');
  } catch (error) {
    console.error('Story subscription failed', error);
  }
};
```

#### 5. Копирование текста

```typescript
const copyToClipboard = async (text: string) => {
  try {
    await bridge.send('VKWebAppCopyText', { text });
  } catch (error) {
    console.error('Copy failed', error);
  }
};
```

#### 6. Сканирование QR-кода

```typescript
const scanQR = async () => {
  try {
    const data = await bridge.send('VKWebAppOpenCodeReader');
    console.log('QR Code data:', data.code_data);
  } catch (error) {
    console.error('QR scan failed', error);
  }
};
```

### Подписка на события

```typescript
// Подписка на события
bridge.subscribe((e) => {
  if (!e.detail) {
    return;
  }

  const { type, data } = e.detail;

  switch (type) {
    case 'VKWebAppUpdateConfig':
      // Обновление конфигурации (тема, язык)
      const scheme = data.scheme; // 'bright_light' | 'space_gray'
      document.body.setAttribute('scheme', scheme);
      break;

    case 'VKWebAppViewHide':
      // Приложение свернуто
      console.log('App hidden');
      break;

    case 'VKWebAppViewRestore':
      // Приложение восстановлено
      console.log('App restored');
      break;
  }
});
```

### Middleware для логирования

```typescript
// Создание middleware для логирования всех событий
const loggingMiddleware = (next) => async (method, props) => {
  console.log(`[VK Bridge] Sending: ${method}`, props);

  try {
    const result = await next(method, props);
    console.log(`[VK Bridge] Success: ${method}`, result);
    return result;
  } catch (error) {
    console.error(`[VK Bridge] Error: ${method}`, error);
    throw error;
  }
};

// Применение middleware
bridge.applyMiddleware(loggingMiddleware);
```

### Best Practices

✅ **DO:**
- Всегда оборачивайте вызовы `bridge.send()` в try-catch
- Используйте тактильную обратную связь для улучшения UX
- Инициализируйте bridge сразу при загрузке приложения
- Используйте middleware для централизованного логирования и обработки ошибок
- Проверяйте поддержку методов через `bridge.supports(method)`

❌ **DON'T:**
- Не храните access_token в localStorage или sessionStorage
- Не игнорируйте ошибки от Bridge API
- Не вызывайте VK API напрямую с клиента - используйте ваш бэкенд
- Не злоупотребляйте haptic feedback (каждая вибрация снижает заряд батареи)

### Тестирование

Для разработки и тестирования используйте `@vkontakte/vk-bridge-mock`:

```typescript
import { mockBridge } from '@vkontakte/vk-bridge-mock';

if (process.env.NODE_ENV === 'development') {
  mockBridge({
    VKWebAppGetUserInfo: {
      id: 123456,
      first_name: 'Test',
      last_name: 'User',
      photo_200: 'https://example.com/photo.jpg'
    }
  });
}
```

---

## 2. VKUI - Компоненты и паттерны

### Описание

VKUI — это библиотека адаптивных React-компонентов для создания интерфейсов, внешне неотличимых от iOS и Android приложений VK. Библиотека автоматически адаптируется к платформе пользователя (iOS, Android, Web).

### Установка

```bash
npm install @vkontakte/vkui
```

### Основная структура приложения

```typescript
import {
  ConfigProvider,
  AdaptivityProvider,
  AppRoot,
  SplitLayout,
  SplitCol,
  View,
  Panel,
  PanelHeader,
  Group,
  Header,
  SimpleCell
} from '@vkontakte/vkui';
import '@vkontakte/vkui/dist/vkui.css';

const App = () => {
  const platform = usePlatform();

  return (
    <ConfigProvider>
      <AdaptivityProvider>
        <AppRoot>
          <SplitLayout
            header={platform !== 'vkcom' && <PanelHeader delimiter="none" />}
          >
            <SplitCol autoSpaced>
              <View activePanel="main">
                <Panel id="main">
                  <PanelHeader>Моё приложение</PanelHeader>

                  <Group header={<Header size="s">Список</Header>}>
                    <SimpleCell>Элемент 1</SimpleCell>
                    <SimpleCell>Элемент 2</SimpleCell>
                    <SimpleCell>Элемент 3</SimpleCell>
                  </Group>
                </Panel>
              </View>
            </SplitCol>
          </SplitLayout>
        </AppRoot>
      </AdaptivityProvider>
    </ConfigProvider>
  );
};
```

### Иерархия компонентов

```
ConfigProvider (глобальные настройки)
  └─ AdaptivityProvider (адаптивность)
      └─ AppRoot (корневой контейнер)
          └─ SplitLayout (макет с колонками)
              └─ SplitCol (колонка)
                  └─ View (группа панелей)
                      └─ Panel (экран)
                          └─ PanelHeader (шапка)
                          └─ Group (группа контента)
```

### Адаптивность

VKUI автоматически адаптируется под размер экрана через 4 свойства:

```typescript
import { useAdaptivity } from '@vkontakte/vkui';

const MyComponent = () => {
  const { viewWidth, viewHeight, sizeX, sizeY } = useAdaptivity();

  // viewWidth: 'mobile' | 'tablet' | 'desktop'
  // sizeX: 'compact' | 'regular'

  return (
    <div>
      {viewWidth === 'mobile' ? (
        <MobileLayout />
      ) : (
        <DesktopLayout />
      )}
    </div>
  );
};
```

### Навигация между панелями

```typescript
import { useState } from 'react';
import { View, Panel, PanelHeader, Button } from '@vkontakte/vkui';

const Navigation = () => {
  const [activePanel, setActivePanel] = useState('panel1');

  return (
    <View activePanel={activePanel}>
      <Panel id="panel1">
        <PanelHeader>Панель 1</PanelHeader>
        <Button onClick={() => setActivePanel('panel2')}>
          Перейти на панель 2
        </Button>
      </Panel>

      <Panel id="panel2">
        <PanelHeader before={<PanelHeaderBack onClick={() => setActivePanel('panel1')} />}>
          Панель 2
        </PanelHeader>
        <div>Содержимое второй панели</div>
      </Panel>
    </View>
  );
};
```

### Модальные окна

```typescript
import { useState } from 'react';
import {
  ModalRoot,
  ModalPage,
  ModalPageHeader,
  Button,
  Group,
  Div
} from '@vkontakte/vkui';

const App = () => {
  const [activeModal, setActiveModal] = useState(null);

  const modal = (
    <ModalRoot activeModal={activeModal} onClose={() => setActiveModal(null)}>
      <ModalPage
        id="modal1"
        header={<ModalPageHeader>Модальное окно</ModalPageHeader>}
      >
        <Group>
          <Div>Содержимое модального окна</Div>
        </Group>
      </ModalPage>
    </ModalRoot>
  );

  return (
    <View activePanel="main" modal={modal}>
      <Panel id="main">
        <PanelHeader>Главная</PanelHeader>
        <Button onClick={() => setActiveModal('modal1')}>
          Открыть модальное окно
        </Button>
      </Panel>
    </View>
  );
};
```

### Snackbar для уведомлений

```typescript
import { useState } from 'react';
import { Snackbar, Avatar, Button } from '@vkontakte/vkui';
import { Icon16Done } from '@vkontakte/icons';

const MyComponent = () => {
  const [snackbar, setSnackbar] = useState(null);

  const showSuccess = () => {
    setSnackbar(
      <Snackbar
        onClose={() => setSnackbar(null)}
        before={
          <Avatar size={24} style={{ background: 'var(--vkui--color_background_positive)' }}>
            <Icon16Done fill="#fff" width={14} height={14} />
          </Avatar>
        }
      >
        Операция выполнена успешно
      </Snackbar>
    );
  };

  return (
    <div>
      <Button onClick={showSuccess}>Показать уведомление</Button>
      {snackbar}
    </div>
  );
};
```

### Формы и инпуты

```typescript
import {
  FormItem,
  Input,
  Textarea,
  Checkbox,
  Radio,
  Select,
  Button,
  FormLayout
} from '@vkontakte/vkui';

const MyForm = () => {
  return (
    <FormLayout>
      <FormItem top="Имя" required>
        <Input placeholder="Введите имя" />
      </FormItem>

      <FormItem top="Email">
        <Input type="email" placeholder="example@mail.ru" />
      </FormItem>

      <FormItem top="Описание">
        <Textarea placeholder="Расскажите о себе" />
      </FormItem>

      <FormItem top="Город">
        <Select>
          <option value="1">Москва</option>
          <option value="2">Санкт-Петербург</option>
        </Select>
      </FormItem>

      <FormItem>
        <Checkbox>Согласен с условиями</Checkbox>
      </FormItem>

      <FormItem>
        <Button size="l" stretched>
          Отправить
        </Button>
      </FormItem>
    </FormLayout>
  );
};
```

### Списки

```typescript
import { Group, Header, SimpleCell, Avatar } from '@vkontakte/vkui';
import { Icon28UserCircleOutline } from '@vkontakte/icons';

const UserList = ({ users }) => {
  return (
    <Group header={<Header>Пользователи</Header>}>
      {users.map(user => (
        <SimpleCell
          key={user.id}
          before={<Avatar size={40} src={user.photo} />}
          after={<Icon28UserCircleOutline />}
          subtitle={user.subtitle}
        >
          {user.name}
        </SimpleCell>
      ))}
    </Group>
  );
};
```

### Цветовые токены и темизация

```typescript
// Использование CSS-переменных VKUI
const StyledComponent = styled.div`
  background: var(--vkui--color_background);
  color: var(--vkui--color_text_primary);
  border: 1px solid var(--vkui--color_separator_primary);
  padding: var(--vkui--spacing_size_m);
`;

// Динамическое изменение темы
import { useAppearance } from '@vkontakte/vkui';

const ThemeSwitcher = () => {
  const appearance = useAppearance(); // 'light' | 'dark'

  return <div>Текущая тема: {appearance}</div>;
};
```

### Icons

VKUI содержит более 2500 иконок:

```typescript
import {
  Icon28HomeOutline,
  Icon28MessageOutline,
  Icon28NotificationOutline,
  Icon28UserCircleOutline,
  Icon24Add,
  Icon16Done
} from '@vkontakte/icons';

// Использование
<Icon28HomeOutline />
<Icon24Add width={20} height={20} />
```

### Best Practices

✅ **DO:**
- Используйте `ConfigProvider` и `AdaptivityProvider` в корне приложения
- Следуйте иерархии: View → Panel → Group → Cell
- Используйте встроенные компоненты вместо кастомных
- Применяйте цветовые токены (CSS-переменные) вместо хардкода цветов
- Используйте иконки из `@vkontakte/icons`
- Добавляйте Taptic Feedback при нажатиях на кнопки
- Показывайте Snackbar для подтверждения действий

❌ **DON'T:**
- Не используйте сторонние UI-библиотеки (Material UI, Ant Design)
- Не нарушайте структуру компонентов (Panel внутри Panel)
- Не используйте фиксированные цвета (#fff, #000) - используйте токены
- Не создавайте слишком глубокую вложенность модальных окон (более 2-3)

### Миграция на v8

VKUI v8 находится в разработке. Основные изменения:
- Переход на `slotProps` в компонентах (Switch, Checkbox, Radio, Input и др.)
- Улучшенная типизация TypeScript
- Оптимизация производительности

Следите за обновлениями: https://github.com/VKCOM/VKUI/releases

---

## 3. Оптимизация производительности

### Bundle Size

VK рекомендует держать размер bundle не более **10 MB** (после сжатия).

### Code Splitting

#### Динамические импорты

```typescript
// Ленивая загрузка компонентов
import { lazy, Suspense } from 'react';

const HeavyComponent = lazy(() => import('./HeavyComponent'));

const App = () => (
  <Suspense fallback={<Spinner />}>
    <HeavyComponent />
  </Suspense>
);
```

#### Route-based splitting

```typescript
import { lazy } from 'react';

const routes = [
  {
    path: '/home',
    component: lazy(() => import('./pages/Home'))
  },
  {
    path: '/profile',
    component: lazy(() => import('./pages/Profile'))
  }
];
```

### Vite оптимизация

```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          // Разделение vendor библиотек
          'react-vendor': ['react', 'react-dom'],
          'vk-vendor': ['@vkontakte/vkui', '@vkontakte/vk-bridge'],
        }
      }
    },
    // Минимальный размер chunk для split
    chunkSizeWarningLimit: 500,
    // Включение минификации
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true, // Удаление console.log в production
      }
    }
  },
  // CSS code splitting
  css: {
    modules: {
      localsConvention: 'camelCase'
    }
  }
});
```

### Webpack оптимизация

```javascript
// webpack.config.js
module.exports = {
  optimization: {
    splitChunks: {
      chunks: 'all',
      cacheGroups: {
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          name: 'vendors',
          priority: 10
        },
        vk: {
          test: /[\\/]node_modules[\\/]@vkontakte[\\/]/,
          name: 'vk',
          priority: 20
        }
      }
    },
    runtimeChunk: 'single',
    minimize: true,
    usedExports: true // Tree shaking
  }
};
```

### Анализ bundle

```bash
# Для Vite
npm install --save-dev rollup-plugin-visualizer

# Для Webpack
npm install --save-dev webpack-bundle-analyzer
```

```typescript
// vite.config.ts
import { visualizer } from 'rollup-plugin-visualizer';

export default defineConfig({
  plugins: [
    visualizer({
      open: true,
      filename: 'dist/stats.html'
    })
  ]
});
```

### Оптимизация изображений

```typescript
// Используйте современные форматы
<img src="photo.webp" alt="Photo" />

// Lazy loading для изображений
<img src="photo.jpg" loading="lazy" />

// VKUI Avatar с fallback
<Avatar size={48} src={user.photo_200} fallbackIcon={<Icon28UserCircleOutline />} />
```

### Оптимизация анимаций

```css
/* GPU-ускоренные свойства */
.animated-element {
  /* ✅ Используйте transform и opacity */
  transform: translateX(100px);
  opacity: 0.5;
  will-change: transform, opacity;

  /* ❌ Избегайте */
  /* box-shadow: 0 0 10px rgba(0,0,0,0.5); */
  /* filter: blur(5px); */
}

/* Использование CSS-переменных для transition */
.button {
  transition: background-color var(--vkui--animation_duration_m);
}
```

### Debounce и Throttle

```typescript
import { useDebouncedCallback } from 'use-debounce';

const SearchInput = () => {
  const [query, setQuery] = useState('');

  const debouncedSearch = useDebouncedCallback(
    (value) => {
      // API запрос
      searchAPI(value);
    },
    500 // 500ms задержка
  );

  const handleChange = (e) => {
    const value = e.target.value;
    setQuery(value);
    debouncedSearch(value);
  };

  return <Input value={query} onChange={handleChange} />;
};
```

### Мемоизация

```typescript
import { memo, useMemo, useCallback } from 'react';

// Мемоизация компонента
const ExpensiveComponent = memo(({ data }) => {
  return <div>{/* render */}</div>;
});

// Мемоизация вычислений
const MyComponent = ({ items }) => {
  const sortedItems = useMemo(() => {
    return items.sort((a, b) => a.name.localeCompare(b.name));
  }, [items]);

  // Мемоизация callback
  const handleClick = useCallback(() => {
    console.log('clicked');
  }, []);

  return <div>{sortedItems.map(item => <div key={item.id}>{item.name}</div>)}</div>;
};
```

### Виртуализация списков

```typescript
import { FixedSizeList } from 'react-window';

const VirtualizedList = ({ items }) => {
  const Row = ({ index, style }) => (
    <div style={style}>
      <SimpleCell>{items[index].name}</SimpleCell>
    </div>
  );

  return (
    <FixedSizeList
      height={600}
      itemCount={items.length}
      itemSize={48}
      width="100%"
    >
      {Row}
    </FixedSizeList>
  );
};
```

### Performance Monitoring

```typescript
import { useEffect } from 'react';

// Отслеживание производительности
const usePerformanceMonitor = (componentName: string) => {
  useEffect(() => {
    const start = performance.now();

    return () => {
      const end = performance.now();
      console.log(`${componentName} render time: ${end - start}ms`);
    };
  });
};

// Web Vitals
import { getCLS, getFID, getFCP, getLCP, getTTFB } from 'web-vitals';

getCLS(console.log);
getFID(console.log);
getFCP(console.log);
getLCP(console.log);
getTTFB(console.log);
```

### Checklist оптимизации

- [ ] Bundle size < 10 MB
- [ ] Code splitting по роутам
- [ ] Tree shaking включен
- [ ] Минификация в production
- [ ] Lazy loading для изображений
- [ ] Виртуализация для длинных списков (>100 элементов)
- [ ] Debounce для поиска и автокомплита
- [ ] Мемоизация тяжелых вычислений
- [ ] GPU-ускоренные анимации (transform, opacity)
- [ ] Удалены console.log в production
- [ ] Используются современные форматы изображений (WebP, AVIF)

---

## 4. VK-специфичные требования

### Размер приложения

- **Максимальный размер**: до 10 MB (после сжатия)
- **Рекомендуемый**: до 3-5 MB для быстрой загрузки

### HTTPS обязательно

SSL-сертификаты являются **обязательным требованием**. Мини-приложения должны открываться через HTTPS, иначе они не будут работать внутри мобильного клиента VK.

```bash
# Для разработки используйте mkcert
brew install mkcert
mkcert -install
mkcert localhost
```

```typescript
// vite.config.ts для локальной разработки с HTTPS
import { defineConfig } from 'vite';
import fs from 'fs';

export default defineConfig({
  server: {
    https: {
      key: fs.readFileSync('./localhost-key.pem'),
      cert: fs.readFileSync('./localhost.pem')
    },
    host: '0.0.0.0',
    port: 10888
  }
});
```

### SPA роутинг

Приложение должно поддерживать SPA-роутинг с обработкой 404:

```nginx
# nginx конфигурация
location / {
  try_files $uri $uri/ /index.html;
}
```

### Проверка параметров запуска (launch_params)

Каждый запуск приложения содержит параметры, которые **обязательно** нужно проверять на бэкенде:

```typescript
// Клиент (получение параметров)
const params = new URLSearchParams(window.location.search);
const launchParams = params.toString();

// Отправка на бэкенд для проверки
const response = await fetch('/api/auth', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ launchParams })
});
```

```go
// Бэкенд (Go) - проверка подписи
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "errors"
    "net/url"
    "sort"
    "strings"
)

func VerifyLaunchParams(launchParams string, secretKey string) (bool, error) {
    params, err := url.ParseQuery(launchParams)
    if err != nil {
        return false, err
    }

    sign := params.Get("sign")
    if sign == "" {
        return false, errors.New("sign parameter is missing")
    }

    // Удаляем sign из параметров
    params.Del("sign")

    // Сортируем параметры
    keys := make([]string, 0, len(params))
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    // Формируем строку для проверки
    var queryString strings.Builder
    for i, k := range keys {
        if i > 0 {
            queryString.WriteString("&")
        }
        queryString.WriteString(k)
        queryString.WriteString("=")
        queryString.WriteString(params.Get(k))
    }

    // Вычисляем HMAC-SHA256
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(queryString.String()))
    expectedSign := base64.URLEncoding.EncodeToString(h.Sum(nil))

    // Удаляем padding и заменяем символы для URL-safe base64
    expectedSign = strings.TrimRight(expectedSign, "=")
    expectedSign = strings.ReplaceAll(expectedSign, "+", "-")
    expectedSign = strings.ReplaceAll(expectedSign, "/", "_")

    return sign == expectedSign, nil
}
```

Для проверки подписи есть готовые библиотеки:
- Go: https://github.com/kravetsone/vk-launch-params
- Node.js: https://github.com/VKCOM/vk-apps-launch-params

### VK ID (авторизация)

Параметр `sign` — это HMAC-хэш, созданный на стороне VK. Задача разработчика — проверить подпись на сервере и подтвердить, что пользователь подлинный.

**Критические нарушения безопасности:**
- ❌ Хранение `access_token` в localStorage
- ❌ Проверка `sign` на клиенте вместо сервера
- ❌ Отсутствие HTTPS
- ❌ Отправка открытых данных без подписи

### Cookie Policy

Сервер мини-приложения **не должен** устанавливать cookie-файлы в ответах на запросы загрузки веб-приложения и последующие Same Origin запросы.

### Требования к UX

#### Модель View → Panel → Group → Cell

Мини-приложения должны следовать структуре:
- **View** — группа панелей
- **Panel** — отдельный экран
- **Group** — группа контента
- **Cell** — элемент списка

Это упрощает управление навигацией и снижает нагрузку на DOM.

#### Типографика

Используйте динамическое масштабирование:

```css
.text {
  font-size: clamp(14px, 2vw, 18px);
}
```

#### Обратная связь

- Добавляйте `TapticFeedback` для подтверждения действий
- Используйте `Snackbar` для уведомлений
- Показывайте `Spinner` при загрузке

### Размер экранов

Поддерживайте все размеры экранов:
- **Mobile**: < 768px
- **Tablet**: 768px - 1024px
- **Desktop**: > 1024px

```typescript
import { useAdaptivity } from '@vkontakte/vkui';

const MyComponent = () => {
  const { viewWidth } = useAdaptivity();

  return (
    <div>
      {viewWidth === 'mobile' && <MobileLayout />}
      {viewWidth === 'tablet' && <TabletLayout />}
      {viewWidth === 'desktop' && <DesktopLayout />}
    </div>
  );
};
```

### Checklist требований

- [ ] HTTPS включен
- [ ] Bundle size < 10 MB
- [ ] Проверка `sign` на бэкенде
- [ ] `access_token` не хранится в браузере
- [ ] SPA роутинг с 404 → index.html
- [ ] Cookie не устанавливаются
- [ ] Поддержка всех размеров экранов
- [ ] Используется VKUI для UI
- [ ] Taptic feedback добавлен
- [ ] Snackbar для подтверждений

---

## 5. Монетизация и VK Pay

### Обзор

VK Pay — это платежная система VK, интегрированная в мини-приложения. Она позволяет принимать платежи с карт, VK кошельков и других методов.

### Статистика

- **Рост конверсии**: +23% после интеграции VK Pay (данные Cloud Mail.ru)
- **Рост in-app ads**: +550% YoY (Q3)
- **MAU топовых Mini Apps**: >1 млн (VK Classifieds, AliExpress, VK Taxi, VK Food)

### Интеграция VK Pay

#### 1. Регистрация в VK Pay

1. Зарегистрируйтесь на https://pay.vk.com
2. Получите Merchant ID
3. Настройте webhook для уведомлений о платежах

#### 2. Инициация платежа

```typescript
import bridge from '@vkontakte/vk-bridge';

const initPayment = async () => {
  try {
    const result = await bridge.send('VKWebAppOpenPayForm', {
      app_id: YOUR_APP_ID,
      action: 'pay-to-service',
      params: {
        amount: 100, // в рублях
        description: 'Покупка премиум подписки',
        merchant_id: YOUR_MERCHANT_ID,
        sign: calculateSign() // подпись от бэкенда
      }
    });

    if (result.status === 'success') {
      console.log('Payment successful:', result);
    }
  } catch (error) {
    console.error('Payment failed:', error);
  }
};
```

#### 3. Обработка платежей на бэкенде

```go
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "net/http"
)

type PaymentNotification struct {
    NotificationType string `json:"notification_type"`
    AppID            int    `json:"app_id"`
    UserID           int    `json:"user_id"`
    ReceiverID       int    `json:"receiver_id"`
    OrderID          int    `json:"order_id"`
    Amount           int    `json:"amount"`
    ItemID           string `json:"item_id"`
    Sign             string `json:"sig"`
}

func verifyPaymentSign(notification PaymentNotification, secretKey string) bool {
    data := fmt.Sprintf(
        "%s_%d_%d_%d_%d_%d_%s",
        notification.NotificationType,
        notification.AppID,
        notification.UserID,
        notification.ReceiverID,
        notification.OrderID,
        notification.Amount,
        notification.ItemID,
    )

    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(data))
    expectedSign := hex.EncodeToString(h.Sum(nil))

    return notification.Sign == expectedSign
}

func handlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
    var notification PaymentNotification
    if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if !verifyPaymentSign(notification, SECRET_KEY) {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    // Обработка успешного платежа
    if notification.NotificationType == "get_item" {
        // Выдаем товар пользователю
        grantItemToUser(notification.UserID, notification.ItemID)
    }

    // Отправляем подтверждение VK
    json.NewEncoder(w).Encode(map[string]interface{}{
        "response": map[string]interface{}{
            "order_id": notification.OrderID,
            "item_id":  notification.ItemID,
        },
    })
}
```

### Форматы монетизации

1. **Прямые платежи**
   - Подписки
   - Внутриигровые покупки
   - Премиум функции

2. **Реклама**
   - VK Ads SDK
   - Нативная реклама
   - Rewarded ads (вознаграждение за просмотр)

3. **Комиссионная модель**
   - Маркетплейсы
   - Бронирования
   - Заказ услуг

### Best Practices

✅ **DO:**
- Используйте VK Pay для платежей (лучшая конверсия)
- Проверяйте подпись платежей на бэкенде
- Храните историю транзакций
- Предлагайте разные методы оплаты
- Показывайте прозрачные цены
- Используйте подписки для recurring revenue

❌ **DON'T:**
- Не принимайте платежи напрямую (в обход VK Pay)
- Не храните данные карт пользователей
- Не обрабатывайте платежи на клиенте
- Не скрывайте комиссии и дополнительные платежи

### Аналитика платежей

Отслеживайте ключевые метрики:
- **ARPU** (Average Revenue Per User)
- **LTV** (Lifetime Value)
- **Conversion Rate** (% пользователей, совершивших покупку)
- **ARPPU** (Average Revenue Per Paying User)
- **Churn Rate** (% отказавшихся от подписки)

```typescript
// Отправка события конверсии в VK Ads
bridge.send('VKWebAppConversionHit', {
  pixel_code: 'YOUR_PIXEL_CODE',
  event: 'purchase',
  value: 100,
  currency: 'RUB'
});
```

---

## 6. Аналитика и метрики

### Доступные платформы

1. **VK Ads Analytics** — метрики CTR, CPA, Retention
2. **MyTracker** — аналитика активности пользователей (beta)
3. **Яндекс.Метрика** — события внутри Mini App
4. **Roistat** — сквозная аналитика до продаж

### MyTracker интеграция

MyTracker предоставляет быстрые отчеты по мини-приложению:
- Активность пользователей
- Информация об устройствах
- География
- Запуски
- Retention

**Важно**: Revenue tracking пока не поддерживается.

#### Установка MyTracker

```typescript
// Инициализация MyTracker
import { initTracker } from '@vkontakte/mytracker';

initTracker({
  trackerId: 'YOUR_TRACKER_ID'
});

// Отправка события
import { trackEvent } from '@vkontakte/mytracker';

trackEvent('button_clicked', {
  button_name: 'buy_premium',
  screen: 'home'
});
```

### VK Bridge события

```typescript
import bridge from '@vkontakte/vk-bridge';

// Трекинг экранов
bridge.send('VKWebAppViewRestore').then(() => {
  // Пользователь вернулся в приложение
  trackEvent('app_restored');
});

bridge.send('VKWebAppViewHide').then(() => {
  // Пользователь свернул приложение
  trackEvent('app_hidden');
});
```

### Яндекс.Метрика

```typescript
// Установка
<script type="text/javascript">
   (function(m,e,t,r,i,k,a){m[i]=m[i]||function(){(m[i].a=m[i].a||[]).push(arguments)};
   m[i].l=1*new Date();
   for (var j = 0; j < document.scripts.length; j++) {if (document.scripts[j].src === r) { return; }}
   k=e.createElement(t),a=e.getElementsByTagName(t)[0],k.async=1,k.src=r,a.parentNode.insertBefore(k,a)})
   (window, document, "script", "https://mc.yandex.ru/metrika/tag.js", "ym");

   ym(YOUR_COUNTER_ID, "init", {
        clickmap:true,
        trackLinks:true,
        accurateTrackBounce:true,
        webvisor:true
   });
</script>

// Отправка событий
ym(YOUR_COUNTER_ID, 'reachGoal', 'purchase', {
  value: 100,
  currency: 'RUB'
});
```

### Ключевые метрики для отслеживания

#### Воронка пользователя

```typescript
// Пример отслеживания воронки
const trackFunnel = () => {
  // 1. Запуск приложения
  trackEvent('app_launch');

  // 2. Просмотр товара
  trackEvent('product_view', { product_id: '123' });

  // 3. Добавление в корзину
  trackEvent('add_to_cart', { product_id: '123' });

  // 4. Начало оформления
  trackEvent('checkout_start');

  // 5. Успешная покупка
  trackEvent('purchase', {
    product_id: '123',
    value: 100,
    currency: 'RUB'
  });
};
```

#### Retention

```typescript
// Day 1, Day 7, Day 30 retention
const trackRetention = () => {
  const installDate = localStorage.getItem('install_date');
  if (!installDate) {
    localStorage.setItem('install_date', new Date().toISOString());
    trackEvent('first_launch');
  } else {
    const daysSinceInstall = Math.floor(
      (Date.now() - new Date(installDate).getTime()) / (1000 * 60 * 60 * 24)
    );

    if ([1, 7, 30].includes(daysSinceInstall)) {
      trackEvent(`retention_day_${daysSinceInstall}`);
    }
  }
};
```

#### Performance метрики

```typescript
import { getCLS, getFID, getFCP, getLCP, getTTFB } from 'web-vitals';

const sendToAnalytics = ({ name, value }) => {
  trackEvent('web_vitals', {
    metric_name: name,
    metric_value: value
  });
};

getCLS(sendToAnalytics);
getFID(sendToAnalytics);
getFCP(sendToAnalytics);
getLCP(sendToAnalytics);
getTTFB(sendToAnalytics);
```

### VK Ads Pixel

```typescript
// Установка пикселя VK Ads
<script type="text/javascript">
!function(){var t=document.createElement("script");t.type="text/javascript",t.async=!0,t.src='https://vk.com/js/api/openapi.js?169',t.onload=function(){VK.Retargeting.Init("VK-RTRG-XXXXX-XXXXX"),VK.Retargeting.Hit()},document.head.appendChild(t)}();
</script>

// Отправка события конверсии
VK.Retargeting.Event('purchase');
```

### Настройка трекера в VK Ads

1. Перейдите в VK Ads кабинет
2. Добавьте приложение
3. Настройте события (установка, регистрация, покупка)
4. Получите код пикселя
5. Интегрируйте в приложение

### Best Practices

✅ **DO:**
- Отслеживайте всю воронку пользователя
- Используйте несколько систем аналитики (VK Ads + Метрика)
- Трекайте performance метрики (Web Vitals)
- Сегментируйте пользователей по источникам
- A/B тестируйте изменения
- Отслеживайте ошибки (Sentry, Bugsnag)

❌ **DON'T:**
- Не отправляйте PII (personally identifiable information) в аналитику
- Не блокируйте UI при отправке событий (используйте async)
- Не отправляйте слишком много событий (перегрузка)
- Не игнорируйте GDPR/CCPA требования

### Чеклист событий для трекинга

- [ ] App launch
- [ ] Screen views
- [ ] Button clicks
- [ ] Form submissions
- [ ] Errors
- [ ] Search queries
- [ ] Product views
- [ ] Add to cart
- [ ] Purchase
- [ ] Share
- [ ] Retention (D1, D7, D30)

---

## 7. PWA и Offline Support

### Обзор PWA поддержки в 2026

В 2026 году все основные браузеры полностью поддерживают PWA:
- **Chrome/Edge**: полная поддержка
- **Firefox**: полная поддержка
- **Safari (iOS)**: service workers с ограничениями

### Service Worker

Service Workers — это основа PWA, позволяющая:
- Офлайн функциональность
- Background sync
- Push уведомления

#### Базовый Service Worker

```typescript
// public/sw.js
const CACHE_NAME = 'vk-miniapp-v1';
const urlsToCache = [
  '/',
  '/index.html',
  '/static/css/main.css',
  '/static/js/main.js',
  '/manifest.json'
];

// Установка
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => cache.addAll(urlsToCache))
  );
});

// Активация
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => name !== CACHE_NAME)
          .map((name) => caches.delete(name))
      );
    })
  );
});

// Fetch (стратегия Cache First)
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request)
      .then((response) => {
        if (response) {
          return response; // Возвращаем из кэша
        }

        return fetch(event.request).then((response) => {
          // Кэшируем новые запросы
          if (event.request.method === 'GET') {
            const responseToCache = response.clone();
            caches.open(CACHE_NAME).then((cache) => {
              cache.put(event.request, responseToCache);
            });
          }
          return response;
        });
      })
  );
});
```

#### Регистрация Service Worker

```typescript
// main.tsx
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/sw.js')
      .then((registration) => {
        console.log('SW registered:', registration);
      })
      .catch((error) => {
        console.log('SW registration failed:', error);
      });
  });
}
```

### Web App Manifest

```json
{
  "name": "My VK Mini App",
  "short_name": "MiniApp",
  "description": "Awesome VK Mini App",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#4680c2",
  "icons": [
    {
      "src": "/icons/icon-72x72.png",
      "sizes": "72x72",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-96x96.png",
      "sizes": "96x96",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-128x128.png",
      "sizes": "128x128",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-144x144.png",
      "sizes": "144x144",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-152x152.png",
      "sizes": "152x152",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-192x192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-384x384.png",
      "sizes": "384x384",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-512x512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}
```

### Стратегии кэширования

#### 1. Cache First (для статики)

```javascript
// Сначала кэш, потом сеть
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request)
      .then((response) => response || fetch(event.request))
  );
});
```

#### 2. Network First (для API)

```javascript
// Сначала сеть, потом кэш
self.addEventListener('fetch', (event) => {
  event.respondWith(
    fetch(event.request)
      .catch(() => caches.match(event.request))
  );
});
```

#### 3. Stale While Revalidate

```javascript
// Возвращаем из кэша, но обновляем в фоне
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.match(event.request).then((response) => {
        const fetchPromise = fetch(event.request).then((networkResponse) => {
          cache.put(event.request, networkResponse.clone());
          return networkResponse;
        });
        return response || fetchPromise;
      });
    })
  );
});
```

### Vite PWA Plugin

```bash
npm install vite-plugin-pwa -D
```

```typescript
// vite.config.ts
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
  plugins: [
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.ico', 'robots.txt', 'icons/*.png'],
      manifest: {
        name: 'My VK Mini App',
        short_name: 'MiniApp',
        theme_color: '#4680c2',
        icons: [
          {
            src: '/icons/icon-192x192.png',
            sizes: '192x192',
            type: 'image/png'
          },
          {
            src: '/icons/icon-512x512.png',
            sizes: '512x512',
            type: 'image/png'
          }
        ]
      },
      workbox: {
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/api\.vk\.com\/.*/i,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'vk-api-cache',
              expiration: {
                maxEntries: 10,
                maxAgeSeconds: 60 * 60 // 1 час
              },
              cacheableResponse: {
                statuses: [0, 200]
              }
            }
          }
        ]
      }
    })
  ]
});
```

### Offline UI индикатор

```typescript
import { useEffect, useState } from 'react';
import { Snackbar } from '@vkontakte/vkui';

const OfflineIndicator = () => {
  const [isOnline, setIsOnline] = useState(navigator.onLine);

  useEffect(() => {
    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  if (isOnline) return null;

  return (
    <Snackbar>
      Нет подключения к интернету
    </Snackbar>
  );
};
```

### iOS PWA ограничения

На iOS есть ограничения:
- Service Workers работают, но с лимитами
- Storage квоты более строгие чем в Chrome
- Background sync не поддерживается
- Push уведомления не работают в PWA (только в Safari)

### Best Practices

✅ **DO:**
- Кэшируйте статические файлы (JS, CSS, images)
- Используйте Network First для API запросов
- Показывайте offline UI индикатор
- Версионируйте кэш (vk-miniapp-v1, v2, ...)
- Удаляйте старые версии кэша при активации нового SW
- Тестируйте offline режим

❌ **DON'T:**
- Не кэшируйте персональные данные
- Не кэшируйте POST/PUT/DELETE запросы
- Не забывайте про ограничения iOS
- Не делайте кэш слишком большим (>50MB может быть очищен)

---

## 8. Реклама (Ads SDK)

### Обзор VK Ads в 2026

VK Ads — это внутренняя рекламная платформа VK Group с точным таргетингом, аналитикой и интеграцией с Mini Apps.

**Рост**: +550% YoY in-app ads revenue (Q3)

### Типы рекламы

1. **Баннерная реклама** — в ленте новостей
2. **Нативная реклама** — внутри Mini App
3. **Rewarded ads** — вознаграждение за просмотр
4. **Interstitial** — полноэкранная реклама
5. **Продвижение Mini Apps** — через VK Ads кабинет

### Интеграция рекламы

#### 1. Настройка в VK Ads кабинете

1. Зарегистрируйтесь на https://ads.vk.com
2. Создайте рекламный кабинет
3. Добавьте Mini App
4. Настройте форматы рекламы
5. Установите бюджет и таргетинг

#### 2. Показ рекламы через VK Bridge

```typescript
import bridge from '@vkontakte/vk-bridge';

const showAd = async () => {
  try {
    const result = await bridge.send('VKWebAppShowNativeAds', {
      ad_format: 'interstitial' // или 'reward'
    });

    if (result.result) {
      console.log('Ad shown successfully');

      // Если rewarded ad, выдаем награду
      if (result.ad_format === 'reward') {
        grantRewardToUser();
      }
    }
  } catch (error) {
    console.error('Failed to show ad:', error);
  }
};
```

#### 3. Rewarded Video

```typescript
const showRewardedAd = async () => {
  try {
    const result = await bridge.send('VKWebAppShowNativeAds', {
      ad_format: 'reward'
    });

    if (result.result) {
      // Пользователь досмотрел рекламу до конца
      const reward = {
        coins: 100,
        type: 'currency'
      };

      // Выдаем награду
      await grantReward(reward);

      // Показываем уведомление
      showSnackbar('Вы получили 100 монет!');
    }
  } catch (error) {
    if (error.error_data?.error_code === 'ADS_NOT_LOADED') {
      showSnackbar('Реклама еще не загрузилась, попробуйте позже');
    }
  }
};
```

### Таргетинг

VK Ads поддерживает детальный таргетинг:

- **Демография**: пол, возраст, семейное положение
- **Геолокация**: страна, город, радиус
- **Интересы**: категории, сообщества, приложения
- **Устройства**: iOS, Android, Web
- **Аудитории**: retargeting VK Pay, подписчиков сообществ
- **События Mini App**: custom events (регистрация, покупка)

### Аналитика рекламы

```typescript
// Трекинг конверсий
bridge.send('VKWebAppConversionHit', {
  pixel_code: 'VK-RTRG-XXXXX-XXXXX',
  event: 'purchase',
  value: 100,
  currency: 'RUB'
});

// Кастомное событие
bridge.send('VKWebAppConversionHit', {
  pixel_code: 'VK-RTRG-XXXXX-XXXXX',
  event: 'level_up',
  value: 10
});
```

### Метрики для отслеживания

- **CTR** (Click-Through Rate) — процент кликов
- **CPA** (Cost Per Action) — цена за действие
- **CPI** (Cost Per Install) — цена за установку
- **Retention** — удержание пользователей
- **ROAS** (Return On Ad Spend) — возврат инвестиций
- **eCPM** (effective Cost Per Mille) — эффективная цена за 1000 показов

### Продвижение через VK Ads

#### Воронка продвижения

```
VK Ads → Landing Page → Mini App Launch → Регистрация → Активация → Покупка
```

#### Оптимизация воронки

1. **Landing Page**
   - Четкий CTA (Call To Action)
   - Быстрая загрузка
   - Превью функционала

2. **Onboarding**
   - Короткий (3-5 экранов)
   - Понятный ценностный посыл
   - Возможность пропустить

3. **Активация**
   - Первое полезное действие
   - Быстрый успех (quick win)
   - Награда за регистрацию

### Тренды 2026

По данным прогнозов на 2026 год:

- **Предиктивные отчеты** — прогноз конверсии по типу аудитории
- **Единые панели аналитики** — расширенные фильтры и прогнозы эффективности
- **Автоматизация рекламы** — оптимизация кампаний на основе бизнес-KPI (лиды, продажи, звонки)
- **AI-оптимизация** — автоматический подбор креативов и таргетинга

### Best Practices

✅ **DO:**
- Используйте rewarded ads для монетизации free users
- Показывайте рекламу в естественных паузах (между уровнями, после действий)
- Настройте frequency capping (не более 3-5 показов в день)
- Тестируйте разные форматы рекламы
- Отслеживайте retention после показа рекламы
- Используйте A/B тесты для креативов

❌ **DON'T:**
- Не показывайте рекламу сразу при запуске приложения
- Не перегружайте пользователя рекламой (снижает retention)
- Не показывайте одну и ту же рекламу слишком часто
- Не блокируйте функционал рекламой (кроме rewarded)
- Не забывайте про UX — реклама не должна раздражать

### Интеграция с MyTracker

```typescript
import { trackEvent } from '@vkontakte/mytracker';

// Трекинг показа рекламы
trackEvent('ad_impression', {
  ad_format: 'interstitial',
  ad_network: 'vk_ads'
});

// Трекинг клика
trackEvent('ad_click', {
  ad_format: 'reward',
  ad_network: 'vk_ads'
});

// Трекинг конверсии после рекламы
trackEvent('ad_conversion', {
  ad_format: 'interstitial',
  ad_network: 'vk_ads',
  conversion_type: 'purchase',
  value: 100
});
```

---

## 9. Требования безопасности

### Критические нарушения безопасности

Даже одно из этих нарушений может привести к блокировке Mini App:

❌ Хранение `access_token` в localStorage
❌ Проверка `sign` на клиенте вместо сервера
❌ Отсутствие HTTPS
❌ Отправка открытых данных без подписи

### 1. HTTPS обязательно

SSL-сертификаты — **обязательное требование**. Мини-приложения должны открываться через HTTPS.

```bash
# Для production используйте Let's Encrypt
certbot certonly --webroot -w /var/www/html -d yourdomain.com
```

### 2. Проверка подписи (sign)

Параметр `sign` — это HMAC-хэш, созданный на стороне VK. **Обязательно** проверяйте его на сервере.

#### Клиент

```typescript
// Получаем параметры запуска
const params = new URLSearchParams(window.location.search);
const launchParams = params.toString();

// Отправляем на бэкенд
const response = await fetch('/api/auth', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ launchParams })
});

const { sessionToken } = await response.json();

// Сохраняем session token (НЕ access_token!)
localStorage.setItem('session_token', sessionToken);
```

#### Сервер (Go)

```go
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "net/url"
    "sort"
    "strings"
)

func VerifyLaunchParams(launchParams string, secretKey string) (map[string]string, error) {
    params, err := url.ParseQuery(launchParams)
    if err != nil {
        return nil, err
    }

    sign := params.Get("sign")
    if sign == "" {
        return nil, errors.New("sign parameter is missing")
    }

    params.Del("sign")

    // Сортируем параметры
    keys := make([]string, 0, len(params))
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    // Формируем строку
    var queryString strings.Builder
    for i, k := range keys {
        if i > 0 {
            queryString.WriteString("&")
        }
        queryString.WriteString(k)
        queryString.WriteString("=")
        queryString.WriteString(params.Get(k))
    }

    // HMAC-SHA256
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(queryString.String()))
    expectedSign := base64.URLEncoding.EncodeToString(h.Sum(nil))

    expectedSign = strings.TrimRight(expectedSign, "=")
    expectedSign = strings.ReplaceAll(expectedSign, "+", "-")
    expectedSign = strings.ReplaceAll(expectedSign, "/", "_")

    if sign != expectedSign {
        return nil, errors.New("invalid signature")
    }

    // Возвращаем проверенные параметры
    result := make(map[string]string)
    for k := range params {
        result[k] = params.Get(k)
    }

    return result, nil
}
```

### 3. Управление токенами

#### Access Token — НИКОГДА не храните в браузере

```typescript
// ❌ НЕПРАВИЛЬНО
const getAccessToken = async () => {
  const data = await bridge.send('VKWebAppGetAuthToken', {
    app_id: APP_ID,
    scope: 'friends'
  });

  // НЕ ДЕЛАЙТЕ ТАК!
  localStorage.setItem('access_token', data.access_token);
};

// ✅ ПРАВИЛЬНО
const getAccessToken = async () => {
  const data = await bridge.send('VKWebAppGetAuthToken', {
    app_id: APP_ID,
    scope: 'friends'
  });

  // Отправляем токен на бэкенд
  await fetch('/api/save-token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${sessionToken}`
    },
    body: JSON.stringify({ access_token: data.access_token })
  });
};
```

#### Session Token — можно хранить в localStorage

```typescript
// Session token можно хранить
const sessionToken = localStorage.getItem('session_token');

// И использовать для API запросов
const response = await fetch('/api/user/profile', {
  headers: {
    'Authorization': `Bearer ${sessionToken}`
  }
});
```

### 4. Безопасное логирование

```typescript
// ❌ ОПАСНО: логирование чувствительных данных
console.log('User data:', user); // может содержать токены

// ✅ БЕЗОПАСНО: редактирование чувствительных полей
const logUser = (user) => {
  console.log('User:', {
    id: user.id,
    name: user.name,
    // Не логируем: access_token, email, phone
  });
};
```

### 5. Cookie Policy

Сервер мини-приложения **не должен** устанавливать cookie-файлы:

```go
// ❌ НЕПРАВИЛЬНО
http.SetCookie(w, &http.Cookie{
    Name:  "session_id",
    Value: sessionID,
})

// ✅ ПРАВИЛЬНО: используйте Authorization header
w.Header().Set("Authorization", "Bearer "+sessionToken)
```

### 6. CORS

```go
// Настройка CORS на бэкенде
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Разрешаем только VK домены
        origin := r.Header.Get("Origin")
        if strings.HasSuffix(origin, ".vk.com") {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        }

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### 7. Rate Limiting

```go
import "golang.org/x/time/rate"

var limiters = make(map[string]*rate.Limiter)

func rateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Header.Get("X-User-ID")

        limiter, exists := limiters[userID]
        if !exists {
            // 10 запросов в секунду на пользователя
            limiter = rate.NewLimiter(10, 20)
            limiters[userID] = limiter
        }

        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### 8. Input Validation

```typescript
// Клиент
import { z } from 'zod';

const userSchema = z.object({
  name: z.string().min(2).max(50),
  email: z.string().email(),
  age: z.number().min(0).max(150)
});

const validateUser = (data: unknown) => {
  try {
    return userSchema.parse(data);
  } catch (error) {
    throw new Error('Invalid user data');
  }
};
```

```go
// Сервер
import "github.com/go-playground/validator/v10"

type User struct {
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"gte=0,lte=150"`
}

func validateUser(user User) error {
    validate := validator.New()
    return validate.Struct(user)
}
```

### 9. XSS Protection

```typescript
// Используйте React (автоматическое экранирование)
const UserName = ({ name }) => <div>{name}</div>; // ✅ безопасно

// Если нужен HTML, используйте DOMPurify
import DOMPurify from 'dompurify';

const SafeHTML = ({ html }) => {
  const clean = DOMPurify.sanitize(html);
  return <div dangerouslySetInnerHTML={{ __html: clean }} />;
};
```

### 10. CSRF Protection

```go
import "github.com/gorilla/csrf"

func main() {
    CSRF := csrf.Protect(
        []byte("32-byte-long-secret-key"),
        csrf.Secure(true),
        csrf.SameSite(csrf.SameSiteStrictMode),
    )

    http.Handle("/", CSRF(handler))
}
```

### Security Checklist

- [ ] HTTPS включен
- [ ] Проверка `sign` на бэкенде
- [ ] `access_token` не хранится в браузере
- [ ] Используется session token для API
- [ ] Cookie не устанавливаются
- [ ] CORS настроен правильно
- [ ] Rate limiting включен
- [ ] Input validation на клиенте и сервере
- [ ] XSS protection (React + DOMPurify)
- [ ] CSRF protection
- [ ] Чувствительные данные не логируются
- [ ] Secrets не хранятся в коде (используйте .env)

### Безопасное хранение секретов

```bash
# .env (НЕ коммитьте в git!)
VK_APP_SECRET=your_secret_key
VK_APP_ID=12345678
DATABASE_URL=postgresql://...
```

```go
// Загрузка из .env
import "github.com/joho/godotenv"

func init() {
    godotenv.Load()
}

func main() {
    secretKey := os.Getenv("VK_APP_SECRET")
    if secretKey == "" {
        log.Fatal("VK_APP_SECRET is not set")
    }
}
```

---

## 10. Процесс модерации

### Общая информация

Модерация — это встроенный механизм качества, который защищает платформу от хаоса, а пользователей от ненадежных или небрежных приложений.

**Срок модерации**: до 7 дней
**Релизы**: еженедельно по четвергам

### Что проверяет модерация

1. **Функциональность** — приложение работает без критических ошибок
2. **UX** — понятный интерфейс, объяснение действий
3. **Безопасность** — проверка sign, HTTPS, отсутствие уязвимостей
4. **Доступы** — корректный запрос прав доступа
5. **Платежи** — прозрачность цен, нет скрытых комиссий
6. **Контент** — соответствие правилам VK

### Основные причины отказа

1. **Непредсказуемое поведение при ошибках**
   - Приложение падает при потере сети
   - Белый экран при ошибке API
   - Нет обработки edge cases

2. **Скрытие смысла действий**
   - Непонятные кнопки без описания
   - Неочевидная навигация
   - Отсутствие подтверждений для критических действий

3. **Избыточный запрос доступов**
   - Запрос friends при отсутствии социальных функций
   - Запрос photos без функционала загрузки
   - Необоснованный запрос прав

4. **Неясное объяснение обработки данных**
   - Отсутствие политики конфиденциальности
   - Нет объяснения зачем нужны данные
   - Непрозрачная обработка платежей

5. **Агрессивная воронка**
   - Навязывание платежей
   - Блокировка функционала без оплаты
   - Spam уведомлений

### Подготовка к модерации

#### 1. Проверьте функциональность

```typescript
// Обработка ошибок
const fetchData = async () => {
  try {
    const response = await fetch('/api/data');
    const data = await response.json();
    return data;
  } catch (error) {
    // Показываем понятное сообщение
    showSnackbar('Не удалось загрузить данные. Попробуйте позже.');

    // Логируем ошибку для анализа
    console.error('Failed to fetch data:', error);

    // Не показываем пустой экран
    return getFallbackData();
  }
};
```

#### 2. Обработка offline состояния

```typescript
import { useEffect, useState } from 'react';

const useNetworkStatus = () => {
  const [isOnline, setIsOnline] = useState(navigator.onLine);

  useEffect(() => {
    const handleOnline = () => {
      setIsOnline(true);
      showSnackbar('Соединение восстановлено');
    };

    const handleOffline = () => {
      setIsOnline(false);
      showSnackbar('Нет подключения к интернету');
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  return isOnline;
};
```

#### 3. Понятные тексты и подтверждения

```typescript
import { Alert } from '@vkontakte/vkui';

const DeleteButton = ({ onDelete }) => {
  const [showAlert, setShowAlert] = useState(false);

  return (
    <>
      <Button onClick={() => setShowAlert(true)}>
        Удалить аккаунт
      </Button>

      {showAlert && (
        <Alert
          actions={[
            {
              title: 'Отмена',
              mode: 'cancel',
              autoClose: true
            },
            {
              title: 'Удалить',
              mode: 'destructive',
              autoClose: true,
              action: onDelete
            }
          ]}
          onClose={() => setShowAlert(false)}
        >
          <h2>Удалить аккаунт?</h2>
          <p>
            Все ваши данные будут удалены без возможности восстановления.
            Это действие необратимо.
          </p>
        </Alert>
      )}
    </>
  );
};
```

#### 4. Запрос только нужных прав

```typescript
// ❌ НЕПРАВИЛЬНО: запрос всех прав
const requestPermissions = async () => {
  await bridge.send('VKWebAppGetAuthToken', {
    app_id: APP_ID,
    scope: 'friends,photos,video,audio,pages,groups,wall'
  });
};

// ✅ ПРАВИЛЬНО: только необходимые
const requestPermissions = async () => {
  await bridge.send('VKWebAppGetAuthToken', {
    app_id: APP_ID,
    scope: 'friends' // только если реально нужно
  });
};
```

#### 5. Политика конфиденциальности

Обязательно добавьте страницу с политикой:

```typescript
const PrivacyPolicy = () => {
  return (
    <Panel id="privacy">
      <PanelHeader>Политика конфиденциальности</PanelHeader>

      <Group>
        <Div>
          <Title level="3">Какие данные мы собираем</Title>
          <Text>
            • Публичная информация из профиля VK (имя, фото)
            • История ваших действий в приложении
            • Техническая информация (тип устройства, браузер)
          </Text>
        </Div>

        <Div>
          <Title level="3">Как мы используем данные</Title>
          <Text>
            • Для предоставления функционала приложения
            • Для улучшения качества сервиса
            • Для отправки уведомлений (с вашего согласия)
          </Text>
        </Div>

        <Div>
          <Title level="3">Как мы защищаем данные</Title>
          <Text>
            • Используем HTTPS для всех соединений
            • Не храним токены доступа в браузере
            • Не передаем данные третьим лицам
          </Text>
        </Div>

        <Div>
          <Title level="3">Ваши права</Title>
          <Text>
            • Вы можете запросить удаление данных
            • Вы можете отозвать права доступа
            • Вы можете экспортировать свои данные
          </Text>
        </Div>
      </Group>
    </Panel>
  );
};
```

### Чеклист перед подачей на модерацию

#### Функциональность
- [ ] Приложение работает без ошибок
- [ ] Все кнопки и формы функциональны
- [ ] Обработаны все сценарии ошибок
- [ ] Нет белых экранов и зависаний
- [ ] Приложение работает на всех устройствах (iOS, Android, Web)

#### Безопасность
- [ ] HTTPS включен
- [ ] Проверка sign на бэкенде
- [ ] access_token не хранится в браузере
- [ ] Cookie не устанавливаются
- [ ] Нет критических уязвимостей

#### UX
- [ ] Понятная навигация
- [ ] Все кнопки подписаны
- [ ] Подтверждения для критических действий
- [ ] Обратная связь на каждое действие (Snackbar, Taptic)
- [ ] Offline состояние обработано
- [ ] Загрузки показывают Spinner

#### Контент
- [ ] Политика конфиденциальности добавлена
- [ ] Пользовательское соглашение добавлено
- [ ] Описание приложения понятно
- [ ] Иконка и обложка соответствуют функционалу
- [ ] Нет запрещенного контента (18+, насилие и т.д.)

#### Монетизация
- [ ] Цены указаны четко
- [ ] Нет скрытых комиссий
- [ ] Возможность отмены подписки
- [ ] Прозрачные условия возврата

#### Доступы
- [ ] Запрашиваются только необходимые права
- [ ] Объяснено зачем нужны права
- [ ] Приложение работает без прав (если возможно)

#### Производительность
- [ ] Время загрузки < 3 секунд
- [ ] Bundle size < 10 MB
- [ ] Нет утечек памяти
- [ ] Плавные анимации (60 FPS)

### Процесс подачи

1. **Подготовка**
   - Заполните все поля в настройках приложения
   - Загрузите иконку (512x512 px)
   - Загрузите обложку (1590x400 px)
   - Напишите описание (до 140 символов)

2. **Тестирование**
   - Протестируйте на всех платформах
   - Попросите коллег протестировать
   - Проверьте все сценарии использования

3. **Отправка**
   - Нажмите "Отправить на модерацию"
   - Дождитесь результата (до 7 дней)

4. **Реакция на отказ**
   - Внимательно прочитайте причину
   - Исправьте указанные проблемы
   - Повторно отправьте на модерацию

### Советы для успешной модерации

✅ **DO:**
- Объясняйте каждое действие пользователю
- Показывайте прогресс загрузки
- Добавьте onboarding для новых пользователей
- Используйте стандартные компоненты VKUI
- Тестируйте на реальных устройствах

❌ **DON'T:**
- Не скрывайте ошибки
- Не запрашивайте лишние права
- Не навязывайте платежи
- Не используйте агрессивные уведомления
- Не игнорируйте комментарии модераторов

---

## Рекомендации для нашего проекта

### Технологический стек

```
Frontend:
- React 18
- TypeScript
- Vite (для быстрой сборки)
- VKUI v7+ (UI компоненты)
- @vkontakte/vk-bridge (интеграция с VK)
- @vkontakte/vk-mini-apps-router (навигация)
- React Query (кэширование API)
- Zustand (state management)

Backend:
- Go 1.21+
- Echo/Gin (HTTP framework)
- PostgreSQL (база данных)
- Redis (кэширование, сессии)
- Docker + Docker Compose

Инфраструктура:
- VK Hosting (статика)
- Yandex Object Storage (CDN)
- GitHub Actions (CI/CD)
```

### Архитектура

```
/frontend
  /src
    /app           # Главный App компонент
    /pages         # Страницы (Home, Profile, Settings)
    /components    # Переиспользуемые компоненты
    /hooks         # Кастомные хуки
    /services      # API клиенты
    /utils         # Утилиты
    /types         # TypeScript типы
  /public
    /icons         # Иконки для PWA
    manifest.json
    sw.js
  vite.config.ts
  tsconfig.json

/backend
  /cmd
    /api         # Точка входа
  /internal
    /handler     # HTTP handlers
    /service     # Бизнес-логика
    /repository  # Database access
    /middleware  # Middlewares
  /pkg
    /vkbridge    # VK Bridge utils
    /validator   # Validation
```

### CI/CD Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy to VK Hosting

on:
  push:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install dependencies
        run: npm ci

      - name: Build
        run: npm run build
        env:
          VITE_VK_APP_ID: ${{ secrets.VK_APP_ID }}

      - name: Deploy to VK Hosting
        run: npx @vkontakte/vk-miniapps-deploy
        env:
          VK_ADMIN_ACCESS_TOKEN: ${{ secrets.VK_ADMIN_ACCESS_TOKEN }}
```

### Monitoring & Logging

```typescript
// Интеграция Sentry для отслеживания ошибок
import * as Sentry from '@sentry/react';

Sentry.init({
  dsn: 'YOUR_SENTRY_DSN',
  environment: import.meta.env.MODE,
  integrations: [
    new Sentry.BrowserTracing(),
    new Sentry.Replay()
  ],
  tracesSampleRate: 1.0,
  replaysSessionSampleRate: 0.1,
  replaysOnErrorSampleRate: 1.0
});
```

---

## Источники

### Официальная документация VK

- [VK Bridge GitHub](https://github.com/VKCOM/vk-bridge)
- [VK Bridge NPM](https://www.npmjs.com/package/@vkontakte/vk-bridge)
- [VKUI GitHub](https://github.com/VKCOM/VKUI)
- [VKUI NPM](https://www.npmjs.com/package/@vkontakte/vkui)
- [VKUI Styleguide](https://vkcom.github.io/vkui-styleguide/)
- [VK Mini Apps Deploy](https://github.com/VKCOM/vk-miniapps-deploy)
- [VK Mini Apps Router](https://github.com/VKCOM/vk-mini-apps-router)
- [VK Apps Launch Params](https://github.com/VKCOM/vk-apps-launch-params)

### Статьи и гайды

- [Дизайн-гайд VK Mini Apps: технический подход к UX](https://spark.ru/user/261378/blog/276937/dizajn-gajd-vk-mini-apps-tehnicheskij-podhod-k-ux)
- [Сценарии авторизации в VK Mini Apps: VK ID, токены и best-practice безопасности](https://spark.ru/user/261378/blog/277666/stsenarii-avtorizatsii-v-vk-mini-apps-vk-id-tokeni-i-best-practice-bezopasnosti)
- [Продвижение VK Mini Apps через VK Ads: как строить воронку внутри экосистемы](https://spark.ru/user/261378/blog/277466/prodvizhenie-vk-mini-apps-cherez-vk-ads-kak-stroit-voronku-vnutri-ekosistemi)
- [CI/CD для VK Mini Apps: инженерный подход к стабильной доставке](https://spark.ru/user/261378/blog/280321/ci-cd-dlya-vk-mini-apps-inzhenernij-podhod-k-stabilnoj-dostavke)
- [Модерация VK Mini Apps: ошибки новичков при публикации](https://onskills.ru/blog/moderaciya-vk-mini-apps-oshibki-publikacii/)
- [VK Mini Apps: как создать мини-приложение ВКонтакте с нуля](https://timeweb.cloud/tutorials/react/kak-sozdat-mini-prilozhenie-vk-mini-apps)
- [Как создать веб-приложение на базе VK Mini Apps](https://selectel.ru/blog/tutorials/vk-mini-apps/)

### Аналитика и трекинг

- [MyTracker - VK Mini Apps Integration](https://docs.tracker.my.com/en/tracking/platforms/vk-mini-apps)
- [VK Pay](https://tadviser.com/index.php/Product:VK_Pay)
- [VK Statistics 2026](https://bayelsawatch.com/vk-statistics/)

### Performance и оптимизация

- [Complete Guide to JavaScript Performance Optimization (2026)](https://needlecode.gitlab.io/blog/javascript/complete-guide-to-javascript-performance.html)
- [Vite vs. Webpack: A Head-to-Head Comparison](https://kinsta.com/blog/vite-vs-webpack/)
- [Code Splitting | webpack](https://webpack.js.org/guides/code-splitting/)
- [Progressive Web Apps 2026: PWA Performance Guide](https://www.digitalapplied.com/blog/progressive-web-apps-2026-pwa-performance-guide)

### PWA и Service Workers

- [Making the PWA work offline with service workers - MDN](https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/Tutorials/js13kGames/Offline_Service_workers)
- [PWA iOS Limitations and Safari Support [2026]](https://www.magicbell.com/blog/pwa-ios-limitations-safari-support-complete-guide)
- [How to Implement Service Workers for Offline Support in React](https://oneuptime.com/blog/post/2026-01-15-service-workers-offline-support-react/view)

### Безопасность

- [Сценарии авторизации в VK Mini Apps: best-practice безопасности](https://spark.ru/user/261378/blog/277666/stsenarii-avtorizatsii-v-vk-mini-apps-vk-id-tokeni-i-best-practice-bezopasnosti)
- [VK WorkSpace Security](https://workspace.vk.ru/security/)

### Монетизация и реклама

- [VK / VK ecosystem development highlights](https://vk.company/ru/press/releases/11047/)
- [Telegram mini apps 2026 monetization guide](https://merge.rocks/blog/telegram-mini-apps-2026-monetization-guide-how-to-earn-from-telegram-mini-apps)
- [Таргетированная реклама ВКонтакте 2026: тренды, ИИ и будущее VK Ads](https://barsukov.by/blog/budushhee-targetirovannoj-reklamy-v-vkontakte-trendy-i-prognozy-na-2026-god/)

---

**Последнее обновление**: 29 марта 2026

**Для вопросов и предложений**: создайте issue в репозитории проекта
