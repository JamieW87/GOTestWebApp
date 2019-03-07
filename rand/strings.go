//Package rand is a wrapping around crypto/rand.
//Crypto/rand is a package that helps generate random strings cryptographically, so they are secure.
//We will use this to simplify the crypto/rand package and use it for our package.
package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// RememberToken is a helper function designed to generate
// remember tokens of a predetermined byte size.
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

// Bytes will help us generate n random bytes, or will
// return an error if there was one.
//Takes in an integer, n, and created a byte slice of that length.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//String returns a random string that is base64 url encoded.
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
