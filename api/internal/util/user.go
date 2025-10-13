package util

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
