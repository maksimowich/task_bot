package handlers

import (
	"fmt"
	"strconv"
	"strings"
)

func (h *Handler) HandleAssign() {
	currChatId := h.Update.Message.Chat.ID

	commandArg := strings.TrimPrefix(h.Update.Message.Text, "/assign_")
	taskId, err := strconv.Atoi(commandArg)
	if err != nil {
		text := "Ошибка при парсинге id задачи"
		h.Reply(currChatId, text)
		return
	}

	task, err := h.Storage.GetTaskById(int64(taskId))
	if err != nil {
		text := "Ошибка при поиске указанной задачи"
		h.Reply(currChatId, text)
		return
	}
	var prevAssigneeUserId, prevAssigneeChatId int64
	if task.AsigneeUser != nil {
		prevAssigneeUserId = task.AsigneeUser.Id
		prevAssigneeChatId = task.AsigneeUser.ChatId
	}

	task, err = h.Storage.AssignTask(int64(taskId), h.UserId)
	if err != nil {
		text := "Ошибка при назначении исполнителя задачи"
		h.Reply(currChatId, text)
		return
	}

	successTextToAssignee := fmt.Sprintf(
		`Задача "%s" назначена на вас`,
		task.Name,
	)
	h.Reply(currChatId, successTextToAssignee)

	if prevAssigneeUserId != 0 && h.UserId != prevAssigneeUserId {
		successTextToPrevAssignee := fmt.Sprintf(
			`Задача "%s" назначена на @%s`,
			task.Name, task.AsigneeUser.UserName,
		)
		h.Reply(prevAssigneeChatId, successTextToPrevAssignee)
	} else if prevAssigneeUserId == 0 && h.UserId != task.CreatorUser.Id {
		successTextToCreator := fmt.Sprintf(
			`Задача "%s" назначена на @%s`,
			task.Name, task.AsigneeUser.UserName,
		)
		h.Reply(task.CreatorUser.ChatId, successTextToCreator)
	}
}
