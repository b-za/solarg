package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/b-za/solarg/internal/fox"
	"github.com/b-za/solarg/internal/tuya"
)

const foxApiKey = FOX_API_KEY
const foxInverterSerialNumber = FOX_INVERTER_SERIAL_NUMBER

const TuyaClientId = TUYA_CLIENT_ID
const TuyaClientSecret = TUYA_CLIENT_SECRET
const TuyaDeviceID = TUYA_DEVICE_ID

const mailtrapFromEmail = MAILTRAP_FROM_EMAIL
const mailtrapAPIURL = MAILTRAP_API_URL
const mailtrapAPIToken = MAILTRAP_API_TOKEN

var mailtrapToEmails = MAILTRAP_TO_EMAILS

// --- Configuration ---
// Set the start and end times for the active window in "HH:MM" format.
const (
	startTimeStr = "09:00"
	endTimeStr   = "15:30"
	batteryMin   = 60.00
	batteryMax   = 80.00
	systemName   = "SolarG"
	// Set the location to ensure the time comparison is accurate.
	// SAST (South Africa Standard Time) corresponds to "Africa/Johannesburg".
	locationName = "Africa/Johannesburg"
)

var batMinStr = strconv.FormatFloat(batteryMin, 'f', 2, 64)
var batMaxStr = strconv.FormatFloat(batteryMax, 'f', 2, 64)
var batMinStr2 = strconv.FormatFloat(batteryMin, 'f', 2, 64)
var batMaxStr2 = strconv.FormatFloat(batteryMax, 'f', 2, 64)

func main() {

	// Load the specified time zone to ensure comparisons are correct.
	location, err := time.LoadLocation(locationName)
	if err != nil {
		log.Printf("Could not load location %s: %v", locationName, err)
	}

	log.Printf("Application started. Checking time every 5 minutes.")
	log.Printf("Active window is between %s and %s (%s).", startTimeStr, endTimeStr, locationName)

	//sendHtmlEmailStart()

	// Create a ticker that fires every 5 minutes.
	ticker := time.NewTicker(5 * time.Minute)

	// Ticker every 5 seconds
	// ticker := time.NewTicker(30 * time.Second)

	// Ensure the ticker is stopped when the function exits to clean up resources.
	defer ticker.Stop()

	// Run the check immediately on startup, then wait for the ticker.
	checkTime(location)

	// Start an infinite loop to listen for ticks.
	for range ticker.C {
		checkTime(location)
	}
}

func inactiveWindowLoop() {

	geyserStatusResponse, err := tuya.GetSwitchStatus(TuyaDeviceID, TuyaClientId, TuyaClientSecret)
	if err != nil {
		log.Printf("Failed to get switch state: %v", err)
	}
	if geyserStatusResponse.Success {
		if geyserStatusResponse.Status == true {

			response, err := tuya.SetSwitchState(TuyaDeviceID, TuyaClientId, TuyaClientSecret, false)
			if err != nil {
				log.Printf("Failed to set switch state: %v", err)
			}
			fmt.Println("API Response:")
			fmt.Println(response)

		}
	} else {

		response, err := tuya.SetSwitchState(TuyaDeviceID, TuyaClientId, TuyaClientSecret, false)
		if err != nil {
			log.Printf("Failed to set switch state: %v", err)
		}
		fmt.Println("API Response:")
		fmt.Println(response)

	}
}

func activeWindowLoop() {

	var batteryPercentage float64
	// var residualEnergy float64
	var geyserOnStatus bool

	// Get the status of the geyser switch
	geyserStatusResponse, err := tuya.GetSwitchStatus(TuyaDeviceID, TuyaClientId, TuyaClientSecret)
	if err != nil {
		log.Printf("Failed to get switch state: %v", err)
	}

	if geyserStatusResponse.Success {
		geyserOnStatus = geyserStatusResponse.Status
	} else {
		geyserOnStatus = true
	}

	// This is the formula for estimating the battery percentage
	// Change the divide by as needed to improve the estimate
	//"[(residualEnergy / divideBy) * 100]"
	var divideBy = 6.6

	_, batteryPercentage = fox.GetBatteryStatus(foxApiKey, foxInverterSerialNumber, divideBy)

	fmt.Println("Battery % is:", batteryPercentage)

	if batteryPercentage > batteryMax {
		if geyserOnStatus == false {
			fmt.Printf("Battery is > %s, and the geyser is off. Turn on the geyser", batMaxStr)
			sendGeyserOnEmail(batteryPercentage)
			response, err := tuya.SetSwitchState(TuyaDeviceID, TuyaClientId, TuyaClientSecret, true)
			if err != nil {
				log.Printf("Failed to set switch state: %v", err)
			}
			fmt.Println("API Response:")
			fmt.Println(response)

		} else {
			fmt.Printf("Battery is > %s, and the geyser is on. Leave the geyser alone.", batMaxStr)
		}
	} else if batteryPercentage < batteryMin {
		if geyserOnStatus == true {
			fmt.Printf("Battery is < %s, and the geyser is on. Turn off the geyser", batMinStr)
			sendGeyserOffEmail(batteryPercentage)
			response, err := tuya.SetSwitchState(TuyaDeviceID, TuyaClientId, TuyaClientSecret, false)
			if err != nil {
				log.Printf("Failed to set switch state: %v", err)
			}
			fmt.Println("API Response:")
			fmt.Println(response)
		} else {
			fmt.Printf("Battery is < %s, and the geyser is off. Leave the geyser alone.", batMinStr)
		}

	} else {
		fmt.Printf("Battery is between %s and %s. Leave the geyser alone.", batMinStr, batMaxStr)
	}

}
