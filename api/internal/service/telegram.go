package service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type TgRouter interface {
	GetToken() string
	HandleUpdate(upd tgbotapi.Update, llmName string)
}
