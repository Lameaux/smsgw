package repos

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
)

const (
	tableNameInboundMessages        = "inbound_messages"
	constraintNameProviderMessageID = "inbound_provider_message_id"
)

type InboundMessageRepo struct {
	db db.Conn
}

func NewInboundMessageRepo(db db.Conn) *InboundMessageRepo {
	return &InboundMessageRepo{db}
}

func (r *InboundMessageRepo) Save(im *models.InboundMessage) error {
	sb := dbQueryBuilder().Insert(tableNameInboundMessages).
		Columns(
			"shortcode",
			"status",
			"msisdn",
			"body",
			"provider_id",
			"provider_message_id",
			"next_attempt_at",
			"attempt_counter",
			"created_at",
			"updated_at",
		).
		Values(
			im.Shortcode,
			im.Status,
			im.MSISDN,
			im.Body,
			im.ProviderID,
			im.ProviderMessageID,
			im.NextAttemptAt,
			im.AttemptCounter,
			im.CreatedAt,
			im.UpdatedAt,
		).
		Suffix(`RETURNING "id"`)

	if err := dbQuerySingle(r.db, &im.ID, sb); err != nil {
		return r.wrapError(err)
	}

	return nil
}

func (r *InboundMessageRepo) Update(im *models.InboundMessage) error {
	im.UpdatedAt = models.TimeNow()

	sb := dbQueryBuilder().Update(tableNameInboundMessages).SetMap(
		map[string]interface{}{
			"status":     im.Status,
			"updated_at": im.UpdatedAt,
		},
	).Where("id = ?", im.ID)

	return dbExec(r.db, sb)
}

func (r *InboundMessageRepo) FindByShortcodeAndID(shortcode, id string) (*models.InboundMessage, error) {
	var msg models.InboundMessage

	sb := r.selectBase().
		Where("shortcode = ?", shortcode).
		Where("id = ?", id)

	err := dbQuerySingle(r.db, &msg, sb)

	return &msg, err
}

func (r *InboundMessageRepo) FindByQuery(q *inputs.InboundMessageSearchParams) ([]*models.InboundMessage, error) {
	sb := r.selectBase().Where("shortcode = ?", q.Shortcode)

	sb = appendMessageParams(q.MessageParams, sb)
	sb = appendSearchParams(q.SearchParams, sb)

	messages := []*models.InboundMessage{}
	err := dbQueryAll(r.db, &messages, sb)

	return messages, err
}

func (r *InboundMessageRepo) selectBase() sq.SelectBuilder {
	return dbQueryBuilder().Select(
		"id",
		"shortcode",
		"status",
		"msisdn",
		"body",
		"provider_id",
		"provider_message_id",
		"next_attempt_at",
		"attempt_counter",
		"created_at",
		"updated_at",
	).From(tableNameInboundMessages)
}

func (r *InboundMessageRepo) wrapError(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.ConstraintName == constraintNameProviderMessageID {
			return models.ErrDuplicateProviderMessageID
		}
	}

	return err
}
