package models

import "errors"

var (
	ErrMissingNotificationURL = errors.New("missing notification url")
	ErrSendFailed             = errors.New("failed to send")
)
