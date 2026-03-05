package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"

	"telegram-go/internal/domain"
)

type InstagramDownloader struct{}

func NewInstagramDownloader() *InstagramDownloader {
	return &InstagramDownloader{}
}

func (i *InstagramDownloader) CanHandle(videoURL string) bool {
	u, err := url.Parse(videoURL)
	if err != nil {
		return false
	}
	return strings.Contains(u.Host, "instagram.com") || strings.Contains(u.Host, "ddinstagram.com")
}

func (i *InstagramDownloader) Download(videoURL string) (io.ReadCloser, *domain.VideoInfo, error) {
	// 1. Получаем метаданные через yt-dlp -j
	cmdInfo := exec.Command("yt-dlp", "-j", videoURL)
	var infoBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmdInfo.Stdout = &infoBuf
	cmdInfo.Stderr = &errBuf
	if err := cmdInfo.Run(); err != nil {
		return nil, nil, fmt.Errorf("failed to get video info: %w (stderr: %s)", err, errBuf.String())
	}

	var metadata struct {
		Title    string  `json:"title"`
		Duration float64 `json:"duration"`
	}
	if err := json.Unmarshal(infoBuf.Bytes(), &metadata); err != nil {
		return nil, nil, fmt.Errorf("failed to parse metadata: %w, output: %s", err, infoBuf.String())
	}

	// 2. Запускаем скачивание в stdout
	cmdDownload := exec.Command("yt-dlp", "-o", "-", videoURL)
	stdout, err := cmdDownload.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmdDownload.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start download: %w", err)
	}

	info := &domain.VideoInfo{
		Title:    metadata.Title,
		Duration: int(metadata.Duration),
	}

	// Оборачиваем pipe, чтобы процесс завершался при закрытии ридера
	return &commandReadCloser{
		ReadCloser: stdout,
		cmd:        cmdDownload,
	}, info, nil
}

type commandReadCloser struct {
	io.ReadCloser
	cmd *exec.Cmd
}

func (c *commandReadCloser) Close() error {
	err := c.ReadCloser.Close()
	_ = c.cmd.Process.Kill()
	_ = c.cmd.Wait()
	return err
}
