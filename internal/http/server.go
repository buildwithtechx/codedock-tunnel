package http

import (
	"fmt"
	"strings"

	"codedock.run/codedock-tunnel/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	app *fiber.App
}

func NewServer(cfg config.APIConfig, deps Dependencies) (*Server, error) {
	deps.PublicAPIURL = cfg.App.PublicAPIURL
	handlers, err := buildHandlers(deps, cfg.Auth.CookieName, cfg.Auth.CookieSecure)
	if err != nil {
		return nil, err
	}
	app := fiber.New(fiber.Config{AppName: cfg.App.Name, DisableStartupMessage: true, ErrorHandler: errorHandler})
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{AllowOrigins: cfg.App.AllowedOrigins, AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Internal-Secret", AllowCredentials: true}))
	if err := RegisterRoutes(app, handlers, RouterOptions{CookieName: cfg.Auth.CookieName, CookieSecure: cfg.Auth.CookieSecure, InternalAPISecret: cfg.Service.InternalAPISecret, BillingWebhookSecret: cfg.Billing.WebhookSecret}); err != nil {
		return nil, err
	}
	return &Server{app: app}, nil
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) Listen(address string) error {
	if s == nil || s.app == nil {
		return fmt.Errorf("http server is not initialized")
	}
	return s.app.Listen(address)
}

func (s *Server) Shutdown() error {
	if s == nil || s.app == nil {
		return nil
	}
	return s.app.Shutdown()
}

func errorHandler(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	if fiberErr, ok := err.(*fiber.Error); ok {
		status = fiberErr.Code
	}
	message := strings.TrimSpace(err.Error())
	if status >= fiber.StatusInternalServerError {
		message = "internal server error"
	}
	return c.Status(status).JSON(fiber.Map{"error": message})
}
