package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/b-za/solarg/internal/tuya"
)

func main() {
	// --- Define command-line flags ---
	getSpec := flag.Bool("spec", false, "Get the device's command and status specifications.")
	getStatus := flag.Bool("status-all", false, "Get the device's current status.")
	getSwitch := flag.Bool("status", false, "Get the switch state.a")
	setSwitch := flag.String("switch", "", "Set the switch state. Accepts 'on' or 'off'.")
	getToken := flag.Bool("get-token", false, "Force fetch a new API token and save it.")

	flag.Parse()

	// If the -get-token flag is used, just get the token and exit.
	if *getToken {
		fmt.Println("Forcing retrieval of a new API token...")
		_, err := tuya.GetTuyaToken(TuyaClientId, TuyaClientSecret)
		if err != nil {
			log.Fatalf("Failed to retrieve a new token: %v", err)
		}
		fmt.Println("New token successfully retrieved and saved.")
		return
	}

	// Perform action based on flags ---
	switch {
	case *getSpec:
		fmt.Printf("\n--- Getting command specifications for device %s ---\n", TuyaDeviceID)
		spec, err := tuya.GetDeviceSpecification(TuyaDeviceID, TuyaClientId, TuyaClientSecret)
		if err != nil {
			log.Fatalf("Failed to get device specification: %v", err)
		}
		specJSON, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			log.Fatalf("Failed to format specification JSON for printing: %v", err)
		}
		fmt.Println("Device Specifications:")
		fmt.Println(string(specJSON))

	case *getSwitch:
		fmt.Printf("\n--- Getting status for device %s ---\n", TuyaDeviceID)
		status, err := tuya.GetSwitchStatus(TuyaDeviceID, TuyaClientId, TuyaClientSecret)
		if err != nil {
			log.Fatalf("Failed to get device status: %v", err)
		}
		statusJSON, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			log.Fatalf("Failed to format status JSON: %v", err)
		}
		fmt.Println("Device Status:")
		fmt.Println(string(statusJSON))

	case *getStatus:
		fmt.Printf("\n--- Getting status for device %s ---\n", TuyaDeviceID)
		status, err := tuya.GetDeviceStatus(TuyaDeviceID, TuyaClientId, TuyaClientSecret)
		if err != nil {
			log.Fatalf("Failed to get device status: %v", err)
		}
		statusJSON, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			log.Fatalf("Failed to format status JSON: %v", err)
		}
		fmt.Println("Device Status:")
		fmt.Println(string(statusJSON))

	case *setSwitch != "":
		var turnOn bool
		switch strings.ToLower(*setSwitch) {
		case "on":
			turnOn = true
		case "off":
			turnOn = false
		default:
			log.Fatalf("Invalid value for -switch flag. Please use 'on' or 'off'.")
		}
		fmt.Printf("\n--- Setting switch state to '%s' for device %s ---\n", *setSwitch, TuyaDeviceID)
		response, err := tuya.SetSwitchState(TuyaDeviceID, TuyaClientId, TuyaClientSecret, turnOn)
		if err != nil {
			log.Fatalf("Failed to set switch state: %v", err)
		}
		fmt.Println("API Response:")
		fmt.Println(response)

	default:
		// Default action: get status
		fmt.Printf("\n--- Getting status for device %s ---\n", TuyaDeviceID)
		status, err := tuya.GetDeviceStatus(TuyaDeviceID, TuyaClientId, TuyaClientSecret)
		if err != nil {
			log.Fatalf("Failed to get device status: %v", err)
		}
		fmt.Println("Device Status:")
		fmt.Println(status.Result[0].Code)
		fmt.Println(status.Result[0].Value)
	}
}
