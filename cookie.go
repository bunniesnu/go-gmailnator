package gmailnator

import (
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type CookieData struct {
	XSRFToken         string
	GmailnatorSession string
}

func NewCookie(baseURL string) (*CookieData, *http.Client, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	resp, err := client.Get(baseURL)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	cookies := jar.Cookies(resp.Request.URL)
	data := new(CookieData)
	for _, cookie := range cookies {
		if cookie.Name == "XSRF-TOKEN" {
			data.XSRFToken = strings.ReplaceAll(cookie.Value, "%3D", "=")
		}
		if cookie.Name == "gmailnator_session" {
			data.GmailnatorSession = cookie.Value
		}
	}
	return data, client, nil
}