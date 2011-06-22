// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"crypto/hmac"
	"encoding/base64"
)

const loginOpAlphabet = ",.-_=+"

func IsValidLogin(login string) bool {
	if len(login) == 0 {
		return false
	}
	for _, c := range login {
		switch {
		case 32 <= c && c <= 126:
		default:
			return false
		}
	}
	return true
}

// passwordHMACKey is the HMAC key for the one-way transformation of plaintext passwords,
// before they are stored in the user database. This key does not have to be secret.
var passwordHMACKey = []byte{ 0x12, 0x13, 0x16, 0x18 }

func HashPassword(password string) string {
	sha256 := hmac.NewSHA256(passwordHMACKey)
	sha256.Write([]byte(password))
	hmac := sha256.Sum()
	return textify(hmac)
}

// textify converts a byte slice into textual representation, using base64 encoding
func textify(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}
