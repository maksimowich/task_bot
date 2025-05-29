package handlers

import (
	storage "github.com/maksimowich/task_bot/storage"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

type Handler struct {
	Bot     *tgbotapi.BotAPI
	Update  *tgbotapi.Update
	Storage *storage.Storage
	UserId  int64
}

func (h *Handler) Reply(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	h.Bot.Send(msg)
}
