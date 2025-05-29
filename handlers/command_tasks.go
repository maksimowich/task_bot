package handlers

import (
	"fmt"
)

func (h *Handler) HandleTasks() {
	chatId := h.Update.Message.Chat.ID

	tasks, err := h.Storage.GetTasks(0, 0)
	if err != nil {
		text := "Ошибка при получении всех задач"
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
		if task.AsigneeUser == nil {
			successText += fmt.Sprintf(
				"%d. %s by @%s\n/assign_%d\n\n",
				task.Id, task.Name, task.CreatorUser.UserName, task.Id,
			)
		} else if task.AsigneeUser.Id == h.UserId {
			successText += fmt.Sprintf(
				"%d. %s by @%s\nassignee: я\n/unassign_%d /resolve_%d\n\n",
				task.Id, task.Name, task.CreatorUser.UserName, task.Id, task.Id,
			)
		} else {
			successText += fmt.Sprintf(
				"%d. %s by @%s\nassignee: @%s\n\n",
				task.Id, task.Name, task.CreatorUser.UserName, task.AsigneeUser.UserName,
			)
		}
	}
	if len(successText) >= 2 {
		successText = successText[:len(successText)-2]
	}
	h.Reply(chatId, successText)
}
