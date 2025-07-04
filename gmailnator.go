package gmailnator

const defaultURL = "https://www.emailnator.com/"

type Gmailnator struct {
	BaseURL string
	CookieData *CookieData
	Email string
}

func NewGmailnator() *Gmailnator {
	cookie := NewCookie(defaultURL)
	if cookie == nil {
		return nil
	}	
	return &Gmailnator{
		BaseURL: defaultURL,
		CookieData: cookie,
		Email: "",
	}
}	