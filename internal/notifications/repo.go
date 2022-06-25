package notifications

import (
	"euromoby.com/smsgw/internal/notifications/models"
	sq "github.com/Masterminds/squirrel"

	"euromoby.com/core/db"
	coremodels "euromoby.com/core/models"
	corerepos "euromoby.com/core/repos"
)

const (
	tableNameDeliveryNotifications = "delivery_notifications"
)

type Repo struct {
	db db.Conn
}

func NewRepo(db db.Conn) *Repo {
	return &Repo{db}
}

func (r *Repo) Save(dn *models.DeliveryNotification) error {
	sb := corerepos.DBQueryBuilder().Insert(tableNameDeliveryNotifications).
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

	return corerepos.DBQuerySingle(r.db, &dn.ID, sb)
}

func (r *Repo) Update(dn *models.DeliveryNotification) error {
	dn.UpdatedAt = coremodels.TimeNow()

	sb := corerepos.DBQueryBuilder().Update(tableNameDeliveryNotifications).SetMap(
		map[string]interface{}{
			"status":          dn.Status,
			"last_response":   dn.LastResponse,
			"next_attempt_at": dn.NextAttemptAt,
			"attempt_counter": dn.AttemptCounter,
			"updated_at":      dn.UpdatedAt,
		},
	).Where("id = ?", dn.ID)

	return corerepos.DBExec(r.db, sb)
}

func (r *Repo) FindOneForProcessing(messageType models.DeliveryNotificationType) (*models.DeliveryNotification, error) {
	var dn models.DeliveryNotification

	sb := selectBase().
		Where("message_type = ?", messageType).
		Where("status = ?", models.DeliveryNotificationStatusNew).
		Where("next_attempt_at < ?", coremodels.TimeNow()).
		Suffix("for update skip locked").
		Limit(1)

	err := corerepos.DBQuerySingle(r.db, &dn, sb)

	return &dn, err
}

func selectBase() sq.SelectBuilder {
	return corerepos.DBQueryBuilder().Select(
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
