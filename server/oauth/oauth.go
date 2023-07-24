package oauth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	PAYPAL_API_BASE = "https://api.sandbox.paypal.com"                                                   // Replace with the actual PayPal API base URL
	CLIENT_ID       = "AfO8JyqOwNtRMq-3X9jr583UkVF10hxeG9Ifku6354w4Xh6eNOSClKl_6lLGi8FEDxseWsDwd9TdmGFG" // Replace with your PayPal client ID
	APP_SECRET      = "EBv2wZGG1L46fHNC7AZJcq_De-OqJrEjRarQeBiLnHBIoILGiNfgphPnxrwMCNIVSj_xpMQed1bHpVMI" // Replace with your PayPal app secret
)

func GetAccessToken() (string, error) {
	credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", CLIENT_ID, APP_SECRET)))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/oauth2/token", PAYPAL_API_BASE), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", credentials))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// func main() {
// 	accessToken, err := getAccessToken()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Access Token:", accessToken)
// }
