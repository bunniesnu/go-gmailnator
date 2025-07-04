package gmailnator

import "fmt"

type Address struct {
	Email string
}

func NewAddress() *Address {
	// TODO: Implement random address generation
	return &Address{
		Email: fmt.Sprintf("randomuser123@gmail.com"),
	}
}

func (a *Address) FetchInbox() ([]string, error) {
	// TODO: Implement inbox fetching
	return []string{}, nil
}
