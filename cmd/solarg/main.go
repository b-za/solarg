package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/b-za/solarg/internal/fox"
)

// --- Configuration ---
// Set the start and end times for the active window in "HH:MM" format.
const (
	startTimeStr = "09:00"
	endTimeStr   = "17:00"
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
	ticker := time.NewTicker(5 * time.Minute)
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
	} else {
		log.Printf("[%s] The current time is NOT within the active window.", now.Format("15:04:05"))
	}
}

func getBatteryStatus() {

	foxApiKey := FOX_API_KEY
	foxInverterSerialNumber := FOX_INVERTER_SERIAL_NUMBER

	// This is the formula for estimating the battery percentage
	// Change the divide by as needed to improve the estimate
	//"[(residualEnergy / divideBy) * 100]"
	divideBy := 6.6

	if foxApiKey == "" || foxInverterSerialNumber == "" {
		fmt.Println("❌ Error: No foxApiKey or foxInverterSerialNumber")
		flag.Usage()
		os.Exit(1)
	}

	residualEnergy, batteryStatus := fox.GetBatteryStatus(foxApiKey, foxInverterSerialNumber, divideBy)

	fmt.Println("residualEnergy")
	fmt.Println(residualEnergy)
	fmt.Println("divideBy")
	fmt.Println(divideBy)
	fmt.Println("batteryStatus")
	fmt.Println("[(residualEnergy / divideBy) * 100]")
	fmt.Println(batteryStatus)
}
