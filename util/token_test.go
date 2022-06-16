package util

import (
	"fmt"
	"testing"
)

func TestCreateToken(t *testing.T) {
	token, err := CreateToken(102)
	if err != nil {
		t.Error("err", err)
		return
	}
	fmt.Println(token)
}

func TestParseToken(t *testing.T) {
	claim, err := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQ4NzkxNDMsIlVzZXJJRCI6NjY2Nn0.DZqICNnAj7GJ4yDs0E0LiVIwVySQpTM_3MsmQXwuzn8")
	if err != nil {
		t.Error("err", err)
		return
	}
	fmt.Println(claim.UserID)
}
