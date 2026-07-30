package mail

import (
	"strings"
	"testing"
)

func TestTemplateRendererEscapesHTMLData(t *testing.T) {
	renderer, err := newTemplateRenderer()
	if err != nil {
		t.Fatalf("create template renderer: %v", err)
	}
	html, err := renderer.render("account-update", AccountUpdateData{Event: "<disabled>", DashboardURL: "https://localhost/dashboard"})
	if err != nil {
		t.Fatalf("render account update: %v", err)
	}
	if strings.Contains(html, "<disabled>") || !strings.Contains(html, "&lt;disabled&gt;") {
		t.Fatalf("expected escaped event in html: %s", html)
	}
}
