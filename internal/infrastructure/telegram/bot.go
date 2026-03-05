package telegram

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"telegram-go/internal/usecase"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	bot     *telego.Bot
	videoUC *usecase.VideoUseCase
	cancel  context.CancelFunc
}

func New(token string, videoUC *usecase.VideoUseCase) (*Bot, error) {
	bot, err := telego.NewBot(token)
	if err != nil {
		return nil, err
	}

	return &Bot{bot: bot, videoUC: videoUC}, nil
}

func (b *Bot) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	// Получаем обновления через Long Polling
	updates, err := b.bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot started...")

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) Stop() {
	if b.cancel != nil {
		b.cancel()
	}
}

func (b *Bot) handleMessage(message *telego.Message) {
	chatID := message.Chat.ID
	text := message.Text

	log.Printf("Received message from %d: %s", chatID, text)

	if text == "/start" {
		b.sendMessage(chatID, "Привет! Отправь мне ссылку на видео с YouTube.")
		return
	}

	if !strings.Contains(text, "http") {
		b.sendMessage(chatID, "Пожалуйста, отправь ссылку.")
		return
	}

	// Отправляем статус "загружает"
	_ = b.bot.SendChatAction(context.Background(), tu.ChatAction(tu.ID(chatID), telego.ChatActionUploadVideo))

	// Скачивание
	stream, info, err := b.videoUC.Download(context.Background(), text)
	if err != nil {
		log.Printf("Download error for URL %s: %v", text, err)
		b.sendMessage(chatID, fmt.Sprintf("Ошибка: %v", err))
		return
	}
	defer stream.Close()

	// Сохраняем во временный файл (Telegram API требует файл для отправки видео)
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("video_%d.mp4", message.MessageID))
	file, err := os.Create(tmpFile)
	if err != nil {
		b.sendMessage(chatID, "Ошибка сохранения файла")
		return
	}
	defer os.Remove(tmpFile)

	if _, err := io.Copy(file, stream); err != nil {
		b.sendMessage(chatID, "Ошибка записи файла")
		file.Close()
		return
	}

	// Закрываем файл перед его отправкой, чтобы гарантировать что все данные сброшены на диск
	file.Close()

	videoFile, err := os.Open(tmpFile)
	if err != nil {
		b.sendMessage(chatID, "Ошибка открытия файла")
		return
	}
	defer videoFile.Close()

	// Отправка видео
	caption := fmt.Sprintf("*%s*", info.Title)
	_, err = b.bot.SendVideo(context.Background(), tu.Video(tu.ID(chatID), tu.File(videoFile)).
		WithCaption(caption).
		WithParseMode(telego.ModeMarkdown),
	)

	if err != nil {
		b.sendMessage(chatID, "Ошибка отправки видео")
	}
}

func (b *Bot) sendMessage(chatID int64, text string) {
	_, _ = b.bot.SendMessage(context.Background(), tu.Message(tu.ID(chatID), text))
}
