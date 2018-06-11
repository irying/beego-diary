package util

import "testing"

func TestParseToken(t *testing.T) {
	username := "hello"
	password := "world"
	token, _ := GenerateToken(username, password)
	_, err2 := ParseToken(token)

	if err2 != nil {
		t.Error("fail token parsing")
	} else {
		t.Log("pass")
	}
}
