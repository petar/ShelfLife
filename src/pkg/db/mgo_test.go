// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package db

import (
	"testing"
	"github.com/petar/ShelfLife/thirdparty/mgo"
)

type IDType []byte

type Doc struct {
	ID []byte
}

func TestByteSlices(t *testing.T) {
	s, err := mgo.Mongo("localhost:22000")
	if err != nil {
		t.Fatalf("connect: %s", err)
	}
	c := s.DB("tester").C("mgo_test")
	doc := &Doc{ []byte{1,2,3,4} }
	if err = c.Insert(doc); err != nil {
		t.Errorf("insert: %s", err)
	}
}
