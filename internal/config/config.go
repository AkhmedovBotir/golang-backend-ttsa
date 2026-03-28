package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	JWTSecret        string
	JWTExpireHours   int
	GeneralAdminName string
	GeneralAdminRole string
	GeneralAdminPhone string
	GeneralAdminUsername string
	GeneralAdminPassword string
	GeneralAdminStatus string
}

func Load() Config {
	_ = godotenv.Load()

	expire, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	if err != nil {
		log.Println("JWT_EXPIRE_HOURS noto'g'ri, default 24 ishlatildi")
		expire = 24
	}

	return Config{
		AppPort:             getEnv("APP_PORT", "8081"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "123456"),
		DBName:              getEnv("DB_NAME", "admin_db"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
		JWTSecret:           getEnv("JWT_SECRET", "very-secret-key"),
		JWTExpireHours:      expire,
		GeneralAdminName:    getEnv("GENERAL_ADMIN_NAME", "Super Admin"),
		GeneralAdminRole:    getEnv("GENERAL_ADMIN_ROLE", "general"),
		GeneralAdminPhone:   getEnv("GENERAL_ADMIN_PHONE", "+998901234567"),
		GeneralAdminUsername: getEnv("GENERAL_ADMIN_USERNAME", "general_admin"),
		GeneralAdminPassword: getEnv("GENERAL_ADMIN_PASSWORD", "12345678"),
		GeneralAdminStatus:  getEnv("GENERAL_ADMIN_STATUS", "active"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
