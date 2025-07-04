package gmailnator

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const defaultURL = "https://www.emailnator.com/"

type Gmailnator struct {
	BaseURL string
	Cookie *CookieData
	Client *http.Client
	Email string
}

func NewGmailnator() (*Gmailnator, error) {
	cookie, client, err := NewCookie(defaultURL)
	if err != nil {
		return nil, err
	}	
	return &Gmailnator{
		BaseURL: defaultURL,
		Cookie: cookie,
		Client: client,
		Email: "",
	}, nil
}	
func (g *Gmailnator) GenerateEmail() (string, error) {
	payload := map[string][]string{"email": {"dotGmail"}}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", g.BaseURL + "generate-email", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", g.Cookie.XSRFToken)
	req.Header.Set("Origin", g.BaseURL)
	req.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: g.Cookie.XSRFToken})
	req.AddCookie(&http.Cookie{Name: "gmailnator_session", Value: g.Cookie.GmailnatorSession})

	resp, err := g.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string][]string
	json.NewDecoder(resp.Body).Decode(&result)

	return result["email"][0], nil
}