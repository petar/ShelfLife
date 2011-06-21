// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"os"
	//"github.com/dchest/authcookie"
	"github.com/petar/GoHTTP/server/rpc"
)

func (a *API) Login(args *rpc.ShortCookieArgs, r *rpc.SetCookieRet) os.Error {
	return nil
}

func (a *API) Register(args *rpc.ShortCookieArgs, r *rpc.ShortSetCookieRet) os.Error {
	return nil
}
