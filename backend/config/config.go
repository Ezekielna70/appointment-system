package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	DBName   string
	JWTSecret string
}

func LoadConfig() Config {
	godotenv.Load() 

	return Config{
		MongoURI:  os.Getenv("MONGO_URI"),
		DBName:    os.Getenv("MONGO_DB_NAME"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}