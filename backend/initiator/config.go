package initiator

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DbUrl       string
	FrontendURL string
	JWT         JWTConfig
	SMTP        SMTPConfig
	SMS         SMSConfig
	S3          S3Config
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
		Port:        Getenv("PORT", "8080"),
		DbUrl:       Getenv("DATABASE_URL", ""),
		FrontendURL: Getenv("FRONTEND_URL", "http://localhost:3000"),
		JWT:         loadJWTConfig(),
		SMTP:        loadSMTPConfig(),
		SMS:         loadSMSConfig(),
		S3:          loadS3Config(),
	}

	if cfg.DbUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return cfg
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

func loadJWTConfig() JWTConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return JWTConfig{
		SecretKey:                 secret,
		AccessTokenTTL:            parseDuration(Getenv("JWT_ACCESS_TOKEN_TTL", "5m"), 5*time.Minute),
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
