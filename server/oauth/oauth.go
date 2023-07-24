package oauth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	PAYPAL_API_BASE = "https://api.paypal.com" // Replace this with your actual API base URL
	CLIENT_ID       = "your_client_id"         // Replace this with your actual client ID
	APP_SECRET      = "your_app_secret"        // Replace this with your actual app secret
)

func getAccessToken() (map[string]interface{}, error) {
	credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", CLIENT_ID, APP_SECRET)))

	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/oauth2/token", PAYPAL_API_BASE), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", credentials))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func ReadAccessToken() {
	accessToken, err := getAccessToken()
	if err != nil {
		fmt.Println("Error getting access token:", err)
		return
	}

	fmt.Println("Access Token:", accessToken)
}
