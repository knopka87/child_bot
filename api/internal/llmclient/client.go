package llmclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
		hc:   &http.Client{Timeout: 60 * time.Second},
	}
}

// Detect отправляет СЫРЫЕ данные 1-в-1 как входы метода Engine.Detect
// и возвращает структуру DetectResult, которую должен вернуть LLM-сервис.
func (c *Client) Detect(ctx context.Context, llmName string, img []byte, mime string, gradeHint int) (ocr.DetectResult, error) {
	in := detectRequest{
		LLMName:   llmName,
		ImageB64:  base64.StdEncoding.EncodeToString(img),
		Mime:      mime,
		GradeHint: gradeHint,
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
		LLMName:  llmName,
		ImageB64: base64.StdEncoding.EncodeToString(image),
		Options:  options,
	}
	var out ocr.ParseResult
	if err := c.post(ctx, "/v1/parse", in, &out); err != nil {
		return ocr.ParseResult{}, err
	}
	return out, nil
}

// Hint отправляет СЫРЫЙ ocr.HintInput и ожидает ocr.HintResult.
func (c *Client) Hint(ctx context.Context, in ocr.HintInput) (ocr.HintResult, error) {
	var out ocr.HintResult
	if err := c.post(ctx, "/v1/hint", in, &out); err != nil {
		return ocr.HintResult{}, err
	}
	return out, nil
}

// --- внутренности ------------------------------------------------------------

type detectRequest struct {
	LLMName   string `json:"llm_name"`
	ImageB64  string `json:"image_b64"`
	Mime      string `json:"mime,omitempty"`
	GradeHint int    `json:"grade_hint,omitempty"`
}

type parseRequest struct {
	LLMName  string           `json:"llm_name"`
	ImageB64 string           `json:"image_b64"`
	Options  ocr.ParseOptions `json:"options"`
}

func (c *Client) post(ctx context.Context, path string, body interface{}, out interface{}) error {
	buf, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	res, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		// Попробуем вытащить текст ошибки
		var e struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(res.Body).Decode(&e)
		if strings.TrimSpace(e.Error) == "" {
			return fmt.Errorf("llm server http %d", res.StatusCode)
		}
		return errors.New(e.Error)
	}
	return json.NewDecoder(res.Body).Decode(out)
}
