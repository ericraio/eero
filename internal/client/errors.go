package client

import "fmt"

type ErrResponse struct {
	Msg string
}

func (er *ErrResponse) Error() string {
	return fmt.Sprintf("eero client error: %v", er.Msg)
}
