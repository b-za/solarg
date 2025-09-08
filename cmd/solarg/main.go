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

// --- Configuration ---
// Set the start and end times for the active window in "HH:MM" format.
const (
	startTimeStr = "09:00"
	endTimeStr   = "15:30"
	batteryMin   = 60.00
	batteryMax   = 80.00
	// Set the location to ensure the time comparison is accurate.
	// SAST (South Africa Standard Time) corresponds to "Africa/Johannesburg".
	locationName = "Africa/Johannesburg"
)

func main() {
	// Load the specified time zone to ensure comparisons are correct.
	location, err := time.LoadLocation(locationName)
	if err != nil {
		log.Fatalf("Fatal: Could not load location %s: %v", locationName, err)
	}

	log.Printf("Application started. Checking time every 5 minutes.")
	log.Printf("Active window is between %s and %s (%s).", startTimeStr, endTimeStr, locationName)

	// Create a ticker that fires every 5 minutes.
	//ticker := time.NewTicker(5 * time.Minute)

	// Ticker every 5 seconds
	ticker := time.NewTicker(30 * time.Second)

	// Ensure the ticker is stopped when the function exits to clean up resources.
	defer ticker.Stop()

	// Run the check immediately on startup, then wait for the ticker.
	checkTime(location)

	// Start an infinite loop to listen for ticks.
	for range ticker.C {
		checkTime(location)
	}
}

// checkTime gets the current time and determines if it is within the active window.
func checkTime(location *time.Location) {
	// Get the current time in the specified location.
	now := time.Now().In(location)

	// Parse the start and end time strings.
	// The date part is ignored, only the time of day is used.
	layout := "15:04" // "HH:MM" format
	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		log.Printf("Error: Could not parse start time: %v", err)
		return
	}
	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		log.Printf("Error: Could not parse end time: %v", err)
		return
	}

	// Construct the full start and end time for the *current* day.
	// This ensures the comparison is always against today's window.
	year, month, day := now.Date()
	activeStartTime := time.Date(year, month, day, startTime.Hour(), startTime.Minute(), 0, 0, location)
	activeEndTime := time.Date(year, month, day, endTime.Hour(), endTime.Minute(), 0, 0, location)

	// Check if the current time is after the start and before the end.
	if now.After(activeStartTime) && now.Before(activeEndTime) {
		log.Printf("[%s] The current time is WITHIN the active window.", now.Format("15:04:05"))
		activeWindowLoop()
	} else {
		log.Printf("[%s] The current time is NOT within the active window.", now.Format("15:04:05"))
	}
}

func activeWindowLoop() {

	batMinStr := strconv.FormatFloat(batteryMin, 'f', 2, 64)
	batMaxStr := strconv.FormatFloat(batteryMax, 'f', 2, 64)

	var batteryPercentage float64
	// var residualEnergy float64
	var geyserOnStatus bool

	// This is the formula for estimating the battery percentage
	// Change the divide by as needed to improve the estimate
	//"[(residualEnergy / divideBy) * 100]"
	var divideBy = 6.6

	_, batteryPercentage = fox.GetBatteryStatus(foxApiKey, foxInverterSerialNumber, divideBy)

	fmt.Println("Battery % is:", batteryPercentage)

	if batteryPercentage > batteryMax {
		if geyserOnStatus == false {
			fmt.Printf("Battery % is > %s, and the geyser is off. Turn on the geyser", batMaxStr)
			response, err := tuya.SetSwitchState(TuyaDeviceID, TuyaClientId, TuyaClientSecret, false)
			if err != nil {
				log.Fatalf("Failed to set switch state: %v", err)
			}
			fmt.Println("API Response:")
			fmt.Println(response)

		} else {
			fmt.Printf("Battery % is > %s, and the geyser is on. Leave the geyser alone.", batMaxStr)
		}
	} else if batteryPercentage < batteryMin {
		if geyserOnStatus == true {
			fmt.Printf("Battery % is < %s, and the geyser is on. Turn off the geyser", batMinStr)
			response, err := tuya.SetSwitchState(TuyaDeviceID, TuyaClientId, TuyaClientSecret, false)
			if err != nil {
				log.Fatalf("Failed to set switch state: %v", err)
			}
			fmt.Println("API Response:")
			fmt.Println(response)
		} else {
			fmt.Printf("Battery % is < %s, and the geyser is off. Leave the geyser alone.", batMinStr)
		}

	} else {
		fmt.Printf("Battery % is between %s and %s. Leave the geyser alone.", batMinStr, batMaxStr)
	}

	// float64 to string
	// fmt.Println("Battery % is:", strconv.FormatFloat(batteryPercentage, 'f', 2, 64))

}
