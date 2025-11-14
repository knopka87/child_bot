package llmclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"child-bot/api/internal/llmclient"
	"child-bot/api/internal/v1/types"
)

type Client struct {
	client *llmclient.Client
}

func New(client *llmclient.Client) *Client {
	return &Client{client: client}
}

// Detect отправляет СЫРЫЕ данные 1-в-1 как входы метода Engine.Detect
// и возвращает структуру DetectResponse, которую должен вернуть LLMClient-сервис.
func (c *Client) Detect(ctx context.Context, llmName string, din types.DetectRequest) (types.DetectResponse, error) {
	in := detectRequest{
		LLMName:       llmName,
		DetectRequest: din,
	}
	var out types.DetectResponse
	if err := c.post(ctx, "/v1/detect", in, &out); err != nil {
		return types.DetectResponse{}, err
	}
	return out, nil
}

// Parse отправляет СЫРЫЕ данные 1-в-1 как входы метода Engine.Parse
// image передаётся как base64, остальные параметры — внутри options.
func (c *Client) Parse(ctx context.Context, llmName string, pin types.ParseRequest) (types.ParseResponse, error) {
	in := parseRequest{
		LLMName:      llmName,
		ParseRequest: pin,
	}
	var out types.ParseResponse
	if err := c.post(ctx, "/v1/parse", in, &out); err != nil {
		return types.ParseResponse{}, err
	}
	return out, nil
}

// Hint отправляет СЫРЫЙ ocr.HintRequest и ожидает ocr.HintResponse.
func (c *Client) Hint(ctx context.Context, llmName string, hin types.HintRequest) (types.HintResponse, error) {
	in := hintRequest{
		LLMName:     llmName,
		HintRequest: hin,
	}
	var out types.HintResponse
	if err := c.post(ctx, "/v1/hint", in, &out); err != nil {
		return types.HintResponse{}, err
	}
	return out, nil
}

func (c *Client) OCR(ctx context.Context, llmName string, oin types.OCRRequest) (types.OCRResponse, error) {
	in := ocrRequest{
		LLMName:    llmName,
		OCRRequest: oin,
	}
	var out types.OCRResponse
	if err := c.post(ctx, "/v1/ocr", in, &out); err != nil {
		return types.OCRResponse{}, err
	}
	return out, nil
}

func (c *Client) Normalize(ctx context.Context, llmName string, nin types.NormalizeRequest) (types.NormalizeResponse, error) {
	in := normalizeRequest{
		LLMName:          llmName,
		NormalizeRequest: nin,
	}
	var out types.NormalizeResponse
	if err := c.post(ctx, "/v1/normalize", in, &out); err != nil {
		return types.NormalizeResponse{}, err
	}
	return out, nil
}

func (c *Client) CheckSolution(ctx context.Context, llmName string, cin types.CheckRequest) (types.CheckResponse, error) {
	in := checkRequest{
		LLMName:      llmName,
		CheckRequest: cin,
	}
	var out types.CheckResponse
	if err := c.post(ctx, "/v1/check_solution", in, &out); err != nil {
		return types.CheckResponse{}, err
	}

	return out, nil
}

func (c *Client) AnalogueSolution(ctx context.Context, llmName string, ain types.AnalogueRequest) (types.AnalogueResponse, error) {
	in := analogueRequest{
		LLMName:         llmName,
		AnalogueRequest: ain,
	}
	var out types.AnalogueResponse
	if err := c.post(ctx, "/v1/analogue_solution", in, &out); err != nil {
		return types.AnalogueResponse{}, err
	}

	return out, nil
}

// --- внутренности ------------------------------------------------------------

type detectRequest struct {
	LLMName string `json:"llm_name"`
	types.DetectRequest
}

type parseRequest struct {
	LLMName string `json:"llm_name"`
	types.ParseRequest
}

type hintRequest struct {
	LLMName string `json:"llm_name"`
	types.HintRequest
}

type ocrRequest struct {
	LLMName string `json:"llm_name"`
	types.OCRRequest
}

type normalizeRequest struct {
	LLMName string `json:"llm_name"`
	types.NormalizeRequest
}

type checkRequest struct {
	LLMName string `json:"llm_name"`
	types.CheckRequest
}

type analogueRequest struct {
	LLMName string `json:"llm_name"`
	types.AnalogueRequest
}

// addTimeoutSec appends ?timeoutSec=N (or &timeoutSec=N) to the given path.
func addTimeoutSec(path string, seconds int) string {
	if seconds <= 0 {
		return path
	}
	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	return path + sep + "timeoutSec=" + fmt.Sprintf("%d", seconds)
}

func (c *Client) post(ctx context.Context, path string, body interface{}, out interface{}) error {
	// Установим per-request timeout, если его ещё нет
	const defaultTotalTimeout = 3 * time.Minute
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultTotalTimeout)
		defer cancel()
	}

	// Вычисляем оставшееся время для передачи downstream
	var timeoutSec int
	if dl, ok := ctx.Deadline(); ok {
		rem := time.Until(dl)
		if rem > 0 {
			timeoutSec = int(rem.Seconds())
		}
	}
	pathWithTimeout := addTimeoutSec(path, timeoutSec)

	buf, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.client.Base+pathWithTimeout, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if timeoutSec > 0 {
		// Дружелюбный заголовок — сервер может читать либо header, либо query (?timeoutSec=)
		req.Header.Set("X-Request-Timeout", fmt.Sprintf("%d", timeoutSec))
	}
	res, err := c.client.HC.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		// Пробуем аккуратно извлечь текст ошибки: JSON (несколько форматов) или простой текст
		b, _ := io.ReadAll(res.Body)

		// 1) Попытка распарсить как простой {"error": "..."} или {"message": "..."}
		var e1 struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(b, &e1); err == nil {
			if msg := strings.TrimSpace(e1.Error); msg != "" {
				return errors.New(msg)
			}
			if msg := strings.TrimSpace(e1.Message); msg != "" {
				return errors.New(msg)
			}
		}

		// 2) Попытка nested-формата: {"error": {"message": "..."}}
		var e2 struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(b, &e2); err == nil {
			if msg := strings.TrimSpace(e2.Error.Message); msg != "" {
				return errors.New(msg)
			}
		}

		// 3) Фоллбэк: использовать тело как простой текст
		if msg := strings.TrimSpace(string(b)); msg != "" {
			return errors.New(msg)
		}

		// 4) Совсем ничего не удалось вытащить — вернуть код HTTP
		return fmt.Errorf("llm server http %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(out)
}
