// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package sociability

import (
	"encoding/hex"
	"github.com/petar/ShelfLife/thirdparty/bson"
)

func WebStringOfObjectID(id bson.ObjectId) string {
	return hex.EncodeToString([]byte(id))
}

func ObjectIDOfWebString(s string) bson.ObjectId {
	b, err := hex.DecodeString(s)
	if err != nil {
		return ""
	}
	return bson.ObjectId(b)
}
