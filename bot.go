package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	storage "github.com/maksimowich/task_bot/storage"
	utils "github.com/maksimowich/task_bot/utils"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var (
	WebhookURL = ""
	BotPort    = "8081"
	BotToken   = ""
)

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return fmt.Errorf("NewBotAPI failed: %w", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		return fmt.Errorf("NewWebhook failed: %w", err)
	}

	if _, err = bot.Request(wh); err != nil {
		return fmt.Errorf("SetWebhook failed: %w", err)
	}

	updates := bot.ListenForWebhook("/")

	storage := &storage.Storage{
		Users:      make(map[int64]*storage.User),
		Tasks:      make(map[int64]*storage.Task),
		NextTaskId: 1,
		Mu:         sync.RWMutex{},
	}

	server := &http.Server{Addr: ":" + BotPort}
	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all is working"))
	})

	serverErr := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Printf("HTTP server shutdown error: %v", err)
			}
			return ctx.Err()

		case err := <-serverErr:
			return err

		case update := <-updates:
			go func(update tgbotapi.Update) {
				select {
				case <-ctx.Done():
					return
				default:
					utils.ProcessUpdate(&update, storage, bot)
				}
			}(update)
		}
	}
}

func main() {
	if err := startTaskBot(context.Background()); err != nil {
		log.Fatalf("Bot failed: %v", err)
	}
}
