//Package hash is a wrapper around the Go hmac package.
//Simplifying it for the purpose of this app
//This package is used to hash our remember tokens for the user cookies.
package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

//HMAC type
type HMAC struct {
	hmac hash.Hash
}

//NewHMAC created and returns a new HMAC object
func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{
		hmac: h,
	}
}

//Hash takes the provided string and hashes it using Hmac.
func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
