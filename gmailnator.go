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
	Client *http.Client
	Email string
	RapidAPI string
	Key string
}

type newEmailResult struct {
    Email     string `json:"email"`
    Timestamp int64 `json:"timestamp"`
}

func NewGmailnator() (*Gmailnator, error) {
	cookie, client, err := NewCookie("https://smailpro.com/temp-gmail")
	if err != nil {
		return nil, err
	}
	return &Gmailnator{
		XSRFToken: cookie.XSRFToken,
		Client: client,
		Email: "",
		RapidAPI: "",
	}, nil
}	
func (g *Gmailnator) GenerateEmail() (string, error) {
	recaptcha, err := gocaptcha.NewRecaptchaV3(RecaptchaAnchorURL, nil, 0)
	if err != nil {
		return "", errors.New("failed to create recaptcha instance: " + err.Error())
	}
	captchaToken, err := recaptcha.Solve()
	if err != nil {
		return "", errors.New("failed to solve recaptcha: " + err.Error())
	}
	if g.RapidAPI == "" {
		req, err := http.NewRequest(http.MethodGet, "https://smailpro.com/js/chunks/smailpro_v2_email.js", nil)
		if err != nil {
			return "", errors.New("failed to create request for rapidapi key: " + err.Error())
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
		req.Header.Set("x-xsrf-token", g.XSRFToken)
		req.Header.Set("x-g-token", captchaToken)
		resp, err := g.Client.Do(req)
        if err != nil {
            return "", errors.New("failed to fetch rapidapi key: " + err.Error())
        }
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)
        js := string(body)
        parts := strings.Split(js, `rapidapi_key:"`)
        if len(parts) < 2 {
            return "", errors.New("rapidapi key not found")
        }
        g.RapidAPI = strings.Split(parts[1], `"`)[0]
	}
	if g.Key == "" {
        jsonData := `{"domain":"gmail.com","username":"random","server":"server-1","type":"alias"}`
        reqBody := strings.NewReader(jsonData)
		req, err := http.NewRequest(http.MethodPost, "https://smailpro.com/app/key", reqBody)
        if err != nil {
			return "", errors.New("failed to create request for new email key: " + err.Error())
        }
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
		req.Header.Set("x-xsrf-token", g.XSRFToken)
		req.Header.Set("x-g-token", captchaToken)
		resp, err := g.Client.Do(req)
		if err != nil {
			return "", errors.New("failed to fetch new email key: " + err.Error())
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New("failed to read response body for new email key: " + err.Error())
		}
		resp.Body = io.NopCloser(strings.NewReader(string(body)))
        defer resp.Body.Close()
        var parsed struct {
            Items string `json:"items"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
            return "", errors.New("failed to decode response for new email key: " + err.Error())
        }
        g.Key = parsed.Items
    }
	endpoint := fmt.Sprintf("https://public-sonjj.p.rapidapi.com/email/gm/get?key=%s&rapidapi-key=%s&domain=gmail.com&username=random&server=server-1&type=alias", g.Key, g.RapidAPI)
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to read response body for new email key: " + err.Error())
	}
	resp.Body = io.NopCloser(strings.NewReader(string(body)))
    defer resp.Body.Close()
    var emailResp struct {
        Items newEmailResult `json:"items"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&emailResp); err != nil {
        return "", errors.New("failed to decode response for new email: " + err.Error())
    }
    g.Email = emailResp.Items.Email
    return emailResp.Items.Email, nil
}