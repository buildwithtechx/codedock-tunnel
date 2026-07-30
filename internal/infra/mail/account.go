package mail

import (
	"context"
	"fmt"
)

type AccountMailer struct{ client *ZeptoClient }

func NewAccountMailer(client *ZeptoClient) (*AccountMailer, error) {
	if client == nil {
		return nil, fmt.Errorf("zepto client is required")
	}
	return &AccountMailer{client: client}, nil
}

func (m *AccountMailer) SendAccountUpdate(ctx context.Context, email, event string) error {
	return m.client.Send(ctx, Message{To: email, Subject: "Tunnel account update", Text: "Your tunnel account was " + event + "."})
}
