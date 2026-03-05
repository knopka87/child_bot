package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	// Визуальный разделитель
	Divider = "━━━━━━━━━━━━━━━━━━━━"

	OkText                  = "✅ OK"
	UnderFoundCommandText   = "Неизвестная команда. Я знаю только команду /start"
	StartMessageText        = "🎒 Поехали!\nОтправь фото задания — разберёмся вместе 🤓"
	NewTaskText             = "🎒 Поехали!\nОтправь фото задания — разберёмся вместе 🤓"
	ReadTaskText            = "🧠 Читаю твоё задание…"
	ParseProgress1          = "📷 Получил фото"
	ParseProgress2          = "🔍 Анализирую задание…"
	ParseProgress3          = "🧠 Распознаю текст…"
	ParseProgress4          = "✨ Почти готово…"
	HintProgress1           = "🤔 Думаю над подсказкой…"
	HintProgress2           = "💡 Подбираю объяснение…"
	HintProgress3           = "✨ Формулирую…"
	CheckProgress1          = "👀 Смотрю решение…"
	CheckProgress2          = "🔍 Проверяю шаги…"
	CheckProgress3          = "🧮 Сверяю ответ…"
	CheckProgress4          = "✨ Почти готово…"
	TaskViewText            = "📸 *Вот что я прочитал:*\n" + Divider + "\n\n%s\n\n" + Divider + "\n\nВсё верно? 🤔"
	HintNotFoundText        = "🤔 Подсказки недоступны: сначала пришлите фото задания 📸"
	HintFinishText          = "🤔 Все подсказки уже показаны. Могу показать аналогичную задачу 🧩"
	HINT1Text               = "✨ *Подсказка 1*\n" + Divider + "\n\n%s\n\n" + Divider + "\n\nПопробуй выполнить сам 😉"
	HINT2Text               = "💡 *Подсказка 2*\n" + Divider + "\n\n%s\n\n" + Divider + "\n\nПопробуй выполнить задание 😉"
	HINT3Text               = "🌟 *Подсказка 3*\n" + Divider + "\n\n%s\n\n" + Divider + "\n\nПопробуй применить это к своему заданию 🚀"
	AnalogueTaskWaitingText = "🎯 Молодец!\n⏳ Думаю и подбираю похожее задание, чтобы объяснить на примере…"
	AnalogueAlert1          = "🔍 Подбираю пример…"
	AnalogueAlert2          = "🤔 Ищу лучший вариант…"
	AnalogueAlert3          = "💡 Выбираю похожее задание…"
	AnalogueTaskText        = "🧩 *Похожее задание*\n" + Divider + "\n\n%s\n\n" + Divider + "\n\nПопробуй вернуться к заданию и решить его 💪"
	CheckAnswerClick        = "🔎 Проверим твой ответ? ✨\n📸 Пришли фото своего решения — я посмотрю, всё ли верно 😊\nНо я ещё учусь проверять, поэтому могу ошибаться."
	CheckAnswerText         = "🤓 Вижу твоё решение!"
	CheckAlert              = "🤔 Проверяю каждый шаг…"
	AnswerCorrectText       = "🎉 *ВЕРНО!*\n" + Divider + "\n\nОтличная работа! Ты понял, как это работает 💪🌟"
	AnswerIncorrectText     = "💪 *Почти получилось!*\n" + Divider + "\n\n%s\n\n" + Divider
	ReportText              = "Отлично, спасибо, что заметил! 📝\nЧтобы я стал лучше, напиши, в чем ошибка?\nНапример:\n• я неправильно прочитал часть задания;\n• подсказка не помогла или была непонятной;\n• я объяснил не то задание;\n• другое (опиши своими словами).\n💬 Напиши коротко, что именно было не так.\n👉 После сообщения мы продолжим разбор задания."
	SendReportText          = "👋 Спасибо, что помогаешь мне стать лучше!\n🎒Давай продолжим. Отправь фото задания — разберёмся вместе 🤓"
	DontLikeHint            = "😌 Понял тебя!\nСпасибо за отзыв — попробую объяснить по-другому 💡"
	ErrorText               = "😅 Ой, что-то пошло не так… Уже чиню 🔧\nПопробуй чуть позже 📝"
	DetectErrorText         = "😥 Не удалось обработать фото."
	SubjectNotSupportedText = "📚 Пока я умею помогать только с математикой.\nНо я учусь! Скоро смогу помочь и с другими предметами 🚀\n\n🎒 Отправь фото задания по математике — разберёмся вместе 🤓"
	GradePreviewText        = "Чтобы я мог давать подсказки подходящего уровня, выбери свой класс 🧩"
	AwaitSolutionText       = "📸 Пришли фото твоего решения — и я посмотрю, всё ли правильно 😊"
	AwaitNewTaskText        = "📸 Скидывай своё задание — и разберёмся вместе! 🤓"
	StepSolutionText        = "\n\n\n\n📘 Шаги решения\n\n"

	YesButton          = "✅ Да, направь подсказку"
	CheckAnswerButton  = "🔎 Проверь мой ответ"
	SendReportButton   = "📝 Сообщить об ошибке"
	NextHintButton     = "➡️ Следующая подсказка"
	DontLikeHintButton = "👎 Не нравится подсказка"
	NewTaskButton      = "🆕 Новое задание"
	AnalogueTaskButton = "🧩 Похожее задание с решением"
	Grade1Button       = "📕 1 класс"
	Grade2Button       = "📗 2 класс"
	Grade3Button       = "📘 3 класс"
	Grade4Button       = "📙 4 класс"
)

