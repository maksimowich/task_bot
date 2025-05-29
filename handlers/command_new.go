package handlers

import (
	"fmt"
)

func (h *Handler) HandleNew() {
	chatId := h.Update.Message.Chat.ID

	taskName := h.Update.Message.CommandArguments()
	if taskName == "" {
		text := "Укажите название задачи после /new"
		h.Reply(chatId, text)
		return
	}

	task, err := h.Storage.CreateTask(taskName, h.UserId)
	if err != nil {
		text := "Ошибка при создании задачи"
		h.Reply(chatId, text)
		return
	}

	successText := fmt.Sprintf(
		`Задача "%s" создана, id=%d`,
		task.Name, task.Id,
	)
	h.Reply(chatId, successText)
}
