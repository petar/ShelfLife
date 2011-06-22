// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"os"
	"github.com/dchest/authcookie"
	"github.com/petar/GoHTTP/server/rpc"
)

// RPC/SignInUser logs in a user, specified by their username
// Args:
//   "Login" string = username
//   "HashPassword" string = hashed password
// Ret:
//    
func (a *API) SignInUser(args *rpc.ShortCookieArgs, r *rpc.SetCookieRet) os.Error {
	return nil
}

func (a *API) SignInEmail(args *rpc.ShortCookieArgs, r *rpc.SetCookieRet) os.Error {
	return nil
}

func (a *API) SignUp(args *rpc.ShortCookieArgs, r *rpc.ShortSetCookieRet) os.Error {
	return nil
}

func (a *API) HaveLogin(args *rpc.ShortArgs, r *rpc.ShortRet) os.Error {
	login, err := rpc.GetBool(args, "Login")
	if err != nil {
		return err
	}
	if !IsValidLogin(login) {
		rpc.SetBool(r, "Have", false)
		return nil
	}
	u, err := a.db.FindUserByLogin(login)
	if err != nil {
		return ErrDb
	}
	rpc.SetBool(r, "Have", u != nil)
	return nil
}

func (a *API) getCookieCredentials(cookie *Cookie) (login, email string, err os.Error) {
}
