package outputs

import (
	"fmt"
	"github.com/Lameaux/smsgw/internal/outbound/models"
)

type GroupView struct {
	models.MessageGroup
	Messages []*MessageView
	HREF     string `json:"href"`
}

func NewGroupView(messageGroup *models.MessageGroup, messages []*MessageView) *GroupView {
	return &GroupView{
		MessageGroup: *messageGroup,
		Messages:     messages,
		HREF:         fmt.Sprintf("/messages/group/%s", messageGroup.ID),
	}
}
