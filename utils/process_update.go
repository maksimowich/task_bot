package utils

import (
	"log"
	"strings"

	handlers "github.com/maksimowich/task_bot/handlers"
	storage "github.com/maksimowich/task_bot/storage"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

func ProcessUpdate(
	update *tgbotapi.Update,
	storage *storage.Storage,
	bot *tgbotapi.BotAPI,
) {
	if update.Message == nil {
		return
	}

	userId := update.Message.From.ID
	userName := update.Message.From.UserName
	userChatId := update.Message.Chat.ID

	_, err := storage.GetOrCreateUser(userId, userName, userChatId)
	if err != nil {
		text := "Бот не может обработать ваш запрос"
		msg := tgbotapi.NewMessage(userChatId, text)
		bot.Send(msg)
		return
	}

	log.Printf(
		"update msg text: '%s'; userId: %d; userName: %ss\n",
		update.Message.Text, userId, userName,
	)

	if update.Message.IsCommand() {
		handler := handlers.Handler{
			Bot:     bot,
			Update:  update,
			Storage: storage,
			UserId:  userId,
		}

		command := update.Message.Command()
		baseCommand := strings.Split(command, "_")[0]

		switch baseCommand {

		case "new":
			handler.HandleNew()

		case "tasks":
			handler.HandleTasks()

		case "assign":
			handler.HandleAssign()

		case "unassign":
			handler.HandleUnassign()

		case "resolve":
			handler.HandleResolve()

		case "my":
			handler.HandleMy()

		case "owner":
			handler.HandleOwner()

		default:
			chatId := update.Message.Chat.ID
			text := "Вы ввели неизвестную команду. Попробуйте ещё раз"
			handler.Reply(chatId, text)
		}

	}
}
