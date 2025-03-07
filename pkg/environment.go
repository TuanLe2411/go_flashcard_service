package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	var envFile string
	switch env {
	case "production":
		envFile = ".env.production"
	case "development":
		envFile = ".env.development"
	default:
		log.Fatalf("ENV không hợp lệ: %s. Chỉ hỗ trợ 'development' hoặc 'production'", env)
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Lỗi khi load file %s: %v", envFile, err)
	}
}
