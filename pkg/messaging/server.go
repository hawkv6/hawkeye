package messaging

const Subsystem = "messaging"

type MessagingServer interface {
	Start() error
	Stop()
}
