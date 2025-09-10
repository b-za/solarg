package main

import (
	"log"
	"path/filepath"
	"runtime"
	"time"
)

// getBasePath returns the absolute path to the project's root directory.
func getBasePath() string {
	// Get the path of the file that calls this function.
	_, b, _, _ := runtime.Caller(0)

	// The project root is the directory of the caller's file.
	// You might need to adjust this if your go files are in a subdirectory e.g. filepath.Dir(b) + "/.."
	return filepath.Dir(b)
}

// checkTime gets the current time and determines if it is within the active window.
func checkTime(location *time.Location) {
	// Get the current time in the specified location.
	now := time.Now().In(location)

	// Parse the start and end time strings.
	// The date part is ignored, only the time of day is used.
	layout := "15:04" // "HH:MM" format
	startTime, err := time.Parse(layout, win1Start)
	if err != nil {
		log.Printf("Error: Could not parse start time: %v", err)
		return
	}
	endTime, err := time.Parse(layout, win1End)
	if err != nil {
		log.Printf("Error: Could not parse end time: %v", err)
		return
	}

	// Construct the full start and end time for the *current* day.
	// This ensures the comparison is always against today's window.
	year, month, day := now.Date()
	win1ActiveStartTime := time.Date(year, month, day, startTime.Hour(), startTime.Minute(), 0, 0, location)
	win1ActiveEndTime := time.Date(year, month, day, endTime.Hour(), endTime.Minute(), 0, 0, location)

	// Check if the current time is after the start and before the end.
	if now.After(win1ActiveStartTime) && now.Before(win1ActiveEndTime) {
		log.Printf("[%s] The current time is WITHIN the active window.", now.Format("15:04:05"))
		activeWindowLoop()
	} else {
		log.Printf("[%s] The current time is NOT within the active window.", now.Format("15:04:05"))
		inactiveWindowLoop()
	}
}
