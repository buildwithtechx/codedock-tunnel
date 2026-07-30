package config

import "time"

type APIConfig struct {
	App      AppConfig      `envPrefix:"CODEDOCK_"`
	Auth     AuthConfig     `envPrefix:"CODEDOCK_"`
	Database DatabaseConfig `envPrefix:"CODEDOCK_"`
	Redis    RedisConfig    `envPrefix:"CODEDOCK_"`
	Mail     MailConfig     `envPrefix:"CODEDOCK_"`
	Service  ServiceConfig  `envPrefix:"CODEDOCK_"`
	Billing  BillingConfig  `envPrefix:"CODEDOCK_"`
	Tunnel   TunnelConfig   `envPrefix:"CODEDOCK_"`
}

type RelayConfig struct {
	App     AppConfig     `envPrefix:"CODEDOCK_"`
	Redis   RedisConfig   `envPrefix:"CODEDOCK_"`
	Tunnel  TunnelConfig  `envPrefix:"CODEDOCK_"`
	Service ServiceConfig `envPrefix:"CODEDOCK_"`
	RelayID string        `env:"RELAY_ID"`
}

type CronConfig struct {
	App      AppConfig      `envPrefix:"CODEDOCK_"`
	Database DatabaseConfig `envPrefix:"CODEDOCK_"`
	Redis    RedisConfig    `envPrefix:"CODEDOCK_"`
	Service  ServiceConfig  `envPrefix:"CODEDOCK_"`
}

type CheckConfig struct {
	App     AppConfig     `envPrefix:"CODEDOCK_"`
	Service ServiceConfig `envPrefix:"CODEDOCK_"`
}

type CLIConfig struct {
	APIURL     string `env:"CODEDOCK_TUNNEL_API_URL" envDefault:"http://localhost:8080"`
	RelayURL   string `env:"CODEDOCK_TUNNEL_RELAY_URL" envDefault:"ws://localhost:8081"`
	PublicDomain string `env:"CODEDOCK_TUNNEL_DOMAIN" envDefault:"tunnel.codedock-tunnel.dev"`
	APIKey     string `env:"CODEDOCK_TUNNEL_API_KEY"`
	AgentToken string `env:"CODEDOCK_TUNNEL_AGENT_TOKEN"`
	Password   string `env:"CODEDOCK_TUNNEL_PASSWORD"`
	ConfigPath string `env:"CODEDOCK_TUNNEL_CONFIG_PATH" envDefault:".config/codedock-tunnel/config.json"`
}

type ServiceConfig struct {
	InternalAPIURL    string `env:"INTERNAL_API_URL" envDefault:"http://localhost:8080"`
	InternalAPISecret string `env:"INTERNAL_API_SECRET"`
}

type AppConfig struct {
	Port             string        `env:"PORT" envDefault:"8080"`
	Name             string        `env:"APP_NAME" envDefault:"codedock-tunnel"`
	Environment      string        `env:"ENV" envDefault:"development"`
	LogLevel         string        `env:"LOG_LEVEL" envDefault:"info"`
	ShutdownTimeout  time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	AllowedOrigins   string        `env:"ALLOWED_ORIGINS" envDefault:"http://localhost:3000,http://localhost:3001"`
	CORSOrigin       string        `env:"CORS_ORIGIN" envDefault:"http://localhost:3000"`
	PublicAPIURL     string        `env:"PUBLIC_API_URL" envDefault:"http://localhost:8080"`
	DashboardURL     string        `env:"DASHBOARD_URL" envDefault:"http://localhost:3000"`
	ACMEEmail        string        `env:"ACME_EMAIL"`
	ACMEDirectory    string        `env:"ACME_DIRECTORY"`
	CertificateCache string        `env:"CERTIFICATE_CACHE_DIR" envDefault:".data/acme"`
	RequireTLS       bool          `env:"REQUIRE_TLS" envDefault:"false"`
	TLSCertFile      string        `env:"TLS_CERT_FILE"`
	TLSKeyFile       string        `env:"TLS_KEY_FILE"`
}

