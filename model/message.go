package model

const DefaultOrigin = "Telepush"

type Message struct {
	Text     string `json:"text"`
	Origin   string `json:"origin"`
	File     string `json:"file"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
}

type MessageWithOptions struct {
	Message `mapstructure:",squash"`
	Options MessageOptions `json:"options" mapstructure:",squash"`
}

type MessageOptions struct {
	DisableLinkPreviews bool `json:"disable_link_previews" mapstructure:"disable_link_previews"`
	DisableMarkdown     bool `json:"disable_markdown" mapstructure:"disable_markdown"`
}

func (o *MessageOptions) ParseMode() string {
	if o.DisableMarkdown {
		return ""
	}
	return "Markdown"
}
