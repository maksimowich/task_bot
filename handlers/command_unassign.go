package handlers

import (
	"fmt"
	"strconv"
	"strings"
)

func (h *Handler) HandleUnassign() {
	currChatId := h.Update.Message.Chat.ID

	commandArg := strings.TrimPrefix(h.Update.Message.Text, "/unassign_")
	taskId, err := strconv.Atoi(commandArg)
	if err != nil {
		text := "Ошибка при парсинге id задачи"
		h.Reply(currChatId, text)
		return
	}

	task, err := h.Storage.GetTaskById(int64(taskId))
	if err != nil {
		text := "Ошибка при попытке получить задачу"
		h.Reply(currChatId, text)
		return
	}

	if h.UserId != task.AsigneeUser.Id {
		text := "Задача не на вас"
		h.Reply(currChatId, text)
		return
	}

	task, err = h.Storage.UnassignTask(int64(taskId))
	if err != nil {
		text := "Ошибка при снятии исполнителя задачи"
		h.Reply(currChatId, text)
		return
	}

	successTextToCurrUser := "Принято"
	h.Reply(currChatId, successTextToCurrUser)

	successTextToCreatorUser := fmt.Sprintf(
		`Задача "%s" осталась без исполнителя`,
		task.Name,
	)
	h.Reply(task.CreatorUser.ChatId, successTextToCreatorUser)
}
