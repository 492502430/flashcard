package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	WxAppID       string
	WxAppSecret   string
	AdminPassword string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://localhost:5433/flashcard?sslmode=disable&user=flashcard"),
		JWTSecret:     getEnv("JWT_SECRET", "flashcard-dev-secret"),
		WxAppID:       getEnv("WX_APP_ID", "wxebcc02862ca09ed2"),
		WxAppSecret:   getEnv("WX_APP_SECRET", "22e3e089d7840f9c85dfde0423c5271c"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "flashcard123"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
