package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	Port    string
}

var Cfg Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file not found, using default/system env")
	}

	Cfg = Config{
		AppName: getEnv("APP_NAME", "SumunarPOS-CORE"),
		Port:    getEnv("APP_PORT", "1323"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
