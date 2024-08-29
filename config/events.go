package config

const (
	EventOnRequestFailed     = "on_request_failed"
	EventOnRequestSuccessful = "on_request_successful"
	EventOnMessageDelivered  = "on_message_delivered"
	EventOnTokenIssued       = "on_user_registered"
)

const (
	FieldMessageOrigin    = "origin"
	FieldMessageType      = "type"
	FieldMessageRecipient = "chat_id"
	FieldTokenToken       = "token"
	FieldTokenChat        = "chat_id"
	FieldTokenUser        = "user_id"
)

func AllEvents() []string {
	return []string{
		EventOnRequestFailed,
		EventOnRequestSuccessful,
		EventOnMessageDelivered,
		EventOnTokenIssued,
	}
}
