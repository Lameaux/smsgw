package repos

import (
	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
)

type DeliveryNotificationRepo struct {
	db db.Conn
}

func NewDeliveryNotificationRepo(db db.Conn) *DeliveryNotificationRepo {
	return &DeliveryNotificationRepo{db}
}

func (r *DeliveryNotificationRepo) Save(in *models.DeliveryNotification) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	stmt := `insert into delivery_notifications (
		message_type, message_id,
		status, last_response,
		next_attempt_at, attempt_counter,
		created_at, updated_at
	)
	values ($1, $2, $3, $4, $5, $6, $7, $8)
	returning id
	`
	var insertedID string
	err := r.db.QueryRow(ctx, stmt,
		in.MessageType,
		in.MessageID,
		in.Status,
		in.LastResponse,
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
