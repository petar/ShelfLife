// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"os"
	"github.com/petar/GoHTTP/http"
	"github.com/petar/GoHTTP/server/rpc"
	"github.com/petar/ShelfLife/db"
)

// IsFollow returns true if the logged user follows the given object
func (a *API) IsFollow(args *rpc.Args, r *rpc.Ret) (err os.Error) {

	user, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrApp
	}

	what, err := args.String("What")
	if err != nil {
		return ErrApp
	}
	attr, _ := args.String("Attr")

	ok, err = a.db.IsFollow(user.ID, what, attr)
	if err != nil {
		return err
	}
	r.SetBool("IsFollow", ok)

	return nil
}

// Follow records that the currently logged user is following the given object
func (a *API) Follow(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	
	user, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrApp
	}

	what, err := args.String("What")
	if err != nil {
		return ErrApp
	}
	attr, _ := args.String("Attr")

	ok, err = a.db.IsFollow(user.ID, what, attr)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	if err = a.db.Follow(user.ID, what, attr); err != nil {
		return err
	}

	return nil
}
