package repos

import (
	sq "github.com/Masterminds/squirrel"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
)

const (
	tableNameDeliveryNotifications = "delivery_notifications"
)

type DeliveryNotificationRepo struct {
	db db.Conn
}

func NewDeliveryNotificationRepo(db db.Conn) *DeliveryNotificationRepo {
	return &DeliveryNotificationRepo{db}
}

func (r *DeliveryNotificationRepo) Save(dn *models.DeliveryNotification) error {
	sb := dbQueryBuilder().Insert(tableNameDeliveryNotifications).
		Columns(
			"message_type",
			"message_id",
			"status",
			"last_response",
			"next_attempt_at",
			"attempt_counter",
			"created_at",
			"updated_at",
		).
		Values(
			dn.MessageType,
			dn.MessageID,
			dn.Status,
			dn.LastResponse,
			dn.NextAttemptAt,
			dn.AttemptCounter,
			dn.CreatedAt,
			dn.UpdatedAt,
		).
		Suffix(`RETURNING "id"`)

	return dbQuerySingle(r.db, &dn.ID, sb)
}

func (r *DeliveryNotificationRepo) Update(dn *models.DeliveryNotification) error {
	dn.UpdatedAt = models.TimeNow()

	sb := dbQueryBuilder().Update(tableNameOutboundMessages).SetMap(
		map[string]interface{}{
			"status":          dn.Status,
			"last_response":   dn.LastResponse,
			"next_attempt_at": dn.NextAttemptAt,
			"attempt_counter": dn.AttemptCounter,
			"updated_at":      dn.UpdatedAt,
		},
	).Where("id = ?", dn.ID)

	return dbExec(r.db, sb)
}

func (r *DeliveryNotificationRepo) FindOneForProcessing(messageType models.MessageType) (*models.DeliveryNotification, error) {
	var msg models.DeliveryNotification

	sb := r.selectBase().
		Where("message_type = ?", messageType).
		Where("status = ?", models.DeliveryNotificationStatusNew).
		Where("next_attempt_at < ?", models.TimeNow()).
		Suffix("for update skip locked").
		Limit(1)

	err := dbQuerySingle(r.db, &msg, sb)

	return &msg, err
}

func (r *DeliveryNotificationRepo) selectBase() sq.SelectBuilder {
	return dbQueryBuilder().Select(
		"id",
		"message_type",
		"message_id",
		"status",
		"last_response",
		"next_attempt_at",
		"attempt_counter",
		"created_at",
		"updated_at",
	).From(tableNameDeliveryNotifications)
}
