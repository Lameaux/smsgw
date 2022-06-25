package inbound

import (
	"errors"
	commonrepos "euromoby.com/smsgw/internal/common/repos"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"

	"euromoby.com/core/db"
	"euromoby.com/smsgw/internal/inbound/models"

	coremodels "euromoby.com/core/models"
	corerepos "euromoby.com/core/repos"
)

const (
	tableNameInboundMessages        = "inbound_messages"
	constraintNameProviderMessageID = "inbound_provider_message_id"
)

type Repo struct {
	db db.Conn
}

func NewRepo(db db.Conn) *Repo {
	return &Repo{db}
}

func (r *Repo) Save(im *models.Message) error {
	sb := corerepos.DBQueryBuilder().Insert(tableNameInboundMessages).
		Columns(
			"merchant_id",
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
			im.MerchantID,
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

	if err := corerepos.DBQuerySingle(r.db, &im.ID, sb); err != nil {
		return wrapError(err)
	}

	return nil
}

func (r *Repo) Update(m *models.Message) error {
	m.UpdatedAt = coremodels.TimeNow()

	sb := corerepos.DBQueryBuilder().Update(tableNameInboundMessages).SetMap(
		map[string]interface{}{
			"status":     m.Status,
			"updated_at": m.UpdatedAt,
		},
	).Where("id = ?", m.ID)

	return corerepos.DBExec(r.db, sb)
}

func (r *Repo) FindByMerchantAndID(merchantID, id string) (*models.Message, error) {
	var m models.Message

	sb := selectBase().
		Where("merchant_id = ?", merchantID).
		Where("id = ?", id)

	err := corerepos.DBQuerySingle(r.db, &m, sb)

	return &m, err
}

func (r *Repo) FindByQuery(q *SearchParams) ([]*models.Message, error) {
	sb := selectBase().Where("merchant_id = ?", q.MerchantID)

	if q.Shortcode != nil {
		sb = sb.Where("shortcode = ?", q.Shortcode)
	}

	sb = commonrepos.AppendMessageParams(q.MessageParams, sb)
	sb = commonrepos.AppendSearchParams(q.SearchParams, sb)

	var messages []*models.Message
	err := corerepos.DBQueryAll(r.db, &messages, sb)

	return messages, err
}

func selectBase() sq.SelectBuilder {
	return corerepos.DBQueryBuilder().Select(
		"id",
		"merchant_id",
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

func wrapError(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.ConstraintName == constraintNameProviderMessageID {
			return models.ErrDuplicateProviderMessageID
		}
	}

	return err
}
