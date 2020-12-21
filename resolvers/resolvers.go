package resolvers

import (
	"github.com/leandro-lugaresi/hub"
	"github.com/muety/webhook2telegram/config"
	"github.com/muety/webhook2telegram/model"
)

const (
	TextType = "TEXT"
	FileType = "FILE"
)

type MessageResolver interface {
	IsValid(*model.DefaultMessage) error
	Resolve(string, *model.DefaultMessage, *model.MessageParams) error
}

func GetResolver(ttype string) MessageResolver {
	switch ttype {
	case FileType:
		return &FileResolver{}
	}
	return &TextResolver{}
}

func logMessage(m *model.DefaultMessage) {
	ttype := m.Type
	if ttype == "" {
		ttype = TextType
	}

	config.GetHub().Publish(hub.Message{
		Name: config.EventOnMessageDelivered,
		Fields: map[string]interface{}{
			"origin": m.Origin,
			"type":   ttype,
		},
	})
}
