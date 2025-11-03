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
func makeActionsKeyboardRow(level int, showAnalogue bool) [][]tgbotapi.InlineKeyboardButton {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 4)
	if level < 3 {
		btnHint := tgbotapi.NewInlineKeyboardButtonData("Показать подсказку", "hint_next")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnHint))
	} else if showAnalogue {
		btnAnalogue := tgbotapi.NewInlineKeyboardButtonData("Похожее задание", "analogue_solution")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnAnalogue))
	}

	btnReady := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Готов дать решение", "ready_solution"))
	btnNew := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Перейти к новой задаче", "new_task"))
	btnReport := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
	rows = append(rows, btnReady, btnNew, btnReport)

	return rows
}

// лёгкое экранирование для Markdown (если функции ещё нет)
func esc(s string) string {
	s = strings.ReplaceAll(s, "`", "'")
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "*", "\\*")
	s = strings.ReplaceAll(s, "[", "\\[")
	return s
}
