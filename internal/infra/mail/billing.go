package mail

import (
	"context"
	"fmt"
)

type BillingMailer struct {
	client  *ZeptoClient
	resolve func(context.Context, string) (string, error)
}

func NewBillingMailer(client *ZeptoClient, resolve func(context.Context, string) (string, error)) (*BillingMailer, error) {
	if client == nil || resolve == nil {
		return nil, fmt.Errorf("zepto client and recipient resolver are required")
	}
	return &BillingMailer{client: client, resolve: resolve}, nil
}

func (m *BillingMailer) SendBillingUpdate(ctx context.Context, organizationID, status string) error {
	to, err := m.resolve(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("resolve billing recipient: %w", err)
	}
	return m.client.Send(ctx, Message{To: to, Subject: "Tunnel subscription update", Text: "Your tunnel subscription status is now " + status + "."})
}
