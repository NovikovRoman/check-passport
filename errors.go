package check_passport

import (
	"fmt"
	"strings"
)

// ErrWebsiteUnavailable при ошибках при загрузке страницы сайта
type ErrWebsiteUnavailable struct{}

func (e ErrWebsiteUnavailable) Error() string {
	return "Website unavailable"
}

// ErrArchiveDownload при ошибках при загрузке архива
type ErrArchiveDownload struct {
	msg string
}

func NewErrArchiveDownload(msg ...string) *ErrArchiveDownload {
	return &ErrArchiveDownload{
		msg: strings.Join(msg, "\n"),
	}
}

func (e *ErrArchiveDownload) Error() string {
	return fmt.Sprintf("Archive download error %v", e.msg)
}
