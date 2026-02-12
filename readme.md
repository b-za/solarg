# Solarg

Control smartlife wifi geyser switch via Tuya based on the battery percentage of the FoxESS-Cloud inverter.

## Overview

SolarG is a Go-based automation script designed to intelligently control a smart switch (e.g., for a geyser/water heater) by leveraging real-time data from a FoxESS solar battery system. The primary goal is to maximize the use of surplus solar energy for heating water while preserving battery charge for essential loads.

The application runs continuously, checking your solar system's status every five minutes. It operates within a configurable time window (the "active window"), which should be set to coincide with peak solar production hours.

Based on the battery's state of charge (SoC), SolarG makes automated decisions to turn a Tuya-compatible smart switch on or off, ensuring that you heat your water with free energy from the sun without unnecessarily draining your battery.

### How It Works

The core logic follows a simple set of rules:

1.  **Time-Based Operation**: The script checks the system every 5 minutes.
2.  **Active Window**: It first determines if the current time is within the defined active window (e.g., 09:00 - 15:30).
3.  **Decision Logic**:
    - **Inside the Active Window:**
      - If the battery charge is **above the `batteryMax` threshold** (e.g., 80%) and the geyser is off, it turns the **geyser ON**.
      - If the battery charge drops **below the `batteryMin` threshold** (e.g., 60%) and the geyser is on, it turns the **geyser OFF**.
      - If the battery charge is between the min and max thresholds, it takes no action.
    - **Outside the Active Window:**
      - The script ensures the **geyser is turned OFF** to conserve power overnight.
4.  **Notifications**: The application uses Mailtrap to send email alerts whenever it turns the switch on or off, keeping you informed of its activity.

#### Core Logic Summary

| Condition                                         | Geyser Status | Action Taken |
| ------------------------------------------------- | ------------- | ------------ |
| **Inside** active window & Battery > `batteryMax` | OFF           | Turn **ON**  |
| **Inside** active window & Battery < `batteryMin` | ON            | Turn **OFF** |
| **Outside** active window                         | ON            | Turn **OFF** |
| _All other states_                                | -             | No Action    |

### Integrations

This application connects to three external services:

- **FoxESS Cloud API**: To retrieve real-time battery status and residual energy data.
- **Tuya IoT Platform**: To control the state (on/off) of the smart switch.
- **Mailtrap**: To send transactional email notifications.

### Configuration

All configuration is managed via constants in the `main.go` file. You will need to provide your own API keys, secrets, and device IDs for the services above.

Key configurable parameters include:

- `startTimeStr` / `endTimeStr`: The start and end of the active window.
- `batteryMin` / `batteryMax`: The float values for the battery charge thresholds.
- `locationName`: Your timezone (e.g., "Africa/Johannesburg") to ensure correct time comparisons.
- API keys and device IDs for FoxESS, Tuya, and Mailtrap.

## Some notes on running the app on a linux server

https://gist.github.com/b-za/75035ae7168eb40fc721038c6e9e76d9

export the go path

```bash
export PATH=$PATH:/usr/local/go/bin
```

Setup the systemd service

```bash
sudo vi /lib/systemd/system/solarg.service
```

The config for the service

```yml
[Unit]
Description=solarg

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/home/rootpi/apps/solarg/cmd/solarg/solarg

[Install]
WantedBy=multi-user.target
```

Start the service

```bash
sudo service solarg start
```

Check the status of the service

```bash
sudo service solarg status
sudo systemctl status solarg


```

To make the service start automatically every time you boot:

````bash

sudo systemctl enable solarg

```

Restart the service

```bash
sudo systemctl daemon-reexec

````

## Below is some explanations for the two tester apps

There is not a proper written explanation for the main app

### FoxESS-Cloud Get Inverter Battery Charge Percentage

./cmd/test_fox/

This is a simple program to get the battery charge percentage from a FoxESS-Cloud inverter using Golang

In order to get access to the data you need and private token to use as API key from the API management section after logging into FoxESS-Cloud
https://www.foxesscloud.com/user/center

Create a file named secrets.go with the following content

```go

package main

var INVERTER_SERIAL_NUMBER = "YOUR_INVERTER_SERIAL_NUMBER"
var API_KEY = "YOUR_PRIVATE_TOKEN"

```

To run the program, run the following command

```bash
go run .
```

### Tuya Smartlife Wifi Geyser Switch Control via Tuya

./cmd/test_tuya/

Create a file named secrets.go with the following content

```go
package main

const TuyaClientId = "YOUR_CLIENT_ID"
const TuyaClientSecret = "YOUR_CLIENT_SECRET"
const TuyaDeviceID = "YOUR_DEVICE_ID"

```

Get the access token

```bash
go run . --get-token
```

See the current status of the wifi switch

```bash
go run . --status
go run .
```

Toggle the wifi switch

```bash
go run . --switch=on
go run . --switch=off
```

Get the spec for the wifi switch

```bash
go run . --spec
```

Get the full status of the wifi switch not just the current on/off status

```bash
go run . --status-all
```

## Sources

https://www.foxesscloud.com/

https://github.com/TonyM1958/FoxESS-Cloud

https://github.com/SoftXperience/home-assistant-foxess-api/blob/main/custom_components/foxess_api/fox_ess_cloud_api.py

https://www.foxesscloud.com/public/i18n/en/OpenApiDocument.html

https://pvoutput.org/help/api_specification.html#csv-data-parameter
