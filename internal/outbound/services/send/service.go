package send

import (
	coremodels "euromoby.com/core/models"
	corerepos "euromoby.com/core/repos"
	"euromoby.com/smsgw/internal/billing"
	"euromoby.com/smsgw/internal/config"
	ois "euromoby.com/smsgw/internal/outbound/inputs/send"
	"euromoby.com/smsgw/internal/outbound/models"
	"euromoby.com/smsgw/internal/outbound/outputs"
	org "euromoby.com/smsgw/internal/outbound/repos/group"
	orm "euromoby.com/smsgw/internal/outbound/repos/message"
)

type Service struct {
	app *config.App
	b   billing.Billing
}

func NewService(app *config.App, b billing.Billing) *Service {
	return &Service{app, b}
}

func (s *Service) SendMessage(params *ois.Params) (*outputs.GroupView, error) {
	if err := s.b.CheckBalance(params.MerchantID); err != nil {
		return nil, err
	}

	tx, err := corerepos.Begin(s.app.Config.DBPool)
	if err != nil {
		return nil, err
	}

	defer corerepos.Rollback(tx)

	messageGroup := makeMessageGroup(params)

	if err := org.NewRepo(tx).Save(messageGroup); err != nil {
		return nil, err
	}

	messages := make([]*models.Message, 0, len(params.Recipients))

	for _, msisdn := range params.Recipients {
		outboundMessage := models.NewMessage(params.MerchantID, messageGroup.ID, msisdn)
		messages = append(messages, outboundMessage)
	}

	repo := orm.NewRepo(tx)

	for i := range messages {
		err = repo.Save(messages[i])
		if err != nil {
			return nil, err
		}
	}

	if err := corerepos.Commit(tx); err != nil {
		return nil, err
	}

	var views []*outputs.MessageView
	for _, message := range messages {
		views = append(views, outputs.NewMessageView(message, nil))
	}

	return outputs.NewGroupView(messageGroup, views), nil
}

func makeMessageGroup(p *ois.Params) *models.MessageGroup {
	now := coremodels.TimeNow()

	return &models.MessageGroup{
		MerchantID:          p.MerchantID,
		Sender:              p.Sender,
		Body:                p.Body,
		NotificationURL:     p.NotificationURL,
		ClientTransactionID: p.ClientTransactionID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
