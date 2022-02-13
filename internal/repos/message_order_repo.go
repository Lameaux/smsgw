package repos

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
)

type MessageOrderRepo struct {
	db db.Conn
}

const (
	constraintNameClientTransactionID = "message_orders_client_transaction_id"
)

func NewMessageOrderRepo(db db.Conn) *MessageOrderRepo {
	return &MessageOrderRepo{db}
}

func (r *MessageOrderRepo) Save(mo *models.MessageOrder) error {
	stmt, args, err := DBQueryBuilder().Insert("message_orders").
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
			mo.MerchantID,
			mo.Sender,
			mo.Body,
			mo.ClientTransactionID,
			mo.NotificationURL,
			mo.CreatedAt,
			mo.UpdatedAt,
		).
		Suffix(`RETURNING "id"`).ToSql()
	if err != nil {
		return err
	}

	if err = DBQueryGet(r.db, &mo.ID, stmt, args...); err != nil {
		return r.wrapError(err)
	}

	return nil
}

func (r *MessageOrderRepo) FindByID(id string) (*models.MessageOrder, error) {
	stmt, args, err := r.selectMessageOrdersBase().Where("id = ?", id).ToSql()
	if err != nil {
		return nil, err
	}

	return r.querySingle(stmt, args...)
}

func (r *MessageOrderRepo) FindByMerchantAndID(merchantID, id string) (*models.MessageOrder, error) {
	stmt, args, err := r.selectMessageOrdersBase().Where("merchant_id = ? AND id = ?", merchantID, id).ToSql()
	if err != nil {
		return nil, err
	}

	return r.querySingle(stmt, args...)
}

func (r *MessageOrderRepo) FindByQuery(q *inputs.MessageOrderSearchParams) ([]*models.MessageOrder, error) {
	sb := r.selectMessageOrdersBase()

	sb = sb.Where("merchant_id = ?", q.MerchantID)

	if q.ClientTransactionID != nil {
		sb = sb.Where("client_transaction_id = ?", q.ClientTransactionID)
	}

	sb = appendSearchParams(q.SearchParams, sb)

	stmt, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	return r.query(stmt, args...)
}

func (r *MessageOrderRepo) query(stmt string, args ...interface{}) ([]*models.MessageOrder, error) {
	orders := []*models.MessageOrder{}

	err := DBQuerySelect(r.db, &orders, stmt, args...)

	return orders, err
}

func (r *MessageOrderRepo) querySingle(stmt string, args ...interface{}) (*models.MessageOrder, error) {
	var order models.MessageOrder

	err := DBQueryGet(r.db, &order, stmt, args...)

	return &order, err
}

func (r *MessageOrderRepo) wrapError(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.ConstraintName == constraintNameClientTransactionID {
			return models.ErrDuplicateClientTransactionID
		}
	}

	return err
}

func (r *MessageOrderRepo) selectMessageOrdersBase() sq.SelectBuilder {
	return DBQueryBuilder().Select(
		"id",
		"merchant_id",
		"sender",
		"body",
		"client_transaction_id",
		"notification_url",
		"created_at",
		"updated_at",
	).From("message_orders")
}
