package main

import (
	"fmt"
	"html/template"
	"log"
	"time"

	"bytes"
	"encoding/json"
	"net/http"
)

type EmailData struct {
	BatteryMax int
	BatteryMin int
	StartTime  string
	EndTime    string
	SystemName string
}

type GeyserStatusEmailData struct {
	BatteryPercentage float64
	BatteryMax        float64
	BatteryMin        float64
	SystemName        string
	Timestamp         string
}

// sendHtmlEmail prepares and sends an email with HTML content.
func sendHtmlEmailStart() {

	subject := "Application Started"

	data := EmailData{
		BatteryMax: batteryMax,
		BatteryMin: batteryMin,
		StartTime:  startTimeStr,
		EndTime:    endTimeStr,
		SystemName: systemName,
	}

	htmlBody, err := generateEmailBody("email-start.html", data)
	if err != nil {
		log.Printf("ERROR: Failed to generate 'Geyser OFF' email body: %v", err)
		return
	}

	sendHtmlEmail(subject, htmlBody)
}

// sendGeyserOnEmail prepares and sends the "Geyser ON" notification.
func sendGeyserOnEmail(currentPercentage float64) {
	log.Println("Preparing 'Geyser ON' email...")
	data := GeyserStatusEmailData{
		BatteryPercentage: currentPercentage,
		BatteryMax:        batteryMax,
		SystemName:        systemName,
		Timestamp:         time.Now().Format("2006-01-02 15:04 SAST"),
	}

	htmlBody, err := generateEmailBody("email-on.html", data)
	if err != nil {
		log.Printf("ERROR: Failed to generate 'Geyser ON' email body: %v", err)
		return
	}

	sendHtmlEmail("✅ Geyser Switched ON", htmlBody)
}

func sendGeyserOffEmail(currentPercentage float64) {
	log.Println("Preparing 'Geyser OFF' email...")
	data := GeyserStatusEmailData{
		BatteryPercentage: currentPercentage,
		BatteryMin:        batteryMin,
		SystemName:        systemName,
		Timestamp:         time.Now().Format("2006-01-02 15:04 SAST"),
	}

	htmlBody, err := generateEmailBody("email-off.html", data)
	if err != nil {
		log.Printf("ERROR: Failed to generate 'Geyser OFF' email body: %v", err)
		return
	}

	sendHtmlEmail("🔌 Geyser Switched OFF", htmlBody)
}

func sendHtmlEmail(subject, htmlBody string) {

	var recipients []map[string]string
	for _, email := range mailtrapToEmails {
		recipients = append(recipients, map[string]string{"email": email})
	}

	emailPayload := map[string]interface{}{
		"from":    map[string]string{"email": mailtrapFromEmail, "name": "SolarG"},
		"to":      recipients,
		"subject": subject,
		"html":    htmlBody,
	}
	sendMailtrapRequest(emailPayload)

}

// sendPlainTextEmail prepares and sends an email with plain text content.
func sendPlainTextEmail(subject, textBody string) {

	var recipients []map[string]string
	for _, email := range mailtrapToEmails {
		recipients = append(recipients, map[string]string{"email": email})
	}
	emailPayload := map[string]interface{}{
		"from":    map[string]string{"email": mailtrapFromEmail, "name": "SolarG"},
		"to":      recipients,
		"subject": subject,
		"text":    textBody,
	}
	sendMailtrapRequest(emailPayload)
}

// sendMailtrapRequest handles the marshalling and sending of the email payload.
func sendMailtrapRequest(payload map[string]interface{}) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error: Could not marshal email payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", mailtrapAPIURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error: Could not create HTTP request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mailtrapAPIToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: Could not send email: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Email sent successfully!")
	} else {
		log.Printf("Error: Failed to send email. Status: %s", resp.Status)
	}
}

func generateEmailBody(templateFile string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", fmt.Errorf("could not parse template file %s: %w", templateFile, err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", fmt.Errorf("could not execute template %s: %w", templateFile, err)
	}

	return tpl.String(), nil
}
