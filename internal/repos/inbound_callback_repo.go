package repos

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
)

type InboundCallbackRepo struct {
	db db.Conn
}

const (
	constraintNameMerchantIDShortcode = "inbound_callbacks_merchant_id_shortcode"
	tableNameInboundCallbacks         = "inbound_callbacks"
)

func NewInboundCallbackRepo(db db.Conn) *InboundCallbackRepo {
	return &InboundCallbackRepo{db}
}

func (r *InboundCallbackRepo) Save(callback *models.InboundCallback) error {
	sb := dbQueryBuilder().Insert(tableNameInboundCallbacks).
		Columns(
			"merchant_id",
			"shortcode",
			"url",
			"created_at",
			"updated_at",
		).
		Values(
			callback.MerchantID,
			callback.Shortcode,
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

func (r *InboundCallbackRepo) Update(callback *models.InboundCallback) error {
	callback.UpdatedAt = models.TimeNow()

	sb := dbQueryBuilder().Update(tableNameInboundCallbacks).SetMap(
		map[string]interface{}{
			"url":        callback.URL,
			"updated_at": callback.UpdatedAt,
		},
	).Where("merchant_id = ?", callback.MerchantID).
		Where("shortcode = ?", callback.Shortcode)

	return dbExec(r.db, sb)
}

func (r *InboundCallbackRepo) Delete(callback *models.InboundCallback) error {
	callback.UpdatedAt = models.TimeNow()

	sb := dbQueryBuilder().Delete(tableNameInboundCallbacks).
		Where("merchant_id = ?", callback.MerchantID).
		Where("shortcode = ?", callback.Shortcode)

	return dbExec(r.db, sb)
}

func (r *InboundCallbackRepo) FindByMerchant(merchantID string) ([]*models.InboundCallback, error) {
	sb := r.selectBase().Where("merchant_id = ?", merchantID)

	callbacks := []*models.InboundCallback{}
	err := dbQueryAll(r.db, &callbacks, sb)

	return callbacks, err
}

func (r *InboundCallbackRepo) FindByMerchantAndShortcode(merchantID, shortcode string) (*models.InboundCallback, error) {
	var callback models.InboundCallback

	sb := r.selectBase().Where("merchant_id = ?", merchantID).
		Where("shortcode = ?", shortcode)
	err := dbQuerySingle(r.db, &callback, sb)

	return &callback, err
}

func (r *InboundCallbackRepo) selectBase() sq.SelectBuilder {
	return dbQueryBuilder().Select(
		"id",
		"merchant_id",
		"shortcode",
		"url",
		"created_at",
		"updated_at",
	).From(tableNameInboundCallbacks)
}

func (r *InboundCallbackRepo) wrapError(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.ConstraintName == constraintNameMerchantIDShortcode {
			return models.ErrDuplicateCallback
		}
	}

	return err
}
