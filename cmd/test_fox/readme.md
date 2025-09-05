# FoxESS-Cloud Get Inverter Battery Charge Percentage

This is a simple program to get the battery charge percentage from a FoxESS-Cloud inverter using Golang

In order to get access to the data you need and private token to use as API key from the API management section after logging into FoxESS-Cloud
https://www.foxesscloud.com/user/center

Create a files named secrets.go with the following content

```go

package main

var INVERTER_SERIAL_NUMBER = "YOUR_INVERTER_SERIAL_NUMBER"
var API_KEY = "YOUR_PRIVATE_TOKEN"

```

To run the program, run the following command

```bash
go run .
```

# Sources

https://www.foxesscloud.com/

https://github.com/TonyM1958/FoxESS-Cloud

https://github.com/SoftXperience/home-assistant-foxess-api/blob/main/custom_components/foxess_api/fox_ess_cloud_api.py

https://www.foxesscloud.com/public/i18n/en/OpenApiDocument.html

https://pvoutput.org/help/api_specification.html#csv-data-parameter
