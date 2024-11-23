package service

type MessageService interface {
	SendMessage(message string) error
}
