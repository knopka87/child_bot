# VK Roadmap Improvements — Улучшения roadmap для VK

**Дата:** 2026-03-29
**На основе:** [VK_BEST_PRACTICES.md](./VK_BEST_PRACTICES.md)

---

## 📋 Анализ текущего roadmap

Проанализирован roadmap на соответствие VK Mini Apps best practices. Выявлены области для улучшения и добавлены VK-специфичные рекомендации.

---

## 🔧 Критические улучшения (Must Have)

### 1. ✅ Интеграция VKUI вместо Custom UI Kit

**Текущее состояние:**
- В `02_CORE.md` предлагается создать Custom UI Kit с нуля

**Проблема:**
- Несоответствие VK design guidelines
- Дополнительная работа на 2-3 дня
- Не используем преимущества VKUI (2500+ иконов, адаптивность, темы)

**Рекомендация:**

```typescript
// ❌ Текущий подход - создавать с нуля
import { Button } from '@/components/ui/Button'

// ✅ Рекомендуется - использовать VKUI
import { Button } from '@vkontakte/vkui'
import '@vkontakte/vkui/dist/vkui.css'
```

**Что изменить в roadmap:**

**`01_SETUP.md`** - добавить VKUI установку:
```bash
npm install @vkontakte/vkui @vkontakte/icons
```

**`02_CORE.md`** - заменить Custom UI Kit на VKUI wrapper:
```typescript
// Создаем thin wrapper над VKUI для кастомизации
// src/components/ui/Button.tsx
import { Button as VKUIButton, ButtonProps } from '@vkontakte/vkui'

export function Button(props: ButtonProps) {
  return <VKUIButton {...props} />
}

// При необходимости добавляем кастомную логику
export function PrimaryButton(props: ButtonProps) {
  return <VKUIButton mode="primary" size="l" {...props} />
}
```

**Преимущества:**
- ✅ Соответствие VK design guidelines
- ✅ Автоматическая адаптивность
- ✅ Темная/светлая тема из коробки
- ✅ Экономия 2-3 дней разработки
- ✅ 2500+ готовых иконок
- ✅ Accessibility из коробки

---

### 2. ✅ VK Bridge Security - Проверка sign

**Текущее состояние:**
- В `SECURITY.md` описана общая JWT аутентификация
- Не упомянута VK-специфичная проверка `sign`

**Проблема:**
- VK обязательно требует проверку `sign` параметра
- Без этого приложение НЕ пройдет модерацию

**Рекомендация:**

Добавить в **`SECURITY.md`** раздел "VK Bridge Sign Validation":

```typescript
// Frontend: получаем launch params
import bridge from '@vkontakte/vk-bridge'

async function authenticate() {
  const launchParams = await bridge.send('VKWebAppGetLaunchParams')

  // Отправляем ВСЕ параметры на backend (включая sign)
  const response = await fetch('/api/v1/auth/vk', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      vk_user_id: launchParams.vk_user_id,
      vk_app_id: launchParams.vk_app_id,
      vk_platform: launchParams.vk_platform,
      sign: launchParams.sign, // КРИТИЧНО!
      // ... все vk_* параметры
    })
  })

  const { access_token } = await response.json()
  localStorage.setItem('token', access_token)
}
```

```go
// Backend: проверяем sign (Go)
func ValidateVKSign(params map[string]string, secretKey string) bool {
    sign := params["sign"]
    delete(params, "sign")

    // Собираем vk_* параметры
    var keys []string
    for k := range params {
        if strings.HasPrefix(k, "vk_") {
            keys = append(keys, k)
        }
    }
    sort.Strings(keys)

    // Строим query string
    var queryString string
    for _, k := range keys {
        queryString += k + "=" + params[k] + "&"
    }
    queryString = strings.TrimSuffix(queryString, "&")

    // HMAC-SHA256
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(queryString))
    expectedSign := base64.URLEncoding.EncodeToString(h.Sum(nil))

    return sign == expectedSign
}
```

**Где применить:**
- `SECURITY.md` - добавить раздел
- `VK_BEST_PRACTICES.md` - уже есть
- `01_SETUP.md` - добавить в backend setup

