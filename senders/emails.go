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
	"strings"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/models"
	"github.com/shopspring/decimal"
	"gopkg.in/mail.v2"
)

type EmailTypeChoice string

const (
	ET_ACTIVATE  EmailTypeChoice = "activate"
	ET_WELCOME EmailTypeChoice = "welcome"
	ET_RESET     EmailTypeChoice = "reset"
	ET_RESET_SUCC EmailTypeChoice = "reset-success"
	ET_PAYMENT_SUCC EmailTypeChoice = "payment-succeeded"
	ET_PAYMENT_FAIL EmailTypeChoice = "payment-failed"
	ET_PAYMENT_CANCEL EmailTypeChoice = "payment-canceled"
	ET_SUBSCRIPTION_EXPIRING EmailTypeChoice = "subscription-expiring"
	ET_SUBSCRIPTION_EXPIRED EmailTypeChoice = "subscription-expired"
)

func sortEmail(cfg config.Config, emailType EmailTypeChoice, otp *uint, tokenString *string, extraData map[string]interface{}) map[string]interface{} {
	templateFile := "templates/welcome.html"
	subject := "Account verified"
	data := make(map[string]interface{})
	data["template_file"] = templateFile
	data["subject"] = subject
	data["text"] = "Your Verification was completed."

	// Sort different templates and subject for respective email types
	switch emailType {
	case ET_ACTIVATE:
		templateFile = "templates/email-activation.html"
		subject = "Activate your account"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = "Please use the code below to verify your email."
		data["code"] = otp

	case ET_RESET:
		templateFile = "templates/password-reset.html"
		subject = "Reset your password"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = "Please click the button below to reset your password."
		data["url"] = fmt.Sprintf("%s://reset-password?token=%s", cfg.AppScheme, *tokenString)

	case ET_RESET_SUCC:
		templateFile = "templates/password-reset-success.html"
		subject = "Password reset successfully"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = "Your password was reset successfully."
	case ET_PAYMENT_SUCC:
		templateFile = "templates/payment-success.html"
		amount := extraData["amount"].(decimal.Decimal)
		subject = "Payment successful"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = fmt.Sprintf("Your payment of %s was successful.", amount) 
	case ET_PAYMENT_FAIL:
		templateFile = "templates/payment-failed.html"
		subject = "Payment failed"
		amount := extraData["amount"].(decimal.Decimal)
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = fmt.Sprintf("Your payment of %s was unsuccessful. Please contact support", amount) 
	case ET_PAYMENT_CANCEL:
		templateFile = "templates/payment-canceled.html"
		subject = "Payment canceled"
		amount := extraData["amount"].(decimal.Decimal)
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = fmt.Sprintf("Your payment of %s was canceled.", amount)
	case ET_SUBSCRIPTION_EXPIRING:
		templateFile = "templates/subscription-expiring.html"
		subject = "Subscription close to expiry"
		subscriptionType := strings.ToLower(extraData["subscriptionType"].(string)) 
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = fmt.Sprintf("Your %s book subscription is about to expire.", subscriptionType)
	case ET_SUBSCRIPTION_EXPIRED:
		templateFile = "templates/subscription-expired.html"
		subject = "Subscription expired"
		subscriptionType := strings.ToLower(extraData["subscriptionType"].(string)) 
		data["template_file"] = templateFile
		data["subject"] = subject
		data["text"] = fmt.Sprintf("Your %s book subscription has expired. Please renew your subscription", subscriptionType)
	}
	return data
}

type EmailContext struct {
	Name string
	Url  *string
	Code *uint
	Text string
}

func SendEmail(user *models.User, emailType EmailTypeChoice, otp *uint, tokenString *string, paymentData map[string]interface{}) {
	if os.Getenv("ENVIRONMENT") == "TESTING" {
		return
	}
	cfg := config.GetConfig()
	emailData := sortEmail(cfg, emailType, otp, tokenString, paymentData)
	templateFile := emailData["template_file"]
	subject := emailData["subject"]

	// Create a context with dynamic data
	data := EmailContext{
		Name: user.Username,
		Text: emailData["text"].(string),
	}
	if url, ok := emailData["url"]; ok {
		url := url.(string)
		data.Url = &url
	}

	if code, ok := emailData["code"]; ok {
		code := code.(*uint)
		data.Code = code
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
		log.Println("Error sending email:", err)
	}
}

type ContactPayload struct {
	Email         string `json:"email"`
	ListIds       []int  `json:"listIds"`
	UpdateEnabled bool   `json:"updateEnabled"`
	Attributes  map[string]string `json:"attributes"`
}

func AddEmailToBrevo(name string, email string) {
	cfg := config.GetConfig()

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
