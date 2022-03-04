package webmentionio

type WebmentionMessage struct {
	Source string `json:"source" binding:"required"`
	Target string `json:"target" binding:"required"`
}
