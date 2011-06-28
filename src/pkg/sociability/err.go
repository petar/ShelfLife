// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package sociability

import (
	"os"
)

var (
	ErrParse = os.NewError("bad or missing RPC call arguments")
	ErrDb    = os.NewError("database error")
	ErrApp   = os.NewError("operation denied by app")
	ErrSec   = os.NewError("operation denied for security reasons")
)