type AuthConfig struct {
	SessionTTL         time.Duration `env:"SESSION_TTL" envDefault:"720h"`
	DeviceLoginTTL     time.Duration `env:"DEVICE_LOGIN_TTL" envDefault:"10m"`
	InvitationTTL      time.Duration `env:"INVITATION_TTL" envDefault:"168h"`
	OAuthStateTTL      time.Duration `env:"OAUTH_STATE_TTL" envDefault:"10m"`
	CookieName         string        `env:"AUTH_COOKIE_NAME" envDefault:"codedock_session"`
	CookieSecure       bool          `env:"AUTH_COOKIE_SECURE" envDefault:"false"`
	GoogleClientID     string        `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string        `env:"GOOGLE_CLIENT_SECRET"`
	GitHubClientID     string        `env:"GITHUB_CLIENT_ID"`
	GitHubClientSecret string        `env:"GITHUB_CLIENT_SECRET"`
	EncryptionKey      string        `env:"AUTH_ENCRYPTION_KEY"`
}

type DatabaseConfig struct {
	URL         string        `env:"DATABASE_URL" envDefault:"postgres://codedock:codedock@localhost:5432/codedock?sslmode=disable"`
	MaxConns    int           `env:"DATABASE_MAX_CONNS" envDefault:"25"`
	MaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"30m"`
	MaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"5m"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     string `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type MailConfig struct {
	FromAddress string `env:"MAIL_FROM" envDefault:"noreply@localhost"`
	ZeptoAPIKey string `env:"ZEPTO_API_KEY"`
	ZeptoURL    string `env:"ZEPTO_URL" envDefault:"https://api.zeptomail.com/v1.1/email"`
}

type TunnelConfig struct {
	Domain          string        `env:"TUNNEL_DOMAIN" envDefault:"tunnel.codedock-tunnel.dev"`
	TokenTTL        time.Duration `env:"TUNNEL_TOKEN_TTL" envDefault:"24h"`
	MaxConnections  int           `env:"TUNNEL_MAX_CONNECTIONS" envDefault:"1000"`
	MaxTunnels      int           `env:"TUNNEL_MAX_TUNNELS" envDefault:"1000"`
	MaxBytes        int64         `env:"TUNNEL_MAX_BYTES" envDefault:"0"`
	MaxBandwidth    int64         `env:"TUNNEL_MAX_BANDWIDTH_BYTES" envDefault:"0"`
	RequireTLS      bool          `env:"TUNNEL_REQUIRE_TLS" envDefault:"false"`
	AgentInactivity time.Duration `env:"AGENT_INACTIVITY_TIMEOUT" envDefault:"90s"`
	Heartbeat       time.Duration `env:"TUNNEL_HEARTBEAT_INTERVAL" envDefault:"20s"`
	ReadTimeout     time.Duration `env:"TUNNEL_READ_TIMEOUT" envDefault:"90s"`
	DrainTimeout    time.Duration `env:"TUNNEL_DRAIN_TIMEOUT" envDefault:"10s"`
	MaxFrameBytes   int64         `env:"TUNNEL_MAX_FRAME_BYTES" envDefault:"16777216"`
}

type BillingConfig struct {
	GracePeriod             time.Duration `env:"BILLING_GRACE_PERIOD" envDefault:"72h"`
	PolarServer             string        `env:"POLAR_SERVER" envDefault:"sandbox"`
	PolarBaseURL            string        `env:"POLAR_BASE_URL" envDefault:"https://sandbox-api.polar.sh"`
	PolarAccessToken        string        `env:"POLAR_ACCESS_TOKEN"`
	PolarWebhookSecret      string        `env:"POLAR_WEBHOOK_SECRET"`
	PolarProductRay         string        `env:"POLAR_PRODUCT_RAY"`
	PolarProductBeam        string        `env:"POLAR_PRODUCT_BEAM"`
	PolarProductPulse       string        `env:"POLAR_PRODUCT_PULSE"`
	PolarProductRayYearly   string        `env:"POLAR_PRODUCT_RAY_YEARLY"`
	PolarProductBeamYearly  string        `env:"POLAR_PRODUCT_BEAM_YEARLY"`
	PolarProductPulseYearly string        `env:"POLAR_PRODUCT_PULSE_YEARLY"`
	PaystackBaseURL         string        `env:"PAYSTACK_BASE_URL" envDefault:"https://api.paystack.co"`
	PaystackSecret          string        `env:"PAYSTACK_SECRET_KEY"`
	WebhookSecret           string        `env:"BILLING_WEBHOOK_SECRET"`
}
