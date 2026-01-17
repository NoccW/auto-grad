package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	JWTSecret  string
	UploadPath string

	// API Keys
	BaiduAPIKey    string
	BaiduSecretKey string
	DeepSeekAPIKey string
}

func LoadConfig() *Config {
	// 尝试加载本地 .env，环境变量优先生效
	if err := godotenv.Load(".env"); err != nil {
		// 不阻断启动，保持可用
		log.Printf("warning: could not load .env file: %v", err)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "Root123456!"),
		DBName:     getEnv("DB_NAME", "auto_grad_web"),
		ServerPort: getEnv("SERVER_PORT", "3000"),
		JWTSecret:  getEnv("JWT_SECRET", "your-jwt-secret-key-change-in-production"),
		UploadPath: getEnv("UPLOAD_PATH", "./uploads"),

		BaiduAPIKey:    getEnv("BAIDU_API_KEY", ""),
		BaiduSecretKey: getEnv("BAIDU_SECRET_KEY", ""),
		DeepSeekAPIKey: getEnv("DEEPSEEK_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
