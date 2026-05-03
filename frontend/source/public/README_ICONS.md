# Иконки и изображения для PWA

## Требуемые файлы

Для корректной работы PWA и улучшения Lighthouse score нужно добавить следующие файлы в эту директорию:

### Обязательные иконки

1. **favicon-16x16.png** (16x16px)
   - Favicon для старых браузеров

2. **favicon-32x32.png** (32x32px)
   - Стандартный favicon

3. **apple-touch-icon.png** (180x180px)
   - Иконка для iOS Safari
   - Должна быть без прозрачности (solid background)

4. **android-chrome-192x192.png** (192x192px)
   - Иконка для Android

5. **android-chrome-512x512.png** (512x512px)
   - Большая иконка для Android
   - Используется для splash screen

### Дополнительные файлы

6. **og-image.png** (1200x630px)
   - Open Graph изображение для превью в соц.сетях
   - Рекомендуемый размер: 1200x630px

7. **favicon.ico** (необязательно)
   - Multi-size ICO файл (16x16, 32x32, 48x48)

## Как создать иконки

### Вариант 1: Онлайн генераторы
- https://realfavicongenerator.net/
- https://www.favicon-generator.org/

### Вариант 2: Figma/Photoshop
1. Создайте квадратное изображение 512x512px
2. Добавьте логотип/символ приложения
3. Экспортируйте в нужных размерах

### Вариант 3: ImageMagick (командная строка)
```bash
# Из исходного файла logo.png создать все размеры
convert logo.png -resize 16x16 favicon-16x16.png
convert logo.png -resize 32x32 favicon-32x32.png
convert logo.png -resize 180x180 apple-touch-icon.png
convert logo.png -resize 192x192 android-chrome-192x192.png
convert logo.png -resize 512x512 android-chrome-512x512.png
convert logo.png -resize 1200x630 og-image.png
```

## Рекомендации по дизайну

### Для иконок приложения:
- ✅ Используйте простой, узнаваемый символ
- ✅ Яркие цвета (совпадающие с брендом)
- ✅ Хорошая контрастность
- ✅ Работает в маленьком размере (32x32)
- ❌ Избегайте мелких деталей
- ❌ Не используйте текст (нечитаем в маленьком размере)

### Для OG image:
- ✅ Название приложения крупным шрифтом
- ✅ Краткое описание
- ✅ Визуальные элементы (маскот, скриншот)
- ✅ Брендовые цвета

## Временное решение (для разработки)

Пока иконки не созданы, можно использовать:
- Однотонные квадраты с буквой "Д" (ДЗ Объяснитель)
- Цвет: #0077FF (theme color)
- Белая буква на синем фоне

## Проверка после добавления

1. Запустите Lighthouse audit:
```bash
npm run build
npm run preview
# Откройте DevTools → Lighthouse → Generate report
```

2. Проверьте что:
   - ✅ PWA installable
   - ✅ Icons загружаются без 404
   - ✅ Manifest корректный
   - ✅ OG image отображается в превью

## Статус

- [ ] favicon-16x16.png
- [ ] favicon-32x32.png
- [ ] apple-touch-icon.png
- [ ] android-chrome-192x192.png
- [ ] android-chrome-512x512.png
- [ ] og-image.png
- [x] manifest.json (создан)
