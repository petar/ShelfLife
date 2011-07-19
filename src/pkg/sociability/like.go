// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package sociability

import (
	//"log"
	"os"
	"github.com/petar/GoHTTP/server/rpc"
)

func (a *API) LikeInfo(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	fid, _ := args.QueryString("FID")
	_, uid, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	likes, err := a.db.Likes(uid, fid)
	if err != nil {
		return err
	}
	n, err := a.db.LikeCount(fid)
	if err != nil {
		return err
	}
	r.SetBool("Likes", likes)
	r.SetInt("Count", n)
	return nil
}

func (a *API) Like(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	fid, _ := args.QueryString("FID")
	_, uid, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	return a.db.Like(uid, fid)
}

func (a *API) Unlike(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	fid, _ := args.QueryString("FID")
	_, uid, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	return a.db.Unlike(uid, fid)
}
