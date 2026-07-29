package handlers

import (
	"fmt"
	"strings"

	"codedock.run/codedock-tunnel/internal/services"
	"github.com/gofiber/fiber/v2"
)

type OAuthHandler struct {
	oauth        *services.OAuthService
	publicAPIURL string
	cookieName   string
	cookieSecure bool
}

func NewOAuthHandler(oauth *services.OAuthService, publicAPIURL, cookieName string, cookieSecure bool) (*OAuthHandler, error) {
	if oauth == nil || strings.TrimSpace(publicAPIURL) == "" || strings.TrimSpace(cookieName) == "" {
		return nil, fmt.Errorf("oauth service, public api url, and cookie name are required")
	}
	return &OAuthHandler{oauth: oauth, publicAPIURL: strings.TrimRight(publicAPIURL, "/"), cookieName: cookieName, cookieSecure: cookieSecure}, nil
}

func (h *OAuthHandler) Start(c *fiber.Ctx) error {
	provider := c.Params("provider")
	redirectURI := h.publicAPIURL + "/api/v1/auth/oauth/" + provider + "/callback"
	url, err := h.oauth.Start(c.UserContext(), provider, redirectURI)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.Redirect(url, fiber.StatusFound)
}

func (h *OAuthHandler) Callback(c *fiber.Ctx) error {
	raw, session, err := h.oauth.Callback(c.UserContext(), c.Query("state"), c.Query("code"), c.Get("User-Agent"), c.IP())
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err)
	}
	c.Cookie(&fiber.Cookie{Name: h.cookieName, Value: raw, HTTPOnly: true, Secure: h.cookieSecure, SameSite: "Lax", Path: "/", Expires: session.ExpiresAt})
	return c.JSON(fiber.Map{"userId": session.UserID, "expiresAt": session.ExpiresAt})
}
