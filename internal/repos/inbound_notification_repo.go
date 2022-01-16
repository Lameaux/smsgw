package repos

import (
	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
)

type InboundNotificationRepo struct {
	db db.Conn
}

func NewInboundNotificationRepo(db db.Conn) *InboundNotificationRepo {
	return &InboundNotificationRepo{db}
}

func (r *InboundNotificationRepo) Save(in *models.InboundNotification) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	stmt := `insert into inbound_notifications (
		message_id, status,
		provider_response,
		next_attempt_at, attempt_counter,
		created_at, updated_at
	)
	values ($1, $2, $3, $4, $5, $6, $7)
	returning id
	`
	var insertedID string
	err := r.db.QueryRow(ctx, stmt,
		in.MessageID,
		in.Status,
		in.ProviderResponse,
		in.NextAttemptAt,
		in.AttemptCounter,
		in.CreatedAt,
		in.UpdatedAt,
	).Scan(&insertedID)

	if err != nil {
		return err
	}
	in.ID = insertedID

	return nil
}
