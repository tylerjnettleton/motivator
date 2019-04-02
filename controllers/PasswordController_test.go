package controllers

import (
	"testing"
)

func init () {

}

var TestPassword = []byte("password")

func TestRandBytes(t *testing.T) {
	salt, err := randbytes()
	if err != nil {
		t.Fatal("Rand Bytes returned nil.")
	}

	if len(salt) != 16 {
		t.Fatal("Length of returned random bytes was not 16.")
	}
}

func TestHashPassword(t *testing.T) {
	salt, _ := randbytes()
	hashedPassword, err := HashPassword(TestPassword, salt)

	if err != nil {
		t.Fatal("Rand Bytes returned nil.")
	}

	if len(hashedPassword) == 0 {
		t.Fatal("Hashed password returned a length of 0.")
	}
}

func TestAuthenticate(t *testing.T) {
	salt, _ := randbytes()
	hashedPassword, _ := HashPassword(TestPassword, salt)

	if ok, err := Authenticate(TestPassword, salt, hashedPassword); ok {
		if err != nil {
			t.Fatal("Rand Bytes returned nil.")
		}

		if ok == false {
			t.Fatal("Authentication failed.")
		}
	}
}