// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package sociability

import (
	//"log"
	"os"
	"github.com/petar/GoHTTP/http"
	"github.com/petar/GoHTTP/server/rpc"
	"github.com/petar/ShelfLife/thirdparty/authcookie"
	"github.com/petar/ShelfLife/db"
)

// RPC/SignInLogin logs in a user, specified by their login (aka username)
// Args:
//   "L" string
//   "P" string = HMAC-hashed password
// Err:
//   ErrApp:  If the sign-in information is incorrect
//   non-nil: If a technical problem occured
//
func (a *API) SignInLogin(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	
	// Validate and sanitize arguments
	login, _ := args.QueryString("L")
	if login, err = SanitizeLogin(login); err != nil {
		return ErrApp
	}
	hpass, _ := args.QueryString("P")

	// Fetch user for this login
	u, err := a.db.FindUserByLogin(login)
	if err != nil {
		return ErrDb
	}
	if u == nil {
		return ErrApp
	}

	// Verify credentials
	if !VerifyPassword(hpass, u.Password) {
		return ErrSec
	}

	// Set authentication cookie
	r.AddSetCookie(a.newUserAuthCookie(u))
	// Set user info cookie
	r.AddSetCookie(a.newUserInfoCookie(u))

	return nil
}

const (
	OneDayInSec  = 60*60*24
	OneWeekInSec = OneDayInSec*7
)

// newUserAuthCookie returns a new cookie authenticating that the given 
// user is signed in
func (a *API) newUserAuthCookie(u *db.UserDoc) *http.Cookie {
	duration := OneWeekInSec
	return &http.Cookie{
		Name:   "SS-UserAuth",
		Value:  authcookie.NewSinceNow(u.Login, int64(duration), a.loginSecret),
		Path:   "/",
		MaxAge: duration,
	}
}

// newUserInfoCookie returns a new cookie with user information
func (a *API) newUserInfoCookie(u *db.UserDoc) *http.Cookie {
	duration := 2*OneWeekInSec
	return &http.Cookie{
		Name:   "SS-UserInfo",
		Value:  u.Name,
		Path:   "/",
		MaxAge: duration,
	}
}

// verifySignInCookie checks that cookie is a valid authentication cookie,
// and if so returns the user who is logged in with this cookie, or nil otherwise.
// A non-nil error indicates a technical problem.
func (a *API) verifySignInCookie(cookie *http.Cookie) (user *db.UserDoc, err os.Error) {
	if cookie == nil || cookie.Name != "SS-UserAuth" {
		return nil, nil
	}
	login := authcookie.Login(cookie.Value, a.loginSecret)
	if login, err = SanitizeLogin(login); err != nil {
		return nil, nil
	}
	user, err = a.db.FindUserByLogin(login)
	if err != nil {
		return nil, ErrDb
	}
	return user, nil
}

func (a *API) whoAmI(args *rpc.Args) (user *db.UserDoc, err os.Error) {
	for _, cookie := range args.Cookies {
		user, err = a.verifySignInCookie(cookie)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}
	return nil, nil
}

// WhoAmI returns the login of the currently signed user
func (a *API) WhoAmI(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	user, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	login := ""
	if user != nil {
		login = user.Login
	}
	r.SetString("Login", login)

	return nil
}

// RPC/SignInEmail logs in a user, specified by their email
// Args:
//   "E" string
//   "P" string = HMAC-hashed password
// Err:
//   ErrApp:  If the sign-in information is incorrect
//   non-nil: If a technical problem occured
//
func (a *API) SignInEmail(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	
	// Validate and sanitize arguments
	email, _ := args.QueryString("E")
	if email, err = SanitizeEmail(email); err != nil {
		return ErrApp
	}
	hpass, _ := args.QueryString("P")

	// Fetch user for this login
	u, err := a.db.FindUserByEmail(email)
	if err != nil {
		return ErrDb
	}
	if u == nil {
		return ErrApp
	}

	// Verify credentials
	if !VerifyPassword(hpass, u.Password) {
		return ErrSec
	}

	// Set authentication cookie
	r.AddSetCookie(a.newUserAuthCookie(u))
	// Set user info cookie
	r.AddSetCookie(a.newUserInfoCookie(u))

	return nil
}

// SignUp registers a new user
func (a *API) SignUp(args *rpc.Args, r *rpc.Ret) (err os.Error) {

	// Validate and sanitize arguments
	name, _ := args.QueryString("N")
	if name, err = SanitizeName(name); err != nil {
		return ErrApp
	}
	email, _ := args.QueryString("E")
	if email, err = SanitizeEmail(email); err != nil {
		return ErrApp
	}
	login, _ := args.QueryString("L")
	if login, err = SanitizeLogin(login); err != nil {
		return ErrApp
	}
	hpass, _ := args.QueryString("P")

	// Check that a user like this doesn't already exist
	u, err := a.db.FindUserByLogin(login)
	if err != nil {
		return ErrDb
	}
	if u != nil {
		return ErrApp
	}
	u, err = a.db.FindUserByEmail(email)
	if err != nil {
		return ErrDb
	}
	if u != nil {
		return ErrApp
	}

	// Add the user
	u = &db.UserDoc{
		Name:     name,
		Login:    login,
		Email:    email,
		Password: hpass,
	}
	if _, err = a.db.AddUser(u); err != nil {
		return ErrDb
	}

	return nil
}

// IsLoginAvailable checks if this login (i.e. username) is already taken
func (a *API) IsLoginAvailable(args *rpc.Args, r *rpc.Ret) os.Error {
	login, err := args.QueryString("L")
	if err != nil {
		return err
	}
	if login, err = SanitizeLogin(login); err != nil {
		return rpc.ErrArg
	}
	u, err := a.db.FindUserByLogin(login)
	if err != nil {
		return ErrDb
	}
	r.SetBool("Available", u == nil)
	return nil
}
