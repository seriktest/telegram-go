package downloader

import (
	"telegram-go/internal/domain"
)

type Factory struct {
	downloaders []domain.Downloader
}

func NewFactory() *Factory {
	return &Factory{
		downloaders: []domain.Downloader{
			NewYouTubeDownloader(),
			NewInstagramDownloader(),
		},
	}
}

func (f *Factory) GetDownloader(url string) domain.Downloader {
	for _, d := range f.downloaders {
		if d.CanHandle(url) {
			return d
		}
	}
	return nil
}
