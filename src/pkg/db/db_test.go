// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package db

import (
	"log"
	"testing"
)

func TestUser(t *testing.T) {
	db, err := NewDb("localhost:22000", "tester")
	if err != nil {
		t.Fatalf("new db: %s", err)
	}
	
	u := &UserDoc{
		Name:         "petar maymounkov",
		Login:        "ma51",
		Email:        "petar@5ttt.org",
		Password:     "aaa",
	}
	uid, err := db.AddUser(u)
	if err != nil {
		t.Errorf("add user: %s", err)
	}
	log.Printf("added #%v", uid)

	if err = db.Close(); err != nil {
		t.Fatalf("close: %s", err)
	}
}