var (
	btnYes           = tgbotapi.NewInlineKeyboardButtonData(YesButton, "parse_yes")
	btnCheckAnswer   = tgbotapi.NewInlineKeyboardButtonData(CheckAnswerButton, "ready_solution")
	btnNextHint      = tgbotapi.NewInlineKeyboardButtonData(NextHintButton, "hint_next")
	btnReport        = tgbotapi.NewInlineKeyboardButtonData(SendReportButton, "report")
	btnDontLikeHint  = tgbotapi.NewInlineKeyboardButtonData(DontLikeHintButton, "dont_like_hint")
	btnNewTask       = tgbotapi.NewInlineKeyboardButtonData(NewTaskButton, "new_task")
	btnReadySolution = tgbotapi.NewInlineKeyboardButtonData(CheckAnswerButton, "ready_solution")
	btnAnalogue      = tgbotapi.NewInlineKeyboardButtonData(AnalogueTaskButton, "analogue_task")
	btnGrade1        = tgbotapi.NewInlineKeyboardButtonData(Grade1Button, "grade1")
	btnGrade2        = tgbotapi.NewInlineKeyboardButtonData(Grade2Button, "grade2")
	btnGrade3        = tgbotapi.NewInlineKeyboardButtonData(Grade3Button, "grade3")
	btnGrade4        = tgbotapi.NewInlineKeyboardButtonData(Grade4Button, "grade4")
)

func makeGradeListButtons() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(btnGrade1, btnGrade2),
		tgbotapi.NewInlineKeyboardRow(btnGrade3, btnGrade4),
	}
}

func makeErrorButtons() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(btnNewTask),
	}
}

// Кнопки подтверждения PARSE
func makeParseConfirmButtons() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(btnYes),
		tgbotapi.NewInlineKeyboardRow(btnCheckAnswer),
		tgbotapi.NewInlineKeyboardRow(btnReport),
	}
}

func makeFinishHintButtons() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(btnCheckAnswer),
		tgbotapi.NewInlineKeyboardRow(btnNewTask),
		tgbotapi.NewInlineKeyboardRow(btnReport),
	}
}

func makeHintButtons(level, maxHints int, showAnalogue bool) [][]tgbotapi.InlineKeyboardButton {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 4)
	if level < maxHints {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnNextHint))
	} else if showAnalogue {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnAnalogue))
	}

	btnReady := tgbotapi.NewInlineKeyboardRow(btnReadySolution)
	btnDontLike := tgbotapi.NewInlineKeyboardRow(btnDontLikeHint)
	btnNew := tgbotapi.NewInlineKeyboardRow(btnNewTask)
	rows = append(rows, btnReady, btnDontLike, btnNew)

	return rows
}

func makeAnalogueButtons() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(btnCheckAnswer),
		tgbotapi.NewInlineKeyboardRow(btnNewTask),
		tgbotapi.NewInlineKeyboardRow(btnReport),
	}
}

func makeCheckAnswerClickButtons() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(btnNewTask),
	}
}

func makeCorrectAnswerButtons() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(btnNewTask),
		tgbotapi.NewInlineKeyboardRow(btnReport),
	}
}

func makeIncorrectAnswerButtons() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(btnAnalogue),
		tgbotapi.NewInlineKeyboardRow(btnNewTask),
		tgbotapi.NewInlineKeyboardRow(btnReport),
	}
}
