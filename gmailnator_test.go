package gmailnator

import (
	"testing"
)

func TestNewGmailnator(t *testing.T) {
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
}