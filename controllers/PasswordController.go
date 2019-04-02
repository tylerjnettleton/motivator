package controllers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"os"
	"unicode"
)

const (
	// scrypt is used for strong keys
	// these are the recommended scrypt parameters
	scryptN      = 16384
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 32
)

func randbytes() ([]byte, error) {
	// Generate a salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return salt, err
}

func hmac_sha256(in, salt []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, salt)
	_, err := mac.Write(in)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func enc_scrypt(in, salt []byte) ([]byte, error) {
	return scrypt.Key(in, salt, scryptN, scryptR, scryptP, scryptKeyLen)
}

func HashPassword(password, salt []byte) ([]byte, error) {
	peppered, _ := hmac_sha256(password, []byte(os.Getenv("TS_PEPPER")))
	cur, _ := enc_scrypt(peppered, salt)

	return cur, nil
}

func Authenticate(password, salt, hash []byte) (bool, error) {
	h, _ := HashPassword(password, salt)

	if subtle.ConstantTimeCompare(h, hash) != 1 {
		return false, fmt.Errorf("Invalid username or password")
	}

	return true, nil
}

// Todo: Unit test this!
func VerifyPassword(s string) (sevenOrMore, number, upper, special bool) {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			//return false, false, false, false
		}
	}
	sevenOrMore = letters >= 7
	return
}
