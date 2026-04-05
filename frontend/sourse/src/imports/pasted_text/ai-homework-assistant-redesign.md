You are a Senior Product Designer and UX Architect.

You already generated a mobile mini-app UI.

Now your task is to UPDATE and RESTRUCTURE the existing design based on new product requirements.

IMPORTANT:
• Do NOT recreate from scratch
• Update structure, flows, and screens
• Keep mobile-first design (390×844)
• Maintain consistency of design system
• Focus on UX improvements and flow clarity

---

CONTEXT

The product is an AI homework assistant for children.

Primary goal:
FAST ACTION → HELP OR CHECK HOMEWORK

Gamification is SECONDARY and must NOT dominate UI.

---

CRITICAL UX PRINCIPLES

• One screen = one main goal
• No competing actions
• Clear linear flows
• Support ALL real user scenarios
• Interface must remain simple for children

---

GLOBAL NAVIGATION UPDATE

Rename tab:

Штаб → Главная

Tabs remain:

1. Главная
2. Достижения
3. Друзья
4. Профиль

---

HOME SCREEN (MAJOR UPDATE)

Transform Home into ACTION-FIRST screen.

Must include:

1. Two primary CTA buttons (dominant):
   • Помоги разобраться
   • Проверка ДЗ

2. Replace "Совет дня" with:
   → Short mascot joke / playful text

3. ADD villain block:
   • character
   • health bar (3 segments)
   • short phrase
   • tappable → opens Villain Screen

4. ADD compact blocks:
   • summary (progress overview)
   • recent attempts (short list)

IMPORTANT:
• Mascot must NOT dominate
• Avoid visual clutter
• Focus on actions

---

NEW SCREEN: VILLAIN SCREEN

Create new screen.

Elements:

• Villain (large)
• Health bar
• Explanation:
  how to defeat villain
  what user gets

CTA:
• continue learning

---

NEW LOGIC: UNFINISHED TASK

If user has unfinished session:

Show MODAL:

• Продолжить
• Новое задание

Trigger:

• entering new flow
• returning to Home

---

HELP FLOW (UPDATED)

Support:

• 1 image
• 2 images of SAME task (multi-page)

UPLOAD OPTIONS:

Mobile:
• Выбрать изображение
• Сфотографировать

Desktop:
• Upload
• Paste
• Drag & Drop

---

NEW SCREEN: IMAGE SET MANAGER

For multi-image scenarios.

Elements:

• thumbnails
• labels:
  задание / ответ / страница 1 / страница 2

Actions:

• replace
• delete
• reorder
• continue

---

NEW SCREEN: IMAGE QUALITY CHECK

Before processing.

Title:
Всё ли видно?

Elements:

• preview
• confirm button
• recapture button
• crop button

---

NEW SCREEN: CROP

Mandatory step after upload.

Text:
Обрежь до одного задания

Buttons:

• Подтвердить
• Переснять

---

WAITING SCREEN (UPDATED)

Now must support LONG WAIT states.

States:

Normal:
• Думаю…
• Проверяю решение…

Long wait actions:

• Сохранить и подождать
• Повторить
• Отменить

Processing continues on server.

---

HELP RESULT SCREEN (UPDATED)

REMOVE:
• repeated task text
• "Похожее задание"

KEEP:

• hints L1 / L2 / L3 (sequential unlock)
• villain health
• actions:

Buttons:

• Отправить ответ
• Новое задание

---

CHECK FLOW (MAJOR UPDATE)

ADD ENTRY SCREEN:

User must SELECT scenario:

1. 1 фото (задание + ответ)
2. 2 фото (задание отдельно, ответ отдельно)
3. 2 фото задания (multi-page task)

IMPORTANT:
• No auto switching
• User stays in selected flow

---

CHECK RESULT SCREEN (REDESIGN)

Create 3 SEPARATE UX STATES:

1. Верно
2. Почти верно
3. Неверно

---

STATE: ВЕРНО

If villain still alive:

• reduce health
• happy mascot

CTA:
Новое задание

---

STATE: ПОЧТИ ВЕРНО

• show correction hint

Buttons:

• Исправил
• Новое задание

---

STATE: НЕВЕРНО

• encouraging message
• mascot support

Buttons:

• Попробовать ещё раз
• Новое задание

---

NEW SCREEN: VICTORY

Trigger:

• villain health = 0
• last answer correct

Show:

• defeated villain
• celebration
• reward (sticker / achievement)

CTA:
continue

---

ACHIEVEMENTS (UPDATE)

REMOVE mascot.

Focus:

• trophies
• progress
• stats

---

FRIENDS (UPDATE)

ADD mechanic:

Invite 5 friends → get rare sticker

Show:

• progress counter
• reward preview

---

PROFILE (SIMPLIFY)

REMOVE:

• email editing
• sound
• language

KEEP:

• child profile
• history
• parent report
• help
• access to adult flow

---

HISTORY (IMPROVED)

Each item must include:

• short task description
• status

Statuses:

• Решено верно
• Решено неверно
• Незакончено
• Использованы подсказки
• В обработке

Add:

• detailed attempt screen

Actions:

• Повторить
• Исправить и проверить

---

PARENT REPORT (SIMPLIFY)

Keep only:

• email
• archive
• toggle on/off

REMOVE:

• time/day selection

Fix schedule:

Monday, 10:00

---

REGISTRATION (UPDATE)

Add explicit FIRST ENTRY FLOW:

• consent screen (adult)
• privacy acceptance
• then child profile

---

PAYMENT (EXPAND)

Must include states:

• Paywall
• Plan selection
• Success
• Error

After payment:
return user to previous flow

---

ARCHITECTURE REQUIREMENTS

You MUST reflect:

• all user scenarios:
  - help
  - check
  - 1 photo
  - 2 photos
  - unfinished task
  - long wait

• all new screens:
  - villain screen
  - image set manager
  - quality check
  - crop
  - scenario selector
  - victory screen

---

DELIVERABLE

Update Figma file:

• modify existing screens
• add new screens
• update flows
• connect interactions
• include all states
• ensure mobile layout