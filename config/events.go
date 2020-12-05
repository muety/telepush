package config

const (
	EventOnRequestFailed     = "on_request_failed"
	EventOnRequestSuccessful = "on_request_successful"
	EventOnMessageDelivered  = "on_message_delivered"
)

const (
	FieldMessageOrigin = "origin"
	FieldMessageType   = "type"
)

func AllEvents() []string {
	return []string{
		EventOnRequestFailed,
		EventOnRequestSuccessful,
		EventOnMessageDelivered,
	}
}
