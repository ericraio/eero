package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const apiEndpoint = "https://api-user.e2ro.com/2.2/%s"

type Client interface {
	Post(string, map[string]string) (map[string]string, error)
	Get(string) (map[string]string, error)
	SetAuth(userToken string) error
}

type client struct {
	cookieStore CookieStore
	http        *http.Client
}

func New(cookieStore CookieStore) Client {
	return &client{
		cookieStore: cookieStore,
		http: &http.Client{
			Jar: cookieStore.Jar(),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

type response struct {
	Error struct {
		RequestID int    `json:"requestId"`
		Message   string `json:"message"`
	} `json:"error"`
	Meta struct {
		Code       int       `json:"code"`
		ServerTime time.Time `json:"server_time"`
		Error      string    `json:"error"`
	} `json:"meta"`
	Data map[string]string
}

func (c *client) parse(r *http.Response) (map[string]string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, &ErrResponse{Msg: err.Error()}
	}
	resp := &response{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, &ErrResponse{Msg: err.Error()}
	}

	if resp.Error.Message != "" {
		return nil, &ErrResponse{Msg: resp.Error.Message}
	}

	if resp.Meta.Code != 200 && resp.Meta.Code != 201 {
		return nil, &ErrResponse{Msg: resp.Meta.Error}
	}

	return resp.Data, nil
}

func (c *client) SetAuth(userToken string) error {
	err := c.cookieStore.Add(&http.Cookie{
		Name:  "s",
		Value: userToken,
	})

	if err != nil {
		return err
	}

	c.http.Jar = c.cookieStore.Jar()
	return nil
}

func (c *client) Post(action string, body map[string]string) (map[string]string, error) {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		return nil, &ErrResponse{Msg: err.Error()}
	}

	r, err := c.http.Post(fmt.Sprintf(apiEndpoint, action), "application/json", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return nil, &ErrResponse{Msg: err.Error()}
	}

	return c.parse(r)
}

func (c *client) Get(action string) (map[string]string, error) {
	r, err := c.http.Get(fmt.Sprintf(apiEndpoint, action))
	if err != nil {
		return nil, &ErrResponse{Msg: err.Error()}
	}
	return c.parse(r)
}
