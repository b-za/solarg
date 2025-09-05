package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/b-za/solarg/internal/fox"
)

func main() {

	apiKey := API_KEY
	inverterSerialNumber := INVERTER_SERIAL_NUMBER
	divideBy := 6.6

	if apiKey == "" || inverterSerialNumber == "" {
		fmt.Println("❌ Error: currentAccessToken or currentStationID")
		flag.Usage()
		os.Exit(1)
	}

	residualEnergy, batteryStatus := fox.GetBatteryStatus(apiKey, inverterSerialNumber, divideBy)

	fmt.Println("residualEnergy")
	fmt.Println(residualEnergy)
	fmt.Println("divideBy")
	fmt.Println(divideBy)
	fmt.Println("batteryStatus")
	fmt.Println("[(residualEnergy / divideBy) * 100]")
	fmt.Println(batteryStatus)

}
