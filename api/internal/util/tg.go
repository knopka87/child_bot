package util

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetUserIDFromTgUpdate(up tgbotapi.Update) *int64 {
	if up.Message != nil {
		userID := GetUserIDFromTgMessage(*up.Message)
		if userID != nil {
			return userID
		}
	}
	if up.CallbackQuery != nil {
		userID := GetUserIDFromTgCB(*up.CallbackQuery)
		if userID != nil {
			return userID
		}
	}

	return nil
}

func GetUserIDFromTgMessage(m tgbotapi.Message) *int64 {
	if m.Contact != nil {
		return &m.Contact.UserID
	}
	if m.From != nil {
		return &m.From.ID
	}
	return nil
}

func GetUserIDFromTgCB(cb tgbotapi.CallbackQuery) *int64 {
	if cb.Message != nil {
		if cb.Message.Contact != nil {
			return &cb.Message.Contact.UserID
		}
		if cb.Message.From != nil {
			return &cb.Message.From.ID
		}
	}
	if cb.From != nil {
		return &cb.From.ID
	}

	return nil
}

func GetChatIDByTgUpdate(up tgbotapi.Update) int64 {
	if up.Message != nil {
		chatID := GetChatIDFromTgMessage(*up.Message)
		if chatID > 0 {
			return chatID
		}
	}
	if up.CallbackQuery != nil {
		chatID := GetChatIDFromTgCB(*up.CallbackQuery)
		if chatID > 0 {
			return chatID
		}
	}
	PrintInfo("GetChatIDByTgUpdate", "", 0, fmt.Sprintf("not found chatID; up: %v", up))
	return 0
}

func GetChatIDFromTgCB(c tgbotapi.CallbackQuery) int64 {
	if c.Message != nil {
		if c.Message.Chat != nil {
			return c.Message.Chat.ID
		}
	}

	PrintInfo("GetChatIDFromTgCB", "", 0, fmt.Sprintf("not found chatID; c: %v", c))
	return 0
}

func GetChatIDFromTgMessage(m tgbotapi.Message) int64 {
	if m.Chat != nil {
		return m.Chat.ID
	}

	PrintInfo("GetChatIDFromTgMessage", "", 0, fmt.Sprintf("not found chatID; m: %v", m))
	return 0
}
