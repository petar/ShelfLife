// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package sociability

import (
	"fmt"
	"github.com/petar/ShelfLife/thirdparty/bson"
	"testing"
)

func TestEnc(t *testing.T) {
	//b = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 2, 2}
	id := bson.NewObjectId()
	fmt.Printf("%s %v", id, id)
}
