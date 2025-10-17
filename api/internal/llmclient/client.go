package llmclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	ocr "child-bot/api/internal/ocr"
)

type Client struct {
	base string
	hc   *http.Client
}

func New(base string) *Client {
	base = strings.TrimRight(base, "/")
	return &Client{
		base: base,
		hc: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				ResponseHeaderTimeout: 120 * time.Second, // ждём заголовки до 2 минут
			},
			Timeout: 0, // общий таймаут управляем per-request через ctx
		},
	}
}

// Detect отправляет СЫРЫЕ данные 1-в-1 как входы метода Engine.Detect
// и возвращает структуру DetectResult, которую должен вернуть LLM-сервис.
func (c *Client) Detect(ctx context.Context, llmName string, img []byte, mime string, gradeHint int) (ocr.DetectResult, error) {
	in := detectRequest{
		LLMName: llmName,
		DetectInput: ocr.DetectInput{
			ImageB64:  base64.StdEncoding.EncodeToString(img),
			Mime:      mime,
			GradeHint: gradeHint,
		},
	}
	var out ocr.DetectResult
	if err := c.post(ctx, "/v1/detect", in, &out); err != nil {
		return ocr.DetectResult{}, err
	}
	return out, nil
}

// Parse отправляет СЫРЫЕ данные 1-в-1 как входы метода Engine.Parse
// image передаётся как base64, остальные параметры — внутри options.
func (c *Client) Parse(ctx context.Context, llmName string, image []byte, options ocr.ParseOptions) (ocr.ParseResult, error) {
	in := parseRequest{
		LLMName: llmName,
		ParseInput: ocr.ParseInput{
			ImageB64: base64.StdEncoding.EncodeToString(image),
			Options:  options,
		},
	}
	var out ocr.ParseResult
	if err := c.post(ctx, "/v1/parse", in, &out); err != nil {
		return ocr.ParseResult{}, err
	}
	return out, nil
}

// Hint отправляет СЫРЫЙ ocr.HintInput и ожидает ocr.HintResult.
func (c *Client) Hint(ctx context.Context, llmName string, hin ocr.HintInput) (ocr.HintResult, error) {
	in := hintRequest{
		LLMName:   llmName,
		HintInput: hin,
	}
	var out ocr.HintResult
	if err := c.post(ctx, "/v1/hint", in, &out); err != nil {
		return ocr.HintResult{}, err
	}
	return out, nil
}

func (c *Client) Normalize(ctx context.Context, llmName string, nin ocr.NormalizeInput) (ocr.NormalizeResult, error) {
	in := normalizeRequest{
		LLMName:        llmName,
		NormalizeInput: nin,
	}
	var out ocr.NormalizeResult
	if err := c.post(ctx, "/v1/normalize", in, &out); err != nil {
		return ocr.NormalizeResult{}, err
	}
	return out, nil
}

func (c *Client) CheckSolution(ctx context.Context, llmName string, cin ocr.CheckSolutionInput) (ocr.CheckSolutionResult, error) {
	in := checkSolutionRequest{
		LLMName:            llmName,
		CheckSolutionInput: cin,
	}
	var out ocr.CheckSolutionResult
	if err := c.post(ctx, "/v1/check_solution", in, &out); err != nil {
		return ocr.CheckSolutionResult{}, err
	}

	return out, nil
}

func (c *Client) AnalogueSolution(ctx context.Context, llmName string, ain ocr.AnalogueSolutionInput) (ocr.AnalogueSolutionResult, error) {
	in := analogueSolutionRequest{
		LLMName:               llmName,
		AnalogueSolutionInput: ain,
	}
	var out ocr.AnalogueSolutionResult
	if err := c.post(ctx, "/v1/analogue_solution", in, &out); err != nil {
		return ocr.AnalogueSolutionResult{}, err
	}

	return out, nil
}

// --- внутренности ------------------------------------------------------------

type detectRequest struct {
	LLMName string `json:"llm_name"`
	ocr.DetectInput
}

type parseRequest struct {
	LLMName string `json:"llm_name"`
	ocr.ParseInput
}

type hintRequest struct {
	LLMName string `json:"llm_name"`
	ocr.HintInput
}

type normalizeRequest struct {
	LLMName string `json:"llm_name"`
	ocr.NormalizeInput
}

type checkSolutionRequest struct {
	LLMName string `json:"llm_name"`
	ocr.CheckSolutionInput
}

type analogueSolutionRequest struct {
	LLMName string `json:"llm_name"`
	ocr.AnalogueSolutionInput
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
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.base+pathWithTimeout, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if timeoutSec > 0 {
		// Дружелюбный заголовок — сервер может читать либо header, либо query (?timeoutSec=)
		req.Header.Set("X-Request-Timeout", fmt.Sprintf("%d", timeoutSec))
	}
	res, err := c.hc.Do(req)
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
