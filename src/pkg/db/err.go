// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package db

import (
	"os"
)

var (
	ErrSec   = os.NewError("database denied for security reasons")
)
