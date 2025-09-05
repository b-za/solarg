package fox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"

	"time"
)

func GetBatteryStatus(apiKey, inverterSerialNumber string, divideBy float64) (float64, float64) {

	var residualEnergy float64
	var batteryStatus float64

	if apiKey == "" || inverterSerialNumber == "" {
		fmt.Println("❌ Fox Error: no apiKey or inverterSerialNumber")
		os.Exit(1)
	}

	const (
		foxESSCloudDomain  = "https://www.foxesscloud.com"
		foxRealDataURLPath = "/op/v0/device/real/query"
	)

	requestBody := RealDataRequestBody{
		SN:        inverterSerialNumber,
		Variables: []string{"ResidualEnergy", "batVolt"},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("❌ Fox Error marshalling request body: %v\n", err)
		os.Exit(1)
	}

	fullRequestURL := fmt.Sprintf("%s%s", foxESSCloudDomain, foxRealDataURLPath)
	requestPath := foxRealDataURLPath

	req, err := http.NewRequest("POST", fullRequestURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("❌ Fox Error creating real data request: %v\n", err)
		os.Exit(1)
	}

	headers, err := buildSignedHeaders(requestPath, apiKey)
	if err != nil {
		fmt.Printf("❌ Fox Error building signed headers: %v\n", err)
		os.Exit(1)
	}
	req.Header = headers

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Fox Error making real data request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Fox Error : Failed to get battery status. Status: %s, Response: %s\n", resp.Status, string(bodyBytes))
		os.Exit(1)
	}

	var realDataResp RealDataResponse

	if err := json.NewDecoder(resp.Body).Decode(&realDataResp); err != nil {
		fmt.Printf("❌ Fox Error decoding real data response: %v\n", err)
		os.Exit(1)
	}

	if realDataResp.Errno != 0 {
		fmt.Printf("❌ Fox Error :FoxESS Cloud API returned an error. (Code: %d)\n", realDataResp.Errno)
		fmt.Println(realDataResp)
		os.Exit(1)
	}

	for _, v := range realDataResp.Result[0].Datas {
		//fmt.Println("----------")
		//fmt.Println(v.Name)
		//fmt.Printf("%s: %.2f %s\n", v.Variable, v.Value, v.Unit)
		if v.Variable == "ResidualEnergy" {

			residualEnergy = v.Value

			batteryStatus = math.Round((residualEnergy / divideBy) * 100)

		}
	}
	return residualEnergy, batteryStatus
}
