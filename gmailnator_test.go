package gmailnator

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"testing"
	"time"
)

func SendTestMail(subject, body, destination string) error {
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	to := []string{destination}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	headers := []string{
		"From: " + from,
		"To: " + strings.Join(to, ","),
		fmt.Sprintf("Subject: %s", subject),
		body,
	}
	message := []byte(strings.Join(headers, "\r\n"))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func TestNewGmailnator(t *testing.T) {
	// Check if environment variables are set
	if os.Getenv("SMTP_FROM") == "" || os.Getenv("SMTP_PASSWORD") == "" {
		t.Fatal("Failed test because SMTP_FROM or SMTP_PASSWORD is not set in environment variables")
	}

	// Test NewGmailnator
	gmailnator, err := NewGmailnator()
	if err != nil {
		t.Fatalf("Failed to create Gmailnator instance: %v", err)
	}
	if gmailnator.XSRFToken == "" {
		t.Error("Expected Cookie to be initialized, got nil")
	}
	if gmailnator.Client == nil {
		t.Error("Expected Client to be initialized, got nil")
	}
	t.Log("Gmailnator instance create success")

	// Test GenerateEmail
	err = gmailnator.GenerateEmail()
	if err != nil {
		t.Errorf("Failed to generate email: %v", err)
	}
	if gmailnator.XSRFToken == "" {
		t.Error("Expected XSRFToken to be set, got empty string")
	}
	if gmailnator.Client == nil {
		t.Error("Expected Client to be set, got nil")
	}
	if gmailnator.Email == nil {
		t.Error("Expected Email to be generated, got nil")
	}
	if gmailnator.Email.Email == "" {
		t.Error("Expected Email.Email to be set, got empty string")
	}
	if gmailnator.Email.Timestamp == 0 {
		t.Error("Expected Email.Timestamp to be set, got zero value")
	}
	if gmailnator.RapidAPI == "" {
		t.Error("Expected RapidAPI to be set, got empty string")
	}
	t.Logf("Generated email: %s at %d", gmailnator.Email.Email, gmailnator.Email.Timestamp)
	
	// Send a test email
	destination := gmailnator.Email.Email
	subject := "Greetings from Go\r\n"
	body := "This is a test email from a Go program.\r\n"
	err = SendTestMail(subject, body, destination)
	if err != nil {
		t.Errorf("Failed to send test email: %v", err)
	}
	trimmedExpectedSubject := strings.TrimSpace(subject)
	trimmedExpectedBody := strings.TrimSpace(body)
	t.Logf("Test email sent to %s with subject '%s'", destination, trimmedExpectedSubject)

	// Test GetMails
	for range 5 {
		email, err := gmailnator.GetMails()
		if err != nil {
			t.Errorf("Failed to get mails: %v", err)
		}
		if email == nil {
			t.Error("Expected email to be returned, got nil")
		}
		for _, mail := range email {
			messageId := mail.Mid
			mailDetails, err := gmailnator.GetMailBody(messageId)
			if err != nil {
				t.Errorf("Failed to get mail details for message ID %s: %v", messageId, err)
			}
			if mailDetails == "" {
				t.Errorf("Expected mail details for message ID %s, got empty", messageId)
			}
			trimmedMailSubject := strings.TrimSpace(mail.Subject)
			trimmedMailDetails := strings.TrimSpace(mailDetails)
			if trimmedMailSubject != trimmedExpectedSubject {
				t.Errorf("Expected mail subject '%s', got '%s'", trimmedExpectedSubject, trimmedMailSubject)
			}
			if trimmedMailDetails != trimmedExpectedBody {
				t.Errorf("Expected mail body '%s', got '%s'", trimmedExpectedBody, trimmedMailDetails)
			}
			t.Logf("Received email from %s with subject '%s' and body '%s'", mail.From, trimmedMailSubject, trimmedMailDetails)
			return
		}
		time.Sleep(5 * time.Second) // Wait for 5 seconds before the next iteration
	}
	t.Errorf("Failed to retrieve any emails after multiple attempts")
}