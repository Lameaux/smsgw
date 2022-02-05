package repos

import (
	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
	"github.com/jackc/pgx/v4"
)

const (
	selectDeliveryNotificationBase = `select
	id, message_type, message_id, status,
	last_response,
	next_attempt_at, attempt_counter,
	created_at, updated_at
	from delivery_notifications
	`
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

func (r *DeliveryNotificationRepo) Update(om *models.DeliveryNotification) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	stmt := `update delivery_notifications
	set status = $1,
	last_response = $2,
	next_attempt_at = $3, attempt_counter = $4,
	updated_at = $5
	where id = $6`

	om.UpdatedAt = models.TimeNow()

	_, err := r.db.Exec(ctx, stmt,
		om.Status,
		om.LastResponse,
		om.NextAttemptAt,
		om.AttemptCounter,
		om.UpdatedAt,
		om.ID,
	)

	return err
}

func (r *DeliveryNotificationRepo) FindOneForProcessing(messageType models.MessageType) (*models.DeliveryNotification, error) {
	stmt := selectDeliveryNotificationBase + `
	where message_type = $1
	and status = $2
	and next_attempt_at < $3
	for update skip locked
 	limit 1
	`
	return r.querySingle(stmt, messageType, models.OutboundMessageStatusNew, models.TimeNow())
}

func (r *DeliveryNotificationRepo) querySingle(stmt string, args ...interface{}) (*models.DeliveryNotification, error) {
	ctx, cancel := DBQueryContext()
	defer cancel()
	row := r.db.QueryRow(ctx, stmt, args...)

	var d models.DeliveryNotification
	switch err := r.scanRow(row, &d); err {
	case pgx.ErrNoRows:
		return nil, nil
	case nil:
		return &d, nil
	default:
		return nil, err
	}
}

func (r *DeliveryNotificationRepo) scanRow(row pgx.Row, m *models.DeliveryNotification) error {
	return row.Scan(
		&m.ID,
		&m.MessageType,
		&m.MessageID,
		&m.Status,
		&m.LastResponse,
		&m.NextAttemptAt,
		&m.AttemptCounter,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
}
