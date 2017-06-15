package main

import (
	"fmt"
	// "log"
	// "net/http"
)

type ServiceError struct {
	Msg string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

func Error(msg string) error {
	return &ServiceError{msg}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	wx := wxweb{}
	wx.start()

}
