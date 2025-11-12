package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	OkText                  = "✅ OK"
	UnderFoundCommandText   = "Неизвестная команда. Я знаю только команду /start"
	StartMessageText        = "👋 Ура, мы начинаем!\n\n\nПогнали! 🎒\nСкидывай своё задание — и разберёмся вместе! 🤓"
	NewTaskText             = "Погнали! 🎒\nСкидывай своё задание — и разберёмся вместе! 🤓"
	GetPhotoText            = "📸 Отлично, я получил твоё задание!"
	ReadTaskText            = "🧠 Читаю твоё задание…"
	TaskViewText            = "Проверь меня, вот что я понял с картинки:\n%s\nПравильно я прочитал? 🤔 Что делаем дальше?"
	HintNotFoundText        = "🤔 Подсказки недоступны: сначала пришлите фото задания 📸"
	HintFinishText          = "🤔 Все подсказки уже показаны. Могу показать аналогичную задачу 🧩"
	HINT1Text               = "✨ Отлично! Тогда начинаем разбираться вместе 🧩\nВот первая подсказка:\n%s\nПопробуй подумать и решить сам 😉\nЕсли нужно — я дам следующую подсказку!"
	HINT2Text               = "💡 Отлично! Давай попробуем посмотреть с другой стороны 👀\nВот вторая подсказка:\n%s\nТы молодец, что не сдаешься 💪\nПопробуй решить — у тебя точно получится!"
	HINT3Text               = "🌟 Отлично, мы почти разобрались!\n🤔 Думаю, как объяснить это самым понятным способом...\nВот третья подсказка — она поможет тебе окончательно понять задание:\n%s\nПодумай, что получится, если применить это к твоему заданию 😉\nТы уже у цели! 🚀"
	AnalogueTaskWaitingText = "🎯 Ты молодец, давай разберём похожее задание!\n\n\n⏳ Придумываю похожее задание, чтобы объяснить на примере."
	AnalogueTaskText        = "%s\n\nПопробуй вернуться к заданию и решить его. 💪\n"
	CheckAnswerClick        = "🔎 Отлично! Давай проверим твой ответ ✨\n📸 Пришли фото твоего решения — и я посмотрю, всё ли правильно 😊"
	CheckAnswerText         = "🤓 Отлично, вижу твоё решение!\nПодожди немного — я внимательно проверяю, как ты решил 🧐"
	AnswerCorrectText       = "🎉 Здорово! Всё правильно!\nТы отлично справился 💪\nТы не просто решил — ты понял, как это работает 🌟\nДавай продолжим с новым заданием!"
	AnswerIncorrectText     = "Почти получилось! 💪\nТы был очень близок к правильному ответу 👀\nДавай я подскажу, что можно исправить 💡\n%s\nДавай продолжим с новым заданием!"
	ReportText              = "Отлично, спасибо, что заметил! 📝\nЧтобы я стал лучше, напиши, в чем ошибка?\nНапример:\n• я неправильно прочитал часть задания;\n• подсказка не помогла или была непонятной;\n• я объяснил не то задание;\n• другое (опиши своими словами).\n💬 Напиши коротко, что именно было не так — и я обязательно учту это, чтобы стать умнее и помогать ещё точнее 💡\n👉 После твоего сообщения мы сразу продолжим разбор твоего следующего задания."
	SendReportText          = "👋 Спасибо, что помогаешь мне стать лучше!\nЯ готов продолжить. Скидывай своё задание — и разберёмся вместе! 🤓"
	DontLikeHint            = "😌 Спасибо, я понял!\nТвоя оценка очень важна — я постараюсь объяснить по-другому 💡"
	ErrorText               = "Ой! 😅 Похоже, что-то пошло не так...\nЯ уже стараюсь всё исправить 🔧\n Попробуй позже или нажми 📝 Сообщить об ошибке,\n чтобы рассказать, что случилось, или нажми на 🆕 Новое задание"
	DetectErrorText         = "😥 Не удалось обработать фото."
	GradePreviewText        = "Чтобы я мог давать подсказки подходящего уровня, выбери свой класс 🧩"
	AwaitSolutionText       = "📸 Пришли фото твоего решения — и я посмотрю, всё ли правильно 😊"
	AwaitNewTaskText        = "📸 Скидывай своё задание — и разберёмся вместе! 🤓"
	StepSolutionText        = "\n\n\n\n📘 Шаги решения\n\n"

	YesButton          = "✅ Да, направь подсказку"
	CheckAnswerButton  = "🔎 Проверь мой ответ"
	SendReportButton   = "📝 Сообщить об ошибке"
	NextHintButton     = "➡️ Следующая подсказка"
	DontLikeHintButton = "👎 Не нравиться подсказка"
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
		tgbotapi.NewInlineKeyboardRow(btnGrade1),
		tgbotapi.NewInlineKeyboardRow(btnGrade2),
		tgbotapi.NewInlineKeyboardRow(btnGrade3),
		tgbotapi.NewInlineKeyboardRow(btnGrade4),
	}
}

func makeErrorButtons() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(btnNewTask),
		tgbotapi.NewInlineKeyboardRow(btnReport),
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

func makeHintButtons(level int, showAnalogue bool) [][]tgbotapi.InlineKeyboardButton {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 4)
	if level < 3 {
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
