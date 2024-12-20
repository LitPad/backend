package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"github.com/shopspring/decimal"
	"gopkg.in/mail.v2"
)

var cfg = config.GetConfig()

func sortEmail(emailType string, tokenString *string, url *string, paymentData map[string]interface{}) map[string]interface{} {
	templateFile := "templates/welcome.html"
	subject := "Account verified"
	data := make(map[string]interface{})
	data["template_file"] = templateFile
	data["subject"] = subject
	data["text"] = "Your Verification was completed."

	// Sort different templates and subject for respective email types
	switch emailType {
	case "activate":
		templateFile = "templates/email-activation.html"
		subject = "Activate your account"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = "Please click the button below to verify your email."
		data["url"] = fmt.Sprintf("%s%s%s", *url, cfg.EmailVerificationPath, *tokenString)

	case "reset":
		templateFile = "templates/password-reset.html"
		subject = "Reset your password"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = "Please click the button below to reset your password."
		data["url"] = fmt.Sprintf("%s%s%s", *url, cfg.PasswordResetPath, *tokenString)

	case "reset-success":
		templateFile = "templates/password-reset-success.html"
		subject = "Password reset successfully"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = "Your password was reset successfully."
	case "payment-succeeded":
		templateFile = "templates/payment-success.html"
		amount := paymentData["amount"].(decimal.Decimal)
		subject = "Payment successful"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = fmt.Sprintf("Your payment of %s was successful.", amount) 
	case "payment-failed":
		templateFile = "templates/payment-failed.html"
		subject = "Payment failed"
		amount := paymentData["amount"].(decimal.Decimal)
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = fmt.Sprintf("Your payment of %s was unsuccessful. Please contact support", amount) 
	case "payment-canceled":
		templateFile = "templates/payment-canceled.html"
		subject = "Payment canceled"
		amount := paymentData["amount"].(decimal.Decimal)
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = "Your password was reset successfully."
		data["text"] = fmt.Sprintf("Your payment of %s was canceled.", amount)
	}
	return data
}

type EmailContext struct {
	Name string
	Url  *string
	Text string
}

func SendEmail(user *models.User, emailType string, tokenString *string, url *string, paymentData map[string]interface{}) {
	if os.Getenv("ENVIRONMENT") == "TESTING" {
		return
	}
	cfg := config.GetConfig()
	emailData := sortEmail(emailType, tokenString, url, paymentData)
	templateFile := emailData["template_file"]
	subject := emailData["subject"]

	// Create a context with dynamic data
	data := EmailContext{
		Name: user.FirstName,
		Text: emailData["text"].(string),
	}
	if url, ok := emailData["url"]; ok {
		url := url.(string)
		data.Url = &url
	}

	// Read the HTML file content
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("Unable to identify current directory (needed to load templates)", os.Stderr)
		os.Exit(1)
	}
	basepath := filepath.Dir(file)
	tempfile := fmt.Sprintf("../%s", templateFile.(string))
	htmlContent, err := os.ReadFile(filepath.Join(basepath, tempfile))
	if err != nil {
		log.Fatal("Error reading HTML file:", err)
	}

	// Create a new template from the HTML file content
	tmpl, err := template.New("email_template").Parse(string(htmlContent))
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	// Execute the template with the context and set it as the body of the email
	var bodyContent bytes.Buffer
	if err := tmpl.Execute(&bodyContent, data); err != nil {
		log.Fatal("Error executing template:", err)
	}

	// Create a new message
	m := mail.NewMessage()
	m.SetHeader("From", cfg.MailFrom)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject.(string))
	m.SetBody("text/html", bodyContent.String())

	// Create a new SMTP client
	d := mail.NewDialer(cfg.MailSenderHost, cfg.MailSenderPort, cfg.MailSenderEmail, cfg.MailSenderPassword)
	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Error sending email:", err)
	}
}

type ContactPayload struct {
	Email         string `json:"email"`
	ListIds       []int  `json:"listIds"`
	UpdateEnabled bool   `json:"updateEnabled"`
	Attributes  map[string]string `json:"attributes"`
}

func AddEmailToBrevo(name string, email string) {
	// Prepare the payload for Brevo API
	payload := ContactPayload{
		Email:       email,
		ListIds:     []int{cfg.BrevoListID}, // Convert string ListID to int
		UpdateEnabled: true,
		Attributes: map[string]string{
			"FIRSTNAME": name,
		},
	}
    // Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	// Create a new request to Brevo API
	req, err := http.NewRequest("POST", cfg.BrevoContactsUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println(err)
	}

    // Set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", cfg.MailApiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
