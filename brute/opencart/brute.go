package opencart

import (
	"net/http"
	"net/url"
	"strings"
)

func Start(link string, username string, password string) bool {
	data := url.Values{
		"username": {username},
		"password": {password},
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	//client.Transport = &http.Transport{
	//	Proxy: http.ProxyURL(&url.URL{
	//		Scheme: "http",
	//		User:   url.UserPassword("login", "password"),
	//		Host:   "IP:PORT",
	//	}),
	//}

	req, err := http.NewRequest("POST", link+"index.php?route=common/login", strings.NewReader(data.Encode()))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:105.0) Gecko/20100101 Firefox/105.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode == 302 {
		return true
	}
	return false
}
