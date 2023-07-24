package main

import "os"

var (
	NODE_ENV   = os.Getenv("NODE_ENV")
	CLIENT_ID  = os.Getenv("CLIENT_ID")
	APP_SECRET = os.Getenv("APP_SECRET")
)

func main() {
	isProd := NODE_ENV == "production"

	var PAYPAL_API_BASE string
	if isProd {
		PAYPAL_API_BASE = "https://api.paypal.com"
	} else {
		PAYPAL_API_BASE = "https://api.sandbox.paypal.com"
	}

	// Now you can use isProd, PAYPAL_API_BASE, CLIENT_ID, and APP_SECRET in your code as needed.
	// For example, you can print them:
	println("isProd:", isProd)
	println("PAYPAL_API_BASE:", PAYPAL_API_BASE)
	println("CLIENT_ID:", CLIENT_ID)
	println("APP_SECRET:", APP_SECRET)
}
