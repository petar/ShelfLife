// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"os"
	"strings"
)

// SanitizeLogin verifies the syntax, and if successful, sanitizes the login string
func SanitizeLogin(login string) (string, os.Error) {
	login = strings.TrimSpace(login)
	if len(login) == 0 || len(login) > 200 {
		return "", ErrParse
	}
	for _, c := range login {
		switch {
		case 32 <= c && c <= 126:
		default:
			return "", ErrParse
		}
	}
	return login, nil
}

// SanitizeLogin verifies the syntax, and if successful, sanitizes the name string
func SanitizeName(name string) (string, os.Error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 || len(name) > 200 {
		return "", ErrParse
	}
	return name, nil
}

// SanitizeLogin verifies the syntax, and if successful, sanitizes the email string
// TODO: This should be made more specific according to the appropriate RFC
func SanitizeEmail(email string) (string, os.Error) {
	email = strings.TrimSpace(email)
	if len(email) == 0 || len(email) > 200 {
		return "", ErrParse
	}
	i := strings.Index(email, "@")
	if i < 0 {
		return "", ErrParse
	}
	j := strings.LastIndex(email, "@")
	if i != j {
		return "", ErrParse
	}
	return email, nil
}
