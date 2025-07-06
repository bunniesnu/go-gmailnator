# Go-Gmailnator

![Go Package](https://img.shields.io/github/release/bunniesnu/go-gmailnator.svg)
![Go Tests](https://github.com/bunniesnu/go-gmailnator/actions/workflows/test-schedule.yml/badge.svg)

A [Go](https://go.dev) package for generating random gmail address and receiving emails.

## Installation

```
go get github.com/bunniesnu/go-gmailnator
```

## Usage

```
// Initialize Gmailnator
gmailnator, err := NewGmailnator()

// Generate a new random Gmail address
err = gmailnator.GenerateEmail()

// Get lists of received emails
mailList, err := gmailnator.GetMails()

// Get the body of each email
for _, mail := range email {
    messageId := mail.Mid
    mailDetails, err := gmailnator, GetMailBody(messageId)
}
```