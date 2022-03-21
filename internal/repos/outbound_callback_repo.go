package repos

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
)

type OutboundCallbackRepo struct {
	db db.Conn
}

const (
	constraintNameMerchantID   = "outbound_callbacks_merchant_id"
	tableNameOutboundCallbacks = "outbound_callbacks"
)

func NewOutboundCallbackRepo(db db.Conn) *OutboundCallbackRepo {
	return &OutboundCallbackRepo{db}
}

func (r *OutboundCallbackRepo) Save(callback *models.OutboundCallback) error {
	sb := dbQueryBuilder().Insert(tableNameOutboundCallbacks).
		Columns(
			"merchant_id",
			"url",
			"created_at",
			"updated_at",
		).
		Values(
			callback.MerchantID,
			callback.URL,
			callback.CreatedAt,
			callback.UpdatedAt,
		).
		Suffix(`RETURNING "id"`)

	if err := dbQuerySingle(r.db, &callback.ID, sb); err != nil {
		return r.wrapError(err)
	}

	return nil
}

func (r *OutboundCallbackRepo) Update(callback *models.OutboundCallback) error {
	callback.UpdatedAt = models.TimeNow()

	sb := dbQueryBuilder().Update(tableNameOutboundCallbacks).SetMap(
		map[string]interface{}{
			"url":        callback.URL,
			"updated_at": callback.UpdatedAt,
		},
	).Where("merchant_id = ?", callback.MerchantID)

	return dbExec(r.db, sb)
}

func (r *OutboundCallbackRepo) Delete(callback *models.OutboundCallback) error {
	callback.UpdatedAt = models.TimeNow()

	sb := dbQueryBuilder().Delete(tableNameOutboundCallbacks).
		Where("merchant_id = ?", callback.MerchantID)

	return dbExec(r.db, sb)
}

func (r *OutboundCallbackRepo) FindByMerchant(merchantID string) (*models.OutboundCallback, error) {
	var callback models.OutboundCallback

	sb := r.selectBase().Where("merchant_id = ?", merchantID)
	err := dbQuerySingle(r.db, &callback, sb)

	return &callback, err
}

func (r *OutboundCallbackRepo) selectBase() sq.SelectBuilder {
	return dbQueryBuilder().Select(
		"id",
		"merchant_id",
		"url",
		"created_at",
		"updated_at",
	).From(tableNameOutboundCallbacks)
}

func (r *OutboundCallbackRepo) wrapError(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.ConstraintName == constraintNameMerchantID {
			return models.ErrDuplicateCallback
		}
	}

	return err
}
