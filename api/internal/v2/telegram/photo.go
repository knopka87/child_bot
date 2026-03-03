package telegram

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/util"
)

const (
	// Максимальное время жизни батча фото (защита от зависших батчей)
	batchMaxAge = 5 * time.Minute
)

var batches sync.Map // key -> *photoBatch

// createBatchTimer создаёт таймер с трекингом для graceful shutdown
// Возвращает таймер, который автоматически отслеживается ShutdownManager
func (r *Router) createBatchTimer(delay time.Duration, key string, userID *int64, cid int64, source string) *time.Timer {
	shutdown := GetShutdownManager()

	// Не создаём таймер если идёт shutdown
	if shutdown.IsShutdown() {
		return time.NewTimer(0) // dummy timer
	}

	// Регистрируем горутину ДО создания таймера
	done := shutdown.TrackGoroutine()

	return time.AfterFunc(delay, func() {
		defer done()

		// Проверяем shutdown перед выполнением
		if shutdown.IsShutdown() {
			return
		}

		defer func() {
			if rec := recover(); rec != nil {
				util.PrintError(source, "", cid, "panic in timer callback", fmt.Errorf("%v", rec))
			}
		}()
		r.processBatch(key, userID)
	})
}

type photoBatch struct {
	ChatID       int64
	Key          string // "grp:<mediaGroupID>" | "chat:<chatID>"
	MediaGroupID string

	mu        sync.Mutex
	images    [][]byte
	timer     *time.Timer
	createdAt time.Time // время создания батча
	lastAt    time.Time
	processed bool // флаг, что batch уже обработан
}

// cleanupStaleBatches удаляет зависшие батчи старше batchMaxAge
func cleanupStaleBatches() int {
	var count int
	now := time.Now()

	batches.Range(func(key, value interface{}) bool {
		b, ok := value.(*photoBatch)
		if !ok {
			// Unexpected type - delete corrupted entry
			batches.Delete(key)
			count++
			return true
		}
		b.mu.Lock()
		age := now.Sub(b.createdAt)
		b.mu.Unlock()

		if age > batchMaxAge {
			// Останавливаем таймер если он ещё активен
			if b.timer != nil {
				b.timer.Stop()
			}
			batches.Delete(key)
			count++
		}
		return true
	})

	return count
}

func (r *Router) acceptPhoto(cid int64, msg tgbotapi.Message) {
	ph := msg.Photo[len(msg.Photo)-1]
	imgBytes, err := r.Bot.DownloadFile(ph.FileID)
	if err != nil {
		r.sendError(cid, err)
		return
	}

	key := "chat:" + fmt.Sprint(cid)
	if msg.MediaGroupID != "" {
		key = "grp:" + msg.MediaGroupID
	}

	bi, _ := batches.LoadOrStore(key, &photoBatch{
		ChatID:       cid,
		Key:          key,
		MediaGroupID: msg.MediaGroupID,
		images:       make([][]byte, 0, 2),
		createdAt:    time.Now(),
	})
	b, ok := bi.(*photoBatch)
	if !ok {
		r.sendError(cid, fmt.Errorf("internal error: invalid batch type"))
		batches.Delete(key)
		return
	}

	b.mu.Lock()
	// Проверяем, не обработан ли уже batch
	if b.processed {
		b.mu.Unlock()
		// Batch уже обработан — создаём новый
		newBatch := &photoBatch{
			ChatID:       cid,
			Key:          key,
			MediaGroupID: msg.MediaGroupID,
			images:       [][]byte{imgBytes},
			createdAt:    time.Now(),
			lastAt:       time.Now(),
		}
		userID := util.GetUserIDFromTgMessage(msg)
		newBatch.timer = r.createBatchTimer(debounce, key, userID, cid, "acceptPhoto.newBatch")
		batches.Store(key, newBatch)
		return
	}
	b.images = append(b.images, imgBytes)
	b.lastAt = time.Now()
	if b.timer != nil {
		b.timer.Stop()
	}
	userID := util.GetUserIDFromTgMessage(msg)
	b.timer = r.createBatchTimer(debounce, key, userID, cid, "acceptPhoto")
	b.mu.Unlock()
}

