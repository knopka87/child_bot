You are a Senior Product Designer working in Figma Make.

You already have a generated mobile mini-app.
Your task is to APPLY UI and UX UPDATES (do not rebuild from scratch).

Focus on:
• mobile-first layout
• game-like UI (similar to children mobile games)
• clean structure
• no overload

---

# 1. HOME SCREEN — FULL TOP RESTRUCTURE (CRITICAL)

Redesign the TOP AREA using a GAME UI approach.

REFERENCE:
Use separate UI elements like in mobile games (level badge, progress bar, currency counters).

---

## 1.1 TOP GAME BAR (replace single strip)

Instead of one line → create SEPARATE ELEMENTS:

LEFT:
• Level badge (e.g. "3")
• Progress bar to next level (with %)

CENTER or LEFT-MID:
• Cups counter (number of correctly solved tasks)

RIGHT:
• Coins counter (with icon + number)

Style:
• Rounded containers
• Bright but clean
• Game-like but not noisy
• Each element visually separated (like chips/cards)

---

# 2. CHARACTER BLOCK (UNDER TOP BAR)

Create a horizontal composition:

LEFT:
• Mascot (use uploaded file: Герой_Лидер.png)

RIGHT:
• Villain (use uploaded file: Монстр 1.png)

IMPORTANT:
• Replace ALL previous characters with these assets
• Do NOT generate new characters
• Preserve proportions and recognizability

---

## 2.1 CHARACTER INTERACTION

Villain:
• health bar (3 segments or smooth bar)
• aggressive pose (attacking feeling)

Mascot:
• speech bubble (dialog style)
• supportive message

Layout:
• characters at same vertical level
• no overlap with CTA buttons

---

# 3. CTA BUTTONS (POSITION FIX)

Move buttons LOWER (thumb zone):

• Помоги разобраться
• Проверка ДЗ

Style:
• large
• high contrast
• easy tap

---

# 4. ACHIEVEMENTS — COMPLETE REDESIGN

Change layout to SHELVES system.

Structure:

• horizontal shelves (like a room / rack)
• items placed on shelves

States:

OPENED STICKERS:
• colorful
• active

LOCKED STICKERS:
• grayscale / faded
• slightly transparent

---

## INTERACTION

On click or hover:

IF unlocked:
→ show modal:
"За что получено"

IF locked:
→ show modal:
"Что нужно сделать"

---

# 5. HINT FLOW CHANGE (VERY IMPORTANT UX FIX)

Modify logic in HELP MODE:

When user presses:

→ "Отправить ответ"

NEW FLOW:

1. Open image source selection:
   • Выбрать изображение
   • Сфотографировать

2. Continue standard flow:
   upload → confirm → crop → processing → result

---

## RESULT AFTER WRONG ANSWER

Instead of generic error:

Show:

• friendly message
• clear indication WHERE the mistake is
• short hint what to fix

NO harsh wording.

Examples tone:
• "Посмотри этот шаг"
• "Здесь небольшая ошибка"
• "Попробуй ещё раз"

---

# 6. GENERAL STYLE UPDATE

Apply:

• mobile game UI inspiration
• but keep minimalism
• avoid overload

Use:

• soft shadows
• rounded corners
• bright accents

DO NOT:
• turn UI into cartoon chaos
• overload with decorative elements

---

# 7. FINAL REQUIREMENTS

You must:

• update existing screens
• replace characters with provided assets
• redesign top bar into game UI
• rebuild achievements screen as shelves
• update hint → answer → photo flow
• keep all flows connected

---

# FINAL CHECK

Make sure:

• top looks like mobile game UI
• characters feel alive but not dominant
• CTA buttons are easy to reach
• hint flow is logical
• achievements feel collectible

If something feels overloaded → simplify.