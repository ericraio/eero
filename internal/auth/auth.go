package auth

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

func login() {
	cookies := make(map[string][]*http.Cookie)
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	r, err := client.Get("https://www.microsoft.com/en-us/")
	if err != nil {
		fmt.Println(err)
		return
	}
	siteCookies := jar.Cookies(r.Request.URL)
	cookies[r.Request.URL.String()] = siteCookies
}
