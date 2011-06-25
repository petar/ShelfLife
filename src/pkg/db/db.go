// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package db

import (
	"log"
	"os"
	"launchpad.net/mgo"
)

// Db encapsulates all interaction with the database
type Db struct {
	// Session with the Db server
	s *mgo.Session

	// Users collection
	u mgo.Collection

	// Follow collection
	f mgo.Collection
}

// NewDb creates a new Db interface to the database.
// addr is the IP address and port, in string form, of the database server.
func NewDb(addr string) (db *Db, err os.Error) {
	s, err := mgo.Mongo(addr)
	if err != nil {
		log.Printf("Problem connecting to MongoDB: %s", err)
		return nil, err
	}
	return &Db{ 
		s: s,
		u: s.DB("shelflife").C("users"),
		f: s.DB("shelflife").C("follow"),
	}, nil
}

func (db *Db) Close() os.Error {
	db.s.Close()
	return nil
}
