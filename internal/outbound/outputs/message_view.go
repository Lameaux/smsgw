package outputs

import (
	"fmt"
	"github.com/Lameaux/smsgw/internal/outbound/models"
)

type MessageView struct {
	models.Message
	MessageGroup *GroupView `json:"message_group,omitempty"`
	HREF         string     `json:"href"`
}

func NewMessageView(m *models.Message, messageGroup *models.MessageGroup) *MessageView {
	var messageGroupView *GroupView

	if messageGroup != nil {
		messageGroupView = NewGroupView(messageGroup, nil)
	}

	return &MessageView{
		Message:      *m,
		MessageGroup: messageGroupView,
		HREF:         fmt.Sprintf("/messages/outbound/%s", m.ID),
	}
}
