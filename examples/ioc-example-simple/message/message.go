package message

type MessageService interface {
	SendMessage(msg string) string
}
