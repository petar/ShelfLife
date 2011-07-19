// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package db

import (
	"log"
	"os"
	"github.com/petar/ShelfLife/thirdparty/mgo"
)

// Db encapsulates all interaction with the database
type Db struct {
	s  *mgo.Session
	kp *KPartite
}

// NewDb creates a new Db interface to the database.
// addr is the IP address and port, in string form, of the database server.
func NewDb(addr, dbname string) (db *Db, err os.Error) {
	s, err := mgo.Mongo(addr)
	if err != nil {
		log.Printf("Problem connecting to MongoDB: %s", err)
		return nil, err
	}
	db = &Db{ 
		s:  s,
		kp: NewKPartite(dbname, s),
	}
	// Initialize user system
	if err := db.initUser(); err != nil {
		db.Close()
		return nil, err
	}
	// Initialize like system
	if err := db.initLike(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (db *Db) Close() os.Error {
	db.s.Close()
	return nil
}
