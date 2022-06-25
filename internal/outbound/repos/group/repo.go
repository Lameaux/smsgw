package group

import (
	"errors"
	"euromoby.com/smsgw/internal/outbound/inputs/group"
	"euromoby.com/smsgw/internal/repos"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"

	"euromoby.com/core/db"
	corerepos "euromoby.com/core/repos"

	"euromoby.com/smsgw/internal/outbound/models"
)

type Repo struct {
	db db.Conn
}

const (
	constraintNameClientTransactionID = "message_groups_client_transaction_id"
	tableNameMessageGroups            = "message_groups"
)

func NewRepo(db db.Conn) *Repo {
	return &Repo{db}
}

func (r *Repo) Save(mg *models.MessageGroup) error {
	sb := corerepos.DBQueryBuilder().Insert(tableNameMessageGroups).
		Columns(
			"merchant_id",
			"sender",
			"body",
			"client_transaction_id",
			"notification_url",
			"created_at",
			"updated_at",
		).
		Values(
			mg.MerchantID,
			mg.Sender,
			mg.Body,
			mg.ClientTransactionID,
			mg.NotificationURL,
			mg.CreatedAt,
			mg.UpdatedAt,
		).
		Suffix(`RETURNING "id"`)

	if err := corerepos.DBQuerySingle(r.db, &mg.ID, sb); err != nil {
		return wrapError(err)
	}

	return nil
}

func (r *Repo) FindByID(id string) (*models.MessageGroup, error) {
	var group models.MessageGroup

	sb := selectBase().Where("id = ?", id)
	err := corerepos.DBQuerySingle(r.db, &group, sb)

	return &group, err
}

func (r *Repo) FindByMerchantAndID(merchantID, id string) (*models.MessageGroup, error) {
	var group models.MessageGroup

	sb := selectBase().Where("merchant_id = ? AND id = ?", merchantID, id)
	err := corerepos.DBQuerySingle(r.db, &group, sb)

	return &group, err
}

func (r *Repo) FindByQuery(q *group.SearchParams) ([]*models.MessageGroup, error) {
	sb := selectBase().Where("merchant_id = ?", q.MerchantID)

	if q.ClientTransactionID != nil {
		sb = sb.Where("client_transaction_id = ?", q.ClientTransactionID)
	}

	sb = repos.AppendSearchParams(q.SearchParams, sb)

	var groups []*models.MessageGroup
	err := corerepos.DBQueryAll(r.db, &groups, sb)

	return groups, err
}

func selectBase() sq.SelectBuilder {
	return corerepos.DBQueryBuilder().Select(
		"id",
		"merchant_id",
		"sender",
		"body",
		"client_transaction_id",
		"notification_url",
		"created_at",
		"updated_at",
	).From(tableNameMessageGroups)
}

func wrapError(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.ConstraintName == constraintNameClientTransactionID {
			return models.ErrDuplicateClientTransactionID
		}
	}

	return err
}
