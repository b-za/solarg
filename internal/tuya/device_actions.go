package tuya

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func calculateSignatureWithToken(clientID, timestamp, accessToken, stringToSignForHMAC, secret string) string {
	finalStringToSign := clientID + accessToken + timestamp + stringToSignForHMAC
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(finalStringToSign))
	hashed := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(hashed))
}

func getToken(ClientID, ClientSecret string) string {
	token, err := ReadToken()
	if err != nil {
		log.Println("Could not read existing token, retrieving a new one...")
		_, err := GetTuyaToken(ClientID, ClientSecret)
		if err != nil {
			log.Fatalf("Failed to retrieve a new token: %v", err)
		}
		// Try reading the newly saved token again.
		token, err = ReadToken()
		if err != nil {
			log.Fatalf("Failed to read newly saved token: %v", err)
		}
	} else {
		fmt.Println("Successfully read existing token from file.")
	}

	if token.AccessToken == "" {
		log.Fatal("Token is invalid or empty.")
	}
	return token.AccessToken
}

func GetDeviceStatus(deviceID, clientID, clientSecret string) (DeviceStatusResponse, error) {

	accessToken := getToken(clientID, clientSecret)

	var statusResponse DeviceStatusResponse
	method := "GET"
	path := fmt.Sprintf("/v1.0/devices/%s/status", deviceID)
	url := TuyaBaseURL + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return statusResponse, fmt.Errorf("failed to create request: %v", err)
	}

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	stringToSign := method + "\n" + hex.EncodeToString(sha256.New().Sum(nil)) + "\n" + "" + "\n" + path
	sign := calculateSignatureWithToken(clientID, timestamp, accessToken, stringToSign, clientSecret)

	req.Header.Set("client_id", clientID)
	req.Header.Set("access_token", accessToken)
	req.Header.Set("sign", sign)
	req.Header.Set("t", timestamp)
	req.Header.Set("sign_method", "HMAC-SHA256")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return statusResponse, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return statusResponse, fmt.Errorf("failed to read response body: %v", err)
	}

	if err := json.Unmarshal(body, &statusResponse); err != nil {
		return statusResponse, fmt.Errorf("failed to unmarshal device status response: %v. Body: %s", err, string(body))
	}

	if !statusResponse.Success {
		return statusResponse, fmt.Errorf("API indicated failure to get status: %s", string(body))
	}

	return statusResponse, nil
}

func SetSwitchState(deviceID, clientID, clientSecret string, turnOn bool) (string, error) {

	accessToken := getToken(clientID, clientSecret)

	const switchCode = "switch"

	log.Printf("Sending command to set switch '%s' to state: %t", switchCode, turnOn)

	commandPayload := map[string]interface{}{
		"commands": []map[string]interface{}{
			{
				"code":  switchCode,
				"value": turnOn,
			},
		},
	}
	commandBody, err := json.Marshal(commandPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal command payload: %v", err)
	}

	method := "POST"
	path := fmt.Sprintf("/v1.0/devices/%s/commands", deviceID)
	url := TuyaBaseURL + path

	req, err := http.NewRequest(method, url, bytes.NewBuffer(commandBody))
	if err != nil {
		return "", fmt.Errorf("failed to create command request: %v", err)
	}

	bodyHash := sha256.Sum256(commandBody)
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	stringToSign := method + "\n" + hex.EncodeToString(bodyHash[:]) + "\n" + "" + "\n" + path
	sign := calculateSignatureWithToken(clientID, timestamp, accessToken, stringToSign, clientSecret)

	req.Header.Set("client_id", clientID)
	req.Header.Set("access_token", accessToken)
	req.Header.Set("sign", sign)
	req.Header.Set("t", timestamp)
	req.Header.Set("sign_method", "HMAC-SHA256")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send command request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read command response body: %v", err)
	}

	return string(body), nil
}

func GetDeviceSpecification(deviceID, clientID, clientSecret string) (DeviceSpecificationResponse, error) {

	accessToken := getToken(clientID, clientSecret)

	var specResponse DeviceSpecificationResponse
	method := "GET"
	path := fmt.Sprintf("/v1.0/devices/%s/specifications", deviceID)
	url := TuyaBaseURL + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return specResponse, fmt.Errorf("failed to create specification request: %v", err)
	}

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	stringToSign := method + "\n" + hex.EncodeToString(sha256.New().Sum(nil)) + "\n" + "" + "\n" + path
	sign := calculateSignatureWithToken(clientID, timestamp, accessToken, stringToSign, clientSecret)

	req.Header.Set("client_id", clientID)
	req.Header.Set("access_token", accessToken)
	req.Header.Set("sign", sign)
	req.Header.Set("t", timestamp)
	req.Header.Set("sign_method", "HMAC-SHA256")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return specResponse, fmt.Errorf("failed to send specification request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return specResponse, fmt.Errorf("failed to read specification response body: %v", err)
	}

	if err := json.Unmarshal(body, &specResponse); err != nil {
		return specResponse, fmt.Errorf("failed to unmarshal device specification response: %v. Body: %s", err, string(body))
	}

	if !specResponse.Success {
		return specResponse, fmt.Errorf("API indicated failure to get specification: %s", string(body))
	}

	return specResponse, nil
}
