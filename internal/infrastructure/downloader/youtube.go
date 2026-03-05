package downloader

import (
	"errors"
	"io"
	"net/url"
	"strings"

	"telegram-go/internal/domain"

	"github.com/kkdai/youtube/v2"
)

type YouTubeDownloader struct {
	client *youtube.Client
}

func NewYouTubeDownloader() *YouTubeDownloader {
	return &YouTubeDownloader{client: &youtube.Client{}}
}

func (y *YouTubeDownloader) CanHandle(videoURL string) bool {
	u, err := url.Parse(videoURL)
	if err != nil {
		return false
	}
	return strings.Contains(u.Host, "youtube.com") || strings.Contains(u.Host, "youtu.be")
}

func (y *YouTubeDownloader) Download(videoURL string) (io.ReadCloser, *domain.VideoInfo, error) {
	video, err := y.client.GetVideo(videoURL)
	if err != nil {
		return nil, nil, err
	}

	// Выбор формата (видео + аудио)
	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return nil, nil, errors.New("no formats with audio found")
	}
	format := formats[0] // Простой выбор, можно улучшить

	stream, _, err := y.client.GetStream(video, &format)
	if err != nil {
		return nil, nil, err
	}

	info := &domain.VideoInfo{
		Title:    video.Title,
		Duration: int(video.Duration.Seconds()),
	}

	return stream, info, nil
}
