package jwt

import (
	"testing"
	"unicode/utf8"
)

func Test_manager_generateToken(t *testing.T) {
	m := manager{}
	token := m.generateToken(0)
	//fmt.Println(token)
	if !utf8.ValidString(token) {
		t.Errorf("invalid utf8 sequanse %s", token)
	}
}
