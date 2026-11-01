package initiator

import (
	"log"
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
	AppEnv      string
	CORSOrigins []string
	// MigrateOnStart applies pending migrations before serving (default true).
	MigrateOnStart bool
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
	SecretKey                 string
	AccessTokenTTL            time.Duration
	RefreshTokenTTL           time.Duration
	RememberMeRefreshTokenTTL time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:           Getenv("PORT", "8080"),
		DbUrl:          Getenv("DATABASE_URL", ""),
		FrontendURL:    Getenv("FRONTEND_URL", "http://localhost:3000"),
		AppEnv:         Getenv("APP_ENV", "production"),
		CORSOrigins:    parseList(Getenv("CORS_ORIGINS", defaultCORSOrigins)),
		MigrateOnStart: Getenv("MIGRATE_ON_START", "true") == "true",
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
		SecretKey:                 secret,
		AccessTokenTTL:            parseDuration(Getenv("JWT_ACCESS_TOKEN_TTL", "15m"), 15*time.Minute),
		RefreshTokenTTL:           parseDuration(Getenv("JWT_REFRESH_TOKEN_TTL", "168h"), 7*24*time.Hour),
		RememberMeRefreshTokenTTL: parseDuration(Getenv("JWT_REMEMBER_ME_REFRESH_TTL", "720h"), 30*24*time.Hour),
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
