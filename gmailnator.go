package gmailnator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bunniesnu/gocaptcha"
)

type Gmailnator struct {
	XSRFToken string
	Client    *http.Client
	Email     *Email
	RapidAPI  string
}

type Email struct {
    Email     string `json:"email"`
    Timestamp int64  `json:"timestamp"`
}

type ReceivedEmail struct {
	Mid         string  `json:"mid"`
	From    string  `json:"textFrom"`
	Date    string  `json:"textDate"`
	Subject string  `json:"textSubject"`
	To      string  `json:"textTo"`
}

type ReceivedEmailDetail struct {
	Body string `json:"body"`
}

func NewGmailnator() (*Gmailnator, error) {
	cookie, client, err := NewCookie("https://smailpro.com/temp-gmail")
	if err != nil {
		return nil, err
	}
	return &Gmailnator{
		XSRFToken: cookie.XSRFToken,
		Client: client,
		Email: nil,
		RapidAPI: "",
	}, nil
}	
func (g *Gmailnator) GenerateEmail() error {
	recaptcha, err := gocaptcha.NewRecaptchaV3(RecaptchaAnchorURL, nil, 0)
	if err != nil {
		return errors.New("failed to create recaptcha instance: " + err.Error())
	}
	captchaToken, err := recaptcha.Solve()
	if err != nil {
		return errors.New("failed to solve recaptcha: " + err.Error())
	}
	if g.RapidAPI == "" {
		req, err := http.NewRequest(http.MethodGet, "https://smailpro.com/js/chunks/smailpro_v2_email.js", nil)
		if err != nil {
			return errors.New("failed to create request for rapidapi key: " + err.Error())
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
		req.Header.Set("x-xsrf-token", g.XSRFToken)
		req.Header.Set("x-g-token", captchaToken)
		resp, err := g.Client.Do(req)
        if err != nil {
            return errors.New("failed to fetch rapidapi key: " + err.Error())
        }
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)
        js := string(body)
        parts := strings.Split(js, `rapidapi_key:"`)
        if len(parts) < 2 {
            return errors.New("rapidapi key not found")
        }
        g.RapidAPI = strings.Split(parts[1], `"`)[0]
	}
	key, err := g.GetKey(`{"domain":"gmail.com","username":"random","server":"server-1","type":"alias"}`, captchaToken)
	if err != nil {
		return errors.New("failed to get key for new email: " + err.Error())
	}
	endpoint := fmt.Sprintf("https://public-sonjj.p.rapidapi.com/email/gm/get?key=%s&rapidapi-key=%s&domain=gmail.com&username=random&server=server-1&type=alias", key, g.RapidAPI)
    req, err := http.NewRequest(http.MethodGet, endpoint, nil)
    if err != nil {
        return errors.New("failed to create request for new email: " + err.Error())
    }
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("x-xsrf-token", g.XSRFToken)
	req.Header.Set("x-g-token", captchaToken)
	resp, err := g.Client.Do(req)
	if err != nil {
		return errors.New("failed to fetch new email: " + err.Error())
	}
    defer resp.Body.Close()
    var emailResp struct {
        Items Email `json:"items"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&emailResp); err != nil {
        return errors.New("failed to decode response for new email: " + err.Error())
    }
    g.Email = &emailResp.Items
    return nil
}
func (g *Gmailnator) GetMails() ([]ReceivedEmail, error) {
	recaptcha, err := gocaptcha.NewRecaptchaV3(RecaptchaAnchorURL, nil, 0)
	if err != nil {
		return nil, errors.New("failed to create recaptcha instance: " + err.Error())
	}
	captchaToken, err := recaptcha.Solve()
	if err != nil {
		return nil, errors.New("failed to solve recaptcha: " + err.Error())
	}
	key, err := g.GetKey(fmt.Sprintf(`{"email": "%s", "timestamp": %d}`, g.Email.Email, g.Email.Timestamp), captchaToken)
	if err != nil {
		return nil, errors.New("failed to get key for new email: " + err.Error())
	}
	endpoint := fmt.Sprintf("https://public-sonjj.p.rapidapi.com/email/gm/check?key=%s&rapidapi-key=%s&email=%s&timestamp=%d", key, g.RapidAPI, g.Email.Email, g.Email.Timestamp)
    req, err := http.NewRequest(http.MethodGet, endpoint, nil)
    if err != nil {
        return nil, errors.New("failed to create request for new email: " + err.Error())
    }
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("x-xsrf-token", g.XSRFToken)
	req.Header.Set("x-g-token", captchaToken)
	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, errors.New("failed to fetch new email: " + err.Error())
	}
	defer resp.Body.Close()
    var emailResp struct {
        Items []ReceivedEmail `json:"items"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&emailResp); err != nil {
        return nil, errors.New("failed to decode response for new email: " + err.Error())
    }
	return emailResp.Items, nil
}
func (g *Gmailnator) GetMailBody(messageId string) (string, error) {
	recaptcha, err := gocaptcha.NewRecaptchaV3(RecaptchaAnchorURL, nil, 0)
	if err != nil {
		return "", errors.New("failed to create recaptcha instance: " + err.Error())
	}
	captchaToken, err := recaptcha.Solve()
	if err != nil {
		return "", errors.New("failed to solve recaptcha: " + err.Error())
	}
	key, err := g.GetKey(fmt.Sprintf(`{"email": "%s", "message_id": "%s"}`, g.Email.Email, messageId), captchaToken)
	if err != nil {
		return "", errors.New("failed to get key for new email: " + err.Error())
	}
	endpoint := fmt.Sprintf("https://public-sonjj.p.rapidapi.com/email/gm/read?key=%s&rapidapi-key=%s&email=%s&message_id=%s", key, g.RapidAPI, g.Email.Email, messageId)
    req, err := http.NewRequest(http.MethodGet, endpoint, nil)
    if err != nil {
        return "", errors.New("failed to create request for new email: " + err.Error())
    }
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("x-xsrf-token", g.XSRFToken)
	req.Header.Set("x-g-token", captchaToken)
	resp, err := g.Client.Do(req)
	if err != nil {
		return "", errors.New("failed to fetch new email: " + err.Error())
	}
	defer resp.Body.Close()
	var emailResp struct {
		Items ReceivedEmailDetail `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emailResp); err != nil {
		return "", errors.New("failed to decode response for new email: " + err.Error())
	}
	return emailResp.Items.Body, nil
}