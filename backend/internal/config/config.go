package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	WxAppID     string
	WxAppSecret string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/flashcard?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "flashcard-dev-secret"),
		WxAppID:     os.Getenv("WX_APP_ID"),
		WxAppSecret: os.Getenv("WX_APP_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
