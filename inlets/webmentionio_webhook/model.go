package webmentionio_webhook

type WebmentionMessage struct {
	Secret string `json:"secret" binding:"required"`
	Source string `json:"source" binding:"required"`
	Target string `json:"target" binding:"required"`
}
