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
	email, err := gmailnator.GenerateEmail()
	if err != nil {
		t.Errorf("Failed to generate email: %v", err)
	}
	if email == "" {
		t.Error("Expected generated email to be non-empty, got empty string")
	}
}