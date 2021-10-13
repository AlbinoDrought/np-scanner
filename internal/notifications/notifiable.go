package notifications

type Notifiable interface {
	ID() string
	Message() string
}

type DiscordNotifiable interface {
	DiscordMessage() string
}
