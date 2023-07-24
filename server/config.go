package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	isProd        bool
	paypalAPIBase string
	clientID      string
	appSecret     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	isProd = os.Getenv("NODE_ENV") == "development"
	paypalAPIBase = getPaypalAPIBase()
	clientID = os.Getenv("CLIENT_ID")
	appSecret = os.Getenv("APP_SECRET")
}

func getPaypalAPIBase() string {
	if isProd {
		return "https://api.paypal.com"
	}
	return "https://api.sandbox.paypal.com"
}

func Config() {
	fmt.Println("isProd:", isProd)
	fmt.Println("PAYPAL_API_BASE:", paypalAPIBase)
	fmt.Println("CLIENT_ID:", clientID)
	fmt.Println("APP_SECRET:", appSecret)
}
