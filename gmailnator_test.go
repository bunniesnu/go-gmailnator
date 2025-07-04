package gmailnator

import (
	"net/smtp"
	"os"
	"testing"
)

func TestNewGmailnator(t *testing.T) {
	// Test NewGmailnator
	gmailnator, err := NewGmailnator()
	if err != nil {
		t.Fatalf("Failed to create Gmailnator instance: %v", err)
	}
	if gmailnator.BaseURL != defaultURL {
		t.Errorf("Expected BaseURL to be %s, got %s", defaultURL, gmailnator.BaseURL)
	}
	if gmailnator.Cookie == nil {
		t.Error("Expected Cookie to be initialized, got nil")
	}
	if gmailnator.Client == nil {
		t.Error("Expected Client to be initialized, got nil")
	}

	// Test GenerateEmail
	email, err := gmailnator.GenerateEmail()
	if err != nil {
		t.Errorf("Failed to generate email: %v", err)
	}
	if email == "" {
		t.Error("Expected generated email to be non-empty, got empty string")
	}
	
	// Send test email
	login := os.Getenv("SMTP_LOGIN")
    from := os.Getenv("SMTP_FROM")
    password := os.Getenv("SMTP_PASSWORD")
    to := []string{email}
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
	if login == "" || from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		t.Fatal("SMTP environment variables are not set")
	}
    message := []byte("Subject: Test Email\n\nThis is the email body.")
    auth := smtp.PlainAuth("", login, password, smtpHost)
    err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
    if err != nil {
        t.Errorf("Failed to send test email: %v", err)
    }
}