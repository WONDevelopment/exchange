package pkg

import "fmt"

type WonError struct {
	Code    string `json:"error"`
	Message string `json:"error_description"`
}

func (e WonError) Error() string {
	return fmt.Sprintf("code:%s,message:%s", e.Code, e.Message)
}
