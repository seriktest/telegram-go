package usecase

import (
	"context"
	"fmt"
	"io"

	"telegram-go/internal/domain"
)

type VideoUseCase struct {
	dlFactory DownloaderFactory
}

// Для простоты здесь используем интерфейс поиска, но в app.go мы передадим фабрику
type DownloaderFactory interface {
	GetDownloader(url string) domain.Downloader
}

func NewVideoUseCase(factory DownloaderFactory) *VideoUseCase {
	return &VideoUseCase{dlFactory: factory}
}

func (uc *VideoUseCase) Download(ctx context.Context, url string) (io.ReadCloser, *domain.VideoInfo, error) {
	downloader := uc.dlFactory.GetDownloader(url)
	if downloader == nil {
		return nil, nil, fmt.Errorf("unsupported URL")
	}

	return downloader.Download(url)
}
