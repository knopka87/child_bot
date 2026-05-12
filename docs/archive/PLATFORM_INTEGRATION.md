# Platform Integration Guide

## Overview

Приложение поддерживает несколько платформ:
- **VK Mini Apps** (vk.com) - VK Bridge API
- **Max** (max.ru) - MAX Bridge API (отдельный мессенджер Mail.ru Group)
- **Telegram Mini Apps** - Telegram WebApp API
- **Web** (fallback) - стандартный веб-браузер

## Architecture

### Platform Bridge

`PlatformBridge` - центральный класс для работы с платформами. Автоматически определяет текущую платформу и создаёт соответствующий адаптер.

```typescript
const bridge = new PlatformBridge();
await bridge.init();

const info = bridge.getInfo();
const user = await bridge.getUser();
const theme = bridge.getTheme();
```

### Platform Adapters

Каждая платформа имеет свой адаптер, реализующий интерфейс `IPlatformAdapter`:

- `WebAdapter` - базовая реализация для веб-версии
- `VKAdapter` - интеграция с VK Bridge API (TODO)
- `MaxAdapter` - интеграция с MAX Bridge API (TODO) - https://dev.max.ru/docs/webapps/bridge
- `TelegramAdapter` - интеграция с Telegram WebApp API (TODO)

## Platform Detection

Платформа определяется автоматически при инициализации:

1. Проверка VK Bridge (`window.vkBridge` или URL параметр `vk_platform`)
2. Проверка MAX Bridge (`window.MaxBridge` или referrer содержит 'max.ru')
3. Проверка Telegram WebApp (`window.Telegram.WebApp`)
4. Fallback на Web

## Features Support Matrix

| Feature | VK | Max | Telegram | Web |
|---------|-----|-----|----------|-----|
| Native Share | ✅ | ✅ | ✅ | ⚠️ (Web Share API) |
| Haptic Feedback | ✅ | ✅ | ✅ | ❌ |
| Native Payment | ✅ | ✅ | ✅ | ❌ |
| Cloud Storage | ✅ | ✅ | ✅ | ❌ |
| Push Notifications | ✅ | ✅ | ✅ | ❌ |
| Camera Access | ❌ | ✅ | ✅ | ⚠️ (getUserMedia) |
| Biometric Auth | ❌ | ❌ | ✅ | ❌ |

## Usage Examples

### Sharing Content

```typescript
await bridge.share({
  title: 'Check this out!',
  text: 'Amazing content',
  url: 'https://example.com'
});
```

### Haptic Feedback

```typescript
bridge.hapticFeedback({
  type: 'impact',
  style: 'medium'
});
```

### Opening Links

```typescript
bridge.openLink('https://example.com');
```

### Getting User Info

```typescript
const user = await bridge.getUser();
console.log(user.firstName, user.lastName);
```

### Theme Detection

```typescript
const theme = bridge.getTheme();
console.log(theme.colorScheme); // 'light' или 'dark'
```

## Implementation Status

- ✅ Platform type definitions
- ✅ PlatformBridge core
- ✅ WebAdapter implementation
- ⏳ VKAdapter (uses existing VK Bridge integration)
- ⏳ MaxAdapter (requires MAX Bridge SDK setup)
- ⏳ TelegramAdapter (requires Telegram Bot API setup)

## Platform-Specific Documentation

### VK Bridge
- Docs: https://dev.vk.com/bridge/getting-started
- SDK: `@vkontakte/vk-bridge`
- Methods: VKWebAppInit, VKWebAppShare, VKWebAppGetUserInfo

### MAX Bridge
- Docs: https://dev.max.ru/docs/webapps/bridge
- SDK: MAX Bridge (встроен в WebApp)
- Methods: MaxBridge.init(), MaxBridge.getUserInfo(), MaxBridge.share()

### Telegram WebApp
- Docs: https://core.telegram.org/bots/webapps
- SDK: `window.Telegram.WebApp`
- Methods: initData, expand(), close(), MainButton

## Next Steps

Для полной интеграции со всеми платформами необходимо:

1. Реализовать `VKAdapter` с полной поддержкой VK Bridge API
2. Реализовать `MaxAdapter` с полной поддержкой MAX Bridge API (https://dev.max.ru/docs)
3. Реализовать `TelegramAdapter` с полной поддержкой Telegram WebApp API
4. Добавить platform-specific UI компоненты
5. Настроить feature flags для каждой платформы
6. Провести тестирование на реальных устройствах (VK клиент, Max app, Telegram)
