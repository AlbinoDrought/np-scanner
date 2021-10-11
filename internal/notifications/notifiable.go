package notifications

type Notifiable interface {
	ID() string
	Message() string
}