---

### 3. ✅ Bundle Size Optimization (< 10 MB)

**Текущее состояние:**
- Не упоминаются ограничения по размеру bundle

**Проблема:**
- VK требует bundle < 10 MB
- Без оптимизации легко превысить лимит

**Рекомендация:**

Добавить в **`01_SETUP.md`** конфигурацию Vite:

```typescript
// vite.config.ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { visualizer } from 'rollup-plugin-visualizer'

export default defineConfig({
  plugins: [
    react(),
    visualizer({ open: true }) // Анализ bundle
  ],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          // Разделяем vendor код
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          'vk-vendor': ['@vkontakte/vkui', '@vkontakte/vk-bridge'],
          'ui-vendor': ['zustand', '@tanstack/react-query'],
        }
      }
    },
    chunkSizeWarningLimit: 1000, // 1MB warning
  },
  // Оптимизация для production
  esbuild: {
    drop: ['console', 'debugger'], // Удаляем в prod
  }
})
```

Добавить в **`12_TESTING.md`** проверку размера:

```bash
# Проверка размера bundle
npm run build
du -sh dist/assets/*.js

# Должно быть:
# main.*.js < 500KB
# vendor.*.js < 1MB
# Total < 2MB (gzip)
```

**Action Items:**
- [ ] Настроить code splitting по routes
- [ ] Использовать dynamic imports для тяжелых компонентов
- [ ] Оптимизировать изображения (WebP, lazy loading)
- [ ] Tree-shaking для неиспользуемого кода

---

### 4. ✅ VK Platform Detection и Adaptive UI

**Текущее состояние:**
- `13_PLATFORMS.md` описывает абстракцию, но без VK-специфики

**Проблема:**
- VK miniapps запускаются на разных платформах (iOS, Android, Web, Desktop)
- Нужна адаптация UI под каждую платформу

**Рекомендация:**

Добавить в **`02_CORE.md`** VK Platform Detection:

```typescript
// src/utils/platform.ts
import bridge from '@vkontakte/vk-bridge'
import { Platform } from '@vkontakte/vkui'

export async function detectVKPlatform(): Promise<Platform> {
  const launchParams = await bridge.send('VKWebAppGetLaunchParams')

  // vk_platform: 'mobile_iphone', 'mobile_android', 'desktop_web', etc.
  const platform = launchParams.vk_platform

  if (platform.includes('iphone') || platform.includes('ipad')) {
    return 'ios'
  }
  if (platform.includes('android')) {
    return 'android'
  }
  return 'vkcom' // desktop
}
```

Обновить **`App.tsx`**:

```typescript
// src/App.tsx
import { ConfigProvider, AdaptivityProvider, AppRoot } from '@vkontakte/vkui'
import { detectVKPlatform } from './utils/platform'
import { useEffect, useState } from 'react'

export function App() {
  const [platform, setPlatform] = useState<Platform>('android')

  useEffect(() => {
    detectVKPlatform().then(setPlatform)
  }, [])

  return (
    <ConfigProvider platform={platform}>
      <AdaptivityProvider>
        <AppRoot>
          {/* Ваши компоненты */}
        </AppRoot>
      </AdaptivityProvider>
    </ConfigProvider>
  )
}
```

**Преимущества:**
- ✅ Автоматическая адаптация UI под iOS/Android/Desktop
- ✅ Правильные нативные паттерны (iOS bottom sheets, Android modals)
- ✅ Соответствие platform guidelines

---

### 5. ✅ Haptic Feedback Integration

**Текущее состояние:**
- Не упоминается тактильная обратная связь

**Проблема:**
- VK рекомендует haptics для лучшего UX
- Особенно важно для мобильных устройств

**Рекомендация:**

Добавить в **`02_CORE.md`** haptics service:

