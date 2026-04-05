# Screen Map — Карта экранов приложения

**Проект:** Объяснятель ДЗ MiniApp
**Дата:** 2026-03-29

---

## Визуальная карта навигации

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           ONBOARDING FLOW                                │
└─────────────────────────────────────────────────────────────────────────┘

[Start] → [Onboarding1] → [Onboarding2] → [Onboarding3]
            ↓
       [Registration]
            ↓
       [Grade Selection] (1-11 класс)
            ↓
       [Avatar Selection] (выбор персонажа)
            ↓
       [Display Name] (ввод имени)
            ↓
       [Consent Screen] (согласия + email)
            ↓
       [Email Verification] (подтверждение кода)
            ↓
       [Onboarding Complete] → [Home]


┌─────────────────────────────────────────────────────────────────────────┐
│                              MAIN APP                                    │
└─────────────────────────────────────────────────────────────────────────┘

                            ┌─────────────┐
                            │    HOME     │
                            │             │
                            │ ┌─────────┐ │
                            │ │ Header  │ │ Level, Coins, Tasks
                            │ └─────────┘ │
                            │             │
                            │  🐻 ⚔️ 👑  │ Mascot vs Villain
                            │             │
                            │ [Помоги]    │
                            │ [Проверка]  │
                            │             │
                            │ Recent      │
                            │ Attempts    │
                            │             │
                            └──┬──────┬───┘
                               │      │
              ┌────────────────┘      └────────────────┐
              │                                        │
              ▼                                        ▼
    ┌──────────────────┐                    ┌──────────────────┐
    │  HELP FLOW       │                    │  CHECK FLOW      │
    └──────────────────┘                    └──────────────────┘
              │                                        │
              │                                        │
              ▼                                        ▼
    [Upload Source]                          [Scenario Select]
    - Camera                                 - 1 Photo
    - File                                   - 2 Photos
    - Clipboard                                      │
    - Drag&Drop                                      │
              │                                      ▼
              ▼                              [Upload Task Photo]
    [Image Upload]                                   │
              │                                      ▼
              ▼                              [Upload Answer Photo]
    [Quality Check]                                  │
    - ✓ Confirm                                      ▼
    - ✗ Reshoot                              [Quality Check]
    - ✂ Manual Crop                                  │
              │                                      ▼
              ▼                                 [Processing]
    [Processing]                                     │
    - Loading                                        ▼
    - Long Wait                               [Check Result]
    - Save & Wait                             - ✅ All Correct
              │                               - ⚠️ Has Errors
              ▼                               - ❌ All Wrong
    [Result]                                         │
    - Hint Level 1                                   │
    - Hint Level 2                                   ▼
    - Hint Level 3                            [Error Details]
    - Submit Answer                           - Step number
              │                               - Line reference
              ▼                               - Hint
    [Reward]                                         │
    + Coins                                          ▼
    + XP                                      [Fix & Resubmit]
    - Villain Damage                          or [New Task]


┌─────────────────────────────────────────────────────────────────────────┐
│                          BOTTOM NAVIGATION                               │
└─────────────────────────────────────────────────────────────────────────┘

    [Главная]     [Достижения]     [Друзья]     [Профиль]
        │              │              │              │
        │              │              │              │
        ▼              ▼              ▼              ▼
      HOME       ACHIEVEMENTS      FRIENDS       PROFILE


┌─────────────────────────────────────────────────────────────────────────┐
│                          ACHIEVEMENTS                                    │
└─────────────────────────────────────────────────────────────────────────┘

[Achievements Screen]
│
├─ [Achievement Shelf 1] (Streak & Daily)
│   ├─ 🔥 5 дней подряд [locked/unlocked]
│   ├─ ✅ 10 проверок ДЗ [locked/unlocked]
│   ├─ ⭐ 5 ошибок исправлено [locked/unlocked]
│   └─ 🎯 Первое задание [locked/unlocked]
│
├─ [Achievement Shelf 2] (Mastery)
│   ├─ ⚡ Скоростной решатель [locked/unlocked]
│   ├─ 🏆 Победитель злодеев [locked/unlocked]
│   ├─ 🦸 Мудрая сова [locked/unlocked]
│   └─ 💎 Коллекционер [locked/unlocked]
│
└─ [Achievement Shelf 3] (Social & Advanced)
    ├─ 🚀 Ракета знаний [locked/unlocked]
    ├─ 🌟 Суперзвезда [locked/unlocked]
    ├─ 🎊 Марафонец [locked/unlocked]
    └─ 🧠 Гений [locked/unlocked]
        │
        ▼
[Achievement Detail Modal]
- Icon & Name
- Description
- Progress (if locked)
- Reward (coins, sticker)
- Claim button (if unlocked & not claimed)


┌─────────────────────────────────────────────────────────────────────────┐
│                             FRIENDS                                      │
└─────────────────────────────────────────────────────────────────────────┘

[Friends Screen]
│
├─ [Referral Progress Card]
│   ├─ "Пригласи 5 друзей — получи редкий стикер!"
│   ├─ Progress Indicators: ✓✓345 (2 из 5)
│   └─ Reward Preview: ⭐ Редкий стикер «Дружба»
│
├─ [Invite Friend Card]
│   ├─ Referral Link: https://homework.app/invite/abc123
│   ├─ [Скопировать] button
│   └─ [Отправить] button
│       │
│       ▼
│   [Share Sheet]
│   - Telegram
│   - VK
│   - WhatsApp
│   - Copy Link
│
└─ [Invited Friends List]
    └─ "Приглашено друзей: 2"


