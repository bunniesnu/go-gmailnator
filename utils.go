package gmailnator

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func (g *Gmailnator) GetKey(jsonData, captchaToken string) (string, error) {
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
	return parsed.Items, nil
}