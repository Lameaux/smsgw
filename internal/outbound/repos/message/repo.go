package message

import (
	commonrepos "euromoby.com/smsgw/internal/common/repos"
	"euromoby.com/smsgw/internal/outbound/inputs/message"
	"euromoby.com/smsgw/internal/outbound/models"
	sq "github.com/Masterminds/squirrel"

	"euromoby.com/core/db"
	coremodels "euromoby.com/core/models"
	corerepos "euromoby.com/core/repos"
)

const (
	tableNameOutboundMessages = "outbound_messages"
)

type Repo struct {
	db db.Conn
}

func NewRepo(db db.Conn) *Repo {
	return &Repo{db}
}

func (r *Repo) Save(om *models.Message) error {
	sb := corerepos.DBQueryBuilder().Insert(tableNameOutboundMessages).
		Columns(
			"merchant_id",
			"message_group_id",
			"status",
			"msisdn",
			"next_attempt_at",
			"attempt_counter",
			"created_at",
			"updated_at",
		).
		Values(
			om.MerchantID,
			om.MessageGroupID,
			om.Status,
			om.MSISDN,
			om.NextAttemptAt,
			om.AttemptCounter,
			om.CreatedAt,
			om.UpdatedAt,
		).
		Suffix(`RETURNING "id"`)

	return corerepos.DBQuerySingle(r.db, &om.ID, sb)
}

func (r *Repo) Update(om *models.Message) error {
	om.UpdatedAt = coremodels.TimeNow()

	sb := corerepos.DBQueryBuilder().Update(tableNameOutboundMessages).SetMap(
		map[string]interface{}{
			"status":              om.Status,
			"provider_id":         om.ProviderID,
			"provider_response":   om.ProviderResponse,
			"provider_message_id": om.ProviderMessageID,
			"next_attempt_at":     om.NextAttemptAt,
			"attempt_counter":     om.AttemptCounter,
			"updated_at":          om.UpdatedAt,
		},
	).Where("id = ?", om.ID)

	return corerepos.DBExec(r.db, sb)
}

func (r *Repo) FindByID(id string) (*models.Message, error) {
	var m models.Message

	sb := selectBase().Where("id = ?", id)
	err := corerepos.DBQuerySingle(r.db, &m, sb)

	return &m, err
}

func (r *Repo) FindByMerchantAndID(merchantID, id string) (*models.Message, error) {
	var m models.Message

	sb := selectBase().
		Where("merchant_id = ?", merchantID).
		Where("id = ?", id)

	err := corerepos.DBQuerySingle(r.db, &m, sb)

	return &m, err
}

func (r *Repo) FindByMerchantAndGroupID(merchantID, groupID string) ([]*models.Message, error) {
	sb := selectBase().
		Where("merchant_id = ?", merchantID).
		Where("message_group_id = ?", groupID)

	var messages []*models.Message
	err := corerepos.DBQueryAll(r.db, &messages, sb)

	return messages, err
}

func (r *Repo) FindByProviderAndMessageID(providerID, messageID string) (*models.Message, error) {
	var m models.Message

	sb := selectBase().
		Where("provider_id = ?", providerID).
		Where("provider_message_id = ?", messageID)

	err := corerepos.DBQuerySingle(r.db, &m, sb)

	return &m, err
}

func (r *Repo) FindByQuery(q *message.SearchParams) ([]*models.Message, error) {
	sb := selectBase().Where("merchant_id = ?", q.MerchantID)

	sb = commonrepos.AppendMessageParams(q.MessageParams, sb)
	sb = commonrepos.AppendSearchParams(q.SearchParams, sb)

	var messages []*models.Message
	err := corerepos.DBQueryAll(r.db, &messages, sb)

	return messages, err
}

func (r *Repo) FindOneForProcessing() (*models.Message, error) {
	var m models.Message

	sb := selectBase().
		Where("status = ?", models.MessageStatusNew).
		Where("next_attempt_at < ?", coremodels.TimeNow()).
		Suffix("for update skip locked").
		Limit(1)

	err := corerepos.DBQuerySingle(r.db, &m, sb)

	return &m, err
}

func selectBase() sq.SelectBuilder {
	return corerepos.DBQueryBuilder().Select(
		"id",
		"merchant_id",
		"message_group_id",
		"status",
		"msisdn",
		"provider_id",
		"provider_message_id",
		"provider_response",
		"next_attempt_at",
		"attempt_counter",
		"created_at",
		"updated_at",
	).From(tableNameOutboundMessages)
}
