package handlers

import (
	"fmt"
)

func (h *Handler) HandleMy() {
	chatId := h.Update.Message.Chat.ID

	tasks, err := h.Storage.GetTasks(h.UserId, 0)
	if err != nil {
		text := "Ошибка при получении нзначенных на вас задач"
		h.Reply(chatId, text)
		return
	}

	if len(tasks) == 0 {
		text := "Нет задач"
		h.Reply(chatId, text)
		return
	}

	var successText string
	for _, task := range tasks {
		successText += fmt.Sprintf(
			"%d. %s by @%s\n/unassign_%d /resolve_%d\n\n",
			task.Id, task.Name, task.CreatorUser.UserName, task.Id, task.Id,
		)
	}
	if len(successText) >= 2 {
		successText = successText[:len(successText)-2]
	}
	h.Reply(chatId, successText)
}
