// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"crypto/hmac"
	"encoding/base64"
)

// SanitizeLogin verifies the syntax, and if successful, sanitizes the login string
func SanitizeLogin(login string) (string, os.Error) {
	?
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

// SanitizeLogin verifies the syntax, and if successful, sanitizes the name string
func SanitizeName(s string) (string, os.Error) {
	?
}

// SanitizeLogin verifies the syntax, and if successful, sanitizes the email string
func SanitizeEmail(s string) (string, os.Error) {
	?
}