┌─────────────────────────────────────────────────────────────────────────┐
│                             PROFILE                                      │
└─────────────────────────────────────────────────────────────────────────┘

[Profile Screen]
│
├─ [Profile Card]
│   ├─ 🦊 Avatar
│   ├─ Артём (display name)
│   ├─ 2 класс (grade)
│   └─ Пробный период — 5 дней (trial status)
│
├─ [История] → [History Screen]
│   │
│   ├─ [Filter Tabs]
│   │   ├─ Все
│   │   ├─ Помощь
│   │   ├─ Проверка
│   │   ├─ Правильно
│   │   └─ Ошибки
│   │
│   ├─ [History Items List]
│   │   └─ [Item] → [History Detail Modal]
│   │       ├─ Thumbnail
│   │       ├─ Mode (help/check)
│   │       ├─ Result status
│   │       ├─ Date
│   │       ├─ [Повторить]
│   │       └─ [Исправить и проверить]
│   │
│   └─ [Pagination]
│
├─ [Отчёт родителю] → [Report Settings]
│   │
│   ├─ [Email Input]
│   ├─ [Weekly Report Toggle]
│   ├─ [Archive Toggle]
│   └─ [Archive] → [Reports Archive]
│       └─ [Report Items]
│           ├─ Period
│           ├─ Summary
│           └─ [Download PDF]
│
├─ [Помощь] → [Support Screen]
│   │
│   ├─ [FAQ]
│   └─ [Contact Form]
│       ├─ Message textarea
│       └─ [Отправить]
│
└─ [Подписка] → [Subscription Screen]
    │
    ├─ Current Plan Info
    ├─ Trial Days Left
    └─ [Управлять подпиской] → [Paywall]
        │
        ├─ [Pricing Plans]
        │   ├─ Monthly
        │   ├─ Quarterly
        │   └─ Annual (popular)
        │
        ├─ [Features List]
        │   ├─ ✓ Unlimited tasks
        │   ├─ ✓ All hints
        │   ├─ ✓ Priority support
        │   └─ ✓ Weekly reports
        │
        └─ [Оформить подписку]
            │
            ▼
        [Payment Screen]
        │
        └─ [Success] or [Error]


┌─────────────────────────────────────────────────────────────────────────┐
│                          VILLAIN FLOW                                    │
└─────────────────────────────────────────────────────────────────────────┘

[Home: Villain Preview]
- 👑 Злодей Кракозябра
- Health Bar (0-100%)
- Click → [Villain Detail]

[Villain Detail Screen]
│
├─ [Villain Character]
│   └─ White character with crown
│
├─ [Health Bar]
│   └─ Current: 67% / Max: 100%
│
├─ [Taunt Message]
│   └─ "Ха-ха! Попробуй-ка реши задачки!"
│
├─ [Defeat Requirements]
│   └─ "Реши 3 задания правильно, чтобы победить"
│   └─ Progress: 2/3 ✓✓○
│
└─ [Reward Preview]
    ├─ 500 Coins
    ├─ Редкий стикер
    └─ Достижение "Победитель злодеев"

[Victory Trigger]
- After 3 correct answers
    │
    ▼
[Victory Screen]
│
├─ [Confetti Animation] 🎊
│
├─ [Victory Message]
│   └─ "Ты победил Кракозябру!"
│
├─ [Rewards]
│   ├─ + 500 Coins
│   ├─ 🌟 Редкий стикер
│   └─ 🏆 Достижение
│
└─ [Продолжить учиться] → [Home]


┌─────────────────────────────────────────────────────────────────────────┐
│                            MODALS                                        │
└─────────────────────────────────────────────────────────────────────────┘

[Unfinished Attempt Modal] (Home)
├─ "У тебя есть незаконченное задание. Хочешь продолжить?"
├─ [Продолжить] → Resume attempt
└─ [Новое задание] → Discard and start new

[Paywall Modal] (любой экран при блокировке)
├─ "Разблокируй все возможности!"
├─ Blocked Feature context
├─ [Pricing Plans]
└─ [Оформить подписку] or [Закрыть]

[Achievement Unlock Notification]
├─ [Confetti Animation] 🎊
├─ Icon & Name
├─ Description
├─ Reward
└─ [Получить] or [Закрыть]

[Error Modal]
├─ Error Icon
├─ Error Message
└─ [OK] or [Повторить]


┌─────────────────────────────────────────────────────────────────────────┐
│                     NAVIGATION SUMMARY                                   │
└─────────────────────────────────────────────────────────────────────────┘

Entry Points:
├─ [First Launch] → Onboarding
├─ [Returning User] → Home
└─ [Deep Link] → Specific Screen

Bottom Tab Navigation (always visible in main app):
├─ [Главная] → Home
├─ [Достижения] → Achievements
├─ [Друзья] → Friends
└─ [Профиль] → Profile

Modal Navigation (overlay, can close):
├─ Unfinished Attempt Modal
├─ Paywall Modal
├─ Achievement Unlock Notification
├─ Error Modal
├─ Villain Detail
├─ Victory Screen
├─ History Detail
└─ Achievement Detail

Full-Screen Flows (replace entire screen, have back button):
├─ Onboarding Flow
├─ Help Flow
├─ Check Flow
├─ History Screen
├─ Support Screen
├─ Report Settings
└─ Subscription Screen
