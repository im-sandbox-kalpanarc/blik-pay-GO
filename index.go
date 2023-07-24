package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

const paypalApiBase = "https://api.paypal.com"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	clientDir := "../client"
	clientDirAbs, err := filepath.Abs(clientDir)
	if err != nil {
		log.Fatalf("Error resolving client directory: %v", err)
	}

	r.StaticFS("/", http.Dir(clientDirAbs))
	r.POST("/capture/:orderId", captureOrder)
	r.POST("/webhook", webhookHandler)

	url := fmt.Sprintf("http://localhost:%s", port)
	cmd := exec.Command("open", url)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error opening the browser: %v", err)
	}

	fmt.Printf("Example app listening at %s\n", url)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("Error running the server: %v", err)
	}
}

func getAccessToken() (string, error) {
	// Implement your getAccessToken logic here
	// You can use the http package to send requests to get the access token
	return "", nil
}

func captureOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	accessToken, err := getAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", paypalApiBase, orderID)
	reqBody := make(map[string]interface{})

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
		return
	}

	fmt.Println("💰 Payment captured!")
	c.JSON(http.StatusOK, data)
}

func webhookHandler(c *gin.Context) {
	var requestBody map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body"})
		return
	}

	accessToken, err := getAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	eventType, _ := requestBody["event_type"].(string)
	resource := requestBody["resource"].(map[string]interface{})
	orderID := resource["id"].(string)

	fmt.Println("🪝 Received Webhook Event")

	// Verify webhook signature (Implement your verification logic here)

	// Capture the order if event_type is "CHECKOUT.ORDER.APPROVED"
	if eventType == "CHECKOUT.ORDER.APPROVED" {
		url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", paypalApiBase, orderID)
		reqBody := make(map[string]interface{})

		reqBodyJSON, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed"})
			return
		}

		fmt.Println("💰 Payment captured!")
	}

	c.JSON(http.StatusOK, gin.H{})
}
