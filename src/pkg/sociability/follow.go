// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package sociability

import (
	//"log"
	"os"
	"github.com/petar/GoHTTP/server/rpc"
)

// FollowInfo returns true if the logged user follows the given object
func (a *API) FollowInfo(args *rpc.Args, r *rpc.Ret) (err os.Error) {

	_, uid, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	what, _ := args.QueryString("What")

	follows, err := a.db.IsFollow(uid, what)
	if err != nil {
		follows = false
	}

	n, err := a.db.FollowerCount(what)
	if err != nil {
		return err
	}

	r.SetBool("Follows", follows)
	r.SetInt("Count", n)

	return nil
}

// Follow makes the currently logged user follow the given object
func (a *API) SetFollow(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	_, uid, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	what, _ := args.QueryString("What")
	r.SetInt("XPad", 0)
	return a.db.SetFollow(uid, what)
}

func (a *API) UnsetFollow(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	_, uid, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	what, _ := args.QueryString("What")
	r.SetInt("XPad", 0)
	return a.db.UnsetFollow(uid, what)
}
