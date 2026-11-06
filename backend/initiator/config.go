package initiator

import (
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/secrets"
	"adinfluence_service_app/platform/utils"
	"context"
	"log"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DbUrl       string
	FrontendURL string
	// ContactInbox receives messages from the public Contact us page.
	ContactInbox string
	AppEnv       string
	CORSOrigins  []string
	// MigrateOnStart applies pending migrations before serving (default true).
	MigrateOnStart bool
	// RedisURL holds the login attempt counters (e.g. redis://redis:6379/0).
	// Empty keeps them in this process's memory.
	RedisURL string
	// TrustedProxies may set X-Real-IP / X-Forwarded-For (CIDRs).
	TrustedProxies []netip.Prefix
	DB             DBConfig
	JWT            JWTConfig
	SMTP           SMTPConfig
	SMS            SMSConfig
	S3             S3Config
}

type S3Config struct {
	Endpoint     string
	Region       string
	AccessKeyID  string
	SecretKey    string
	Bucket       string
	PublicURL    string
	UsePathStyle bool
}

// DBConfig bounds the connection pool so a traffic spike queues requests
// instead of exhausting Postgres max_connections.
type DBConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type SMSConfig struct {
	Username string
	APIKey   string
	From     string
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

type JWTConfig struct {
	SecretKey      string
	AccessTokenTTL time.Duration
	Sessions       jwt.SessionPolicy
}

// LoadVault runs first at startup: it reads .env, then, when VAULT_ADDR is
// set, applies every key from Vault over it. Without Vault access the API
// refuses to start: running on stale or missing secrets would be worse.
// Docker restarts it, so it comes up by itself once Vault is unsealed.
func LoadVault() *secrets.Vault {
	_ = godotenv.Load()

	vault, err := secrets.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if vault == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	names, err := vault.Apply(ctx)
	if err != nil {
		log.Fatalf("cannot load configuration from Vault %s: %v", vault.Path(), err)
	}
	// Names only, never values.
	log.Printf("loaded %d keys from Vault %s: %s", len(names), vault.Path(), strings.Join(names, ", "))
	return vault
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:           Getenv("PORT", "8080"),
		DbUrl:          Getenv("DATABASE_URL", ""),
		FrontendURL:    Getenv("FRONTEND_URL", "http://localhost:3000"),
		ContactInbox:   Getenv("CONTACT_INBOX_EMAIL", Getenv("SMTP_FROM", os.Getenv("SMTP_USER"))),
		AppEnv:         Getenv("APP_ENV", "production"),
		CORSOrigins:    parseList(Getenv("CORS_ORIGINS", defaultCORSOrigins)),
		MigrateOnStart: Getenv("MIGRATE_ON_START", "true") == "true",
		RedisURL:       Getenv("REDIS_URL", ""),
		TrustedProxies: utils.ParseTrustedProxies(Getenv("TRUSTED_PROXIES", utils.DefaultTrustedProxies)),
		DB:             loadDBConfig(),
		JWT:            loadJWTConfig(),
		SMTP:           loadSMTPConfig(),
		SMS:            loadSMSConfig(),
		S3:             loadS3Config(),
	}

	if cfg.DbUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWT.SecretKey == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if len(cfg.JWT.SecretKey) < minJWTSecretLength {
		log.Printf("WARNING: JWT_SECRET is shorter than %d bytes; use a longer random secret", minJWTSecretLength)
	}

	return cfg
}

// IsDevelopment reports whether the service runs in local development mode,
// where secrets such as OTPs may be written to logs for testing.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development" || c.AppEnv == "dev" || c.AppEnv == "local"
}

const (
	minJWTSecretLength = 32
	defaultCORSOrigins = "http://localhost:3000,http://127.0.0.1:3000,https://adfluence-habesha.vercel.app"
)

func parseList(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, strings.TrimRight(p, "/"))
		}
	}
	return out
}

func loadSMTPConfig() SMTPConfig {
	user := os.Getenv("SMTP_USER")
	return SMTPConfig{
		Host:     Getenv("SMTP_HOST", "smtp.gmail.com"),
		Port:     Getenv("SMTP_PORT", "587"),
		User:     user,
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     Getenv("SMTP_FROM", user),
	}
}

func loadSMSConfig() SMSConfig {
	return SMSConfig{
		Username: os.Getenv("SMS_USERNAME"),
		APIKey:   os.Getenv("SMS_API_KEY"),
		From:     os.Getenv("SMS_FROM"),
	}
}

func loadS3Config() S3Config {
	return S3Config{
		Endpoint:     os.Getenv("S3_ENDPOINT"),
		Region:       Getenv("S3_REGION", "us-east-1"),
		AccessKeyID:  os.Getenv("S3_ACCESS_KEY_ID"),
		SecretKey:    os.Getenv("S3_SECRET_ACCESS_KEY"),
		Bucket:       os.Getenv("S3_BUCKET_NAME"),
		PublicURL:    os.Getenv("S3_PUBLIC_URL"),
		UsePathStyle: Getenv("S3_USE_PATH_STYLE", "true") == "true",
	}
}

func loadDBConfig() DBConfig {
	return DBConfig{
		MaxOpenConns:    parseInt(Getenv("DB_MAX_OPEN_CONNS", "25"), 25),
		MaxIdleConns:    parseInt(Getenv("DB_MAX_IDLE_CONNS", "10"), 10),
		ConnMaxLifetime: parseDuration(Getenv("DB_CONN_MAX_LIFETIME", "30m"), 30*time.Minute),
		ConnMaxIdleTime: parseDuration(Getenv("DB_CONN_MAX_IDLE_TIME", "5m"), 5*time.Minute),
	}
}

func loadJWTConfig() JWTConfig {
	secret := os.Getenv("JWT_SECRET")

	return JWTConfig{
		SecretKey: secret,
		// Short-lived access tokens; sessions end after inactivity (bank-style).
		AccessTokenTTL: parseDuration(Getenv("JWT_ACCESS_TOKEN_TTL", "5m"), 5*time.Minute),
		Sessions: jwt.SessionPolicy{
			Idle:        parseDuration(Getenv("SESSION_IDLE_TIMEOUT", "15m"), 15*time.Minute),
			MaxAge:      parseDuration(Getenv("SESSION_MAX_AGE", "12h"), 12*time.Hour),
			AdminIdle:   parseDuration(Getenv("SESSION_ADMIN_IDLE_TIMEOUT", "10m"), 10*time.Minute),
			AdminMaxAge: parseDuration(Getenv("SESSION_ADMIN_MAX_AGE", "8h"), 8*time.Hour),
		},
	}
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func parseDuration(raw string, fallback time.Duration) time.Duration {
	if raw == "" {
		return fallback
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return d
}

func parseInt(raw string, fallback int) int {
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
