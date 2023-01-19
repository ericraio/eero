package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ericraio/eero"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var response string
	var err error
	e := eero.New()

	for {
		fmt.Println("enter your phone or email address")

		response, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		break
	}

	token, err := e.Login(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("verification key from email or SMS: ")

		response, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		break
	}

	_, err = e.LoginVerify(response, token)
	if err != nil {
		fmt.Println(err)
	}

	_, err = e.Accounts()
	if err != nil {
		fmt.Println(err)
	}

}
