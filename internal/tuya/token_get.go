package tuya

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetTuyaToken(clientID, clientSecret string) (TokenResult, error) {
	var tokenResult TokenResult

	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	signMethod := "HMAC-SHA256"
	httpMethod := "GET"
	requestPath := "/v1.0/token?grant_type=1"

	contentSHA256 := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	headersStr := ""
	stringToSignForHMAC := fmt.Sprintf("%s\n%s\n%s\n%s", httpMethod, contentSHA256, headersStr, requestPath)
	finalStringToSign := clientID + timestamp + stringToSignForHMAC

	sign := calculateSignature(finalStringToSign, clientSecret, signMethod)

	fullURL := TuyaBaseURL + requestPath
	req, err := http.NewRequest(httpMethod, fullURL, nil)
	if err != nil {
		return tokenResult, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("client_id", clientID)
	req.Header.Set("sign", sign)
	req.Header.Set("t", timestamp)
	req.Header.Set("sign_method", signMethod)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return tokenResult, fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tokenResult, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return tokenResult, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return tokenResult, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if !tokenResponse.Success {
		return tokenResult, fmt.Errorf("API returned a failure response: %s", string(body))
	}

	SaveToken(tokenResponse.Result)

	return tokenResponse.Result, nil
}

func calculateSignature(stringToSign, secret, signMethod string) string {
	if signMethod != "HMAC-SHA256" {
		log.Printf("Warning: Unsupported sign method '%s'", signMethod)
		return ""
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	hashed := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(hashed))
}
