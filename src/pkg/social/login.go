// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"os"
	"github.com/dchest/authcookie"
	"github.com/petar/GoHTTP/server/rpc"
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
func (a *API) SignInLogin(args *rpc.Args, r *rpc.Ret) os.Error {
	
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
	r.AddCookie(?)

	return nil
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
func (a *API) SignInEmail(args *rpc.Args, r *rpc.Ret) os.Error {
	?
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
	login, err := args.Bool("Login")
	if err != nil {
		return err
	}
	if !IsValidLogin(login) {
		r.Bool("Have", false)
		return nil
	}
	u, err := a.db.FindUserByLogin(login)
	if err != nil {
		return ErrDb
	}
	r.Bool("Have", u != nil)
	return nil
}

// If cookie is a valid authentication cookie, getCookieCredentials returns the user
// who is logged in with this cookie and nil otherwise. 
// A non-nil error indicates a technical problem.
func (a *API) getCookieCredentials(cookie *Cookie) (user *db.User, err os.Error) {
	?
}
