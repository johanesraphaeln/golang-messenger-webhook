package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PageAccessToken string
	VerifyToken     string
	PageID          string
)

func LoadConfig() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	PageAccessToken = os.Getenv("PAGE_ACCESS_TOKEN")
	VerifyToken = os.Getenv("VERIFY_TOKEN")
	PageID = os.Getenv("PAGE_ID")

	if PageAccessToken == "" || VerifyToken == "" || PageID == "" {
		log.Fatal("ENV is not set")
	}
}
