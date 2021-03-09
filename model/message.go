package model

type DefaultMessage struct {
	RecipientToken string `json:"recipient_token"`
	Text           string `json:"text"`
	Origin         string `json:"origin"`
	File           string `json:"file"`
	Filename       string `json:"filename"`
	Type           string `json:"type"`
}

type ExtendedMessage struct {
	DefaultMessage
	Options *MessageParams `json:"options"`
}

type MessageParams struct {
	DisableLinkPreviews bool `json:"disable_link_previews"`
}