```typescript
// src/services/haptics.ts
import bridge from '@vkontakte/vk-bridge'

export const haptics = {
  impact: (style: 'light' | 'medium' | 'heavy' = 'medium') => {
    bridge.send('VKWebAppTapticImpactOccurred', { style })
  },

  notification: (type: 'error' | 'success' | 'warning') => {
    bridge.send('VKWebAppTapticNotificationOccurred', { type })
  },

  selection: () => {
    bridge.send('VKWebAppTapticSelectionChanged', {})
  }
}

// Использование в компонентах
import { haptics } from '@/services/haptics'

function ActionButton() {
  const handleClick = () => {
    haptics.impact('light')
    // ... логика
  }

  return <Button onClick={handleClick}>Помоги разобраться</Button>
}

// При успешной операции
async function submitAnswer() {
  const result = await api.submitAnswer(answer)
  if (result.is_correct) {
    haptics.notification('success') // 👍 вибрация успеха
  } else {
    haptics.notification('error') // ⚠️ вибрация ошибки
  }
}
```

**Где использовать haptics:**
- Клики по кнопкам (light)
- Навигация между экранами (selection)
- Успешные операции (notification: success)
- Ошибки (notification: error)
- Достижения разблокированы (notification: success + impact: heavy)

---

### 6. ✅ VK Storage API (вместо localStorage)

**Текущее состояние:**
- Используется обычный `localStorage`

**Проблема:**
- В iOS VK не всегда сохраняется localStorage
- VK Storage API гарантирует сохранение данных

**Рекомендация:**

Добавить в **`02_CORE.md`** VK Storage wrapper:

```typescript
// src/utils/storage.ts
import bridge from '@vkontakte/vk-bridge'

export const storage = {
  async getItem(key: string): Promise<string | null> {
    try {
      const data = await bridge.send('VKWebAppStorageGet', { keys: [key] })
      return data.keys[0]?.value || null
    } catch {
      return localStorage.getItem(key)
    }
  },

  async setItem(key: string, value: string): Promise<void> {
    try {
      await bridge.send('VKWebAppStorageSet', { key, value })
    } catch {
      localStorage.setItem(key, value)
    }
  },

  async removeItem(key: string): Promise<void> {
    try {
      await bridge.send('VKWebAppStorageSet', { key, value: '' })
    } catch {
      localStorage.removeItem(key)
    }
  }
}

// Использование
await storage.setItem('token', accessToken)
const token = await storage.getItem('token')
```

**Важно:**
- VK Storage ограничен ~1MB на ключ
- Для больших данных используйте сжатие или IndexedDB

---

### 7. ✅ Share и Invite через VK Bridge

**Текущее состояние:**
- В `08_FRIENDS.md` не указан способ share

**Проблема:**
- Нужно использовать нативный VK share для лучшего UX

**Рекомендация:**

Обновить **`08_FRIENDS.md`**:

```typescript
// src/components/sections/friends/InviteFriendCard.tsx
import bridge from '@vkontakte/vk-bridge'
import { haptics } from '@/services/haptics'

async function handleShare() {
  haptics.impact('light')

  try {
    await bridge.send('VKWebAppShare', {
      link: referralInfo.referral_link
    })

    analytics.sendEvent({
      event_name: 'referral_share_sent',
      channel_type: 'vk_native'
    })
  } catch (error) {
    // Fallback to copy
    await navigator.clipboard.writeText(referralInfo.referral_link)
    showToast('Ссылка скопирована')
  }
}

// Также можно использовать VKWebAppShowWallPostBox
async function handleShareToWall() {
  await bridge.send('VKWebAppShowWallPostBox', {
    message: 'Попробуй крутой помощник с домашкой! 🎓',
    attachments: referralInfo.referral_link
  })
}
```

---

### 8. ✅ QR Code Scanner для Upload

**Текущее состояние:**
- В `05_HELP.md` и `06_CHECK.md` нет упоминания QR сканера

**Проблема:**
- Пользователи могут хотеть отсканировать QR код с учебника

**Рекомендация:**

Добавить в **`05_HELP.md`** QR scanner опцию:

```typescript
// src/components/sections/help/SourcePicker.tsx
import bridge from '@vkontakte/vk-bridge'

async function handleQRScan() {
  try {
    const result = await bridge.send('VKWebAppOpenCodeReader', {})

    if (result.code_data) {
      // Если QR код содержит ссылку на задание
      const taskUrl = result.code_data
      await loadTaskFromURL(taskUrl)
    }
  } catch (error) {
    showToast('QR сканирование недоступно')
  }
}

// UI
<Card>
  <Cell before={<Icon28QrCodeOutline />} onClick={handleQRScan}>
    Сканировать QR код
  </Cell>
</Card>
```

