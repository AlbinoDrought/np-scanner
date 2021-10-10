package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type Sink interface {
	Send(ctx context.Context, notifiable Notifiable) error
}

type nilSink struct{}

func (s *nilSink) Send(ctx context.Context, notifiable Notifiable) error {
	return nil
}

func NewNilSink() Sink {
	return &nilSink{}
}

type discordWebhookSink struct {
	url    string
	client *http.Client
}

type discordWebhook struct {
	Content string `json:"content"`
}

var ErrBadResponse = errors.New("bad response")

func (s *discordWebhookSink) Send(ctx context.Context, notifiable Notifiable) error {
	webhook := discordWebhook{
		Content: notifiable.Message(),
	}
	jsonBytes, err := json.Marshal(&webhook)
	if err != nil {
		return err
	}

	r, err := http.NewRequestWithContext(ctx, "POST", s.url, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(r)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		resp.Body.Close()
	}
	if resp.StatusCode >= 400 {
		return ErrBadResponse
	}
	return nil
}

func NewDiscordWebhookSink(url string, client *http.Client) Sink {
	return &discordWebhookSink{url, client}
}
