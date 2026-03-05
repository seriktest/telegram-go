package domain

import "io"

type VideoInfo struct {
	Title    string
	Duration int
}

type Downloader interface {
	Download(url string) (io.ReadCloser, *VideoInfo, error)
	CanHandle(url string) bool
}