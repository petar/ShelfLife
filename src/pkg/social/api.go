// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"os"
	"github.com/petar/GoHTTP/server/rpc"
	"github.com/petar/ShelfLife/db"
)

type API struct {
	db          *db.Db
	loginSecret []byte
}

func NewAPI(db *db.Db, loginSecret []byte) *API { 
	return &API{ 
		db:          db, 
		loginSecret: loginSecret,
	} 
}

func (a *API) Ping(args *rpc.NoArgs, r *rpc.NoRet) os.Error {
	return nil
}

func (a *API) HelloWorld(args *rpc.NoArgs, r *rpc.ShortRet) os.Error {
	r.Value = make(map[string]string)
	r.Value["Hello"] = "World"
	return nil
}