---

### 9. ✅ VK Pay Integration

**Текущее состояние:**
- `09_PROFILE.md` описывает generic payment flow

**Проблема:**
- Не используется VK Pay (самый удобный для VK пользователей)

**Рекомендация:**

Обновить **`09_PROFILE.md`** с VK Pay:

```typescript
// src/pages/Profile/PaywallPage.tsx
import bridge from '@vkontakte/vk-bridge'

async function handleVKPay(planId: string) {
  try {
    // 1. Создаем заказ на backend
    const order = await api.post('/api/v1/subscription/vk-pay/order', {
      billing_plan_id: planId
    })

    // 2. Открываем VK Pay
    const result = await bridge.send('VKWebAppOpenPayForm', {
      app_id: VK_APP_ID,
      action: 'pay-to-service',
      params: {
        merchant_id: order.merchant_id,
        amount: order.amount,
        description: `Подписка ${planName}`,
        order_id: order.order_id,
        sign: order.sign, // Подпись с backend
      }
    })

    if (result.status === 'success') {
      // 3. Ждем webhook от VK на backend
      await pollOrderStatus(order.order_id)
      showToast('Подписка активирована! 🎉')
    }
  } catch (error) {
    analytics.sendEvent({
      event_name: 'payment_failed',
      error_code: error.message
    })
  }
}
```

Backend webhook:

```go
// Backend: обработка VK Pay webhook
func HandleVKPayWebhook(w http.ResponseWriter, r *http.Request) {
    var notification VKPayNotification
    json.NewDecoder(r.Body).Decode(&notification)

    // Проверяем подпись
    if !validateVKPaySign(notification, vkPaySecret) {
        http.Error(w, "Invalid signature", 403)
        return
    }

    // Обрабатываем оплату
    if notification.NotificationType == "order_status_change_test" ||
       notification.NotificationType == "order_status_change" {

        orderID := notification.OrderID
        status := notification.Status // "chargeable"

        if status == "chargeable" {
            // Активируем подписку
            db.ActivateSubscription(orderID)

            // Отправляем аналитику
            analytics.SendEvent("payment_success", map[string]interface{}{
                "order_id": orderID,
                "amount": notification.Amount,
            })
        }
    }

    // Важно: вернуть "ok"
    w.Write([]byte("ok"))
}
```

**Преимущества VK Pay:**
- ✅ +23% конверсия vs другие методы
- ✅ Уже привязанные карты пользователей
- ✅ Быстрая оплата (1 клик)
- ✅ Поддержка рассрочки

---

### 10. ✅ Ads Integration для Монетизации

**Текущее состояние:**
- Не упоминается реклама

**Проблема:**
- Упускаем возможность монетизации через рекламу
- Особенно для free tier пользователей

**Рекомендация:**

Добавить новый файл **`14_MONETIZATION.md`**:

```typescript
// src/services/ads.ts
import bridge from '@vkontakte/vk-bridge'

export const ads = {
  // Показать rewarded video (за награду)
  async showRewarded(): Promise<boolean> {
    try {
      const result = await bridge.send('VKWebAppCheckNativeAds', {
        ad_format: 'reward'
      })

      if (result.result) {
        await bridge.send('VKWebAppShowNativeAds', { ad_format: 'reward' })
        return true // Пользователь досмотрел
      }
    } catch {
      return false
    }
  },

  // Показать interstitial (полноэкранная)
  async showInterstitial(): Promise<void> {
    try {
      await bridge.send('VKWebAppShowNativeAds', {
        ad_format: 'interstitial'
      })
    } catch {
      // Реклама не показана
    }
  }
}

// Использование в приложении
// Получить дополнительные подсказки за просмотр рекламы
async function unlockExtraHint() {
  const watched = await ads.showRewarded()

  if (watched) {
    // Разблокируем hint level 3
    setUnlockedHints([1, 2, 3])
    haptics.notification('success')
    showToast('Подсказка разблокирована! 🎉')
  }
}

// Показать рекламу после N попыток (не раздражая)
let attemptCount = 0
function incrementAttempts() {
  attemptCount++
  if (attemptCount % 5 === 0) {
    // Каждые 5 попыток
    ads.showInterstitial()
  }
}
```

