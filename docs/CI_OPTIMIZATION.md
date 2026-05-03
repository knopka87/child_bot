# Оптимизация GitHub Actions CI/CD

## Текущее состояние
- ⏱️ Время сборки: **~20 минут**
- Причины медленной сборки:
  1. Multi-platform builds (amd64 + arm64) - удваивает время
  2. Последовательная сборка backend и frontend
  3. Chromium устанавливается в runtime образе каждый раз
  4. node:20-bookworm (полный образ ~1GB) вместо slim

## Оптимизации

### 1. Убрать multi-platform builds ⚡ (-50% времени)

**Проблема:** Сервер работает на `77.222.60.149` (x86_64/amd64), сборка для arm64 не нужна.

**Решение:** Изменить `.github/workflows/docker-publish.yml`

```yaml
# Было
platforms: linux/amd64,linux/arm64

# Стало
platforms: linux/amd64
```

**Результат:** Сборка в 2 раза быстрее

---

### 2. Параллельная сборка backend + frontend ⚡ (-40% времени)

**Проблема:** Backend и frontend собираются последовательно

**Решение:** Использовать два job'а которые выполняются параллельно

Файл: `.github/workflows/docker-publish-optimized.yml` (уже создан)

**Результат:** Backend и frontend собираются одновременно

---

### 3. Оптимизация Frontend Dockerfile ⚡ (-20% времени)

**Проблема:**
- Используется `node:20-bookworm` (~1GB)
- Ретраи при установке npm могут замедлять

**Решение:** Использовать `node:20-bookworm-slim` и nginx:alpine

Файл: `frontend/Dockerfile.optimized` (уже создан)

**Результат:** Меньший размер образа, быстрее скачивание и сборка

---

### 4. Динамический GOARCH в Backend ✅

**Проблема:** Жестко прописан `GOARCH=amd64`, несовместимо с multi-platform

**Решение:** Использовать build args `TARGETOS` и `TARGETARCH`

**Результат:** Гибкость при сборке для разных платформ (если потребуется)

---

## Применение изменений

### Вариант 1: Быстрая оптимизация (рекомендуется)

Просто убираем arm64 из текущего workflow:

```bash
# Редактируем .github/workflows/docker-publish.yml
sed -i '' 's/linux\/amd64,linux\/arm64/linux\/amd64/g' .github/workflows/docker-publish.yml

git add .github/workflows/docker-publish.yml Dockerfile
git commit -m "perf: optimize CI build - remove arm64, add dynamic GOARCH"
git push origin prod-v1
```

**Ожидаемое время сборки:** ~10 минут

---

### Вариант 2: Полная оптимизация

Использовать новые оптимизированные файлы:

```bash
# 1. Заменяем workflow
mv .github/workflows/docker-publish.yml .github/workflows/docker-publish.old.yml
mv .github/workflows/docker-publish-optimized.yml .github/workflows/docker-publish.yml

# 2. Заменяем frontend Dockerfile
mv frontend/Dockerfile frontend/Dockerfile.old
mv frontend/Dockerfile.optimized frontend/Dockerfile

# 3. Коммитим
git add .github/workflows/ frontend/Dockerfile Dockerfile
git commit -m "perf: full CI optimization - parallel builds, slim images"
git push origin prod-v1
```

**Ожидаемое время сборки:** ~5-7 минут

---

## Дополнительные оптимизации (опционально)

### A. Self-hosted runner

Если сборки запускаются часто, можно поднять свой runner:

```yaml
jobs:
  build-backend:
    runs-on: self-hosted  # вместо ubuntu-latest
```

**Плюсы:**
- Персистентный кеш между сборками
- Более мощное железо
- Бесплатно (в отличие от GitHub hosted runners)

**Минусы:**
- Нужно поддерживать сервер
- Настройка безопасности

---

### B. Кеширование npm в GitHub Actions

Добавить перед сборкой frontend:

```yaml
- name: Cache node modules
  uses: actions/cache@v3
  with:
    path: frontend/node_modules
    key: ${{ runner.os }}-node-${{ hashFiles('frontend/package-lock.json') }}
```

**Результат:** +1-2 минуты экономии на npm install

---

### C. Использовать Docker Layer Caching

Уже включено через `cache-from: type=gha`, но можно улучшить:

```yaml
cache-from: |
  type=gha
  type=registry,ref=ghcr.io/${{ github.repository_owner }}/child_bot-backend:latest
cache-to: type=gha,mode=max
```

---

## Сравнение

| Вариант | Время сборки | Сложность внедрения |
|---------|--------------|---------------------|
| Текущий | ~20 минут | - |
| Вариант 1 (быстрый) | ~10 минут | 🟢 Низкая (1 команда) |
| Вариант 2 (полный) | ~5-7 минут | 🟡 Средняя (тестирование) |
| + Self-hosted runner | ~3-4 минуты | 🔴 Высокая (настройка) |

---

## Рекомендация

**Начать с Варианта 1** - убрать arm64 и добавить динамический GOARCH:

```bash
# Одна команда
sed -i '' 's/linux\/amd64,linux\/arm64/linux\/amd64/g' .github/workflows/docker-publish.yml && \
git add .github/workflows/docker-publish.yml Dockerfile && \
git commit -m "perf: optimize CI - remove arm64, use dynamic GOARCH" && \
git push origin prod-v1
```

Это даст **50% ускорение** (20 мин → 10 мин) без рисков.

После проверки можно перейти на Вариант 2 для дальнейшего ускорения.

---

## Мониторинг

После внедрения проверить:
- GitHub Actions → последний workflow run → время выполнения
- Размеры образов: `docker images | grep child_bot`
- Логи сборки на наличие ошибок

---

## Откат

Если что-то пошло не так:

```bash
# Откат на старый workflow
git revert HEAD
git push origin prod-v1

# Или восстановить из бэкапа
mv .github/workflows/docker-publish.old.yml .github/workflows/docker-publish.yml
git add .github/workflows/docker-publish.yml
git commit -m "revert: restore old CI workflow"
git push origin prod-v1
```
