package eero

import (
	"fmt"

	"github.com/ericraio/eero/internal/client"
)

type Eero interface {
	Accounts() (string, error)
	Login(identifier string) (string, error)
	LoginVerify(verificationCode, userToken string) (string, error)
}

type eero struct {
	client client.Client
}

func New() Eero {
	cookieStore := client.NewCookieStore("cookie.session")
	return &eero{
		client: client.New(cookieStore),
	}
}

func (e *eero) Login(identifier string) (string, error) {
	body := map[string]string{
		"login": identifier,
	}

	data, err := e.client.Post("login", body)
	if err != nil {
		return "", err
	}

	return data["user_token"], nil
}

func (e *eero) LoginVerify(verificationCode, userToken string) (string, error) {
	body := map[string]string{
		"code": verificationCode,
	}
	e.client.SetAuth(userToken)
	data, err := e.client.Post("login/verify", body)
	if err != nil {
		return "", err
	}

	fmt.Println(data)
	return data["user_token"], nil
}

func (e *eero) Accounts() (string, error) {
	data, err := e.client.Get("accounts")
	if err != nil {
		return "", err
	}

	fmt.Println(data)
	return data["user_token"], nil
}
