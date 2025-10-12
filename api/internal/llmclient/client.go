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