func (r *Router) acceptDocument(cid int64, msg tgbotapi.Message) {
	if msg.Document == nil {
		r.sendError(cid, fmt.Errorf("document is nil"))
		return
	}
	imgBytes, err := r.Bot.DownloadFile(msg.Document.FileID)
	if err != nil {
		r.sendError(cid, err)
		return
	}

	key := "chat:" + fmt.Sprint(cid)
	if msg.MediaGroupID != "" {
		key = "grp:" + msg.MediaGroupID
	}

	bi, _ := batches.LoadOrStore(key, &photoBatch{
		ChatID:       cid,
		Key:          key,
		MediaGroupID: msg.MediaGroupID,
		images:       make([][]byte, 0, 2),
		createdAt:    time.Now(),
	})
	b, ok := bi.(*photoBatch)
	if !ok {
		r.sendError(cid, fmt.Errorf("internal error: invalid batch type"))
		batches.Delete(key)
		return
	}

	b.mu.Lock()
	// Проверяем, не обработан ли уже batch
	if b.processed {
		b.mu.Unlock()
		// Batch уже обработан — создаём новый
		newBatch := &photoBatch{
			ChatID:       cid,
			Key:          key,
			MediaGroupID: msg.MediaGroupID,
			images:       [][]byte{imgBytes},
			createdAt:    time.Now(),
			lastAt:       time.Now(),
		}
		userID := util.GetUserIDFromTgMessage(msg)
		newBatch.timer = r.createBatchTimer(debounce, key, userID, cid, "acceptDocument.newBatch")
		batches.Store(key, newBatch)
		return
	}
	b.images = append(b.images, imgBytes)
	b.lastAt = time.Now()
	if b.timer != nil {
		b.timer.Stop()
	}
	userID := util.GetUserIDFromTgMessage(msg)
	b.timer = r.createBatchTimer(debounce, key, userID, cid, "acceptDocument")
	b.mu.Unlock()
}

func (r *Router) processBatch(key string, userID *int64) {
	// Используем context с таймаутом вместо Background
	// Также проверяем shutdown статус
	shutdown := GetShutdownManager()
	if shutdown.IsShutdown() {
		return
	}

	// Создаём context с таймаутом для всей операции (3 минуты)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	bi, ok := batches.Load(key)
	if !ok {
		return
	}
	b, ok := bi.(*photoBatch)
	if !ok {
		batches.Delete(key)
		return
	}

	// Захватываем lock и проверяем, не обработан ли уже batch
	b.mu.Lock()
	if b.processed {
		// Batch уже обработан другой горутиной
		b.mu.Unlock()
		return
	}
	// Помечаем как обработанный до разблокировки
	b.processed = true
	images := append([][]byte(nil), b.images...)
	chatID := b.ChatID
	mediaGroupID := b.MediaGroupID
	b.mu.Unlock()

	// Удаляем из map после разблокировки
	batches.Delete(key)

	if len(images) == 0 {
		return
	}

	var merged []byte
	if len(images) == 1 {
		merged = images[0]
	} else {
		var err error
		merged, err = combineAsOne(images)
		if err != nil {
			r.sendError(chatID, fmt.Errorf("склейка: %w", err))
			return
		}
	}
	// r.send(chatID, GetPhotoText, nil)
	r.runDetectThenParse(ctx, chatID, userID, merged, mediaGroupID)
}

func combineAsOne(images [][]byte) ([]byte, error) {
	decoded := make([]image.Image, 0, len(images))
	widths := make([]int, 0, len(images))
	heights := make([]int, 0, len(images))

	for _, b := range images {
		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			if try, err2 := tryDecodeStrict(b); err2 == nil {
				img = try
			} else {
				return nil, err
			}
		}
		decoded = append(decoded, img)
		bounds := img.Bounds()
		widths = append(widths, bounds.Dx())
		heights = append(heights, bounds.Dy())
	}

	maxW := 0
	sumH := 0
	for i := range decoded {
		if widths[i] > maxW {
			maxW = widths[i]
		}
		sumH += heights[i]
	}
	if maxW == 0 || sumH == 0 {
		return nil, fmt.Errorf("пустые изображения")
	}

	dst := image.NewRGBA(image.Rect(0, 0, maxW, sumH))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	y := 0
	for i, img := range decoded {
		w := widths[i]
		h := heights[i]
		x := (maxW - w) / 2
		rect := image.Rect(x, y, x+w, y+h)
		draw.Draw(dst, rect, img, img.Bounds().Min, draw.Over)
		y += h
	}

	totalPx := maxW * sumH
	final := image.Image(dst)
	if totalPx > maxPixels {
		scale := math.Sqrt(float64(maxPixels) / float64(totalPx))
		newW := int(float64(maxW)*scale + 0.5)
		newH := int(float64(sumH)*scale + 0.5)
		if newW < 1 {
			newW = 1
		}
		if newH < 1 {
			newH = 1
		}
		final = scaleDownNN(dst, newW, newH)
	}

	var out bytes.Buffer
	if err := jpeg.Encode(&out, final, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func tryDecodeStrict(b []byte) (image.Image, error) {
	if len(b) >= 2 && b[0] == 0xFF && b[1] == 0xD8 {
		return jpeg.Decode(bytes.NewReader(b))
	}
	if len(b) >= 8 &&
		b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 &&
		b[4] == 0x0D && b[5] == 0x0A && b[6] == 0x1A && b[7] == 0x0A {
		return png.Decode(bytes.NewReader(b))
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	return img, err
}

func scaleDownNN(src image.Image, newW, newH int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	sb := src.Bounds()
	srcW := sb.Dx()
	srcH := sb.Dy()
	for y := 0; y < newH; y++ {
		sy := sb.Min.Y + (y*srcH)/newH
		for x := 0; x < newW; x++ {
			sx := sb.Min.X + (x*srcW)/newW
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}
