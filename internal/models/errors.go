package models

import "errors"

var (
	ErrAlreadyAcked                 = errors.New("message is already acked")
	ErrSendFailed                   = errors.New("failed to send")
	ErrInvalidJSON                  = errors.New("invalid json")
	ErrDeadLetter                   = errors.New("message can not be delivered")
	ErrDuplicateProviderMessageID   = errors.New("message_id already exists")
	ErrDuplicateClientTransactionID = errors.New("client_transaction_id already exists")
	ErrDuplicateCallback            = errors.New("callback already exists")
	ErrNotFound                     = errors.New("not found")
	ErrMissingNotificationURL       = errors.New("missing notification url")
	ErrInsufficientFunds            = errors.New("insufficient funds")
	ErrMessageNotFound              = errors.New("message not found")
	ErrMessageOrderNotFound         = errors.New("message order not found")
	ErrCallbackNotFound             = errors.New("callback not found")
)
