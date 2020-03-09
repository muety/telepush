package model

type StoreObject struct {
	User   TelegramUser `json:"user"`
	ChatId int          `json:"chat_id"`
}

type DefaultMessage struct {
	RecipientToken string `json:"recipient_token"`
	Text           string `json:"text"`
	Origin         string `json:"origin"`
	File           string `json:"file"`
	Filename       string `json:"filename"`
	Type           string `json:"type"`
}

type MessageParams struct {
	DisableLinkPreviews bool
}

type AlertmanagerMessage struct {
	Alerts []*AlertmanagerAlert
}

type AlertmanagerAlert struct {
	Status      string            `json:"status"`
	Url         string            `json:"generatorURL"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// Only required fields are implemented
type TelegramUser struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Only required fields are implemented
type TelegramChat struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}

// Only required fields are implemented
type TelegramOutMessage struct {
	ChatId             string `json:"chat_id"`
	Text               string `json:"text"`
	ParseMode          string `json:"parse_mode"`
	DisableLinkPreview bool   `json:"disable_web_page_preview"`
}

// Only required fields are implemented
type TelegramInMessage struct {
	MessageId int          `json:"message_id"`
	From      TelegramUser `json:"from"`
	Date      int          `json:"date"`
	Chat      TelegramChat `json:"chat"`
	Text      string       `json:"text"`
}

// Only required fields are implemented
type TelegramUpdate struct {
	UpdateId int               `json:"update_id"`
	Message  TelegramInMessage `json:"message"`
}

type TelegramUpdateResponse struct {
	Ok     bool             `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}

type Stats struct {
	TotalRequests int `json:"total_requests"`
	Timestamp     int `json:"timestamp"`
}
