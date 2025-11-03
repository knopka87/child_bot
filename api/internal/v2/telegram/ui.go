package telegram

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Кнопки подтверждения PARSE
func makeParseConfirmKeyboard() tgbotapi.InlineKeyboardMarkup {
	yes := tgbotapi.NewInlineKeyboardButtonData("Да", "parse_yes")
	no := tgbotapi.NewInlineKeyboardButtonData("Нет", "parse_no")
	report := tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(yes, no), tgbotapi.NewInlineKeyboardRow(report))
}

// Три кнопки действий после подсказки/парсинга
func makeActionsKeyboard(level int) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	if level < 3 {
		btnHint := tgbotapi.NewInlineKeyboardButtonData("Показать подсказку", "hint_next")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnHint))
	} else {
		btnHint := tgbotapi.NewInlineKeyboardButtonData("Похожее задание", "analogue_solution")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnHint))
	}

	btnReady := tgbotapi.NewInlineKeyboardButtonData("Готов дать решение", "ready_solution")
	btnNew := tgbotapi.NewInlineKeyboardButtonData("Перейти к новой задаче", "new_task")
	btnReport := tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnReady), tgbotapi.NewInlineKeyboardRow(btnNew), tgbotapi.NewInlineKeyboardRow(btnReport))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// лёгкое экранирование для Markdown (если функции ещё нет)
func esc(s string) string {
	s = strings.ReplaceAll(s, "`", "'")
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "*", "\\*")
	s = strings.ReplaceAll(s, "[", "\\[")
	return s
}
