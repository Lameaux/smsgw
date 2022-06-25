package group

import (
	corerepos "github.com/Lameaux/core/repos"
	"github.com/Lameaux/smsgw/internal/config"
	ig "github.com/Lameaux/smsgw/internal/outbound/inputs/group"
	"github.com/Lameaux/smsgw/internal/outbound/models"
	v "github.com/Lameaux/smsgw/internal/outbound/outputs"
	rg "github.com/Lameaux/smsgw/internal/outbound/repos/group"
	rm "github.com/Lameaux/smsgw/internal/outbound/repos/message"
)

type Service struct {
	app *config.App
}

func NewService(app *config.App) *Service {
	return &Service{app}
}

func (s *Service) FindByMerchantAndID(merchantID, id string) (*v.GroupView, error) {
	ctx, cancel := corerepos.DBConnContext()
	defer cancel()

	conn, err := s.app.Config.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	messageGroup, err := rg.NewRepo(conn).FindByMerchantAndID(merchantID, id)
	if err != nil {
		return nil, err
	}

	messages, err := rm.NewRepo(conn).FindByMerchantAndGroupID(merchantID, messageGroup.ID)
	if err != nil {
		return nil, err
	}

	var views []*v.MessageView
	for _, message := range messages {
		views = append(views, v.NewMessageView(message, nil))
	}

	return v.NewGroupView(messageGroup, views), nil
}

func (s *Service) FindByQuery(p *ig.SearchParams) ([]*models.MessageGroup, error) {
	ctx, cancel := corerepos.DBConnContext()
	defer cancel()

	conn, err := s.app.Config.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	groups, err := rg.NewRepo(conn).FindByQuery(p)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
