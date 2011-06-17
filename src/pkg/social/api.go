// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package social

import (
	"os"
)

type API struct {
	//db *Db
}

func NewAPI() *API { return &API{} }

type Empty struct {}

func (a *API) Ping(args *Empty, r *Empty) os.Error {
	return nil
}

func (a *API) HelloWorld(args *Empty, r *string) os.Error {
	*r = "Hello world!"
	return nil
}
