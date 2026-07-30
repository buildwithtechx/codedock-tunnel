package mail

import (
	"context"
	"fmt"
)

type BillingMailer struct {
	client    *ZeptoClient
	resolve   func(context.Context, string) (string, error)
	renderer  *templateRenderer
	dashboard string
}

func NewBillingMailer(client *ZeptoClient, resolve func(context.Context, string) (string, error), dashboardURL string) (*BillingMailer, error) {
	if client == nil || resolve == nil {
		return nil, fmt.Errorf("zepto client and recipient resolver are required")
	}
	renderer, err := newTemplateRenderer()
	if err != nil {
		return nil, err
	}
	return &BillingMailer{client: client, resolve: resolve, renderer: renderer, dashboard: dashboardURL}, nil
}

func (m *BillingMailer) SendBillingUpdate(ctx context.Context, organizationID, status string) error {
	to, err := m.resolve(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("resolve billing recipient: %w", err)
	}
	html, err := m.renderer.render("billing-update", BillingUpdateData{Status: status, DashboardURL: m.dashboard})
	if err != nil {
		return err
	}
	return m.client.Send(ctx, Message{To: to, Subject: "Codedock Tunnel subscription update", HTML: html})
}
