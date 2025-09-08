package main

import (
	"log"
	"time"
)

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
		inactiveWindowLoop()
	}
}
