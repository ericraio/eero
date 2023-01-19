package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type CookieStore interface {
	Jar() *cookiejar.Jar
	Add(*http.Cookie) error
	Set([]*http.Cookie)
	Get() []*http.Cookie
	Save() error
}

type cookieStore struct {
	cookieFile string
	cookies    []*http.Cookie
	jar        *cookiejar.Jar
}

func NewCookieStore(cookieFile string) CookieStore {
	jar, _ := cookiejar.New(nil)
	cookies := make([]*http.Cookie, 0)

	return &cookieStore{
		cookieFile,
		cookies,
		jar,
	}
}

func (cs *cookieStore) Jar() *cookiejar.Jar {
	return cs.jar
}

func (cs *cookieStore) Add(cookie *http.Cookie) error {
	cs.cookies = append(cs.cookies, cookie)
	cs.jar.SetCookies(&url.URL{}, cs.cookies)
	return cs.Save()
}

func (cs *cookieStore) Set(cookies []*http.Cookie) {
	cs.cookies = cookies
	cs.jar.SetCookies(&url.URL{}, cs.cookies)
}

func (cs *cookieStore) Get() []*http.Cookie {
	return cs.cookies
}

func (cs *cookieStore) Save() error {
	data, err := json.Marshal(cs.cookies)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cs.cookieFile, data, 0644)
}
