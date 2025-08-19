package logto

import (
	"fmt"
	"testing"
)

func TestLogto_Parse(t *testing.T) {
	logto, err := NewLogto("login.randome.chat")
	if err != nil {
		panic(err)
	}
	token, err := logto.Parse("")
	if err != nil {
		panic(err)
	}
	fmt.Println(token.Valid)
	claims := logto.Claims(token)

	fmt.Println(claims)
}
