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
// Check latest test results
available, err := gmailnator.available()

// Initialize Gmailnator
gmail, err := gmailnator.NewGmailnator()

// Generate a new random Gmail address
err = gmail.GenerateEmail()

// Get lists of received emails
mailList, err := gmail.GetMails()

// Get the body of each email
for _, mail := range mailList {
    messageId := mail.Mid
    mailDetails, err := gmail.GetMailBody(messageId)
}
```

## Legal Disclaimer

This was made for educational purposes only, nobody which directly involved in this project is responsible for any damages caused.
**You are responsible for your actions.**

## License

[MIT](https://choosealicense.com/licenses/mit/)