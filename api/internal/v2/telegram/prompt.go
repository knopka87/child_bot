package telegram

import (
	"context"
	"fmt"

	"child-bot/api/internal/v2/types"
)

// postUpdatePrompt sends UpdatePromptRequest to llm-proxy /api/prompt and reports the result back to the chat.
func (r *Router) postUpdatePrompt(ctx context.Context, chatID int64, name, text string) {
	provider := r.LlmManager.Get(chatID)

	// Build request payload
	reqBody := types.UpdatePromptRequest{
		Provider: provider,
		Name:     name,
		Text:     text,
	}

	out, err := r.GetLLMClient().UpdatePrompt(ctx, reqBody)
	if err != nil {
		r.sendDebug(chatID, "update prompt", err)
	}

	if err != nil {
		// Ответ пришёл с ошибкой
		r.send(chatID, fmt.Sprintf("Не удалось обновить промпт '%s' для провайдера '%s': %v", reqBody.Name, reqBody.Provider, err), nil)
		return
	}
	if !out.OK {
		// Ответ пришёл, но ок == false — покажем пользователю
		r.send(chatID, fmt.Sprintf("Не удалось обновить промпт '%s' для провайдера '%s' (путь: %s)", out.Name, out.Provider, out.Path), nil)
		return
	}

	// Успех
	msg := fmt.Sprintf("✅ Промпт обновлён.\nПровайдер: %s\nИмя: %s\nФайл: %s\nРазмер: %d байт\nОбновлён: %s", out.Provider, out.Name, out.Path, out.Size, out.Updated)
	r.send(chatID, msg, nil)
}
