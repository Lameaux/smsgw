package models

import "errors"

var (
	ErrMessageNotFound      = errors.New("message not found")
	ErrMessageGroupNotFound = errors.New("message group not found")

	ErrMaxRecipients = errors.New("too many recipients")

	ErrMissingRecipients = errors.New("missing recipients")

	ErrSendFailed = errors.New("failed to send")

	ErrAlreadyAcked                 = errors.New("message is already acked")
	ErrDuplicateClientTransactionID = errors.New("client_transaction_id already exists")
)
