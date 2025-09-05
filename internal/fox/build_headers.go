package fox

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"strconv"
	"time"
)

// The `requestPath` should be the URI path (e.g., "/op/v0/device/real/get"), not the full URL.
func buildSignedHeaders(requestPath, apiToken string) (http.Header, error) {

	headers := make(http.Header)

	headers.Set("Content-Type", "application/json")

	currentTimestamp := time.Now().UnixMilli()
	timestamp := strconv.FormatInt(currentTimestamp, 10)

	headers.Set("token", apiToken)
	headers.Set("lang", "en")
	//headers.Set("user-agent", UserAgent)

	testPath := "/op/v0/device/real/query"

	signature := generateSignature(testPath, apiToken, timestamp)

	headers.Set("Timestamp", timestamp)
	headers.Set("Signature", signature)

	return headers, nil
}

func generateSignature(path, token, timestamp string) string {
	signatureString := fmt.Sprintf("%s\\r\\n%s\\r\\n%s", path, token, timestamp)
	hasher := md5.New()
	hasher.Write([]byte(signatureString))
	return hex.EncodeToString(hasher.Sum(nil))
}
