package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"telegram-go/internal/config"
	"telegram-go/internal/infrastructure/downloader"
	"telegram-go/internal/infrastructure/telegram"
	"telegram-go/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Инициализация зависимостей
	dlFactory := downloader.NewFactory()
	videoUC := usecase.NewVideoUseCase(dlFactory)

	log.Printf("Loaded token: '%s'", cfg.TelegramToken)

	bot, err := telegram.New(cfg.TelegramToken, videoUC)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	// Запуск бота в горутине
	go bot.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down bot...")
	bot.Stop()
}