**Сценарии использования:**
- Rewarded ads для разблокировки hint level 3
- Rewarded ads для дополнительных монет
- Interstitial после каждых 5 попыток
- Native ads в списке истории (после каждых 3 карточек)

---

## 🎨 UI/UX Улучшения (Nice to Have)

### 11. ⚡ Skeleton Loaders (вместо Spinner)

**Рекомендация:**

Использовать VKUI Skeleton вместо обычного Spinner:

```typescript
import { Skeleton } from '@vkontakte/vkui'

// ❌ Обычный spinner
{isLoading && <Spinner />}

// ✅ Skeleton loader (лучший UX)
{isLoading ? (
  <>
    <Skeleton height={60} borderRadius={12} />
    <Skeleton height={60} borderRadius={12} />
    <Skeleton height={60} borderRadius={12} />
  </>
) : (
  achievementsList
)}
```

---

### 12. 🎯 Pull to Refresh

**Рекомендация:**

Добавить pull-to-refresh для главного экрана:

```typescript
import { PullToRefresh } from '@vkontakte/vkui'

function HomePage() {
  const [isFetching, setIsFetching] = useState(false)

  const handleRefresh = async () => {
    setIsFetching(true)
    await refetchHomeData()
    setIsFetching(false)
  }

  return (
    <PullToRefresh onRefresh={handleRefresh} isFetching={isFetching}>
      {/* Контент */}
    </PullToRefresh>
  )
}
```

---

### 13. 📱 Swipe Actions для списков

**Рекомендация:**

Использовать swipe actions в истории попыток:

```typescript
import { SimpleCell, useActionSheet } from '@vkontakte/vkui'

function HistoryItem({ attempt }) {
  const { showActionSheet } = useActionSheet()

  const handleSwipe = () => {
    showActionSheet([
      {
        title: 'Повторить',
        action: () => retryAttempt(attempt.id),
        mode: 'default'
      },
      {
        title: 'Удалить',
        action: () => deleteAttempt(attempt.id),
        mode: 'destructive'
      }
    ])
  }

  return (
    <SimpleCell
      onSwipeRight={handleSwipe}
      badge={attempt.status}
    >
      {attempt.title}
    </SimpleCell>
  )
}
```

---

## 📊 Performance Улучшения

### 14. ⚡ Virtualized Lists для истории

**Рекомендация:**

Для длинных списков использовать виртуализацию:

```typescript
import { FixedSizeList } from 'react-window'

function HistoryList({ attempts }) {
  const Row = ({ index, style }) => (
    <div style={style}>
      <HistoryCard attempt={attempts[index]} />
    </div>
  )

  return (
    <FixedSizeList
      height={600}
      itemCount={attempts.length}
      itemSize={80}
      width="100%"
    >
      {Row}
    </FixedSizeList>
  )
}
```

---

### 15. 🖼️ Progressive Image Loading

**Рекомендация:**

```typescript
import { Image } from '@vkontakte/vkui'

function AttemptImage({ url }) {
  return (
    <Image
      src={url}
      size={64}
      fallbackIcon={<Icon28PictureOutline />}
      withBorder
    />
  )
}
```

---

## 🔐 Security Улучшения

### 16. 🛡️ Rate Limiting на Frontend

**Рекомендация:**

Добавить debounce для частых действий:

```typescript
import { useDebouncedCallback } from 'use-debounce'

function SearchInput() {
  const search = useDebouncedCallback(
    (query) => api.search(query),
    500 // 500ms задержка
  )

  return <Input onChange={(e) => search(e.target.value)} />
}
```

---

### 17. 📸 Image Upload Security

**Рекомендация:**

Добавить проверку изображений на клиенте:

