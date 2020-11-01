package resolvers

import (
	"github.com/muety/webhook2telegram/model"
)

const (
	TextType = "TEXT"
	FileType = "FILE"
)

type MessageResolver interface {
	IsValid(*model.DefaultMessage) error
	Resolve(string, *model.DefaultMessage, *model.MessageParams) *model.ApiError
}

func GetResolver(ttype string) MessageResolver {
	switch ttype {
	case FileType:
		return &FileResolver{}
	}
	return &TextResolver{}
}
