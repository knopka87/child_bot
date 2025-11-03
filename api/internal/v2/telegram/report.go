package telegram

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Developer/admin chat to receive reports.
var adminsChatID = []int64{
	255509524,
	310452272,
}

type imageFile struct {
	Data []byte
	Mime string
	Name string
}

// SendSessionReport builds a ZIP report for the current task session of chatID
// and sends it to the admin report chat. The report contains human-readable
// Markdown with all steps from timeline_events and extracted images (task & answer).
func (r *Router) SendSessionReport(ctx context.Context, chatID int64) error {
	// 1) Resolve session id
	sid, ok := r.getSession(chatID)
	if !ok || sid == "" {
		return r.sendErrorToAdmin(fmt.Errorf("no active session for chatID=%d", chatID),
			"Не найден активный сеанс для пользователя.")
	}

	// 2) Load history events
	events, err := r.History.FindALLRecordsBySessionID(ctx, sid)
	if err != nil {
		return r.sendErrorToAdmin(err, "Не удалось получить историю для отчёта.")
	}
	if len(events) == 0 {
		return r.sendErrorToAdmin(fmt.Errorf("no events for session_id=%s", sid), "История пуста.")
	}

	// 3) Prepare buffers for ZIP
	now := time.Now().UTC()
	baseName := fmt.Sprintf("report_%d_%s_%s", chatID, sid, now.Format("20060102T150405Z"))
	zipPath := filepath.Join(os.TempDir(), baseName+".zip")

	out, err := os.Create(zipPath)
	if err != nil {
		return r.sendErrorToAdmin(err, "Не удалось создать ZIP для отчёта.")
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	defer zw.Close()

	// 4) Build README.md (Markdown) + collect images
	var md bytes.Buffer
	_, _ = fmt.Fprintf(&md, "# Отчёт по сессии\n\n")
	_, _ = fmt.Fprintf(&md, "- Пользователь (chatID): **%d**\n", chatID)
	_, _ = fmt.Fprintf(&md, "- Session ID: **%s**\n", sid)
	_, _ = fmt.Fprintf(&md, "- Временная зона сервера (UTC): **%s**\n", now.Format(time.RFC3339))
	_, _ = fmt.Fprintf(&md, "- Количество шагов: **%d**\n\n", len(events))

	_, _ = fmt.Fprintf(&md, "## Шаги\n\n")

	// Collect first task/answer images we can find for convenience

	var taskImages []imageFile
	var answerImages []imageFile

	for i, ev := range events {
		idx := i + 1
		_, _ = fmt.Fprintf(&md, "### %02d. %s | %s | провайдер: %s | ok=%v\n\n",
			idx, strings.ToUpper(ev.Direction), ev.EventType, ev.Provider, ev.OK)

		// Timestamp (CreatedAt might be zero if not set in insert)
		if !ev.CreatedAt.IsZero() {
			_, _ = fmt.Fprintf(&md, "- Время: %s\n", ev.CreatedAt.Format(time.RFC3339))
		}
		if ev.LatencyMS != nil {
			_, _ = fmt.Fprintf(&md, "- Время ответа: %d ms\n", *ev.LatencyMS)
		}
		if ev.TgMessageID != nil {
			_, _ = fmt.Fprintf(&md, "- Telegram message id: %d\n", *ev.TgMessageID)
		}
		if strings.TrimSpace(ev.Text) != "" {
			_, _ = fmt.Fprintf(&md, "- Текст: %s\n", sanitizeForMD(ev.Text))
		}
		_, _ = fmt.Fprintln(&md, "")

		// Pretty print payloads
		if ev.InputPayload != nil {
			j, _ := json.MarshalIndent(ev.InputPayload, "", "  ")
			_, _ = fmt.Fprintf(&md, "<details><summary>Input</summary>\n\n```json\n%s\n```\n\n</details>\n\n", string(j))

			// Try to extract images from input
			imgs := extractImagesGeneric(ev.InputPayload)
			classifyAndAppendImages(ev.EventType, imgs, &taskImages, &answerImages)
		}
		if ev.OutputPayload != nil {
			j, _ := json.MarshalIndent(ev.OutputPayload, "", "  ")
			_, _ = fmt.Fprintf(&md, "<details><summary>Output</summary>\n\n```json\n%s\n```\n\n</details>\n\n", string(j))

			// Try to extract images from output (just in case)
			imgs := extractImagesGeneric(ev.OutputPayload)
			classifyAndAppendImages(ev.EventType, imgs, &taskImages, &answerImages)
		}

		if ev.Error != nil {
			_, _ = fmt.Fprintf(&md, "> Ошибка: %s\n\n", sanitizeForMD(ev.Error.Error()))
		}
	}

	// 5) Add README.md to ZIP
	if err := writeZipFile(zw, "Report.md", md.Bytes()); err != nil {
		return r.sendErrorToAdmin(err, "Не удалось записать Report.md в ZIP.")
	}

	// 6) Add images to ZIP
	// Save up to a few images to keep report small
	for i, img := range taskImages {
		if i >= 3 {
			break
		}
		name := img.Name
		if name == "" {
			name = fmt.Sprintf("images/task_%d%s", i+1, mimeToExt(img.Mime))
		} else {
			// ensure under images/
			name = filepath.Join("images", filepath.Base(name))
		}
		if err := writeZipFile(zw, name, img.Data); err != nil {
			return r.sendErrorToAdmin(err, "Не удалось записать изображение задания в ZIP.")
		}
	}
	for i, img := range answerImages {
		if i >= 3 {
			break
		}
		name := img.Name
		if name == "" {
			name = fmt.Sprintf("images/answer_%d%s", i+1, mimeToExt(img.Mime))
		} else {
			name = filepath.Join("images", filepath.Base(name))
		}
		if err := writeZipFile(zw, name, img.Data); err != nil {
			return r.sendErrorToAdmin(err, "Не удалось записать изображение ответа в ZIP.")
		}
	}

	// 7) Close ZIP
	if err := zw.Close(); err != nil {
		return r.sendErrorToAdmin(err, "Не удалось закрыть ZIP.")
	}
	_ = out.Close()

	// 8) Send to admin
	doc := tgbotapi.FilePath(zipPath)
	for _, cid := range adminsChatID {
		msg := tgbotapi.NewDocument(cid, doc)
		msg.Caption = fmt.Sprintf("Отчёт по сессии\nchatID=%d\nsession=%s\nsteps=%d", chatID, sid, len(events))
		if _, err := r.Bot.Send(msg); err != nil {
			return r.sendErrorToAdmin(err, "Не удалось отправить отчёт в Telegram.")
		}
	}

	// Optional: notify user that report has been sent
	r.send(chatID, "Отчёт отправлен разработчику. Спасибо!", nil)

	// Cleanup temp file
	_ = os.Remove(zipPath)
	return nil
}

// sendErrorToAdmin sends an error message to admin and returns the same error.
func (r *Router) sendErrorToAdmin(err error, userMsg string) error {
	for _, cid := range adminsChatID {
		adm := tgbotapi.NewMessage(cid, fmt.Sprintf("report error: %v\n\n%s", err, userMsg))
		_, _ = r.Bot.Send(adm)
	}
	return err
}

// writeZipFile adds a new file to the zip writer.
func writeZipFile(zw *zip.Writer, name string, data []byte) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// classifyAndAppendImages heuristically separates task vs answer images by event type.
func classifyAndAppendImages(eventType string, imgs []imageInfo, tasks *[]imageFile, answers *[]imageFile) {
	l := strings.ToLower(eventType)
	for _, im := range imgs {
		img := imageFile{Data: im.Data, Mime: im.Mime, Name: im.Name}
		switch {
		case strings.Contains(l, "detect"), strings.Contains(l, "parse"):
			*tasks = append(*tasks, img)
		case strings.Contains(l, "normalize"), strings.Contains(l, "solution"):
			*answers = append(*answers, img)
		default:
			// fallback: classify by key hint
			if strings.Contains(strings.ToLower(im.Name), "answer") {
				*answers = append(*answers, img)
			} else {
				*tasks = append(*tasks, img)
			}
		}
	}
}

// imageInfo is an internal helper for extracted image data.
type imageInfo struct {
	Data []byte
	Mime string
	Name string
}

// extractImagesGeneric tries to find any base64 images in arbitrary JSON-like structures.
// It looks for common keys: photo_b64, image_b64, images[].b64, content (data URL), and carries mime when available.
func extractImagesGeneric(payload any) []imageInfo {
	var out []imageInfo

	var walk func(k string, v any)
	walk = func(k string, v any) {
		switch t := v.(type) {
		case map[string]any:
			// Try to pick up mime/name hints early
			curMime := getString(t, "mime")
			curName := getString(t, "name")

			// Primary fields
			if s := getString(t, "photo_b64"); s != "" {
				if b, mimeStr := decodeMaybeDataURL(s, curMime); len(b) > 0 {
					out = append(out, imageInfo{Data: b, Mime: mimeStr, Name: curName})
				}
			}
			if s := getString(t, "image"); s != "" {
				if b, mimeStr := decodeMaybeDataURL(s, curMime); len(b) > 0 {
					out = append(out, imageInfo{Data: b, Mime: mimeStr, Name: curName})
				}
			}
			if s := getString(t, "image_b64"); s != "" {
				if b, mimeStr := decodeMaybeDataURL(s, curMime); len(b) > 0 {
					out = append(out, imageInfo{Data: b, Mime: mimeStr, Name: curName})
				}
			}
			// Some providers may use "content": "data:image/jpeg;base64,...."
			if s := getString(t, "content"); s != "" {
				if b, mimeStr := decodeMaybeDataURL(s, curMime); len(b) > 0 {
					out = append(out, imageInfo{Data: b, Mime: mimeStr, Name: curName})
				}
			}
			// Nested structures
			for kk, vv := range t {
				// Avoid infinite recursion on already processed scalar keys
				if kk == "photo_b64" || kk == "image_b64" || kk == "image" || kk == "content" {
					continue
				}
				walk(kk, vv)
			}
		case []any:
			for _, it := range t {
				walk(k, it)
			}
		}
	}

	walk("", payload)
	return out
}

var dataURLRe = regexp.MustCompile(`^data:([^;]+);base64,(.*)$`)

func decodeMaybeDataURL(s string, fallbackMime string) ([]byte, string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fallbackMime
	}
	if m := dataURLRe.FindStringSubmatch(s); len(m) == 3 {
		mimeType := m[1]
		raw := m[2]
		b, _ := base64.StdEncoding.DecodeString(raw)
		return b, mimeType
	}
	// Not a data URL — try plain base64
	b, err := base64.StdEncoding.DecodeString(s)
	if err == nil && len(b) > 0 {
		return b, fallbackMime
	}
	// Some payloads may be URL-encoded
	if u, err := url.QueryUnescape(s); err == nil {
		if bb, err2 := base64.StdEncoding.DecodeString(u); err2 == nil && len(bb) > 0 {
			return bb, fallbackMime
		}
	}
	return nil, fallbackMime
}

func mimeToExt(m string) string {
	if m == "" {
		return ".jpg"
	}
	exts, _ := mime.ExtensionsByType(m)
	if len(exts) > 0 {
		return exts[0]
	}
	switch strings.ToLower(m) {
	case "image/jpeg", "jpeg", "jpg":
		return ".jpg"
	case "image/png", "png":
		return ".png"
	case "image/heic", "heic":
		return ".heic"
	case "image/avif", "avif":
		return ".avif"
	default:
		return ".bin"
	}
}

func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// sanitizeForMD escapes Markdown special characters for safe display.
func sanitizeForMD(s string) string {
	// Only basic escaping for code/quotes/markdown
	s = strings.ReplaceAll(s, "`", "'")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