```typescript
import imageCompression from 'browser-image-compression'

async function handleImageUpload(file: File) {
  // 1. Проверка типа
  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp']
  if (!allowedTypes.includes(file.type)) {
    throw new Error('Неподдерживаемый формат')
  }

  // 2. Проверка размера (макс 10MB)
  if (file.size > 10 * 1024 * 1024) {
    throw new Error('Файл слишком большой')
  }

  // 3. Сжатие
  const compressed = await imageCompression(file, {
    maxSizeMB: 1,
    maxWidthOrHeight: 1920,
    useWebWorker: true
  })

  // 4. Загрузка
  await api.uploadImage(compressed)
}
```

---

## 📋 Обновленные Checklists

### Pre-Development Checklist

- [ ] Установить VKUI и иконки
- [ ] Настроить VK Bridge
- [ ] Настроить проверку sign на backend
- [ ] Настроить bundle size limits
- [ ] Настроить platform detection
- [ ] Интегрировать haptics
- [ ] Использовать VK Storage вместо localStorage
- [ ] Настроить VK Pay (если монетизация)
- [ ] Настроить Ads SDK (опционально)

### Pre-Release Checklist

- [ ] Bundle size < 10 MB (проверить с `npm run build`)
- [ ] Sign validation работает на backend
- [ ] Все haptics добавлены для основных действий
- [ ] Pull-to-refresh работает на главном экране
- [ ] Share через VK Bridge работает
- [ ] VK Pay integration протестирована
- [ ] Ads показываются корректно
- [ ] Все изображения оптимизированы (WebP)
- [ ] Performance метрики в норме (Lighthouse > 80)
- [ ] Accessibility проверена

### VK Moderation Checklist

- [ ] HTTPS настроен
- [ ] Sign validation на backend
- [ ] Нет критических security уязвимостей
- [ ] Приложение работает на всех платформах (iOS, Android, Web)
- [ ] UI соответствует VK guidelines (используем VKUI)
- [ ] Нет недоступных функций без объяснения
- [ ] Privacy policy добавлена
- [ ] Terms of service добавлены
- [ ] Монетизация соответствует правилам VK
- [ ] Контент не нарушает правила VK

---

## 📁 Новые файлы для roadmap

Рекомендуется создать:

1. **`14_MONETIZATION.md`** - детальный план монетизации
   - VK Pay integration
   - Ads SDK integration
   - Subscription logic
   - Free tier strategy

2. **`15_VK_MODERATION.md`** - подготовка к модерации
   - Checklist
   - Типичные ошибки
   - Timing (релизы по четвергам)
   - Escalation process

3. **`VKUI_MIGRATION.md`** - миграция с Custom UI на VKUI
   - Mapping компонентов
   - Примеры рефакторинга
   - Breaking changes

---

## 🎯 Приоритизация

### Must Have (до первого релиза)

1. ✅ Интеграция VKUI
2. ✅ VK Bridge sign validation
3. ✅ Bundle size optimization
4. ✅ Platform detection
5. ✅ Haptic feedback
6. ✅ VK Storage API

### Should Have (в первой итерации)

7. ✅ VK Pay integration
8. ✅ Share через VK Bridge
9. ✅ Skeleton loaders
10. ✅ Pull to refresh

### Nice to Have (в будущих версиях)

11. ⚡ QR scanner
12. ⚡ Ads integration
13. ⚡ Virtualized lists
14. ⚡ Swipe actions

---

## 📚 Ссылки

- [VK_BEST_PRACTICES.md](./VK_BEST_PRACTICES.md) - полный гайд
- [VK Mini Apps Documentation](https://dev.vk.com/mini-apps)
- [VKUI Documentation](https://vkcom.github.io/VKUI/)
- [VK Bridge API](https://dev.vk.com/bridge/overview)

---

## ✅ Следующие шаги

1. Обновить `01_SETUP.md` с VKUI установкой
2. Переписать `02_CORE.md` с VKUI вместо Custom UI Kit
3. Добавить VK Bridge security в `SECURITY.md`
4. Создать `14_MONETIZATION.md`
5. Создать `15_VK_MODERATION.md`
6. Обновить все roadmap файлы с haptics, storage, share

**Roadmap будет полностью соответствовать VK best practices!** ✨
