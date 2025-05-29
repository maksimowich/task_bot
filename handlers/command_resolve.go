package handlers

import (
	"fmt"
	"strconv"
	"strings"
)

func (h *Handler) HandleResolve() {
	currChatId := h.Update.Message.Chat.ID

	commandArg := strings.TrimPrefix(h.Update.Message.Text, "/resolve_")
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

	task, err = h.Storage.ResolveTask(int64(taskId))
	if err != nil {
		text := "Ошибка при попытке перевести задачу в Выполнено"
		h.Reply(currChatId, text)
		return
	}

	successTextToCurrUser := fmt.Sprintf(
		`Задача "%s" выполнена`,
		task.Name,
	)
	h.Reply(currChatId, successTextToCurrUser)

	if h.UserId != task.CreatorUser.Id {
		successTextToCreatorUser := fmt.Sprintf(
			`Задача "%s" выполнена @%s`,
			task.Name, task.AsigneeUser.UserName,
		)
		h.Reply(task.CreatorUser.ChatId, successTextToCreatorUser)
	}
}
