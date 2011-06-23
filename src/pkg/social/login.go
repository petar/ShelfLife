// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"os"
	"github.com/dchest/authcookie"
	"github.com/petar/GoHTTP/http"
	"github.com/petar/GoHTTP/server/rpc"
	"github.com/petar/ShelfLife/db"
)

// RPC/SignInLogin logs in a user, specified by their login (aka username)
// Args:
//   "Login" string
//   "HPass" string = HMAC-hashed password
// Ret: n/a
// Err:
//   ErrApp:  If the sign-in information is incorrect
//   non-nil: If a technical problem occured
//
func (a *API) SignInLogin(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	
	// Validate and sanitize arguments
	login, _ := args.String("Login")
	if login, err = SanitizeLogin(login); err != nil {
		return ErrApp
	}
	hpass, _ := args.String("HPass")

	// Fetch user for this login
	u, err := a.db.FindUserByLogin(login)
	if err != nil {
		return ErrDb
	}
	if u == nil {
		return ErrApp
	}

	// Verify credentials
	if !VerifyPassword(hpass, u.HashPassword) {
		return ErrSec
	}

	// Set authentication cookie
	r.AddSetCookie(a.newSignInCookie(u))

	return nil
}

const (
	OneDayInSec  = 60*60*24
	OneWeekInSec = OneDayInSec*7
)

// newSignInCookie returns a new cookie authenticating that the given 
// user is signed in
func (a *API) newSignInCookie(u *db.User) *http.Cookie {
	duration := OneWeekInSec
	return &http.Cookie{
		Name:   "Login",
		Value:  authcookie.NewSinceNow(u.Login, int64(duration), a.loginSecret),
		MaxAge: duration,
	}
}

// verifySignInCookie checks that cookie is a valid authentication cookie,
// and if so returns the user who is logged in with this cookie, or nil otherwise.
// A non-nil error indicates a technical problem.
func (a *API) verifySignInCookie(cookie *http.Cookie) (user *db.User, err os.Error) {
	if cookie == nil || cookie.Name != "Login" {
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

// RPC/SignInEmail logs in a user, specified by their email
// Args:
//   "Email" string
//   "HPass" string = HMAC-hashed password
// Ret: n/a
// Err:
//   ErrApp:  If the sign-in information is incorrect
//   non-nil: If a technical problem occured
//
func (a *API) SignInEmail(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	
	// Validate and sanitize arguments
	email, _ := args.String("Email")
	if email, err = SanitizeEmail(email); err != nil {
		return ErrApp
	}
	hpass, _ := args.String("HPass")

	// Fetch user for this login
	u, err := a.db.FindUserByEmail(email)
	if err != nil {
		return ErrDb
	}
	if u == nil {
		return ErrApp
	}

	// Verify credentials
	if !VerifyPassword(hpass, u.HashPassword) {
		return ErrSec
	}

	// Set authentication cookie
	r.AddSetCookie(a.newSignInCookie(u))

	return nil
}

// RPC/SignUp registers a new user
// Args:
//   "Name"  string
//   "Email" string
//   "Login" string
//   "HPass" string = HMAC-hashed password
// Ret: n/a
// Err:
//   ErrApp:  If the application logic prohibits this registration
//   non-nil: If a technical problem occured
//
func (a *API) SignUp(args *rpc.Args, r *rpc.Ret) (err os.Error) {

	// Validate and sanitize arguments
	name, _ := args.String("Name")
	if name, err = SanitizeName(name); err != nil {
		return ErrApp
	}
	email, _ := args.String("Email")
	if email, err = SanitizeEmail(email); err != nil {
		return ErrApp
	}
	login, _ := args.String("Login")
	if login, err = SanitizeLogin(login); err != nil {
		return ErrApp
	}
	hpass, _ := args.String("HPass")

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
	u = &db.User{
		Name: name,
		Login: login,
		Email: email,
		HashPassword: hpass,
	}
	if err = a.db.AddUser(u); err != nil {
		return ErrDb
	}

	return nil
}

// RPC/HaveLogin checks if this login (i.e. username) is already taken
// Args:
//   "Login" string
// Ret:
//   "Have" bool
// Err:
//   non-nil: If a technical problem occured
//
func (a *API) HaveLogin(args *rpc.Args, r *rpc.Ret) os.Error {
	login, err := args.String("Login")
	if err != nil {
		return err
	}
	if login, err = SanitizeLogin(login); err != nil {
		r.SetBool("Have", false)
		return nil
	}
	u, err := a.db.FindUserByLogin(login)
	if err != nil {
		return ErrDb
	}
	r.SetBool("Have", u != nil)
	return nil
}
