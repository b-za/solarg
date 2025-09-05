# Solarg

Control smartlife wifi geyser switch based on the battery percentage of the FoxESS-Cloud inverter.

## FoxESS-Cloud Get Inverter Battery Charge Percentage

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

## Tuya Smartlife Wifi Geyser Switch

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
